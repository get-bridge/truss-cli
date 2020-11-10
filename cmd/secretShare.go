package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretShareCmd = &cobra.Command{
	Use:   "share {content}",
	Short: "Shares a secret via a Vault Wrapped Token",
	Long: `
The "share" command can be used to share sensitive information with the team
without leaking it to third-party services such as Slack. The command can
share the content via argument or Stdin.

	$ truss secret share "my sensitive secret"
	$ cat ./privkey.pem | truss secret share

This will then return a wrapped token to Stdout. You can conveniently copy this
output by piping it to the "pbcopy" command on MacOS:

	$ truss secret share "my sensitive secret" | pbcopy

You can specify how long the secret will be valid for. The default is 30 min:

	$ truss secret share --ttl=1 "better hurry!"
	
With the wrapped token copied, you can then use the "receive" command to decrypt
it.

	$ truss secret receive s.8ZruRfPtTRMHEn7bBoS7Gz2R

This will return the secret content to Stdout, which can then be copied with the
"pbcopy" command on MacOS:

	$ truss secret receive s.8ZruRfPtTRMHEn7bBoS7Gz2R | pbcopy
	
Once a shared secret has been received/unwrapped, it cannot be received again.`,
	Args: cobra.MaximumNArgs(1),
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

		fmt.Print(string(token))
		return nil
	},
}

func init() {
	secretsCmd.AddCommand(secretShareCmd)

	secretShareCmd.Flags().Int("ttl", 30, "TTL of the wrapped secret in minutes")
}
