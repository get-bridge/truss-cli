package truss

// SecretConfig interface for secret config
// can be a file or directory or anything!
type SecretConfig interface {
	Name() string
	Kubeconfig() string
	VaultPath() string
	existsOnDisk() bool
	getDecryptedFromDisk(vault VaultCmd, transitKeyName string) ([]byte, error)
	getMapFromDisk(vault VaultCmd, transitKeyName string) (map[string]map[string]string, error)
	encryptAndSaveToDisk(vault VaultCmd, transitKeyName string, raw []byte) error
	writeMapToDisk(vault VaultCmd, transitKeyName string, secrets map[string]map[string]string) error
	write(vault VaultCmd, transitKeyName string) error
}
