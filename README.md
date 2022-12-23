# Kleister: CLI client

[![General Workflow](https://github.com/kleister/kleister-cli/actions/workflows/general.yml/badge.svg)](https://github.com/kleister/kleister-cli/actions/workflows/general.yml) [![Join the Matrix chat at https://matrix.to/#/#kleister:matrix.org](https://img.shields.io/badge/matrix-%23kleister-7bc9a4.svg)](https://matrix.to/#/#kleister:matrix.org) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/7143ea13bd644aa3be6749ca967be7d0)](https://www.codacy.com/gh/kleister/kleister-cli/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kleister/kleister-cli&amp;utm_campaign=Badge_Grade) [![Go Reference](https://pkg.go.dev/badge/github.com/kleister/kleister-cli.svg)](https://pkg.go.dev/github.com/kleister/kleister-cli) [![GitHub Repo](https://img.shields.io/badge/github-repo-yellowgreen)](https://github.com/kleister/kleister-cli)

Within this repository we are building the command-line client to interact with
the [Kleister API][api] server.

## Install

You can download prebuilt binaries from the [GitHub releases][releases] or from
our [download site][downloads]. If you prefer to use containers you could use
our images published on [Docker Hub][dockerhub] or [Quay][quay]. You are a Mac
user? Just take a look at our [homebrew formula][homebrew]. If you need further
guidance how to install this take a look at our [documentation][docs].

## Development

Make sure you have a working Go environment, for further reference or a guide
take a look at the [install instructions][golang]. This project requires
Go >= v1.18, at least that's the version we are using.

```console
git clone https://github.com/kleister/kleister-cli.git
cd kleister-cli

make generate build

./bin/kleister-cli -h
```

## Security

If you find a security issue please contact
[kleister@webhippie.de](mailto:kleister@webhippie.de) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

-   [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```

[api]: https://github.com/kleister/kleister-cli
[releases]: https://github.com/kleister/kleister-cli/releases
[downloads]: https://dl.kleister.eu/cli
[homebrew]: https://github.com/kleister/homebrew-kleister
[dockerhub]: https://hub.docker.com/r/kleister/kleister-cli/tags/
[quay]: https://quay.io/repository/kleister/kleister-cli?tab=tags
[docs]: https://kleister.eu/
[golang]: http://golang.org/doc/install.html
