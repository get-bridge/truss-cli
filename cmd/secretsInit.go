package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a given service's secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		var service string
		for service == "" {
			service = prompter.Prompt("Service name", "")
		}
		envSecretsPath := prompter.Prompt("Directory to save secrets", "deploy")
		envSecretFileName := prompter.Prompt("Filename to save secrets", "secrets")
		vaultPath := prompter.Prompt("Vault path to save secrets", "secret")
		fileName := prompter.Prompt("File name", "secrets.yaml")

		environments := viper.GetStringMapString("environments")

		fileData := generateSecretsFile(environments, service, envSecretsPath, envSecretFileName, vaultPath, fileName)
		return ioutil.WriteFile(fileName, []byte(fileData), 0644)
	},
}

func init() {
	secretsCmd.AddCommand(secretsInitCmd)
}

func generateSecretsFile(environments map[string]string, service, envSecretsPath, envSecretFileName, vaultPath, fileName string) string {
	envSecrets := ""
	for e, kubeconfig := range environments {
		path := strings.ReplaceAll(e, "-", "/")
		envSecrets = envSecrets + fmt.Sprintf(`- name: %s
  kubeconfig: %s
  vaultPath: %s/%s/%s
  filePath: %s/%s/%s
`,
			e,
			kubeconfig,
			envSecretsPath, path, envSecretFileName,
			vaultPath, path, service,
		)
	}

	return fmt.Sprintf(`# %s
transit-key-name: %s

secrets:
%s`, fileName, service, envSecrets)
}
