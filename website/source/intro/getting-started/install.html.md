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
`vault`, you should see help output similar to the following:

```
$ vault
usage: vault [--version] [--help] <command> [<args>]

Available commands are:
audit-disable    Disable an audit backend
audit-enable     Enable an audit backend
audit-list       Lists enabled audit backends in Vault
auth             Prints information about how to authenticate with Vault
auth-disable     Disable an auth provider
auth-enable      Enable a new auth provider
delete           Delete operation on secrets in Vault
help             Look up the help for a path
init             Initialize a new Vault server
mount            Mount a logical backend
mounts           Lists mounted backends in Vault
policies         List the policies on the server
policy-write     Write a policy to the server
read             Read data or secrets from Vault
remount          Remount a secret backend to a new path
renew            Renew the lease of a secret
revoke           Revoke a secret.
seal             Seals the vault server
seal-status      Outputs status of whether Vault is sealed
server           Start a Vault server
token-create     Create a new auth token
token-revoke     Revoke one or more auth tokens
unmount          Unmount a secret backend
unseal           Unseals the vault server
version          Prints the Vault version
write            Write secrets or configuration into Vault
```

If you get an error that Vault could not be found, then your PATH environment
variable was not setup properly. Please go back and ensure that your PATH
variable contains the directory where Vault was installed.

Otherwise, Vault is installed and ready to go!
