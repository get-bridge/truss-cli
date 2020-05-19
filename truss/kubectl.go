package truss

import (
	"errors"
	"net"
	"os/exec"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// KubectlCmd wrapper for kubectl
type KubectlCmd struct {
	portForwardCmd *exec.Cmd
}

// Kubectl wrapper for kubectl
func Kubectl(context string) (*KubectlCmd, error) {
	kubectl := &KubectlCmd{}
	if context != "" {
		if err := kubectl.UseContext(context); err != nil {
			return nil, err
		}
	}
	return kubectl, nil
}

// UseContext kubectl config use-context
func (*KubectlCmd) UseContext(context string) error {
	log.Debugln("Using context ", context)
	cmd := exec.Command("kubectl", "config", "use-context", context)
	if _, err := cmd.Output(); err != nil {
		return errors.New(string(err.(*exec.ExitError).Stderr))
	}
	return nil
}

// PortForward kubectl port-forward
func (kubectl *KubectlCmd) PortForward(port string, namespace string, target string) error {
	log.Debugln("Opening connection port forward for", port)
	argsWithCmd := []string{"port-forward", "-n=" + namespace, target, port}
	kubectl.portForwardCmd = exec.Command("kubectl", argsWithCmd...)
	if err := kubectl.portForwardCmd.Start(); err != nil {
		return errors.New(string(err.(*exec.ExitError).Stderr))
	}

	waitForPort(port)

	return nil
}

// ClosePortForward sigterm kubectl port-forward
func (kubectl *KubectlCmd) ClosePortForward() error {
	return kubectl.portForwardCmd.Process.Signal(syscall.SIGTERM)
}

func waitForPort(port string) {
	log.Debugln("Waiting for port", port)
	timeout := 15
	for i := 0; i < timeout; i++ {
		conn, err := net.Dial("tcp", ":"+port)
		if conn != nil {
			defer conn.Close()
		}
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}

	log.Warnln("Could not reach port", port)
}
