package cmd

import (
	"github.com/Songmu/prompter"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editPush bool
var secretsEditCmd = &cobra.Command{
	Use:   "edit [name] [kubeconfig] [-y]",
	Short: "Edits a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "edit")
		if err != nil {
			return err
		}

		saved, err := sm.Edit(*secret)
		if err != nil {
			return err
		}
		if !saved {
			return nil
		}

		if editPush || prompter.YesNo("Push to environment "+ secret.Name +"?", false) {
			if len(args) > 0 {
				return secretsPushCmd.RunE(cmd, args)
			} else {
				newArgs := []string{secret.Name}
				return secretsPushCmd.RunE(cmd, newArgs)
			}
		}
		return nil
	},
}
