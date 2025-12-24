package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// createConversation handles requests to create a new direct conversation between two users.
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	creatorId := ctx.UserId
	dec := NewStrictDecoder(r)
	var conversationRequest struct {
		ParticipantIds []int  `json:"participantIds"`
		Category       string `json:"category"`
	}
	if err := dec.Decode(&conversationRequest); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode createConversation request")
		return
	}
	for _, participantId := range conversationRequest.ParticipantIds {
		if participantId == creatorId {
			SendErrorResponse(w, ctx, http.StatusBadRequest, "creator cannot be a participant in the conversation")
			return
		}
		exists, err := rt.db.CheckUserId(participantId)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if participant user id exists")
			return
		}
		if !exists {
			SendErrorResponse(w, ctx, http.StatusNotFound, "participant user id does not exist")
			return
		}
	}
	var conversationResponse struct {
		Id              int    `json:"id"`
		Name            string `json:"name"`
		Image           string `json:"image"`
		Category        string `json:"category"`
		LatestMessageId int    `json:"latestMessageId"`
	}
	var err error
	switch conversationRequest.Category {
	case "group":
		conversationResponse.Id, err = rt.db.CreateGrpConversation(creatorId, conversationRequest.ParticipantIds)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create group conversation")
			return
		}
	case "direct":
		if len(conversationRequest.ParticipantIds) != 1 {
			SendErrorResponse(w, ctx, http.StatusBadRequest, "direct conversation must have exactly one participant")
			return
		}
		conversationResponse.Id, err = rt.db.CreateDirConversation(creatorId, conversationRequest.ParticipantIds[0])
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create direct conversation")
			return
		}
		if conversationResponse.Id == 0 {
			SendErrorResponse(w, ctx, http.StatusConflict, "direct conversation between these users already exists")
			return
		}
		otherUserId, err := rt.db.GetOtherContact(creatorId, conversationResponse.Id)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other contact in direct conversation")
			return
		}
		otherUser, err := rt.db.GetUser(otherUserId)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other user info in direct conversation")
			return
		}
		conversationResponse.Name = otherUser.Name
		conversationResponse.Image = otherUser.Image
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(conversationResponse); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode create conversation response")
		return
	}
	ctx.Logger.Infof("conversation created with id: %d", conversationResponse.Id)
}
