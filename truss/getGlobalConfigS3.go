package truss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// GetGlobalConfigS3Input input for GetGlobalConfigS3
type GetGlobalConfigS3Input struct {
	Bucket string
	Region string
	Key    string
	Role   string
	Dir    string
}

// GetGlobalConfigS3 fetch global config from S3 and put it in home dir
func GetGlobalConfigS3(input *GetGlobalConfigS3Input) (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(input.Region)},
	)
	if input.Role != "" {
		sess.Config.Credentials = stscreds.NewCredentials(sess, input.Role)
	}

	dir := input.Dir
	fileName := path.Join(dir, ".truss.yaml")

	//store existing content
	oldContent, readFileErr := ioutil.ReadFile(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(input.Bucket),
			Key:    &input.Key,
		})

	// if access to s3 throws an error but we had existing content, restore it
	if err != nil && readFileErr == nil {
		ioutil.WriteFile(fileName, oldContent, os.ModeAppend.Perm())
	}
	if err != nil && err.(awserr.Error).Code() == "NoCredentialProviders" {
		fmt.Println("It seems that you forgot to configure aws access")
	}
	return file.Name(), err
}
