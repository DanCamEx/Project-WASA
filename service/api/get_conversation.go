package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {

	// Parse conversation id from path to int
	conversationId, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': could not convert conversation id in path to int")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Could not convert conversation id in path to int."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Retrieve messages from the database
	messages, err := rt.db.GetConversationMessages(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Warn("500 'Internal Server Error': could not retrieve messages from the database")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Could not retrieve messages from the database."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Encode and send the messages as JSON response
	var res []MessagesResponse
	for _, msg := range messages {
		id := msg[0].(int)
		senderId := msg[1].(int)
		content := msg[2].(string)
		timestamp := msg[3].(string)
		status := msg[4].(string)
		isForwarded := msg[5].(bool)

		res = append(res, MessagesResponse{
			Id :         id,
			SenderId:    senderId,
			Content:     content,
			Timestamp:   timestamp,
			Status:      status,
			IsForwarded: isForwarded,
		})
	}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		ctx.Logger.Info("200 'OK': conversation retrieval successful")
		return
	}
}

type MessagesResponse struct {
	Id 		    int    `json:"id"`
	SenderId    int    `json:"senderId"`
	Content     string `json:"content"`
	Timestamp   string `json:"timestamp"`
	Status      string `json:"status"`
	IsForwarded bool   `json:"isForwarded"`
}