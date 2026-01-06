package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// addToGroup handles requests to add a user to a group conversation.
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	dec := NewStrictDecoder(r)
	var req struct {
		UserName string `json:"userName"`
	}
	if err := dec.Decode(&req); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode addToGroup request")
		return
	}

	groupId, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid group id")
		return
	}
	user, err := rt.db.GetUserByName(req.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendErrorResponse(w, ctx, http.StatusBadRequest, "user not found")
			return
		}
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get user by name")
		return
	}
	userId := user.Id
	isGroup, err := rt.db.CheckIfGroup(groupId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if conversation is a group")
		return
	}
	if !isGroup {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "conversation is not a group")
		return
	}
	isParticipant, err := rt.db.CheckParticipant(groupId, userId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if user is already a participant")
		return
	}
	if isParticipant {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user is already a participant in the group conversation")
		return
	}

	if err := rt.db.AddParticipant(groupId, userId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add user to group conversation")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("User added to group conversation successfully")
}
