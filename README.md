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

## Prerequisites

We use [mise][mise] to manage all required tools and their versions. Install it
by following the [official installation instructions][mise-install], then run
the following commands inside the repository to activate mise and install all
tools defined in `mise.toml`:

```console
mise trust
mise install
```

## Build

Since all required commands ar part of our [go-task][gotask] taskfile the
commands you got to execute are quite simple:

```console
git clone https://github.com/kleister/kleister-cli.git
cd kleister-cli

task generate build
./bin/kleister-cli -h
```

## Security

If you find a security issue please contact
[kleister@webhippie.de](mailto:kleister@webhippie.de) first.

## Contributing

Generally we are following [conventional commits][commits] when we apply
changes. That way we are able to generate proper changelogs for every release.
Please use always pull requests to integrate new functionalities or to fix
issues.

For the release process we are following [semantic versioning][semver] which
clearly indicates if a new version just resolves bugs, includes new features or
even includes breaking changes.

After installing the tools via `mise install` as described above set up the
pre-commit hooks so they run automatically on every commit:

```console
pre-commit install --hook-type pre-commit --hook-type commit-msg
```

> `pre-commit` is managed by mise and will be available after `mise install`.

If you have changed something on the source you should simply commit following
the mentioned conventions:

```console
git checkout -b feat/new-feature
git add --all
git commit -m 'feat: added awesome new feature'
git push --set-upstream origin feat/new-feature
```

After pushing your changes into the Git repository you should create a pull
request on GitHub. If the pull request have been merged and everything built
fine it will also create automatically a new release at least once a week.

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
[cloudsmith]: https://cloudsmith.com/
[gotask]: https://taskfile.dev/installation/
[devcontainer]: https://containers.dev/
[mise]: https://mise.jdx.dev/
[mise-install]: https://mise.jdx.dev/getting-started.html
[commits]: https://www.conventionalcommits.org/en/v1.0.0/
[semver]: https://semver.org/
