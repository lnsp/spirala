package main

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "spirala"
	app.Usage = "Management UI for Docker Swarm"
	app.Copyright = "(c) 2017 The Spirala Maintainers"
	app.Author = "The Spirala Maintainers <maintainer@spirala.co>"
	app.Version = "0.1.0"
	app.Action = func(ctx *cli.Context) error {
		endpoint := "unix:///var/run/docker.sock"
		client, err := docker.NewClient(endpoint)
		if err != nil {
			panic(err)
		}
		imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
		if err != nil {
			panic(err)
		}
		for _, img := range imgs {
			fmt.Println("ID: ", img.ID)
			fmt.Println("RepoTags: ", img.RepoTags)
			fmt.Println("Created: ", img.Created)
			fmt.Println("Size: ", img.Size)
			fmt.Println("VirtualSize: ", img.VirtualSize)
			fmt.Println("ParentId: ", img.ParentID)
		}
		return nil
	}
	app.Run(os.Args)
}
