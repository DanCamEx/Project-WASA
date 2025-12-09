package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getMyConversations handles requests to retrieve all conversations for the authenticated user.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	userId := ctx.UserId

	// Retrieve all conversations for the user from the database
	convs, err := rt.db.GetAllUserConversations(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to retrieve conversations from database")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to retrieve conversations from database."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Loop through conversations and convert to response format
	// Parse latestMessageId to int
	var res []ConversationResponse
	for _, conv := range convs {
		id := conv[0].(int)
		name := conv[1].(string)
		image := conv[2].(string)
		category := conv[3].(string)
		latestMessageId := conv[4].(int)


		res = append(res, ConversationResponse{
			Id:              id,
			Name:            name,
			Image:           image,
			Category:        category,
			LatestMessageId: latestMessageId,
		})
	}

	// Encode and send the conversation list as JSON response
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to encode response."}
		EncodeErrResponse(w, errRes, ctx)
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write(buf.Bytes())

		ctx.Logger.Info("200 'OK': conversations retrieval successful")
		w.WriteHeader(http.StatusOK)
		return
	}
}

type ConversationResponse struct {
	Id 				int    `json:"id"`
	Name            string `json:"name"`
	Image           string `json:"image"`
	Category        string `json:"category"`
	LatestMessageId int    `json:"latestMessageId"`	
}