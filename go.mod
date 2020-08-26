module github.com/instructure-bridge/truss-cli

go 1.14

require (
	github.com/Songmu/prompter v0.3.0
	github.com/aws/aws-sdk-go v1.31.2
	github.com/creasty/defaults v1.5.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.2.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.15.12
	k8s.io/apimachinery v0.15.12
	k8s.io/client-go v0.15.12
)

replace (
	k8s.io/api => k8s.io/api v0.15.12
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.12
	k8s.io/client-go => k8s.io/client-go v0.15.12
)
