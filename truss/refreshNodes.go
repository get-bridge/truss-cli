package truss

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

// RefreshNodesCmd is used to refresh Truss nodes
type RefreshNodesCmd struct {
	ASC autoscalingiface.AutoScalingAPI
}

// NewRefreshNodesCmd instantiates a RefreshNodesCmd
func NewRefreshNodesCmd(region, arn string) *RefreshNodesCmd {
	sess := NewAWSSession(region, arn)
	asc := autoscaling.New(sess)

	return &RefreshNodesCmd{
		ASC: asc,
	}
}

// ASGFilterFunc filters ASGs
type ASGFilterFunc func(*autoscaling.Group) bool

// GetFilteredAutoscalingGroups returns a filtered list of ASGs
func (c RefreshNodesCmd) GetFilteredAutoscalingGroups(ff ASGFilterFunc) ([]*autoscaling.Group, error) {
	asgs := []*autoscaling.Group{}

	if err := c.ASC.DescribeAutoScalingGroupsPages(&autoscaling.DescribeAutoScalingGroupsInput{}, func(r *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
		for _, asg := range r.AutoScalingGroups {
			if ff(asg) {
				asgs = append(asgs, asg)
			}
		}
		return !lastPage
	}); err != nil {
		return nil, err
	}

	return asgs, nil
}

// RefreshNodes triggers an Instance Refresh on the provided ASG
func (c RefreshNodesCmd) RefreshNodes(g *autoscaling.Group) error {
	if g == nil {
		return errors.New("nope")
	}

	_, err := c.ASC.StartInstanceRefresh(&autoscaling.StartInstanceRefreshInput{
		AutoScalingGroupName: g.AutoScalingGroupName,
	})

	return err
}
