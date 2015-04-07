---
layout: "intro"
page_title: "Install Vault"
sidebar_current: "gettingstarted-install"
description: |-
  The first step to using Vault is to get it installed.
---

# Install Vault

Vault must first be installed on your machine. Vault is distributed as
a [binary package](/downloads.html) for all supported platforms and
architectures. This page will not cover how to compile Vault from source,
but compiling from source is covered in the [documentation](#)
for those who want to be sure they're compiling source they trust into
the final binary.

## Installing Vault

To install Vault, find the [appropriate package](/downloads.html) for
your system and download it. Vault is packaged as a zip archive.

After downloading Vault, unzip the package. Vault runs as a single binary
named `vault`. Any other files in the package can be safely removed and
Vault will still function.

The final step is to make sure that `vault` is available on the PATH.
See [this page](http://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux)
for instructions on setting the PATH on Linux and Mac.
[This page](http://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows)
contains instructions for setting the PATH on Windows.

## Verifying the Installation

After installing Vault, verify the installation worked by opening a new
terminal session and checking that `vault` is available. By executing
`vault`, you should see help output similar to that below:

```
TODO
```

If you get an error that Vault could not be found, then your PATH environment
variable was not setup properly. Please go back and ensure that your PATH
variable contains the directory where Vault was installed.

Otherwise, Vault is installed and ready to go!
