package truss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretDirConfig(t *testing.T) {
	vault, _ := createTestVault(t)
	//defer server.Stop()

	Convey("TestSecretConfig", t, func() {
		dir, err := ioutil.TempDir("", "")
		defer os.RemoveAll(dir)
		So(err, ShouldBeNil)
		So(dir, ShouldNotBeEmpty)

		transitKey := "dir-config-test"
		aContent := "a is best"
		contentsMap := map[string]string{
			"a": aContent,
		}
		secretYaml := fmt.Sprintf(`a: %v
`, aContent)

		err = ioutil.WriteFile(path.Join(dir, "a"), []byte(aContent), 0644)
		So(err, ShouldBeNil)

		defaultConfig := SecretDirConfig{dirPath: dir, vaultPath: "kv/dir/path"}

		Convey("existsOnDisk", func() {
			Convey("true if file exists", func() {
				ok := defaultConfig.existsOnDisk()
				So(ok, ShouldBeTrue)
			})

			Convey("false if empty", func() {
				ok := SecretDirConfig{}.existsOnDisk()
				So(ok, ShouldBeFalse)
			})

			Convey("false", func() {
				ok := SecretDirConfig{dirPath: "probably-not-a-dir"}.existsOnDisk()
				So(ok, ShouldBeFalse)
			})
		})

		Convey("getDecryptedFromDisk", func() {
			Convey("returns contents", func() {
				bytes, err := SecretDirConfig{}.getDecryptedFromDisk(vault, transitKey)
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, `{}
`)
			})

			Convey("returns default if file doesn't exist", func() {
				bytes, err := defaultConfig.getDecryptedFromDisk(vault, transitKey)
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, secretYaml)
			})
		})

		Convey("getMapFromDisk", func() {
			Convey("returns map of secrets", func() {
				m, err := defaultConfig.getMapFromDisk(vault, transitKey)
				So(err, ShouldBeNil)
				So(m, ShouldResemble, contentsMap)
			})
		})

		Convey("saveToDiskFromVault", func() {
			Convey("writes encrypted secrets", func() {
				newDir, err := ioutil.TempDir("", "")
				defer os.RemoveAll(newDir)
				So(err, ShouldBeNil)
				So(dir, ShouldNotBeEmpty)

				_, err = vault.Write(kv2DataPath(defaultConfig.vaultPath), map[string]interface{}{
					"data": contentsMap,
				})
				So(err, ShouldBeNil)
				config := SecretDirConfig{dirPath: newDir, vaultPath: defaultConfig.vaultPath}
				err = config.saveToDiskFromVault(vault, transitKey)
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(path.Join(newDir, "a"))
				So(err, ShouldBeNil)

				decrypted, err := vault.Decrypt(transitKey, bytes)
				So(err, ShouldBeNil)
				So(string(decrypted), ShouldEqual, aContent)
			})
		})

		Convey("writeToVault", func() {
			Convey("writes to vault", func() {
				err = defaultConfig.writeToVault(vault, transitKey)
				So(err, ShouldBeNil)

				data, err := vault.GetMap(kv2DataPath(defaultConfig.vaultPath))
				So(err, ShouldBeNil)
				So(data, ShouldResemble, map[string]interface{}{
					"a": contentsMap["a"],
				})
			})
		})

		Convey("saveToDisk", func() {
			Convey("writes to disk", func() {
				newDir, err := ioutil.TempDir("", "")
				defer os.RemoveAll(newDir)
				So(err, ShouldBeNil)
				So(dir, ShouldNotBeEmpty)

				err = SecretDirConfig{dirPath: newDir}.saveToDisk(vault, "test-trans", []byte(secretYaml))
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(path.Join(newDir, "a"))
				So(err, ShouldBeNil)
				So(string(bytes), ShouldNotResemble, secretYaml)
			})
		})

		Convey("getFromVault", func() {
			Convey("returns from vault", func() {
				bytes, err := defaultConfig.getFromVault(vault)
				So(err, ShouldBeNil)
				So(string(bytes), ShouldResemble, secretYaml)
			})
		})
	})
}
