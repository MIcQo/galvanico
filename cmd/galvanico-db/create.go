package galvanicodb

import (
	"galvanico/internal/database"
	"galvanico/migrations"
	"strings"

	"github.com/uptrace/bun/migrate"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// createCmd represents the serve command
var createCmd = &cobra.Command{
	Use:   "create {name}",
	Short: "create Go migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		var db = database.Connection()
		var migrator = migrate.NewMigrator(db, migrations.Migrations)

		var ctx = cmd.Context()
		name := strings.Join(args, "_")
		mf, err := migrator.CreateSQLMigrations(ctx, name)
		if err != nil {
			return err
		}
		for _, m := range mf {
			log.Info().Msgf("created migration %s", m.Name)
		}
		return nil
	},
}
