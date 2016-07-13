package stack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	utilyaml "github.com/kubernetes/kubernetes/pkg/util/yaml"
)

type Metadata struct {
	Name        string      `json:"name,omitempty"`
	Namespace   string      `json:"namespace,omitempty"`
	Labels      interface{} `json:"labels,omitempty"`
	Annotations interface{} `json:"annotations,omitempty"`
}

type KubernetesResource struct {
	Kind       string      `json:"kind"`
	APIVersion string      `json:"apiVersion,omitempty"`
	Metadata   Metadata    `json:"metadata,omitempty"`
	Spec       interface{} `json:"spec,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Template   interface{} `json:"template,omitempty"`
	Items      interface{} `json:"items,omitempty"`
}

func InjectNamespaceToString(originalConfig string, namespace string) string {
	if namespace == "" {
		return originalConfig
	}
	return string(InjectNamespace([]byte(originalConfig), namespace))
}

func InjectNamespace(originalConfig []byte, namespace string) []byte {
	if namespace == "" {
		return originalConfig
	}
	decoder := utilyaml.NewYAMLOrJSONDecoder(bytes.NewReader(originalConfig), len(originalConfig))
	resources := []KubernetesResource{}
	var err error
	for {
		out := KubernetesResource{}
		err = decoder.Decode(&out)
		if err != nil {
			break
		}
		resources = append(resources, out)
	}
	if err != io.EOF {
		logrus.Infof("Failed to read config  %v", err)
		return originalConfig
	}

	var toReturn []byte
	for _, resource := range resources {
		if resource.Metadata.Namespace == "" {
			resource.Metadata.Namespace = namespace
		}
		modifiedConfig, err := json.Marshal(&resource)
		if err != nil {
			logrus.Errorf("Error marshaling %v", err)
			return originalConfig
		}
		toReturn = append(toReturn, modifiedConfig...)
	}
	return toReturn
}

type Lookup struct {
	Vars map[string]string
}

func (l *Lookup) Lookup(key, serviceName string, config *project.ServiceConfig) []string {
	ret := l.Vars[key]
	if ret == "" {
		return []string{}
	}
	return []string{fmt.Sprintf("%s=%s", key, ret)}
}
