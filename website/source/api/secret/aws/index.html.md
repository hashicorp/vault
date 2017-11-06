---
layout: "api"
page_title: "AWS Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-aws"
description: |-
  This is the API documentation for the Vault AWS secret backend.
---

# AWS Secret Backend HTTP API

This is the API documentation for the Vault AWS secret backend. For general
information about the usage and operation of the AWS backend, please see the
[Vault AWS backend documentation](/docs/secrets/aws/index.html).

This documentation assumes the AWS backend is mounted at the `/aws` path in
Vault. Since it is possible to mount secret backends at any location, please
update your API calls accordingly.

## Configure Root IAM Credentials

This endpoint configures the root IAM credentials to communicate with AWS. There
are multiple ways to pass root IAM credentials to the Vault server, specified
below with the highest precedence first. If credentials already exist, this will
overwrite them.

The official AWS SDK is used for sourcing credentials from env vars, shared
files, or IAM/ECS instances.

- Static credentials provided to the API as a payload

- Credentials in the `AWS_ACCESS_KEY`, `AWS_SECRET_KEY`, and `AWS_REGION`
  environment variables **on the server**

- Shared credentials files

- Assigned IAM role or ECS task role credentials

At present, this endpoint does not confirm that the provided AWS credentials are
valid AWS credentials with proper permissions.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/aws/config/root`           | `204 (empty body)`     |

### Parameters

- `access_key` `(string: <required>)` – Specifies the AWS access key ID.

- `secret_key` `(string: <required>)` – Specifies the AWS secret access key.

- `region` `(string: <optional>)` – Specifies the AWS region. If not set it
  will use the `AWS_REGION` env var, `AWS_DEFAULT_REGION` env var, or
  `us-east-1` in that order.

- `iam_endpoint` `(string: <optional>)` – Specifies a custom HTTP IAM endpoint to use.

- `sts_endpoint` `(string: <optional>)` – Specifies a custom HTTP STS endpoint to use.

### Sample Payload

```json
{
  "access_key": "AKIA...",
  "secret_key": "2J+...",
  "region": "us-east-1"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/aws/config/root
```

## Configure Lease

This endpoint configures lease settings for the AWS secret backend. It is
optional, as there are default values for `lease` and `lease_max`.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/aws/config/lease`          | `204 (empty body)`     |

### Parameters

- `lease` `(string: <required>)` – Specifies the lease value provided as a
  string duration with time suffix. "h" (hour) is the largest suffix.

- `lease_max` `(string: <required>)` – Specifies the maximum lease value
  provided as a string duration with time suffix. "h" (hour) is the largest
  suffix.

### Sample Payload

```json
{
  "lease": "30m",
  "lease_max": "12h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/aws/config/lease
```

## Read Lease

This endpoint returns the current lease settings for the AWS secret backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/aws/config/lease`          | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/aws/config/lease
```

### Sample Response

```json
{
  "data": {
    "lease": "30m0s",
    "lease_max": "12h0m0s"
  }
}
```

## Create/Update Role

This endpoint creates or updates the role with the given `name`. If a role with
the name does not exist, it will be created. If the role exists, it will be
updated with the new attributes.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/aws/roles/:name`           | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to create. This
  is part of the request URL.

- `policy` `(string: <required unless arn provided>)` – Specifies the IAM policy
  in JSON format.

- `arn` `(string: <required unless policy provided>)` – Specifies the full ARN
  reference to the desired existing policy.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/aws/roles/example-role
```

### Sample Payloads

Using an inline IAM policy:

```json
{
  "policy": "{\"Version\": \"...\"}",
}
```

Using an ARN:

```json
{
  "arn": "arn:aws:iam::123456789012:user/David"
}
```

## Read Role

This endpoint queries an existing role by the given name. If the role does not
exist, a 404 is returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/aws/roles/:name`           | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to read. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/aws/roles/example-role
```

### Sample Responses

For an inline IAM policy:

```json
{
  "data": {
    "policy": "{\"Version\": \"...\"}"
  }
}
```

For an ARN:

```json
{
  "data": {
    "arn": "arn:aws:iam::123456789012:user/David"
  }
}
```

## List Roles

This endpoint lists all existing roles in the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/aws/roles`                 | `200 application/json` |
| `GET`    | `/aws/roles?list=true`       | `200 application/json` |

### Sample Request

```
$ curl
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/aws/roles
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "example-role"
    ]
  }
}
```

## Delete Role

This endpoint deletes an existing role by the given name. If the role does not
exist, a 404 is returned.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE`  | `/aws/roles/:name`           | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to delete. This
  is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/aws/roles/example-role
```

## Generate IAM Credentials

This endpoint generates dynamic IAM credentials based on the named role. This
role must be created before queried.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/aws/creds/:name`           | `200 application/json` |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role to generate
  credentials against. This is part of the request URL.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/aws/creds/example-role
```

### Sample Response

```json
{
  "data": {
    "access_key": "AKIA...",
    "secret_key": "xlCs...",
    "security_token": null
  }
}
```

## Generate IAM with STS

This generates a dynamic IAM credential with an STS token based on the named
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/aws/sts/:name`             | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the role against which
  to create this STS credential. This is part of the request URL.

- `ttl` `(string: "3600s")` – Specifies the TTL for the use of the STS token.
  This is specified as a string with a duration suffix. AWS documentation
  excerpt: `The duration, in seconds, that the credentials should remain valid.
  Acceptable durations for IAM user sessions range from 900 seconds (15
  minutes) to 129600 seconds (36 hours), with 43200 seconds (12 hours) as the
  default. Sessions for AWS account owners are restricted to a maximum of 3600
  seconds (one hour). If the duration is longer than one hour, the session for
  AWS account owners defaults to one hour.`

### Sample Payload

```json
{
  "ttl": "15m"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/aws/sts/example-role
```

### Sample Response

```json
{
  "data": {
    "access_key": "AKIA...",
    "secret_key": "xlCs...",
    "security_token": "429255"
  }
}
```
