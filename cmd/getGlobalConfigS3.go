package cmd

import (
	"github.com/instructure-bridge/truss-cli/truss"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// getGlobalConfigS3Cmd represents the getGlobalConfigS3 command
var getGlobalConfigS3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Get global config from S3",
	Long: `
Fetches .truss.yaml from S3 and puts it in your home directory:

	$ truss get-global-config s3 -b truss-cli-global-confi -k .truss.yaml

Uses S3 under the hood so it only works if you have AWS credentials set in your shell. Use the --role flag if you need to specify an AWS role ARN:

	$ truss get-global-config s3 --role arn:aws:iam::127178877223:role/xacct/ops-admin
		`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket, err := cmd.Flags().GetString("bucket")
		key, err := cmd.Flags().GetString("key")
		region, err := cmd.Flags().GetString("region")
		role, err := cmd.Flags().GetString("role")
		dir, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		if dir == "" {
			home, err := homedir.Dir()
			if err != nil {
				return err
			}
			dir = home
		}

		input := &truss.GetGlobalConfigS3Input{
			Bucket: bucket,
			Key:    key,
			Region: region,
			Role:   role,
			Dir:    dir,
		}

		err = truss.GetGlobalConfigS3(input)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	getGlobalConfigCmd.AddCommand(getGlobalConfigS3Cmd)
	getGlobalConfigS3Cmd.Flags().StringP("bucket", "b", "truss-cli-global-config", "S3 bucket that contains your .truss.yaml file")
	getGlobalConfigS3Cmd.Flags().StringP("key", "k", ".truss.yaml", "Name of the .truss.yaml file in the bucket")
	getGlobalConfigS3Cmd.Flags().StringP("region", "r", "us-east-2", "Region of S3 bucket that contains your .truss.yaml file")
	getGlobalConfigS3Cmd.Flags().StringP("role", "u", "", "Role with access to the S3 bucket with your .truss.yaml file")
	getGlobalConfigS3Cmd.Flags().StringP("out", "o", "", "Output directory where the .truss.yaml file is written")
}
