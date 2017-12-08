package routes

import "net/http"

func (router *Router) showServices(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Services")
	router.templates["services.html"].Execute(w, context)
}
