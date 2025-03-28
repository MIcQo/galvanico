package galvanicodb

import (
	"github.com/spf13/cobra"
)

// initCmd represents the serve command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create migration tables",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return migrator.Init(cmd.Context())
	},
}
