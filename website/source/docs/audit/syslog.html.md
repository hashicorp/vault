---
layout: "docs"
page_title: "Audit Backend: Syslog"
sidebar_current: "docs-audit-syslog"
description: |-
  The "syslog" audit backend writes audit logs to syslog.
---

# Audit Backend: Syslog

Name: `syslog`

The "syslog" audit backend writes audit logs to syslog.

It currently does not support a configurable syslog desination, and
always sends to the local agent. This backend is only supported on Unix systems,
and should not be enabled if any standby Vault instances do not support it.

## Options

When enabling this backend, the following options are accepted:

 * `facility` (optional) - The syslog facility to use. Defaults to "AUTH".
 * `tag` (optional) - The syslog tag to use. Defaults to "vault".
 * `log_raw` (optional) Should security sensitive information be logged raw. Defaults to "false".

## Format

Each line in the audit log is a JSON object. The "type" field specifies
what type of object it is. Currently, only two types exist: "request" and
"response".

The line contains all of the information for any given request and response.

