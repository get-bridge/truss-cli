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

	Convey("getKubeconfigStruct", t, func() {
		viper.Set("environments", map[string]string{
			"test-env": "kubeconfig-truss-dev-cmh",
		})
		viper.Set("kubeconfigfiles.directory", "fixtures")
		viper.Set("TRUSS_ENV", "test-env")

		kc, err := getKubeconfigStruct()
		So(err, ShouldBeNil)
		So(kc.AuthInfos, ShouldHaveLength, 1)
		auth := kc.AuthInfos["cluster-admin@truss-dev-cmh-shared-cluster"]
		So(auth.Exec.Args, ShouldHaveLength, 8)

		Convey("parses cluster name", func() {
			So(must(envClusterName()), ShouldEqual, "truss-dev-cmh-shared-cluster")
		})

		Convey("parses cluster region", func() {
			So(must(envClusterRegion()), ShouldEqual, "us-east-2")
		})

		Convey("parses cluster role arn", func() {
			So(must(envClusterRoleArn()), ShouldEqual, "arn:aws:iam::127178877223:role/xacct/ops-admin")
		})
	})
}
