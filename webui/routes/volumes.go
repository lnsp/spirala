package routes

import (
	"net/http"

	humanize "github.com/dustin/go-humanize"
	docker "github.com/fsouza/go-dockerclient"
)

type VolumeInstance struct {
	Name    string
	Driver  string
	Created string
}

type VolumeContext struct {
	BaseContext
	VolumeCount int
	Volumes     []VolumeInstance
}

func (router *Router) showVolumes(w http.ResponseWriter, r *http.Request) {
	context, err := router.getVolumeContext()
	if err != nil {
		router.showError(w, r, err)
		return
	}
	router.templates["volumes.html"].Execute(w, context)
}

func (router *Router) getVolumeContext() (VolumeContext, error) {
	context := VolumeContext{
		BaseContext: router.getBaseContext(),
		Volumes:     []VolumeInstance{},
	}
	context.SetActive("Volumes")
	for _, ep := range router.endpoints {
		volumes, err := ep.ListVolumes(docker.ListVolumesOptions{})
		if err != nil {
			return VolumeContext{}, err
		}
		for _, vol := range volumes {
			context.Volumes = append(context.Volumes, VolumeInstance{
				Name:    vol.Name,
				Driver:  vol.Driver,
				Created: humanize.Time(vol.CreatedAt),
			})
			context.VolumeCount++
		}

	}
	return context, nil
}
