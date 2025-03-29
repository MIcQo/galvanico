package galvanicodb

import (
	"galvanico/internal/database"
	"galvanico/migrations"

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

		return migrator.Init(cmd.Context())
	},
}
