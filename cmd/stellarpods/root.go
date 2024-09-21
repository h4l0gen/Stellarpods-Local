package main

import (
	"github.com/spf13/cobra"
	"os"
	// "github.com/j4ck4l-24/StellarPods/internal/check_docker"
	"github.com/j4ck4l-24/StellarPods/internal/startrouter"
	"github.com/j4ck4l-24/StellarPods/internal/statuscheck"
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
}
