package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	Convey("getKubeconfigName", t, func() {
		Convey("gets kubeconfig name from env", func() {
			viper.Set("environments", map[string]string{
				"my-env": "my-kube",
			})

			rootCmd.SetArgs([]string{})
			rootCmd.PersistentFlags().Set("env", "my-env")
			kubeconfig, err := getKubeconfigName()
			So(err, ShouldBeNil)
			So(kubeconfig, ShouldEqual, "my-kube")
		})
	})

	Convey("getKubeconfig", t, func() {
		Convey("defaults to empty string", func() {
			viper.Set("environments", map[string]string{})

			rootCmd.SetArgs([]string{})
			rootCmd.PersistentFlags().Set("env", "")
			kubeconfig, err := getKubeconfig()
			So(err, ShouldBeNil)
			So(kubeconfig, ShouldBeEmpty)
		})

		Convey("gets kubeconfig full path from env", func() {
			viper.Set("environments", map[string]string{
				"my-env": "my-kube",
			})

			rootCmd.SetArgs([]string{})
			rootCmd.PersistentFlags().Set("env", "my-env")
			kubeconfig, err := getKubeconfig()
			So(err, ShouldBeNil)
			So(kubeconfig, ShouldEndWith, "my-kube")
		})
	})

	Convey("getKubeDir", t, func() {
		Convey("gets kubeconfig name from env", func() {
			viper.Set("kubeconfigfiles", map[string]interface{}{
				"directory": "my-kube-dir",
			})

			dir, err := getKubeDir()
			So(err, ShouldBeNil)
			So(dir, ShouldEndWith, "my-kube-dir")
		})
	})
}
