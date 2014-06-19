package docker

import (
	"os/exec"
	"strings"
)

const (
	LXC_LIST_COMMAND = "docker ps --no-trunc | grep -v CONTAINER | awk '{print $1}'"
)

func ListContainers() ([]string, error) {
	output, err := exec.Command("sh", "-c", LXC_LIST_COMMAND).Output()
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return []string{}, nil
	}
	return strings.Split(strings.TrimRight(string(output),"\n"), "\n"), nil
}
