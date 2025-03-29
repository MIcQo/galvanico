package galvanicoserver

import (
	"fmt"
	"galvanico/internal/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const defaultPort = 8080

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, _ []string) {
		var srv = server.NewServer()
		var port = cmd.Flag("port").Value.String()

		go startServer(srv, port)

		c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
		signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

		<-c // This blocks the main thread until an interrupt is received
		log.Warn().Msg("Gracefully shutting down...")
		_ = srv.Shutdown()

		log.Info().Msg("Running cleanup tasks...")
	},
}

func init() {
	ServeCmd.PersistentFlags().IntP("port", "p", defaultPort, "port to serve on")
}

func startServer(srv *fiber.App, port string) {
	log.Info().Str("port", port).Msg("Starting server")

	err := srv.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
