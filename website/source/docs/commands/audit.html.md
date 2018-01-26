---
layout: "docs"
page_title: "audit - Command"
sidebar_current: "docs-commands-audit"
description: |-
  The "audit" command groups subcommands for interacting with Vault's audit
  devices. Users can list, enable, and disable audit devices.
---

# audit

The `audit` command groups subcommands for interacting with Vault's audit
devices. Users can list, enable, and disable audit devices.

For more information, please see the [audit device
documentation](/docs/audit/index.html)

## Examples

Enable an audit device:

```text
$ vault audit enable file file_path=/tmp/my-file.txt
Success! Enabled the file audit device at: file/
```

List all audit devices:

```text
$ vault audit list
Path     Type    Description
----     ----    -----------
file/    file    n/a
```

Disable an audit device:

```text
$ vault audit disable file/
Success! Disabled audit device (if it was enabled) at: file/
```

## Usage

```text
Usage: vault audit <subcommand> [options] [args]

  # ...

Subcommands:
    disable    Disables an audit device
    enable     Enables an audit device
    list       Lists enabled audit devices
```

For more information, examples, and usage about a subcommand, click on the name
of the subcommand in the sidebar.
