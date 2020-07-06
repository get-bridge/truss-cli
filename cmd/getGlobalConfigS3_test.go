package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestGetGlobalConfigS3(t *testing.T) {
	Convey("getGlobalConfigS3", t, func() {
		viper.Reset()

		Convey("pulls the file down", func() {
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
				key = ".truss.yaml"
			}
			tmp := os.TempDir()
			cmd := rootCmd
			buffz := bytes.NewBufferString("")
			cmd.SetOut(buffz)
			cmd.SetArgs([]string{
				"get-global-config",
				"s3",
				"--role",
				awsrole,
				"--bucket",
				bucket,
				"--key",
				key,
				"--region",
				region,
				"--out",
				tmp,
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buffz)
			So(string(out), ShouldContainSubstring, fmt.Sprintf("Global config written to %s/%s", path.Clean(tmp), key))
		})
	})
}
