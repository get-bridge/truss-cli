package cmd

import (
	"os"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestSecrets(t *testing.T) {
	Convey("findSecret", t, func() {
		secretName := "secret-name"
		kubeconfigName := "kube-config-name"

		viper.Reset()

		sm, err := truss.NewSecretsManager(os.Getenv("EDITOR"), nil)
		So(err, ShouldBeNil)

		Convey("runs no errors", func() {
			config, err := findSecret(sm, []string{secretName}, "pull")
			So(err, ShouldBeNil)
			So(config.Name, ShouldEqual, "secret-name")
		})

		Convey("provided with explicit kubeconfig name", func() {
			args := []string{secretName, kubeconfigName}

			Convey("runs no errors", func() {
				config, err := findSecret(sm, args, "pull")
				So(err, ShouldBeNil)
				So(config.Name, ShouldEqual, "secret-name")
			})

			Convey("errors if unknown env specified", func() {
				viper.Set("TRUSS_ENV", "test-env")
				_, err := findSecret(sm, args, "pull")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown env test-env. Options: []")
			})

			Convey("errors if env also specified", func() {
				viper.Set("TRUSS_ENV", "test-env")
				viper.Set("environments", map[string]string{"test-env": "kube-config-name"})
				_, err := findSecret(sm, args, "pull")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "do not specify --env and kubeconfig")
			})

			Convey("errors if no secret defined", func() {
				_, err := findSecret(sm, []string{"wrong-secret-name", "kube-config-name"}, "pull")
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "secret named 'wrong-secret-name' in 'kube-config-name' not found")
			})
		})
	})
}
