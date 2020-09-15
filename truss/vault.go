package truss

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/phayes/freeport"
	"gopkg.in/yaml.v2"
)

// VaultCmd Interface for interacting with vault
type VaultCmd interface {
	PortForward() (string, error)
	ClosePortForward() error
	Run(args []string) ([]byte, error)
	Decrypt(transitKeyName string, encrypted []byte) ([]byte, error)
	Encrypt(transitKeyName string, raw []byte) ([]byte, error)
	GetToken() (string, error)
	GetMap(vaultPath string) (map[string]string, error)
	ListPath(vaultPath string) ([]string, error)
}

// VaultCmdImpl wrapper implementation for hashicorp vault
type VaultCmdImpl struct {
	kubectl        *KubectlCmd
	auth           VaultAuth
	portForwarded  *string
	timeoutSeconds int
}

// Vault wrapper for hashicorp vault
func Vault(kubeconfig string, auth VaultAuth) VaultCmd {
	return &VaultCmdImpl{
		kubectl:        Kubectl(kubeconfig),
		auth:           auth,
		timeoutSeconds: 15,
	}
}

// PortForward instantiates a port-forward for Vault
func (vault *VaultCmdImpl) PortForward() (string, error) {
	if vault.portForwarded != nil {
		return *vault.portForwarded, nil
	}

	p, err := freeport.GetFreePort()
	if err != nil {
		return "", err
	}
	port := strconv.Itoa(p)
	vault.portForwarded = &port

	return port, vault.kubectl.PortForward("8200", port, "vault", "service/vault", vault.timeoutSeconds)
}

// ClosePortForward closes the port forward, if any
func (vault *VaultCmdImpl) ClosePortForward() error {
	if vault.portForwarded == nil {
		return nil
	}
	vault.portForwarded = nil
	return vault.kubectl.ClosePortForward()
}

// Run run command
func (vault *VaultCmdImpl) Run(args []string) ([]byte, error) {
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

// GetToken gets a Vault Token
func (vault *VaultCmdImpl) GetToken() (string, error) {
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
func (vault *VaultCmdImpl) Encrypt(transitKeyName string, raw []byte) ([]byte, error) {
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
func (vault *VaultCmdImpl) Decrypt(transitKeyName string, encrypted []byte) ([]byte, error) {
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

// GetMap returns a vaultPath as a map
func (vault *VaultCmdImpl) GetMap(vaultPath string) (map[string]string, error) {
	get, err := vault.Run([]string{
		"kv",
		"get",
		"-format=yaml",
		vaultPath,
	})
	if err != nil {
		return nil, err
	}

	getData := struct {
		Data struct {
			Data map[string]string `yaml:"data"`
		} `yaml:"data"`
	}{}
	if err := yaml.NewDecoder(bytes.NewReader(get)).Decode(&getData); err != nil {
		return nil, err
	}

	return getData.Data.Data, nil
}

// ListPath returns a vaultPath as a map
func (vault *VaultCmdImpl) ListPath(vaultPath string) ([]string, error) {
	list, err := vault.Run([]string{
		"kv",
		"list",
		"-format=yaml",
		vaultPath,
	})
	if err != nil {
		return nil, err
	}

	secrets := []string{}
	if err := yaml.NewDecoder(bytes.NewReader(list)).Decode(&secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}
