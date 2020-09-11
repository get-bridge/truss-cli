package integration

import (
	"os"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVaultAuthAws(t *testing.T) {
	Convey("VaultAuthAWS", t, func() {
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if !ok {
			t.Fatalf("Missing env var TEST_AWS_ROLE")
		}
		vaultrole, ok := os.LookupEnv("TEST_VAULT_ROLE")
		if !ok {
			t.Fatalf("Missing env var TEST_VAULT_ROLE")
		}

		cmd := truss.VaultAuthAWS(vaultrole, awsrole)

		Convey("Login", func() {
			Convey("runs no errors", func() {
				port, err := truss.Vault("", nil).PortForward()
				So(err, ShouldBeNil)

				err = cmd.Login(port)
				So(err, ShouldBeNil)
			})
		})
	})
}
