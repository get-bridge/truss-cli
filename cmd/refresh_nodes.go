package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/instructure-bridge/truss-cli/truss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type asgFilterFunc func(*autoscaling.Group) bool

var refreshNodesCmd = &cobra.Command{
	Use:   "refresh-nodes [-a|--all] [-n|--nodegroup]",
	Short: "Trigger an instance refresh on a Truss Nodegroup's ASG",
	RunE: func(cmd *cobra.Command, args []string) error {
		kc, err := getKubeconfigStruct()
		if err != nil {
			return err
		}

		var clusterName string
		var clusterRegion string
		var clusterRoleArn string
		var auth *clientcmdapi.AuthInfo
		for _, a := range kc.AuthInfos {
			auth = a
			break
		}
		for k, v := range auth.Exec.Args {
			if v == "--cluster-name" {
				clusterName = auth.Exec.Args[k+1]
			}
			if v == "--region" {
				clusterRegion = auth.Exec.Args[k+1]
			}
			if v == "--role" {
				clusterRoleArn = auth.Exec.Args[k+1]
			}
		}

		sess := truss.NewAWSSession(clusterRegion, clusterRoleArn)
		asc := autoscaling.New(sess)

		var ff asgFilterFunc
		if viper.GetBool("refresh_all") {
			ff = func(g *autoscaling.Group) bool {
				for _, t := range g.Tags {
					if *t.Key == "kubernetes.io/cluster/"+clusterName && *t.Value == "owned" {
						return true
					}
				}
				return false
			}
		} else {
			ff = func(g *autoscaling.Group) bool {
				prefix := strings.Replace(clusterName, "cluster", "", 1) + viper.GetString("refresh_node_group")
				return strings.HasPrefix(*g.AutoScalingGroupName, prefix)
			}
		}

		if err := asc.DescribeAutoScalingGroupsPages(&autoscaling.DescribeAutoScalingGroupsInput{}, func(r *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			for _, g := range r.AutoScalingGroups {
				if ff(g) {
					log.Printf("Triggering instance refresh on ASG %s", *g.AutoScalingGroupName)
					// refresh the asg
				}
			}
			return !lastPage
		}); err != nil {
			// some aws error
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(refreshNodesCmd)

	refreshNodesCmd.Flags().BoolP("all", "a", false, "Refresh all node groups")
	viper.BindPFlag("refresh_all", refreshNodesCmd.Flags().Lookup("all"))
	refreshNodesCmd.Flags().StringP("nodegroup", "n", "default", "Node group to refresh")
	viper.BindPFlag("refresh_node_group", refreshNodesCmd.Flags().Lookup("nodegroup"))
}
