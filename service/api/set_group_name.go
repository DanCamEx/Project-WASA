package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// setGroupName handles requests to set or update the name of a group conversation.
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	// Decode new group name from request body
	var req NewGroupName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode new group name")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode new group name."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Parse conversation id from path to int
	var conversationId int
	conversationId, err = strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': could not convert conversation id in path to int")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Could not convert conversation id in path to int."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Check if conversation is a group
	isGroup, err := rt.db.CheckIfGroup(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if conversation is a group")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if conversation is a group."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if !isGroup {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': conversation is not a group")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Conversation is not a group."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Update group name in the database
	err = rt.db.SetGroupName(conversationId, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to set new group name")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to set new group name."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Get the updated conversation
	id, name, image, category, latestMessageId, err := rt.db.GetConversation(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to get updated conversation")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get updated conversation."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Encode and respond with updated conversation info
	res := UpdatedConversationResponse{
		Id:              id,
		Name:            name,
		Image:           image,
		Category:        category,
		LatestMessageId: latestMessageId,
	}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode updated conversation response")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to encode updated conversation response."}
		EncodeErrResponse(w, errRes, ctx)
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		ctx.Logger.Info("200 'OK': group name update successful")
		w.Write(buf.Bytes())
		return
	}
}

type NewGroupName struct {
	Name string `json:"name"`
}

type UpdatedConversationResponse struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Image           string `json:"image"`
	Category        string `json:"category"`
	LatestMessageId int    `json:"latestMessageId"`
}
