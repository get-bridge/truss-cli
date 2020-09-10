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

var secretConfigFactories = map[string]func(map[string]interface{}) (SecretConfig, error){}

// SecretConfigListFromFile reads a config file
func SecretConfigListFromFile(path string) (*SecretConfigList, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	type UnparseSecretConfigList struct {
		Secrets        []map[string]interface{} `yaml:"secrets"`
		TransitKeyName string                   `yaml:"transit-key-name"`
	}

	unparsedList := &UnparseSecretConfigList{}
	if err := yaml.NewDecoder(r).Decode(unparsedList); err != nil {
		return nil, err
	}

	list := &SecretConfigList{
		TransitKeyName: unparsedList.TransitKeyName,
	}
	for _, s := range unparsedList.Secrets {
		secret, err := parseSecretConfig(s)
		if err != nil {
			return nil, err
		}
		list.Secrets = append(list.Secrets, secret)
	}

	return list, nil
}

func parseSecretConfig(s map[string]interface{}) (SecretConfig, error) {
	var secretType string
	secretTypeInterface, ok := s["type"]
	if !ok {
		secretType = "file"
	} else {
		secretType, ok = secretTypeInterface.(string)
		if !ok {
			return nil, fmt.Errorf("unknown secret type: %v", secretTypeInterface)
		}
	}
	factory, ok := secretConfigFactories[secretType]
	if !ok {
		return nil, fmt.Errorf("unknown secret type: %v", secretType)
	}
	return factory(s)
}

// Secret locates a secret by name and kubeconfig
func (l SecretConfigList) Secret(name, kubeconfig string) (SecretConfig, error) {
	for _, s := range l.Secrets {
		if s.Name() == name && s.Kubeconfig() == kubeconfig {
			return s, nil
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
			if s.Name() == n {
				dupe = true
			}
		}
		if dupe {
			continue
		}

		names = append(names, s.Name())
	}

	return names
}

// SecretKubeconfigs returns a list of kubeconfigs defined for a given secret
func (l SecretConfigList) SecretKubeconfigs(name string) []string {
	kubeconfigs := []string{}
	for _, s := range l.Secrets {
		if s.Name() == name {
			kubeconfigs = append(kubeconfigs, s.Kubeconfig())
		}
	}
	return kubeconfigs
}
