package truss

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// SecretConfigList represents a list of named SecretConfigs
type SecretConfigList struct {
	Secrets        []SecretConfig `yaml:"secrets"`
	TransitKeyName string         `yaml:"transit-key-name"`
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

// Secret locates a secret by name and kubeconfig
func (l SecretConfigList) Secret(name, kubeconfig string) (*SecretConfig, error) {
	for _, s := range l.Secrets {
		if s.Name == name && s.Kubeconfig == kubeconfig {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("secret named '%s' in '%s' not found", name, kubeconfig)
}

// SecretNames returns a list of secret names provided in the list
func (l SecretConfigList) SecretNames() []string {
	names := []string{}
	for _, s := range l.Secrets {
		var dupe bool
		for _, n := range names {
			if s.Name == n {
				dupe = true
			}
		}
		if dupe {
			continue
		}

		names = append(names, s.Name)
	}

	return names
}

// SecretKubeconfigs returns a list of kubeconfigs defined for a given secret
func (l SecretConfigList) SecretKubeconfigs(name string) []string {
	kubeconfigs := []string{}
	for _, s := range l.Secrets {
		if s.Name == name {
			kubeconfigs = append(kubeconfigs, s.Kubeconfig)
		}
	}
	return kubeconfigs
}
