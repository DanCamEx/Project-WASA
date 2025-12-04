package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// SearchUsers handles requests to search for all users in the system.
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx RequestContext) {

	// Retrieve all users from the database
	info, err := rt.db.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to retrieve users from database")

		errRes := ErrorResponse{Code: "Internal Server Error", Message: "Failed to retrieve users from database."}
		EncodeErrResponse(w, errRes)
		return
	}

	// Encode and send the user list as JSON response
	res := make([]usersInfoNoID, 0, len(info))
	for _, user := range info {
		res = append(res, usersInfoNoID{Name: user[0], Image: user[1]})
	}
	w.Header().Set("content-type", "application/json")
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("500 'Internal Server Error': failed to encode response")
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		ctx.Logger.Info("200 'OK': users search successful")
		return
	}

}

type usersInfoNoID struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
