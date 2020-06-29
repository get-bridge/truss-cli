package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestEnv(t *testing.T) {
	Convey("env", t, func() {
		viper.Reset()
		viper.Set("kubeconfigfiles.directory", "/tmp/")

		Convey("returns env for a valid environment", func() {
			viper.Set("environments", map[string]interface{}{
				"edge-cmh": "kubeconfig-truss-nonprod-cmh",
			})

			cmd := rootCmd
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs([]string{
				"env",
				"-e",
				"edge-cmh",
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buff)
			So(string(out), ShouldEqual, "export KUBECONFIG=/tmp/kubeconfig-truss-nonprod-cmh\n# Run this command to configure your shell:\n# eval \"$(truss env -e edge-cmh)\"\n")
		})

		Convey("returns error for invalid environment", func() {
			viper.Set("environments", map[string]interface{}{
				"edge-cmh": "kubeconfig-truss-nonprod-cmh",
			})

			cmd := rootCmd
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs([]string{
				"env",
				"-e",
				"no-env",
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buff)
			So(string(out), ShouldContainSubstring, "Error: No kubeconfig found for env no-env")
		})
	})
}
