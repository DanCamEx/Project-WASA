package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	// Decode new group name from request body
	var req NewGroupName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: failed to decode new group name")
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

	// Update group name in the database
	err = rt.db.SetGroupName(conversationId, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to set new group name")
		return
	}
	//TODO RETURN CONVERSATION INFO
	w.WriteHeader(http.StatusOK)
	logrus.Info("group name update successful")
}

type NewGroupName struct {
	Name string `json:"name"`
}
