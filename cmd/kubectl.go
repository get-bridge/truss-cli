package cmd

import (
	"github.com/get-bridge/truss-cli/truss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Proxy commands to kubectl",
	// Long: `TODO...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, err := getKubeconfig()
		if err != nil {
			return err
		}

		output, err := truss.Kubectl(kubeconfig).Run(args...)
		if err != nil {
			return err
		}
		log.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(kubectlCmd)

	kubectlCmd.Flags().SetInterspersed(false)
}
