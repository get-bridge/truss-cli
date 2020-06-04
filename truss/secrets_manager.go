package truss

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
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
	l, err := SecretConfigListFromFile("./secrets.yaml")
	if err != nil {
		return nil, err
	}

	return &SecretsManager{
		SecretConfigList: l,
		Editor:           editor,
		VaultAuth:        vaultAuth,
	}, nil
}

// Kubectl creates a Kubectl client
func (m SecretsManager) Kubectl(environment string) (*KubectlCmd, error) {
	config := viper.GetStringMap("kubeconfigfiles")
	directory, ok := config["directory"].(string)
	if !ok {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		directory = home + "/.kube/"
	}

	environments := viper.GetStringMapString("environments")

	return Kubectl(path.Join(directory, environments[environment])), nil
}

// Vault creates a proxied Vault client
func (m SecretsManager) Vault(environment string) (*VaultCmd, error) {
	kubectl, err := m.Kubectl(environment)
	if err != nil {
		return nil, err
	}

	return Vault(kubectl, m.VaultAuth), nil
}

// Edit edits an environments's ecret
func (m SecretsManager) Edit(environment string) error {
	// load existing disk value
	// decrypt it or provide default
	raw, err := m.GetDecryptedFromDisk(environment)
	if err != nil {
		log.Fatal(err)
	}

	// save to tmp file
	tmpFile, err := ioutil.TempFile("", "trussvault-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write(raw)
	tmpFile.Close()

	// vim tmp file
	cmd := exec.Command(viper.GetString("EDITOR"), tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	// TODO: check if saved
	// encrypt and save to disk
	raw, err = ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}

	if err := m.EncryptAndSaveToDisk(environment, raw); err != nil {
		return err
	}

	// prompt to push

	return nil
}

// GetDecryptedFromDisk returns the encrypted yaml configuration from disk
func (m SecretsManager) GetDecryptedFromDisk(environment string) ([]byte, error) {
	e, err := m.Environment(environment)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(e.Secret)
	if err != nil {
		return []byte("secrets: {}"), nil
	}

	encrypted, err := ioutil.ReadFile(e.Secret)
	if err != nil {
		return nil, err
	}

	return m.Decrypt(environment, encrypted)
}

// GetMapFromDisk returns a collection of secrets as a map
func (m SecretsManager) GetMapFromDisk(environment string) (map[string]map[string]interface{}, error) {
	raw, err := m.GetDecryptedFromDisk(environment)
	if err != nil {
		return nil, err
	}

	p := struct {
		Secrets map[string]map[string]interface{} `yaml:"secrets"`
	}{}
	if err := yaml.NewDecoder(bytes.NewReader(raw)).Decode(&p); err != nil {
		return nil, err
	}

	return p.Secrets, nil
}

// EncryptAndSaveToDisk encrypts and saves to disk
func (m SecretsManager) EncryptAndSaveToDisk(environment string, raw []byte) error {
	e, err := m.Environment(environment)
	if err != nil {
		return err
	}

	enc, err := m.Encrypt(environment, raw)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(e.Secret, enc, 0644)
}

// Decrypt shit
func (m SecretsManager) Decrypt(environment string, encrypted []byte) ([]byte, error) {
	vault, err := m.Vault(environment)
	if err != nil {
		return nil, err
	}

	out, err := vault.Run([]string{
		"write",
		"-field=plaintext",
		"transit/decrypt/" + m.TransitKeyName,
		"ciphertext=" + string(encrypted),
	})

	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(out))
}

// Encrypt shit
func (m SecretsManager) Encrypt(environment string, raw []byte) ([]byte, error) {
	vault, err := m.Vault(environment)
	if err != nil {
		return nil, err
	}

	out, err := vault.Run([]string{
		"write",
		"-field=ciphertext",
		"transit/encrypt/" + m.TransitKeyName,
		"plaintext=" + base64.StdEncoding.EncodeToString(raw),
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Push pushes secrets to Vaut
func (m SecretsManager) Push(environment string) error {

	secrets, err := m.GetMapFromDisk(environment)
	if err != nil {
		return err
	}

	for path, data := range secrets {
		m.Write(environment, path, data)
	}
	return nil
}

// Write writes a secret to Vault
func (m SecretsManager) Write(environment, dst string, data map[string]interface{}) error {
	e, err := m.Environment(environment)
	if err != nil {
		return err
	}

	vault, err := m.Vault(environment)
	if err != nil {
		return err
	}

	args := []string{"kv", "put", path.Join(e.Path, dst)}
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	_, err = vault.Run(args)
	return err
}

// PushAll pushes all secrets for all environments
func (m SecretsManager) PushAll() error {
	for env := range m.Environments {
		if err := m.Push(env); err != nil {
			return err
		}
	}
	return nil
}
