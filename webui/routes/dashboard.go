package routes

import "net/http"

func (router *Router) showDashboard(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Dashboard")
	router.templates["index.html"].Execute(w, context)
}
