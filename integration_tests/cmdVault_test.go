package integration

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestVault(t *testing.T) {
	Convey("vault", t, func() {
		c := &cobra.Command{}

		// TODO fragile
		SkipConvey("runs no errors", func() {
			err := vaultCmd.RunE(c, []string{"status"})
			So(err, ShouldBeNil)
		})

		Convey("forwards errors", func() {
			err := vaultCmd.RunE(c, []string{""})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Usage: vault")
		})

		Convey("aws auth", func() {
			awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
			if !ok {
				t.Fatalf("Missing env var TEST_AWS_ROLE")
			}
			vaultrole, ok := os.LookupEnv("TEST_VAULT_ROLE")
			if !ok {
				t.Fatalf("Missing env var TEST_VAULT_ROLE")
			}

			viper.Set("vault.auth.aws.awsrole", awsrole)
			viper.Set("vault.auth.aws.vaultrole", vaultrole)

			// TODO fragile
			SkipConvey("runs no errors", func() {
				err := vaultCmd.RunE(c, []string{"status"})
				So(err, ShouldBeNil)
			})
		})
	})
}
