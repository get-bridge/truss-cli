package truss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretsManager(t *testing.T) {
	Convey("SecretsManager", t, func() {
		tmp := os.TempDir()
		defer os.RemoveAll(tmp)

		secretsFileName := path.Join(tmp, "secrets.file")

		secretsContent := fmt.Sprintf(`
transit-key-name: foo-transit

secrets:
- name: secret-name
  filePath: %s
  vaultPath: secret/bridge/truss-cli-test/file
  kubeconfig: kubeconfig-truss-nonprod-iad
`, secretsFileName)

		secretsPath := tmp + "/secrets.yaml"
		err := ioutil.WriteFile(secretsPath, []byte(secretsContent), 0644)
		So(err, ShouldBeNil)

		sm, err := NewSecretsManager(secretsPath, "", nil)
		So(err, ShouldBeNil)

		Convey("parses secrets", func() {
			firstSecret := sm.SecretConfigList.Secrets[0]
			So(firstSecret, ShouldNotBeNil)
		})
	})
}
