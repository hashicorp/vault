---
layout: "install"
page_title: "Install Vault"
sidebar_current: "docs-install-install"
description: |-
  Learn how to install Vault.
---

# Install Vault

Installing Vault is simple. There are two approaches to installing Vault:
downloading a precompiled binary for your system, or installing from source.

Downloading a precompiled binary is easiest, and we provide downloads over
TLS along with SHA256 sums to verify the binary is what we say it is. We
also distribute a PGP signature with the SHA256 sums that can be verified.
However, we use a 3rd party storage host, and some people feel that
due to the importance of security with Vault, they'd rather compile it
from source.

For this reason, we also document on this page how to compile Vault
from source, from the same versions of all dependent libraries that
we used for the official builds.

## Precompiled Binaries

To install the precompiled binary,
[download](/downloads.html) the appropriate package for your system.
Vault is currently packaged as a zip file. We don't have any near term
plans to provide system packages.

Once the zip is downloaded, unzip it into any directory. The
`vault` binary inside is all that is necessary to run Vault (or
`vault.exe` for Windows). Any additional files, if any, aren't
required to run Vault.

Copy the binary to anywhere on your system. If you intend to access it
from the command-line, make sure to place it somewhere on your `PATH`.

## Compiling from Source

Check [Vault's README on
GitHub](https://github.com/hashicorp/vault#developing-vault) for up-to-date
information about compiling Go from source.
