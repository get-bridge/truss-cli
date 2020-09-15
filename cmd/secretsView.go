package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var secretsViewCmd = &cobra.Command{
	Use:   "view [name] [kubeconfig]",
	Short: "Views a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "view")
		if err != nil {
			return err
		}

		secretString, err := sm.View(secret)
		if err != nil {
			return err
		}
		fmt.Println(secretString)
		return nil
	},
}
