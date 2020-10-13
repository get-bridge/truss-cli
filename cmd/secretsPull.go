package cmd

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
)

var pullAll bool
var forcePull bool
var secretsPullCmd = &cobra.Command{
	Use:   "pull [name] [kubeconfig] [-a]",
	Short: "Pulls a given environment's secrets from its corresponding Vault",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		if pullAll {
			for _, secret := range sm.Secrets {
				if err := pullSecret(sm, secret); err != nil {
					return err
				}
			}
		} else {
			secret, err := findSecret(sm, args, "pull")
			if err != nil {
				return err
			}
			if err := pullSecret(sm, secret); err != nil {
				return err
			}
		}
		return nil
	},
}

func pullSecret(sm *truss.SecretsManager, secret truss.SecretConfig) error {
	areSame, err := secretCompare(sm, secret)
	if err != nil {
		return err
	}
	if areSame && !forcePull {
		fmt.Println("No need to pull, remote and local secrets are equal.")
		return nil
	}

	if forcePull || prompter.YesNo(fmt.Sprintf("Pull from environment %s?", secret.Name()), false) {
		return sm.Pull(secret)
	}
	return nil
}

func init() {
	secretsPullCmd.Flags().BoolVarP(&pullAll, "all", "a", false, "Pull all environments")
	secretsPullCmd.Flags().BoolVarP(&forcePull, "force", "f", false, "Force pull")
}
