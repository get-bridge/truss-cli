package truss

import (
	"fmt"
	"os"
)

// NetworkRouteProvider provides routes by provisioning them directly with the OS
type NetworkRouteProvider struct {
	Routes      []string
	OpenConnect *OpenConnect
}

// Setup creates the routes
func (n NetworkRouteProvider) Setup() error {
	if n.OpenConnect == nil {
		return nil
	}

	cmd := fmt.Sprintf("%s vpn openconnect-vpnc", os.Args[0])
	n.OpenConnect.Script = &cmd

	return nil
}

// Teardown destroys the routes
func (n NetworkRouteProvider) Teardown() error {
	return nil
}
