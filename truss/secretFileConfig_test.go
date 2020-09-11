package truss

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretFileConfig(t *testing.T) {
	Convey("TestSecretConfig", t, func() {
		fileContent := `secrets:
  firstApp:
    firstSecret: "1"
    secondSecret: "2"
  secondApp: {}
`
		contentsMap := map[string]map[string]string{
			"firstApp": {
				"firstSecret":  "1",
				"secondSecret": "2",
			},
			"secondApp": {},
		}
		f, err := ioutil.TempFile("", "")
		So(err, ShouldBeNil)
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		vault := &mockVault{}
		defaultConfig := SecretFileConfig{filePath: f.Name()}
		err = encryptAndSaveToDisk(vault, "", f.Name(), []byte(fileContent))
		So(err, ShouldBeNil)

		Convey("existsOnDisk", func() {
			Convey("true if file exists", func() {
				ok := defaultConfig.existsOnDisk()
				So(ok, ShouldBeTrue)
			})

			Convey("false if empty", func() {
				ok := SecretFileConfig{}.existsOnDisk()
				So(ok, ShouldBeFalse)
			})

			Convey("false", func() {
				ok := SecretFileConfig{filePath: "probably-not-a-file"}.existsOnDisk()
				So(ok, ShouldBeFalse)
			})
		})

		Convey("getDecryptedFromDisk", func() {
			Convey("returns contents", func() {
				bytes, err := SecretFileConfig{}.getDecryptedFromDisk(vault, "")
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, "secrets: {}")
			})

			Convey("returns default if file doesn't exist", func() {
				bytes, err := defaultConfig.getDecryptedFromDisk(vault, "")
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, fileContent)
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
				newFile, err := ioutil.TempFile("", "")
				So(err, ShouldBeNil)
				defer func() {
					f.Close()
					os.Remove(f.Name())
				}()

				secrets := map[string]interface{}{
					"firstApp": map[string]string{
						"firstSecret":  "1",
						"secondSecret": "2",
					},
					"secondApp": map[string]string{},
				}
				config := SecretFileConfig{filePath: newFile.Name()}
				err = config.saveToDiskFromVault(&mockVault{secrets: secrets}, "")
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(newFile.Name())
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, fileContent+"-encrypted")
			})
		})

		Convey("writeToVault", func() {
			Convey("writes to vault", func() {
				err = defaultConfig.writeToVault(vault, "")
				So(err, ShouldBeNil)
				So(vault.commands, ShouldHaveLength, 2)
				So(vault.commands, ShouldContain,
					[]string{"kv", "put", "firstApp", "firstSecret=1", "secondSecret=2"},
				)
				So(vault.commands, ShouldContain,
					[]string{"kv", "put", "secondApp"},
				)
			})
		})

		Convey("saveToDisk", func() {
			Convey("writes to disk", func() {
				newFile, err := ioutil.TempFile("", "")
				So(err, ShouldBeNil)
				defer func() {
					f.Close()
					os.Remove(f.Name())
				}()

				err = SecretFileConfig{filePath: newFile.Name()}.saveToDisk(vault, "", []byte(fileContent))
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(newFile.Name())
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, fileContent+"-encrypted")
			})
		})
	})
}
