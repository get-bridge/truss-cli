package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func getKubeconfigName() (string, error) {
	env := viper.GetString("TRUSS_ENV")
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

func getKubeconfigStruct() (*clientcmdapi.Config, error) {
	kc, err := getKubeconfig()
	if err != nil {
		return nil, err
	}
	kcb, err := ioutil.ReadFile(kc)
	if err != nil {
		return nil, err
	}

	cc, err := clientcmd.NewClientConfigFromBytes(kcb)
	if err != nil {
		return nil, err
	}

	c, err := cc.RawConfig()
	return &c, err
}

func must(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}

func envClusterName() (string, error) {
	return getAuthInfoArg("cluster-name")
}

func envClusterRegion() (string, error) {
	return getAuthInfoArg("region")
}

func envClusterRoleArn() (string, error) {
	return getAuthInfoArg("role")
}

func getAuthInfoArg(arg string) (string, error) {
	kc, err := getKubeconfigStruct()
	if err != nil {
		return "", err
	}

	var auth *clientcmdapi.AuthInfo
	for _, a := range kc.AuthInfos {
		auth = a
		break
	}

	for k, v := range auth.Exec.Args {
		if v == fmt.Sprintf("--%s", arg) {
			return auth.Exec.Args[k+1], nil
		}
	}
	return "", nil
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
