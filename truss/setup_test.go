package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Setup", t, func() {
		Convey("runs no errors", func() {
			dep := []string{}
			err := Setup(&dep)
			So(err, ShouldBeNil)
		})

		Convey("runs error if dependency not found", func() {
			err := Setup(&[]string{"a_mysterious_program"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "missing dependencies: [a_mysterious_program]")
		})
	})
}
