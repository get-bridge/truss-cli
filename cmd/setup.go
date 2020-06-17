package cmd

import (
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
		return truss.Setup(&dependencies)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
