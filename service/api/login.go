package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// doLogin handles user login requests.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	dec := NewStrictDecoder(r)
	var doLoginRequest struct {
		Name string `json:"name"`
	}
	if err := dec.Decode(&doLoginRequest); err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to decode login request.")
		return
	}

	if doLoginRequest.Name == "" {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user name cannot be empty.")
		return
	}
	if len(doLoginRequest.Name) > 16 {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user name cannot be longer than 16 characters.")
		return
	}

	user, err := rt.db.GetUserByName(doLoginRequest.Name)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get or create user id by username.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode login response", http.StatusInternalServerError)
		return
	}

	logrus.Infof("user %s logged in with id: %d", user.Name, user.Id)
}
