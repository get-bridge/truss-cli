package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("runs no errors", func(t *testing.T) {
		rootCmd.SetArgs([]string{})
		Execute()
	})
	t.Run("returns errors", func(t *testing.T) {
		if os.Getenv("RUN_EXECUTE") == "1" {
			rootCmd.SetArgs([]string{"foo"})
			Execute()
			return
		}

		var stderr strings.Builder
		cmd := exec.Command(os.Args[0], "-test.run=TestExecute/returns_errors")
		cmd.Stderr = &stderr
		cmd.Env = append(os.Environ(), "RUN_EXECUTE=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			if !strings.Contains(stderr.String(), "Error: unknown command \"foo\" for \"truss-cli\"") {
				t.Fatalf("process ran with unexpected error: %v", stderr.String())
			}
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})
}

func TestRootCommand(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		rootCmd.SetArgs([]string{})
		err := rootCmd.Execute()
		if err != nil {
			t.Fatalf("execute returned %s", err)
		}
	})

	t.Run("invalid args", func(t *testing.T) {
		rootCmd.SetArgs([]string{"foo"})
		err := rootCmd.Execute()
		if err != nil && err.Error() != "unknown command \"foo\" for \"truss-cli\"" {
			t.Fatalf("execute should fail with invalid args. got %s", err)
		}
	})
}
