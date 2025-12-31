package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	senderId := ctx.UserId
	dec := NewStrictDecoder(r)
	var forwardMessage struct {
		ConversationId int `json:"conversationId"`
	}
	if err := dec.Decode(&forwardMessage); err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to decode forward message request")
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "could not convert message id in path to int or invalid format")
		return
	}

	exists, err := rt.db.CheckMessageOnly(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if message exists")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "message does not exist")
		return
	}
	exists, err = rt.db.CheckConversation(forwardMessage.ConversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if conversation exists")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "conversation does not exist")
		return
	}
	newMessageId, err := rt.db.ForwardMessage(messageId, senderId, forwardMessage.ConversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to forward message")
		return
	}
	participants, err := rt.db.GetParticipants(forwardMessage.ConversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get participants of the conversation")
		return
	}
	for _, userId := range participants {
		if userId != senderId {
			if err = rt.db.AddStatus(messageId, userId); err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add message status for a participant")
				return
			}
		}
	}
	if err = rt.db.UpdateLatestMessage(forwardMessage.ConversationId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update latest message in the database")
		return
	}

	message, err := rt.db.GetMessage(newMessageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get forwarded message details")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(message); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode forward message response")
		return
	}
	ctx.Logger.Infof("message %d forwarded successfully by user %d to conversation %d", messageId, senderId, forwardMessage.ConversationId)
}
