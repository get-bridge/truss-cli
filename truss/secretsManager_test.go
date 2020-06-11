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

		// TODO
		Convey("Edit", nil)
		// TODO
		Convey("PushAll", nil)
		// TODO
		Convey("Push", nil)
		// TODO
		Convey("PullAll", nil)
		// TODO
		Convey("Pull", nil)
		// TODO
		Convey("GetDecryptedFromDisk", nil)
		// TODO
		Convey("GetMapFromDisk", nil)
		// TODO
		Convey("GetMapFromVault", nil)
		// TODO
		Convey("WriteMapToDisk", nil)
		// TODO
		Convey("EncryptAndSaveToDisk", nil)
		// TODO
		Convey("Decrypt", nil)
		// TODO
		Convey("Encrypt", nil)
		// TODO
		Convey("Write", nil)
	})
}
