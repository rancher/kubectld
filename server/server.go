package server

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type Server struct {
	Server string
}

func (s *Server) Get(rw http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	namespaceVar := "--all-namespaces"
	namespace := r.URL.Query()["namespace"]
	if len(namespace) > 0 {
		namespaceVar = "--namespace=" + namespace[0]
	}
	output := Kubectl(nil, "-s", s.Server, "get", namespaceVar, "-o", "yaml", path)
	writeResponse(rw, output, 200)
}

func (s *Server) Post(rw http.ResponseWriter, r *http.Request) {
	command := mux.Vars(r)["command"]
	output := Kubectl(r.Body, "-s", s.Server, command, "-f", "-")
	writeResponse(rw, output, 203)
}

func writeResponse(rw http.ResponseWriter, output Output, goodCode int) {
	rw.Header()["Content-Type"] = []string{"application/json"}

	if output.ExitCode > 0 {
		rw.WriteHeader(500)
	} else {
		rw.WriteHeader(goodCode)
	}

	respText, err := json.Marshal(output)
	if err != nil {
		logrus.Errorf("Failed to marshall response: %v : %v", err, output)
	}
	if _, err = rw.Write(respText); err != nil {
		logrus.Errorf("Failed to write response: %v : %v", err, string(respText))
	}
}
