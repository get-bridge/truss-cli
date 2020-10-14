package cmd

import (
	"fmt"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "A wrapper around hashicorp vault",
	Long: `This is useful when your vault is not exposed publicly.
As it will port-forward to the service and call aws auth`,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, err := getKubeconfig()
		if err != nil {
			return err
		}

		output, err := truss.Vault(kubeconfig, getVaultAuth()).Run(args)
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

func getVaultAuth() truss.VaultAuth {
	vaultRole := viper.GetString("vault.auth.aws.vaultrole")
	awsRole := viper.GetString("vault.auth.aws.awsrole")
	awsRegion := viper.GetString("vault.auth.aws.awsregion")

	if vaultRole == "" || awsRole == "" {
		return nil
	}

	return truss.VaultAuthAWS(vaultRole, awsRole, awsRegion)
}

func init() {
	rootCmd.AddCommand(vaultCmd)
	vaultCmd.Flags().SetInterspersed(false)

	viper.SetDefault("vault.auth.aws.awsregion", "us-east-1")
}
