---
layout: "docs"
page_title: "Syslog - Audit Devices"
sidebar_title: "Syslog"
sidebar_current: "docs-audit-syslog"
description: |-
  The "syslog" audit device writes audit logs to syslog.
---

# Syslog Audit Device

The `syslog` audit device writes audit logs to syslog.

It currently does not support a configurable syslog destination, and always
sends to the local agent. This device is only supported on Unix systems,
and should not be enabled if any standby Vault instances do not support it.

~> **Warning**: Audit messages generated for some operations can be quite
large, and can be larger than a [maximum-size single UDP
packet](https://tools.ietf.org/html/rfc5426#section-3.1). If possible with your
syslog daemon, configure a TCP listener. Otherwise, consider using a `file`
backend and having syslog configured to read entries from the file; or, enable
both `file` and `syslog` so that a failure for a particular message to log
directly to `syslog` will not result in Vault being blocked.

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
