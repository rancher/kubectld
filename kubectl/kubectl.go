package kubectl

import (
	"fmt"

	"github.com/rancher/kubectld/cli"
)

func DeleteNamespace(namespace string) error {
	cmd := "kubectl"
	args := []string{"delete", "namespace", namespace}
	output := cli.Execute(cmd, args...)
	if output.ExitCode > 0 {
		return fmt.Errorf("%s", output.StdErr)
	}
	return output.Err
}
