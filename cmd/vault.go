package cmd

import (
	"os"

	"github.com/instructure-bridge/truss-cli/truss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "A wrapper around hashicorp vault",
	Long: `This is useful when your vault is not exposed publicly.
As it will port-forward to the service and call aws auth`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, err := getKubeconfig(cmd, args)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}

		kubectl := truss.Kubectl(kubeconfig)
		output, err := truss.Vault(kubectl, getVaultAuth()).Run(args)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Println(string(output))
	},
}

func getVaultAuth() truss.VaultAuth {
	vaultRole := viper.GetString("vault.auth.aws.vaultrole")
	awsRole := viper.GetString("vault.auth.aws.awsrole")

	if vaultRole == "" || awsRole == "" {
		return nil
	}

	return truss.VaultAuthAWS(vaultRole, awsRole)
}

func init() {
	rootCmd.AddCommand(vaultCmd)
}
