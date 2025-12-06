package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	senderId := ctx.UserId

	// Parse message id from path to int
	messageId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': could not convert message id in path to int")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Could not convert message id in path to int."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if message exists and belongs to sender
	msgExists, err := rt.db.CheckMessage(messageId, senderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if message exists or belongs to sender")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if message exists or belongs to sender."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !msgExists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': message %d does not exist or does not belong to user %d", messageId, senderId)

		errRes := ErrorResponse{Code: "Not Found", Message: "Message does not exist or does not belong to user."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Delete message from the database
	err = rt.db.DeleteMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to delete message")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to delete message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("204 'No Content': message %d deleted successfully by user %d", messageId, senderId)
}