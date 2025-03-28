package galvanico_db

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

// rollbackCmd represents the serve command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback the last migration group",
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

		group, err := migrator.Rollback(ctx)
		if err != nil {
			return err
		}
		if group.IsZero() {
			log.Info().Msg("there are no groups to roll back")
			return nil
		}
		log.Info().Msgf("rolled back %s\n", group)
		return nil
	},
}
