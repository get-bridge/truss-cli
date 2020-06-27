package cmd

import (
	"github.com/spf13/cobra"
)

// getGlobalConfigCmd represents the getGlobalConfig command
var getGlobalConfigCmd = &cobra.Command{
	Use:   "get-global-config",
	Short: "Fetch global config",
}

func init() {
	rootCmd.AddCommand(getGlobalConfigCmd)
}
