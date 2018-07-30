---
layout: "docs"
page_title: "Vault Agent Auto-Auth AWS Method"
sidebar_current: "docs-agent-autoauth-methods-aws"
description: |-
  AWS Method for Vault Agent Auto-Auth
---

# Vault Agent Auto-Auth AWS Method 

The `aws` method performs authentication against the [AWS Auth
method](https://www.vaultproject.io/docs/auth/aws.html). Both `ec2` and `iam`
authentication types are supported. If `ec2` is used, the agent will store the
reauthentication value in memory and use it for reauthenticating, but will not
persist it to disk.

Due to the complexity of the Trust On First Use (TOFU) model used in the `ec2`
method, we recommend the `iam` method when possible.

## Credentials

Vault will use the AWS SDK's normal credential chain behavior, which means it
will try to source credentials from the assigned instance profile, a
credentials file, the environment, or static credentials. Generally it should
not be required to set the `access_key` and `secret_key` parameters.

## Configuration

- `type` `(string: required)` - The type of authentication; must be `ec2` or `iam`

- `role` `(string: required)` - The role to authenticate against on Vault

- `access_key` `(string: optional)` - When using static credentials, the access key to use

- `secret_key` `(string: optional)` - When using static credentials, the secret key to use

- `session_token` `(string: optional)` - The session token to use for authentication, if needed

- `header_value` `(string: optional)` - If configured in Vault, the value to
  use for
  [`iam_server_id_header_value`](https://www.vaultproject.io/api/auth/aws/index.html#iam_server_id_header_value)
