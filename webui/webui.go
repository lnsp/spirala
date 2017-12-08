package webui

import (
	"net/http"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/lnsp/spirala/webui/routes"
)

const DefaultPort = ":8080"

func ListenAndServe(hostport string, client *docker.Client) error {
	srv := &http.Server{
		Handler:      routes.New(client),
		Addr:         hostport,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
