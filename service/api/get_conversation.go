package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid conversation id in path, or failed parse")
		return
	}

	messages, err := rt.db.GetConversationMessages(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get messages from database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode messages response")
		return
	}
	ctx.Logger.Info("conversation messages got successfully")
}
