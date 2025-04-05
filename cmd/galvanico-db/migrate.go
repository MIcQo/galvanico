package galvanicodb

import (
	"context"
	"galvanico/internal/database"
	"galvanico/migrations"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

// migrateCmd represents the serve command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return createMigrator(cmd.Context(), func(migrator *migrate.Migrator) error {
			group, err := migrator.Migrate(cmd.Context())
			if err != nil {
				return err
			}
			if group.IsZero() {
				log.Info().Msgf("There are no new migrations to run (database is up to date)\n")
				return nil
			}
			log.Info().Msgf("migrated to %s", group)
			return nil
		})
	},
}

func createMigrator(ctx context.Context, m func(migrator *migrate.Migrator) error) error {
	var db = database.Connection()
	var migrator = migrate.NewMigrator(db, migrations.Migrations)

	if err := migrator.Lock(ctx); err != nil {
		return err
	}

	defer func(migrator *migrate.Migrator, ctx context.Context) {
		err := migrator.Unlock(ctx)
		if err != nil {
			log.Panic().Err(err).Msg("failed to unlock migration group")
		}
	}(migrator, ctx)

	return m(migrator)
}
