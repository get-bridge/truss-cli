package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBootstrapConfig(t *testing.T) {
	Convey("BootstrapConfig", t, func() {
		Convey("it parses a valid config file", func() {
			c, err := LoadBootstrapConfig("./fixtures/bootstrap.truss.yaml")

			So(err, ShouldBeNil)
			So(c.TrussDir, ShouldEqual, ".test")
			So(c.Template, ShouldEqual, "test")
			So(c.TemplateSource.Type, ShouldEqual, "local")
		})

		Convey("it errors on missing config file", func() {
			c, err := LoadBootstrapConfig("/tmp/foo")

			So(c, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
