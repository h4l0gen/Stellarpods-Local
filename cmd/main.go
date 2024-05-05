package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func checkEnvVars() {
	if os.Getenv("STELLARPODS_PORT") == "" {
		log.Fatal("error: environment variable STELLARPODS_PORT not found")
	}
}
func main() {
	checkEnvVars()
	Execute()
}
