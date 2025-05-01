package client

import "os/exec"

var DockerCmd *exec.Cmd

func RunDockerSQL() error {
	DockerCmd = exec.Command("docker-compose", "up", "-d")
	return DockerCmd.Run()
}
