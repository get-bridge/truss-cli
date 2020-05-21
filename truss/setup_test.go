package truss

import (
	"testing"
)

func TestSetup(t *testing.T) {
	t.Run("runs no errors", func(t *testing.T) {
		dep := []string{}
		if err := Setup(&dep); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("runs error if dependency not found", func(t *testing.T) {
		dep := []string{"a_mysterious_program"}
		if err := Setup(&dep); err == nil {
			t.Fatal("Expected Setup to fail with invalid dependency")
		}
	})
}
