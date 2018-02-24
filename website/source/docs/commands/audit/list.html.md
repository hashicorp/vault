---
layout: "docs"
page_title: "audit list - Command"
sidebar_current: "docs-commands-audit-list"
description: |-
  The "audit list" command lists the audit devices enabled. The output lists the
  enabled audit devices and options for those devices.
---

# audit list

The `audit list` command lists the audit devices enabled. The output lists the
enabled audit devices and options for those devices.

## Examples

List all audit devices:

```text
$ vault audit list
Path     Type    Description
----     ----    -----------
file/    file    n/a
```

List detailed audit device information:

```text
$ vault audit list -detailed
Path     Type    Description    Replication    Options
----     ----    -----------    -----------    -------
file/    file    n/a            replicated     file_path=/var/log/audit.log
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

### Output Options

- `-format` `(string: "table")` - Print the output in the given format. Valid
  formats are "table", "json", or "yaml". This can also be specified via the
  `VAULT_FORMAT` environment variable.

### Command Options

- `-detailed` `(bool: false)` - Print detailed information such as options and
  replication status about each auth device.
