package truss

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// ErrSecretFileConfigInvalidYaml error if invalid yaml
var ErrSecretFileConfigInvalidYaml = errors.New("Unable to parse secret as yaml or missing required root element `secrets`")

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
func (s SecretFileConfig) getDecryptedFromDisk(vault *VaultCmd, transitKeyName string) ([]byte, error) {
	encrypted, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	ot, err := NewObfuscationTarget(bytes.NewReader(encrypted))
	if err != nil {
		// It's possibly a fully-encrypted file. Try to decrypt the whole thing?
		return vault.Decrypt(transitKeyName, encrypted)
	}

	if err := ot.Decrypt(vault, transitKeyName); err != nil {
		return nil, err
	}

	out := bytes.NewBuffer(nil)
	yaml.NewEncoder(out).Encode(ot)
	return out.Bytes(), nil
}

func (s SecretFileConfig) getFromVault(vault *VaultCmd) ([]byte, error) {
	secrets, err := s.getMapFromVault(vault)
	if err != nil {
		return nil, err
	}
	out := map[string]map[string]map[string]string{
		"secrets": secrets,
	}

	y := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(y).Encode(out); err != nil {
		return nil, err
	}
	return y.Bytes(), nil
}

func (s SecretFileConfig) getMapFromDisk(vault *VaultCmd, transitKeyName string) (map[string]map[string]string, error) {
	raw, err := s.getDecryptedFromDisk(vault, transitKeyName)
	if err != nil {
		return nil, err
	}

	return parseSecretFileYaml(raw)
}

func (s SecretFileConfig) getMapFromVault(vault *VaultCmd) (map[string]map[string]string, error) {
	secretNames, err := vault.ListPath(kv2MetadataPath(s.vaultPath))
	if err != nil {
		return nil, err
	}

	secrets := map[string]map[string]string{}
	for _, name := range secretNames {
		vaultPath := kv2DataPath(path.Join(s.vaultPath, name))
		secret, err := vault.GetMap(vaultPath)
		if err != nil {
			return nil, err
		}

		secretStringMap := map[string]string{}
		for k, v := range secret {
			vString, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("failed to parse secret: %v", secret)
			}
			secretStringMap[k] = vString
		}

		secrets[name] = secretStringMap
	}
	return secrets, nil
}

// saveToDiskFromVault writes encrypted secrets to disk from vault
func (s SecretFileConfig) saveToDiskFromVault(vault *VaultCmd, transitKeyName string) error {
	v, err := s.getFromVault(vault)
	if err != nil {
		return err
	}

	return s.saveToDisk(vault, transitKeyName, v)
}

// writeToVault writes a secret to Vault
func (s SecretFileConfig) writeToVault(vault *VaultCmd, transitKeyName string) error {
	secrets, err := s.getMapFromDisk(vault, transitKeyName)
	if err != nil {
		return err
	}

	for key, data := range secrets {
		vaultPath := kv2DataPath(path.Join(s.vaultPath, key))
		_, err := vault.Write(vaultPath, map[string]interface{}{
			"data": data,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SecretFileConfig) saveToDisk(vault *VaultCmd, transitKeyName string, raw []byte) error {
	ot, err := NewObfuscationTarget(bytes.NewReader(raw))
	if err != nil {
		return err
	}

	if err := ot.Encrypt(vault, transitKeyName); err != nil {
		return err
	}

	// ensure dir exists
	if err := os.MkdirAll(path.Dir(s.filePath), 0744); err != nil {
		return err
	}
	f, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewEncoder(f).Encode(ot)
}

func parseSecretFileYaml(raw []byte) (map[string]map[string]string, error) {
	p := struct {
		Secrets map[string]map[string]string `yaml:"secrets"`
	}{}
	decoder := yaml.NewDecoder(bytes.NewReader(raw))
	decoder.SetStrict(true)
	if err := decoder.Decode(&p); err != nil {
		return nil, ErrSecretFileConfigInvalidYaml
	}
	return p.Secrets, nil
}
