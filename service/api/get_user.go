package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getUser handles requests to get a user's information.
func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
}
