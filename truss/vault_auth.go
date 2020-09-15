package truss

// VaultAuth vault auth
type VaultAuth interface {
	Login(data interface{}, port string) error
	LoadCreds() (interface{}, error)
}
