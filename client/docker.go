package client

import (
	"bytes"
	"os/exec"
	"strings"
)

// Dockerコンテナが起動しているか確認
func IsDockerPostgresRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", "name=local-postgres", "--format", "{{.Names}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}
	containers := strings.Split(out.String(), "\n")
	for _, name := range containers {
		if name == "local-postgres" {
			return true
		}
	}
	return false
}

var DockerCmd *exec.Cmd

func RunDockerSQL() error {
	DockerCmd = exec.Command("docker-compose", "up", "-d")
	return DockerCmd.Run()
}
