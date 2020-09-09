package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplateSourceS3(t *testing.T) {
	Convey("TemplateSourceS3", t, func() {
		ts := NewS3TemplateSource(
			"truss-cli-global-config", "bootstrap-templates", "us-east-2", "arn:aws:iam::127178877223:role/xacct/ops-admin",
		)
		_ = ts

		Convey("it can list templates", func() {
			//
		})

		Convey("it knows the local directory", func() {
			//
		})

		Convey("it can get template manifests", func() {
			//
		})

		Convey("it can't get missing template manifests", func() {
			//
		})
	})
}
