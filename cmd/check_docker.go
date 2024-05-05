package main

import (
    "errors"
    "os/exec"
)

func isDockerInstalled() error {
    _, err := exec.LookPath("docker")
    if err != nil {
        return errors.New("error: Docker is not installed")
    }
    return nil
}