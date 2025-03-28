package cmd

import (
	galvanicodb "galvanico/cmd/galvanico-db"
	galvaniconotification "galvanico/cmd/galvanico-notifications"
	galvanicoserver "galvanico/cmd/galvanico-server"
	galvanicowebsocket "galvanico/cmd/galvanico-websocket"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "galvanico",
	Short: "Game CLI",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		config.FileName = cmd.Flag("config").Value.String() //nolint:reassign // we need to set config

		if _, err := config.Load(); err != nil {
			return err
		}

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
