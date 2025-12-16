package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getMyConversations handles requests to retrieve all conversations for the authenticated user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	userId := ctx.UserId
	conversations, err := rt.db.GetAllUserConversations(userId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get user conversations")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(conversations); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode the conversations response")
		return
	}
	ctx.Logger.Info("conversations got successfully")
}
