package helm

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/kubectld/cli"
	"github.com/rancher/kubectld/kubectl"
)

func InstallHelmStack(stack *Stack) (string, error) {
	args := []string{"install"}
	return executeHelmCreateUpgradeTask(stack, args, false)
}

func DeleteHelmStack(stack *Stack) error {
	releases, err := ListReleases()
	if err != nil {
		log.Errorf("Error obtaining helm releases %v", err)
		return err
	}
	releaseFound := false
	brokenRelease := false
	for _, r := range releases {
		if r.Name == stack.Name {
			if r.Status != "DEPLOYED" && r.Status != "SUPERSEDED" {
				brokenRelease = true
			}
			releaseFound = true
			break
		}
	}
	if !releaseFound {
		return nil
	}
	args := []string{"delete", stack.Name}
	output := cli.Execute(cmd, args...)
	if brokenRelease {
		log.Infof("Tried to delete %s, err: %v", stack.Name, output.Err)
		return nil
	}
	if output.ExitCode > 0 {
		return fmt.Errorf("%s", output.StdErr)
	}
	if output.Err != nil {
		return output.Err
	}
	if stack.Namespace != "" {
		err = kubectl.DeleteNamespace(stack.Namespace)
	}
	return err
}

func UpgradeHelmStack(stack *Stack) (string, error) {
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
