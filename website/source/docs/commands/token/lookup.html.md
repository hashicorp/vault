---
layout: "docs"
page_title: "token lookup - Command"
sidebar_current: "docs-commands-token-lookup"
description: |-
  The "token lookup" displays information about a token or accessor. If a TOKEN
  is not provided, the locally authenticated token is used.
---

# token lookup

The `token lookup` displays information about a token or accessor. If a TOKEN is
not provided, the locally authenticated token is used.

## Examples

Get information about the locally authenticated token (this uses the
`/auth/token/lookup-self` endpoint and permission):

```text
$ vault token lookup
```

Get information about a particular token (this uses the `/auth/token/lookup`
endpoint and permission):

```text
$ vault token lookup 96ddf4bc-d217-f3ba-f9bd-017055595017
```

Get information about a token via its accessor:

```text
$ vault token lookup -accessor 9793c9b3-e04a-46f3-e7b8-748d7da248da
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(default: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-accessor` `(bool: false)` - Treat the argument as an accessor instead of a
  token. When this option is selected, the output will NOT include the token.
