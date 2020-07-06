package cmd

import (
	"errors"
	"fmt"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsViewCmd = &cobra.Command{
	Use:   "view [name] [kubeconfig]",
	Short: "Views a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "view")
		if err != nil {
			return err
		}

		if !sm.Exists(*secret) {
			return errors.New("no such local secrets file exists. try running truss secrets pull")
		}

		vault, err := sm.Vault(*secret)
		if err != nil {
			return err
		}

		out, err := sm.GetDecryptedFromDisk(vault, *secret)
		if err != nil {
			return err
		}

		fmt.Print(string(out))
		return nil
	},
}
