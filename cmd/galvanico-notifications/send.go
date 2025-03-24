package galvanico_notification

import (
	"fmt"
	"github.com/spf13/cobra"
)

// sendCmd represents the serve command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
	},
}

func init() {

}
