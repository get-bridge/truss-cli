package truss

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/vault"
)

var vaultAddr = ""
var vaultToken = "this-is-the-root-token"

// Initialize an authenticated VaultCmd
func createTestVault(t *testing.T) *VaultCmd {
	t.Helper()

	vault := VaultWithToken("", vaultToken)
	vault.addr = vaultAddr

	return vault
}

// This wraps our entire test run to:
// 1. Start and configure Vault with required backends
// 2. Execute tests
// 3. Teardown Vault
// 4. Exit with the exit value as determined by the tests
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start Vault server with kv2 and transit backends enabled
	vaultContainer, err := vault.RunContainer(ctx, vault.WithToken(vaultToken), vault.WithInitCommand(
		"secrets enable -version=2 -path=kv kv",
		"secrets enable -path=transit transit",
	))

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	vaultAddr, err = vaultContainer.HttpHostAddress(ctx)
	if err != nil {
		log.Fatalf("failed to get Vault address: %s", err)
	}

	// Run tests
	exitVal := m.Run()

	// Teardown Vault server
	if err := vaultContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(exitVal)
}
