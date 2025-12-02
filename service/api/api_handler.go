package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.doLogin)
	rt.router.PATCH("/users/:userId", rt.setMyUserName)
	rt.router.GET("/users", rt.searchUsers)
	rt.router.POST("/users/:userId/conversations", rt.createConversation)
	rt.router.POST("/conversations/:conversationId/messages", rt.sendMessage)
	rt.router.PATCH("/conversations/:conversationId", rt.setGroupName)
	return rt.router
}
