package server

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
)

type Output struct {
	ExitCode int    `json:"exitCode"`
	StdOut   string `json:"stdOut"`
	StdErr   string `json:"stdErr"`
	Err      error  `json:"error"`
}

func Kubectl(stdIn io.Reader, args ...string) Output {
	var (
		out, outErr bytes.Buffer
	)

	logrus.Info("kubectl ", strings.Join(args, " "))

	c := exec.Command("kubectl", args...)
	c.Stdin = stdIn
	c.Stdout = &out
	c.Stderr = &outErr

	output := Output{}
	output.Err = c.Run()
	output.StdOut = string(out.Bytes())
	output.StdErr = string(outErr.Bytes())

	if exitError, ok := output.Err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			output.ExitCode = waitStatus.ExitStatus()
		}
	}

	logrus.Infof("output: %#v", output)

	return output
}
