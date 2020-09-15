package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Songmu/prompter"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manages synchronizing secrets between Vault and the filesystem",
}

var defaultSecretsFileName = "secrets.yaml"

func newSecretsManager() (*truss.SecretsManager, error) {
	secretsFile := viper.GetString("TRUSS_SECRETS_FILE")
	if secretsFile == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		foundSecrets := false
		for !foundSecrets {
			fileInfo, err := os.Stat(path.Join(dir, defaultSecretsFileName))
			if err != nil {
				dir = filepath.Dir(dir)
				if dir == "/" {
					return nil, errors.New("Could not find secrets.yaml")
				}
			} else {
				secretsFile = path.Join(dir, fileInfo.Name())
				foundSecrets = true
			}
		}
	}

	return truss.NewSecretsManager(secretsFile, viper.GetString("EDITOR"), getVaultAuth())
}

func findSecret(sm *truss.SecretsManager, args []string, verb string) (truss.SecretConfig, error) {
	var name, kubeconfig string
	if len(args) >= 1 {
		name = args[0]
	} else {
		name = prompter.Choose(fmt.Sprintf("Which secret would you like to %s?", verb), sm.SecretNames(), "")
	}

	var err error
	kubeconfig, err = getKubeconfigName()
	if err != nil {
		return nil, err
	}

	if len(args) >= 2 {
		if kubeconfig != "" {
			return nil, errors.New("do not specify --env and kubeconfig")
		}
		kubeconfig = args[1]
	} else if kubeconfig == "" {
		kubeconfigOptions := sm.SecretKubeconfigs(name)
		switch len(kubeconfigOptions) {
		case 0:
			return nil, errors.New("no kubeconfig found for secret")
		case 1:
			kubeconfig = kubeconfigOptions[0]
		default:
			kubeconfig = prompter.Choose("For which kubeconfig?", kubeconfigOptions, "")
		}
	}

	return sm.Secret(name, kubeconfig)
}

func init() {
	secretsEditCmd.Flags().BoolVarP(&editPush, "push", "y", false, "Push after editing, if saved")
	secretsPushCmd.Flags().BoolVarP(&pushAll, "all", "a", false, "Push all environments")
	secretsPullCmd.Flags().BoolVarP(&pullAll, "all", "a", false, "Pull all environments")

	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsEncryptCmd)
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsPushCmd)
	secretsCmd.AddCommand(secretsPullCmd)
	rootCmd.AddCommand(secretsCmd)
	viper.SetDefault("EDITOR", "vim")
}
