package cmd

import (
	"log"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/instructure-bridge/truss-cli/truss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type asgFilterFunc func(*autoscaling.Group) bool

var refreshNodesCmd = &cobra.Command{
	Use:   "refresh-nodes [-a|--all] [-n|--nodegroup]",
	Short: "Trigger an instance refresh on a Truss Nodegroup's ASG",
	RunE: func(cmd *cobra.Command, args []string) error {
		var ff truss.ASGFilterFunc
		if viper.GetBool("refresh_all") {
			ff = clusterFilterFunc(must(envClusterName()))
		} else {
			ff = nodeGroupFilterFunc(must(envClusterName()), viper.GetString("refresh_node_group"))
		}

		rnc := truss.NewRefreshNodesCmd(must(envClusterRegion()), must(envClusterRoleArn()))

		asgs, err := rnc.GetFilteredAutoscalingGroups(ff)
		if err != nil {
			return err
		}
		for _, asg := range asgs {
			if viper.GetBool("refresh_yes") || prompter.YN("Trigger Instance Refresh on "+*asg.AutoScalingGroupName+"?", false) {
				log.Printf("Triggering instance refresh on ASG %s", *asg.AutoScalingGroupName)
				rnc.RefreshNodes(asg)
			}
		}

		return nil
	},
}

func clusterFilterFunc(clusterName string) truss.ASGFilterFunc {
	return func(g *autoscaling.Group) bool {
		for _, t := range g.Tags {
			if *t.Key == "kubernetes.io/cluster/"+clusterName && *t.Value == "owned" {
				return true
			}
		}
		return false
	}
}

func nodeGroupFilterFunc(clusterName, groupName string) truss.ASGFilterFunc {
	return func(g *autoscaling.Group) bool {
		prefix := strings.Replace(clusterName, "cluster", "", 1) + groupName
		return strings.HasPrefix(*g.AutoScalingGroupName, prefix)
	}
}

func init() {
	rootCmd.AddCommand(refreshNodesCmd)

	refreshNodesCmd.Flags().BoolP("all", "a", false, "Refresh all node groups")
	viper.BindPFlag("refresh_all", refreshNodesCmd.Flags().Lookup("all"))
	refreshNodesCmd.Flags().StringP("nodegroup", "n", "default", "Node group to refresh")
	viper.BindPFlag("refresh_node_group", refreshNodesCmd.Flags().Lookup("nodegroup"))
	refreshNodesCmd.Flags().BoolP("yes", "y", false, "Say yes to prompts")
	viper.BindPFlag("refresh_yes", refreshNodesCmd.Flags().Lookup("yes"))
}
