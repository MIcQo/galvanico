package galvanicowebsocket

import (
	"github.com/spf13/cobra"
)

const defaultPort = 8082

// runCmd represents the serve command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
	},
}

func init() {
	runCmd.PersistentFlags().IntP("port", "p", defaultPort, "port to serve on")
}
