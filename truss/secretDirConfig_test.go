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
	Convey("TestSecretConfig", t, func() {
		dir, err := ioutil.TempDir("", "")
		defer os.Remove(dir)
		So(err, ShouldBeNil)
		So(dir, ShouldNotBeEmpty)

		aContent := "a is best"
		contentsMap := map[string]string{
			"a": aContent,
		}
		secretYaml := fmt.Sprintf(`a: %v
`, aContent)

		err = ioutil.WriteFile(path.Join(dir, "a"), []byte(aContent), 0644)
		So(err, ShouldBeNil)

		_ = aContent

		vault := &mockVault{}
		defaultConfig := SecretDirConfig{dirPath: dir, vaultPath: "secret/path"}

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
				bytes, err := SecretDirConfig{}.getDecryptedFromDisk(vault, "")
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, `{}
`)
			})

			Convey("returns default if file doesn't exist", func() {
				bytes, err := defaultConfig.getDecryptedFromDisk(vault, "")
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, secretYaml)
			})
		})

		Convey("getMapFromDisk", func() {
			Convey("returns map of secrets", func() {
				m, err := defaultConfig.getMapFromDisk(vault, "")
				So(err, ShouldBeNil)
				So(m, ShouldResemble, contentsMap)
			})
		})

		Convey("saveToDiskFromVault", func() {
			Convey("writes encrypted secrets", func() {
				newDir, err := ioutil.TempDir("", "")
				defer os.Remove(dir)
				So(err, ShouldBeNil)
				So(dir, ShouldNotBeEmpty)

				secrets := map[string]interface{}{
					"a": []byte(aContent),
				}
				config := SecretDirConfig{dirPath: newDir}
				err = config.saveToDiskFromVault(&mockVault{secrets: secrets}, "")
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(path.Join(newDir, "a"))
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, aContent+"-encrypted")
			})
		})

		Convey("writeToVault", func() {
			Convey("writes to vault", func() {
				err = defaultConfig.writeToVault(vault, "")
				So(err, ShouldBeNil)
				So(vault.commands, ShouldHaveLength, 1)
				So(vault.commands, ShouldContain,
					[]string{"kv", "put", "secret/path", "a=a is best"},
				)
			})
		})

		Convey("saveToDisk", func() {
			Convey("writes to disk", func() {
				newDir, err := ioutil.TempDir("", "")
				defer os.Remove(dir)
				So(err, ShouldBeNil)
				So(dir, ShouldNotBeEmpty)

				err = SecretDirConfig{dirPath: newDir}.saveToDisk(vault, "", []byte(secretYaml))
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(path.Join(newDir, "a"))
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, aContent+"-encrypted")
			})
		})
	})
}
