package galvanico_websocket

import (
	"fmt"
	"github.com/spf13/cobra"
)

// runCmd represents the serve command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
	},
}

func init() {
	runCmd.PersistentFlags().IntP("port", "p", 8082, "port to serve on")
}
