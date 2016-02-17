package router

import (
	"github.com/gorilla/mux"
	"github.com/rancher/kubectld/server"
)

func New() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/v1-kubectl/{path:.*}").HandlerFunc(server.Get)
	router.Methods("POST").Path("/v1-kubectl/{command}").HandlerFunc(server.Post)
	return router
}
