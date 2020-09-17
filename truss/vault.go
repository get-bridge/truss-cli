package truss

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
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
	GetWrappingToken() (string, error)
	GetMap(vaultPath string) (map[string]string, error)
	ListPath(vaultPath string) ([]string, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
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

func newVaultClient(port string) (*api.Client, error) {
	config := api.Config{Address: "https://localhost:" + port}
	config.ConfigureTLS(&api.TLSConfig{Insecure: true})
	return api.NewClient(&config)
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

	if err := vault.kubectl.PortForward("8200", port, "vault", "service/vault", vault.timeoutSeconds); err != nil {
		return "", err
	}

	vault.portForwarded = &port
	return port, nil
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
	// if we didn't start the port forward, don't close it
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	token, err := vault.getToken()
	if err != nil {
		return nil, err
	}

	output, err := vault.execVault(token, args...)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// GetToken gets a Vault Token
// Caller is responsible for closing port
func (vault *VaultCmdImpl) getToken() (string, error) {
	data, err := vault.auth.LoadCreds()
	if err != nil {
		return "", err
	}

	if vault.portForwarded == nil {
		_, err = vault.PortForward()
		if err != nil {
			return "", err
		}
	}

	return vault.auth.Login(data, *vault.portForwarded)
}

// GetWrappingToken gets a Vault wrapping token
// Caller is responsible for closing port
func (vault *VaultCmdImpl) GetWrappingToken() (string, error) {
	token, err := vault.Run([]string{"write", "-wrap-ttl=3m", "-field=wrapping_token", "-force", "auth/token/create"})
	return string(token), err
}

func (vault *VaultCmdImpl) execVault(token string, arg ...string) ([]byte, error) {
	cmd := exec.Command("vault", arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"VAULT_ADDR=https://localhost:"+*vault.portForwarded,
		"VAULT_SKIP_VERIFY=true",
		"VAULT_TOKEN="+token,
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Vault command [%v] failed: %v", arg, string(err.(*exec.ExitError).Stderr))
	}

	return output, nil
}

// Write to vault
func (vault *VaultCmdImpl) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	// if we didn't start the port forward, don't close it
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	token, err := vault.getToken()
	if err != nil {
		return nil, err
	}

	client, err := newVaultClient(*vault.portForwarded)
	client.SetToken(token)
	return client.Logical().Write(path, data)
}

// Encrypt bytes using transit key
func (vault *VaultCmdImpl) Encrypt(transitKeyName string, raw []byte) ([]byte, error) {
	if transitKeyName == "" {
		return nil, errors.New(("Must provide transitkey to encrypt"))
	}
	out, err := vault.Write("/transit/encrypt/"+transitKeyName, map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString(raw),
	})
	if err != nil {
		return nil, err
	}
	ciphertext, ok := out.Data["ciphertext"].(string)
	if !ok {
		return nil, errors.New("There was an error encyrpting your data")
	}
	return []byte(ciphertext), nil
}

// Decrypt bytes using transit key
func (vault *VaultCmdImpl) Decrypt(transitKeyName string, encrypted []byte) ([]byte, error) {
	if transitKeyName == "" {
		return nil, errors.New(("Must provide transitkey to decrypt"))
	}
	out, err := vault.Write("/transit/decrypt/"+transitKeyName, map[string]interface{}{
		"ciphertext": string(encrypted),
	})
	if err != nil {
		return nil, err
	}

	plaintext, ok := out.Data["plaintext"].(string)
	if !ok {
		return nil, errors.New("There was an error encyrpting your data")
	}

	return base64.StdEncoding.DecodeString(plaintext)
}

// GetMap returns a vaultPath as a map
func (vault *VaultCmdImpl) GetMap(vaultPath string) (map[string]string, error) {
	// if we didn't start the port forward, don't close it
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	token, err := vault.getToken()
	if err != nil {
		return nil, err
	}

	client, err := newVaultClient(*vault.portForwarded)
	client.SetToken(token)
	// TODO SO GROSS
	if !strings.Contains(vaultPath, "secret/data") {
		vaultPath = strings.Replace(vaultPath, "secret/", "secret/data/", 1)
	}
	out, err := client.Logical().Read(vaultPath)
	if err != nil {
		return nil, err
	}

	data, ok := out.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse path [%v] data: %v", vaultPath, out.Data)
	}

	secretStringMap := map[string]string{}
	for k, v := range data {
		vString, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("failed to parse path [%v] data: %v", vaultPath, out.Data)
		}
		secretStringMap[k] = vString
	}

	return secretStringMap, nil
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
