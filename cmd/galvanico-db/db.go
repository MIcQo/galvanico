package galvanico_db

import (
	"galvanico/internal/database"
	"galvanico/migrations"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

var migrator *migrate.Migrator

// DBCmd represents the serve command
var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var db = database.Connection()
		migrator = migrate.NewMigrator(db, migrations.Migrations)
		return nil
	},
}

func init() {
	DBCmd.AddCommand(createCmd)
	DBCmd.AddCommand(initCmd)
	DBCmd.AddCommand(migrateCmd)
	DBCmd.AddCommand(rollbackCmd)
}
