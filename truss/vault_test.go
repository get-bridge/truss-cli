package truss

import (
	"testing"
)

func TestVault(t *testing.T) {
	t.Run("runs no errors", func(t *testing.T) {
		if err := Vault("env"); err != nil {
			t.Fatal(err)
		}
	})
}
