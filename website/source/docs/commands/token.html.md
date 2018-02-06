---
layout: "docs"
page_title: "token - Command"
sidebar_current: "docs-commands-token"
description: |-
  The "token" command groups subcommands for interacting with tokens. Users can
  create, lookup, renew, and revoke tokens.
---

# token

The `token` command groups subcommands for interacting with tokens. Users can
create, lookup, renew, and revoke tokens.

For more information on tokens, please see the [token concepts
page](/docs/concepts/tokens.html).

## Examples

Create a new token:

```text
$ vault token create
```

Revoke a token:

```text
$ vault token revoke 96ddf4bc-d217-f3ba-f9bd-017055595017
```

Renew a token:

```text
$ vault token renew 96ddf4bc-d217-f3ba-f9bd-017055595017
```

## Usage

```text
Usage: vault token <subcommand> [options] [args]

  # ...

Subcommands:
    capabilities    Print capabilities of a token on a path
    create          Create a new token
    lookup          Display information about a token
    renew           Renew a token lease
    revoke          Revoke a token and its children
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
