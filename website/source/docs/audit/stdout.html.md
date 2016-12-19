---
layout: "docs"
page_title: "Audit Backend: Stdout"
sidebar_current: "docs-audit-stdout"
description: |-
  The "stdout" audit backend writes audit logs to stdout.
---

# Audit Backend: Stdout

The `stdout` audit backend writes audit logs stdout. This is useful if you have
a requirement to capture logs from stdout.

## Format

Each line in the audit log is a JSON object. The `type` field specifies what type of
object it is. Currently, only two types exist: `request` and `response`. The line contains
all of the information for any given request and response. By default, all the sensitive
information is first hashed before logging in the audit logs.

## Enabling

#### Via the CLI

Audit `stdout` backend can be enabled by the following command.

```
$ vault audit-enable stdout
```

Use `vault audit-enable -help` to see the command options.
Following are the configuration options available for the backend.

<dl class="api">
  <dt>Backend configuration options</dt>
  <dd>
    <ul>
      <li>
        <span class="param">log_raw</span>
        <span class="param-flags">optional</span>
            A string containing a boolean value ('true'/'false'), if set, logs
            the security sensitive information without hashing, in the raw
            format. Defaults to `false`.
      </li>
      <li>
        <span class="param">hmac_accessor</span>
        <span class="param-flags">optional</span>
            A string containing a boolean value ('true'/'false'), if set,
            enables the hashing of token accessor. Defaults
            to `true`. This option is useful only when `log_raw` is `false`.
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
            Allows selecting the output format. Valid values are `json` (the
            default) and `jsonx`, which formats the normal log entries as XML.
      </li>
    </ul>
  </dd>
</dl>

