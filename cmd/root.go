package cmd

import (
	galvanicodb "galvanico/cmd/galvanico-db"
	galvaniconotification "galvanico/cmd/galvanico-notifications"
	galvanicoserver "galvanico/cmd/galvanico-server"
	galvanicowebsocket "galvanico/cmd/galvanico-websocket"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"galvanico/internal/logging"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "galvanico",
	Short: "Game CLI",
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

		defer log.Info().Msg("Setup complete")

		return database.Connection().Ping()
	},
	PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
		return database.Close()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "config.yaml", "config file path")

	rootCmd.AddCommand(galvanicoserver.ServeCmd)
	rootCmd.AddCommand(galvanicowebsocket.WsCmd)
	rootCmd.AddCommand(galvaniconotification.NotificationCmd)
	rootCmd.AddCommand(galvanicodb.DBCmd)
}
