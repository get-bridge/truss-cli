package truss

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretFileConfig(t *testing.T) {
	vault := createTestVault(t)

	Convey("TestSecretConfig", t, func() {
		transitKey := "file-test-transit"
		fileContent := `secrets:
  firstApp:
    firstSecret: "1"
    secondSecret: "2"
  secondApp:
    thirdSecret: "3"
`
		contentsMap := map[string]map[string]string{
			"firstApp": {
				"firstSecret":  "1",
				"secondSecret": "2",
			},
			"secondApp": {
				"thirdSecret": "3",
			},
		}
		f, err := ioutil.TempFile("", "")
		So(err, ShouldBeNil)
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		defaultConfig := SecretFileConfig{filePath: f.Name(), vaultPath: "kv/file/config"}
		err = encryptAndSaveToDisk(vault, transitKey, f.Name(), []byte(fileContent))
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
				bytes, err := defaultConfig.getDecryptedFromDisk(vault, transitKey)
				So(err, ShouldBeNil)
				So(string(bytes), ShouldEqual, fileContent)
			})

			Convey("errors if file doesn't exist", func() {
				_, err := SecretFileConfig{}.getDecryptedFromDisk(vault, transitKey)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "no such file")
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
				newFile, err := ioutil.TempFile("", "")
				So(err, ShouldBeNil)
				defer os.Remove(newFile.Name())

				_, err = vault.Write(kv2DataPath(defaultConfig.vaultPath)+"/firstApp", map[string]interface{}{
					"data": contentsMap["firstApp"],
				})
				So(err, ShouldBeNil)
				_, err = vault.Write(kv2DataPath(defaultConfig.vaultPath)+"/secondApp", map[string]interface{}{
					"data": contentsMap["secondApp"],
				})
				So(err, ShouldBeNil)

				config := SecretFileConfig{filePath: newFile.Name(), vaultPath: defaultConfig.vaultPath}
				err = config.saveToDiskFromVault(vault, transitKey)
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(newFile.Name())
				So(err, ShouldBeNil)

				decrypted, err := vault.Decrypt(transitKey, bytes)
				So(err, ShouldBeNil)
				So(string(decrypted), ShouldEqual, fileContent)
			})
		})

		Convey("writeToVault", func() {
			Convey("writes to vault", func() {
				err = defaultConfig.writeToVault(vault, transitKey)
				So(err, ShouldBeNil)

				data, err := vault.GetMap(kv2DataPath(defaultConfig.vaultPath) + "/firstApp")
				So(err, ShouldBeNil)
				So(data, ShouldResemble, map[string]interface{}{
					"firstSecret":  "1",
					"secondSecret": "2",
				})
				data, err = vault.GetMap(kv2DataPath(defaultConfig.vaultPath) + "/secondApp")
				So(err, ShouldBeNil)
				So(data, ShouldResemble, map[string]interface{}{
					"thirdSecret": "3",
				})
			})
		})

		Convey("saveToDisk", func() {
			Convey("creates dir if doesn't exist", func() {
				dir, err := ioutil.TempDir("", "")
				So(err, ShouldBeNil)
				defer os.RemoveAll(dir)

				fileName := dir + "/foo/bar/secret"
				err = SecretFileConfig{filePath: fileName}.saveToDisk(vault, transitKey, []byte(fileContent))
				So(err, ShouldBeNil)

				_, err = os.Stat(dir)
				So(err, ShouldBeNil)
			})

			Convey("errors if not valid yaml", func() {
				err = defaultConfig.saveToDisk(vault, transitKey, []byte(`{"asdf": {}}`))
				So(err, ShouldEqual, ErrSecretFileConfigInvalidYaml)
			})

			Convey("writes to disk", func() {
				newFile, err := ioutil.TempFile("", "")
				So(err, ShouldBeNil)
				defer os.Remove(newFile.Name())

				err = SecretFileConfig{filePath: newFile.Name()}.saveToDisk(vault, transitKey, []byte(fileContent))
				So(err, ShouldBeNil)

				bytes, err := ioutil.ReadFile(newFile.Name())
				So(err, ShouldBeNil)

				decrypted, err := vault.Decrypt(transitKey, bytes)
				So(err, ShouldBeNil)
				So(string(decrypted), ShouldEqual, fileContent)
			})
		})
	})
}
