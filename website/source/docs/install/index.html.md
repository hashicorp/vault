---
layout: "docs"
page_title: "Install Vault"
sidebar_current: "docs-install-install"
description: |-
  Installing Vault is simple. You can download a precompiled binary or compile
  from source. This page details both methods.
---

# Install Vault

Installing Vault is simple. There are two approaches to installing Vault:

1. Using a [precompiled binary](#precompiled-binaries)

1. Installing [from source](#compiling-from-source)

Downloading a precompiled binary is easiest, and we provide downloads over TLS
along with SHA256 sums to verify the binary. We also distribute a PGP signature
with the SHA256 sums that can be verified.

## Precompiled Binaries

To install the precompiled binary, [download](/downloads.html) the appropriate
package for your system. Vault is currently packaged as a zip file. We do not
have any near term plans to provide system packages.

Once the zip is downloaded, unzip it into any directory. The `vault` binary
inside is all that is necessary to run Vault (or `vault.exe` for Windows). Any
additional files, if any, aren't required to run Vault.

Copy the binary to anywhere on your system. If you intend to access it from the
command-line, make sure to place it somewhere on your `PATH`.

## Compiling from Source

To compile from source, you will need [Go](https://golang.org) installed and
configured properly (including a `GOPATH` environment variable set), as well
as a copy of [`git`](https://www.git-scm.com/) in your `PATH`.

  1. Clone the Vault repository from GitHub into your `GOPATH`:

    ```shell
    $ mkdir -p $GOPATH/src/github.com/hashicorp && cd $!
    $ git clone https://github.com/hashicorp/vault.git
    $ cd vault
    ```

  1. Bootstrap the project. This will download and compile libraries and tools
  needed to compile Vault:

    ```shell
    $ make bootstrap
    ```

  1. Build Vault for your current system and put the
  binary in `./bin/` (relative to the git checkout). The `make dev` target is
  just a shortcut that builds `vault` for only your local build environment (no
  cross-compiled targets).

    ```shell
    $ make dev
    ```

## Verifying the Installation

To verify Vault is properly installed, run `vault -h` on your system. You should
see help output. If you are executing it from the command line, make sure it is
on your PATH or you may get an error about Vault not being found.

```shell
$ vault -h
```
