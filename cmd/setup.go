package cmd

import (
	"fmt"
	"os"

	"github.com/instructure/truss-cli/truss"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Ensure your truss cli is ready",
	// Long: `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := truss.Setup(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
