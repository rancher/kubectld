package stack

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/utils"
	"github.com/pkg/errors"
	"github.com/rancher/kubectld/cli"
)

var (
	blacklist = map[string]bool{
		"rancher-compose.yaml": true,
		"rancher-compose.yml":  true,
		"docker-compose.yaml":  true,
		"docker-compose.yml":   true,
	}
)

type Input struct {
	Files       map[string]string `yaml:"files"`
	Environment map[string]string `yaml:"-"`
	Server      string            `yaml:"-"`
	Labels      map[string]string `yaml:"-"`
	Namespace   string            `yaml:"-"`
}

func Create(opts Input) error {
	return doOp("apply", true, true, opts)
}

func Remove(opts Input) error {
	var lastErr error
	if err := doOp("delete", false, false, opts); err != nil {
		lastErr = err
	}

	if len(opts.Labels) > 0 {
		args := []string{"-s", opts.Server, fmt.Sprintf("--namespace=%s", opts.Namespace), "delete", "deployment,ep,hpa,ing,rs,rc,svc"}
		for k, v := range opts.Labels {
			args = append(args, "-l", fmt.Sprintf("%s=%s", k, v))
		}
		output := cli.Kubectl(nil, args...)
		if output.ExitCode > 0 {
			lastErr = &cli.ErrExec{output}
		}
	}
	return lastErr
}

func Upgrade(opts Input) error {
	return doOp("apply", true, false, opts)
}

func doOp(action string, label, deleteOnFailure bool, opts Input) error {
	var raw project.RawServiceMap
	var newOpts Input

	if err := utils.Convert(opts, &raw); err != nil {
		return err
	}

	if err := project.Interpolate(&Lookup{opts.Environment}, &raw); err != nil {
		return errors.Wrap(err, "replace variables")
	}

	if err := utils.Convert(raw, &newOpts); err != nil {
		return errors.Wrap(err, "reinterpret after replace variables")
	}

	opts.Files = newOpts.Files

	var output cli.Output
	success := false
	worked := []string{}

	defer func() {
		if !success && deleteOnFailure {
			for _, name := range worked {
				cli.Kubectl(strings.NewReader(opts.Files[name]), "-s", opts.Server, "delete", "-f", "-")
			}
		}
	}()

	for name, file := range opts.Files {
		if blacklist[name] {
			continue
		}
		modifiedConfig := InjectNamespaceToString(file, opts.Namespace)
		output = cli.Kubectl(strings.NewReader(modifiedConfig), "-s", opts.Server, action, "-f", "-")
		if output.ExitCode > 0 {
			logrus.Errorf("Failed with input: %s", modifiedConfig)
			return &cli.ErrExec{output}
		}

		if len(opts.Labels) > 0 && label {
			args := []string{"-s", opts.Server, "label", "--overwrite", "-f", "-"}
			for k, v := range opts.Labels {
				args = append(args, fmt.Sprintf("%s=%s", k, v))
			}
			output = cli.Kubectl(strings.NewReader(modifiedConfig), args...)
			if output.ExitCode > 0 {
				logrus.Errorf("Failed with input: %s", modifiedConfig)
				return &cli.ErrExec{output}
			}
		}

		worked = append(worked, name)
	}

	success = true
	return nil
}
