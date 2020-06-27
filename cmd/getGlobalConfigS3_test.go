package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestGetGlobalConfigS3(t *testing.T) {
	Convey("wrap", t, func() {
		viper.Reset()

		Convey("shows usage if passed --help", func() {
			viper.Set("environments", map[string]interface{}{
				"edge-cmh": "kubeconfig-truss-nonprod-cmh",
			})
			cmd := rootCmd
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs([]string{
				"get-global-config",
				"s3",
				"--help",
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buff)
			So(string(out), ShouldContainSubstring, "Fetches .truss.yaml from S3 and puts it in your home directory")
		})
	})
}
