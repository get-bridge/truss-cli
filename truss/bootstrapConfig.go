package truss

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

// BootstrapConfig represetns a Bootstrap configuration
type BootstrapConfig struct {
	TemplateSource struct {
		Type  string `default:"local"`
		Local struct {
			Directory string `default:"./bootstrap-templates"`
		}
		S3 struct {
			Bucket string `default:"truss-cli-global-config"`
			Region string `default:"us-east-2"`
			Prefix string `default:"bootstrap-templates"`
			Role   string `default:"arn:aws:iam::127178877223:role/xacct/ops-admin"`
		}
		Git struct {
			CloneURL    string `default:"git@github.com:instructure-bridge/truss-cli.git" yaml:"clone_url"`
			Directory   string `default:"bootstrap-templates"`
			CheckoutRef string `yaml:"checkout_ref"`
		}
	} `yaml:"templateSource"`
	TrussDir string `default:"truss" yaml:"trussDir"`
	Template string `default:"default"`
	Params   map[string]interface{}
}

// LoadBootstrapConfig loads a config from disk
func LoadBootstrapConfig(name string) (*BootstrapConfig, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &BootstrapConfig{}
	defaults.Set(c)
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}

// GetTemplateSource gets a TemplateSource from a configuration
func (c BootstrapConfig) GetTemplateSource() (t TemplateSource, err error) {
	switch c.TemplateSource.Type {
	case "local":
		t = NewLocalTemplateSource(c.TemplateSource.Local.Directory)
		return
	case "s3":
		t = NewS3TemplateSource(
			c.TemplateSource.S3.Bucket,
			c.TemplateSource.S3.Prefix,
			c.TemplateSource.S3.Region,
			c.TemplateSource.S3.Role,
		)
		return
	case "git":
		t, err = NewGitTemplateSource(
			c.TemplateSource.Git.CloneURL,
			c.TemplateSource.Git.Directory,
			c.TemplateSource.Git.CheckoutRef,
		)
		return
	}
	return nil, fmt.Errorf("Invalid templateSource.type: %s", c.TemplateSource.Type)
}

// GetBootstrapper gets a Bootstrapper from a configuration
func (c BootstrapConfig) GetBootstrapper() (b *Bootstrapper, err error) {
	ts, err := c.GetTemplateSource()
	if err != nil {
		return nil, err
	}

	return NewBootstrapper(ts, c.TrussDir, c.Template), nil
}
