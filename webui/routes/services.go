package routes

import (
	"net/http"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	humanize "github.com/dustin/go-humanize"
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
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
	Instances    uint64
	Image        string
	ImageTag     string
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

func (router *Router) getBaseImageFromTagDigest(s string) (string, string) {
	components := strings.Split(strings.Split(s, "@")[0], ":")
	return components[0], components[1]
}

func (router *Router) getServiceListContext(limit int) (ServiceListContext, error) {
	serviceList := ServiceListContext{
		BaseContext: router.getBaseContext(),
		Services:    make([]ServiceContext, 0),
	}
	var (
		services []swarm.Service
		err      error
	)
	for i := 0; i < len(router.endpoints) && services == nil; i++ {
		services, err = router.endpoints[i].ListServices(docker.ListServicesOptions{})
	}
	if err != nil {
		return serviceList, errors.Wrap(err, "no endpoint reachable")
	}
	sort.Slice(services, func(i, j int) bool { return services[i].UpdatedAt.After(services[j].UpdatedAt) })
	for i := 0; i < len(services) && (limit == 0 || i < limit); i++ {
		svc := services[i]
		lastUpdate := humanize.Time(svc.UpdatedAt)
		imageName, imageTag := router.getBaseImageFromTagDigest(svc.Spec.TaskTemplate.ContainerSpec.Image)
		instanceMode := "Global"
		instanceCount := uint64(0)
		if svc.Spec.Mode.Replicated != nil {
			instanceMode = "Replicated"
			instanceCount = *svc.Spec.Mode.Replicated.Replicas
		}
		serviceList.Services = append(serviceList.Services, ServiceContext{
			ID:           svc.ID,
			Name:         svc.Spec.Name,
			LastUpdate:   lastUpdate,
			Image:        imageName,
			ImageTag:     imageTag,
			InstanceMode: instanceMode,
			Instances:    instanceCount,
		})
	}
	serviceList.ServiceCount = len(serviceList.Services)
	return serviceList, nil
}

func (router *Router) showServices(w http.ResponseWriter, r *http.Request) {
	context, err := router.getServiceListContext(0)
	if err != nil {
		router.showError(w, r, err)
		return
	}
	context.SetActive("Services")
	router.templates["services.html"].Execute(w, context)
}
