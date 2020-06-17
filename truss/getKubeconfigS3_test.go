package truss

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetKubeconfigS3(t *testing.T) {
	Convey("Fetch", t, func() {
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
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
		tmp := os.TempDir()
		cmd := GetKubeconfigS3(awsrole, bucket, tmp, region)

		Convey("runs with no errors", func() {
			err := cmd.Fetch()
			So(err, ShouldBeNil)
		})

		// Some reason this breaks allthethings
		// Reset(func() {
		// 	os.RemoveAll(tmp)
		// })
	})
}
