package cmd

import (
	"fmt"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func getKubeconfigName() (string, error) {
	env := viper.GetString("env")
	if env == "" {
		return "", nil
	}

	environments := viper.GetStringMapString("environments")
	kubeconfig := environments[env]
	if kubeconfig == "" {
		var keys []string
		for k := range environments {
			keys = append(keys, k)
		}
		return "", fmt.Errorf("unknown env %v. Options: %v", env, keys)
	}
	return kubeconfig, nil
}

func getKubeconfig() (string, error) {
	var kubeconfigDir string
	kubeconfigDir, err := getKubeDir()
	if err != nil {
		return "", err
	}

	kubeconfig, err := getKubeconfigName()
	if err != nil {
		return "", err
	}
	return path.Join(kubeconfigDir, kubeconfig), nil
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
