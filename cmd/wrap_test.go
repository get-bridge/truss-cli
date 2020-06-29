package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestWrap(t *testing.T) {
	Convey("wrap", t, func() {
		viper.Reset()

		Convey("runs subcommand", func() {
			viper.Set("environments", map[string]interface{}{
				"edge-cmh": "kubeconfig-truss-nonprod-cmh",
			})
			cmd := rootCmd
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs([]string{
				"wrap",
				"-e",
				"edge-cmh",
				"--",
				"echo",
				"hello",
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buff)
			So(string(out), ShouldEqual, "hello\n")
		})

		Convey("shows subcommand errors", func() {
			viper.Set("environments", map[string]interface{}{
				"edge-cmh": "kubeconfig-truss-nonprod-cmh",
			})
			cmd := rootCmd
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			errBuff := bytes.NewBufferString("")
			cmd.SetErr(errBuff)
			cmd.SetArgs([]string{
				"wrap",
				"-e",
				"edge-cmh",
				"--",
				"ls",
				"asdf",
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buff)
			errOut, _ := ioutil.ReadAll(errBuff)
			So(string(out), ShouldContainSubstring, "Error: exit status 1\n")
			So(string(errOut), ShouldContainSubstring, "ls: asdf: No such file or directory")
		})

		Convey("shows usage if passed zero args ", func() {
			viper.Set("environments", map[string]interface{}{
				"edge-cmh": "kubeconfig-truss-nonprod-cmh",
			})
			cmd := rootCmd
			buff := bytes.NewBufferString("")
			cmd.SetOut(buff)
			cmd.SetArgs([]string{
				"wrap",
			})
			cmd.Execute()
			out, _ := ioutil.ReadAll(buff)
			So(string(out), ShouldContainSubstring, "Sets KUBECONFIG and then executes the subcommand")
		})
	})
}
