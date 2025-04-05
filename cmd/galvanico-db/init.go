package galvanicodb

import (
	"galvanico/internal/database"
	"galvanico/migrations"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

// initCmd represents the serve command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "create migration tables",
	RunE: func(cmd *cobra.Command, _ []string) error {
		var db = database.Connection()
		var migrator = migrate.NewMigrator(db, migrations.Migrations)

		err := migrator.Init(cmd.Context())
		if err != nil {
			return err
		}

		log.Info().Msg("Migration tables created")

		log.Info().Msg("Running migrations...")

		return createMigrator(cmd.Context(), func(migrator *migrate.Migrator) error {
			group, migrateErr := migrator.Migrate(cmd.Context())
			if migrateErr != nil {
				return migrateErr
			}
			if group.IsZero() {
				log.Info().Msgf("There are no new migrations to run (database is up to date)")
				return nil
			}
			log.Info().Msgf("migrated to %s", group)
			return nil
		})
	},
}
