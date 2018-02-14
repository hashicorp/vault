---
layout: "docs"
page_title: "Socket - Audit Devices"
sidebar_current: "docs-audit-socket"
description: |-
  The "socket" audit device writes audit writes to a TCP or UDP socket.
---

# Socket Audit Device

The `socket` audit device writes to a TCP, UDP, or UNIX socket.

~> **Warning:** Due to the nature of the underlying protocols used in this
device there exists a case when the connection to a socket is lost a single
audit entry could be omitted from the logs and the request will still succeed.
Using this device in conjunction with another audit device will help to improve
accuracy, but the socket device should not be used if strong guarantees are
needed for audit logs.

## Enabling

Enable at the default path:

```text
$ vault audit enable socket
```

Supply configuration parameters via K=V pairs:

```text
$ vault audit enable socket address=127.0.0.1:9090 socket_type=tcp
```

## Configuration

- `address` `(string: "")` - The socket server address to use. Example
  `127.0.0.1:9090` or `/tmp/audit.sock`.

- `socket_type` `(string: "tcp")` - The socket type to use, any type compatible
  with <a href="https://golang.org/pkg/net/#Dial">net.Dial</a> is acceptable.

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
