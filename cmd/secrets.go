package cmd

import (
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
	Use:   "edit [name] [kubeconfig] [-y]",
	Short: "Edits a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			log.Fatal(err)
		}

		secret, err := findSecret(sm, args, "edit")
		if err != nil {
			log.Fatal(err)
		}

		saved, err := sm.Edit(*secret)
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
	Use:   "view [name] [kubeconfig]",
	Short: "Views a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			log.Fatal(err)
		}

		secret, err := findSecret(sm, args, "view")
		if err != nil {
			log.Fatal(err)
		}

		vault, err := sm.Vault(*secret)
		if err != nil {
			log.Fatal(err)
		}

		out, err := sm.GetDecryptedFromDisk(vault, *secret)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(out))
	},
}

var pushAll bool
var secretsPushCmd = &cobra.Command{
	Use:   "push [name] [kubeconfig] [-a]",
	Short: "Pushes a given environment's secrets to its corresponding Vault",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			log.Fatal(err)
		}

		if pushAll {
			err = sm.PushAll()
		} else {
			var secret *truss.SecretConfig
			secret, err = findSecret(sm, args, "push")
			if err != nil {
				log.Fatal(err)
			}
			err = sm.Push(*secret)
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

var pullAll bool
var secretsPullCmd = &cobra.Command{
	Use:   "pull [name] [kubeconfig] [-a]",
	Short: "Pulls a given environment's secrets from its corresponding Vault",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := truss.NewSecretsManager(viper.GetString("EDITOR"), getVaultAuth())
		if err != nil {
			return err
		}

		if pullAll {
			err = sm.PullAll()
		} else {
			var secret *truss.SecretConfig
			secret, err = findSecret(sm, args, "push")
			if err != nil {
				log.Fatal(err)
			}
			err = sm.Pull(*secret)
		}

		if err != nil {
			return err
		}
		return nil
	},
}

func findSecret(sm *truss.SecretsManager, args []string, verb string) (*truss.SecretConfig, error) {
	var name, kubeconfig string
	if len(args) >= 1 {
		name = args[0]
	} else {
		name = prompter.Choose(fmt.Sprintf("Which secret would you like to %s?", verb), sm.SecretNames(), "")
	}
	if len(args) >= 2 {
		kubeconfig = args[1]
	} else {
		kubeconfig = prompter.Choose("For which kubeconfig?", sm.SecretKubeconfigs(name), "")
	}

	return sm.Secret(name, kubeconfig)
}

func init() {
	secretsEditCmd.Flags().BoolVarP(&editPush, "push", "y", false, "Push after editing, if saved")
	secretsPushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push all environments")
	secretsPullCmd.Flags().BoolVarP(&pullAll, "all", "a", false, "Pull all environments")

	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsPushCmd)
	secretsCmd.AddCommand(secretsPullCmd)
	rootCmd.AddCommand(secretsCmd)
	viper.SetDefault("EDITOR", "vim")
}
