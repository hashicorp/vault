---
layout: "docs"
page_title: "Audit Backend: File"
sidebar_current: "docs-audit-file"
description: |-
  The "file" audit backend writes audit logs to a file.
---

# Audit Backend: File

The `file` audit backend writes audit logs to a file. This is a very simple audit
backend: it appends logs to a file. It does not currently assist with any log rotation.

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
            path exists, the audit backend will append to it.
      </li>
      <li>
        <span class="param">log_raw</span>
        <span class="param-flags">optional</span>
            A boolean, if set, logs the security sensitive information without
            hashing, in the raw format. Defaults to `false`.
      </li>
      <li>
        <span class="param">hmac_accessor</span>
        <span class="param-flags">optional</span>
            A boolean, if set, enables the hashing of token accessor. Defaults to `true`. This option
            is useful only when `log_raw` is `false`.
      </li>
    </ul>
  </dd>
</dl>

