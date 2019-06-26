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
	endpointFlag  = "H"
	tlsFlag       = "tls"
	tlsKeyFlag    = "tls-key"
	tlsCertFlag   = "tls-cert"
	tlsCAFlag     = "tls-ca"
	debugFlag     = "debug"
)

func main() {
	app := cli.NewApp()
	app.Name = "spirala"
	app.Usage = "Private Cloud platform based on Docker Swarm"
	app.Copyright = "(c) 2017 - 2019 The Spirala Maintainers"
	app.Author = "The Spirala Maintainers <maintainer@spirala.co>"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   debugFlag,
			EnvVar: "DEBUG",
		},
		cli.StringSliceFlag{
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
			clients []*docker.Client
		)
		// Enable debugging mode if requested
		if ctx.Bool(debugFlag) {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Debug mode enabled")
		}
		// Connect to the first available endpoint
		endpoints := ctx.StringSlice(endpointFlag)
		for _, ep := range endpoints {
			var cli *docker.Client
			var err error
			if ctx.Bool(tlsFlag) {
				cli, err = docker.NewTLSClient(ep, ctx.String(tlsCertFlag), ctx.String(tlsKeyFlag), ctx.String(tlsCAFlag))
			} else {
				cli, err = docker.NewClient(ep)
			}
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":      err,
					"endpoint": ep,
					"tls":      ctx.Bool(tlsFlag),
				}).Warning("Could not connect to endpoint")
			}
			clients = append(clients, cli)
			logrus.WithFields(logrus.Fields{
				"endpoint": ep,
				"tls":      ctx.Bool(tlsFlag),
			}).Info("Connected to endpoint")
		}
		return webui.ListenAndServe(webui.DefaultPort, clients)
	}
	app.Run(os.Args)
}
