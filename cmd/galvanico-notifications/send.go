package galvaniconotification

import (
	"galvanico/internal/notifications"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const sendCmdMaxArgs = 2

// sendCmd represents the serve command
var sendCmd = &cobra.Command{
	Use:   "send {channel} {msg}",
	Short: "Send message to debug",
	Long:  ``,
	Args:  cobra.ExactArgs(sendCmdMaxArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		var repetations, err = cmd.Flags().GetInt("repeat")
		if err != nil {
			return err
		}

		if repetations < 1 {
			repetations = 1
		}

		for range repetations {
			if sendErr := notifications.NewPublisher(args[0], []byte(strings.Join(args[1:], " "))); sendErr != nil {
				log.Panic().Err(sendErr).Msg("unable to send message")
			}
		}

		return nil
	},
}

func init() {
	sendCmd.Flags().IntP("repeat", "r", 1, "Number of times to repeat the message")
}
