package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// setMyPhoto handles requests to set or update the authenticated user's profile photo.
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parses multipart form data and create file
	userId := ctx.UserId
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "failed to parse multipart form")
		return
	}
	file, handler, err := r.FormFile("photo")
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusBadRequest, "photo not found")
		return
	}
	defer file.Close()

	// Creates id for file and saves it
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
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to create photo file")
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to copy uploaded photo file")
		return
	}

	// Updates user photo filename in db
	err = rt.db.SetUserPhoto(userId, filename)
	if err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to set/update profile photo in db")
		return
	}

	// Sends response with filename
	var fileNameRes struct {
		FileName string `json:"fileName"`
	}
	fileNameRes.FileName = filename
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(fileNameRes); err != nil {
		SendErrorResponse(w, ctx, http.StatusInternalServerError, "failed to encode filename response")
		return
	}
	ctx.Logger.Infof("user %d profile photo updated successfully", userId)
}
