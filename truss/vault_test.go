package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	Convey("Vault", t, func() {
		cmd := Vault(Kubectl(""), nil)

		Convey("PortForward", func() {
			Convey("runs no errors", func() {
				port, err := cmd.PortForward()
				So(err, ShouldBeNil)
				So(port, ShouldNotBeEmpty)

				port2, err := cmd.PortForward()
				So(err, ShouldBeNil)
				So(port, ShouldEqual, port2)

				err = cmd.ClosePortForward()
				So(err, ShouldBeNil)
			})
		})

		Convey("Run", func() {
			Convey("runs no errors", func() {
				_, err := cmd.Run([]string{"status"})
				So(err, ShouldBeNil)
			})

			Convey("forwards errors", func() {
				_, err := cmd.Run([]string{})
				So(err, ShouldNotBeNil)
			})
		})
	})
}
