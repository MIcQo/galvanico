package galvanicodb

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

// rollbackCmd represents the serve command.
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback the last migration group",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return createMigrator(cmd.Context(), func(migrator *migrate.Migrator) error {
			group, err := migrator.Rollback(cmd.Context())
			if err != nil {
				return err
			}
			if group.IsZero() {
				log.Info().Msg("there are no groups to roll back")
				return nil
			}
			log.Info().Msgf("rolled back %s\n", group)
			return nil
		})
	},
}
