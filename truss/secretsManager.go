package truss

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// SecretsManager syncrhonizes secrets between the filesystem and Vault
type SecretsManager struct {
	*SecretConfigList
	VaultAuth VaultAuth
	Editor    string
}

// NewSecretsManager creates a new SecretsManager
func NewSecretsManager(editor string, vaultAuth VaultAuth) (*SecretsManager, error) {
	secretsFile := os.Getenv("TRUSS_SECRETS_FILE")
	if secretsFile == "" {
		secretsFile = "./secrets.yaml"
	}
	l, err := SecretConfigListFromFile(secretsFile)
	if err != nil {
		return nil, err
	}

	return &SecretsManager{
		SecretConfigList: l,
		Editor:           editor,
		VaultAuth:        vaultAuth,
	}, nil
}

// Edit edits an environments's secrets
// Returns true if $EDITOR wrote to the temp file
func (m SecretsManager) Edit(secret *SecretConfig) (bool, error) {
	// start port-forward
	vault, err := m.vault(secret)
	if err != nil {
		return false, err
	}
	if _, err = vault.PortForward(); err != nil {
		return false, err
	}
	defer vault.ClosePortForward()

	// load existing disk value
	// decrypt it or provide default
	raw, err := secret.getDecryptedFromDisk(vault)
	if err != nil {
		return false, err
	}

	// save to tmp file
	tmpFile, err := ioutil.TempFile("", "trussvault-*")
	if err != nil {
		return false, err
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write(raw)
	tmpFile.Close()

	info, _ := os.Stat(tmpFile.Name())
	modTimeAtOpen := info.ModTime()

	// vim tmp file
	cmd := exec.Command(viper.GetString("EDITOR"), tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	// check if saved
	info, _ = os.Stat(tmpFile.Name())
	modTimeAtClose := info.ModTime()
	if !modTimeAtClose.After(modTimeAtOpen) {
		return false, nil
	}

	// encrypt and save to disk
	raw, err = ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return true, err
	}
	if err := secret.encryptAndSaveToDisk(vault, raw); err != nil {
		return true, err
	}

	return true, nil
}

// PushAll pushes all secrets for all environments
func (m SecretsManager) PushAll() error {
	for _, secret := range m.Secrets {
		if err := m.Push(secret); err != nil {
			return err
		}
	}
	return nil
}

// Push pushes secrets to Vaut
func (m SecretsManager) Push(secret *SecretConfig) error {
	vault, err := m.vault(secret)
	if err != nil {
		return err
	}
	if _, err := vault.PortForward(); err != nil {
		return err
	}
	defer vault.ClosePortForward()

	return secret.write(vault)
}

// PullAll pulls all environments
func (m SecretsManager) PullAll() error {
	for _, secret := range m.Secrets {
		if err := m.Pull(secret); err != nil {
			return err
		}
	}
	return nil
}

// Pull updates the file on disk with the vaules from Vault (destructive)
func (m SecretsManager) Pull(secret *SecretConfig) error {
	vault, err := m.vault(secret)
	if err != nil {
		return err
	}
	if _, err := vault.PortForward(); err != nil {
		return err
	}
	defer vault.ClosePortForward()

	p, err := m.getMapFromVault(vault, secret)
	if err != nil {
		return err
	}

	return secret.writeMapToDisk(vault, p)
}

// kubectl creates a Kubectl client
func (m SecretsManager) kubectl(secret *SecretConfig) (*KubectlCmd, error) {
	config := viper.GetStringMap("kubeconfigfiles")
	directory, ok := config["directory"].(string)
	if !ok {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		directory = home + "/.kube/"
	}

	return Kubectl(path.Join(directory, secret.Kubeconfig)), nil
}

// vault creates a proxied Vault client
func (m SecretsManager) vault(secret *SecretConfig) (VaultCmd, error) {
	kubectl, err := m.kubectl(secret)
	if err != nil {
		return nil, err
	}

	return Vault(kubectl, m.VaultAuth), nil
}

// getMapFromVault returns a collection of secrets as a map
func (m SecretsManager) getMapFromVault(vault VaultCmd, secret *SecretConfig) (map[string]map[string]string, error) {
	out := map[string]map[string]string{}

	list, err := vault.Run([]string{
		"kv",
		"list",
		"-format=yaml",
		secret.VaultPath,
	})
	if err != nil {
		return nil, err
	}

	secrets := []string{}
	if err := yaml.NewDecoder(bytes.NewReader(list)).Decode(&secrets); err != nil {
		return nil, err
	}

	for _, s := range secrets {
		get, err := vault.Run([]string{
			"kv",
			"get",
			"-format=yaml",
			path.Join(secret.VaultPath, s),
		})
		if err != nil {
			return nil, err
		}

		getData := struct {
			Data struct {
				Data map[string]string `yaml:"data"`
			} `yaml:"data"`
		}{}
		if err := yaml.NewDecoder(bytes.NewReader(get)).Decode(&getData); err != nil {
			return nil, err
		}

		out[s] = getData.Data.Data
	}

	return out, nil
}

// View Secret
func (m SecretsManager) View(secret *SecretConfig) (string, error) {
	if !secret.existsOnDisk() {
		return "", errors.New("no such local secrets file exists. try running truss secrets pull")
	}

	vault, err := m.vault(secret)
	if err != nil {
		return "", err
	}

	out, err := secret.getDecryptedFromDisk(vault)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
