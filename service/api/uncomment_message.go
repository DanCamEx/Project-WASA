package api

import (
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// uncommentMessage handles requests to remove a comment from a message.
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commenterId := ctx.UserId
	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid message id in path, or failed parse")
		return
	}

	exists, err := rt.db.CheckMessageOnly(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if message exists")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "message does not exist")
		return
	}
	exists, err = rt.db.CheckParticipantByMessageId(messageId, commenterId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if user is participant of the conversation containing the message")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusForbidden, "user is not a participant of the conversation containing the message")
		return
	}
	exists, err = rt.db.CheckComment(messageId, commenterId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if comment exists or belongs to commenter")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "comment does not exist or does not belong to user")
		return
	}

	if err = rt.db.DeleteComment(messageId, commenterId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to remove comment from message")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("comment removed successfully from message %d by user %d", messageId, commenterId)
}
