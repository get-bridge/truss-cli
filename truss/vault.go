package truss

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/phayes/freeport"
)

// VaultCmd wrapper implementation for hashicorp vault
type VaultCmd struct {
	kubectl *KubectlCmd
	auth    VaultAuth
	token   string
	addr    string

	portForwarded  *string
	timeoutSeconds int
}

// Vault wrapper for hashicorp vault
func Vault(kubeconfig string, auth VaultAuth) *VaultCmd {
	return &VaultCmd{
		kubectl:        Kubectl(kubeconfig),
		auth:           auth,
		timeoutSeconds: 15,
	}
}

// VaultWithToken wrapper for hashicorp vault with token for auth
func VaultWithToken(kubeconfig string, authToken string) *VaultCmd {
	return &VaultCmd{
		kubectl:        Kubectl(kubeconfig),
		token:          authToken,
		timeoutSeconds: 15,
	}
}

func newVaultClient(addr string) (*api.Client, error) {
	config := api.Config{Address: addr}
	config.ConfigureTLS(&api.TLSConfig{Insecure: true})
	return api.NewClient(&config)
}

func (vault *VaultCmd) newVaultClientWithToken() (*api.Client, error) {
	token, err := vault.getToken()
	if err != nil {
		return nil, err
	}

	client, err := newVaultClient(vault.vaultAddr())
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return client, nil
}

func (vault *VaultCmd) vaultAddr() string {
	addr := vault.addr
	if addr == "" {
		addr = "https://localhost:" + *vault.portForwarded
	}
	return addr
}

// PortForward instantiates a port-forward for Vault
func (vault *VaultCmd) PortForward() (string, error) {
	_, err := vault.auth.LoadCreds()
	if err != nil {
		return "", err
	}

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
func (vault *VaultCmd) ClosePortForward() error {
	if vault.portForwarded == nil {
		return nil
	}
	vault.portForwarded = nil
	return vault.kubectl.ClosePortForward()
}

// Run run command
func (vault *VaultCmd) Run(args []string) ([]byte, error) {
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
func (vault *VaultCmd) getToken() (string, error) {
	if vault.token != "" {
		return vault.token, nil
	}
	if vault.auth == nil {
		return "", errors.New("vault auth must be configured to get token")
	}
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

	return vault.auth.Login(data, vault.vaultAddr())
}

// GetWrappingToken gets a Vault wrapping token
// Caller is responsible for closing port
func (vault *VaultCmd) GetWrappingToken() (string, error) {
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	client, err := vault.newVaultClientWithToken()
	if err != nil {
		return "", err
	}
	client.SetWrappingLookupFunc(func(string, string) string { return "3m" })
	out, err := client.Logical().Write("auth/token/create", map[string]interface{}{})
	if err != nil {
		return "", err
	}

	return out.WrapInfo.Token, nil
}

func (vault *VaultCmd) execVault(token string, arg ...string) ([]byte, error) {
	cmd := exec.Command("vault", arg...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"VAULT_ADDR="+vault.vaultAddr(),
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
func (vault *VaultCmd) Write(vaultPath string, data map[string]interface{}) (*api.Secret, error) {
	// if we didn't start the port forward, don't close it
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	client, err := vault.newVaultClientWithToken()
	if err != nil {
		return nil, err
	}
	return client.Logical().Write(vaultPath, data)
}

// Encrypt bytes using transit key
func (vault *VaultCmd) Encrypt(transitKeyName string, raw []byte) ([]byte, error) {
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
func (vault *VaultCmd) Decrypt(transitKeyName string, encrypted []byte) ([]byte, error) {
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
func (vault *VaultCmd) GetMap(vaultPath string) (map[string]interface{}, error) {
	// if we didn't start the port forward, don't close it
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	client, err := vault.newVaultClientWithToken()
	if err != nil {
		return nil, err
	}
	out, err := client.Logical().Read(vaultPath)
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, nil
	}

	data, ok := out.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse path [%v] data: %v", vaultPath, out.Data)
	}

	return data, nil
}

// ListPath returns a vaultPath as a map
func (vault *VaultCmd) ListPath(vaultPath string) ([]string, error) {
	// if we didn't start the port forward, don't close it
	if vault.portForwarded == nil {
		defer vault.ClosePortForward()
	}

	client, err := vault.newVaultClientWithToken()
	if err != nil {
		return nil, err
	}
	out, err := client.Logical().List(vaultPath)
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, nil
	}

	keys, ok := out.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to list path [%v] data: %v", vaultPath, out.Data)
	}

	keysString := []string{}
	for _, k := range keys {
		kString, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("failed to list path [%v] data: %v", vaultPath, out.Data)
		}
		keysString = append(keysString, kString)
	}

	return keysString, nil
}

func kv2DataPath(vaultPath string) string {
	split := strings.Split(vaultPath, "/")
	if len(split) > 1 && split[1] == "data" {
		return vaultPath
	}
	split = append([]string{split[0], "data"}, split[1:]...)
	return strings.Join(split, "/")
}

func kv2MetadataPath(vaultPath string) string {
	split := strings.Split(vaultPath, "/")
	if len(split) > 1 && split[1] == "metadata" {
		return vaultPath
	}
	split = append([]string{split[0], "metadata"}, split[1:]...)
	return strings.Join(split, "/")
}
