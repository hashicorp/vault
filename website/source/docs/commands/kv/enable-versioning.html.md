---
layout: "docs"
page_title: "kv enable-versioning - Command"
sidebar_title: "<code>enable-versioning</code>"
sidebar_current: "docs-commands-kv-enable-versioning"
description: |-
  The "kv enable-versioning" command turns on versioning for the backend
  at the provided path.
---

# kv enable-versioning

The `kv enable-versioning` command turns on versioning for an existing
non-versioned key/value secrets engine (K/V Version 1) at its path.

## Examples

This command turns on versioning for the K/V Version 1 secrets engine enabled at
"secret".

```text
$ vault kv enable-versioning secret
Success! Tuned the secrets engine at: secret/
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
