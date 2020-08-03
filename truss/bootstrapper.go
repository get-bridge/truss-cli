package truss

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Bootstrapper bootstraps a deployment
type Bootstrapper struct {
	TemplateSource
	TrussDir string
	Template string
}

// BootstrapManifest represents the manifest thingy
type BootstrapManifest struct {
	Params []struct {
		Name    string      `yaml:"name"`
		Type    string      `yaml:"type"`
		Prompt  string      `yaml:"prompt"`
		Default interface{} `yaml:"default"`
	} `yaml:"params"`
}

// TemplateSource sources templates
type TemplateSource interface {
	ListTemplates() ([]string, error)
	LocalDirectory(template string) (string, error)
	GetTemplateManifest(t string) *BootstrapManifest
	Cleanup()
}

// NewBootstrapper returns a new TemplateSource
func NewBootstrapper(ts TemplateSource, trussDir, template string) *Bootstrapper {
	return &Bootstrapper{ts, trussDir, template}
}

// GetTemplateManifest gets a template's manifest
func (b Bootstrapper) GetTemplateManifest() *BootstrapManifest {
	return b.TemplateSource.GetTemplateManifest(b.Template)
}

// Bootstrap does the thing!
func (b Bootstrapper) Bootstrap(params map[string]interface{}) error {
	wd, _ := os.Getwd()
	trussDir := filepath.Join(wd, b.TrussDir)

	if _, err := os.Stat(trussDir); err == nil {
		return fmt.Errorf("the target directory [%s] already exists", trussDir)
	}

	l, err := b.LocalDirectory(b.Template)
	if err != nil {
		return err
	}
	return filepath.Walk(l, func(path string, info os.FileInfo, err error) error {
		rel, _ := filepath.Rel(l, path)
		data := map[string]interface{}{
			"Params":   params,
			"TrussDir": b.TrussDir,
			"Template": b.Template,
		}

		if rel == ".truss-manifest.yaml" {
			return nil
		}

		if info.IsDir() {
			dstDir, err := quickTemplate(filepath.Join(trussDir, rel), data)
			if err != nil {
				return err
			}
			return os.MkdirAll(dstDir, 0766)
		}

		t, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		dstFile, err := quickTemplate(filepath.Join(trussDir, rel), data)
		if err != nil {
			return err
		}

		dst, err := os.Create(dstFile)
		if err != nil {
			return err
		}
		defer dst.Close()

		return t.Execute(dst, data)
	})
}

func quickTemplate(tmpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("dstFile").Parse(tmpl)
	if err != nil {
		return "", err
	}
	out := bytes.NewBuffer(nil)
	if err := t.Execute(out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}
