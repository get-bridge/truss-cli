package truss

import (
	"fmt"
	"os"
	"os/exec"
)

// OpenConnect represents an OpenConnect Connection
type OpenConnect struct {
	User      string
	Authgroup string
	Server    string
	Env       []string
	Script    *string

	hooks map[OpenConnectEvent][]OpenConnectHook
}

// OpenConnectHook is a function that can be invoked when an event occurs
type OpenConnectHook func() error

// OpenConnectEvent represents a connection event
type OpenConnectEvent int

const (
	OpenConnectConnecting   OpenConnectEvent = iota
	OpenConnectConnected    OpenConnectEvent = iota
	OpenConnectDisconnected OpenConnectEvent = iota
)

// NewOpenConnect returns a default OpenConnect instance
func NewOpenConnect(user, server, authGroup string) *OpenConnect {
	return &OpenConnect{
		User:      user,
		Server:    server,
		Authgroup: authGroup,
	}
}

// Start interactively starts the tunnel
func (c *OpenConnect) Start() error {
	if err := c.dispatch(OpenConnectConnecting); err != nil {
		return err
	}

	cmd := exec.Command("sudo", "-E", "openconnect", "--background", "--quiet",
		fmt.Sprintf("--user=%s", c.User),
		fmt.Sprintf("--authgroup=%s", c.Authgroup),
		c.Server,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if c.Script != nil {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--script=%s", *c.Script))
	}

	cmd.Env = os.Environ()
	for _, v := range c.Env {
		cmd.Env = append(cmd.Env, v)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return c.dispatch(OpenConnectConnected)
}

// Stop stops the tunnel
func (c OpenConnect) Stop() error {
	cmd := exec.Command("sudo", "pkill", "-2", "openconnect")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return c.dispatch(OpenConnectDisconnected)
}

// AddHook adds an event hook
func (c *OpenConnect) AddHook(k OpenConnectEvent, h OpenConnectHook) {
	if c.hooks == nil {
		c.hooks = map[OpenConnectEvent][]OpenConnectHook{}
	}
	if c.hooks[k] == nil {
		c.hooks[k] = []OpenConnectHook{}
	}

	c.hooks[k] = append(c.hooks[k], h)
}

func (c OpenConnect) dispatch(k OpenConnectEvent) error {
	if c.hooks == nil {
		return nil
	}

	for _, hook := range c.hooks[k] {
		if err := hook(); err != nil {
			return err
		}
	}
	return nil
}
