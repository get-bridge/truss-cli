package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestSecretsEdit(t *testing.T) {
	Convey("secrets edit", t, func() {
		c := &cobra.Command{}

		Convey("errors if no such configuration", func() {
			err := secretsEditCmd.RunE(c, []string{"secret-name", "kubeconfig-name"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "secret named 'secret-name' in 'kubeconfig-name' not found")
		})

		// TODO how can we mock edit?
		Convey("runs with no errors", nil)
	})
}
