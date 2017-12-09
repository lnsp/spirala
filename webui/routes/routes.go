package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	includeDirectory = "includes/*.html"
	layoutDirectory  = "layouts/*.html"
	staticDirectory  = "static"
	baseTemplate     = "base"
)

type Router struct {
	mux       *mux.Router
	endpoints []*docker.Client
	templates map[string]*template.Template
}

type NavContext struct {
	Link, Name string
	Active     bool
}

type SiteContext struct {
	Title string
	Year  int
	Nav   []NavContext
}

type ErrorContext struct {
	BaseContext
	Err string
}

func (ctx SiteContext) SetActive(nav string) SiteContext {
	for i := range ctx.Nav {
		if ctx.Nav[i].Name == nav {
			ctx.Nav[i].Active = true
		}
	}
	return ctx
}

type BaseContext struct {
	SiteContext
}

func (router *Router) showError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Fprintln(w, err)
}

func (router *Router) getBaseContext() BaseContext {
	return BaseContext{
		SiteContext{
			Title: "Spirala",
			Year:  time.Now().Year(),
			Nav: []NavContext{
				// Swarm overview
				{
					Link: "/",
					Name: "Dashboard",
				},
				// Services and stacks
				{
					Link: "/services",
					Name: "Services",
				},
				// Node and swarm management
				{
					Link: "/nodes",
					Name: "Nodes",
				},
				// Volumes and secrets
				{
					Link: "/volumes",
					Name: "Volumes",
				},
				// Networks
				{
					Link: "/networks",
					Name: "Networks",
				},
				// Docker images
				{
					Link: "/images",
					Name: "Images",
				},
				// Docker containers
				{
					Link: "/containers",
					Name: "Containers",
				},
			},
		},
	}
}

func (router *Router) initRoutes() error {
	// Initialize base routes
	router.mux.HandleFunc("/", router.showDashboard).Methods("GET")
	router.mux.HandleFunc("/services", router.showServices).Methods("GET")
	router.mux.HandleFunc("/nodes", router.showNodes).Methods("GET")
	router.mux.HandleFunc("/volumes", router.showVolumes).Methods("GET")
	router.mux.HandleFunc("/networks", router.showNetworks).Methods("GET")
	router.mux.HandleFunc("/images", router.showImages).Methods("GET")
	router.mux.HandleFunc("/containers", router.showContainers).Methods("GET")

	// Initialize static files
	router.mux.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDirectory))))
	return nil
}

func (router *Router) initTemplates() error {
	router.templates = make(map[string]*template.Template)
	layouts, err := filepath.Glob(layoutDirectory)
	if err != nil {
		return errors.Wrap(err, "could not fetch layouts")
	}
	includes, err := filepath.Glob(includeDirectory)
	if err != nil {
		return errors.Wrap(err, "could not fetch includes")
	}
	for _, layout := range layouts {
		files := append(includes, layout)
		tmp, err := template.ParseFiles(files...)
		if err != nil {
			return errors.Wrap(err, "could not parse template")
		}
		logrus.WithFields(logrus.Fields{
			"name": filepath.Base(layout),
		}).Debug("Added template to route")
		router.templates[filepath.Base(layout)] = tmp
	}
	return nil
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

// New instantiates a new Web UI route handler.
func New(endpoints []*docker.Client) http.Handler {
	// TODO: Add route implementations
	router := &Router{
		mux:       mux.NewRouter(),
		endpoints: endpoints,
	}
	if err := router.initTemplates(); err != nil {
		logrus.Panic(err)
	}
	if err := router.initRoutes(); err != nil {
		logrus.Panic(err)
	}
	return router
}
