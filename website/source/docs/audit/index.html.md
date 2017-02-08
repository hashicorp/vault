---
layout: "docs"
page_title: "Audit Backends"
sidebar_current: "docs-audit"
description: |-
  Audit backends are mountable backends that log requests and responses in Vault.
---

# Audit Backends

Audit backends are the components in Vault that keep a detailed log
of all requests and response to Vault. Because _every_ operation with
Vault is an API request/response, the audit log contains _every_ interaction
with Vault, including errors.

Vault ships with multiple audit backends, depending on the location you want
the logs sent to. Multiple audit backends can be enabled and Vault will send
the audit logs to both. This allows you to not only have a redundant copy,
but also a second copy in case the first is tampered with.

## Sensitive Information

The audit logs contain the full request and response objects for every
interaction with Vault. The request and response can be matched utilizing a
unique identifier assigned to each request. The data in the request and the
data in the response (including secrets and authentication tokens) will be
hashed with a salt using HMAC-SHA256.

The purpose of the hash is so that secrets aren't in plaintext within your
audit logs. However, you're still able to check the value of secrets by
generating HMACs yourself; this can be done with the audit backend's hash
function and salt by using the `/sys/audit-hash` API endpoint (see the
documentation for more details).

## Enabling/Disabling Audit Backends

When a Vault server is first initialized, no auditing is enabled. Audit
backends must be enabled by a root user using `vault audit-enable`.

When enabling an audit backend, options can be passed to it to configure it.
For example, the command below enables the file audit backend:

```
$ vault audit-enable file file_path=/var/log/vault_audit.log
...
```

In the command above, we passed the "file_path" parameter to specify the path
where the audit log will be written to. Each audit backend has its own
set of parameters. See the documentation to the left for more details.

When an audit backend is disabled, it will stop receiving logs immediately.
The existing logs that it did store are untouched.

## Blocked Audit Backends

If there are any audit backends enabled, Vault requires that at least
one be able to persist the log before completing a Vault request.

If you have only one audit backend enabled, and it is blocking (network
block, etc.), then Vault will be _unresponsive_. Vault _will not_ complete
any requests until the audit backend can write.

If you have more than one audit backend, then Vault will complete the request
as long as one audit backend persists the log.

Vault will not respond to requests if audit backends are blocked because
audit logs are critically important and ignoring blocked requests opens
an avenue for attack. Be absolutely certain that your audit backends cannot
block.

## API

### /sys/audit/[path]
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
      Enables audit backend at the specified path.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dd>
    <ul>
      <li>
       <span class="param">type</span>
        <span class="param-flags">required</span>
            The path to where the audit log will be written. If this
            path exists, the audit backend will append to it.
      </li>
      <li>
        <span class="param">description</span>
        <span class="param-flags">optional</span>
            A description.
      </li>
      <li>
        <span class="param">options</span>
        <span class="param-flags">optional</span>
           Configuration options of the backend in JSON format.
           Refer to `syslog`, `file` and `socket` audit backend options.
      </li>
    </ul>
  </dd>
</dl>
