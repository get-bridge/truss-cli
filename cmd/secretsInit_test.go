package cmd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSecretsInit(t *testing.T) {
	Convey("secrets init", t, func() {
		Convey("runs with no errors", nil)
	})
	Convey("generateSecretsFile", t, func() {
		environments := map[string]string{
			"edge-cmh": "kubeconfig-edge-cmh",
		}
		Convey("returns yaml", func() {
			result := generateSecretsFile(environments, "service-foo", "secretsPath", "secretsFile", "vaultPath", "fileName.yaml")
			So(result, ShouldEqual, `# fileName.yaml
transit-key-name: service-foo

secrets:
- name: edge-cmh
  kubeconfig: kubeconfig-edge-cmh
  vaultPath: secretsPath/edge/cmh/secretsFile
  filePath: vaultPath/edge/cmh/service-foo
`)
		})
	})
}
