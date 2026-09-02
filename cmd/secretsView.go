package cmd

import (
	"fmt"

	"github.com/get-bridge/truss-cli/truss"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var secretsViewCmd = &cobra.Command{
	Use:   "view [name] [kubeconfig]",
	Short: "Views a given environment's secrets on disk",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sm, err := newSecretsManager()
		if err != nil {
			return err
		}

		secret, err := findSecret(sm, args, "view")
		if err != nil {
			return err
		}

		_, err = secretCompare(sm, secret, true, false)
		return err
	},
}

// return true if same
func secretCompare(sm *truss.SecretsManager, secret truss.SecretConfig, localToRemote bool, quiet bool) (bool, error) {
	localContent, remoteContent, err := sm.View(secret)
	if err != nil {
		return false, err
	}

	dmp := diffmatchpatch.New()
	var diffs []diffmatchpatch.Diff
	if localToRemote {
		diffs = dmp.DiffMain(remoteContent, localContent, false)
	} else {
		diffs = dmp.DiffMain(localContent, remoteContent, false)
	}
	if !quiet {
		fmt.Println(dmp.DiffPrettyText(diffs))
	}
	return remoteContent == localContent, nil
}
