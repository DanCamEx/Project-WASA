package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// login
	rt.router.POST("/session", rt.doLogin)

	// user
	rt.router.GET("/users", rt.wrap(rt.searchUsers))
	rt.router.PATCH("/users/me", rt.wrap(rt.setMyUserName))

	// conversation
	rt.router.POST("/users/me/conversations", rt.wrap(rt.createConversation))
	rt.router.GET("/users/me/conversations", rt.wrap(rt.getMyConversations))
	rt.router.PATCH("/conversations/:id", rt.wrap(rt.setGroupName))
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))
	rt.router.GET("/conversations/:id/messages", rt.wrap(rt.getConversation))

	// message
	rt.router.DELETE("/conversations/:cid/messages/:mid", rt.wrap(rt.deleteMessage))
	rt.router.POST("/messages/:id", rt.wrap(rt.forwardMessage))
	rt.router.POST("/messages/:id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:id/comments", rt.wrap(rt.uncommentMessage))

	return rt.router
}
