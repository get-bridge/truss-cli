package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Proxy commands to kubectl",
	// Long: `TODO...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubectl called")
	},
}

func init() {
	rootCmd.AddCommand(kubectlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubectlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubectlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
