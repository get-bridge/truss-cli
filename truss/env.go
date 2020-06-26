package truss

import (
	"bytes"
	"fmt"
)

// EnvInput input
type EnvInput struct {
	Env         string
	Kubeconfigs map[string]interface{}
	KubeDir     string
}

// Env returns string to eval that configures shell environment variables
func Env(input *EnvInput) (string, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	kubeconfigName := input.Kubeconfigs[input.Env]

	if kubeconfigName == nil {
		return "", fmt.Errorf("No kubeconfig found for env %s", input.Env)
	}

	kubeconfig := fmt.Sprintf("%s%s", input.KubeDir, kubeconfigName)
	buffer.WriteString(fmt.Sprintf("export KUBECONFIG=%s\n", kubeconfig))
	buffer.WriteString("# Run this command to configure your shell:\n")
	buffer.WriteString(fmt.Sprintf("# eval \"$(truss env -e %s)\"", input.Env))

	return buffer.String(), nil
}
