package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getUser handles requests to get a user's information.
func (rt *_router) getUserById(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid user id in path, or failed parse")
		return
	}
	user, err := rt.db.GetUser(userId)
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
