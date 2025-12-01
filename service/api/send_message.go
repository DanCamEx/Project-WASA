package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Check authorization
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		logrus.Warn("unauthorized access attempt")
		return
	}

	// Parse sender id from token to int
	senderId, err := strconv.Atoi(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: could not convert user id in token to int")
		return
	}

	// Decode message content from request body
	var req SendMessageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: failed to decode message content")
		return
	}

	// Parse conversation id from path to int
	var conversationId int
	conversationId, err = strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: could not convert conversation id in path to int")
		return
	}

	// Check if conversation exists
	var exists bool
	exists, err = rt.db.CheckConversation(conversationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to check if conversation exists")
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		logrus.Warnf("not found: conversation %s does not exist for user %d", ps.ByName("conversationId"), senderId)
		return
	}

	// Create message in the database
	var messageId int
	var content, timestamp, status string
	var isForwarded bool
	messageId, err = rt.db.CreateMessage(senderId, conversationId, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to create message")
		return
	}

	// Encode and respond with message info as JSON
	messageId, senderId, conversationId, content, timestamp, status, isForwarded, err = rt.db.GetMessage(messageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to get message")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := MessageInfoResponse{
		Id:             messageId,
		SenderId:       senderId,
		ConversationId: conversationId,
		Content:        content,
		Timestamp:      timestamp,
		Status:         status,
		IsForwarded:    isForwarded,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to encode message info response")
		return
	}
	logrus.Infof("message successfully sent in conversation %d by user %d", conversationId, senderId)
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
