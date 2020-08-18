package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
)

// VpnCmd manages VPN Connections
var VpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Manage VPN and Network Tunnels for accessing Truss",
	Long: `
Connects to the VPN and configures network routing for accessing Truss clusters using Cisco Anyconnect or sshuttle. This command requires openconnect and optionally sshuttle.

If network traffic doesn't route as you'd expect it to, try adding the -s flag.

Cisco Anyconnect:
During the connection process, the openconnect process will re-invoke the CLI to initialize the connection. This creates the required environment variables for the standard vpnc-script for traffic forwarding.

There is no teardown process when using this method.

sshuttle:
After the connection is established, sshuttle will run as a daemon to forward traffic.

The "truss vpn stop" command will kill this sshuttle process by pid.

By default, all Kubernetes clusters defined in your ~/.truss.yaml will be routed through the VPN. You can specify vpn.forwardHosts and vpn.forwardIPs to add additional static routes as well.
`,
}

// VpnStartCmd starts a VPN connection
var VpnStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a Truss VPN and Network Tunnel",
	RunE: func(cmd *cobra.Command, args []string) error {
		oc := getOC(cmd)
		return oc.Start()
	},
}

// VpnStopCmd starts a VPN connection
var VpnStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a Truss VPN and Network Tunnel",
	RunE: func(cmd *cobra.Command, args []string) error {
		oc := getOC(cmd)
		return oc.Stop()
	},
}

// VpnOpenConnectVpncCmd does openconnect things for routes!
var VpnOpenConnectVpncCmd = &cobra.Command{
	Use:    "openconnect-vpnc",
	Short:  "Used by the NetworkRouteProvider to initialize the VPN Connection",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		splitinc := os.Getenv("CISCO_SPLIT_INC")
		inc, err := strconv.Atoi(splitinc)
		if err != nil {
			inc = 0
		}

		c := exec.Command(vpncScript)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Env = os.Environ()

		for _, r := range getRoutes(cmd) {
			c.Env = append(c.Env,
				fmt.Sprintf("CISCO_SPLIT_INC_%d_ADDR=%s", inc, r),
				fmt.Sprintf("CISCO_SPLIT_INC_%d_MASK=%s", inc, "255.255.255.255"),
				fmt.Sprintf("CISCO_SPLIT_INC_%d_MASKLEN=%s", inc, "32"),
			)
			inc++
		}
		c.Env = append(c.Env, fmt.Sprintf("CISCO_SPLIT_INC=%d", inc))

		return c.Run()
	},
}

func init() {
	VpnCmd.AddCommand(VpnStartCmd)
	VpnCmd.AddCommand(VpnStopCmd)
	VpnCmd.AddCommand(VpnOpenConnectVpncCmd)
	rootCmd.AddCommand(VpnCmd)

	u, _ := user.Current()
	VpnCmd.PersistentFlags().String("server", "vpn.instructure.com", "Cisco Anyconnect Group to join")
	viper.BindPFlag("vpn.server", VpnCmd.PersistentFlags().Lookup("server"))
	VpnCmd.PersistentFlags().StringP("user", "u", u.Username, "User to use to connect to the VPN")
	viper.BindPFlag("vpn.user", VpnCmd.PersistentFlags().Lookup("user"))
	VpnCmd.PersistentFlags().StringP("authgroup", "g", "Employee_VPN", "Cisco Anyconnect Group to join")
	viper.BindPFlag("vpn.authgroup", VpnCmd.PersistentFlags().Lookup("authgroup"))
	VpnCmd.PersistentFlags().BoolP("sshuttle", "s", false, "Use sshuttle instead of Split Tunnel")
	VpnCmd.PersistentFlags().String("ssh-host", "10.0.34.70", "SSH Host to use for sshuttle implementation")
}

func getOC(cmd *cobra.Command) *truss.OpenConnect {
	oc := truss.NewOpenConnect(
		viper.GetString("vpn.user"),
		viper.GetString("vpn.server"),
		viper.GetString("vpn.authgroup"),
	)

	if s, err := cmd.Flags().GetBool("sshuttle"); err == nil && s {
		h, _ := cmd.Flags().GetString("ssh-host")
		rp := truss.SshuttleRouteProvider{
			Host:   h,
			Routes: getRoutes(cmd),
		}

		oc.AddHook(truss.OpenConnectConnected, rp.Setup)
		oc.AddHook(truss.OpenConnectDisconnected, rp.Teardown)
	} else {
		rp := truss.NetworkRouteProvider{
			OpenConnect: oc,
		}
		oc.AddHook(truss.OpenConnectConnecting, rp.Setup)
	}

	return oc
}

func getRoutes(cmd *cobra.Command) []string {
	kd, err := getKubeDir()
	if err != nil {
		return nil
	}

	hosts := viper.GetStringSlice("vpn.forwardHosts")
	for name := range viper.GetStringMap("environments") {
		f := filepath.Join(kd, viper.GetString("environments."+name))
		c, err := clientcmd.LoadFromFile(f)
		if err != nil {
			return nil
		}

		for _, cluster := range c.Clusters {
			hosts = appendDedupe(hosts, strings.Replace(cluster.Server, "https://", "", 1))
		}
	}

	routes := viper.GetStringSlice("vpn.forwardIPs")
	for _, h := range hosts {
		addrs, err := net.LookupHost(h)
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			routes = appendDedupe(routes, addr)
		}
	}

	return routes
}

func appendDedupe(dst []string, v string) []string {
	a := true
	for _, d := range dst {
		if d == v {
			a = false
			break
		}
	}
	if a {
		dst = append(dst, v)
	}
	return dst
}
