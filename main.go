package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/kubectld/events"
	"github.com/rancher/kubectld/router"
	"github.com/rancher/swarm-agent/healthcheck"
)

func main() {
	app := cli.NewApp()
	app.Name = "kubectld"
	app.Action = launch

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Usage: "Kubernetes server address",
			Value: "http://localhost:8080",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "Listen address",
			Value: ":8091",
		},
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
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("Fatal exit: %v", err)
	}
}

func listenToEvents(server, url, accessKey, secretKey string, workers int) error {
	stackHandler := &events.StackHandler{
		Server: server,
	}
	eventHandlers := map[string]revents.EventHandler{
		"kubernetesStack.create":        stackHandler.Create,
		"kubernetesStack.finishupgrade": stackHandler.FinishUpgrade,
		"kubernetesStack.upgrade":       stackHandler.Upgrade,
		"kubernetesStack.rollback":      stackHandler.Rollback,
		"kubernetesStack.remove":        stackHandler.Remove,
		"ping": (&events.PingHandler{}).Handler,
	}

	router, err := revents.NewEventRouter("", 0, url, accessKey, secretKey, nil, eventHandlers, "", workers)
	if err != nil {
		return err
	}
	return router.StartWithoutCreate(nil)
}

func launch(ctx *cli.Context) error {
	hcPort := ctx.Int("health-check-port")
	listen := ctx.String("listen")
	server := ctx.String("server")

	url := ctx.String("cattle-url")
	accessKey := ctx.String("cattle-access-key")
	secretKey := ctx.String("cattle-secret-key")
	workers := ctx.Int("worker-count")

	go func() {
		logrus.Fatalf("Rancher healthcheck exited with error: %v", healthcheck.StartHealthCheck(hcPort))
	}()

	go func() {
		logrus.Fatalf("Failed to listen to events: %v", listenToEvents(server, url, accessKey, secretKey, workers))
	}()

	logrus.Info("Starting kubectld on ", listen)
	r := router.New(server)
	return http.ListenAndServe(listen, r)
}
