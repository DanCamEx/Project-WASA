package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// setMyUserName handles requests to set or update the authenticated user's name.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {
	userId := ctx.UserId

	// Decode new user name from request body and check if it already exists
	var req UserNameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("400 'Bad Request': failed to decode new user name")
		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode new user name."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	nameExists, err := rt.db.CheckUserName(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to check if user name exists")
		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to check if user name exists."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}
	if nameExists {
		w.WriteHeader(http.StatusConflict)
		ctx.Logger.Warn("409 'Conflict': user name already exists")

		errRes := ErrorResponse{Code: "Conflict", Message: "User name already exists."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Update user name in the database
	err = rt.db.SetUserName(userId, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to set new user name")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to set new user name."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Retrieve updated user info
	name, image, err := rt.db.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to get user info")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to get user info."}
		EncodeErrResponse(w, errRes, ctx)
		return
	}

	// Encode and respond with updated user info
	res := UserResponse{Name: name, Image: image}
	var buf bytes.Buffer
	
	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to encode response."}
		EncodeErrResponse(w, errRes, ctx)
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(buf.Bytes())
		ctx.Logger.Infof("201 'Created': user name updated for user id: %d", userId)
		return
	}
}

type UserNameRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
