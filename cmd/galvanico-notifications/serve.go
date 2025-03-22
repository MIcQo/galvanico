package galvanico_server

import (
	"fmt"
	"github.com/spf13/cobra"
)

// NotificationCmd represents the serve command
var NotificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
	},
}

func init() {
	NotificationCmd.PersistentFlags().IntP("port", "p", 8082, "port to serve on")
}
