package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setGroupName handles requests to set or update the name of a group conversation.
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	dec := NewStrictDecoder(r)
	var groupName struct {
		Name string `json:"name"`
	}
	err := dec.Decode(&groupName)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode request")
		return
	}

	var conversationId int
	conversationId, err = strconv.Atoi(ps.ByName("id"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to get conversation id from path or invalid format")
		return
	}

	exists, err := rt.db.CheckIfGroup(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if the conversation it's a group in the database")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "conversation is not a group")
		return
	}

	if err = rt.db.SetGroupName(conversationId, groupName.Name); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to set new group name")
		return
	}
	conversation, err := rt.db.GetConversation(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get updated conversation from database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(conversation); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode updated conversation response")
		return
	}
	ctx.Logger.Info("group name updated successfully")
}
