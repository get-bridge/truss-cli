package cmd

import (
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pushAll bool
var secretsPushCmd = &cobra.Command{
	Use:   "push [name] [kubeconfig] [-a]",
	Short: "Pushes a given environment's secrets to its corresponding Vault",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			return err
		}

		if pushAll {
			return sm.PushAll()
		}

		var secret *truss.SecretConfig
		secret, err = findSecret(sm, args, "push")
		if err != nil {
			return err
		}
		return sm.Push(*secret)
	},
}
