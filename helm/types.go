package helm

const (
	cmd = "helm"
)

type Stack struct {
	Namespace string            `json:"namespace" yaml:"namespace"`
	Name      string            `json:"name" yaml:"name"`
	Files     map[string]string `json:"files" yaml:"files"`
}
