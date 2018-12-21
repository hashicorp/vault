---
layout: "docs"
page_title: "AWS KMS - Seals - Configuration"
sidebar_title: "AWS KMS"
sidebar_current: "docs-configuration-seal-awskms"
description: |-
  The AWS KMS seal configures Vault to use AWS KMS as the seal wrapping
  mechanism.
---

# `awskms` Seal

The AWS KMS seal configures Vault to use AWS KMS as the seal wrapping mechanism.
The AWS KMS seal is activated by one of the following:

* The presence of a `seal "awskms"` block in Vault's configuration file
* The presence of the environment variable `VAULT_SEAL_TYPE` set to `awskms`. If
  enabling via environment variable, all other required values specific to AWS
  KMS (i.e. `VAULT_AWSKMS_SEAL_KEY_ID`) must be also supplied, as well as all
  other AWS-related environment variables that lends to successful
  authentication (i.e. `AWS_ACCESS_KEY_ID`, etc.).

## `awskms` Example

This example shows configuring AWS KMS seal through the Vault configuration file
by providing all the required values:

```hcl
seal "awskms" {
  region     = "us-east-1"
  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  kms_key_id = "19ec80b0-dfdd-4d97-8164-c6examplekey"
  endpoint   = "https://vpce-0e1bb1852241f8cc6-pzi0do8n.kms.us-east-1.vpce.amazonaws.com"
}
```

## `awskms` Parameters

These parameters apply to the `seal` stanza in the Vault configuration file:

- `region` `(string: "us-east-1")`: The AWS region where the encryption key
  lives. May also be specified by the `AWS_REGION` or `AWS_DEFAULT_REGION`
  environment variable or as part of the AWS profile from the AWS CLI or
  instance profile.

- `access_key` `(string: <required>)`: The AWS access key ID to use. May also be
  specified by the `AWS_ACCESS_KEY_ID` environment variable or as part of the
  AWS profile from the AWS CLI or instance profile.

- `secret_key` `(string: <required>)`: The AWS secret access key to use. May
  also be specified by the `AWS_SECRET_ACCESS_KEY` environment variable or as
  part of the AWS profile from the AWS CLI or instance profile.

- `kms_key_id` `(string: <required>)`: The AWS KMS key ID to use for encryption
  and decryption. May also be specified by the `VAULT_AWSKMS_SEAL_KEY_ID`
  environment variable.

- `endpoint` `(string: "")`: The KMS API endpoint to be used to make AWS KMS
  requests. May also be specified by the `AWS_KMS_ENDPOINT` environment
  variable. This is useful, for example, when connecting to KMS over a [VPC
  Endpoint](https://docs.aws.amazon.com/kms/latest/developerguide/kms-vpc-endpoint.html).
  If not set, Vault will use the default API endpoint for your region.

## Authentication

Authentication-related values must be provided, either as environment
variables or as configuration parameters.

~> **Note:** Although the configuration file allows you to pass in
`AWS_ACCESS_KEY_ID` and `AWS_ACCESS_KEY_ID` as part of the seal's parameters, it
is *strongly* recommended to set these values via environment variables.

```text
AWS authentication values:

* `AWS_REGION` or `AWS_DEFAULT_REGION`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
```

Note: The client uses the official AWS SDK and will use the specified
credentials, environment credentials, shared file credentials, or IAM role/ECS
task credentials in that order, if the above AWS specific values are not
provided.

Vault needs the following permissions on the KMS key:

* `kms:Encrypt`
* `kms:Decrypt`
* `kms:DescribeKey`

These can be granted via IAM permissions on the principal that Vault uses, on
the KMS key policy for the KMS key, or via KMS Grants on the key.

## `awskms` Environment Variables

Alternatively, the AWS KMS seal can be activated by providing the following
environment variables:

```text
Vault Seal specific values:

* `VAULT_SEAL_TYPE`
* `VAULT_AWSKMS_SEAL_KEY_ID`
```

## Key Rotation

This seal supports rotating the master keys defined in AWS KMS 
[doc](https://docs.aws.amazon.com/kms/latest/developerguide/rotate-keys.html). Both automatic 
rotation and manual rotation is supported for KMS since the key information is stored with the 
encrypted data.  Old keys must not be disabled or deleted and are used to decrypt older data. 
Any new or updated data will be encrypted with the current key defined in the seal configuration
or set to current under a key alias.
