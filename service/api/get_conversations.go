package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// CreateConversation handles requests to create a new direct conversation between two users.
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
		convName := conv[0]
		convImage := conv[1]
		convType := conv[2]
		latestMessageId := conv[3]

		latestMessageIdInt, err := strconv.Atoi(latestMessageId)
		if err != nil {
			ctx.Logger.Error("500 'Internal Server Error': failed to convert conversation id to int")
			w.WriteHeader(http.StatusInternalServerError)

			errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to convert latest message id to int."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}

		res = append(res, ConversationResponse{
			Name:            convName,
			Image:           convImage,
			Type:            convType,
			LatestMessageId: latestMessageIdInt,
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
	Name            string `json:"name"`
	Image           string `json:"image"`
	Type            string `json:"type"`
	LatestMessageId int    `json:"latestMessageId"`	
}