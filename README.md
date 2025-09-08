# Kleister: CLI client

[![General Workflow](https://github.com/kleister/kleister-cli/actions/workflows/general.yml/badge.svg)](https://github.com/kleister/kleister-cli/actions/workflows/general.yml) [![Join the Matrix chat at https://matrix.to/#/#kleister:matrix.org](https://img.shields.io/badge/matrix-%23kleister-7bc9a4.svg)](https://matrix.to/#/#kleister:matrix.org) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/bf217141b736464297b3aa6001162cbd)](https://app.codacy.com/gh/kleister/kleister-cli/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade) [![Go Reference](https://pkg.go.dev/badge/github.com/kleister/kleister-cli.svg)](https://pkg.go.dev/github.com/kleister/kleister-cli) [![GitHub Repo](https://img.shields.io/badge/github-repo-yellowgreen)](https://github.com/kleister/kleister-cli) [![Hosted By: Cloudsmith](https://img.shields.io/badge/OSS%20hosting%20by-cloudsmith-blue?logo=cloudsmith&style=flat-square)](https://cloudsmith.com)

> [!CAUTION]
> This project is in active development and does not provide any stable release
> yet, you can expect breaking changes until our first real release!

Within this repository we are building the command-line client to interact with
the [Kleister API][api] server.

## Install

You can download prebuilt binaries from the [GitHub releases][releases] or from
our [download site][downloads]. Besides that we also prepared repositories for
DEB and RPM packages which can be found at [Cloudsmith][pkgrepo]. If you prefer
to use containers you could use our images published on [GHCR][ghcr],
[Docker Hub][dockerhub] or [Quay][quay]. You are a Mac user? Just take a look
at our [homebrew formula][homebrew]. If you need further guidance how to
install this take a look at our [documentation][docs].

Package repository hosting is graciously provided by [Cloudsmith][cloudsmith].
Cloudsmith is the only fully hosted, cloud-native, universal package management
solution, that enables your organization to create, store and share packages in
any format, to any place, with total confidence.

## Build

If you are not familiar with [Nix][nix] it is up to you to have a working
environment for Go (>= 1.24.0) as the setup won't we covered within this guide.
Please follow the official install instructions for [Go][golang]. Beside that
we are using [go-task][gotask] to define all commands to build this project.

```console
git clone https://github.com/kleister/kleister-cli.git
cd kleister-cli

task generate build
./bin/kleister-cli -h
```

If you got [Nix][nix] and [Direnv][direnv] configured you can simply execute
the following commands to get al dependencies including [go-task][gotask] and
the required runtimes installed:

```console
cat << EOF > .envrc
use flake . --impure --extra-experimental-features nix-command
EOF

direnv allow
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
[downloads]: https://dl.kleister.eu
[homebrew]: https://github.com/kleister/homebrew-kleister
[ghcr]: https://github.com/orgs/kleister/packages
[dockerhub]: https://hub.docker.com/r/kleister/kleister-cli/tags/
[quay]: https://quay.io/repository/kleister/kleister-cli?tab=tags
[docs]: https://kleister.eu/
[nix]: https://nixos.org/
[golang]: http://golang.org/doc/install.html
[gotask]: https://taskfile.dev/installation/
[direnv]: https://direnv.net/
[pkgrepo]: https://cloudsmith.io/~kleister/repos/general/groups/
[cloudsmith]: https://cloudsmith.com/
