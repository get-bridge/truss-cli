package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVault(t *testing.T) {
	vault := createTestVault(t)

	Convey("Vault", t, func() {
		fooSecretData := map[string]string{
			"secret": "bar",
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
