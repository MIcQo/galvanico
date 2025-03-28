package galvanicodb

import (
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// createCmd represents the serve command
var createCmd = &cobra.Command{
	Use:   "create {name}",
	Short: "create Go migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx = cmd.Context()
		name := strings.Join(args, "_")
		mf, err := migrator.CreateGoMigration(ctx, name)
		if err != nil {
			return err
		}
		log.Info().Msgf("created migration %s (%s)\n", mf.Name, mf.Path)
		return nil
	},
}
