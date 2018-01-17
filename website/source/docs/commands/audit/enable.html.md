---
layout: "docs"
page_title: "audit enable - Command"
sidebar_current: "docs-commands-audit-enable"
description: |-
  The "audit enable" command enables an audit device at a given path.
---

# audit enable

The `audit enable` command enables an audit device at a given path. If an audit
device already exists at the given path, an error is returned. Additional
options for configuring the audit device are provided as `KEY=VALUE`. Each audit
device declares its own set of configuration options.

Once an audit device is enabled, almost every request and response will be
logged to the device.

## Examples

Enable the audit device "file" enabled at "file/":

```text
$ vault audit enable file file_path=/tmp/my-file.txt
Success! Enabled the file audit device at: file/
```

## Usage

The following flags are available in addition to the [standard set of
flags](/docs/commands/index.html) included on all commands.

- `-description` `(string: "")` - Human-friendly description for the purpose of
  this audit device.

-  `-local` `(bool: false)` - Mark the audit device as a local-only device.
   Local devices are not replicated or removed by replication.

- `-path` `(string: "")` - Place where the audit device will be accessible. This
  must be unique across all audit devices. This defaults to the "type" of the
  audit device.
