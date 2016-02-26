package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/kubectld/router"

	flag "github.com/ogier/pflag"
)

var (
	listen = flag.String("listen", ":8091", "Listen address")
	server = flag.String("server", "http://localhost:8080", "Kubernetes server address")
)

func main() {
	flag.Parse()

	logrus.Info("Starting kubectld on ", *listen)
	r := router.New(*server)
	if err := http.ListenAndServe(*listen, r); err != nil {
		logrus.Fatal(err)
	}
}
