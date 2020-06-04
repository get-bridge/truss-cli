package truss

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// SecretConfigList represents a list of named SecretConfigs
type SecretConfigList struct {
	Environments   map[string]SecretConfig `yaml:"environments"`
	TransitKeyName string                  `yaml:"transit-key-name"`
}

// SecretConfigListFromFile reads a config file
func SecretConfigListFromFile(path string) (*SecretConfigList, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	list := &SecretConfigList{}
	if err := yaml.NewDecoder(r).Decode(list); err != nil {
		return nil, err
	}

	return list, nil
}

// Environment retreives an environment config from a config list
func (l SecretConfigList) Environment(name string) (*SecretConfig, error) {
	s, ok := l.Environments[name]
	if !ok {
		return nil, fmt.Errorf("secret configuration %s not found", name)
	}
	return &s, nil
}

// SecretConfig represents a desired Vault synchronization
type SecretConfig struct {
	Path   string `yaml:"path"`
	Secret string `yaml:"secret"`
}
