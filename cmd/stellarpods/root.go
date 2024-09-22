package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/j4ck4l-24/StellarPods/core/github"
	"github.com/j4ck4l-24/StellarPods/core/terraform"
	"github.com/j4ck4l-24/StellarPods/internal/startrouter"
	"github.com/j4ck4l-24/StellarPods/internal/statuscheck"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "StellarPods",
    Short: "A tool for automatically deploying CTF challenges",
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    rootCmd.AddCommand(statuscheck.StatusCheckCmd)
    rootCmd.AddCommand(startrouter.StartRouterCmd)
    rootCmd.AddCommand(deployChallenge())
}

func deployChallenge() *cobra.Command {
    return &cobra.Command{
        Use:   "deploy [owner] [repo] [dir] [path]",
        Short: "Deploy a CTF challenge",
        RunE: func(cmd *cobra.Command, args []string) error {
            var owner, repo, dirToDownload, basePath string

            if len(args) >= 4 {
                owner = args[0]
                repo = args[1]
                dirToDownload = args[2]
                basePath = args[3]
            } else {
                reader := bufio.NewReader(os.Stdin)

                if len(args) > 0 {
                    owner = args[0]
                } else {
                    fmt.Print("Enter repository owner: ")
                    owner, _ = reader.ReadString('\n')
                    owner = strings.TrimSpace(owner)
                }

                if len(args) > 1 {
                    repo = args[1]
                } else {
                    fmt.Print("Enter repository name: ")
                    repo, _ = reader.ReadString('\n')
                    repo = strings.TrimSpace(repo)
                }

                if len(args) > 2 {
                    dirToDownload = args[2]
                } else {
                    fmt.Print("Enter directory to download from repo: ")
                    dirToDownload, _ = reader.ReadString('\n')
                    dirToDownload = strings.TrimSpace(dirToDownload)
                }

                if len(args) > 3 {
                    basePath = args[3]
                } else {
                    fmt.Print("Enter local path to save files: ")
                    basePath, _ = reader.ReadString('\n')
                    basePath = strings.TrimSpace(basePath)
                }
            }

            token := os.Getenv("GITHUB_TOKEN")
            if token == "" {
                return fmt.Errorf("GITHUB_TOKEN environment variable not set")
            }

            fmt.Println("Downloading challenge files...")
            err := github.DownloadChallengeFiles(token, owner, repo, dirToDownload, basePath)
            if err != nil {
                return fmt.Errorf("error downloading challenge files: %w", err)
            }

            // Construct the full path tp the directory container the Terraform files
            tfDir := filepath.Join(basePath, dirToDownload)

            fmt.Println("Running Terraform commands...")
            err = terraform.RunTerraformCommands(tfDir)
            if err != nil {
                if strings.Contains(err.Error(), "no Terraform configuration files found") {
                    fmt.Println("No Terraform files found in the downloaded directory.")
                    fmt.Println("Please ensure that the repository contains .tf files in the specified directory.")
                    return nil
                }
                return fmt.Errorf("error running Terraform commands: %w", err)
            }

            fmt.Println("Challenge deployed successfully!")
            return nil
        },
    }
}

// func createContainerImageCmd() *cobra.Command {
//     return &cobra.Command{
//         Use:   "create-image [dockerfile-path]",
//         Short: "Create a Docker image from a Dockerfile",
//         Args:  cobra.ExactArgs(1),
//         RunE: func(cmd *cobra.Command, args []string) error {
//             return docker.CreateContainerImage(args[0])
//         },
//     }
// }

// func pushImageToGitHubCmd() *cobra.Command {
//     return &cobra.Command{
//         Use:   "push-image [username] [repo-name] [image-name] [tag-name] [token]",
//         Short: "Push a Docker image to GitHub Container Registry",
//         Args:  cobra.ExactArgs(5),
//         RunE: func(cmd *cobra.Command, args []string) error {
//             return github.PushImageToGitHub(args[0], args[1], args[2], args[3], args[4])
//         },
//     }
// }