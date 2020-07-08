package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Launch a shell container in a Truss cluster",
	Long:  "Launches a shell Pod in a Truss cluster. Useful for debugging purposes.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shell called")
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
