package truss

import (
	"errors"
	"os/exec"
	"syscall"
)

// KubectlCmd wrapper for kubectl
type KubectlCmd struct {
	portForwardCmd *exec.Cmd
}

// Kubectl wrapper for kubectl
func Kubectl(context string) (*KubectlCmd, error) {
	kubectl := &KubectlCmd{}
	if err := kubectl.SetContext(context); err != nil {
		return nil, err
	}
	return kubectl, nil
}

// SetContext kubectl config set-context
func (*KubectlCmd) SetContext(context string) error {
	cmd := exec.Command("kubectl", "config", "set-context", context)
	if _, err := cmd.Output(); err != nil {
		return errors.New(string(err.(*exec.ExitError).Stderr))
	}
	return nil
}

// PortForward kubectl port-forward
func (kubectl *KubectlCmd) PortForward(arg ...string) error {
	argsWithCmd := []string{"port-forward"}
	argsWithCmd = append(argsWithCmd, arg...)
	kubectl.portForwardCmd = exec.Command("kubectl", argsWithCmd...)
	if err := kubectl.portForwardCmd.Start(); err != nil {
		return errors.New(string(err.(*exec.ExitError).Stderr))
	}
	return nil
}

// ClosePortForward sigterm kubectl port-forward
func (kubectl *KubectlCmd) ClosePortForward() error {
	return kubectl.portForwardCmd.Process.Signal(syscall.SIGTERM)
}
