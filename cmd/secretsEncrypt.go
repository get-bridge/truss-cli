package cmd

import (
	"github.com/spf13/cobra"
)

var secretsEncryptCmd = &cobra.Command{
	Use:   "encrypt [name] [kubeconfig]",
	Short: "Encrypt a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "encrypt")
		if err != nil {
			return err
		}

		return sm.EncryptSecret(secret)
	},
}
