package cmd

import (
	galvaniconotification "galvanico/cmd/galvanico-notifications"
	galvanicoserver "galvanico/cmd/galvanico-server"
	galvanicowebsocket "galvanico/cmd/galvanico-websocket"
	"galvanico/internal/config"
	"galvanico/internal/database"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "galvanico",
	Short: "Game CLI",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		log.Println("PreRun called")
		config.FileName = cmd.Flag("config").Value.String()

		if _, err := config.Load(); err != nil {
			return err
		}

		return database.Connection().Ping(cmd.Context())
	},
	PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
		log.Println("PostRun called")
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
}
