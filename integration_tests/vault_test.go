package integration

import (
	"os"
	"testing"

	"github.com/get-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	Convey("Vault", t, func() {
		var auth truss.VaultAuth
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if ok {
			vaultrole := os.Getenv("TEST_VAULT_ROLE")
			auth = truss.VaultAuthAWS(vaultrole, awsrole, "us-east-1")
		}
		vault := truss.Vault("", auth)

		Convey("PortForward", func() {
			Convey("runs no errors", func() {
				port, err := vault.PortForward()
				So(err, ShouldBeNil)
				So(port, ShouldNotBeEmpty)

				port2, err := vault.PortForward()
				So(err, ShouldBeNil)
				So(port, ShouldEqual, port2)

				err = vault.ClosePortForward()
				So(err, ShouldBeNil)
			})
		})

		Convey("Run", func() {
			bytes, err := vault.Run([]string{"status"})
			So(err, ShouldBeNil)
			So(bytes, ShouldNotBeEmpty)
		})
	})
}
