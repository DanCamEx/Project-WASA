package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	userId := ctx.UserId

	// Decode new user name from request body and check if it already exists
	var req NewUserName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("bad request: failed to decode new user name")
		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode new user name."}
		EncodeErrResponse(w, errRes)
		return
	}
	nameExists, err := rt.db.CheckUserName(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("internal server error: failed to check if user name exists")
		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if user name exists."}
		EncodeErrResponse(w, errRes)
		return
	}
	if nameExists {
		w.WriteHeader(http.StatusConflict)
		ctx.Logger.Warn("conflict: user name already exists")
		errRes := ErrorResponse{Code: "Conflict", Message: "User name already exists."}
		EncodeErrResponse(w, errRes)
		return
	}

	// Update user name in the database
	err = rt.db.SetUserName(userId, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("internal server error: failed to set new user name")
		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to set new user name."}
		EncodeErrResponse(w, errRes)
		return
	}

	// Retrieve updated user info
	userId, name, image, err := rt.db.GetUserInfoById(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("internal server error: failed to get user info")
		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get user info."}
		EncodeErrResponse(w, errRes)
		return
	}

	// Encode and respond with updated user info
	res := UserInfoResponse{Id: userId, Name: name, Image: image}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(res)
	w.Header().Set("content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write(buf.Bytes())
		ctx.Logger.Infof("201 'Created': user name updated for user id: %d", userId)
		return
	}
}

type NewUserName struct {
	Name string `json:"name"`
}

type UserInfoResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
