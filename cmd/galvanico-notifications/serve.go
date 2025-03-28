package galvaniconotification

import (
	"github.com/spf13/cobra"
)

// NotificationCmd represents the serve command
var NotificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {

	},
}

func init() {
	NotificationCmd.AddCommand(sendCmd)
	NotificationCmd.AddCommand(runCmd)
}
