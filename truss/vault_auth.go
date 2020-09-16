package truss

// VaultAuth vault auth
type VaultAuth interface {
	Login(data interface{}, port string) (token string, err error)
	LoadCreds() (data interface{}, err error)
}
