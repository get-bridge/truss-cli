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
    awsrole: arn:aws:iam::1234567:role/xacct/my-s3-role
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

## Usage

### Secrets

The `truss secrets` command makes it easier to manage secrets in git, and
synchronize them across multiple Truss Vault servers. Start by creating a
`secrets.yaml` file.

```yaml
# secrets.yaml
transit-key-name: my-project
secrets:
  - name: app-1
    kubeconfig: kubeconfig-truss-nonprod-cmh # relative to `kubeconfigfiles.directory` in `~/.truss.yaml`
    filePath: ./secrets/edge-cmh # relative to `pwd`
    vaultPath: secret/bridge/edge/cmh/app-1 # name of folder for multiple vault secrets
  - name: app-1
    kubeconfig: kubeconfig-truss-prod-cmh
    filePath: ./secrets/prod-cmh
    vaultPath: secret/bridge/prod/cmh/app-1
  - name: app-2
    kubeconfig: kubeconfig-truss-nonprod-cmh
    filePath: ./secrets/edge-cmh
    vaultPath: secret/bridge/edge/cmh/app-2
  - name: app-2
    kubeconfig: kubeconfig-truss-prod-cmh
    filePath: ./secrets/prod-cmh
    vaultPath: secret/bridge/prod/cmh/app-2
```

Then, run `truss secrets edit app-1 kubeconfig-truss-nonprod-cmh`. This will
open your `$EDITOR` with a file containing `secrets: {}`. An example secrets
file might look like this:

```yaml
secrets:
  web:
    SOME_API_KEY: a_super_secret_secret
    CSRF_SECRET: a_super_secret_secret
  db:
    DB_USERNAME: root
    DB_PASSWORD: a_super_secret_secret
```

Running `truss secrets push app-1 kubeconfig-truss-nonprod-cmh` will create two
secrets in Vault, each containing their corresponding key-vaule pairs.

- `secrets/bridge/edge/cmh/app-1/web`
- `secrets/bridge/edge/cmh/app-1/db`

Create multiple environments with `secrets.yaml` and
`truss secrets edit <name> <kubeconfig>`, then you can run
`truss secrets push --all` to update all secrets.

When in doubt, you can run `truss secrets pull --all` to update the files on
disk with the values from Vault. Note: this action is destructive!

## Contributing

### Development

We are using [cobra][1] for cli parsing. For ease in
adding new commands install the cobra cli.

### Adding new command

```sh
cobra add my-new-command
```

### Testing
Add a `.envrc` file that looks like this:
```
export TEST_S3_BUCKET=truss-kubeconfig-us-east-2
export TEST_AWS_ROLE=<role with access to bucket>
export TEST_VAULT_ROLE=<vault role with rw access to secrets>
export TEST_S3_BUCKET_REGION=us-east-2
export TEST_GLOBAL_CONFIG_BUCKET=truss-cli-global-config
export TEST_GLOBAL_CONFIG_KEY=.truss.yaml
```

Load the those into your shell:
```sh
source .envrc
```

Run all tests:

```sh
go test go test ./...
```

Run a specific test using the `-run` flag and a regex for the test name:

```
go test go test ./... -run ^TestWrap$
```

[1]: https://github.com/spf13/cobra#installing

### Publish release

Update version in Makefile. Commit to master.

```sh
make release
```
