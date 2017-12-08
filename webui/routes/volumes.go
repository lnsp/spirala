package routes

import "net/http"

func (router *Router) showVolumes(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Volumes")
	router.templates["volumes.html"].Execute(w, context)
}
