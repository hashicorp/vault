<!-- markdownlint-disable first-line-h1 no-inline-html -->

[![Build](https://github.com/vmware/govmomi/actions/workflows/govmomi-build.yaml/badge.svg)][ci-build]
[![Tests](https://github.com/vmware/govmomi/actions/workflows/govmomi-go-tests.yaml/badge.svg)][ci-tests]
[![Go Report Card](https://goreportcard.com/badge/github.com/vmware/govmomi)][go-report-card]
[![Latest Release](https://img.shields.io/github/release/vmware/govmomi.svg?logo=github&style=flat-square)][latest-release]
[![Go Reference](https://pkg.go.dev/badge/github.com/vmware/govmomi.svg)][go-reference]
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/vmware/govmomi)][go-version]

# govmomi

A Go library for interacting with VMware vSphere APIs (ESXi and/or vCenter Server).

In addition to the vSphere API client, this repository includes:

* [govc][govc] - vSphere CLI
* [vcsim][vcsim] - vSphere API mock framework
* [toolbox][toolbox] - VM guest tools framework

## Compatibility

vSphere 7.0 and higher.

## Documentation

The APIs exposed by this library closely follow the API described in the [VMware vSphere API Reference Documentation][reference-api]. Refer to the documentation to become familiar with the upstream API.

The code in the `govmomi` package is a wrapper for the code that is generated from the vSphere API description. It primarily provides convenience functions for working with the vSphere API. See [godoc.org][reference-godoc] for documentation.

## Installation

### Binaries and Docker Images for `govc` and `vcsim`

Installation instructions, released binaries, and Docker images are documented in the respective README files of [`govc`][govc] and [`vcsim`][vcsim].

## Discussion

Collaborate with the community using GitHub [discussions][govmomi-github-discussions] and GitHub [issues][govmomi-github-issues].

## Status

Changes to the API are subject to [semantic versioning][reference-semver].

Refer to the [CHANGELOG][govmomi-changelog] for version to version changes.

## Related Projects

* [pyvmomi][reference-pyvmomi]
* [rbvmomi][reference-rbvmomi]

## License

govmomi is available under the [Apache 2 License][govmomi-license].

## Name

Pronounced: _go·​v·​mom·​e_

Follows pyvmomi and rbvmomi: language prefix + the vSphere acronym "VM Object Management Infrastructure".

[//]: Links

[ci-build]: https://github.com/vmware/govmomi/actions/workflows/govmomi-build.yaml
[ci-tests]: https://github.com/vmware/govmomi/actions/workflows/govmomi-go-tests.yaml
[latest-release]: https://github.com/vmware/govmomi/releases/latest
[govc]: govc/README.md
[govmomi-github-issues]: https://github.com/vmware/govmomi/issues
[govmomi-github-discussions]: https://github.com/vmware/govmomi/discussions
[govmomi-changelog]: CHANGELOG.md
[govmomi-license]: LICENSE.txt
[go-reference]: https://pkg.go.dev/github.com/vmware/govmomi
[go-report-card]: https://goreportcard.com/report/github.com/vmware/govmomi
[go-version]: https://github.com/vmware/govmomi
[reference-api]: https://developer.broadcom.com/xapis/vsphere-web-services-api/latest/
[reference-godoc]: https://pkg.go.dev/github.com/vmware/govmomi
[reference-pyvmomi]: https://github.com/vmware/pyvmomi
[reference-rbvmomi]: https://github.com/vmware/rbvmomi
[reference-semver]: http://semver.org
[toolbox]: toolbox/README.md
[vcsim]: vcsim/README.md
