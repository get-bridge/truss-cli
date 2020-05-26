# truss-cli

CLI to help you manage many k8s clusters

## Configuration

Defaults to using a file `~/.truss.yaml` for configuration.

Example:

```yaml
dependencies:
  - kubectl
  - sshuttle
  - vault
kubeconfigfiles:
  s3:
    bucket: my-bucket-with-kubeconfigs
    region: us-east-2
environments:
  edge-cmh: kubeconfig-truss-nonprod-cmh
  staging-cmh: kubeconfig-truss-nonprod-cmh
  staging-dub: kubeconfig-truss-nonprod-dub
  staging-syd: kubeconfig-truss-nonprod-syd
  prod-cmh: kubeconfig-truss-prod-cmh
  prod-dub: kubeconfig-truss-prod-dub
  prod-syd: kubeconfig-truss-prod-syd
vault:
  auth:
    aws:
      vaultrole: admin
      awsrole: arn:aws:iam::1234567:role/xacct/my-role
```

## Install

### Mac Homebrew

```sh
brew install instructure-bridge/truss-cli/truss-cli
```

### GO

```sh
go get github.com/instructure-bridge/truss-cli truss
```

## Contributing

### Development

We are using [cobra][1] for cli parsing. For ease in
adding new commands install the cobra cli.

### Adding new command

```sh
cobra add my-new-command
```

[1]: https://github.com/spf13/cobra#installing

### Publish release

Update version in Makefile. Commit to master.

```sh
make release
```
