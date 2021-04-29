package truss

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitTemplateSource is a Local Filesystem Template Source
type GitTemplateSource struct {
	TemplateSource
	CloneURL    string `default:"git@github.com:get-bridge/truss-cli.git" yaml:"clone_url"`
	Directory   string `default:"bootstrap-templates"`
	CheckoutRef string `yaml:"checkout_ref"`
	tmpDir      string
}

// NewGitTemplateSource returns a new TemplateSource
func NewGitTemplateSource(cloneURL, directory, checkoutRef string) (*GitTemplateSource, error) {
	wd, _ := os.Getwd()
	tmpDir, err := ioutil.TempDir(wd, ".bootstrap-template-")
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(tmpDir, 0766); err != nil {
		return nil, err
	}
	o := &git.CloneOptions{
		URL: cloneURL,
	}
	if checkoutRef != "" {
		o.ReferenceName = plumbing.ReferenceName(checkoutRef)
	}

	if _, err := git.PlainClone(tmpDir, false, o); err != nil {
		return nil, err
	}

	ts := NewLocalTemplateSource(filepath.Join(tmpDir, directory))
	return &GitTemplateSource{
		TemplateSource: ts,
		CloneURL:       cloneURL,
		Directory:      directory,
		CheckoutRef:    checkoutRef,
		tmpDir:         tmpDir,
	}, nil
}

// Cleanup cleans up
func (s GitTemplateSource) Cleanup() {
	os.RemoveAll(s.tmpDir)
}
