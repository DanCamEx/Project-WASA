package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getConversation handles requests to retrieve messages from a specific conversation.
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	conversationId, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid conversation id in path, or failed parse")
		return
	}

	conversation, err := rt.db.GetConversation(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get conversation from database")
		return
	}

	if conversation.Category == direct {
		otherUserId, err := rt.db.GetOtherContact(userId, (conversation).Id)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other contact in direct conversation")
			return
		}
		otherUser, err := rt.db.GetUser(otherUserId)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other user info in direct conversation")
			return
		}
		conversation.Name = otherUser.Name
		conversation.Photo = otherUser.Photo
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode conversation response")
		return
	}
	ctx.Logger.Info("conversation got successfully")
}
