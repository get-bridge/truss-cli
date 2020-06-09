package cmd

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manages synchronizing secrets between Vault and the filesystem",
}

func findSecret(sm *truss.SecretsManager, args []string, verb string) (*truss.SecretConfig, error) {
	var name, kubeconfig string
	if len(args) >= 1 {
		name = args[0]
	} else {
		name = prompter.Choose(fmt.Sprintf("Which secret would you like to %s?", verb), sm.SecretNames(), "")
	}
	if len(args) >= 2 {
		kubeconfig = args[1]
	} else {
		kubeconfig = prompter.Choose("For which kubeconfig?", sm.SecretKubeconfigs(name), "")
	}

	return sm.Secret(name, kubeconfig)
}

func init() {
	secretsEditCmd.Flags().BoolVarP(&editPush, "push", "y", false, "Push after editing, if saved")
	secretsPushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push all environments")
	secretsPullCmd.Flags().BoolVarP(&pullAll, "all", "a", false, "Pull all environments")

	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsPushCmd)
	secretsCmd.AddCommand(secretsPullCmd)
	rootCmd.AddCommand(secretsCmd)
	viper.SetDefault("EDITOR", "vim")
}
