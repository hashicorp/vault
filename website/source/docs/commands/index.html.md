---
layout: "docs"
page_title: "Commands (CLI)"
sidebar_current: "docs-commands"
description: |-
  Vault can be controlled via a command-line interface. This page documents all the commands Vault accepts.
---

# Vault Commands (CLI)

Vault is controlled via a very easy to use command-line interface (CLI).
Vault is only a single command-line application: `vault`. This application
then takes a subcommand such as "read" or "write". The complete list of
subcommands is in the navigation to the left.

The Vault CLI is a well-behaved command line application. In erroneous cases,
a non-zero exit status will be returned. It also responds to `-h` and `--help`
as you'd most likely expect.

To view a list of the available commands at any time, just run Vault
with no arguments. To get help for any specific subcommand, run the subcommand
with the `-h` argument.

The help output is very comprehensive, so we defer you to that for documentation.
We've included some guides to the left of common interactions with the
CLI.

## Autocompletion

The `vault` command features opt-in subcommand autocompletion that you can
enable for your shell with `vault -autocomplete-install`. After doing so,
you can invoke a new shell and use the feature.

For example, assume a tab is typed at the end of each prompt line:

```
$ vault au
audit-disable  audit-enable  audit-list  auth  auth-disable  auth-enable

$ vault s
seal  server  ssh  status  step-down
```
