package truss

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3TemplateSource is an Amazon S3 Template Source
type S3TemplateSource struct {
	Bucket string
	Folder string
	Region string
	Role   string
}

// NewS3TemplateSource returns a new TemplateSource
func NewS3TemplateSource(bucket, folder, region, role string) *S3TemplateSource {
	return &S3TemplateSource{bucket, folder, region, role}
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
func (s S3TemplateSource) LocalDirectory(template string) string {
	return ""
}

// GetTemplateManifest parses the template's manifest
func (s S3TemplateSource) GetTemplateManifest(t string) *BootstrapManifest {
	return nil
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
