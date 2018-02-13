---
layout: "docs"
page_title: "write - Command"
sidebar_current: "docs-commands-write"
description: |-
  The "write" command writes data to Vault at the given path. The data can be
  credentials, secrets, configuration, or arbitrary data. The specific behavior
  of this command is determined at the thing mounted at the path.
---

# write

The `write` command writes data to Vault at the given path. The data can be
credentials, secrets, configuration, or arbitrary data. The specific behavior of
this command is determined at the thing mounted at the path.

Data is specified as "key=value" pairs. If the value begins with an "@", then it
is loaded from a file. If the value is "-", Vault will read the value from
stdin.

For a full list of examples and paths, please see the documentation that
corresponds to the secrets engines in use.

## Examples

Persist data in the KV secrets engine:

```text
$ vault write secret/my-secret foo=bar
```

Create a new encryption key in the transit secrets engine:

```text
$ vault write -f transit/keys/my-key
```

Upload an AWS IAM policy from a file on disk:

```text
$ vault write aws/roles/ops policy=@policy.json
```

Configure access to Consul by providing an access token:

```text
$ echo $MY_TOKEN | vault write consul/config/access token=-
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

### Command Options

- `-force` `(bool: false)` - Allow the operation to continue with no key=value
  pairs. This allows writing to keys that do not need or expect data. This is
  aliased as "-f".
