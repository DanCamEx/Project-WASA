package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.doLogin)
	rt.router.PATCH("/users/:userId", rt.SetMyUserName)
	rt.router.GET("/users", rt.SearchUsers)
	rt.router.POST("/users/:userId/conversations", rt.CreateConversation)
	rt.router.POST("/conversations/:conversationId/messages", rt.sendMessage)

	return rt.router
}
