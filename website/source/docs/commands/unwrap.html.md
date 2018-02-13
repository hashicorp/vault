---
layout: "docs"
page_title: "unwrap - Command"
sidebar_current: "docs-commands-unwrap"
description: |-
  The "unwrap" command unwraps a wrapped secret from Vault by the given token.
  The result is the same as the "vault read" operation on the non-wrapped
  secret. If no token is given, the data in the currently authenticated token is
  unwrapped.
---

# unwrap

The `unwrap` command unwraps a wrapped secret from Vault by the given token. The
result is the same as the "vault read" operation on the non-wrapped secret. If
no token is given, the data in the currently authenticated token is unwrapped.

## Examples

Unwrap the data in the cubbyhole secrets engine for a token:

```text
$ vault unwrap 3de9ece1-b347-e143-29b0-dc2dc31caafd
```

Unwrap the data in the active token:

```text
$ vault login 848f9ccf-7176-098c-5e2b-75a0689d41cd
$ vault unwrap # unwraps 848f9ccf...
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-field` `(string: "")` - Print only the field with the given name. Specifying
  this option will take precedence over other formatting directives. The result
  will not have a trailing newline making it ideal for piping to other processes.

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
