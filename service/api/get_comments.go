package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// getComments handles requests to retrieve comments for a specific message.
func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	messageId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid message id in path, or failed parse")
		return
	}
	comments, err := rt.db.GetAllCommentsByMessageId(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get comments from database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode comments response")
		return
	}
	ctx.Logger.Info("message comments got successfully")
}
