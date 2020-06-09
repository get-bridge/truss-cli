package cmd

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getKubeconfig(cmd *cobra.Command, args []string) (kubeconfig string, err error) {
	env := viper.GetString("env")

	if env != "" {
		var kubeconfigDir string
		kubeconfigDir, err = getKubeDir()
		if err != nil {
			return
		}

		environments := viper.GetStringMapString("environments")
		kubeconfig = path.Join(kubeconfigDir, environments[env])
	}
	return
}

func getKubeDir() (string, error) {
	config := viper.GetStringMap("kubeconfigfiles")

	directory, ok := config["directory"].(string)
	if !ok {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		directory = home + "/.kube/"
	}
	return directory, nil
}
