package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/utils"
	"github.com/gorilla/mux"
)

var (
	blacklist = map[string]bool{
		"rancher-compose.yaml": true,
		"rancher-compose.yml":  true,
		"docker-compose.yaml":  true,
		"docker-compose.yml":   true,
	}
)

type Server struct {
	Server string
}

type catalogCreate struct {
	Files       map[string]string `json:"files"`
	Environment map[string]string `json:"environment"`
}

func (s *Server) Catalog(rw http.ResponseWriter, r *http.Request) {
	var opts catalogCreate
	var raw project.RawServiceMap

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		writeError(rw, err)
		return
	}

	if err := utils.Convert(opts, &raw); err != nil {
		writeError(rw, err)
		return
	}

	if err := project.Interpolate(&lookup{opts.Environment}, &raw); err != nil {
		writeError(rw, err)
		return
	}

	if err := utils.Convert(raw, &opts); err != nil {
		writeError(rw, err)
		return
	}

	var output Output
	success := false
	worked := []string{}

	defer func() {
		if !success {
			for _, name := range worked {
				Kubectl(strings.NewReader(opts.Files[name]), "-s", s.Server, "delete", "-f", "-")
			}
		}
	}()

	for name, file := range opts.Files {
		if blacklist[name] {
			continue
		}

		output = Kubectl(strings.NewReader(file), "-s", s.Server, "create", "-f", "-")
		if output.ExitCode > 0 {
			writeResponse(rw, output, 203)
			return
		}

		worked = append(worked, name)
	}

	success = true
	writeResponse(rw, output, 203)
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

func writeError(rw http.ResponseWriter, err error) {
	writeResponse(rw, Output{
		ExitCode: -1,
		StdErr:   err.Error(),
		Err:      err,
	}, 400)
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

type lookup struct {
	Vars map[string]string
}

func (l *lookup) Lookup(key, serviceName string, config *project.ServiceConfig) []string {
	ret := l.Vars[key]
	if ret == "" {
		return []string{}
	}
	return []string{fmt.Sprintf("%s=%s", key, ret)}
}
