package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// deleteMessage handles requests to delete a message in a conversation.
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	senderId := ctx.UserId
	messageId, err := strconv.Atoi(ps.ByName("mid"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid message id in path, or failed parse")
		return
	}
	conversationId, err := strconv.Atoi(ps.ByName("cid"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid conversation id in path, or failed parse")
		return
	}

	exists, err := rt.db.CheckMessage(messageId, senderId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if message exists in database")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "message does not exist or does not belong to user")
		return
	}

	if err = rt.db.DeleteMessage(messageId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to delete message")
		return
	}
	if err = rt.db.UpdateLatestMessageId(conversationId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update latest message id in the database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("message %d deleted successfully by user %d", messageId, senderId)
}
