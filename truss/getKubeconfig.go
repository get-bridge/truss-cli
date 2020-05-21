package truss

// GetKubeconfigCmd command for managing kubeconfigs
type GetKubeconfigCmd interface {
	Fetch() error
}
