package api

import (
	"encoding/json"
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/DanCamEx/Project-WASA/service/database"
	"github.com/julienschmidt/httprouter"
)

// createConversation handles requests to create a new direct conversation between two users.
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Decodes request body
	creatorId := ctx.UserId
	dec := NewStrictDecoder(r)
	var req struct {
		Name            string `json:"name"`
		ParticipantsIds []int  `json:"participantsIds"`
		Category        string `json:"category"`
	}
	if err := dec.Decode(&req); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to decode createConversation request")
		return
	}

	// Loop through participant IDs to check if they exist or one of them is the creator.
	for _, participantId := range req.ParticipantsIds {
		if participantId == creatorId {
			SendErrorResponse(w, ctx, http.StatusBadRequest, "creator cannot be a participant in the conversation")
			return
		}
		exists, err := rt.db.CheckUserId(participantId)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if participant user id exists")
			return
		}
		if !exists {
			SendErrorResponse(w, ctx, http.StatusNotFound, "participant user id does not exist")
			return
		}
	}

	// Creates res struct and check conversation category to create the right conversation. If the category is direct,
	// it also checks that there is exactly one participant and set the res conversation name and photo to the other user's info,
	// depending on who made the request.
	var res struct {
		Id            int                    `json:"id"`
		Name          string                 `json:"name"`
		Photo         string                 `json:"photo"`
		Category      string                 `json:"category"`
		LatestMessage database.LatestMessage `json:"latestMessage"`
	}
	var err error
	if req.Category == "group" {
		res.Category = "group"
		res.Id, err = rt.db.CreateGrpConversation(creatorId, req.ParticipantsIds, req.Name)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create group conversation")
			return
		}
	}
	if req.Category == direct {
		res.Category = direct
		if len(req.ParticipantsIds) != 1 {
			SendErrorResponse(w, ctx, http.StatusBadRequest, "direct conversation must have exactly one participant")
			return
		}
		res.Id, err = rt.db.CreateDirConversation(creatorId, req.ParticipantsIds[0])
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create conversation")
			return
		}
		if res.Id == 0 {
			SendErrorResponse(w, ctx, http.StatusConflict, "direct conversation between these users already exists")
			return
		}
		otherUserId, err := rt.db.GetOtherContact(creatorId, res.Id)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other contact in direct conversation")
			return
		}
		otherUser, err := rt.db.GetUser(otherUserId)
		if err != nil {
			SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed db operation to get other user info in direct conversation")
			return
		}
		res.Name = otherUser.Name
		res.Photo = otherUser.Photo
	}

	// Sends response with created conversation data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode create conversation response")
		return
	}
	ctx.Logger.Infof("conversation created with id: %d", res.Id)
}
