package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// commentMessage handles requests to add a comment to a message.
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commenterId := ctx.UserId
	dec := NewStrictDecoder(r)
	var commentRequest struct {
		Emoji string `json:"emoji"`
	}
	if err := dec.Decode(&commentRequest); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode commentMessage request")
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid message id in path, or failed parse")
		return
	}
	exists, err := rt.db.CheckParticipantByMessageId(messageId, commenterId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if user is participant of the conversation containing the message")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusForbidden, "user is not a participant of the conversation containing the message")
		return
	}

	if commentRequest.Emoji != "👍" && commentRequest.Emoji != "😂" && commentRequest.Emoji != "💀" && commentRequest.Emoji != "😭" && commentRequest.Emoji != "🔥" {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid emoji for comment")
		return
	}

	exists, err = rt.db.CheckMessageOnly(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if message exists")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusNotFound, "message does not exist")
		return
	}

	var emojiResponse struct {
		SenderId  int    `json:"sender_id"`
		MessageId int    `json:"message_id"`
		Emoji     string `json:"emoji"`
	}
	emojiResponse.Emoji, err = rt.db.CreateComment(messageId, commenterId, commentRequest.Emoji)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add comment to message")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(emojiResponse)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode comment message response")
		return
	}
	ctx.Logger.Infof("200 'OK': comment added to message %d by user %d", messageId, commenterId)
}
