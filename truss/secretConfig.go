package truss

// SecretConfig interface for secret config
// can be a file or directory or anything!
type SecretConfig interface {
	Name() string
	Kubeconfig() string
	VaultPath() string
	existsOnDisk() bool
	getDecryptedFromDisk(vault *VaultCmd, transitKeyName string) ([]byte, error)
	saveToDisk(vault *VaultCmd, transitKeyName string, raw []byte) error
	saveToDiskFromVault(vault *VaultCmd, transitKeyName string) error
	writeToVault(vault *VaultCmd, transitKeyName string) error
}
