package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// doLogin handles user login requests.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	emptyCtx := RequestContext{}

	// Decode login request
	var dec = json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req LoginRequest
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("400 'Bad Request': failed to decode")

		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode."}
		EncodeErrResponse(w, errRes, emptyCtx)
		return
	}

	// Get or create user id by username
	userId, err := rt.db.GetIdByUserName(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("500 'Internal Server Error': failed to get or create user id")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get or create user id."}
		EncodeErrResponse(w, errRes, emptyCtx)
		return
	}

	// Encode and respond with user id as JSON
	res := IdResponse{Id: userId}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("500 'Internal Server Error': failed to encode response")
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(buf.Bytes())
		logrus.Infof("201 'Created': user logged in with id: %d", userId)
	}
}

type LoginRequest struct {
	Name string `json:"name"`
}

type IdResponse struct {
	Id int `json:"id"`
}
