package routes

import "net/http"

func (router *Router) showContainers(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Containers")
	router.templates["containers.html"].Execute(w, context)
}
