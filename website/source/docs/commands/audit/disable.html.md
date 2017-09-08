---
layout: "docs"
page_title: "audit disable - Command"
sidebar_current: "docs-commands-audit-disable"
description: |-
  The "audit disable" command disables an audit device at a given path, if one
  exists. This command is idempotent, meaning it succeeds even if no audit
  device is enabled at the path.
---

# audit disable

The `audit disable` command disables an audit device at a given path, if one
exists. This command is idempotent, meaning it succeeds even if no audit device
is enabled at the path.

Once an audit device is disabled, no future audit logs are dispatched to it. The
data associated with the audit device is unaffected. For example, if you
disabled an audit device that was logging to a file, the file would still exist
and have stored contents.

## Examples

Disable the audit device enabled at "file/":

```text
$ vault audit disable file/
Success! Disabled audit device (if it was enabled) at: file/
```

## Usage

There are no flags beyond the [standard set of flags](/docs/commands/index.html)
included on all commands.
