package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func (rt *_router) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Check authorization, ensure the token matches the user id in the path
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

	// Decode new username from request body and check if it already exists
	var req NewUserName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: failed to decode new username")
		return
	}
	exists, err := rt.db.CheckUserName(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to check if username exists")
		return
	}
	if exists {
		http.Error(w, "username already exists", http.StatusConflict)
		logrus.Warn("conflict: username already exists")
		return
	}

	// Parse user id from path to int
	var id int
	id, err = strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: could not convert user id in path to int")
		return
	}

	// Update username in the database
	err = rt.db.SetUserName(id, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to set new username")
		return
	}

	// Retrieve updated user info
	userId, name, image, err := rt.db.GetUserInfoById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to get user info")
		return
	}

	// Encode and respond with updated user info
	w.Header().Set("content-type", "application/merge-patch+json")
	logrus.Infof("user %d name successfully updated", id)
	response := UserInfoResponse{Id: userId, Name: name, Image: image}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to encode user info response")
		return
	}
	w.WriteHeader(http.StatusOK)
	logrus.Infof("user %d info obtained and response sent", id)
}

type NewUserName struct {
	Name string `json:"name"`
}

type UserInfoResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
