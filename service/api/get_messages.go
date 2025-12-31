package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMessages handles requests to retrieve messages from a specific conversation.
func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	conversationId, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid conversation id in path, or failed parse")
		return
	}

	messages, err := rt.db.GetConversationMessages(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get messages from database")
		return
	}
	for i := range *messages {
		if (*messages)[i].SenderId != userId && (*messages)[i].Status == "sent" {
			if err := rt.db.UpdateTracking(userId, (*messages)[i].Id); err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update message status")
				return
			}
		}
		if readByAll, err := rt.db.IsMessageReadByAll((*messages)[i].Id, conversationId); err == nil && readByAll {
			if err := rt.db.UpdateStatus((*messages)[i].Id); err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update message status")
				return
			}
		} else if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if message is read by all")
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode messages response")
		return
	}
	ctx.Logger.Info("conversation messages got successfully")
}
