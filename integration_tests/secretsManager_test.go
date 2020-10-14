package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/instructure-bridge/truss-cli/truss"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretsManager(t *testing.T) {
	Convey("SecretsManager", t, func() {
		tmp := os.TempDir()
		defer os.RemoveAll(tmp)

		secretsFileName := path.Join(tmp, "secrets.file")
		secretsFileContent := `secrets:
  foo:
    a: b
`

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

		var auth truss.VaultAuth
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if ok {
			vaultrole := os.Getenv("TEST_VAULT_ROLE")
			auth = truss.VaultAuthAWS(vaultrole, awsrole, "us-east-1")
		}
		sm, err := truss.NewSecretsManager(secretsPath, "", auth)
		So(err, ShouldBeNil)

		firstSecret := sm.SecretConfigList.Secrets[0]
		So(firstSecret, ShouldNotBeNil)

		err = ioutil.WriteFile(secretsFileName, []byte(secretsFileContent), 0644)
		So(err, ShouldBeNil)
		sm.EncryptSecret(firstSecret)
		So(err, ShouldBeNil)

		// TODO how do we deal with $EDITOR
		Convey("Edit", nil)

		Convey("PushAll", func() {
			err := sm.PushAll()
			So(err, ShouldBeNil)
		})

		Convey("PullAll", func() {
			err := sm.PullAll()
			So(err, ShouldBeNil)
		})

		Convey("Pull", func() {
			Convey("errors if secret invalid", func() {
				secondSecret := &truss.SecretFileConfig{}
				err := sm.Pull(secondSecret)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Push and Pull", func() {
			err := sm.Push(firstSecret)
			So(err, ShouldBeNil)

			vault, err := sm.Vault(firstSecret)
			So(err, ShouldBeNil)
			vaultData, err := vault.GetMap("secret/data/bridge/truss-cli-test/file/foo")
			So(err, ShouldBeNil)
			So(vaultData, ShouldResemble, map[string]interface{}{
				"a": "b",
			})

			err = sm.Pull(firstSecret)
			So(err, ShouldBeNil)

			localContent, remoteContent, err := sm.View(firstSecret)
			So(err, ShouldBeNil)
			So(localContent, ShouldEqual, secretsFileContent)
			So(remoteContent, ShouldEqual, secretsFileContent)
		})
	})
}
