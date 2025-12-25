package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// SearchUsers handles requests to search for all users in the system.
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	users, err := rt.db.GetAllUsers()
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get users from database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode users response")
		return
	}
	ctx.Logger.Info("users search successful")
}
