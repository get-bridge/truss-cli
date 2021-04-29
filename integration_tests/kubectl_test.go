package integration

import (
	"testing"

	"github.com/get-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKubectl(t *testing.T) {
	Convey("Kubectl", t, func() {
		kubectl := truss.Kubectl("")

		Convey("Run", func() {
			Convey("runs with no errors", func() {
				_, err := kubectl.Run("get", "pods")
				So(err, ShouldBeNil)
			})
		})
	})
}
