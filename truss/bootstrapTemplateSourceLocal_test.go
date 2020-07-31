package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplateSourceLocal(t *testing.T) {
	Convey("TemplateSourceLocal", t, func() {
		ts := NewLocalTemplateSource("../bootstrap-templates")

		Convey("it can list templates", func() {
			t, err := ts.ListTemplates()

			So(t, ShouldContain, "default")
			So(err, ShouldBeNil)
		})

		Convey("it knows the local directory", func() {
			So(ts.LocalDirectory("default"), ShouldEqual, "../bootstrap-templates/default")
		})

		Convey("it can get template manifests", func() {
			m := ts.GetTemplateManifest("default")

			So(m, ShouldNotBeNil)
			So(len(m.Params), ShouldNotEqual, 0)
		})

		Convey("it can't get missing template manifests", func() {
			m := ts.GetTemplateManifest("should-never-exist")

			So(m, ShouldBeNil)
		})
	})
}
