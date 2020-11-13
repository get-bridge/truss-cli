package cmd

import (
	"fmt"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretReceiveCmd = &cobra.Command{
	Use:   "receive {token}",
	Short: "Receives a shared secret by unwrapping a Vault Wrapped Token",
	Long:  secretShareCmd.Long,
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

		fmt.Print(string(output))
		return nil
	},
}

func init() {
	secretsCmd.AddCommand(secretReceiveCmd)
}
