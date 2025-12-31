package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	senderId := ctx.UserId
	conversationId, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "invalid conversation id")
		return
	}
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to parse multipart form")
		return
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "file non trovato")
		return
	}
	defer file.Close()

	uuidStr := uuid.Must(uuid.NewV4()).String()
	ext := ""
	for i := len(handler.Filename) - 1; i >= 0; i-- {
		if handler.Filename[i] == '.' {
			ext = handler.Filename[i:]
			break
		}
	}
	filename := uuidStr + ext
	path := "images/" + filename

	dst, err := os.Create(path)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create file")
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to copy uploaded file")
		return
	}

	messageId, err := rt.db.CreateMessage(senderId, conversationId, filename)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to set/update profile photo in db")
		return
	}

	participants, err := rt.db.GetParticipants(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to get participants of the conversation")
		return
	}
	for _, participantId := range participants {
		if participantId != senderId {
			if err = rt.db.AddStatus(messageId, participantId); err != nil {
				SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to add message status for a participant")
				return
			}
		}
	}
	if err = rt.db.UpdateLatestMessage(conversationId); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to update latest message in the database")
		return
	}

	message, err := rt.db.GetMessage(messageId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to retrieve the created message from the database")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(message); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode response")
		return
	}
	ctx.Logger.Infof("message created with id: %d sent by user %d in conversation %d", messageId, senderId, conversationId)
}
