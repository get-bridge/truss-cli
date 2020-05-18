package cmd

import (
	"testing"
)

func TestVault(t *testing.T) {
	t.Run("runs no errors", func(t *testing.T) {
		rootCmd.SetArgs([]string{"vault"})
		Execute()
	})
}
