package api

import (
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// deleteMessage handles requests to delete a message in a conversation.
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	senderId := ctx.UserId
	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid message id in path, or failed parse")
		return
	}
	conversationId, err := strconv.Atoi(ps.ByName("conversationId"))
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
	if err = rt.db.DeleteAllCommentsByMessageId(messageId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to delete comments related to the message")
		return
	}
	participants, err := rt.db.GetParticipants(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get participants of the conversation")
		return
	}
	for _, userId := range participants {
		if userId != senderId {
			if err = rt.db.RemoveStatus(messageId, userId); err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add message status for a participant")
				return
			}
		}
	}
	if err = rt.db.UpdateLatestMessage(conversationId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update latest message in the database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Infof("message %d deleted successfully by user %d", messageId, senderId)
}
