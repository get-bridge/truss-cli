package integration

import (
	"strconv"
	"testing"

	"github.com/phayes/freeport"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKubectl(t *testing.T) {
	Convey("Kubectl", t, func() {
		cmd := Kubectl("")

		Convey("PortForward", func() {
			Convey("runs with no errors", func() {
				p, err := freeport.GetFreePort()
				So(err, ShouldBeNil)
				port := strconv.Itoa(p)

				err = cmd.PortForward(port, port, "namespace", "target")
				So(err, ShouldBeNil)

				err = cmd.ClosePortForward()
				So(err, ShouldBeNil)
			})
		})

		Convey("Run", func() {
			Convey("runs with no errors", func() {
				_, err := cmd.Run()
				So(err, ShouldBeNil)
			})
		})
	})
}
