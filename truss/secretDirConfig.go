package truss

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func init() {
	secretConfigFactories["dir"] = parseSecretDirConfig
}

// SecretDirConfig represents a desired Vault synchronization
type SecretDirConfig struct {
	name       string `yaml:"name"`
	kubeconfig string `yaml:"kubeconfig"`
	vaultPath  string `yaml:"vaultPath"`
	dirPath    string `yaml:"dirPath"`
}

func parseSecretDirConfig(d map[string]interface{}) (SecretConfig, error) {
	config := SecretDirConfig{}
	if val, ok := d["name"]; ok {
		config.name = val.(string)
	}
	if val, ok := d["kubeconfig"]; ok {
		config.kubeconfig = val.(string)
	}
	if val, ok := d["vaultPath"]; ok {
		config.vaultPath = val.(string)
	}
	if val, ok := d["dirPath"]; ok {
		config.dirPath = val.(string)
	}
	return config, nil
}

// Name name
func (s SecretDirConfig) Name() string {
	return s.name
}

// Kubeconfig Kubeconfig
func (s SecretDirConfig) Kubeconfig() string {
	return s.kubeconfig
}

// VaultPath VaultPath
func (s SecretDirConfig) VaultPath() string {
	return s.vaultPath
}

// exists determines whether the local secrets file exists
func (s SecretDirConfig) existsOnDisk() bool {
	_, err := os.Stat(s.dirPath)

	return err == nil
}

// getDecryptedFromDisk returns the decrypted yaml from disk
func (s SecretDirConfig) getDecryptedFromDisk(vault *VaultCmd, transitKeyName string) ([]byte, error) {
	data, err := s.getMapFromDisk(vault, transitKeyName)
	if err != nil {
		return nil, err
	}
	y := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(y).Encode(data); err != nil {
		return nil, err
	}
	return y.Bytes(), nil
}

// getMapFromDisk returns a collection of secrets as a map
func (s SecretDirConfig) getMapFromDisk(vault *VaultCmd, transitKeyName string) (map[string]string, error) {
	dirData := make(map[string]string)
	if !s.existsOnDisk() {
		return dirData, nil
	}
	err := filepath.Walk(s.dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		encrypted, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		bytes, err := vault.Decrypt(transitKeyName, encrypted)
		if err != nil {
			// if we fail to decrypt, might not be encypted
			dirData[info.Name()] = string(encrypted)
		} else {
			dirData[info.Name()] = string(bytes)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirData, nil
}

// writeMapToDisk serializes a collection of secrets and writes them encrypted to disk
func (s SecretDirConfig) saveToDiskFromVault(vault *VaultCmd, transitKeyName string) error {
	secrets, err := vault.GetMap(kv2DataPath(s.vaultPath))
	if err != nil {
		return err
	}

	secretStringMap := map[string]string{}
	for k, v := range secrets {
		vString, ok := v.(string)
		if !ok {
			return fmt.Errorf("failed to parse secrets: %v", secrets)
		}
		secretStringMap[k] = vString
	}

	for name, secretData := range secretStringMap {
		err := encryptAndSaveToDisk(vault, transitKeyName, path.Join(s.dirPath, name), []byte(secretData))
		if err != nil {
			return err
		}
	}

	return nil
}

// writeToVault writes a secret to Vault
func (s SecretDirConfig) writeToVault(vault *VaultCmd, transitKeyName string) error {
	secrets, err := s.getMapFromDisk(vault, transitKeyName)
	if err != nil {
		return err
	}

	_, err = vault.Write(kv2DataPath(s.vaultPath), map[string]interface{}{
		"data": secrets,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s SecretDirConfig) saveToDisk(vault *VaultCmd, transitKeyName string, raw []byte) error {
	secrets := map[string]string{}
	if err := yaml.NewDecoder(bytes.NewReader(raw)).Decode(&secrets); err != nil {
		return err
	}
	for name, data := range secrets {
		err := encryptAndSaveToDisk(vault, transitKeyName, path.Join(s.dirPath, name), []byte(data))
		if err != nil {
			return err
		}
	}
	return nil
}
