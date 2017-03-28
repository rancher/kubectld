package helm

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/kubectld/cli"
)

func ListReleases() ([]Release, error) {
	args := []string{"ls"}
	output := cli.Execute(cmd, args...)
	if output.ExitCode > 0 {
		return nil, &cli.ErrExec{output}
	}
	if output.Err != nil {
		return nil, output.Err
	}

	releases := []Release{}
	lines := strings.Split(output.StdOut, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parseableString := strings.Map(stripContiguousSpaces, line)
		parseableString = strings.TrimSpace(parseableString)
		parts := strings.Split(parseableString, " ")
		if len(parts) == 5 {
			//Skips the table headers
			continue
		}

		if len(parts) != 9 {
			return nil, fmt.Errorf("Error parsing the output of helm ls: unknown format")
		}
		rev, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Errorf("Error parsing revision %v", err)
			continue
		}
		t, err := time.Parse(time.ANSIC, fmt.Sprintf("%s %s %s %s %s", parts[2], parts[3], parts[4], parts[5], parts[6]))
		if err != nil {
			//just log an error and move on
			log.Errorf("Error parsing release timestamp %v", err)
		}
		release := Release{
			Name:     parts[0],
			Revision: rev,
			Updated:  t,
			Status:   parts[7],
			Chart:    parts[8],
		}
		releases = append(releases, release)
	}
	return releases, nil
}
