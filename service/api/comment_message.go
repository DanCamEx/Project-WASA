package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// commentMessage handles requests to add a comment to a message.
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	commenterId := ctx.UserId
	// Decode comment text from request body
	var dec = json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req CommentMessageRequest

	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode comment message request")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Parse message id from path to int
	messageId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': could not convert message id in path to int")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Could not convert message id in path to int."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if user is participant of the conversation containing the message
	isParticipant, err := rt.db.CheckParticipantByMessageId(messageId, commenterId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if user is participant of the conversation containing the message")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if user is participant of the conversation containing the message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !isParticipant {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Warnf("403 'Forbidden': user %d is not a participant of the conversation containing message %d", commenterId, messageId)

		errRes := ErrorResponse{Code: "Forbidden", Message: "User is not a participant of the conversation containing the message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if the comment is one of the selected emoji
	if req.Emoji != "👍" && req.Emoji != "😂" && req.Emoji != "💀" && req.Emoji != "😭" && req.Emoji != "🔥" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': invalid emoji for comment")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Invalid emoji for comment."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if message exists
	messageExists, err := rt.db.CheckMessageOnly(messageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if message exists")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if message exists."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !messageExists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': message %d does not exist", messageId)

		errRes := ErrorResponse{Code: "Not Found", Message: "Message does not exist."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Add comment to the database
	emojiAdded, err := rt.db.CreateComment(messageId, commenterId, req.Emoji)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to add comment to message")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to add comment to message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Send the comment created as json response
	res := CommentMessageResponse{
		Emoji:     emojiAdded,
	}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to encode response."}
		EncodeErrResponse(w, errRes, ctx)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return
	}
	
}

type CommentMessageRequest struct {
	Emoji string `json:"emoji"`
}	

type CommentMessageResponse struct {
	Emoji     string `json:"emoji"`
}