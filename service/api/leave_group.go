package api

import (
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// leaveGroup handles requests to remove a user from a group conversation.
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	dec := NewStrictDecoder(r)
	var leaveRequest struct {
		UserId int `json:"user_id"`
	}
	if err := dec.Decode(&leaveRequest); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode leaveGroup request")
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
	exists, err = rt.db.CheckParticipant(groupId, leaveRequest.UserId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if user is a participant")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user is not a participant in the group conversation")
		return
	}

	if err := rt.db.RemoveParticipant(groupId, leaveRequest.UserId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to remove user from group conversation")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx.Logger.Info("User removed from group conversation successfully")
}
