package truss

// VaultAuth vault auth
type VaultAuth interface {
	Login(port string) error
}
