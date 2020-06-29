package truss

import (
	"bytes"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrap(t *testing.T) {
	Convey("Wrap", t, func() {
		env := "edge-cmh"
		kubeconfigs := map[string]interface{}{
			"edge-cmh": "kubeconfig-truss-nonprod-cmh",
		}
		kubeDir := "/home/test/.kube/"

		out := bytes.NewBufferString("")

		Convey("runs subcommand", func() {
			input := &WrapInput{
				Env:         env,
				Kubeconfigs: kubeconfigs,
				KubeDir:     kubeDir,
				Stdout:      out,
				Stderr:      os.Stderr,
				Stdin:       os.Stdin,
			}

			err := Wrap(input, "echo", "hello")
			So(err, ShouldBeNil)
			So(out.String(), ShouldContainSubstring, "hello")
		})

		Convey("runs subcommand when kubeconfig is not found", func() {
			input := &WrapInput{
				Env:         "no-env",
				Kubeconfigs: kubeconfigs,
				KubeDir:     kubeDir,
				Stdout:      out,
				Stderr:      os.Stderr,
				Stdin:       os.Stdin,
			}

			err := Wrap(input, "echo", "hello")
			So(err, ShouldBeNil)
			So(out.String(), ShouldContainSubstring, "hello")
		})

		Convey("sets KUBECONFIG for the command", func() {
			input := &WrapInput{
				Env:         env,
				Kubeconfigs: kubeconfigs,
				KubeDir:     kubeDir,
				Stdout:      out,
				Stderr:      bytes.NewBufferString(""),
				Stdin:       os.Stdin,
			}
			err := Wrap(input, "printenv")
			So(err, ShouldBeNil)
			So(out.String(), ShouldContainSubstring, "KUBECONFIG=/home/test/.kube/kubeconfig-truss-nonprod-cmh")
		})

		Convey("returns error and sends it to stderr", func() {
			errOut := bytes.NewBufferString("")

			input := &WrapInput{
				Env:         env,
				Kubeconfigs: kubeconfigs,
				KubeDir:     kubeDir,
				Stdout:      out,
				Stderr:      errOut,
				Stdin:       os.Stdin,
			}
			err := Wrap(input, "ls", "asdf")
			So(err.Error(), ShouldEqual, "exit status 1")
			So(errOut.String(), ShouldContainSubstring, "ls: asdf: No such file or directory")
		})
	})
}
