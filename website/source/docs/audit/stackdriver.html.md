---
layout: "docs"
page_title: "Stackdriver - Audit Devices"
sidebar_title: "Stackdriver"
sidebar_current: "docs-audit-stackdriver"
description: |-
  The "stackdriver" audit device writes audit logs to Stackdriver Logging.
---

# Stackdriver Audit Device

The Stackdriver audit device uses the official Google Cloud Golang SDK. This means
it supports the common ways of
[providing credentials to Google Cloud](https://cloud.google.com/docs/authentication/production#providing_credentials_to_your_application).

To use this audit device, the service account must have the following
minimum scope(s):

```text
https://www.googleapis.com/auth/cloud-platform
https://www.googleapis.com/auth/logging.write
```

And the following IAM role(s) over the `parent` where logs will be written:

```text
roles/monitoring.metricWriter
```

## Enabling

Enable at the default path:

```text
$ vault audit enable stackdriver parent=foo log_id=bar
```

## Configuration

- `parent` `(string: "")` - The location to write logs to. Example
  `projects/PROJECT_ID` or `folders/FOLDER_ID`.

- `log_id` `(string: "")` - The name of the log to write to. Example
  `vault` or `vault-prod`.

- `async` `(bool: false)` - If enabled, logs are written to Stackdriver
  in batches asynchronously. This is more performant, but less secure,
  since logs may be lost in the event that the vault process terminates.
  This should not be enabled if strong guarantees are needed.

- `log_raw` `(bool: false)` - If enabled, logs the security sensitive
  information without hashing, in the raw format.

- `hmac_accessor` `(bool: true)` - If enabled, enables the hashing of token
  accessor.

