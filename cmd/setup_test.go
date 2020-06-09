package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestSetup(t *testing.T) {
	Convey("setup", t, func() {
		c := &cobra.Command{}

		Convey("runs no errors", func() {
			err := setupCmd.RunE(c, []string{})
			So(err, ShouldBeNil)
		})

		Convey("with dependencies", func() {
			Convey("runs no errors", func() {
				viper.Set("dependencies", []string{"bash"})
				err := setupCmd.RunE(c, []string{})
				So(err, ShouldBeNil)
			})
			Convey("returns missing dependencies", func() {
				viper.Set("dependencies", []string{"foo"})
				err := setupCmd.RunE(c, []string{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "missing dependencies: [foo]")
			})
		})
	})
}
