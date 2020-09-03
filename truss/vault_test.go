package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	Convey("Vault", t, func() {
		vault := Vault(Kubectl(""), nil)

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
			})
		})

		Convey("ClosePortForward", func() {
			Convey("set portForwarded to nil", func() {
				portForwarded := "true"
				vault.portForwarded = &portForwarded
				err := vault.ClosePortForward()
				So(err, ShouldNotBeNil)
				So(vault.portForwarded, ShouldNotBeNil)
			})
		})

		// TODO
		Convey("Decrypt", nil)
		Convey("Encrypt", nil)
	})
}
