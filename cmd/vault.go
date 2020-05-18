package cmd

import (
	"fmt"
	"os"

	"github.com/instructure/truss-cli/truss"
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "A wrapper around hashicorp vault",
	Long: `This is useful when your vault is not exposed publicly.
As it will port-forward to the service and call aws auth`,
	Run: func(cmd *cobra.Command, args []string) {

		env, err := cmd.Flags().GetString("env")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// println("env", env)
		// TODO pull from config
		contexts := map[string]string{
			"": "cluster-admin@truss-nonprod-cmh-shared-cluster",
		}
		context := contexts[env]
		err = truss.Vault(context).Run(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
}
