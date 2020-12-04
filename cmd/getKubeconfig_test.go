package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestGetKubeconfig(t *testing.T) {
	Convey("get-kubeconfig", t, func() {
		c := &cobra.Command{}
		viper.Reset()

		Convey("runs no errors", func() {
			err := getKubeconfigCmd.RunE(c, []string{})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "No global config file found")
		})

		Convey("s3 configured", func() {
			Convey("returns errors", func() {
				viper.Set("kubeconfigfiles.s3", map[string]interface{}{
					"bucket": "foo-bar",
				})
				err := getKubeconfigCmd.RunE(c, []string{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "s3 config must have region")
			})
		})
	})
}
