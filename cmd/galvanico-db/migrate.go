package galvanico_db

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

// migrateCmd represents the serve command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	RunE: func(cmd *cobra.Command, args []string) error {
		var ctx = cmd.Context()
		if err := migrator.Lock(ctx); err != nil {
			return err
		}
		defer func(migrator *migrate.Migrator, ctx context.Context) {
			err := migrator.Unlock(ctx)
			if err != nil {
				log.Panic().Err(err).Msg("failed to unlock migration group")
			}
		}(migrator, ctx) //nolint:errcheck

		group, err := migrator.Migrate(ctx)
		if err != nil {
			return err
		}
		if group.IsZero() {
			log.Info().Msgf("there are no new migrations to run (database is up to date)\n")
			return nil
		}
		log.Info().Msgf("migrated to %s\n", group)
		return nil
	},
}
