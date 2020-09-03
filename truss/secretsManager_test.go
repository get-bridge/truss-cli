package truss

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretsManager(t *testing.T) {
	Convey("SecretsManager", t, func() {
		secretsContent := `
transit-key-name: foo

secrets:
- name: secret-name
  kubeconfig: kube-config-name
`

		tmp := os.TempDir()
		secretsPath := tmp + "/secrets.yaml"
		err := ioutil.WriteFile(secretsPath, []byte(secretsContent), 0644)
		So(err, ShouldBeNil)

		err = os.Setenv("TRUSS_SECRETS_FILE", secretsPath)
		So(err, ShouldBeNil)
		sm, err := NewSecretsManager("", nil)
		So(err, ShouldBeNil)

		// TODO remove when used
		_ = sm

		Convey("Edit", nil)
		Convey("PushAll", nil)
		Convey("Push", nil)
		Convey("PullAll", nil)
		Convey("Pull", nil)
	})
}
