package truss

import (
	"fmt"
	"os/exec"
)

// Setup set up
func Setup(dependencies *[]string) error {
	if err := checkDependencies(dependencies); err != nil {
		return err
	}
	return nil
}

func checkDependencies(dependencies *[]string) error {
	missingDependencies := []string{}
	for _, d := range *dependencies {
		if _, err := exec.LookPath(d); err != nil {
			missingDependencies = append(missingDependencies, d)
		}
	}
	if len(missingDependencies) > 0 {
		return fmt.Errorf("missing dependencies: %s", missingDependencies)
	}
	return nil
}
