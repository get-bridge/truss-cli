package truss

import (
	"bytes"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrap(t *testing.T) {
	Convey("Wrap", t, func() {
		out := bytes.NewBufferString("")

		Convey("runs subcommand", func() {
			input := &WrapInput{
				Kubeconfig: "/kube/config",
				Stdout:     out,
				Stderr:     os.Stderr,
				Stdin:      os.Stdin,
			}

			err := Wrap(input, "echo", "hello")
			So(err, ShouldBeNil)
			So(out.String(), ShouldContainSubstring, "hello")
		})

		Convey("sets KUBECONFIG for the command", func() {
			input := &WrapInput{
				Kubeconfig: "/kube/config",
				Stdout:     out,
				Stderr:     bytes.NewBufferString(""),
				Stdin:      os.Stdin,
			}
			err := Wrap(input, "printenv")
			So(err, ShouldBeNil)
			So(out.String(), ShouldContainSubstring, "KUBECONFIG=/kube/config")
		})

		Convey("returns error and sends it to stderr", func() {
			errOut := bytes.NewBufferString("")

			input := &WrapInput{
				Kubeconfig: "/kube/config",
				Stdout:     out,
				Stderr:     errOut,
				Stdin:      os.Stdin,
			}
			err := Wrap(input, "ls", "asdf")
			So(err.Error(), ShouldContainSubstring, "exit status")
			So(err.Error(), ShouldNotEqual, "exit status 0")
			So(errOut.String(), ShouldContainSubstring, "ls")
			So(errOut.String(), ShouldContainSubstring, "asdf")
			So(errOut.String(), ShouldContainSubstring, "No such file or directory")
		})
	})
}
