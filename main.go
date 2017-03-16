package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/rancher/kubectld/events"
	"github.com/rancher/swarm-agent/healthcheck"
)

func main() {
	app := cli.NewApp()
	app.Name = "kubectld"
	app.Action = launch

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "cattle-url",
			Usage:  "URL for cattle API",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "cattle-access-key",
			Usage:  "Cattle API Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "cattle-secret-key",
			Usage:  "Cattle API Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:   "worker-count",
			Value:  50,
			Usage:  "Number of workers for handling events",
			EnvVar: "WORKER_COUNT",
		},
		cli.IntFlag{
			Name:   "health-check-port",
			Value:  10240,
			Usage:  "Port to configure an HTTP health check listener on",
			EnvVar: "HEALTH_CHECK_PORT",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debug logs",
			EnvVar: "DEBUG",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("Fatal exit: %v", err)
	}
}

func launch(ctx *cli.Context) error {
	hcPort := ctx.Int("health-check-port")

	url := ctx.String("cattle-url")
	accessKey := ctx.String("cattle-access-key")
	secretKey := ctx.String("cattle-secret-key")
	workers := ctx.Int("worker-count")

	if ctx.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	go func() {
		logrus.Fatalf("Rancher healthcheck exited with error: %v", healthcheck.StartHealthCheck(hcPort))
	}()

	return events.StartEventHandler(url, accessKey, secretKey, workers)
}
