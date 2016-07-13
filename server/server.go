package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rancher/kubectld/cli"
	"github.com/rancher/kubectld/stack"
)

const defaultNamespace = "defaultNamespace"

type Server struct {
	Server string
}

func (s *Server) Catalog(rw http.ResponseWriter, r *http.Request) {
	var opts stack.Input
	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		writeError(rw, err)
		return
	}

	opts.Namespace = r.FormValue(defaultNamespace)
	opts.Server = s.Server

	err = stack.Create(opts)
	if cliErr, ok := err.(*cli.ErrExec); ok {
		writeResponse(rw, cliErr.Output, 203)
		return
	} else if err != nil {
		writeError(rw, err)
		return
	}

	writeResponse(rw, cli.Output{}, 203)
}

func (s *Server) Get(rw http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	namespaceVar := "--all-namespaces"
	namespace := r.URL.Query()["namespace"]
	if len(namespace) > 0 {
		namespaceVar = "--namespace=" + namespace[0]
	}
	output := cli.Kubectl(nil, "-s", s.Server, "get", namespaceVar, "-o", "yaml", path)
	writeResponse(rw, output, 200)
}

func (s *Server) Post(rw http.ResponseWriter, r *http.Request) {
	command := mux.Vars(r)["command"]
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(rw, err)
		return
	}
	modifiedConfig := string(stack.InjectNamespace(b, r.FormValue(defaultNamespace)))
	output := cli.Kubectl(strings.NewReader(modifiedConfig), "-s", s.Server, command, "-f", "-")
	writeResponse(rw, output, 203)
}

func writeError(rw http.ResponseWriter, err error) {
	writeResponse(rw, cli.Output{
		ExitCode: -1,
		StdErr:   err.Error(),
		Err:      err,
	}, 400)
}

func writeResponse(rw http.ResponseWriter, output cli.Output, goodCode int) {
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
