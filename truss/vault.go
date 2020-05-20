package truss

import (
	"errors"
	"os"
	"os/exec"
)

// VaultCmd wrapper for hashicorp vault
type VaultCmd struct {
	kubectl *KubectlCmd
	auth    VaultAuth
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
	if err := vault.kubectl.PortForward("8200", "vault", "service/vault"); err != nil {
		return nil, err
	}
	defer vault.kubectl.ClosePortForward()

	if vault.auth != nil {
		if err := vault.auth.Login(); err != nil {
			return nil, err
		}
	}

	output, err := execVault(args...)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func execVault(arg ...string) ([]byte, error) {
	cmd := exec.Command("vault", arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_ADDR=https://localhost:8200", "VAULT_SKIP_VERIFY=true")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New(string(err.(*exec.ExitError).Stderr))
	}

	return output, nil
}
