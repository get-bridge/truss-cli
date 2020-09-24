package cmd

import (
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
)

var editPush bool
var secretsEditCmd = &cobra.Command{
	Use:   "edit [name] [kubeconfig] [-y]",
	Short: "Edits a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "edit")
		if err != nil {
			return err
		}

		saved, err := sm.Edit(secret)
		if err != nil {
			return err
		}
		if !saved {
			return nil
		}

		if editPush || prompter.YesNo("Push to environment "+secret.Name()+"?", false) {
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
	secretsEditCmd.Flags().BoolVarP(&editPush, "push", "y", false, "Push after editing, if saved")
}
