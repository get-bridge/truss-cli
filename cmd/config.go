package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
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
	kubeconfig, err := getKubeconfigName()
	if err != nil {
		return "", err
	}
	if kubeconfig == "" {
		return "", nil
	}

	kubeconfigDir, err := getKubeDir()
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

func getSSHPublicKeyPath() (string, error) {
	home, err := homedir.Dir()

	if err != nil {
		return "", errors.Wrap(err, "Unable to locate user's homedir")
	}

	viper.SetDefault("publicKeyPath", home+"/.ssh/id_rsa.pub")
	return viper.GetString("publicKeyPath"), nil
}

func getSSHPublicKey() (string, error) {
	publicKeyPath, err := getSSHPublicKeyPath()

	if err != nil {
		return "", err
	}

	publicKeyFile, err := ioutil.ReadFile(publicKeyPath)

	if err != nil {
		return "", errors.Wrap(err, "Unable to read public key from "+publicKeyPath)
	}

	return string(publicKeyFile), nil
}
