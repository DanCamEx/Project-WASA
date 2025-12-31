package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMessage handles requests to get a message's information.
func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid message id in path, or failed parse")
		return
	}
	message, err := rt.db.GetMessage(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to retrieve message")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode message response")
		return
	}
	ctx.Logger.Info("message got successfully")

}
