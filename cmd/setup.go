package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Ensure your truss cli is ready",
	// Long: `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called", args)
		fmt.Println("clusters : ", viper.Get("clusters"))
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
