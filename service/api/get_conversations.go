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
	// Parse string ids to int
	var res []ConversationRes
	for _, conv := range convs {
		convId := conv[0]
		convName := conv[1]
		convImage := conv[2]
		convType := conv[3]
		latestMessageId := conv[4]

		convIdInt, err := strconv.Atoi(convId)
		if err != nil {
			ctx.Logger.Error("500 'Internal Server Error': failed to convert conversation id to int")
			w.WriteHeader(http.StatusInternalServerError)

			errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to convert conversation id to int."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}

		latestMessageIdInt, err := strconv.Atoi(latestMessageId)
		if err != nil {
			ctx.Logger.Error("500 'Internal Server Error': failed to convert conversation id to int")
			w.WriteHeader(http.StatusInternalServerError)

			errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to convert latest message id to int."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}

		res = append(res, ConversationRes{
			Id:   convIdInt,
			Name: convName,
			Image: convImage,
			Type: convType,
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

type ConversationRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Type string `json:"type"`
	LatestMessageId int `json:"latest_message_id"`	
}