---
layout: "docs"
page_title: "Audit Backend: File"
sidebar_current: "docs-audit-file"
description: |-
  The "file" audit backend writes audit logs to a file.
---

# Audit Backend: File

Name: `file`

The "file" audit backend writes audit logs to a file.

This is a very simple audit backend: it appends logs to a file. It does
not currently assist with any log rotation.

## Options

When enabling this backend, the following options are accepted:

  * `path` (required) - The path to where the file will be written. If
      this path exists, the audit backend will append to it.
  * `log_raw` (optional) Should security sensitive information be logged raw. Defaults to "false".

## Format

Each line in the audit log is a JSON object. The "type" field specifies
what type of object it is. Currently, only two types exist: "request" and
"response".

The line contains all of the information for any given request and response.

If `log_raw` if false, as is default, all sensitive information is first hashed
before logging. If explicitly enabled, all values are logged raw without hashing.

