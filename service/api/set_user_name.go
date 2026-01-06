package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName handles requests to set or update the authenticated user's name.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Gets user id from context and decodes request body
	userId := ctx.UserId
	dec := NewStrictDecoder(r)
	var req struct {
		Name string `json:"name"`
	}
	err := dec.Decode(&req)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode updateUserName request")
		return
	}

	// Checks if user name is valid
	if regexp.MustCompile(`^[A-Za-z0-9 _.-]{3,16}$`).MatchString(req.Name) == false {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user name can be min 3, max 16 and only contain letters, numbers, spaces, and special characters (_ . -)")
		return
	}

	// Check if another user in the db with the same name already exists
	exists, err := rt.db.CheckUserName(req.Name)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check user name in db")
		return
	}
	if exists {
		SendErrorResponse(w, ctx, http.StatusConflict, "another user with the same name already exists")
		return
	}

	// Sets the new user name
	if err = rt.db.SetUserName(userId, req.Name); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to set new user name")
		return
	}

	// Get the user with the updated name and sends it as response
	user, err := rt.db.GetUser(userId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get updated user from database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode updated user response")
		return
	}
	ctx.Logger.Infof("user name updated for user id: %d", userId)
}
