---
layout: "docs"
page_title: "Kafka - Audit Devices"
sidebar_title: "Kafka"
sidebar_current: "docs-audit-kafka"
description: |-
  The "kafka" audit device writes audit logs to Apache Kafka.
---

# Kafka Audit Device

The `kafka` audit device writes audit logs to a Apache Kafka topic.
Messages will be written in a json format, with the message key being the
Request ID of the audit log.

## Examples

Enable at the default path:

```text
$ vault audit enable kafka topic=vault-audit-logs address=kafka.service.consul:9092
```

## Configuration

- `topic` `(string: "vault")` -  The name of the Kafka topic

- `address` `(string: "kafka.service.consul:9092)"` - The address of a Kafka
  broker.

- `tls_disabled` (bool: "false") - Disable TLS

- `ca_cert` (string: "-----BEGIN CERTIFICATE----- ...") - The root certificate to trust

- `client_cert` (string:  "-----BEGIN CERTIFICATE----- ...") The producers
  certificate.

- `client_private_key` (string:  "-----BEGIN RSA PRIVATE KEY----- ...") - The
  producers private key.

- `hmac_accessor` `(bool: true)` - If enabled, enables the hashing of token

- `log_raw` `(bool: false)` - If enabled, logs the security sensitive
  information without hashing, in the raw format.
  accessor.

