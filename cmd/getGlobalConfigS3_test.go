package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestGetGlobalConfigS3(t *testing.T) {
	Convey("GetGlobalConfigS3", t, func() {
		viper.Reset()

		Convey("return errors", func() {
			err := getGlobalConfigS3Cmd.RunE(getGlobalConfigS3Cmd, []string{})
			So(err, ShouldNotBeNil)
		})
	})
}
