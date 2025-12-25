package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// doLogin handles user login requests.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := NewStrictDecoder(r)
	var doLoginRequest struct {
		Name string `json:"name"`
	}
	if err := dec.Decode(&doLoginRequest); err != nil {
		SendErrorResponse(w, reqcontext.RequestContext{}, http.StatusBadRequest, "failed to decode login request.")
		return
	}

	userId, err := rt.db.GetIdByUserName(doLoginRequest.Name)
	if err != nil {
		SendErrorResponse(w, reqcontext.RequestContext{}, http.StatusInternalServerError, "failed to get or create user id by username.")
		return
	}

	var doLoginResponse struct {
		Id int `json:"id"`
	}
	doLoginResponse.Id = userId
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(doLoginResponse); err != nil {
		http.Error(w, "failed to encode login response", http.StatusInternalServerError)
		return
	}

	logrus.Infof("user logged in with id: %d", userId)
}
