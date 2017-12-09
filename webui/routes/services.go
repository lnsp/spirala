package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	humanize "github.com/dustin/go-humanize"
	"github.com/fsouza/go-dockerclient"
)

type ServiceListContext struct {
	BaseContext
	Services     []ServiceContext
	ServiceCount int
}

type ServiceContext struct {
	ID           string
	Name         string
	LastUpdate   string
	InstanceMode string
	Image        string
}

func (router *Router) parseServiceUpdateStatus(status *swarm.UpdateStatus) string {
	if status == nil {
		return "Unknown"
	}

	switch status.State {
	case swarm.UpdateStateCompleted:
		return "Update completed"
	case swarm.UpdateStatePaused:
		return "Update paused"
	case swarm.UpdateStateUpdating:
		return "Performing update"
	case swarm.UpdateStateRollbackCompleted:
		return "Rollback completed"
	case swarm.UpdateStateRollbackPaused:
		return "Rollback paused"
	case swarm.UpdateStateRollbackStarted:
		return "Performing rollback"
	}

	return "Unknown"
}

func (router *Router) getBaseImage(s string) string {
	return strings.Split(s, "@")[0]
}

func (router *Router) getServiceListContext() (ServiceListContext, error) {
	serviceList := ServiceListContext{
		BaseContext: router.getBaseContext(),
		Services:    make([]ServiceContext, 0),
	}
	services, err := router.client.ListServices(docker.ListServicesOptions{})
	if err != nil {
		return serviceList, err
	}
	for _, svc := range services {
		lastUpdate := humanize.Time(svc.UpdatedAt)
		baseImage := router.getBaseImage(svc.Spec.TaskTemplate.ContainerSpec.Image)
		instanceMode := "Global"
		if svc.Spec.Mode.Replicated != nil {
			instanceMode = fmt.Sprintf("Replicated (%d)", svc.Spec.Mode.Replicated.Replicas)
		}
		serviceList.Services = append(serviceList.Services, ServiceContext{
			ID:           svc.ID,
			Name:         svc.Spec.Name,
			LastUpdate:   lastUpdate,
			Image:        baseImage,
			InstanceMode: instanceMode,
		})
	}
	serviceList.ServiceCount = len(serviceList.Services)
	return serviceList, nil
}

func (router *Router) showServices(w http.ResponseWriter, r *http.Request) {
	context, err := router.getServiceListContext()
	if err != nil {
		router.showError(w, r, err)
		return
	}
	context.SetActive("Services")
	router.templates["services.html"].Execute(w, context)
}
