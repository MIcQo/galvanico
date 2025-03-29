package galvanicodb

import (
	"galvanico/internal/config"
	"galvanico/internal/database"
	"galvanico/internal/logging"

	"github.com/spf13/cobra"
)

// DBCmd represents the serve command
var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		config.FileName = cmd.Flag("config").Value.String() //nolint:reassign // we need to set config

		if err := logging.Setup(); err != nil {
			return err
		}

		var cfg, cfgErr = config.Load()
		if cfgErr != nil {
			return cfgErr
		}

		if err := logging.SetupLevel(cfg); err != nil {
			return err
		}

		return database.Connection().Ping()
	},
}

func init() {
	DBCmd.AddCommand(createCmd)
	DBCmd.AddCommand(initCmd)
	DBCmd.AddCommand(migrateCmd)
	DBCmd.AddCommand(rollbackCmd)
}
