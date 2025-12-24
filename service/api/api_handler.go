package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// login
	rt.router.POST("/session", rt.doLogin)

	// user
	rt.router.GET("/users/search", rt.wrap(rt.searchUsers))
	rt.router.PATCH("/users/me", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users", rt.wrap(rt.getUser))

	// conversation
	rt.router.POST("/users/me/conversations", rt.wrap(rt.createConversation))
	rt.router.GET("/users/me/conversations", rt.wrap(rt.getMyConversations))
	rt.router.PATCH("/conversations/:id", rt.wrap(rt.setGroupName))
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))
	rt.router.GET("/conversations/:id/messages", rt.wrap(rt.getConversation))
	rt.router.POST("/conversations/:id/participants", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/conversations/:id/participants", rt.wrap(rt.leaveGroup))

	// message
	rt.router.DELETE("/conversations/:id/messages/:id", rt.wrap(rt.deleteMessage))
	rt.router.POST("/messages/:id", rt.wrap(rt.forwardMessage))
	rt.router.POST("/messages/:id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:id/comments", rt.wrap(rt.uncommentMessage))
	rt.router.GET("/messages/:id", rt.wrap(rt.getMessage))
	rt.router.GET("/messages/:id/comments", rt.wrap(rt.getComments))

	return rt.router
}
