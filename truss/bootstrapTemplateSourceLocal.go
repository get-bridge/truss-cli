package truss

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// LocalTemplateSource is a Local Filesystem Template Source
type LocalTemplateSource struct {
	Directory string
}

// NewLocalTemplateSource returns a new TemplateSource
func NewLocalTemplateSource(directory string) *LocalTemplateSource {
	return &LocalTemplateSource{directory}
}

// ListTemplates returns a list of temlpates in the local directory
func (s LocalTemplateSource) ListTemplates() ([]string, error) {
	l, err := ioutil.ReadDir(s.Directory)
	if err != nil {
		return nil, err
	}

	out := []string{}
	for _, t := range l {
		if t.IsDir() {
			out = append(out, t.Name())
		}
	}

	return out, nil
}

// LocalDirectory returns a local directory for the template
func (s LocalTemplateSource) LocalDirectory(template string) string {
	return filepath.Join(s.Directory, template)
}

// GetTemplateManifest gets the template's manifest
func (s LocalTemplateSource) GetTemplateManifest(t string) *BootstrapManifest {
	f, err := os.Open(filepath.Join(s.Directory, t, ".truss-manifest.yaml"))
	if err != nil {
		return nil
	}
	defer f.Close()

	m := &BootstrapManifest{}
	if err := yaml.NewDecoder(f).Decode(m); err != nil {
		return nil
	}

	return m
}
