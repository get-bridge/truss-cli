package truss

import (
	"io"
	"os"
	"os/exec"
)

// WrapInput input for Wrap
type WrapInput struct {
	Kubeconfig string
	Stdout     io.Writer
	Stderr     io.Writer
	Stdin      io.Reader
}

// Wrap exports relevant kubeconfig and runs command
func Wrap(input *WrapInput, bin string, arg ...string) error {
	cmd := exec.Command(bin, arg...)
	cmd.Stdout = input.Stdout
	cmd.Stdin = input.Stdin
	cmd.Stderr = input.Stderr

	if input.Kubeconfig != "" {
		cmd.Env = append(os.Environ(), "KUBECONFIG="+input.Kubeconfig)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
