package galvanico_server

import (
	"fmt"
	"github.com/spf13/cobra"
)

// WsCmd represents the serve command
var WsCmd = &cobra.Command{
	Use:   "ws",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
	},
}

func init() {
	WsCmd.PersistentFlags().IntP("port", "p", 8081, "port to serve on")
}
