package truss

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	tFoo = "foo"
	tBar = "bar"
	tBaz = "baz"
)

func TestRefreshNodes(t *testing.T) {
	Convey("RefreshNodes", t, func() {

		Convey("It instantiates the cmd", func() {
			rnc := NewRefreshNodesCmd("us-east-2", "some-arn")
			So(rnc, ShouldNotBeNil)
			So(rnc.ASC, ShouldNotBeNil)
		})

		state := &asapiState{}
		rnc := &RefreshNodesCmd{
			ASC: asapi{
				asapiState: state,
				ASGs: []*autoscaling.Group{
					{AutoScalingGroupName: &tFoo},
					{AutoScalingGroupName: &tBar},
					{AutoScalingGroupName: &tBaz},
				},
			},
		}

		Convey("It gets filtered groups", func() {
			g, err := rnc.GetFilteredAutoscalingGroups(func(asg *autoscaling.Group) bool {
				return strings.HasPrefix(*asg.AutoScalingGroupName, "b")
			})

			So(err, ShouldBeNil)
			So(g, ShouldHaveLength, 2)
			So(*g[0].AutoScalingGroupName, ShouldEqual, "bar")

			g, err = rnc.GetFilteredAutoscalingGroups(func(*autoscaling.Group) bool { return true })

			So(err, ShouldBeNil)
			So(g, ShouldHaveLength, 3)
			So(*g[0].AutoScalingGroupName, ShouldEqual, "foo")
		})

		Convey("It doesn't allow refreshing of nothing", func() {
			err := rnc.RefreshNodes(nil)
			So(err, ShouldBeError)
			So(state.InstanceRefreshed, ShouldBeFalse)
		})

		Convey("It does allow refreshing of nodegroups", func() {
			err := rnc.RefreshNodes(&autoscaling.Group{
				AutoScalingGroupName: &tFoo,
			})
			So(err, ShouldBeNil)
			So(state.InstanceRefreshed, ShouldBeTrue)
			So(*state.StartInstanceRefresh.AutoScalingGroupName, ShouldEqual, "foo")
		})
	})
}

type asapi struct {
	autoscalingiface.AutoScalingAPI
	*asapiState
	ASGs []*autoscaling.Group
}

type asapiState struct {
	InstanceRefreshed    bool
	StartInstanceRefresh *autoscaling.StartInstanceRefreshInput
}

func (a asapi) StartInstanceRefresh(in *autoscaling.StartInstanceRefreshInput) (*autoscaling.StartInstanceRefreshOutput, error) {
	(*a.asapiState).InstanceRefreshed = true
	(*a.asapiState).StartInstanceRefresh = in
	return nil, nil
}

func (a asapi) DescribeAutoScalingGroupsPages(in *autoscaling.DescribeAutoScalingGroupsInput, cb func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool) error {
	out := &autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: a.ASGs,
	}
	cb(out, true)

	return nil
}

var _ autoscalingiface.AutoScalingAPI = asapi{}
