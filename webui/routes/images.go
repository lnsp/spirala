package routes

import (
	"net/http"
	"sort"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	docker "github.com/fsouza/go-dockerclient"
)

type ImageContext struct {
	ID         string
	Repository string
	Tag        string
	Size       string
	Created    string
}

type ImageListContext struct {
	BaseContext
	Images     []ImageContext
	ImageCount int
}

func (router *Router) showImages(w http.ResponseWriter, r *http.Request) {
	context, err := router.getImageListContext(0)
	if err != nil {
		router.showError(w, r, err)
		return
	}
	context.SetActive("Images")
	router.templates["images.html"].Execute(w, context)
}

func (router *Router) getBaseImageFromTag(tag string) (string, string) {
	components := strings.Split(tag, ":")
	if len(components) < 2 {
		return tag, ""
	}
	return components[0], components[1]
}

const NoneTag = "<none>"

func (router *Router) getImageListContext(limit int) (ImageListContext, error) {
	imageListContext := ImageListContext{
		BaseContext: router.getBaseContext(),
		Images:      []ImageContext{},
	}
	var (
		images []docker.APIImages
		err    error
	)
	for i := 0; i < len(router.endpoints) && images == nil; i++ {
		images, err = router.endpoints[i].ListImages(docker.ListImagesOptions{})
	}
	if err != nil {
		return imageListContext, err
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].Created < images[i].Created
	})
	for i := 0; i < len(images) && (limit == 0 || i < limit); i++ {
		// ignore untagged images
		if len(images[i].RepoTags) < 1 {
			continue
		}
		_, digest := router.getBaseImageFromTag(images[i].ID)
		base, tag := router.getBaseImageFromTag(images[i].RepoTags[0])
		// ignore non-tagged images
		if base == NoneTag {
			continue
		}
		imageListContext.Images = append(imageListContext.Images, ImageContext{
			ID:         digest[:16],
			Repository: base,
			Tag:        tag,
			Size:       humanize.Bytes(uint64(images[i].Size)),
			Created:    humanize.Time(time.Unix(images[i].Created, 0)),
		})
	}
	imageListContext.ImageCount = len(imageListContext.Images)
	return imageListContext, nil
}
