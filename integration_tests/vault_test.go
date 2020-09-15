package integration

import (
	"os"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	Convey("Vault", t, func() {
		var auth truss.VaultAuth
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if ok {
			vaultrole := os.Getenv("TEST_VAULT_ROLE")
			auth = truss.VaultAuthAWS(vaultrole, awsrole)
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

		Convey("Encrypt then Decrypt", func() {
			transitKeyName := "test-transit-key"
			input := "my-great-stuff"
			encrypted, err := vault.Encrypt(transitKeyName, []byte(input))
			So(err, ShouldBeNil)
			So(encrypted, ShouldNotBeNil)
			decrypted, err := vault.Decrypt(transitKeyName, encrypted)
			So(err, ShouldBeNil)
			So(string(decrypted), ShouldEqual, input)
		})
	})
}
