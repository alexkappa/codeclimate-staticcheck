# Code Climate Staticcheck Engine

`codeclimate-staticcheck` is a Code Climate engine that performs static analysis for Go programs, finding all kinds of bugs – crashes as well as incorrect behaviour. You can run it on your command line using the Code Climate CLI, or on our hosted analysis platform.

### Installation

1. If you haven't already, [install the Code Climate CLI](https://github.com/codeclimate/codeclimate).
2. Run `codeclimate engines:enable staticcheck`. This command both installs the engine and enables it in your `.codeclimate.yml` file.
3. You're ready to analyze! Browse into your project's folder and run `codeclimate analyze`.

### Configuration

`staticcheck` needs to build your code so it can perform the analysis. For that reason a valid `GOPATH` is necessary within the engine. By default `staticcheck` will copy your code into `$GOPATH/src/app`. This can work for single package projects.

For multiple package projects such as application, you may need to configure the root package as you would for any project being built by `go build`. You can specify the root package in your `.codeclimate.yml`:

```yaml
engines:
  staticcheck:
    enabled: true
    config:
      package: github.com/myname/myapp
```

### Building

```console
docker build -t codeclimate/codeclimate-staticcheck .
```

During development you may find it useful to run `codeclimate analyze --dev` which will look for the engines image locally.

### Need help?

For help with `staticcheck`, [check out the documentation](https://staticcheck.io).

If you're running into a Code Climate issue, first look over this project's [GitHub Issues](https://github.com/alexkappa/codeclimate-staticcheck/issues), as your question may have already been covered. If not, [go ahead and open a support ticket with us](https://codeclimate.com/help).
