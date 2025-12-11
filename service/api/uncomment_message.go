package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// uncommentMessage handles requests to remove a comment from a message.
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	commenterId := ctx.UserId

	// Parse message id from path to int
	var messageId int
	var err error
	messageId, err = strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': could not convert message id in path to int")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Could not convert message id in path to int."}
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

	// Check if user is participant of the conversation
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

	// Check if comment exists
	commentExists, err := rt.db.CheckComment(messageId, commenterId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if comment exists or belongs to commenter")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if comment exists or belongs to commenter."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !commentExists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': comment on message %d does not exist or does not belong to user %d", messageId, commenterId)

		errRes := ErrorResponse{Code: "Not Found", Message: "Comment does not exist or does not belong to user."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Remove comment from the database
	err = rt.db.DeleteComment(messageId, commenterId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to remove comment from message")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to remove comment from message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("204 'No Content': comment removed successfully from message %d by user %d", messageId, commenterId)
}
