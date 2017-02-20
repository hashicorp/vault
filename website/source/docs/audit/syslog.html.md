---
layout: "docs"
page_title: "Audit Backend: Syslog"
sidebar_current: "docs-audit-syslog"
description: |-
  The "syslog" audit backend writes audit logs to syslog.
---

# Audit Backend: Syslog

The `syslog` audit backend writes audit logs to syslog.

It currently does not support a configurable syslog destination, and always
sends to the local agent. This backend is only supported on Unix systems,
and should not be enabled if any standby Vault instances do not support it.

## Format

Each line in the audit log is a JSON object. The `type` field specifies what type of
object it is. Currently, only two types exist: `request` and `response`. The line contains
all of the information for any given request and response. By default, all the sensitive
information is first hashed before logging in the audit logs.

## Enabling

#### Via the CLI

Audit `syslog` backend can be enabled by the following command.

```
$ vault audit-enable syslog
```

Backend configuration options can also be provided from command-line.

```
$ vault audit-enable syslog tag="vault" facility="AUTH"
```

Following are the configuration options available for the backend.

<dl class="api">
  <dt>Backend configuration options</dt>
  <dd>
    <ul>
      <li>
        <span class="param">facility</span>
        <span class="param-flags">optional</span>
            The syslog facility to use. Defaults to `AUTH`.
      </li>
      <li>
        <span class="param">tag</span>
        <span class="param-flags">optional</span>
            The syslog tag to use. Defaults to `vault`.
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
        <span class="param">prefix</span>
        <span class="param-flags">optional</span>
            Allows a customizable string prefix to write before the actual log
            line. Defaults to an empty string.
      </li>
    </ul>
  </dd>
</dl>
