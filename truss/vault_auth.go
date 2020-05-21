package truss

// VaultAuth vault auth
type VaultAuth interface {
	Login() error
}
