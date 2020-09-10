package truss

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	Convey("Vault", t, func() {

		var auth VaultAuth
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if ok {
			vaultrole := os.Getenv("TEST_VAULT_ROLE")
			auth = VaultAuthAWS(vaultrole, awsrole)
		}
		vault := VaultCmdImpl{
			kubectl: Kubectl(""),
			auth:    auth,
		}

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
				So(vault.portForwarded, ShouldBeNil)
			})
		})

		Convey("Run", func() {
			Convey("runs no errors", func() {
				_, err := vault.Run([]string{"status"})
				So(err, ShouldBeNil)
			})

			Convey("forwards errors", func() {
				_, err := vault.Run([]string{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "Usage: vault")
			})
		})

		Convey("Decrypt", func() {
			Convey("errors if no transitKeyName provided", func() {
				transitKeyName := ""
				encrypted := []byte{}
				_, err := vault.Decrypt(transitKeyName, encrypted)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Must provide transitkey to decrypt")
			})
		})

		Convey("Encrypt", func() {
			Convey("errors if no transitKeyName provided", func() {
				transitKeyName := ""
				raw := []byte{}
				_, err := vault.Encrypt(transitKeyName, raw)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Must provide transitkey to encrypt")
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
