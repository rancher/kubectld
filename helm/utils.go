package helm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/rancher/kubectld/cli"
)

func executeHelmCreateUpgradeTask(stack *Stack, args []string, isUpgrade bool) (string, error) {
	if stack.Namespace != "" {
		args = append(args, "--namespace", stack.Namespace)
	}
	if stack.Name != "" {
		if isUpgrade {
			args = append(args, stack.Name)
		} else { //create
			args = append(args, "--name", stack.Name)
		}
	} else {
		return "", fmt.Errorf("KubernetesStack.Name cannot be empty")
	}
	dir, err := ioutil.TempDir("", "helm-templates")
	if err != nil {
		return "", err
	}
	err = os.Chmod(dir, 0700)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	helmPath := ""
	for name, data := range stack.Files {
		index := strings.LastIndex(name, "/")
		if index == -1 {
			index = 0
		}
		cd := name[:index]
		err = os.MkdirAll(path.Join(dir, cd), 0700)
		if err != nil {
			return "", err
		}
		if strings.HasSuffix(name, "Chart.yaml") {
			helmPath = path.Join(dir, cd)
		}
		f, err := os.Create(path.Join(dir, name))
		if err != nil {
			return "", err
		}
		defer f.Close()
		_, err = f.WriteString(data)
		if err != nil {
			return "", err
		}
	}
	args = append(args, helmPath)
	depArgs := []string{"dependency", "update", helmPath}
	output := cli.Execute(cmd, depArgs...)
	if output.ExitCode > 0 {
		return "", fmt.Errorf("%s", output.StdErr)
	}
	if output.Err != nil {
		return "", output.Err
	}
	output = cli.Execute(cmd, args...)
	if output.ExitCode > 0 {
		return "", fmt.Errorf("%s", output.StdErr)
	}
	return output.StdOut, output.Err
}

var contiguousSpace = false

func stripContiguousSpaces(r rune) rune {
	if unicode.IsSpace(r) {
		if contiguousSpace {
			return -1
		}
		contiguousSpace = true
		return ' '
	}
	contiguousSpace = false
	return r
}
