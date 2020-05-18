package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestSetup(t *testing.T) {
	t.Run("runs no errors", func(t *testing.T) {
		rootCmd.SetArgs([]string{"setup"})
		viper.Set("dependencies", []interface{}{})
		Execute()
	})
}
