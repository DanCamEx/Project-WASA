package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var ctx = RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Checks the Authorization header for a valid Bearer token.
		// If valid extracts the user id and sets it in the context
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
		ctx.UserId = id

		// Update the logger with the user id if the auth is successful
		ctx.Logger = ctx.Logger.WithField("user-id", ctx.UserId)

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

// RequestContext is the context of the request, for request-dependent parameters
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID uuid.UUID

	// UserId is the authenticated user id for the request, or 0 if not authenticated
	UserId int

	// Logger is a custom field logger for the request
	Logger logrus.FieldLogger
}
