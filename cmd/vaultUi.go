package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
)

var vaultUICmd = &cobra.Command{
	Use:   "ui",
	Short: "Open the Vault UI in your browser",
	Long: `This is useful when your vault is not exposed publicly.
As it will port-forward to the service, authenticate with aws auth,
and open the UI in your browser.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, err := getKubeconfig()
		if err != nil {
			return err
		}

		vault := truss.Vault(kubeconfig, getVaultAuth())

		port, err := vault.PortForward()
		if err != nil {
			return err
		}
		defer vault.ClosePortForward()

		token, err := vault.GetWrappingToken()
		if err != nil {
			return err
		}
		vaultURL := fmt.Sprintf("https://localhost:%s/ui/vault/auth?with=token&wrapped_token=%s", port, token)
		log.Printf("Opening Vault UI at %s", vaultURL)

		openbrowser(vaultURL)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for range c {
			log.Println("Received SIGINT, cleaning up...")
			return nil
		}
		return nil
	},
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	vaultCmd.AddCommand(vaultUICmd)
}
