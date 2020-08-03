package truss

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

// S3TemplateSource is an Amazon S3 Template Source
type S3TemplateSource struct {
	Bucket string
	Folder string
	Region string
	Role   string

	tmpDirs []string
}

// NewS3TemplateSource returns a new TemplateSource
func NewS3TemplateSource(bucket, folder, region, role string) *S3TemplateSource {
	return &S3TemplateSource{bucket, folder, region, role, []string{}}
}

// ListTemplates returns a list of temlpates in the template bucket
func (s S3TemplateSource) ListTemplates() ([]string, error) {
	sess, err := s.awsSession()
	if err != nil {
		return nil, err
	}

	api := s3.New(sess)
	r, err := api.ListObjects(&s3.ListObjectsInput{
		Bucket: &s.Bucket,
		Prefix: &s.Folder,
	})
	if err != nil {
		return nil, err
	}

	out := make([]string, 0)
	for _, t := range r.Contents {
		if name := s.getTemplateNameFromS3Key(*t.Key); name != nil {
			out = append(out, *name)
		}
	}

	return out, nil
}

// LocalDirectory returns a local cache of the S3 Template
func (s *S3TemplateSource) LocalDirectory(template string) (string, error) {
	wd, _ := os.Getwd()
	tmpDir, err := ioutil.TempDir(wd, ".bootstrap-template-")
	if err != nil {
		return "", err
	}

	s.tmpDirs = append(s.tmpDirs, tmpDir)

	if err := os.MkdirAll(tmpDir, 0766); err != nil {
		return "", err
	}

	sess, err := s.awsSession()
	if err != nil {
		return "", err
	}
	api := s3.New(sess)

	prefix := filepath.Join(s.Folder, template)
	list, err := api.ListObjects(&s3.ListObjectsInput{
		Bucket: &s.Bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return "", err
	}

	for _, f := range list.Contents {
		rel := strings.Replace(*f.Key, filepath.Join(s.Folder, template), "", -1)
		dst := filepath.Join(tmpDir, rel)

		if err := os.MkdirAll(filepath.Dir(dst), 0766); err != nil {
			return "", err
		}

		if strings.HasSuffix(*f.Key, "/") {
			continue
		}

		file, err := os.Create(dst)
		if err != nil {
			return "", err
		}
		defer file.Close()

		obj, err := api.GetObject(&s3.GetObjectInput{
			Bucket: &s.Bucket,
			Key:    f.Key,
		})
		if err != nil {
			return "", err
		}
		defer obj.Body.Close()

		io.Copy(file, obj.Body)
	}

	return tmpDir, err
}

// GetTemplateManifest parses the template's manifest
func (s S3TemplateSource) GetTemplateManifest(t string) *BootstrapManifest {
	sess, err := s.awsSession()
	if err != nil {
		return nil
	}

	key := filepath.Join(s.Folder, t, ".truss-manifest.yaml")
	o, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil
	}
	defer o.Body.Close()

	m := &BootstrapManifest{}
	if err := yaml.NewDecoder(o.Body).Decode(m); err != nil {
		return nil
	}

	return m
}

// Cleanup removes tmpDirs
func (s *S3TemplateSource) Cleanup() {
	for _, d := range s.tmpDirs {
		os.RemoveAll(d)
	}
}

func (s S3TemplateSource) awsSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.Region)},
	)
	if err != nil {
		return nil, err
	}
	if s.Role != "" {
		sess.Config.Credentials = stscreds.NewCredentials(sess, s.Role)
	}

	return sess, nil
}

func (s S3TemplateSource) getTemplateNameFromS3Key(key string) *string {
	rex := regexp.MustCompile(fmt.Sprintf(`^%s\/([\w-]*)\/$`, s.Folder))
	if !rex.MatchString(key) {
		return nil
	}

	match := rex.FindStringSubmatch(key)
	return &match[1]
}
