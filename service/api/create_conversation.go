package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// createConversation handles requests to create a new direct conversation between two users.
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	creatorId := ctx.UserId

	// Decode participants ids and conversation category from request body
	var dec = json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req CreateConversationRequest

	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode create conversation request")


		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode."}
		EncodeErrResponse(w, errRes, ctx)

		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode create conversation request")
		return
	}

	// Check participants ids are not the same as creator id
	for _, participantId := range req.ParticipantIds {
		if participantId == creatorId {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Warn("400 'Bad Request': cannot create conversation with oneself")

			errRes := ErrorResponse{Code: "Bad Request", Message: "Cannot create conversation with oneself."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}
	}

	// Check if participants exist
	for _, participantId := range req.ParticipantIds {
		idExists, err := rt.db.CheckUserId(participantId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("500 'Internal Server Error': failed to check if participant exists")

			errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if participant exists."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}
		if !idExists {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.Warnf("404 'Not Found': participant user id %d does not exist", participantId)

			errRes := ErrorResponse{Code: "Not Found", Message: "Participant user id does not exist."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}
	}

	// Check conversation category and create conversation in the database
	var conversationId int
	if req.Category != "direct" {
		conversationId, err = rt.db.CreateGrpConversation(creatorId, req.ParticipantIds)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("500 'Internal Server Error': failed to create group conversation")

			errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to create group conversation."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}
	} else {
		// Check that there is exactly one participant for direct conversations
		if len(req.ParticipantIds) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Warn("400 'Bad Request': direct conversations must have exactly one participant")

			errRes := ErrorResponse{Code: "Bad Request", Message: "Direct conversations must have exactly one participant."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}

		// Check if a direct conversation between the two users already exists.
		conversationId, err = rt.db.CreateDirConversation(creatorId, req.ParticipantIds[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.Error("500 'Internal Server Error': failed to create direct conversation")

			errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to create direct conversation."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}
		if conversationId == 0 {
			w.WriteHeader(http.StatusConflict)
			ctx.Logger.Warn("409 'Conflict': direct conversation between the two users already exists")

			errRes := ErrorResponse{Code: "Conflict", Message: "Direct conversation between the two users already exists."}
			EncodeErrResponse(w, errRes, ctx)
			return
		}
	}

	// Encode and respond with conversation id as JSON.
	res := ConversationIdResponse{Id: conversationId}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(buf.Bytes())
		ctx.Logger.Infof("201 'Created': conversation created with id: %d", conversationId)
		return
	}
}

type CreateConversationRequest struct {
	ParticipantIds []int  `json:"participantIds"`
	Category       string `json:"category"`
}

type ConversationIdResponse struct {
	Id int `json:"id"`
}
