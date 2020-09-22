package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestSecrets(t *testing.T) {
	secretName := "secret-name"
	kubeconfigName := "kube-config-name"
	secretsFileContent := fmt.Sprintf(`
transit-key-name: omg-bbq

secrets:
- name: %v
  kubeconfig: %v
`, secretName, kubeconfigName)

	Convey("newSecretsManager", t, func() {
		viper.Reset()

		Convey("accepts TRUSS_SECRETS_FILE", func() {
			viper.Set("TRUSS_SECRETS_FILE", "foo")

			_, err := newSecretsManager()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "open foo")
		})

		Convey("find secret in current directory", func() {
			dir, err := ioutil.TempDir("", "")
			So(err, ShouldBeNil)
			defer os.RemoveAll(dir)

			err = ioutil.WriteFile(path.Join(dir, defaultSecretsFileName), []byte(secretsFileContent), 0644)
			So(err, ShouldBeNil)

			err = os.Chdir(dir)
			So(err, ShouldBeNil)

			_, err = newSecretsManager()
			So(err, ShouldBeNil)
		})

		Convey("find secret in parent directory", func() {
			dir, err := ioutil.TempDir("", "")
			So(err, ShouldBeNil)
			defer os.RemoveAll(dir)

			err = ioutil.WriteFile(path.Join(dir, defaultSecretsFileName), []byte(secretsFileContent), 0644)
			So(err, ShouldBeNil)

			subDir := path.Join(dir, "foo/bar/baz")
			err = os.MkdirAll(subDir, 0777)
			So(err, ShouldBeNil)
			err = os.Chdir(subDir)
			So(err, ShouldBeNil)

			_, err = newSecretsManager()
			So(err, ShouldBeNil)
		})
	})

	Convey("findSecret", t, func() {
		viper.Reset()

		// save to tmp secret
		tmpFile, err := ioutil.TempFile("", "")
		So(err, ShouldBeNil)
		defer os.Remove(tmpFile.Name())
		tmpFile.WriteString(secretsFileContent)
		tmpFile.Close()

		sm, err := truss.NewSecretsManager(tmpFile.Name(), viper.GetString("EDITOR"), getVaultAuth())
		So(err, ShouldBeNil)

		Convey("runs no errors", func() {
			config, err := findSecret(sm, []string{secretName}, "pull")
			So(err, ShouldBeNil)
			So(config.Name(), ShouldEqual, "secret-name")
		})

		Convey("provided with explicit kubeconfig name", func() {
			args := []string{secretName, kubeconfigName}

			Convey("runs no errors", func() {
				config, err := findSecret(sm, args, "pull")
				So(err, ShouldBeNil)
				So(config.Name(), ShouldEqual, "secret-name")
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
