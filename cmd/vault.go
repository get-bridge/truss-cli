package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/instructure/truss-cli/truss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "A wrapper around hashicorp vault",
	Long: `This is useful when your vault is not exposed publicly.
As it will port-forward to the service and call aws auth`,
	Run: func(cmd *cobra.Command, args []string) {

		env, err := cmd.Flags().GetString("env")
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		region, err := cmd.Flags().GetString("region")
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}

		var context string
		if env != "" {
			contexts, err := getContexts()
			if err != nil {
				log.Errorln(err)
				os.Exit(1)
			}
			context = (*contexts)[fmt.Sprintf("%s-%s", env, region)]
		}

		err = truss.Vault(context).Run(args)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
}

func getContexts() (*map[string]string, error) {
	contextsPtr, ok := viper.Get("contexts").(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid context configuration")
	}
	contexts := map[string]string{}
	for k, v := range contextsPtr {
		contextStr, ok := v.(string)
		if !ok {
			log.Errorln("invalid dependency type", v)
			os.Exit(1)
		}
		contexts[k] = contextStr
	}

	return &contexts, nil
}
