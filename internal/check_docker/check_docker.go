package check_docker

import (
    "errors"
    "os/exec"
)


// IsDockerInstalled checks if Docker is installed on the system
func IsDockerInstalled() error {
    _, err := exec.LookPath("docker")
    if err != nil {
        return errors.New("error: Docker is not installed")
    }
    return nil
}