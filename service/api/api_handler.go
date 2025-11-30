package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.doLogin)
	rt.router.PATCH("/users/:userId", rt.SetMyUserName)

	return rt.router
}
