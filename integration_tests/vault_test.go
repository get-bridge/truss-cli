package integration

import (
	"os"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVaultIntegration(t *testing.T) {
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

		Convey("Run", func() {
			Convey("runs no errors", func() {
				_, err := vault.Run([]string{"status"})
				So(err, ShouldBeNil)
			})

			Convey("forwards errors", func() {
				_, err := vault.Run([]string{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "Vault command failed:")
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

		Convey("GetMap", func() {
			vaultPath := "secret/bridge/truss-cli-test/getMap"

			_, err := vault.Run([]string{"kv", "put", vaultPath, "foo=bar"})
			So(err, ShouldBeNil)

			data, err := vault.GetMap(vaultPath)
			So(err, ShouldBeNil)
			So(data, ShouldResemble, map[string]string{"foo": "bar"})
		})
	})
}
