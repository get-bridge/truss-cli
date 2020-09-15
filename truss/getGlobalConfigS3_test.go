package truss

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetGlobalConfigS3(t *testing.T) {
	Convey("GetGlobalConfigS3", t, func() {
		dir, err := ioutil.TempDir("", "")
		So(err, ShouldBeNil)
		defer os.Remove(dir)

		input := &GetGlobalConfigS3Input{
			Dir: dir,
		}

		Convey("fails ", func() {
			Convey("fails ", func() {
				_, err := GetGlobalConfigS3(input)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "MissingRegion")
			})
		})
	})
}
