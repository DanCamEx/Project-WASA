package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// doLogin handles user login requests.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	// Decodes request body
	dec := NewStrictDecoder(r)
	var req struct {
		Name string `json:"name"`
	}
	if err := dec.Decode(&req); err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to decode login request")
		return
	}

	// Checks if user name is valid
	if regexp.MustCompile(`^[A-Za-z0-9 _.-]{3,16}$`).MatchString(req.Name) == false {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user name can be min 3, max 16 and only contain letters, numbers, spaces, and special characters (_ . -)")
		return
	}

	// Get user data from db or add the user if not exists
	user, err := rt.db.GetUserByName(req.Name)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get or create user id")
		return
	}

	// Sends response with all user data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode login response", http.StatusInternalServerError)
		return
	}

	logrus.Infof("user %s logged in with id: %d", user.Name, user.Id)
}
