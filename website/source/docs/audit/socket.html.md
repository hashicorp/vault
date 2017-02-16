---
layout: "docs"
page_title: "Audit Backend: Socket"
sidebar_current: "docs-audit-socket"
description: |-
  The "socket" audit backend writes audit writes to a TCP or UDP socket.
---

# Audit Backend: Socket

The `socket` audit backend writes to a TCP, UDP, or UNIX socket.

~> **Warning:** Due to the nature of the underlying protocols used in this backend there exists a case when the connection to a socket is lost a single audit entry could be omitted from the logs and the request will still succeed. Using this backend in conjunction with another audit backend will help to improve accuracy, but the socket backend should not be used if strong guarantees are needed for audit logs.

## Format

Each line in the audit log is a JSON object. The `type` field specifies what type of
object it is. Currently, only two types exist: `request` and `response`. The line contains
all of the information for any given request and response. By default, all the sensitive
information is first hashed before logging in the audit logs.

## Enabling

#### Via the CLI

Audit `socket` backend can be enabled by the following command.

```
$ vault audit-enable socket
```

Backend configuration options can also be provided from command-line.

```
$ vault audit-enable socket address="127.0.0.1:9090" socket_type="tcp"
```

Following are the configuration options available for the backend.

<dl class="api">
  <dt>Backend configuration options</dt>
  <dd>
    <ul>
      <li>
        <span class="param">address</span>
        <span class="param-flags">required</span>
            The socket server address to use. Example `127.0.0.1:9090` or `/tmp/audit.sock`.
      </li>
      <li>
        <span class="param">socket_type</span>
        <span class="param-flags">optional</span>
            The socket type to use, any type compatible with <a href="https://golang.org/pkg/net/#Dial">net.Dial</a> is acceptable. Defaults to `tcp`.
      </li>
      <li>
        <span class="param">log_raw</span>
        <span class="param-flags">optional</span>
            A string containing a boolean value ('true'/'false'), if set, logs the security sensitive information without
            hashing, in the raw format. Defaults to `false`.
      </li>
      <li>
        <span class="param">hmac_accessor</span>
        <span class="param-flags">optional</span>
            A string containing a boolean value ('true'/'false'), if set, enables the hashing of token accessor. Defaults
            to `true`. This option is useful only when `log_raw` is `false`.
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
            Allows selecting the output format. Valid values are `json` (the
            default) and `jsonx`, which formats the normal log entries as XML.
      </li>
      <li>
        <span class="param">write_timeout</span>
        <span class="param-flags">optional</span>
            Sets the timeout for writes to the socket. Defaults to "2s" (2 seconds).
        </li>
      <li>
        <span class="param">prefix</span>
        <span class="param-flags">optional</span>
            Allows a customizable string prefix to write before the actual log
            line. Defaults to an empty string.
      </li>
    </ul>
  </dd>
</dl>
