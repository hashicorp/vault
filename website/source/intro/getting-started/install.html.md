---
layout: "intro"
page_title: "Install Vault - Getting Started"
sidebar_current: "gettingstarted-install"
description: |-
  The first step to using Vault is to get it installed.
---

# Install Vault

Vault must first be installed on your machine. Vault is distributed as
a [binary package](/downloads.html) for all supported platforms and
architectures. This page will not cover how to compile Vault from source,
but compiling from source is covered in the [documentation](/docs/install/index.html)
for those who want to be sure they're compiling source they trust into
the final binary.

## Installing Vault

To install Vault, find the [appropriate package](/downloads.html) for
your system and download it. Vault is packaged as a zip archive.

After downloading Vault, unzip the package. Vault runs as a single binary
named `vault`. Any other files in the package can be safely removed and
Vault will still function.

The final step is to make sure that the `vault` binary is available on the `PATH`.
See [this page](https://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux)
for instructions on setting the PATH on Linux and Mac.
[This page](https://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows)
contains instructions for setting the PATH on Windows.

## Verifying the Installation

After installing Vault, verify the installation worked by opening a new
terminal session and checking that the `vault` binary is available. By executing
`vault`, you should see help output similar to the following:

```text
$ vault
Usage: vault <command> [args]

Common commands:
    read        Read data and retrieves secrets
    write       Write data, configuration, and secrets
    delete      Delete secrets and configuration
    list        List data or secrets
    login       Authenticate locally
    server      Start a Vault server
    status      Print seal and HA status
    unwrap      Unwrap a wrapped secret

Other commands:
    audit          Interact with audit devices
    auth           Interact with auth methods
    lease          Interact with leases
    operator       Perform operator-specific tasks
    path-help      Retrieve API help for paths
    policy         Interact with policies
    secrets        Interact with secrets engines
    ssh            Initiate an SSH session
    token          Interact with tokens
```

If you get an error that the binary could not be found, then your `PATH`
environment variable was not setup properly. Please go back and ensure that your
`PATH` variable contains the directory where Vault was installed.

Otherwise, Vault is installed and ready to go!

## Command Completion

Vault also includes command-line completion for subcommands, flags, and path
arguments where supported. To install command-line completion, you must be using
Bash, ZSH or Fish. Unfortunately other shells are not supported at this time.

To install completions, run:

```sh
$ vault -autocomplete-install
```

This will automatically install the helpers in your `~/.bashrc` or `~/.zshrc`, or to
`~/.config/fish/completions/vault.fish` for Fish users. Then restart your terminal
or reload your shell:
```sh
$ exec $SHELL
```

Now when you type `vault <tab>`, Vault will suggest options. This is very
helpful for beginners and advanced Vault users.

## Next

Now Vault is installed we can start our first Vault server! [Let's do
that now](/intro/getting-started/dev-server.html).
