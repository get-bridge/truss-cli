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

Grab yourself a [personal access token](https://github.com/settings/tokens/new?scopes=repo&description=Homebrew%20for%20Bridge%20VPN%20CLI). Then...

```sh
brew install get-bridge/truss-cli/truss-cli
```

### GO

```sh
go get github.com/get-bridge/truss-cli truss
```

## Usage

### Bootstrapping a Truss Deployment

With the `truss bootstrap` command, you can quickly and easily bootstrap your
Truss deployment configuration using one of our pre-made templates. To get
started, you'll need to configure the bootstrapper by creating a
`bootstrap.truss.yaml` file. It should look something like this:

```yaml
# The following templateSource configuration represents the default values,
# which means you'll need to have insopshub credentials loaded in order to
# assume the ops-admin role. You only need to include any of this in your local
# configuration file if you intend to override these defaults.
templateSource:
  type: s3
  local:
    directory: ./bootstrap-templates
  s3:
    bucket: truss-cli-global-config
    region: us-east-2
    prefix: bootstrap-templates
    role: arn:aws:iam::127178877223:role/xacct/ops-admin
  git:
    clone_url: git@github.com:get-bridge/truss-cli.git
    directory: bootstrap-templates
    checkout_ref: refs/heads/master

# template represents which template to render. The default is "default"
template: default
# trussDir is the directory where deployment configuration will be rendered. The
# default is "truss"
trussDir: truss
# params represent your values for the given template's parameters. They are
# defined in the template in the `.truss-manifest.yaml` file. The default values
# for the default template are included here.
params:
  name: ""
  role: ""
  httpPort: ""
  healthCheckPath: ""
  image: ""
  slackChannel: ""
```

> You can also provide param values by passing `--set name=value` to the
> `truss bootstrap` command.

With your configuration file in place at the root of your project, simply run
`truss bootstrap` to create your local `./truss` directory. This deployment
config will serve as a starting point for your project, and it is expected that
you will make changes per your application's needs. Thus, your
`bootstrap.truss.yaml` is no longer necessary.

In the future, we might strive to make this template customizable enough such
that you could keep your `bootstrap.truss.yaml` and re-generate as we publish
updates to the template. For now, it's one-and-done!

### Secrets

The `truss secrets` command makes it easier to manage secrets in git, and
synchronize them across multiple Truss Vault servers. Start by creating a
`secrets.yaml` file or running `truss secrets init`

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
