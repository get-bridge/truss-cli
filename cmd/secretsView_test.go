package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestSecretsView(t *testing.T) {
	Convey("secrets view", t, func() {
		c := &cobra.Command{}

		tmpFile, err := ioutil.TempFile("", "")
		So(err, ShouldBeNil)
		defer os.Remove(tmpFile.Name())
		tmpFile.WriteString("transit-key-name: omg-bbq")
		viper.Set("TRUSS_SECRETS_FILE", tmpFile.Name())

		Convey("errors if no such configuration", func() {
			err := secretsViewCmd.RunE(c, []string{"secret-name", "kubeconfig-name"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "secret named 'secret-name' in 'kubeconfig-name' not found")
		})

		// TODO how can we mock view?
		Convey("runs with no errors", nil)
	})
}
