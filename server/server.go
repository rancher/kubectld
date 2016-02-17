package server

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func Get(rw http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	output := Kubectl(nil, "get", "-o", "yaml", path)
	writeResponse(rw, output, 200)
}

func Post(rw http.ResponseWriter, r *http.Request) {
	command := mux.Vars(r)["command"]
	output := Kubectl(r.Body, command, "-f", "-")
	writeResponse(rw, output, 203)
}

func writeResponse(rw http.ResponseWriter, output Output, goodCode int) {
	if output.ExitCode > 0 {
		rw.WriteHeader(500)
	} else {
		rw.WriteHeader(goodCode)
	}

	rw.Header()["Content-Type"] = []string{"application/json"}

	respText, err := json.Marshal(output)
	if err != nil {
		logrus.Errorf("Failed to marshall response: %v : %v", err, output)
	}
	if _, err = rw.Write(respText); err != nil {
		logrus.Errorf("Failed to write response: %v : %v", err, string(respText))
	}
}
