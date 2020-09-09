package truss

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretConfig(t *testing.T) {
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
		defaultConfig := SecretConfig{FilePath: f.Name()}
		err = defaultConfig.encryptAndSaveToDisk(vault, []byte(fileContent))
		So(err, ShouldBeNil)

		Convey("existsOnDisk", func() {
			Convey("true if file exists", func() {
				ok := defaultConfig.existsOnDisk()
				So(ok, ShouldBeTrue)
			})

			Convey("false if empty", func() {
				ok := SecretConfig{}.existsOnDisk()
				So(ok, ShouldBeFalse)
			})

			Convey("false", func() {
				ok := SecretConfig{FilePath: "probably-not-a-file"}.existsOnDisk()
				So(ok, ShouldBeFalse)
			})
		})

		Convey("getDecryptedFromDisk", func() {
			Convey("returns contents", func() {
				bytes, err := SecretConfig{}.getDecryptedFromDisk(vault)
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, "secrets: {}")
			})

			Convey("returns default if file doesn't exist", func() {
				bytes, err := defaultConfig.getDecryptedFromDisk(vault)
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, fileContent)
			})
		})

		Convey("getMapFromDisk", func() {
			Convey("returns map of secrets", func() {
				m, err := defaultConfig.getMapFromDisk(vault)
				So(err, ShouldBeNil)
				So(m, ShouldResemble, contentsMap)
			})
		})

		Convey("writeMapToDisk", func() {
			Convey("writes encrypted secrets", func() {
				f, err := ioutil.TempFile("", "")
				So(err, ShouldBeNil)
				defer func() {
					f.Close()
					os.Remove(f.Name())
				}()

				err = SecretConfig{FilePath: f.Name()}.writeMapToDisk(vault, contentsMap)
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(f.Name())
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, fileContent+"-encrypted")
			})
		})

		Convey("write", func() {
			Convey("writes to vault", func() {
				err = defaultConfig.write(vault)
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
	})
}
