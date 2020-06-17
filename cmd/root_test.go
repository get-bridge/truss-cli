package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExecute(t *testing.T) {
	Convey("Execute", t, func() {
		Convey("runs no errors", func() {
			rootCmd.SetArgs([]string{})
			Execute()
		})

		Convey("returns errors", func() {
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
				So(stderr.String(), ShouldContainSubstring, "Error: unknown command \"foo\" for \"truss-cli\"")
				return
			}
			t.Fatalf("process ran with err %v, want exit status 1", err)
		})
	})
}

func TestRootCommand(t *testing.T) {
	Convey("Execute", t, func() {
		Convey("no args", func() {
			rootCmd.SetArgs([]string{})
			err := rootCmd.Execute()
			So(err, ShouldBeNil)
		})

		Convey("invalid args", func() {
			rootCmd.SetArgs([]string{"foo"})
			err := rootCmd.Execute()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "unknown command \"foo\" for \"truss-cli\"")
		})
	})
}
