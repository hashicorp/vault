---
layout: "docs"
page_title: "token capabilities - Command"
sidebar_current: "docs-commands-token-capabilities"
description: |-
  The "token capabilities" command fetches the capabilities of a token for a
  given path.
---

# token capabilities

The `token capabilities` command fetches the capabilities of a token for a given
path.

If a TOKEN is provided as an argument, this command uses the "/sys/capabilities"
endpoint and permission. If no TOKEN is provided, this command uses the
"/sys/capabilities-self" endpoint and permission with the locally authenticated
token.

## Examples

List capabilities for the local token on the "secret/foo" path:

```text
$ vault token capabilities secret/foo
read
```

List capabilities for a token on the "cubbyhole/foo" path:

```text
$ vault token capabilities 96ddf4bc-d217-f3ba-f9bd-017055595017 database/creds/readonly
deny
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.