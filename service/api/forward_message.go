package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	senderId := ctx.UserId

	// Decode conversation id where to forward the message
	var dec = json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req ForwardMessageRequest

	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode forward message request")

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

	// Check if message exists
	msgExists, err := rt.db.CheckMessageOnly(messageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if message exists")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if message exists."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !msgExists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': message %d does not exist", messageId)

		errRes := ErrorResponse{Code: "Not Found", Message: "Message does not exist."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	// Check if conversation exists
	conversationExists, err := rt.db.CheckConversation(req.ConversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if conversation exists")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if conversation exists."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !conversationExists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': conversation %d does not exist", req.ConversationId)

		errRes := ErrorResponse{Code: "Not Found", Message: "Conversation does not exist."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}


	// Forward the message using the database method and get new message id
	newMessageId, err := rt.db.ForwardMessage(messageId, senderId, req.ConversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to forward message")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to forward message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Encode new message and send response
	senderId, conversationId, content, timestamp, status, isForwarded, err:= rt.db.GetMessage(newMessageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to get forwarded message details")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get forwarded message details."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	res := ForwardMessageResponse{
		SenderId:		senderId,
		ConversationId: conversationId,
		Content:        content,
		Timestamp:      timestamp,
		Status:         status,
		IsForwarded:    isForwarded,
	}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode forward message response")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to encode forward message response."}
		EncodeErrResponse(w, errRes, ctx)
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		ctx.Logger.Infof("200 'OK': message %d forwarded successfully by user %d to conversation %d", messageId, senderId, req.ConversationId)
	}
}

type ForwardMessageRequest struct {
	ConversationId  int       `json:"conversationId"`
}

type ForwardMessageResponse struct {
	SenderId        int       `json:"senderId"`
	ConversationId  int       `json:"conversationId"`
	Content	        string    `json:"content"`
	Timestamp       string    `json:"timestamp"`
	Status          string    `json:"status"`
	IsForwarded     bool      `json:"isForwarded"`
}