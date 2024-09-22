package terraform

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// RunTerraformCommands executes Terraform init, plan, and apply
func RunTerraformCommands(workingDir string) error {
	fmt.Printf("Checking for Terraform files in: %s\n", workingDir)
    // Check if Terraform files exist
    tfFiles, err := filepath.Glob(filepath.Join(workingDir, "*.tf"))
    if err != nil {
        return fmt.Errorf("error checking for Terraform files: %w", err)
    }
    if len(tfFiles) == 0 {
        return fmt.Errorf("no Terraform configuration files found in %s", workingDir)
    }

	fmt.Printf("Found %d Terraform files\n", len(tfFiles))

    commands := []string{"init", "plan", "apply --auto-approve"}

    for _, cmd := range commands {
        fmt.Printf("Running: terraform %s\n", cmd)
        command := exec.Command("terraform", strings.Fields(cmd)...)
        command.Dir = workingDir
        command.Stdout = os.Stdout
        command.Stderr = os.Stderr

        err := command.Run()
        if err != nil {
            return fmt.Errorf("error running 'terraform %s': %w", cmd, err)
        }
    }

    return nil
}