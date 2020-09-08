package truss

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/phayes/freeport"
)

// VaultCmd wrapper for hashicorp vault
type VaultCmd struct {
	kubectl       *KubectlCmd
	auth          VaultAuth
	portForwarded *string
}

// Vault wrapper for hashicorp vault
func Vault(kubectl *KubectlCmd, auth VaultAuth) *VaultCmd {
	return &VaultCmd{
		kubectl: kubectl,
		auth:    auth,
	}
}

// PortForward instantiates a port-forward for Vault
func (vault *VaultCmd) PortForward() (string, error) {
	if vault.portForwarded != nil {
		return *vault.portForwarded, nil
	}

	p, err := freeport.GetFreePort()
	if err != nil {
		return "", err
	}
	port := strconv.Itoa(p)
	vault.portForwarded = &port

	return port, vault.kubectl.PortForward("8200", port, "vault", "service/vault")
}

// ClosePortForward closes the port forward, if any
func (vault *VaultCmd) ClosePortForward() error {
	if vault.portForwarded == nil {
		return nil
	}
	vault.portForwarded = nil
	return vault.kubectl.ClosePortForward()
}

// Run run command
func (vault *VaultCmd) Run(args []string) ([]byte, error) {
	var port string
	var err error

	if vault.portForwarded != nil {
		port = *vault.portForwarded
	} else {
		port, err = vault.PortForward()
		if err != nil {
			return nil, err
		}
		defer vault.ClosePortForward()
	}

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

// GetToken gets a Vaut Token
func (vault *VaultCmd) GetToken() (string, error) {
	if vault.auth == nil {
		return "", errors.New("vault auth not configured")
	}

	var port string
	var err error
	if vault.portForwarded != nil {
		port = *vault.portForwarded
	} else {
		port, err = vault.PortForward()
		if err != nil {
			return "", err
		}
		defer vault.ClosePortForward()
	}

	if err := vault.auth.Login(port); err != nil {
		return "", err
	}

	out, err := vault.Run([]string{"write", "-wrap-ttl=3m", "-field=wrapping_token", "-force", "auth/token/create"})

	return string(out), err
}

func execVault(port string, arg ...string) ([]byte, error) {
	cmd := exec.Command("vault", arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "VAULT_ADDR=https://localhost:"+port, "VAULT_SKIP_VERIFY=true")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Vault command failed: %v", string(err.(*exec.ExitError).Stderr))
	}

	return output, nil
}

// Encrypt shit
func (vault *VaultCmd) Encrypt(transitKeyName string, raw []byte) ([]byte, error) {
	if transitKeyName == "" {
		return nil, errors.New(("Must provide transitkey to encrypt"))
	}
	out, err := vault.Run([]string{
		"write",
		"-field=ciphertext",
		"transit/encrypt/" + transitKeyName,
		"plaintext=" + base64.StdEncoding.EncodeToString(raw),
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Decrypt shit
func (vault *VaultCmd) Decrypt(transitKeyName string, encrypted []byte) ([]byte, error) {
	if transitKeyName == "" {
		return nil, errors.New(("Must provide transitkey to decrypt"))
	}
	out, err := vault.Run([]string{
		"write",
		"-field=plaintext",
		"transit/decrypt/" + transitKeyName,
		"ciphertext=" + string(encrypted),
	})

	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(out))
}
