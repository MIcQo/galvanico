package galvanico_server

import (
	"fmt"
	"github.com/spf13/cobra"
)

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
	},
}

func init() {
	ServeCmd.PersistentFlags().IntP("port", "p", 8080, "port to serve on")
}
