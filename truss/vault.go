package truss

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	"github.com/phayes/freeport"
)

// VaultCmd wrapper for hashicorp vault
type VaultCmd struct {
	kubectl       *KubectlCmd
	auth          VaultAuth
	portForwarded bool
}

// Vault wrapper for hashicorp vault
func Vault(kubectl *KubectlCmd, auth VaultAuth) *VaultCmd {
	return &VaultCmd{
		kubectl: kubectl,
		auth:    auth,
	}
}

// Run run command
func (vault *VaultCmd) Run(args []string) ([]byte, error) {
	p, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	port := strconv.Itoa(p)

	if err := vault.kubectl.PortForward("8200", port, "vault", "service/vault"); err != nil {
		return nil, err
	}
	defer vault.kubectl.ClosePortForward()

	if vault.auth != nil {
		if err := vault.auth.Login(port); err != nil {
			return nil, err
		}
	}

	output, err := execVault(port, args...)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func execVault(port string, arg ...string) ([]byte, error) {
	cmd := exec.Command("vault", arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_ADDR=https://localhost:"+port, "VAULT_SKIP_VERIFY=true")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New(string(err.(*exec.ExitError).Stderr))
	}

	return output, nil
}
