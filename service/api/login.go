package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// doLogin handles user login requests.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Decode login request
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: failed to decode login request")
		return
	}

	// Get or create user ID
	userId, err := rt.db.GetIDByUserName(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to get or create user ID")
		return
	}

	// Set headers and respond with user ID
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	logrus.Infof("user logged in with ID: %d", userId)
	response := IdResponse{Id: userId}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to encode login response")
		return
	}
}

type LoginRequest struct {
	Name string `json:"name"`
}

type IdResponse struct {
	Id int `json:"id"`
}
