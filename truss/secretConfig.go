package truss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// SecretConfig represents a desired Vault synchronization
type SecretConfig struct {
	Name       string `yaml:"name"`
	Kubeconfig string `yaml:"kubeconfig"`
	VaultPath  string `yaml:"vaultPath"`
	FilePath   string `yaml:"filePath"`

	transitKeyName string
}

// exists determines whether the local secrets file exists
func (s SecretConfig) existsOnDisk() bool {
	_, err := os.Stat(s.FilePath)

	return err == nil
}

// getDecryptedFromDisk returns the decrypted yaml from disk
func (s SecretConfig) getDecryptedFromDisk(vault VaultCmd) ([]byte, error) {
	if !s.existsOnDisk() {
		return []byte("secrets: {}"), nil
	}

	encrypted, err := ioutil.ReadFile(s.FilePath)
	if err != nil {
		return nil, err
	}

	return vault.Decrypt(s.transitKeyName, encrypted)
}

// getMapFromDisk returns a collection of secrets as a map
func (s SecretConfig) getMapFromDisk(vault VaultCmd) (map[string]map[string]string, error) {
	raw, err := s.getDecryptedFromDisk(vault)
	if err != nil {
		return nil, err
	}

	p := struct {
		Secrets map[string]map[string]string `yaml:"secrets"`
	}{}
	if err := yaml.NewDecoder(bytes.NewReader(raw)).Decode(&p); err != nil {
		return nil, err
	}

	return p.Secrets, nil
}

// encryptAndSaveToDisk encrypts and saves to disk
func (s SecretConfig) encryptAndSaveToDisk(vault VaultCmd, raw []byte) error {
	enc, err := vault.Encrypt(s.transitKeyName, raw)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.FilePath, enc, 0644)
}

// writeMapToDisk serializes a collection of secrets and writes them encrypted to disk
func (s SecretConfig) writeMapToDisk(vault VaultCmd, secrets map[string]map[string]string) error {
	out := map[string]map[string]map[string]string{
		"secrets": secrets,
	}

	y := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(y).Encode(out); err != nil {
		return err
	}

	enc, err := vault.Encrypt(s.transitKeyName, y.Bytes())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.FilePath, enc, 0644)
}

// Write writes a secret to Vault
func (s SecretConfig) write(vault VaultCmd) error {
	secrets, err := s.getMapFromDisk(vault)
	if err != nil {
		return err
	}

	for key, data := range secrets {
		args := []string{"kv", "put", path.Join(s.VaultPath, key)}
		for k, v := range data {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}

		_, err := vault.Run(args)
		if err != nil {
			return err
		}
	}

	return nil
}
