---
layout: "docs"
page_title: "File - Audit Devices"
sidebar_current: "docs-audit-file"
description: |-
  The "file" audit device writes audit logs to a file.
---

# File Audit Device

The `file` audit device writes audit logs to a file. This is a very simple audit
device: it appends logs to a file.

The device does not currently assist with any log rotation. There are very
stable and feature-filled log rotation tools already, so we recommend using
existing tools.

Sending a `SIGHUP` to the Vault process will cause `file` audit devices to close
and re-open their underlying file, which can assist with log rotation needs.

## Examples

Enable at the default path:

```text
$ vault audit enable file file_path=/var/log/vault_audit.log
```

Enable at a different path. It is possible to enable multiple copies of an audit
device:

```text
$ vault audit enable -path="vault_audit_1" file file_path=/home/user/vault_audit.log
```

## Configuration

- `file_path` `(string: "")` - The path to where the audit log will be written.
  If this path exists, the audit device will append to it. Specify `"stdout"` to
  write audit log to standard output. Specify `"discard"` to discard output
  (useful in testing scenarios).

- `log_raw` `(bool: false)` - If enabled, logs the security sensitive
  information without hashing, in the raw format.

- `hmac_accessor` `(bool: true)` - If enabled, enables the hashing of token
  accessor.

- `mode` `(string: "0600")` - A string containing an octal number representing
  the bit pattern for the file mode, similar to `chmod`.

- `format` `(string: "json")` - Allows selecting the output format. Valid values
  are `"json"` and `"jsonx"`, which formats the normal log entries as XML.

- `prefix` `(string: "")` - A customizable string prefix to write before the
  actual log line.
