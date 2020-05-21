package truss

import (
	"errors"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// GetKubeconfigS3Cmd command for managing kubeconfigs
type GetKubeconfigS3Cmd struct {
	bucket string
	dest   string
}

// GetKubeconfigS3 return command
func GetKubeconfigS3(bucket string, dest string) GetKubeconfigCmd {
	return &GetKubeconfigS3Cmd{
		bucket: bucket,
		dest:   dest,
	}
}

// Fetch kubeconfigs
func (config *GetKubeconfigS3Cmd) Fetch() error {
	log.Infoln("Fetching kubeconfig from s3")
	bucketSrc := fmt.Sprintf("s3://%s", config.bucket)
	cmd := exec.Command("aws", "s3", "cp", "--recursive", bucketSrc, config.dest)
	if _, err := cmd.Output(); err != nil {
		return errors.New(string(err.(*exec.ExitError).Stderr))
	}
	return nil
}
