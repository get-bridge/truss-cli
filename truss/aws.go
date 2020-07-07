package truss

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewAWSSession creates an AWS session with role assumed
func NewAWSSession(region string, assumeRole string) *session.Session {
	cfg := &aws.Config{Region: aws.String(region)}
	sess, _ := session.NewSession(cfg)

	if assumeRole != "" {
		sess.Config.Credentials = stscreds.NewCredentials(sess, assumeRole)
	}

	return sess
}
