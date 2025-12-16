package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// addToGroup handles requests to add a user to a group conversation.
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	dec := NewStrictDecoder(r)
	var addRequest struct {
		UserId int `json:"user_id"`
	}
	if err := dec.Decode(&addRequest); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode addToGroup request")
		return
	}

	groupId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid group id")
		return
	}

	exists, err := rt.db.CheckIfGroup(groupId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if conversation is a group")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "conversation is not a group")
		return
	}
	exists, err = rt.db.CheckParticipant(groupId, addRequest.UserId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if user is already a participant")
		return
	}
	if exists {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user is already a participant in the group conversation")
		return
	}

	if err := rt.db.AddParticipant(groupId, addRequest.UserId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add user to group conversation")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("User added to group conversation successfully")
}
