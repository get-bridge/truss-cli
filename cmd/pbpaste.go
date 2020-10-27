package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pbpasteCmd = &cobra.Command{
	Use:   "pbpaste {token}",
	Short: "Unwrapps a shared secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("TRUSS_ENV", "edge-cmh")
		kubeconfig, err := getKubeconfig()
		if err != nil {
			return err
		}

		output, err := truss.Vault(kubeconfig, getVaultAuth()).Run([]string{
			"unwrap",
			"-field=pb",
			args[0],
		})
		if err != nil {
			return err
		}

		copy := fmt.Sprintf("\"%s\" has been copied to your clipboard!", output)
		fmt.Println(copy)
		return clipboard.WriteAll(copy)
	},
}

func init() {
	rootCmd.AddCommand(pbpasteCmd)
}
