package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Skip Authorization for login route
		if r.URL.Path != "/session" && !strings.HasPrefix(r.URL.Path, "/images/") {
			// Check Authorization header
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				SendErrorResponse(w, ctx, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}
			strId := strings.TrimPrefix(auth, "Bearer ")
			id, err := strconv.Atoi(strId)
			if err != nil {
				SendErrorResponse(w, ctx, http.StatusUnauthorized, "invalid user id in Authorization header")
				return
			}
			exists, err := rt.db.CheckUserId(id)
			if err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check user id in database")
				return
			}
			if !exists {
				SendErrorResponse(w, ctx, http.StatusUnauthorized, "user id in Authorization header does not exist")
				return
			}
			ctx.UserId = id
			ctx.Logger = ctx.Logger.WithField("user-id", ctx.UserId)
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
