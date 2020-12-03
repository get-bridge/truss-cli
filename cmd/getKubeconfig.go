package cmd

import (
	"errors"

	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getKubeconfigCmd represents the getKubeconfig command
var getKubeconfigCmd = &cobra.Command{
	Use:   "get-kubeconfig",
	Short: "Get Kubeconfigs from source",
	Long: `
kubeconfigfiles:
  s3:
    bucket: my-aws-bucket
    region: us-east-1
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, err := getKubeDir()
		if err != nil {
			return err
		}

		s3bucket := viper.GetString("kubeconfigfiles.s3.bucket")
		if s3bucket == "" {
			if viper.ConfigFileUsed() == "" {
				return errors.New("No global config file found")
			}
			return errors.New(viper.ConfigFileUsed() + " does not contain a valid 'kubeconfigfiles'")
		}
		awsrole := viper.GetString("kubeconfigfiles.s3.awsrole")
		region := viper.GetString("kubeconfigfiles.s3.region")
		if region == "" {
			return errors.New("s3 config must have region")
		}
		return truss.GetKubeconfigS3(awsrole, s3bucket, dest, region).Fetch()
	},
}

func init() {
	rootCmd.AddCommand(getKubeconfigCmd)
}
