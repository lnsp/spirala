package routes

import "net/http"

func (router *Router) showDashboard(w http.ResponseWriter, r *http.Request) {
	context, err := router.getDashboardContext()
	if err != nil {
		router.showError(w, r, err)
		return
	}
	context.SetActive("Dashboard")
	router.templates["dashboard.html"].Execute(w, context)
}

type SystemContext struct {
	ServerVersion   string
	KernelVersion   string
	OperatingSystem string
	Architecture    string
	Swarm           string
	SwarmVersion    uint64
	Nodes, Managers int
}

type DashboardContext struct {
	BaseContext
	System SystemContext
}

func (router *Router) getSystemContext() (SystemContext, error) {
	dockerInfo, err := router.client.Info()
	if err != nil {
		return SystemContext{}, err
	}
	return SystemContext{
		ServerVersion:   dockerInfo.ServerVersion,
		KernelVersion:   dockerInfo.KernelVersion,
		OperatingSystem: dockerInfo.OperatingSystem,
		Architecture:    dockerInfo.Architecture,
		Swarm:           dockerInfo.Swarm.Cluster.ID,
		SwarmVersion:    dockerInfo.Swarm.Cluster.Version.Index,
		Nodes:           dockerInfo.Swarm.Nodes,
		Managers:        dockerInfo.Swarm.Managers,
	}, nil
}

func (router *Router) getDashboardContext() (DashboardContext, error) {
	dashboard := DashboardContext{
		BaseContext: router.getBaseContext(),
	}
	system, err := router.getSystemContext()
	if err != nil {
		return dashboard, err
	}
	dashboard.System = system
	return dashboard, nil
}
