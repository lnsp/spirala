package routes

import "net/http"

func (router *Router) showNetworks(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Networks")
	router.templates["networks.html"].Execute(w, context)
}
