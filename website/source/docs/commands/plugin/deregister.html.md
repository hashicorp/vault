---
layout: "docs"
page_title: "plugin deregister - Command"
sidebar_title: "<code>deregister</code>"
sidebar_current: "docs-commands-plugin-deregister"
description: |-
  The "plugin deregister" command deregisters a new plugin in Vault's plugin
  catalog.
---

# plugin deregister

The `plugin deregister` command deregisters an existing plugin from Vault's
plugin catalog. If the plugin does not exist, no error is returned. The plugin's
type of "auth", "database", or "secret" must be included.

## Examples

Deregister a plugin:

```text
$ vault plugin deregister auth my-custom-plugin
Success! Deregistered plugin (if it was registered): my-custom-plugin
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
