package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// CreateConversation handles requests to create a new direct conversation between two users.
func (rt *_router) CreateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Check authorization, ensure the token matches the user ID in the path
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		logrus.Warn("unauthorized access attempt")
		return
	} else if strings.TrimPrefix(auth, "Bearer ") != ps.ByName("userId") {
		w.WriteHeader(http.StatusForbidden)
		logrus.Warn("forbidden access attempt")
		return
	}

	// Decode participant ID and conversation type from request body
	var req CreateConversationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: failed to decode participant id and conversation type from request body")
		return
	}

	// Parse user id and participant id from path to int
	var creatorId int
	creatorId, err = strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: could not convert user id in path to int")
		return
	}

	// Check participant id is not the same as creator id
	if req.ParticipantId == creatorId {
		http.Error(w, "cannot create conversation with oneself", http.StatusBadRequest)
		logrus.Warn("bad request: cannot create conversation with oneself")
		return
	}

	// Check if participant exists
	var exists bool
	exists, err = rt.db.CheckUserId(req.ParticipantId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to check if participant exists")
		return
	}
	if !exists {
		http.Error(w, "participant user id does not exist", http.StatusNotFound)
		logrus.Warnf("not found: participant user id %d does not exist", req.ParticipantId)
		return
	}

	// Check conversation type and create conversation in the database
	var conversationId int
	if req.Type != "direct" {
		http.Error(w, "unsupported conversation type", http.StatusBadRequest)
		logrus.Warn("TODO: unsupported conversation type")
		return
	} else {
		conversationId, err = rt.db.CreateDirConversation(creatorId, req.ParticipantId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logrus.Error("internal server error: failed to create direct conversation")
			return
		}
	}

	// Encode and respond with conversation id as JSON
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	logrus.Infof("direct conversation created with id: %d", conversationId)
	response := ConversationIdResponse{Id: conversationId}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to encode create conversation response")
		return
	}
}

type CreateConversationRequest struct {
	ParticipantId int    `json:"participantId"`
	Type          string `json:"type"`
}

type ConversationIdResponse struct {
	Id int `json:"id"`
}
