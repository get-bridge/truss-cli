package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestVault(t *testing.T) {
	// TODO tests take 15s to timeout
	SkipConvey("vault", t, func() {
		c := &cobra.Command{}

		Convey("forwards errors", func() {
			err := vaultCmd.RunE(c, []string{""})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Usage: vault")
		})

		Convey("aws auth", func() {
			Convey("errors", func() {
				err := vaultCmd.RunE(c, []string{"status"})
				So(err, ShouldNotBeNil)
			})
		})
	})
}
