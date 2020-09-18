package truss

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecrets(t *testing.T) {
	Convey("SecretConfigList", t, func() {
		secretsContent := `
transit-key-name: foo

secrets:
- name: secret-name
  kubeconfig: kube-config-name
- name: secret-name
  kubeconfig: kube-config-name-2
- name: secret-name-2
  kubeconfig: kube-config-name-2
`
		tmp := os.TempDir()
		defer os.RemoveAll(tmp)

		secretsPath := tmp + "/secrets.yaml"
		err := ioutil.WriteFile(secretsPath, []byte(secretsContent), 0644)
		So(err, ShouldBeNil)

		list, err := SecretConfigListFromFile(secretsPath)
		So(err, ShouldBeNil)

		Convey("Secret", func() {
			Convey("returns secret", func() {
				config, err := list.Secret("secret-name", "kube-config-name")
				So(err, ShouldBeNil)
				So(config.Name(), ShouldEqual, "secret-name")
				So(config.Kubeconfig(), ShouldEqual, "kube-config-name")
			})
		})

		Convey("SecretNames", func() {
			Convey("returns all secret names", func() {
				secrets := list.SecretNames()
				So(secrets, ShouldContain, "secret-name")
				So(secrets, ShouldContain, "secret-name-2")
			})
		})

		Convey("SecretKubeconfigs", func() {
			Convey("returns all kubeconfigs for secret", func() {
				kubeconfigs := list.SecretKubeconfigs("secret-name")
				So(kubeconfigs, ShouldContain, "kube-config-name")
				So(kubeconfigs, ShouldContain, "kube-config-name-2")
			})
		})

		// Some reason this breaks allthethings
		// Reset(func() {
		// 	os.RemoveAll(tmp)
		// })
	})
}
