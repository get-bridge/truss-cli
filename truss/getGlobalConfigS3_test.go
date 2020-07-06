package truss

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetGlobalConfigS3(t *testing.T) {
	Convey("GetGlobalConfigS3", t, func() {
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if !ok {
			t.Fatalf("Missing env var TEST_AWS_ROLE")
		}
		bucket, ok := os.LookupEnv("TEST_GLOBAL_CONFIG_BUCKET")
		if !ok {
			t.Fatalf("Missing env var TEST_GLOBAL_CONFIG_BUCKET")
		}
		region, ok := os.LookupEnv("TEST_S3_BUCKET_REGION")
		if !ok {
			region = "us-east-2"
		}
		key, ok := os.LookupEnv("TEST_GLOBAL_CONFIG_KEY")
		if !ok {
			region = ".truss.yaml"
		}
		tmp := os.TempDir()
		input := &GetGlobalConfigS3Input{
			Bucket: bucket,
			Region: region,
			Role:   awsrole,
			Dir:    tmp,
			Key:    key,
		}

		Convey("runs with no errors", func() {
			_, err := GetGlobalConfigS3(input)
			So(err, ShouldBeNil)
			_, err = os.Stat(fmt.Sprintf("%s/%s", tmp, key))
			So(err, ShouldBeNil)
		})
	})
}
