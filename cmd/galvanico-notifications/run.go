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
		var channel, err = cmd.Flags().GetString("channel")
		if err != nil {
			return err
		}
		return notifications.NewConsumer(cmd.Context(), channel)
	},
}

func init() {
	runCmd.Flags().String("channel", "", "channel name")
}
