package truss

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

// SshuttleRouteProvider provides routes by tunneling with sshuttle
type SshuttleRouteProvider struct {
	Host   string
	Routes []string
}

// Setup creates the routes
func (n SshuttleRouteProvider) Setup() error {
	time.Sleep(3 * time.Second)
	pidfile := "/tmp/truss-sshuttle.pid"
	cmd := exec.Command("sshuttle", "-D", "--pidfile", pidfile, "-r", n.Host)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	for _, r := range n.Routes {
		cmd.Args = append(cmd.Args, fmt.Sprintf("%s/32", r))
	}

	return cmd.Run()
}

// Teardown destroys the routes
func (n SshuttleRouteProvider) Teardown() error {
	pidfile := "/tmp/truss-sshuttle.pid"
	pid, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return err
	}

	return exec.Command("kill", string(pid)).Run()
}
