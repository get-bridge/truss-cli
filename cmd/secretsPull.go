package cmd

import (
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pullAll bool
var secretsPullCmd = &cobra.Command{
	Use:   "pull [name] [kubeconfig] [-a]",
	Short: "Pulls a given environment's secrets from its corresponding Vault",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			return err
		}

		if pullAll {
			return sm.PullAll()
		}

		var secret *truss.SecretConfig
		secret, err = findSecret(sm, args, "pull")
		if err != nil {
			return err
		}
		return sm.Pull(*secret)
	},
}
