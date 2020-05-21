package cmd

import (
	"os"

	"github.com/instructure/truss-cli/truss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Proxy commands to kubectl",
	// Long: `TODO...`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, err := getKubeconfig(cmd, args)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}

		output, err := truss.Kubectl(kubeconfig).Run(args...)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(kubectlCmd)

	kubectlCmd.Flags().SetInterspersed(false)
}
