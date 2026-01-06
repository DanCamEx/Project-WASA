/*
Package api exposes the main API engine. All HTTP APIs are handled here - so-called "business logic" should be here, or
in a dedicated package (if that logic is complex enough).

To use this package, you should create a new instance with New() passing a valid Config. The resulting Router will have
the Router.Handler() function that returns a handler that can be used in a http.Server (or in other middlewares).

Example:

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger,
		Database: appdb,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	router := apirouter.Handler()

	// ... other stuff here, like middleware chaining, etc.

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           router,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()

See the `main.go` file inside the `cmd/webapi` for a full usage example.
*/
package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/DanCamEx/Project-WASA/service/database"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Config is used to provide dependencies and configuration to the New function.
type Config struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger

	// Database is the instance of database.AppDatabase where data are saved
	Database database.AppDatabase
}

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler() http.Handler

	// Close terminates any resource used in the package
	Close() error
}

// New returns a new Router instance
func New(cfg Config) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}, nil
}

type _router struct {
	router *httprouter.Router

	// baseLogger is a logger for non-requests contexts, like goroutines or background tasks not started by a request.
	// Use context logger if available (e.g., in requests) instead of this logger.
	baseLogger logrus.FieldLogger

	db database.AppDatabase
}

// SendErrorResponse is a helper function to send a JSON error response with the given status code and message.
func SendErrorResponse(w http.ResponseWriter, ctx reqcontext.RequestContext, statusCode int, errorMessage string) {
	response := ErrorResponse{
		Status:  http.StatusText(statusCode),
		Code:    statusCode,
		Message: errorMessage,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to encode error response")
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
		return
	}

	ctx.Logger.WithFields(logrus.Fields{
		"status": response.Status,
		"code":   response.Code,
	}).Error(response.Message)
}

// NewStrictDecoder creates a simple JSON decoder that disallows unknown fields.
func NewStrictDecoder(r *http.Request) *json.Decoder {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec
}

// ErrorResponse represents the standard error response structure for the API.
// It's used in requests when returning an error to the client.
type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	direct = "direct"
)
