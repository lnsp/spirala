package routes

import "net/http"

func (router *Router) showImages(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Images")
	router.templates["images.html"].Execute(w, context)
}
