package startrouter

import (
	"github.com/j4ck4l-24/StellarPods/api"
	"github.com/spf13/cobra"
)

// StartRouterCmd represents the start-router command
var StartRouterCmd = &cobra.Command{
	Use: "start-router",
	Run: func(cmd *cobra.Command, args []string) {
		api.StartRouter()
	},
	PreRun: func(cmd *cobra.Command, args []string) {

	},
}
