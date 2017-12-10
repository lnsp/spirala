package routes

import (
	"net/http"
	"sort"

	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"

	"github.com/docker/docker/api/types/swarm"
)

type NodeContext struct {
	Status       string
	ID           string
	Hostname     string
	Availability string
	Role         string
}

type NodeListContext struct {
	BaseContext
	Nodes     []NodeContext
	NodeCount int
}

func (router *Router) getNodeListContext(limit int) (NodeListContext, error) {
	nodeList := NodeListContext{
		BaseContext: router.getBaseContext(),
		Nodes:       []NodeContext{},
	}
	var (
		nodes []swarm.Node
		err   error
	)
	for i := 0; i < len(router.endpoints) && nodes == nil; i++ {
		nodes, err = router.endpoints[i].ListNodes(docker.ListNodesOptions{})
	}
	if err != nil {
		return nodeList, errors.Wrap(err, "could not fetch nodes")
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].UpdatedAt.After(nodes[j].UpdatedAt) })
	for i := 0; i < len(nodes) && (limit == 0 || i < limit); i++ {
		n := nodes[i]
		nodeList.Nodes = append(nodeList.Nodes, NodeContext{
			Status:       string(n.Status.State),
			ID:           n.ID,
			Hostname:     n.Description.Hostname,
			Role:         string(n.Spec.Role),
			Availability: string(n.Spec.Availability),
		})
	}
	nodeList.NodeCount = len(nodeList.Nodes)
	return nodeList, nil
}

func (router *Router) showNodes(w http.ResponseWriter, r *http.Request) {
	context, err := router.getNodeListContext(0)
	if err != nil {
		router.showError(w, r, err)
		return
	}
	context.SetActive("Nodes")
	router.templates["nodes.html"].Execute(w, context)
}
