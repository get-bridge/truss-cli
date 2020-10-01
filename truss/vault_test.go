package truss

import (
	"encoding/base64"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	vault := createTestVault(t)

	Convey("Vault", t, func() {
		binaryContent := []byte{0x0, 0xe8, 0x03, 0xd0, 0x07}
		fooSecretData := map[string]interface{}{
			"secret": "bar",
			"binary": binaryContent,
		}
		_, err := vault.Write("kv/data/foo", map[string]interface{}{
			"data": fooSecretData,
		})
		So(err, ShouldBeNil)

		Convey("GetWrappingToken", func() {
			Convey("returns token", func() {
				token, err := vault.GetWrappingToken()
				So(err, ShouldBeNil)
				So(token, ShouldNotEqual, vault.token)
			})
		})

		Convey("Encrypt and Decrypt", func() {
			toEncrypt := []byte("to-encrypt")
			encrypted, err := vault.Encrypt("my-trans", toEncrypt)
			So(err, ShouldBeNil)
			So(encrypted, ShouldNotResemble, toEncrypt)

			decrypted, err := vault.Decrypt("my-trans", encrypted)
			So(err, ShouldBeNil)
			So(decrypted, ShouldResemble, toEncrypt)
		})

		Convey("GetMap", func() {
			Convey("returns path as map of strings", func() {
				list, err := vault.GetMap("kv/data/foo")
				So(err, ShouldBeNil)
				So(list, ShouldResemble, map[string]interface{}{
					"secret": "bar",
					"binary": base64.StdEncoding.EncodeToString(binaryContent),
				})
			})
		})

		Convey("List", func() {
			Convey("returns keys as strings", func() {
				list, err := vault.ListPath("kv/metadata")
				So(err, ShouldBeNil)
				So(list, ShouldResemble, []string{"foo"})
			})
		})
	})
}
