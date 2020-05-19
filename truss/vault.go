package truss

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// VaultCmd wrapper for hashicorp vault
type VaultCmd struct {
	Context string
}

// Vault wrapper for hashicorp vault
func Vault(context string) *VaultCmd {
	return &VaultCmd{
		Context: context,
	}
}

// Run run command
func (vault *VaultCmd) Run(args []string) error {
	kubectl, err := Kubectl(vault.Context)
	if err != nil {
		return err
	}

	err = kubectl.PortForward("8200", "vault", "service/vault")
	if err != nil {
		return err
	}
	defer kubectl.ClosePortForward()
	// TODO make configurable
	// rapture assume arn:aws:iam::127178877223:role/xacct/ops-admin
	err = login("login", "-method=aws", "role=admin")
	if err != nil {
		return err
	}

	output, err := execVault(args...)
	if err != nil {
		return err
	}
	log.Println(string(output))
	return nil
}

func login(arg ...string) error {
	_, err := execVault(arg...)
	return err
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
