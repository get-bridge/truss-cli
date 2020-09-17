package truss

import (
	"testing"
	"time"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/transit"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	hashivault "github.com/hashicorp/vault/vault"
)

// creates test vault server
func createTestVault(t *testing.T) *VaultCmd {
	t.Helper()

	coreConfig := &hashivault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv":      kv.Factory,
			"transit": transit.Factory,
		},
	}
	cluster := hashivault.NewTestCluster(t, coreConfig, &hashivault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	// Create KV V2 mount
	sys := cluster.Cores[0].Client.Sys()
	if err := sys.Mount("kv", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}
	// Create transit mount
	if err := sys.Mount("transit", &api.MountInput{
		Type: "transit",
	}); err != nil {
		t.Fatal(err)
	}

	vault := VaultWithToken("", cluster.Cores[0].Client.Token())
	vault.addr = cluster.Cores[0].Client.Address()

	timeout := 0
	for timeout < 20 {
		_, err := vault.ListPath("kv/metadata")
		if err == nil {
			return vault
		}
		time.Sleep(100 * time.Millisecond)
		timeout += 1
	}
	t.Fatal("vault engine not started")
	return nil
}
