package cmd

import (
	"fmt"
	"os"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// wrapCmd represents the wrap command
var wrapCmd = &cobra.Command{
	Use: "wrap",
	Long: `
Sets KUBECONFIG and then executes the subcommand:

	$ truss wrap -e edge-cmh -- printenv

This allows you to do this:

	$ truss wrap -e edge-cmh -- k9s
`,
	Short: "Wraps a subcommand with truss environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return err
		}
		environments := viper.GetStringMapString("environments")

		if env == "" {
			return fmt.Errorf("-e flag is required. Options: %v", getEnvironmentKeys(environments))
		}

		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		bin := args[0]
		binargs := args[1:]

		kubeconfig, err := getKubeconfig()

		if err != nil {
			return err
		}

		input := &truss.WrapInput{
			Kubeconfig: kubeconfig,
			Stdout:     cmd.OutOrStdout(),
			Stdin:      os.Stdin,
			Stderr:     cmd.ErrOrStderr(),
		}

		err = truss.Wrap(input, bin, binargs...)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(wrapCmd)
}
