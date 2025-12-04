package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.doLogin)
	rt.router.PATCH("/users/me", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users", rt.wrap(rt.searchUsers))
	rt.router.POST("/users/me/conversations", rt.wrap(rt.createConversation))
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))
	rt.router.PATCH("/conversations/:id", rt.wrap(rt.setGroupName))

	return rt.router
}