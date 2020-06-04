package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/Songmu/prompter"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manages synchronizing secrets between Vault and the filesystem",
}

var editPush bool
var secretsEditCmd = &cobra.Command{
	Use:   "edit <environment> [-y]",
	Short: "Edits a given environment's secrets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			log.Fatal(err)
		}

		saved, err := sm.Edit(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if !saved {
			return
		}

		if editPush || prompter.YesNo("Push to environment "+args[0]+"?", false) {
			secretsPushCmd.Run(cmd, args)
		}
	},
}

var secretsViewCmd = &cobra.Command{
	Use:   "view <environment>",
	Short: "Views a given environment's secrets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			log.Fatal(err)
		}

		vault, err := sm.Vault(args[0])
		if err != nil {
			log.Fatal(err)
		}

		out, err := sm.GetDecryptedFromDisk(vault, args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(out))
	},
}

var pushAll bool
var secretsPushCmd = &cobra.Command{
	Use:   "push [environment] [-a]",
	Short: "Pushes a given environment's secrets to its corresponding Vault",
	Args: func(cmd *cobra.Command, args []string) error {
		if !pushAll && len(args) != 1 {
			return errors.New("must specify an environment or --all")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			log.Fatal(err)
		}

		if pushAll {
			err = sm.PushAll()
		} else {
			err = sm.Push(args[0])
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	secretsEditCmd.Flags().BoolVarP(&editPush, "push", "y", false, "Push after editing, if saved")
	secretsPushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push all environments")

	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsPushCmd)
	rootCmd.AddCommand(secretsCmd)
	viper.SetDefault("EDITOR", "vim")
}
