package cmd

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
)

var pushAll bool
var forcePush bool
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
			for _, secret := range sm.Secrets {
				if err := promptPushSecret(sm, secret); err != nil {
					return err
				}
			}
		} else {
			secret, err := findSecret(sm, args, "push")
			if err != nil {
				return err
			}
			if err := promptPushSecret(sm, secret); err != nil {
				return err
			}
		}
		return nil
	},
}

func promptPushSecret(sm *truss.SecretsManager, secret truss.SecretConfig) error {
	areSame, err := secretCompare(sm, secret, true)
	if err != nil {
		return err
	}
	if areSame && !forcePush {
		fmt.Println("No need to push, remote and local secrets are equal.")
		return nil
	}

	if forcePush || prompter.YesNo(fmt.Sprintf("Push to environment %s?", secret.Name()), false) {
		return sm.Push(secret)
	}
	return nil
}

func init() {
	secretsPushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push all environments")
	secretsPushCmd.Flags().BoolVarP(&forcePush, "force", "f", false, "Force push")
}
