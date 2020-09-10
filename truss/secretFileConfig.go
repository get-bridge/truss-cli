package truss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func init() {
	secretConfigFactories["file"] = parseSecretFileConfig
}

// SecretFileConfig represents a desired Vault synchronization
type SecretFileConfig struct {
	name       string `yaml:"name"`
	kubeconfig string `yaml:"kubeconfig"`
	vaultPath  string `yaml:"vaultPath"`
	filePath   string `yaml:"filePath"`
}

func parseSecretFileConfig(d map[string]interface{}) (SecretConfig, error) {
	config := SecretFileConfig{}
	if val, ok := d["name"]; ok {
		config.name = val.(string)
	}
	if val, ok := d["kubeconfig"]; ok {
		config.kubeconfig = val.(string)
	}
	if val, ok := d["vaultPath"]; ok {
		config.vaultPath = val.(string)
	}
	if val, ok := d["filePath"]; ok {
		config.filePath = val.(string)
	}

	return config, nil
}

// Name name
func (s SecretFileConfig) Name() string {
	return s.name
}

// Kubeconfig Kubeconfig
func (s SecretFileConfig) Kubeconfig() string {
	return s.kubeconfig
}

// VaultPath VaultPath
func (s SecretFileConfig) VaultPath() string {
	return s.vaultPath
}

// exists determines whether the local secrets file exists
func (s SecretFileConfig) existsOnDisk() bool {
	_, err := os.Stat(s.filePath)

	return err == nil
}

// getDecryptedFromDisk returns the decrypted yaml from disk
func (s SecretFileConfig) getDecryptedFromDisk(vault VaultCmd, transitKeyName string) ([]byte, error) {
	if !s.existsOnDisk() {
		return []byte("secrets: {}"), nil
	}

	encrypted, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	return vault.Decrypt(transitKeyName, encrypted)
}

func (s SecretFileConfig) getMapFromDisk(vault VaultCmd, transitKeyName string) (map[string]map[string]string, error) {
	raw, err := s.getDecryptedFromDisk(vault, transitKeyName)
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

// saveToDiskFromVault writes encrypted secrets to disk from vault
func (s SecretFileConfig) saveToDiskFromVault(vault VaultCmd, transitKeyName string) error {
	secretNames, err := vault.ListPath(s.vaultPath)
	if err != nil {
		return err
	}

	secrets := map[string]map[string]string{}
	for _, name := range secretNames {
		secret, err := vault.GetMap(path.Join(s.vaultPath, name))
		if err != nil {
			return err
		}
		secrets[name] = secret
	}

	out := map[string]map[string]map[string]string{
		"secrets": secrets,
	}

	y := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(y).Encode(out); err != nil {
		return err
	}

	return encryptAndSaveToDisk(vault, transitKeyName, s.filePath, y.Bytes())
}

// writeToVault writes a secret to Vault
func (s SecretFileConfig) writeToVault(vault VaultCmd, transitKeyName string) error {
	secrets, err := s.getMapFromDisk(vault, transitKeyName)
	if err != nil {
		return err
	}

	for key, data := range secrets {
		args := []string{"kv", "put", path.Join(s.vaultPath, key)}
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

func (s SecretFileConfig) saveToDisk(vault VaultCmd, transitKeyName string, raw []byte) error {
	return encryptAndSaveToDisk(vault, transitKeyName, s.filePath, raw)
}
