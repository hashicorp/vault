---
layout: "docs"
page_title: "Syslog - Audit Devices"
sidebar_current: "docs-audit-syslog"
description: |-
  The "syslog" audit device writes audit logs to syslog.
---

# Syslog Audit Device

The `syslog` audit device writes audit logs to syslog.

It currently does not support a configurable syslog destination, and always
sends to the local agent. This device is only supported on Unix systems,
and should not be enabled if any standby Vault instances do not support it.

## Examples

Audit `syslog` device can be enabled by the following command:

```text
$ vault audit enable syslog
```

Supply configuration parameters via K=V pairs:

```text
$ vault audit enable syslog tag="vault" facility="AUTH"
```

## Configuration

- `facility` `(string: "AUTH")` - The syslog facility to use.

- `tag` `(string: "vault")` - The syslog tag to use.

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
