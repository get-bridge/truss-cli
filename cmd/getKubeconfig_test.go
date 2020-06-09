package cmd

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestGetKubeconfig(t *testing.T) {
	Convey("get-kubeconfig", t, func() {
		c := &cobra.Command{}
		viper.Reset()

		Convey("runs no errors", func() {
			err := getKubeconfigCmd.RunE(c, []string{})
			So(err, ShouldBeNil)
		})

		Convey("s3 configured", func() {
			role, ok := os.LookupEnv("TEST_AWS_ROLE")
			if !ok {
				t.Fatalf("Missing env var TEST_AWS_ROLE")
			}
			bucket, ok := os.LookupEnv("TEST_S3_BUCKET")
			if !ok {
				t.Fatalf("Missing env var TEST_S3_BUCKET")
			}
			region, ok := os.LookupEnv("TEST_S3_BUCKET_REGION")
			if !ok {
				region = "us-east-2"
			}

			Convey("returns errors", func() {
				viper.Set("kubeconfigfiles.s3", map[string]interface{}{
					"bucket": bucket,
				})
				err := getKubeconfigCmd.RunE(c, []string{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "s3 config must have region")
			})

			Convey("runs no errors", func() {
				viper.Set("kubeconfigfiles.s3", map[string]interface{}{
					"awsrole": role,
					"bucket":  bucket,
					"region":  region,
				})

				err := getKubeconfigCmd.RunE(c, []string{})
				So(err, ShouldBeNil)
			})
		})
	})
}
