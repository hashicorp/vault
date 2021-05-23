# Changelog

## [0.1.1] - 2021-03-11

This bugfix release backports fixes for a few issues:

* Keep leading empty lines in func bodies if they help readability

* Avoid breaking comment alignment on empty field lists

* Add support for `//go-sumtype:` directives

## [0.1.0] - 2021-01-05

This is gofumpt's first release, based on Go 1.15.x. It solidifies the features
which have worked well for over a year.

This release will be the last to include `gofumports`, the fork of `goimports`
which applies `gofumpt`'s rules on top of updating the Go import lines. Users
who were relying on `goimports` in their editors or IDEs to apply both `gofumpt`
and `goimports` in a single step should switch to gopls, the official Go
language server. It is supported by many popular editors such as VS Code and
Vim, and already bundles gofumpt support. Instructions are available [in the
README](https://github.com/mvdan/gofumpt).

`gofumports` also added maintenance work and potential confusion to end users.
In the future, there will only be one way to use `gofumpt` from the command
line. We also have a [Go API](https://pkg.go.dev/mvdan.cc/gofumpt/format) for
those building programs with gofumpt.

Finally, this release adds the `-version` flag, to print the tool's own version.
The flag will work for "master" builds too.

[0.1.1]: https://github.com/mvdan/gofumpt/releases/tag/v0.1.1
[0.1.0]: https://github.com/mvdan/gofumpt/releases/tag/v0.1.0
