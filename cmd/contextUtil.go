package cmd

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

func getKubeContext(cmd *cobra.Command, args []string) (string, error) {
	var context string

	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return context, err
	}
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		return context, err
	}

	if env != "" {
		contexts, err := getContexts()
		if err != nil {
			return context, err
		}
		context = (*contexts)[fmt.Sprintf("%s-%s", env, region)]
	}
	return context, nil
}
