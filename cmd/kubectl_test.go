package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestKubectl(t *testing.T) {
	Convey("kubectl", t, func() {
		c := &cobra.Command{}

		Convey("runs no errors", func() {
			err := kubectlCmd.RunE(c, []string{})
			So(err, ShouldBeNil)
		})

		Convey("reports errors", func() {
			err := kubectlCmd.RunE(c, []string{"foo"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "unknown command \"foo\" for \"kubectl\"")
		})
	})
}
