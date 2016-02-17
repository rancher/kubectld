package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/kubectld/router"
)

func main() {
	listen := ":8091"
	logrus.Info("Starting kubectld on ", listen)
	r := router.New()
	if err := http.ListenAndServe(listen, r); err != nil {
		logrus.Fatal(err)
	}
}
