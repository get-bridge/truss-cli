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

		secretsContent := fmt.Sprintf(`
transit-key-name: foo-transit

secrets:
- name: secret-name
  filePath: %s/secrets.file
  vaultPath: secret/bridge/truss-cli-test/file
  kubeconfig: kubeconfig-truss-nonprod-iad
`, tmp)

		secretsFileContent := `secrets:
  foo:
    a: b
`

		secretsPath := tmp + "/secrets.yaml"
		err := ioutil.WriteFile(secretsPath, []byte(secretsContent), 0644)
		So(err, ShouldBeNil)

		err = os.Setenv("TRUSS_SECRETS_FILE", secretsPath)
		So(err, ShouldBeNil)

		var auth VaultAuth
		awsrole, ok := os.LookupEnv("TEST_AWS_ROLE")
		if ok {
			vaultrole := os.Getenv("TEST_VAULT_ROLE")
			auth = VaultAuthAWS(vaultrole, awsrole)
		}
		sm, err := NewSecretsManager("", auth)
		So(err, ShouldBeNil)

		firstSecret := sm.SecretConfigList.Secrets[0]
		So(firstSecret, ShouldNotBeNil)

		vault, err := sm.vault(firstSecret)
		So(err, ShouldBeNil)
		err = encryptAndSaveToDisk(vault, sm.TransitKeyName, path.Join(tmp, "secrets.file"), []byte(secretsFileContent))
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
				secondSecret := &SecretFileConfig{}
				err := sm.Pull(secondSecret)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Push and Pull", func() {
			err := sm.Push(firstSecret)
			So(err, ShouldBeNil)

			vaultData, err := vault.Run([]string{"kv", "get", "-field=a", "secret/bridge/truss-cli-test/file/foo"})
			So(err, ShouldBeNil)
			So(string(vaultData), ShouldEqual, "b")

			err = sm.Pull(firstSecret)
			So(err, ShouldBeNil)

			bytes, err := firstSecret.getDecryptedFromDisk(vault, sm.TransitKeyName)
			So(err, ShouldBeNil)
			So(string(bytes), ShouldEqual, secretsFileContent)
		})

		Convey("View", func() {
			secretString, err := sm.View(firstSecret)
			So(err, ShouldBeNil)
			So(secretString, ShouldEqual, secretsFileContent)
		})

		Convey("EncryptSecret", func() {
			newFile, err := ioutil.TempFile("", "")
			So(err, ShouldBeNil)
			defer os.Remove(newFile.Name())

			newFile.WriteString(secretsFileContent)
			newFile.Close()

			err = sm.EncryptSecret(SecretFileConfig{filePath: newFile.Name(), kubeconfig: firstSecret.Kubeconfig()})
			So(err, ShouldBeNil)

			bytes, err := ioutil.ReadFile(newFile.Name())
			So(err, ShouldBeNil)
			So(string(bytes), ShouldNotEqual, secretsFileContent)
		})
	})
}
