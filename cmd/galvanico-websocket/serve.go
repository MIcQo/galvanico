package galvanicowebsocket

import (
	"github.com/spf13/cobra"
)

// WsCmd represents the serve command
var WsCmd = &cobra.Command{
	Use:   "ws",
	Short: "",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
	},
}

func init() {
	WsCmd.AddCommand(sendCmd)
	WsCmd.AddCommand(runCmd)
}
