package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/atotto/clipboard"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pbcopyCmd = &cobra.Command{
	Use:   "pbcopy {content}",
	Short: "Wraps a secret for sharing",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("TRUSS_ENV", "edge-cmh")
		kubeconfig, err := getKubeconfig()
		if err != nil {
			return err
		}

		content := bytes.NewBuffer(nil)
		if len(args) == 1 {
			content.Write([]byte(args[0]))
		} else {
			io.Copy(content, os.Stdin)
		}
		if content.Len() == 0 {
			return errors.New("No content provided")
		}

		ttl, _ := cmd.Flags().GetInt("ttl")
		token, err := truss.Vault(kubeconfig, getVaultAuth()).Run([]string{
			"write",
			"-field=wrapping_token",
			fmt.Sprintf("-wrap-ttl=%dm", ttl),
			"/sys/wrapping/wrap",
			fmt.Sprintf("pb=%s", content),
		})
		if err != nil {
			return err
		}

		copy := fmt.Sprintf("\"truss pbpaste %s\" has been copied to your clipboard!", token)
		fmt.Println(copy)
		return clipboard.WriteAll(copy)
	},
}

func init() {
	rootCmd.AddCommand(pbcopyCmd)

	pbcopyCmd.Flags().Int("ttl", 30, "TTL of the wrapped secret in minutes")
}
