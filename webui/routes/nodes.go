package routes

import "net/http"

func (router *Router) showNodes(w http.ResponseWriter, r *http.Request) {
	context := router.getBaseContext().SetActive("Nodes")
	router.templates["nodes.html"].Execute(w, context)
}
