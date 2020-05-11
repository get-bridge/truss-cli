package cmd

import (
	"testing"
)

func TestSetup(t *testing.T) {
	t.Run("runs no errors", func(t *testing.T) {
		rootCmd.SetArgs([]string{"setup"})
		Execute()
	})
}
