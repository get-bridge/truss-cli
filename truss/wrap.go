package truss

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// WrapInput input for Wrap
type WrapInput struct {
	Env         string
	Kubeconfigs map[string]interface{}
	KubeDir     string
	Stdout      io.Writer
	Stderr      io.Writer
	Stdin       io.Reader
}

// Wrap exports relevant kubeconfig and runs command
func Wrap(input *WrapInput, bin string, arg ...string) error {
	cmd := exec.Command(bin, arg...)
	cmd.Stdout = input.Stdout
	cmd.Stdin = input.Stdin
	cmd.Stderr = input.Stderr
	envKubeconfig := input.Kubeconfigs[input.Env]
	var kubeconfig string

	if envKubeconfig != nil {
		kubeconfigName := fmt.Sprintf("%s", envKubeconfig)
		kubeconfig = fmt.Sprintf("%s%s", input.KubeDir, kubeconfigName)
		cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfig)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
