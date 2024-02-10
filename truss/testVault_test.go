package truss

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type VaultDevServer struct {
	Addr  string
	Token string
	cmd   *exec.Cmd
}

var port = 8200

func NewVaultDevServer() *VaultDevServer {
	listenAddr := fmt.Sprintf("http://localhost:%s", strconv.Itoa(port))
	return &VaultDevServer{
		Addr:  listenAddr,
		Token: "",
	}
}

func (v *VaultDevServer) Start() error {
	v.cmd = exec.Command("vault", "server", "-dev", fmt.Sprintf("-address=%s", v.Addr))

	stdout, err := v.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	err = v.cmd.Start()
	if err != nil {
		return err
	}

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

func (v *VaultDevServer) Stop() {
	v.cmd.Process.Kill()
	v.cmd.Wait()
}

func (v *VaultDevServer) Client() (*vault.Client, error) {
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

// creates test vault server
func createTestVault(t *testing.T) (*VaultCmd, *VaultDevServer) {
	t.Helper()

	server := NewVaultDevServer()

	err := server.Start()
	if err != nil {
		t.Fatalf("failed to start Vault server: %s", err)
	}

	client, err := server.Client()
	if err != nil {
		t.Fatal("failed to initialize Vault client")
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
		t.Fatal("failed to enable kv engine")
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
		t.Fatal("failed to enable transit engine")
	}

	vault := VaultWithToken("", server.Token)
	vault.addr = server.Addr

	timeout := 0
	for timeout < 20 {
		_, err := vault.ListPath("kv/metadata")
		if err == nil {
			return vault, server
		}
		time.Sleep(100 * time.Millisecond)
		timeout++
	}
	t.Fatal("vault engine not started")
	return nil, nil
}
