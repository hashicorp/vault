---
layout: "docs"
page_title: "plugin list - Command"
sidebar_current: "docs-commands-plugin-list"
description: |-
  The "plugin list" command lists all available plugins in the plugin catalog.
---

# plugin list

The `plugin list` command lists all available plugins in the plugin catalog.

## Examples

List all available plugins in the catalog.

```text
$ vault plugin list

Plugins
-------
my-custom-plugin
# ...
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.
