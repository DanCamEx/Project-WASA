package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// sendMessage handles requests to send a message in a conversation.
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	senderId := ctx.UserId

	// Decode message content from request body
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req SendMessageRequest

	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Parse conversation id from path to int
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': could not convert conversation id in path to int")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Could not convert conversation id in path to int."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if conversation exists
	exists, err := rt.db.CheckConversation(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if conversation exists")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if conversation exists."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': conversation %d does not exist", conversationId)
		
		errRes := ErrorResponse{Code: "Not Found", Message: "Conversation does not exist."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if sender is a participant of the conversation
	exists, err = rt.db.CheckParticipant(conversationId, senderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if sender is a participant of the conversation")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if sender is a participant of the conversation."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Warnf("403 'Forbidden': sender %d is not a participant of conversation %d", senderId, conversationId)

		errRes := ErrorResponse{Code: "Forbidden", Message: "Sender is not a participant of the conversation."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Create message in the database
	messageId, err := rt.db.CreateMessage(senderId, conversationId, req.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to create message")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to create message."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Update conversation's latestMessageId
	err = rt.db.UpdateLatestMessageId(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to update conversation's latestMessageId")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to update conversation's latestMessageId."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}


	// Encode and respond with message info as JSON
	id, senderId, conversationId, content, timestamp, status, isForwarded, err := rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to get message")
		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get message."}

		EncodeErrResponse(w, errRes, ctx)
		return
	}
	// Encode and respond with conversation id as JSON.
	res := MessageSentResponse{
		Id:             id,
		SenderId:       senderId,
		ConversationId: conversationId,
		Content:        content,
		Timestamp:      timestamp,
		Status:         status,
		IsForwarded:    isForwarded,
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to encode response."}
		EncodeErrResponse(w, errRes, ctx)
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(buf.Bytes())
		ctx.Logger.Infof("201 'Created': message created with id: %d sent by user %d in conversation %d", messageId, senderId, conversationId)
		return
	}
}

type SendMessageRequest struct {
	Content string `json:"content"`
}

type MessageSentResponse struct {
	Id 		       int    `json:"id"`
	SenderId       int    `json:"senderId"`
	ConversationId int    `json:"conversationId"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
	Status         string `json:"status"`
	IsForwarded    bool   `json:"isForwarded"`
}
