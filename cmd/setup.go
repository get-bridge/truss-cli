package cmd

import (
	"fmt"
	"os"

	"github.com/instructure/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Ensure your truss cli is ready",
	Long: `Dependencies are configured using 'dependencies' field in configfile.

dependencies:
  - kubectl
  - sshuttle
	- vault
`,
	Run: func(cmd *cobra.Command, args []string) {
		dependenciesPtr, ok := viper.Get("dependencies").([]interface{})
		if !ok {
			fmt.Println("invalid dependency configuration")
			os.Exit(1)
		}
		dependencies := []string{}
		for _, d := range dependenciesPtr {
			dependencyStr, ok := d.(string)
			if !ok {
				fmt.Println("invalid dependency type", d)
				os.Exit(1)
			}
			dependencies = append(dependencies, dependencyStr)
		}
		if err := truss.Setup(&dependencies); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
