package truss

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type VaultServer struct {
	Addr  string
	Token string
	cmd   *exec.Cmd
}

var server *VaultServer = &VaultServer{
	Addr:  "http://localhost:8200",
	Token: "",
}

func (v *VaultServer) Start() error {
	v.cmd = exec.Command("vault", "server", "-dev", fmt.Sprintf("-address=%s", v.Addr))

	// Attach to Vault's stdout and setup a scanner to read it
	stdout, err := v.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	// Start Vault
	err = v.cmd.Start()
	if err != nil {
		return err
	}

	// Scan stdout until we read the root token
	for scanner.Scan() {
		output := scanner.Text()

		re := regexp.MustCompile(`Root Token: (.*)$`)

		match := re.FindStringSubmatch(output)

		if len(match) > 0 {
			v.Token = match[1]
			break
		}
	}

	return nil
}

// Send Vault a KILL signal and wait for it to stop
func (v *VaultServer) Stop() {
	v.cmd.Process.Kill()
	v.cmd.Wait()
}

// Initialize and authenticate a Vault client
func (v *VaultServer) Client() (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress(v.Addr),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	if err := client.SetToken(v.Token); err != nil {
		return nil, err
	}

	return client, nil
}

func SetupVaultServer() error {
	err := server.Start()
	if err != nil {
		return fmt.Errorf("failed to start Vault server: %s", err)
	}

	client, err := server.Client()
	if err != nil {
		return fmt.Errorf("failed to initialize Vault client: %s", err)
	}

	// Create KV V2 mount
	_, err = client.System.MountsEnableSecretsEngine(
		context.Background(),
		"kv",
		schema.MountsEnableSecretsEngineRequest{
			Type: "kv",
			Options: map[string]interface{}{
				"version": "2",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to enable kv engine: %s", err)
	}

	// Create transit mount
	_, err = client.System.MountsEnableSecretsEngine(
		context.Background(),
		"transit",
		schema.MountsEnableSecretsEngineRequest{
			Type: "transit",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to enable transit engine: %s", err)
	}

	return nil
}

func TeardownVaultServer() {
	server.Stop()
}

// Initialize an authenticated VaultCmd
func createTestVault(t *testing.T) *VaultCmd {
	t.Helper()

	vault := VaultWithToken("", server.Token)
	vault.addr = server.Addr

	return vault
}

func TestMain(m *testing.M) {
	err := SetupVaultServer()
	if err != nil {
		fmt.Printf("Failed to setup Vault server: %s\n", err)
		os.Exit(1)
	}

	// Run tests
	exitVal := m.Run()

	TeardownVaultServer()
	os.Exit(exitVal)
}
