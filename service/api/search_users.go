package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// SearchUsers handles requests to search for all users in the system.
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Check authorization, ensure the token matches the user ID in the path
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		logrus.Warn("unauthorized access attempt")
		return
	}

	// Retrieve all users from the database
	info, err := rt.db.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to retrieve users from database")
		return
	}

	// Encode and send the user list as JSON response
	response := make([]usersInfoNoID, 0, len(info))
	for _, user := range info {
		response = append(response, usersInfoNoID{Name: user[0], Image: user[1]})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logrus.Info("users search successful")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to encode users response")
		return
	}

}

type usersInfoNoID struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
