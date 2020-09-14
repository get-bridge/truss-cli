package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKubectl(t *testing.T) {
	Convey("Kubectl", t, func() {
		cmd := Kubectl("")

		Convey("Run", func() {
			Convey("runs with no errors", func() {
				_, err := cmd.Run()
				So(err, ShouldBeNil)
			})
		})
	})
}
