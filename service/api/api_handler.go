package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login route
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User routes
	rt.router.GET("/users/search", rt.wrap(rt.searchUsers))
	rt.router.PATCH("/users/me", rt.wrap(rt.setMyUserName))
	rt.router.POST("/users/me/conversations", rt.wrap(rt.createConversation))
	rt.router.GET("/users/me/conversations", rt.wrap(rt.getMyConversations))
	rt.router.PUT("/users/me/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users/by_name/:name", rt.wrap(rt.getUserByName))
	rt.router.GET("/users/by_id/:id", rt.wrap(rt.getUserById))

	// Conversation routes
	rt.router.PATCH("/conversations/:conversationId", rt.wrap(rt.setGroupName))
	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.sendMessage))
	rt.router.GET("/conversations/:conversationId/messages", rt.wrap(rt.getMessages))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation))
	rt.router.POST("/conversations/:conversationId/participants", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/conversations/:conversationId/participants", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/conversations/:conversationId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.POST("/conversations/:conversationId/images", rt.wrap(rt.uploadImage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.deleteMessage))

	// Message routes
	rt.router.POST("/messages/:messageId/forward", rt.wrap(rt.forwardMessage))
	rt.router.POST("/messages/:messageId/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:messageId/comments", rt.wrap(rt.uncommentMessage))
	rt.router.GET("/messages/:messageId", rt.wrap(rt.getMessage))
	rt.router.GET("/messages/:messageId/comments", rt.wrap(rt.getComments))

	// Other routes
	rt.router.GET("/images/:imageId", rt.wrap(rt.getImage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
