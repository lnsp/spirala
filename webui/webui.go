package webui

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/lnsp/spirala/webui/routes"
)

const DefaultPort = ":8080"

type LogRequestMiddleware struct {
	handler http.Handler
}

func (middleware LogRequestMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"host": r.Host,
		"uri":  r.RequestURI,
	}).Debug("Incoming request")
	middleware.handler.ServeHTTP(w, r)
}

func ListenAndServe(hostport string, clients []*docker.Client) error {
	srv := &http.Server{
		Handler:      &LogRequestMiddleware{routes.New(clients)},
		Addr:         hostport,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
