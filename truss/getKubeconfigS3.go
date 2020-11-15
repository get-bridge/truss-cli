package truss

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// GetKubeconfigS3Cmd command for managing kubeconfigs
type GetKubeconfigS3Cmd struct {
	awsRole string
	bucket  string
	dest    string
	region  string
}

// GetKubeconfigS3 return command
func GetKubeconfigS3(awsRole string, bucket string, dest string, region string) GetKubeconfigCmd {
	return &GetKubeconfigS3Cmd{
		awsRole: awsRole,
		bucket:  bucket,
		dest:    dest,
		region:  region,
	}
}

// Fetch kubeconfigs
func (config *GetKubeconfigS3Cmd) Fetch() error {
	log.Infoln("Fetching kubeconfig from s3")

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.region)},
	)
	if config.awsRole != "" {
		sess.Config.Credentials = stscreds.NewCredentials(sess, config.awsRole)
	}

	s3Client := s3.New(sess)
	objects, err := s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(config.bucket),
	})
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(sess)
	for _, key := range objects.Contents {
		if err := os.MkdirAll(config.dest, 0755); err != nil && !os.IsExist(err) {
			return err
		}
		file, err := os.Create(path.Join(config.dest, *key.Key))
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(config.bucket),
				Key:    key.Key,
			})
		if err != nil {
			return fmt.Errorf("Unable to download from %q, %v", config.bucket, err)
		}
		log.Infoln("Downloaded", *key.Key)
	}

	return nil
}
