---
layout: "docs"
page_title: "kv list - Command"
sidebar_title: "<code>list</code>"
sidebar_current: "docs-commands-kv-list"
description: |-
  The "kv list" command lists data from Vault's K/V secrets engine at the given
  path.
---

# kv list

The `kv list` command returns a list of key names at the specified location.
Folders are suffixed with /. The input must be a folder; list on a file will not
return a value. Note that no policy-based filtering is performed on keys; do not
encode sensitive information in key names. The values themselves are not
accessible via this command.

Use this command to list all existing key names at a specific path.

## Examples

List values under the key "my-app":

```text
$ vault kv list secret/my-app/
Keys
----
admin_creds
domain
eng_creds
qa_creds
release
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
