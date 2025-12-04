package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {

	

	// Decode new group name from request body
	var req NewGroupName
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Warn("bad request: failed to decode new group name")
		errRes := ErrorResponse{Code: "Bad Request", Message: "Failed to decode new group name."}
		EncodeErrResponse(w, errRes)
		return
	}

	// Parse conversation id from path to int
	var conversationId int
	conversationId, err = strconv.Atoi(ps.ByName("conversationId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Warn("bad request: could not convert conversation id in path to int")
		return
	}

	// Update group name in the database
	err = rt.db.SetGroupName(conversationId, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.Error("internal server error: failed to set new group name")
		return
	}
	//TODO RETURN CONVERSATION INFO
	w.WriteHeader(http.StatusOK)
	logrus.Info("group name update successful")
}

type NewGroupName struct {
	Name string `json:"name"`
}
