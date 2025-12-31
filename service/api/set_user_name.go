package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName handles requests to set or update the authenticated user's name.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	dec := NewStrictDecoder(r)
	var userName struct {
		Name string `json:"name"`
	}
	err := dec.Decode(&userName)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode updateUserName request")
		return
	}
	if userName.Name == "" {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user name cannot be empty")
		return
	}
	if len(userName.Name) > 16 {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "user name cannot be longer than 16 characters")
		return
	}

	exists, err := rt.db.CheckUserName(userName.Name)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check user name in db")
		return
	}
	if exists {
		SendErrorResponse(w, ctx, http.StatusConflict, "another user with the same name already exists")
		return
	}

	if err = rt.db.SetUserName(userId, userName.Name); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to set new user name")
		return
	}
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
