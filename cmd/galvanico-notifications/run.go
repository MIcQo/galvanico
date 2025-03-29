package galvaniconotification

import (
	"galvanico/internal/notifications"

	"github.com/spf13/cobra"
)

// runCmd represents the serve command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return notifications.NewConsumer(cmd.Context())
	},
}

func init() {
}
