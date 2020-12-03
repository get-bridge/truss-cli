package cmd

import (
	"fmt"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Ensure your truss cli is ready",
	Long: `Dependencies are configured using 'dependencies' field in configfile.

dependencies:
- kubectl
- sshuttle
- vault
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dependencies := viper.GetStringSlice("dependencies")

		err := truss.Setup(&dependencies)
		if err == nil {
			fmt.Println("No problems detected")
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
