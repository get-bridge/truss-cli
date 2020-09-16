module github.com/instructure-bridge/truss-cli

go 1.14

require (
	github.com/Songmu/prompter v0.3.0
	github.com/aws/aws-sdk-go v1.31.2
	github.com/creasty/defaults v1.5.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/hashicorp/vault v1.5.3
	github.com/hashicorp/vault/api v1.0.5-0.20200630205458-1a16f3c699c6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools/v3 v3.0.2 // indirect
	k8s.io/api v0.15.12
	k8s.io/apimachinery v0.15.12
	k8s.io/client-go v0.15.12
)

replace (
	github.com/hashicorp/vault/api => github.com/hashicorp/vault/api v1.0.5-0.20200916152957-e1fe191e6ba0
	k8s.io/api => k8s.io/api v0.15.12
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.12
	k8s.io/client-go => k8s.io/client-go v0.15.12
)
