---
layout: "docs"
page_title: "plugin register - Command"
sidebar_title: "<code>register</code>"
sidebar_current: "docs-commands-plugin-register"
description: |-
  The "plugin register" command registers a new plugin in Vault's plugin
  catalog.
---

# plugin register

The `plugin register` command registers a new plugin in Vault's plugin catalog.
The plugin's type of "auth", "database", or "secret" must be included.

## Examples

Register a plugin:

```text
$ vault plugin register \
    -sha256=d3f0a8be02f6c074cf38c9c99d4d04c9c6466249 \
    auth my-custom-plugin
Success! Registered plugin: my-custom-plugin
```

Register a plugin with custom args:

```text
$ vault plugin register \
    -sha256=d3f0a8be02f6c074cf38c9c99d4d04c9c6466249 \
    -args=--with-glibc,--with-curl-bindings \
    auth my-custom-plugin
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-sha256` `(string: <required>)` - Checksum (SHA256) of the plugin binary.

- `-args` `(string: "")` - List of arguments to pass to the binary plugin during
  each invocation. Specify multiple arguments with commas.

- `-command` `(string: "")` - Name of the command to run to invoke the binary.
  By default, this is the name of the plugin.
