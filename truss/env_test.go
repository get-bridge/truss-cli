package truss

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEnv(t *testing.T) {
	Convey("Env", t, func() {
		env := "edge-cmh"
		kubeconfigs := map[string]interface{}{
			"edge-cmh":    "kubeconfig-truss-nonprod-cmh",
			"staging-cmh": "kubeconfig-truss-nonprod-cmh",
			"staging-dub": "kubeconfig-truss-nonprod-dub",
		}
		kubeDir := "/home/test/.kube/"

		Convey("when kubeconfig is found", func() {
			input := &EnvInput{
				Env:         env,
				Kubeconfigs: kubeconfigs,
				KubeDir:     kubeDir,
			}
			output, err := Env(input)
			So(err, ShouldBeNil)
			So(output.Kubeconfig, ShouldEqual, "/home/test/.kube/kubeconfig-truss-nonprod-cmh")
		})

		Convey("handles case when no kubeconfig found", func() {
			input := &EnvInput{
				Env:         "bogus",
				Kubeconfigs: kubeconfigs,
				KubeDir:     kubeDir,
			}
			output, err := Env(input)
			So(output.Kubeconfig, ShouldBeEmpty)
			So(err.Error(), ShouldEqual, "No kubeconfig found for env bogus")
		})
	})
}

func TestBashFormat(t *testing.T) {
	Convey("formats Envs into bash command that can be eval'd", t, func() {
		environmentVars := &EnvironmentVars{
			Kubeconfig: "/home/test/.kube/kubeconfig-truss-nonprod-cmh",
		}

		output := environmentVars.BashFormat("edge-cmh")
		So(output, ShouldStartWith, "export KUBECONFIG=/home/test/.kube/kubeconfig-truss-nonprod-cmh")
		So(output, ShouldContainSubstring, "# Run this command to configure your shell:")
		So(output, ShouldContainSubstring, "# eval \"$(truss env -e edge-cmh)")
	})
}
