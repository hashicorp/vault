---
layout: "docs"
page_title: "Audit Backend: File"
sidebar_current: "docs-audit-file"
description: |-
  The "file" audit backend writes audit logs to a file.
---

# Audit Backend: File

The `file` audit backend writes audit logs to a file. This is a very simple audit
backend: it appends logs to a file.

## Rotation

The backend does not currently assist with any log rotation. There are very
stable and feature-filled log rotation tools already, so we recommend using
existing tools.

As of 0.6.2, sending a `SIGHUP` to the Vault process will cause `file` audit
backends to close and re-open their underlying file, which can assist with log
rotation needs.

## Format

Each line in the audit log is a JSON object. The `type` field specifies what type of
object it is. Currently, only two types exist: `request` and `response`. The line contains
all of the information for any given request and response. By default, all the sensitive
information is first hashed before logging in the audit logs.

## Enabling

#### Via the CLI

Audit `file` backend can be enabled by the following command.

```
$ vault audit-enable file file_path=/var/log/vault_audit.log
```

Any number of `file` audit logs can be created by enabling it with different `path`s.

```
$ vault audit-enable -path="vault_audit_1" file file_path=/home/user/vault_audit.log
```

Note the difference between `audit-enable` command options and the `file` backend
configuration options. Use `vault audit-enable -help` to see the command options.
Following are the configuration options available for the backend.

<dl class="api">
  <dt>Backend configuration options</dt>
  <dd>
    <ul>
      <li>
        <span class="param">file_path</span>
        <span class="param-flags">required</span>
            The path to where the audit log will be written. If this
            path exists, the audit backend will append to it. Specify `"stdout"` to write audit log to standard output; specify `"discard"` to discard output (useful in testing scenarios).
      </li>
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
        <span class="param">mode</span>
        <span class="param-flags">optional</span>
            A string containing an octal number representing the bit pattern
            for the file mode, similar to `chmod`. This option defaults to
            `0600`.
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

