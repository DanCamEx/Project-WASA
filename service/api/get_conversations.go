package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyConversations handles requests to retrieve all conversations for the authenticated user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	conversations, err := rt.db.GetAllUserConversations(userId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get user conversations")
		return
	}
	for i := range *conversations {
		if (*conversations)[i].Category == "direct" {
			otherUserId, err := rt.db.GetOtherContact(userId, (*conversations)[i].Id)
			if err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other contact in direct conversation")
				return
			}
			otherUser, err := rt.db.GetUser(otherUserId)
			if err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other user info in direct conversation")
				return
			}
			(*conversations)[i].Name = otherUser.Name
			(*conversations)[i].Photo = otherUser.Photo
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(conversations); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode the conversations response")
		return
	}
	ctx.Logger.Info("conversations got successfully")
}
