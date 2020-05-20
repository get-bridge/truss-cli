# truss-cli

CLI to help you manage many k8s clusters

## Contributing

### Install

#### Mac Homebrew

```sh
brew install instructure/truss-cli/truss-cli
```

#### GO

```sh
go get github.com/instructure/truss-cli truss
```

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
