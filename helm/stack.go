package helm

import (
	"fmt"

	"github.com/rancher/kubectld/cli"
)

func InstallHelmStack(stack *Stack) error {
	args := []string{"install"}
	return executeHelmCreateUpgradeTask(stack, args, false)
}

func DeleteHelmStack(stack *Stack) error {
	args := []string{"delete", stack.Name}
	output := cli.Execute(cmd, args...)
	if output.ExitCode > 0 {
		return fmt.Errorf("%s", output.StdErr)
	}
	return output.Err
}

func UpgradeHelmStack(stack *Stack) error {
	args := []string{"upgrade"}
	return executeHelmCreateUpgradeTask(stack, args, true)
}

func RollbackHelmStack(stack *Stack) error {
	args := []string{"rollback", stack.Name, "0"}
	output := cli.Execute(cmd, args...)
	if output.ExitCode > 0 {
		return fmt.Errorf("%s", output.StdErr)
	}
	return output.Err
}
