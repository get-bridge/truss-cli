package cmd

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
)

var encryptPush bool
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

		if err = sm.EncryptSecret(secret); err != nil {
			return err
		}

		if encryptPush || prompter.YesNo(fmt.Sprintf("Push to environment %s?", secret.Name()), false) {
			if len(args) > 0 {
				return secretsPushCmd.RunE(cmd, args)
			}
			newArgs := []string{secret.Name()}
			return secretsPushCmd.RunE(cmd, newArgs)
		}
		return nil
	},
}

func init() {
	secretsEncryptCmd.Flags().BoolVarP(&encryptPush, "push", "y", false, "Push after encrypt")
}
