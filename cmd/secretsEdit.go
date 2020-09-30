package cmd

import (
	"fmt"

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
			fmt.Println("secrets not saved.")
			return nil
		}

		promptPushSecret(sm, secret)
		return nil
	},
}

func init() {
	secretsEditCmd.Flags().BoolVarP(&editPush, "push", "y", false, "Push after editing, if saved")
}
