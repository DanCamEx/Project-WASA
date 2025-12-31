package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getUserByName handles requests to get a user's information.
func (rt *_router) getUserByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userName := ps.ByName("name")
	var err error
	if userName == "" {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "missing user identifier in path")
		return
	}

	user, err := rt.db.GetUserByName(userName)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get user from database")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode user response")
		return
	}
	ctx.Logger.Info("user got successfully")
}
