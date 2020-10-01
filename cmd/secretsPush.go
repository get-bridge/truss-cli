package cmd

import (
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
)

var pushAll bool
var secretsPushCmd = &cobra.Command{
	Use:   "push [name] [kubeconfig] [-a]",
	Short: "Pushes a given environment's secrets to its corresponding Vault",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		if pushAll {
			return sm.PushAll()
		}

		var secret truss.SecretConfig
		secret, err = findSecret(sm, args, "push")
		if err != nil {
			return err
		}
		return sm.Push(secret)
	},
}

func init() {
	secretsPushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push all environments")
}
