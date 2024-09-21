package env

import (
	"fmt"
	"os"
)

// CheckEnvVars checks if required environment variables are set
func CheckEnvVars() error {
	if os.Getenv("STELLARPODS_PORT") == "" {
		return fmt.Errorf("error: environment variable STELLARPODS_PORT not found")
	}
	return nil
}