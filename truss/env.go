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

// EnvironmentVars key/value pairs of environment variables that should be set in the shell
type EnvironmentVars struct {
	Kubeconfig string
}

// Env configures environment variables that should be set in the bash shell
func Env(input *EnvInput) (EnvironmentVars, error) {
	kubeconfigName := input.Kubeconfigs[input.Env]

	if kubeconfigName == nil {
		return EnvironmentVars{}, fmt.Errorf("No kubeconfig found for env %s", input.Env)
	}

	kubeconfig := fmt.Sprintf("%s%s", input.KubeDir, kubeconfigName)

	environmentVars := EnvironmentVars{
		Kubeconfig: kubeconfig,
	}

	return environmentVars, nil
}

// BashFormat formats environment variables for bash
func (environmentVars *EnvironmentVars) BashFormat(env string) string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("export KUBECONFIG=%s\n", environmentVars.Kubeconfig))
	buffer.WriteString("# Run this command to configure your shell:\n")
	buffer.WriteString(fmt.Sprintf("# eval \"$(truss env -e %s)\"", env))

	return buffer.String()
}
