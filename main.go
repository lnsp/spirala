package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/lnsp/spirala/webui"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	localEndpoint = "unix:///var/run/docker.sock"
	endpointFlag  = "endpoint"
	tlsFlag       = "tls"
	tlsKeyFlag    = "tls-key"
	tlsCertFlag   = "tls-cert"
	tlsCAFlag     = "tls-ca"
)

func main() {
	app := cli.NewApp()
	app.Name = "spirala"
	app.Usage = "Private Cloud platform based on Docker Swarm"
	app.Copyright = "(c) 2017 The Spirala Maintainers"
	app.Author = "The Spirala Maintainers <maintainer@spirala.co>"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   endpointFlag,
			EnvVar: "DOCKER_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   tlsFlag,
			EnvVar: "DOCKER_TLS",
		},
		cli.StringFlag{
			Name:   tlsCertFlag,
			EnvVar: "DOCKER_TLS_CERT",
		},
		cli.StringFlag{
			Name:   tlsKeyFlag,
			EnvVar: "DOCKER_TLS_KEY",
		},
		cli.StringFlag{
			Name:   tlsCAFlag,
			EnvVar: "DOCKER_TLS_CA",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		var (
			client   *docker.Client
			err      error
			endpoint string
		)
		if declaredEndpoint := ctx.String(endpointFlag); endpoint != declaredEndpoint {
			endpoint = declaredEndpoint
		} else {
			endpoint = localEndpoint
		}
		if enableTLS := ctx.Bool(tlsFlag); enableTLS {
			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint,
			})
			client, err = docker.NewTLSClient(endpoint, ctx.String(tlsCertFlag), ctx.String(tlsKeyFlag), ctx.String(tlsCAFlag))
		} else {
			client, err = docker.NewClient(endpoint)
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Panic("Could not connect to endpoint")
		} else {
			logrus.WithFields(logrus.Fields{
				"endpoint": endpoint,
				"tls":      ctx.Bool(tlsFlag),
			}).Info("Connected to endpoint")
		}
		return webui.ListenAndServe(webui.DefaultPort, client)
	}
	app.Run(os.Args)
}
