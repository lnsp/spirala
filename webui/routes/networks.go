package routes

import "net/http"

type NetworkInstance struct {
	ID     string
	Name   string
	Driver string
}

type NetworkContext struct {
	BaseContext
	Networks []NetworkInstance
}

func (router *Router) showNetworks(w http.ResponseWriter, r *http.Request) {
	context, err := router.getNetworkContext()
	if err != nil {
		router.showError(w, r, err)
		return
	}
	router.templates["networks.html"].Execute(w, context)
}

func (router *Router) getNetworkContext() (NetworkContext, error) {
	context := NetworkContext{
		BaseContext: router.getBaseContext(),
		Networks:    []NetworkInstance{},
	}
	context.SetActive("Networks")
	for _, ep := range router.endpoints {
		networks, err := ep.ListNetworks()
		if err != nil {
			return NetworkContext{}, err
		}
		for _, net := range networks {
			context.Networks = append(context.Networks, NetworkInstance{
				ID:     net.ID,
				Name:   net.Name,
				Driver: net.Driver,
			})
		}
	}
	return context, nil
}
