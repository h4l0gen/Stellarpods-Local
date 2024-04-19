package main

import (
	"github.com/j4ck4l-24/StellarPods/api"
	"github.com/spf13/cobra"
)

var startRouterCmd = &cobra.Command{
	Use: "start-router",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartRouter()
	},
	PreRun: func(cmd *cobra.Command, args []string) {

	},
}
