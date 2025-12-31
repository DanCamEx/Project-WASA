package api

import (
	"net/http"

	"github.com/DanCamEx/Project-WASA/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// get handles requests to retrieve an  by its filename
func (rt *_router) get(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	fileName := ps.ByName("imageId")
	Path := "images/" + fileName
	http.ServeFile(w, r, Path)
	ctx.Logger.Infof(" %s sent successfully", fileName)
}
