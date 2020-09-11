package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestSecretsEncrypt(t *testing.T) {
	Convey("secrets encrypt", t, func() {
		c := &cobra.Command{}

		Convey("errors if no such configuration", func() {
			err := secretsViewCmd.RunE(c, []string{"secret-name", "kubeconfig-name"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "secret named 'secret-name' in 'kubeconfig-name' not found")
		})

		// TODO how can we mock encrypt?
		Convey("runs with no errors", nil)
	})
}
