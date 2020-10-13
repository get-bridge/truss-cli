package truss

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// SecretsManager syncrhonizes secrets between the filesystem and Vault
type SecretsManager struct {
	*SecretConfigList
	VaultAuth VaultAuth
	Editor    string
}

// NewSecretsManager creates a new SecretsManager
func NewSecretsManager(secretsFile string, editor string, vaultAuth VaultAuth) (*SecretsManager, error) {
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
func (m SecretsManager) Edit(secret SecretConfig) (bool, error) {
	// start port-forward
	vault, err := m.Vault(secret)
	if err != nil {
		return false, err
	}
	if _, err = vault.PortForward(); err != nil {
		return false, err
	}
	defer vault.ClosePortForward()

	var raw []byte
	if !secret.existsOnDisk() {
		// test that we can encrypt
		_, err = vault.Encrypt(m.TransitKeyName, []byte{})
		if err != nil {
			return false, err
		}
		raw = []byte("secrets: {}")
	} else {
		// load existing disk value
		raw, err = secret.getDecryptedFromDisk(vault, m.TransitKeyName)
		if err != nil {
			return false, err
		}
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
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", viper.GetString("EDITOR"), tmpFile.Name()))
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
	err = secret.saveToDisk(vault, m.TransitKeyName, raw)

	return true, err
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
func (m SecretsManager) Push(secret SecretConfig) error {
	vault, err := m.Vault(secret)
	if err != nil {
		return err
	}
	if _, err := vault.PortForward(); err != nil {
		return err
	}
	defer vault.ClosePortForward()

	return secret.writeToVault(vault, m.TransitKeyName)
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
func (m SecretsManager) Pull(secret SecretConfig) error {
	vault, err := m.Vault(secret)
	if err != nil {
		return err
	}
	if _, err := vault.PortForward(); err != nil {
		return err
	}
	defer vault.ClosePortForward()

	return secret.saveToDiskFromVault(vault, m.TransitKeyName)
}

// get kubeconfig name
func (m SecretsManager) kubeconfig(secret SecretConfig) (string, error) {
	config := viper.GetStringMap("kubeconfigfiles")
	directory, ok := config["directory"].(string)
	if !ok {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		directory = home + "/.kube/"
	}

	return path.Join(directory, secret.Kubeconfig()), nil
}

// Vault creates a proxied Vault client
func (m SecretsManager) Vault(secret SecretConfig) (*VaultCmd, error) {
	kubeconfig, err := m.kubeconfig(secret)
	if err != nil {
		return nil, err
	}

	return Vault(kubeconfig, m.VaultAuth), nil
}

// View Secret
func (m SecretsManager) View(secret SecretConfig) (localContent string, remoteContent string, err error) {
	vault, err := m.Vault(secret)
	if err != nil {
		return
	}

	local, err := secret.getDecryptedFromDisk(vault, m.TransitKeyName)
	if err != nil {
		return
	}
	remote, err := secret.getFromVault(vault)
	if err != nil {
		return
	}

	return string(local), string(remote), nil
}

// EncryptSecret on disk with cypher text from vault
func (m SecretsManager) EncryptSecret(secret SecretConfig) error {
	vault, err := m.Vault(secret)
	if err != nil {
		return err
	}
	if _, err := vault.PortForward(); err != nil {
		return err
	}
	defer vault.ClosePortForward()

	secretData, err := secret.getDecryptedFromDisk(vault, m.TransitKeyName)
	if err != nil {
		return err
	}

	return secret.saveToDisk(vault, m.TransitKeyName, secretData)
}
