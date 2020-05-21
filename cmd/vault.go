package cmd

import (
	"os"

	"github.com/instructure/truss-cli/truss"
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

		var vaultAuth truss.VaultAuth

		vaultConfig := viper.GetStringMap("vault")
		authConfig, ok := vaultConfig["auth"].(map[string]interface{})
		if ok {
			awsConfig, ok := authConfig["aws"].(map[string]interface{})
			if ok {
				vaultRole := awsConfig["vaultrole"].(string)
				awsRole := awsConfig["awsrole"].(string)
				vaultAuth = truss.VaultAuthAWS(vaultRole, awsRole)
			}
		}

		kubectl := truss.Kubectl(kubeconfig)
		output, err := truss.Vault(kubectl, vaultAuth).Run(args)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
}
