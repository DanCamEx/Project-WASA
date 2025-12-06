package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// login
	rt.router.POST("/session", rt.doLogin)

	// user
	rt.router.PATCH("/users/me", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users", rt.wrap(rt.searchUsers))

	// conversation
	rt.router.POST("/users/me/conversations", rt.wrap(rt.createConversation))
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))
	rt.router.PATCH("/conversations/:id", rt.wrap(rt.setGroupName))
	rt.router.GET("/users/me/conversations", rt.wrap(rt.getMyConversations))

	// message
	rt.router.DELETE("/messages/:id", rt.wrap(rt.deleteMessage))

	return rt.router
}