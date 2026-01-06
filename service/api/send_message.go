package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// sendMessage handles requests to send a message in a conversation.
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	senderId := ctx.UserId
	dec := NewStrictDecoder(r)
	var req struct {
		Content string `json:"content"`
	}
	err := dec.Decode(&req)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode sendMessage request")
		return
	}

	conversationId, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to get conversation id from path or invalid format")
		return
	}
	exists, err := rt.db.CheckConversation(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if conversation exists")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "conversation doesn't exist")
		return
	}
	exists, err = rt.db.CheckParticipant(conversationId, senderId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if sender is a participant of the conversation")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "sender is not a participant of the conversation")
		return
	}
	messageId, err := rt.db.CreateMessage(senderId, conversationId, req.Content)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create the message in the database")
		return
	}
	participants, err := rt.db.GetParticipants(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get participants of the conversation")
		return
	}
	for _, participantId := range participants {
		if participantId != senderId {
			if err = rt.db.AddStatus(messageId, participantId); err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add message status for a participant")
				return
			}
		}
	}
	if err = rt.db.UpdateLatestMessage(conversationId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update latest message in the database")
		return
	}

	message, err := rt.db.GetMessage(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to retrieve the created message from the database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(message); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode response")
		return
	}
	ctx.Logger.Infof("message created with id: %d sent by user %d in conversation %d", messageId, senderId, conversationId)
}
