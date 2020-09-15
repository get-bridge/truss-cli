package integration

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
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
		tmp, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)
		defer os.Remove(tmp)

		cmd := truss.GetKubeconfigS3(awsrole, bucket, tmp, region)

		Convey("runs with no errors", func() {
			err := cmd.Fetch()
			So(err, ShouldBeNil)
		})
	})
}
