package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBootstrapParameter(t *testing.T) {
	Convey("BootstrapParameter", t, func() {
		Convey("Strings are given the type `string`", func() {
			p := NewBootstrapParameter("anything")

			So(p.Type, ShouldEqual, "string")
		})

		Convey("Booleans are given the type `bool` and return a string", func() {
			p := NewBootstrapParameterBool(true)

			So(p.Type, ShouldEqual, "bool")
			So(p.String(), ShouldEqual, "true")
		})

		Convey("Simple case conversions work", func() {
			inputs := []string{"simple-test", "simpleTest", "SimpleTest", "simple_test"}

			for _, input := range inputs {
				p := NewBootstrapParameter(input)

				So(p.CamelCase, ShouldEqual, "simpleTest")
				So(p.PascalCase, ShouldEqual, "SimpleTest")
				So(p.KebabCase, ShouldEqual, "simple-test")
				So(p.SnakeCase, ShouldEqual, "simple_test")
				So(p.FlatCase, ShouldEqual, "simpletest")
			}
		})

		Convey("Even bools can be converted (useful for eg. Python)", func() {
			So(NewBootstrapParameterBool(true).CamelCase, ShouldEqual, "True")
			So(NewBootstrapParameterBool(false).CamelCase, ShouldEqual, "False")
		})
	})
}
