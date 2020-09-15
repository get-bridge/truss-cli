package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	// Vault Not installed on CI
	SkipConvey("Vault", t, func() {
		vault := Vault("", nil)
		vault.(*VaultCmdImpl).timeoutSeconds = 0

		Convey("Run", func() {
			Convey("shows error", func() {
				_, err := vault.Run([]string{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "Vault command failed: Usage")
			})
		})
	})
}
