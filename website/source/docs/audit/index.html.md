---
layout: "docs"
page_title: "Audit Devices"
sidebar_current: "docs-audit"
description: |-
  Audit devices are mountable devices that log requests and responses in Vault.
---

# Audit Devices

Audit devices are the components in Vault that keep a detailed log
of all requests and response to Vault. Because _every_ operation with
Vault is an API request/response, the audit log contains _every_ interaction
with Vault, including errors.

Multiple audit devices can be enabled and Vault will send the audit logs to
both. This allows you to not only have a redundant copy, but also a second copy
in case the first is tampered with.

## Format

Each line in the audit log is a JSON object. The `type` field specifies what
type of object it is. Currently, only two types exist: `request` and `response`.
The line contains all of the information for any given request and response. By
default, all the sensitive information is first hashed before logging in the
audit logs.

## Sensitive Information

The audit logs contain the full request and response objects for every
interaction with Vault. The request and response can be matched utilizing a
unique identifier assigned to each request. The data in the request and the
data in the response (including secrets and authentication tokens) will be
hashed with a salt using HMAC-SHA256.

The purpose of the hash is so that secrets aren't in plaintext within your
audit logs. However, you're still able to check the value of secrets by
generating HMACs yourself; this can be done with the audit device's hash
function and salt by using the `/sys/audit-hash` API endpoint (see the
documentation for more details).

## Enabling/Disabling Audit Devices

When a Vault server is first initialized, no auditing is enabled. Audit
devices must be enabled by a root user using `vault audit enable`.

When enabling an audit device, options can be passed to it to configure it.
For example, the command below enables the file audit device:

```text
$ vault audit enable file file_path=/var/log/vault_audit.log
```

In the command above, we passed the "file_path" parameter to specify the path
where the audit log will be written to. Each audit device has its own
set of parameters. See the documentation to the left for more details.

When an audit device is disabled, it will stop receiving logs immediately.
The existing logs that it did store are untouched.

## Blocked Audit Devices

If there are any audit devices enabled, Vault requires that at least
one be able to persist the log before completing a Vault request.

!> If you have only one audit device enabled, and it is blocking (network
block, etc.), then Vault will be _unresponsive_. Vault **will not** complete
any requests until the audit device can write.

If you have more than one audit device, then Vault will complete the request
as long as one audit device persists the log.

Vault will not respond to requests if audit devices are blocked because
audit logs are critically important and ignoring blocked requests opens
an avenue for attack. Be absolutely certain that your audit devices cannot
block.

## API

Audit devices also have a full HTTP API. Please see the [Audit device API
docs](/api/system/audit.html) for more details.
