package galvaniconotification

import (
	"github.com/spf13/cobra"
)

// runCmd represents the serve command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
	},
}

func init() {
}
