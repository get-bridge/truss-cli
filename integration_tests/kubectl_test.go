package integration

import (
	"strconv"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/phayes/freeport"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKubectl(t *testing.T) {
	Convey("Kubectl", t, func() {
		kubectl := truss.Kubectl("")

		Convey("PortForward", func() {
			Convey("runs with no errors", func() {
				p, err := freeport.GetFreePort()
				So(err, ShouldBeNil)
				port := strconv.Itoa(p)

				err = kubectl.PortForward(port, port, "namespace", "target", 15)
				So(err, ShouldBeNil)

				err = kubectl.ClosePortForward()
				So(err, ShouldBeNil)
			})
		})
	})
}
