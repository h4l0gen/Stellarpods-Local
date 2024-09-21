package main

import (
	"github.com/j4ck4l-24/StellarPods/internal/check_docker"
	"github.com/j4ck4l-24/StellarPods/internal/env"
	log "github.com/sirupsen/logrus"
	// "os"
	
)

func main() {
	if err := check_docker.IsDockerInstalled(); err != nil {
		log.Fatal(err)
	}

	if err := env.CheckEnvVars(); err != nil {
		log.Fatal(err)
	}

	Execute()
}
