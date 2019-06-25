package routes

import (
	"errors"
	"net/http"

	"github.com/docker/docker/api/types/swarm"
	humanize "github.com/dustin/go-humanize"
	docker "github.com/fsouza/go-dockerclient"
)

type ContainerInstance struct {
	ID              string
	Service         string
	Image, ImageTag string
	Node            string
	Created         string
}

type ContainerContext struct {
	BaseContext
	ContainerCount int
	Containers     []ContainerInstance
}

func (router *Router) showContainers(w http.ResponseWriter, r *http.Request) {
	context, err := router.getContainerContext()
	if err != nil {
		router.showError(w, r, err)
		return
	}
	router.templates["containers.html"].Execute(w, context)
}

func (router *Router) getServiceByID(id string) (*swarm.Service, error) {
	for _, ep := range router.endpoints {
		svc, err := ep.InspectService(id)
		if err != nil {
			continue
		}
		return svc, nil
	}
	return nil, errors.New("failed to contact instance")
}

func (router *Router) getContainerContext() (ContainerContext, error) {
	ctx := ContainerContext{
		BaseContext: router.getBaseContext(),
		Containers:  []ContainerInstance{},
	}
	for _, ep := range router.endpoints {
		tasks, err := ep.ListTasks(docker.ListTasksOptions{})
		if err != nil {
			continue
		}
		for _, t := range tasks {
			svc, err := router.getServiceByID(t.ServiceID)
			if err != nil {
				continue
			}
			base, tag := router.getBaseImageFromTagDigest(t.Spec.ContainerSpec.Image)
			ctx.Containers = append(ctx.Containers, ContainerInstance{
				ID:       t.ID,
				Service:  svc.Spec.Name,
				Image:    base,
				ImageTag: tag,
				Node:     t.NodeID,
				Created:  humanize.Time(t.CreatedAt),
			})
			ctx.ContainerCount++
		}
	}
	ctx.SetActive("Containers")
	return ctx, nil
}
