package truss

import (
	"fmt"
	"os/exec"

	"github.com/spf13/viper"
)

// Setup set up
func Setup() error {
	if err := checkDependencies(); err != nil {
		return err
	}
	return nil
}

func checkDependencies() error {
	missingDependencies := []string{}
	for _, dPtr := range viper.Get("dependencies").([]interface{}) {
		d := dPtr.(string)
		if _, err := exec.LookPath(d); err != nil {
			missingDependencies = append(missingDependencies, d)
		}
	}
	if len(missingDependencies) > 0 {
		return fmt.Errorf("missing dependencies: %s", missingDependencies)
	}
	return nil
}
