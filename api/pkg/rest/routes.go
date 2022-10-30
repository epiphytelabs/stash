package rest

import "github.com/ddollar/stdapi"

func (r *REST) Routes(router *stdapi.Router) {
	router.Route("GET", "/blobs", r.BlobList)
	router.Route("POST", "/blobs", r.BlobCreate)
	router.Route("HEAD", "/blobs/{hash}", r.BlobExists)
	router.Route("GET", "/blobs/{hash}", r.BlobGet)
	router.Route("DELETE", "/blobs/{hash}", r.BlobDelete)

	router.Route("GET", "/blobs/{hash}/labels", r.LabelList)
	router.Route("POST", "/blobs/{hash}/labels", r.LabelCreate)
	router.Route("GET", "/blobs/{hash}/labels/{key:.*}", r.LabelGet)
	router.Route("DELETE", "/blobs/{hash}/labels", r.LabelDelete)

	router.Route("POST", "/users", r.UserCreate)
}
