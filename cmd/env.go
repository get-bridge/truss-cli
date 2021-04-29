package cmd

import (
	"fmt"

	"github.com/get-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EnvCmd represents the env command
var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Display the commands to set up the shell environment for truss",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return err
		}
		kubeconfigs := viper.GetStringMap("environments")
		kubeDir, err := getKubeDir()
		if err != nil {
			return err
		}
		input := &truss.EnvInput{
			Env:         env,
			Kubeconfigs: kubeconfigs,
			KubeDir:     kubeDir,
		}
		environmentVars, err := truss.Env(input)
		if err != nil {
			return err
		}

		output := environmentVars.BashFormat(env)

		fmt.Fprintln(cmd.OutOrStdout(), output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(EnvCmd)
}
