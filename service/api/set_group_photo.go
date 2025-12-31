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

// setMyPhoto handles requests to set or update the authenticated user's profile photo.
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationId, err := strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to get conversation id from path or invalid format")
		return
	}
	exists, err := rt.db.CheckConversation(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if the conversation exists in the database")
		return
	}
	if !exists {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "conversation does not exist")
		return
	}
	isGroup, err := rt.db.CheckIfGroup(conversationId)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to check if the conversation is a group in the database")
		return
	}
	if !isGroup {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "conversation is not a group")
		return
	}
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to parse multipart form")
		return
	}
	file, handler, err := r.FormFile("photo")
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

	err = rt.db.SetGroupPhoto(conversationId, filename)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to set/update group photo in db")
		return
	}

	var fileNameRes struct {
		FileName string `json:"fileName"`
	}
	fileNameRes.FileName = filename
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(fileNameRes); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode updated user response")
		return
	}
	ctx.Logger.Infof("group conversation %d photo updated successfully", conversationId)
}
