package router

import (
	"github.com/gorilla/mux"
	"github.com/rancher/kubectld/server"
)

func New(serverURL string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	s := &server.Server{
		Server: serverURL,
	}
	router.Methods("GET").Path("/v1-kubectl/{path:.*}").HandlerFunc(s.Get)
	router.Methods("POST").Path("/v1-kubectl/{command}").HandlerFunc(s.Post)
	return router
}
