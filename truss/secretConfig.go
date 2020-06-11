package truss

// SecretConfig represents a desired Vault synchronization
type SecretConfig struct {
	Name       string `yaml:"name"`
	Kubeconfig string `yaml:"kubeconfig"`
	VaultPath  string `yaml:"vaultPath"`
	FilePath   string `yaml:"filePath"`
}
