package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

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

	// Check if conversation exists and sender is a participant
	convExists, err := rt.db.CheckConversation(conversationId, senderId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if conversation exists and sender is a participant")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if conversation exists and sender is a participant."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !convExists {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Warnf("404 'Not Found': conversation %d does not exist or sender %d is not a participant", conversationId, senderId)
		
		errRes := ErrorResponse{Code: "Not Found", Message: "Conversation does not exist or sender is not a participant."}
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

	// Encode and respond with message info as JSON
	messageId, senderId, conversationId, content, timestamp, status, isForwarded, err := rt.db.GetMessage(messageId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to get message")
		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get message."}

		EncodeErrResponse(w, errRes, ctx)
		return
	}
	// Encode and respond with conversation id as JSON.
	res := MessageInfoResponse{
		Id:             messageId,
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

type MessageInfoResponse struct {
	Id             int    `json:"id"`
	SenderId       int    `json:"senderId"`
	ConversationId int    `json:"conversationId"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
	Status         string `json:"status"`
	IsForwarded    bool   `json:"isForwarded"`
}
