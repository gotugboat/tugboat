# Tugboat

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gotugboat/tugboat)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gotugboat/tugboat)
![GitHub](https://img.shields.io/github/license/gotugboat/tugboat?color=blue)

Tugboat is an open source command-line tool for building multi-architecture container images. Supported on popular platforms, Tugboat simplifies the process of creating images for different architectures. Create and deploy your images to popular container registries like Docker Hub, Quay (and more), or a self hosted private registry using just a few simple commands.

> Please note that the `main` branch may be in an unstable or even broken state during development. To get a stable version of the binary, see the [releases](https://github.com/gotugboat/tugboat/releases) page to download the latest version.

## Getting Started
<!-- TODO -->

## Development

To develop on Tugboat follow these instructions to get your local environment all set up.

### Prerequisites

- Review the [community code of conduct](./.github/CODE_OF_CONDUCT.md) and [contributions guidelines](./.github/CONTRIBUTING.md).
- Installation of [Golang](https://go.dev/dl/).
- An IDE that supports Go, such as [Visual Studio Code](https://code.visualstudio.com/) with the [Go plugin](https://code.visualstudio.com/docs/languages/go).

### Install dependencies

Install the Go dependencies with:

```bash
go mod download
```

### Run the tests

```
make test
```

### Build Tugboat locally

```
make build
```

## Contributing

Contributions are welcomed! Please refer to the [contributions guidelines](./.github/CONTRIBUTING.md) for information about our code of conduct and the process for submitting pull requests.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [releases page](https://github.com/gotugboat/tugboat/releases).

## Licensing

Tugboat is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for the full license text.
