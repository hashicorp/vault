---
layout: "docs"
page_title: "plugin info - Command"
sidebar_title: "<code>info</code>"
sidebar_current: "docs-commands-plugin-info"
description: |-
  The "plugin info" command displays information about a plugin in the catalog.
---

# plugin info

The `plugin info` displays information about a plugin in the catalog.
The plugin's type of "auth", "database", or "secret" must be included.

## Examples

Display information about a plugin

```text
$ vault plugin info auth my-custom-plugin

Key        Value
---        -----
args       []
builtin    false
command    my-custom-plugin
name       my-custom-plugin
sha256     d3f0a8be02f6c074cf38c9c99d4d04c9c6466249
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-field` `(string: "")` - Print only the field with the given name. Specifying
  this option will take precedence over other formatting directives. The result
  will not have a trailing newline making it ideal for piping to other
  processes.

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
