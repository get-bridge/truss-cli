package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVaultIntegration(t *testing.T) {
	Convey("Vault", t, func() {
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
