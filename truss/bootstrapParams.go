package truss

import (
	"errors"
	"fmt"
	"strings"
)

// BootstrapParams represents user-provided parameters for a template
type BootstrapParams map[string]interface{}

// LoadFromConfig reads params from a given configuration
func (p *BootstrapParams) LoadFromConfig(c *BootstrapConfig) error {
	for k, v := range c.Params {
		(*p)[k] = v
	}
	return nil
}

// LoadFromFlags reads params from a collection of flag values
func (p *BootstrapParams) LoadFromFlags(s map[string]string) error {
	for k, v := range s {
		switch v {
		case "true":
			(*p)[k] = true
			break
		case "false":
			(*p)[k] = false
			break
		default:
			(*p)[k] = v
		}
	}
	return nil
}

// Validate validates the given parameters against a manifest
func (p BootstrapParams) Validate(m *BootstrapManifest) (errs []string, err error) {
	for _, param := range m.Params {
		value := p[param.Name]

		if p == nil {
			errs = append(errs, fmt.Sprintf("missing required param [%s]", param.Name))
			continue
		}

		switch strings.ToLower(param.Type) {
		case "bool":
			if _, ok := value.(bool); !ok {
				errs = append(errs, fmt.Sprintf("param [%s] must be a boolean", param.Name))
				continue
			}
		case "string":
			if _, ok := value.(string); !ok {
				errs = append(errs, fmt.Sprintf("param [%s] must be a string", param.Name))
				continue
			}
		}
	}

	if len(errs) > 0 {
		err = errors.New("invalid params")
	}

	return
}
