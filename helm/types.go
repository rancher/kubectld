package helm

import (
	"time"
)

const (
	cmd = "helm"
)

type Stack struct {
	Namespace string            `json:"namespace" yaml:"namespace"`
	Name      string            `json:"name" yaml:"name"`
	Files     map[string]string `json:"files" yaml:"files"`
}

type Release struct {
	Name      string                 `json:"name" yaml:"name"`
	Chart     string                 `json:"chart" yaml:"chart"`
	Revision  int                    `json:"revision" yaml:"revision"`
	Status    string                 `json:"status" yaml:"status"`
	Updated   time.Time              `json:"updated" yaml:"updated"`
	Values    map[string]interface{} `json:"values" yaml:"values"`
	Manifests map[string]string      `json:"manifests" yaml:"manifests"`
}
