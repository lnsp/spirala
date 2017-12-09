package routes

import (
	"net/http"

	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

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
	ID              string
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
	var (
		info *docker.DockerInfo
		err  error
	)
	for i := 0; i < len(router.endpoints) && info == nil; i++ {
		info, err = router.endpoints[i].Info()
	}
	if err != nil {
		return SystemContext{}, errors.Wrap(err, "no endpoints reachable")
	}
	return SystemContext{
		ID:              info.ID,
		ServerVersion:   info.ServerVersion,
		KernelVersion:   info.KernelVersion,
		OperatingSystem: info.OperatingSystem,
		Architecture:    info.Architecture,
		Swarm:           info.Swarm.Cluster.ID,
		SwarmVersion:    info.Swarm.Cluster.Version.Index,
		Nodes:           info.Swarm.Nodes,
		Managers:        info.Swarm.Managers,
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
