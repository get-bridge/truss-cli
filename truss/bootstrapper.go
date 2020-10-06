package truss

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

// IncludeTemplateRE - regular expression to match the template inclusion expression
var IncludeTemplateRE = regexp.MustCompile(
	"\\{\\{include ([A-Za-z0-9_-]+) as ((\\.[A-Za-z_][A-Za-z0-9_]+)+)\\}\\}",
)

// Bootstrapper bootstraps a deployment
type Bootstrapper struct {
	TemplateSource
	TrussDir string
	Template string

	overrideUseTrussDir *bool
}

// BootstrapManifest represents the manifest thingy
type BootstrapManifest struct {
	Settings struct {
		UseTrussDir bool `yaml:"useTrussDir"`
	} `yaml:"settings"`

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
	return &Bootstrapper{ts, trussDir, template, nil}
}

// GetTemplateManifest gets a template's manifest
func (b Bootstrapper) GetTemplateManifest() *BootstrapManifest {
	manifest := b.TemplateSource.GetTemplateManifest(b.Template)

	if b.overrideUseTrussDir != nil {
		manifest.Settings.UseTrussDir = *b.overrideUseTrussDir
	}

	return manifest
}

// Bootstrap does the thing!
func (b Bootstrapper) Bootstrap(params *BootstrapParams) error {
	var trussDir string
	settings := b.GetTemplateManifest().Settings

	wd, _ := os.Getwd()
	trussDir = wd

	if settings.UseTrussDir {
		trussDir = filepath.Join(wd, b.TrussDir)
	}

	if _, err := os.Stat(trussDir); settings.UseTrussDir && err == nil {
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
			if IncludeTemplateRE.MatchString(rel) {
				submatch := IncludeTemplateRE.FindStringSubmatch(rel)
				newTemplate := submatch[1]
				nameExpr := fmt.Sprintf("{{ %s }}", submatch[2])

				newRel, err := quickTemplate(nameExpr, data)
				if err != nil {
					return err
				}

				boot := NewBootstrapper(b.TemplateSource, b.TrussDir, newTemplate)
				boot.overrideUseTrussDir = newBoolPtr(false)
				subDir := filepath.Join(wd, newRel)

				os.MkdirAll(subDir, 0766)
				os.Chdir(subDir)
				defer os.Chdir(wd)

				return boot.Bootstrap(params)
			}

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

func newBoolPtr(val bool) *bool {
	b := val
	return &b
}
