---
layout: "api"
page_title: "AWS Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-aws"
description: |-
  This is the API documentation for the Vault AWS authentication backend.
---

# AWS Auth Backend HTTP API

This is the API documentation for the Vault AWS authentication backend. For
general information about the usage and operation of the AWS backend, please
see the [Vault AWS backend documentation](/docs/auth/aws.html).

This documentation assumes the AWS backend is mounted at the `/auth/aws`
path in Vault. Since it is possible to mount auth backends at any location,
please update your API calls accordingly.

## Configure Client

Configures the credentials required to perform API calls to AWS as well as
custom endpoints to talk to AWS APIs. The instance identity document
fetched from the PKCS#7 signature will provide the EC2 instance ID. The
credentials configured using this endpoint will be used to query the status
of the instances via DescribeInstances API. If static credentials are not
provided using this endpoint, then the credentials will be retrieved from
the environment variables `AWS_ACCESS_KEY`, `AWS_SECRET_KEY` and
`AWS_REGION` respectively. If the credentials are still not found and if the
backend is configured on an EC2 instance with metadata querying
capabilities, the credentials are fetched automatically.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/config/client`    | `204 (empty body)`     |

### Parameters

- `access_key` `(string: "")` - AWS Access key with permissions to query AWS
  APIs. The permissions required depend on the specific configurations. If using
  the `iam` auth method without inferencing, then no credentials are necessary.
  If using the `ec2` auth method or using the `iam` auth method with
  inferencing, then these credentials need access to `ec2:DescribeInstances`. If
  additionally a `bound_iam_role` is specified, then these credentials also need
  access to `iam:GetInstanceProfile`. If, however, an alternate sts
  configuration is set for the target account, then the credentials must be
  permissioned to call `sts:AssumeRole` on the configured role, and that role
  must have the permissions described here.
- `secret_key` `(string: "")` - AWS Secret key with permissions to query AWS
  APIs.
- `endpoint` `(string: "")` - URL to override the default generated endpoint for
  making AWS EC2 API calls.
- `iam_endpoint` `(string: "")` - URL to override the default generated endpoint
  for making AWS IAM API calls.
- `sts_endpoint` `(string: "")` - URL to override the default generated endpoint
  for making AWS STS API calls.
- `iam_server_id_header_value` `(string: "")` - The value to require in the
  `X-Vault-AWS-IAM-Server-ID` header as part of GetCallerIdentity requests that
  are used in the iam auth method. If not set, then no value is required or
  validated. If set, clients must include an X-Vault-AWS-IAM-Server-ID header in
  the headers of login requests, and further this header must be among the
  signed headers validated by AWS. This is to protect against different types of
  replay attacks, for example a signed request sent to a dev server being resent
  to a production server. Consider setting this to the Vault server's DNS name.

### Sample Payload

```json
{
  "access_key": "VKIAJBRHKH6EVTTNXDHA",
  "secret_key": "vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/config/client
```

## Read Config

Returns the previously configured AWS access credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/config/client`     | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/config/client
```

### Sample Response

```json
{
  "data": {
    "secret_key": "vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj",
    "access_key": "VKIAJBRHKH6EVTTNXDHA",
    "endpoint": "",
    "iam_endpoint": "",
    "sts_endpoint": "",
    "iam_server_id_header_value": ""
  }
}
```

## Delete Config

Deletes the previously configured AWS access credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/aws/config/client`  | `204 (empty body)`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/config/client
```

## Create Certificate Configuration

Registers an AWS public key to be used to verify the instance identity
documents. While the PKCS#7 signature of the identity documents have DSA
digest, the identity signature will have RSA digest, and hence the public
keys for each type varies respectively. Indicate the type of the public key
using the "type" parameter.

| Method   | Path                                         | Produces               |
| :------- | :------------------------------------------- | :--------------------- |
| `POST`   | `/auth/aws/config/certificate/:cert_name`    | `204 (empty body)`     |

### Parameters

- `cert_name` `(string: <required>)` - Name of the certificate.
- `aws_public_cert` `(string: <required>)` - Base64 encoded AWS Public key required to verify
  PKCS7 signature of the EC2 instance metadata.
- `type` `(string: "pkcs7")` - Takes the value of either "pkcs7" or "identity",
  indicating the type of document which can be verified using the given
  certificate. The PKCS#7 document will have a DSA digest and the identity
  signature will have an RSA signature, and accordingly the public certificates
  to verify those also vary. Defaults to "pkcs7".

### Sample Payload

```json
{
  "aws_public_cert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM3VENDQXEwQ0NRQ1d1a2paNVY0YVp6QUpCZ2NxaGtqT09BUURNRnd4Q3pBSkJnTlZCQVlUQWxWVE1Sa3cKRndZRFZRUUlFeEJYWVhOb2FXNW5kRzl1SUZOMFlYUmxNUkF3RGdZRFZRUUhFd2RUWldGMGRHeGxNU0F3SGdZRApWUVFLRXhkQmJXRjZiMjRnVjJWaUlGTmxjblpwWTJWeklFeE1RekFlRncweE1qQXhNRFV4TWpVMk1USmFGdzB6Ck9EQXhNRFV4TWpVMk1USmFNRnd4Q3pBSkJnTlZCQVlUQWxWVE1Sa3dGd1lEVlFRSUV4QlhZWE5vYVc1bmRHOXUKSUZOMFlYUmxNUkF3RGdZRFZRUUhFd2RUWldGMGRHeGxNU0F3SGdZRFZRUUtFeGRCYldGNmIyNGdWMlZpSUZObApjblpwWTJWeklFeE1RekNDQWJjd2dnRXNCZ2NxaGtqT09BUUJNSUlCSHdLQmdRQ2prdmNTMmJiMVZRNHl0LzVlCmloNU9PNmtLL24xTHpsbHI3RDhad3RRUDhmT0VwcDVFMm5nK0Q2VWQxWjFnWWlwcjU4S2ozbnNzU05wSTZiWDMKVnlJUXpLN3dMY2xuZC9Zb3pxTk5tZ0l5WmVjTjdFZ2xLOUlUSEpMUCt4OEZ0VXB0M1FieVlYSmRtVk1lZ042UApodmlZdDVKSC9uWWw0aGgzUGExSEpkc2tnUUlWQUxWSjNFUjExK0tvNHRQNm53dkh3aDYrRVJZUkFvR0JBSTFqCmsrdGtxTVZIdUFGY3ZBR0tvY1Rnc2pKZW02LzVxb216SnVLRG1iSk51OVF4dzNyQW90WGF1OFFlK01CY0psL1UKaGh5MUtIVnBDR2w5ZnVlUTJzNklMMENhTy9idXljVTFDaVlRazQwS05IQ2NIZk5pWmJkbHgxRTlycFVwN2JuRgpsUmEydjFudE1YM2NhUlZEZGJ0UEVXbWR4U0NZc1lGRGs0bVpyT0xCQTRHRUFBS0JnRWJtZXZlNWY4TElFL0dmCk1ObVA5Q001ZW92UU9HeDVobzhXcUQrYVRlYnMrazJ0bjkyQkJQcWVacXBXUmE1UC8ranJkS21sMXF4NGxsSFcKTVhyczNJZ0liNitoVUlCK1M4ZHo4L21tTzBicHI3NlJvWlZDWFlhYjJDWmVkRnV0N3FjM1dVSDkrRVVBSDVtdwp2U2VEQ09VTVlRUjdSOUxJTll3b3VISXppcVFZTUFrR0J5cUdTTTQ0QkFNREx3QXdMQUlVV1hCbGs0MHhUd1N3CjdIWDMyTXhYWXJ1c2U5QUNGQk5HbWRYMlpCclZOR3JOOU4yZjZST2swazlLCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/config/certificate/test-cert
```

## Read Certificate Configuration

Returns the previously configured AWS public key.

| Method   | Path                                     | Produces               |
| :------- | :--------------------------------------- | :--------------------- |
| `GET`   | `/auth/aws/config/certificate/:cert_name` | `200 application/json` |

### Parameters

- `cert_name` `(string: <required>)` - Name of the certificate.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/config/certificate/test-cert
```

### Sample Response

```json
{
        "data": {
                "aws_public_cert": "-----BEGIN CERTIFICATE-----\nMIIC7TCCAq0CCQCWukjZ5V4aZzAJBgcqhkjOOAQDMFwxCzAJBgNVBAYTAlVTMRkw\nFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYD\nVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQzAeFw0xMjAxMDUxMjU2MTJaFw0z\nODAxMDUxMjU2MTJaMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9u\nIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNl\ncnZpY2VzIExMQzCCAbcwggEsBgcqhkjOOAQBMIIBHwKBgQCjkvcS2bb1VQ4yt/5e\nih5OO6kK/n1Lzllr7D8ZwtQP8fOEpp5E2ng+D6Ud1Z1gYipr58Kj3nssSNpI6bX3\nVyIQzK7wLclnd/YozqNNmgIyZecN7EglK9ITHJLP+x8FtUpt3QbyYXJdmVMegN6P\nhviYt5JH/nYl4hh3Pa1HJdskgQIVALVJ3ER11+Ko4tP6nwvHwh6+ERYRAoGBAI1j\nk+tkqMVHuAFcvAGKocTgsjJem6/5qomzJuKDmbJNu9Qxw3rAotXau8Qe+MBcJl/U\nhhy1KHVpCGl9fueQ2s6IL0CaO/buycU1CiYQk40KNHCcHfNiZbdlx1E9rpUp7bnF\nlRa2v1ntMX3caRVDdbtPEWmdxSCYsYFDk4mZrOLBA4GEAAKBgEbmeve5f8LIE/Gf\nMNmP9CM5eovQOGx5ho8WqD+aTebs+k2tn92BBPqeZqpWRa5P/+jrdKml1qx4llHW\nMXrs3IgIb6+hUIB+S8dz8/mmO0bpr76RoZVCXYab2CZedFut7qc3WUH9+EUAH5mw\nvSeDCOUMYQR7R9LINYwouHIziqQYMAkGByqGSM44BAMDLwAwLAIUWXBlk40xTwSw\n7HX32MxXYruse9ACFBNGmdX2ZBrVNGrN9N2f6ROk0k9K\n-----END CERTIFICATE-----\n",
                "type": "pkcs7"
        }
}
```

## Delete Certificate Configuration

Removes the previously configured AWS public key.

| Method   | Path                                      | Produces               |
| :------- | :---------------------------------------- | :--------------------- |
| `DELETE` | `/auth/aws/config/certificate/:cert_name` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/config/certificate/test-cert
```

## List Certificate Configurations

Lists all the AWS public certificates that are registered with the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/aws/config/certificates` | `200 application/json` |
| `GET`   | `/auth/aws/config/certificates?list=true` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/aws/config/certificates
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "cert1"
    ]
  }
}
```

## Create STS Role

Allows the explicit association of STS roles to satellite AWS accounts
(i.e. those which are not the account in which the Vault server is
running.) Login attempts from EC2 instances running in these accounts will
be verified using credentials obtained by assumption of these STS roles.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/config/sts/:account_id` | `204 (empty body)`     |

### Parameters

- `account_id` `(string: <required>)` - AWS account ID to be associated with
  STS role. If set, Vault will use assumed credentials to verify any login
  attempts from EC2 instances in this account.
- `sts_role` `(string: <required>)` - AWS ARN for STS role to be assumed when
  interacting with the account specified.  The Vault server must have
  permissions to assume this role.

### Sample Payload

```json
{
  "sts_role": "arn:aws:iam:111122223333:role/myRole"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/config/sts/111122223333
```

## Read STS Role

Returns the previously configured STS role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/config/sts/:account_id` | `200 application/json` |

### Parameters

- `account_id` `(string: <required>)` - AWS account ID to be associated with
  STS role. If set, Vault will use assumed credentials to verify any login
  attempts from EC2 instances in this account.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/config/sts/111122223333
```

### Sample Response

```json
{
  "data": {
    "sts_role ": "arn:aws:iam:111122223333:role/myRole"
  }
}
```

## List STS Roles

Lists all the AWS Account IDs for which an STS role is registered.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/aws/config/sts`       | `200 application/json` |
| `GET`   | `/auth/aws/config/sts?list=true`       | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/aws/config/sts
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "111122223333",
      "999988887777"
    ]
  }
}
```

## Delete STS Role

Deletes a previously configured AWS account/STS role association.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/aws/config/sts`       | `204 (empty body)`  |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/config/sts
```

## Configure Identity Whitelist Tidy Operation

Configures the periodic tidying operation of the whitelisted identity entries.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/config/tidy/identity-whitelist` | `204 (empty body)`     |

### Parameters

- `safety_buffer` `(string: "72h")` - The amount of extra time that must have
  passed beyond the `roletag` expiration, before it is removed from the backend
  storage. Defaults to 72h.
- `disable_periodic_tidy` `(bool: false)` - If set to 'true', disables the
  periodic tidying of the `identity-whitelist/<instance_id>` entries.

### Sample Payload

```json
{
  "safety_buffer": "48h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/config/tidy/identity-whitelist
```

## Read Identity Whitelist Tidy Settings

Returns the previously configured periodic whitelist tidying settings.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/config/tidy/identity-whitelist` | `200 application/json`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/config/tidy/identity-whitelist
```

### Sample Response

```json
{
  "data": {
    "safety_buffer": 600,
    "disable_periodic_tidy": false
  }
}
```

## Delete Identity Whitelist Tidy Settings

Deletes the previously configured periodic whitelist tidying settings.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE`   | `/auth/aws/config/tidy/identity-whitelist` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/config/tidy/identity-whitelist
```

## Configure Role Tag Blacklist Tidy Operation

Configures the periodic tidying operation of the blacklisted role tag entries.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/config/tidy/roletag-blacklist` | `204 (empty body)`     |

### Parameters

- `safety_buffer` `(string: "72h")` - The amount of extra time that must have
  passed beyond the `roletag` expiration, before it is removed from the backend
  storage. Defaults to 72h.
- `disable_periodic_tidy` `(bool: false)` - If set to 'true', disables the
  periodic tidying of the `roletag-blacklist/<instance_id>` entries.

### Sample Payload

```json
{
  "safety_buffer": "48h"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/config/tidy/roletag-blacklist
```

## Read Role Tag Blackist Tidy Settings

Returns the previously configured periodic blacklist tidying settings.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/config/tidy/roletag-blacklist` | `200 application/json`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/config/tidy/roletag-blacklist
```

### Sample Response

```json
{
  "data": {
    "safety_buffer": 600,
    "disable_periodic_tidy": false
  }
}
```

## Delete Role Tag Blackist Tidy Settings

Deletes the previously configured periodic blacklist tidying settings.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE`   | `/auth/aws/config/tidy/roletag-blacklist` | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/config/tidy/roletag-blacklist
```

## Create Role

Registers a role in the backend. Only those instances or principals which
are using the role registered using this endpoint, will be able to perform
the login operation. Contraints can be specified on the role, that are
applied on the instances or principals attempting to login. At least one
constraint should be specified on the role. The available constraints you
can choose are dependent on the `auth_type` of the role and, if the
`auth_type` is `iam`, then whether inferencing is enabled. A role will not
let you configure a constraint if it is not checked by the `auth_type` and
inferencing configuration of that role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/role/:role`       | `204 (empty body)`     |

### Parameters

- `role` `(string: <required>)` - Name of the role.
- `auth_type` `(string: "iam")` - The auth type permitted for this role. Valid
  choices are "ec2" or "iam".  If no value is specified, then it will default to
  "iam" (except for legacy `aws-ec2` auth types, for which it will default to
  "ec2"). Only those bindings applicable to the auth type chosen will be allowed
  to be configured on the role.
- `bound_ami_id` `(string: "")` - If set, defines a constraint on the EC2
  instances that they should be using the AMI ID specified by this parameter.
  This constraint is checked during ec2 auth as well as the iam auth method only
  when inferring an EC2 instance.
- `bound_account_id` `(string: "")` - If set, defines a constraint on the EC2
  instances that the account ID in its identity document to match the one
  specified by this parameter. This constraint is checked during ec2 auth as
  well as the iam auth method only when inferring an EC2 instance.
- `bound_region` `(string: "")` - If set, defines a constraint on the EC2
  instances that the region in its identity document must match the one
  specified by this parameter. This constraint is only checked by the ec2 auth
  method as well as the iam auth method only when inferring an ec2 instance.
- `bound_vpc_id` `(string: "")` - If set, defines a constraint on the EC2
  instance to be associated with the VPC ID that matches the value specified by
  this parameter. This constraint is only checked by the ec2 auth method as well
  as the iam auth method only when inferring an ec2 instance.
- `bound_subnet_id` `(string: "")` - If set, defines a constraint on the EC2
  instance to be associated with the subnet ID that matches the value specified
  by this parameter. This constraint is only checked by the ec2 auth method as
  well as the iam auth method only when inferring an ec2 instance.
- `bound_iam_role_arn` `(string: "")` - If set, defines a constraint on the
  authenticating EC2 instance that it must match the IAM role ARN specified by
  this parameter.  The value is refix-matched (as though it were a glob ending
  in `*`).  The configured IAM user or EC2 instance role must be allowed to
  execute the `iam:GetInstanceProfile` action if this is specified. This
  constraint is checked by the ec2 auth method as well as the iam auth method
  only when inferring an EC2 instance.
- `bound_iam_instance_profile_arn` `(string: "")` - If set, defines a constraint
  on the EC2 instances to be associated with an IAM instance profile ARN which
  has a prefix that matches the value specified by this parameter. The value is
  prefix-matched (as though it were a glob ending in `*`). This constraint is
  checked by the ec2 auth method as well as the iam auth method only when
  inferring an ec2 instance.
- `role_tag` `(string: "")` - If set, enables the role tags for this role. The
  value set for this field should be the 'key' of the tag on the EC2 instance.
  The 'value' of the tag should be generated using `role/<role>/tag` endpoint.
  Defaults to an empty string, meaning that role tags are disabled. This
  constraint is valid only with the ec2 auth method and is not allowed when an
  auth_type is iam.
- `bound_iam_principal_arn` `(string: "")` - Defines the IAM principal that must
  be authenticated using the iam auth method. It should look like
  "arn:aws:iam::123456789012:user/MyUserName" or
  "arn:aws:iam::123456789012:role/MyRoleName". Wildcards are supported at the
  end of the ARN, e.g., "arn:aws:iam::123456789012:\*" will match any IAM
  principal in the AWS account 123456789012. This constraint is only checked by
  the iam auth method. Wildcards are supported at the end of the ARN, e.g.,
  "arn:aws:iam::123456789012:role/\*" will match all roles in the AWS account.
- `inferred_entity_type` `(string: "")` -  When set, instructs Vault to turn on
  inferencing. The only current valid value is "ec2\_instance" instructing Vault
  to infer that the role comes from an EC2 instance in an IAM instance profile.
  This only applies to the iam auth method. If you set this on an existing role
  where it had not previously been set, tokens that had been created prior will
  not be renewable; clients will need to get a new token.
- `inferred_aws_region` `(string: "")` - When role inferencing is activated, the
  region to search for the inferred entities (e.g., EC2 instances). Required if
  role inferencing is activated. This only applies to the iam auth method.
- `resolve_aws_unique_ids` `(bool: false)` - When set, resolves the
  `bound_iam_principal_arn` to the
  [AWS Unique ID](http://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-unique-ids)
  for the bound principal ARN. This field is ignored when
  `bound_iam_principal_arn` ends with a wildcard character.
  This requires Vault to be able to call `iam:GetUser` or `iam:GetRole` on the
  `bound_iam_principal_arn` that is being bound. Resolving to internal AWS IDs
  more closely mimics the behavior of AWS services in that if an IAM user or
  role is deleted and a new one is recreated with the same name, those new users
  or roles won't get access to roles in Vault that were permissioned to the
  prior principals of the same name. The default value for new roles is true,
  while the default value for roles that existed prior to this option existing
  is false (you can check the value for a given role using the GET method on the
  role). Any authentication tokens created prior to this being supported won't
  verify the unique ID upon token renewal.  When this is changed from false to
  true on an existing role, Vault will attempt to resolve the role's bound IAM
  ARN to the unique ID and, if unable to do so, will fail to enable this option.
  Changing this from `true` to `false` is not supported; if absolutely
  necessary, you would need to delete the role and recreate it explicitly
  setting it to `false`. However; the instances in which you would want to do
  this should be rare. If the role creation (or upgrading to use this) succeed,
  then Vault has already been able to resolve internal IDs, and it doesn't need
  any further IAM permissions to authenticate users. If a role has been deleted
  and recreated, and Vault has cached the old unique ID, you should just call
  this endpoint specifying the same `bound_iam_principal_arn` and, as long as
  Vault still has the necessary IAM permissions to resolve the unique ID, Vault
  will update the unique ID. (If it does not have the necessary permissions to
  resolve the unique ID, then it will fail to update.) If this option is set to
  false, then you MUST leave out the path component in bound_iam_principal_arn
  for **roles** only, but not IAM users. That is, if your IAM role ARN is of the
  form `arn:aws:iam::123456789012:role/some/path/to/MyRoleName`, you **must**
  specify a bound_iam_principal_arn of
  `arn:aws:iam::123456789012:role/MyRoleName` for authentication to work.
- `ttl` `(string: "")` - The TTL period of tokens issued using this role,
  provided as "1h", where hour is the largest suffix.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens issued using
  this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.  The maximum allowed lifetime of tokens issued using
  this role.
- `policies` `(array: [])` - Policies to be set on tokens issued using this
  role.
- `allow_instance_migration` `(bool: false)` - If set, allows migration of the
  underlying instance where the client resides. This keys off of pendingTime in
  the metadata document, so essentially, this disables the client nonce check
  whenever the instance is migrated to a new host and pendingTime is newer than
  the previously-remembered time. Use with caution. This only applies to
  authentications via the ec2 auth method. This is mutually exclusive with
  `disallow_reauthentication`.
- `disallow_reauthentication` `(bool: false)` - If set, only allows a single
  token to be granted per instance ID. In order to perform a fresh login, the
  entry in whitelist for the instance ID needs to be cleared using
  'auth/aws/identity-whitelist/<instance_id>' endpoint. Defaults to 'false'.
  This only applies to authentications via the ec2 auth method. This is mutually
  exclusive with `allow_instance_migration`.

### Sample Payload

```json
{
  "bound_ami_id": "ami-fce36987",
  "role_tag": "",
  "policies": [
    "default",
    "dev",
    "prod"
  ],
  "max_ttl": 1800000,
  "disallow_reauthentication": false,
  "allow_instance_migration": false
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/role/dev-role
```

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/role/:role`        | `200 application/json` |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/role/dev-role
```

### Sample Response

```json
{
  "data": {
    "bound_ami_id": "ami-fce36987",
    "role_tag": "",
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "max_ttl": 1800000,
    "disallow_reauthentication": false,
    "allow_instance_migration": false
  }
}
```

## List Roles

Lists all the roles that are registered with the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/aws/roles`       | `200 application/json` |
| `GET`   | `/auth/aws/roles?list=true`       | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/aws/roles
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "dev-role",
      "prod-role"
    ]
  }
}
```

## Delete Role

Deletes the previously registered role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/aws/role/:role`       | `204 (empty body)`  |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/role/dev-role
```

## Create Role Tags

Creates a role tag on the role, which help in restricting the capabilities
that are set on the role. Role tags are not tied to any specific ec2
instance unless specified explicitly using the `instance_id` parameter. By
default, role tags are designed to be used across all instances that
satisfies the constraints on the role. Regardless of which instances have
role tags on them, capabilities defined in a role tag must be a strict
subset of the given role's capabilities.  Note that, since adding and
removing a tag is often a widely distributed privilege, care needs to be
taken to ensure that the instances are attached with correct tags to not
let them gain more privileges than what were intended.  If a role tag is
changed, the capabilities inherited by the instance will be those defined
on the new role tag. Since those must be a subset of the role
capabilities, the role should never provide more capabilities than any
given instance can be allowed to gain in a worst-case scenario.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/role/:role/tag`   | `200 application/json` |

### Parameters

- `role` `(string: <required>)` - Name of the role.
- `policies` `(array: [])` - Policies to be associated with the tag. If set,
  must be a subset of the role's policies. If set, but set to an empty value,
  only the 'default' policy will be given to issued tokens.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens issued using
  this role.
- `instance_id` `(string: "")` - Instance ID for which this tag is intended for.
  If set, the created tag can only be used by the instance with the given ID.
- `allow_instance_migration` `(bool: false)` - If set, allows migration of the
  underlying instance where the client resides. This keys off of pendingTime in
  the metadata document, so essentially, this disables the client nonce check
  whenever the instance is migrated to a new host and pendingTime is newer than
  the previously-remembered time. Use with caution. Defaults to 'false'.
  Mutually exclusive with `disallow_reauthentication`.
- `disallow_reauthentication` `(bool: false)` - If set, only allows a single
  token to be granted per instance ID. This can be cleared with the
  auth/aws/identity-whitelist endpoint. Defaults to 'false'. Mutually exclusive
  with `allow_instance_migration`.

### Sample Payload

```json
{
  "policies": ["default", "prod"]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/role/dev-role/tag
```

### Sample Response

```json
{
  "data": {
    "tag_value": "v1:09Vp0qGuyB8=:r=dev-role:p=default,prod:d=false:t=300h0m0s:uPLKCQxqsefRhrp1qmVa1wsQVUXXJG8UZP/pJIdVyOI=",
    "tag_key": "VaultRole"
  }
}
```

## Login

Fetch a token. This endpoint verifies the pkcs7 signature of the instance
identity document or the signature of the signed GetCallerIdentity request.
With the ec2 auth method, or when inferring an EC2 instance, verifies that
the instance is actually in a running state.  Cross checks the constraints
defined on the role with which the login is being performed. With the ec2
auth method, as an alternative to pkcs7 signature, the identity document
along with its RSA digest can be supplied to this endpoint.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/login`            | `200 application/json` |

### Sample Payload

- `role` `(string: "")` - Name of the role against which the login is being
  attempted.  If `role` is not specified, then the login endpoint looks for a
  role bearing the name of the AMI ID of the EC2 instance that is trying to
  login if using the ec2 auth method, or the "friendly name" (i.e., role name or
  username) of the IAM principal authenticated. If a matching role is not found,
  login fails.
- `identity` `(string: <required-ec2>)` - Base64 encoded EC2 instance identity
  document. This needs to be supplied along with the `signature` parameter. If
  using `curl` for fetching the identity document, consider using the option
  `-w 0` while piping the output to `base64` binary.
- `signature` `(string: <required-ec2>)` - Base64 encoded SHA256 RSA signature of
  the instance identity document. This needs to be supplied along with
  `identity` parameter when using the ec2 auth method.
- `pkcs7` `(string: <required-ec2>)` - PKCS7 signature of the identity document with
  all `\n` characters removed.  Either this needs to be set *OR* both `identity`
  and `signature` need to be set when using the ec2 auth method.
- `nonce` `(string: "")` - The nonce to be used for subsequent login requests.
  If this parameter is not specified at all and if reauthentication is allowed,
  then the backend will generate a random nonce, attaches it to the instance's
  identity-whitelist entry and returns the nonce back as part of auth metadata.
  This value should be used with further login requests, to establish client
  authenticity. Clients can choose to set a custom nonce if preferred, in which
  case, it is recommended that clients provide a strong nonce.  If a nonce is
  provided but with an empty value, it indicates intent to disable
  reauthentication. Note that, when `disallow_reauthentication` option is
  enabled on either the role or the role tag, the `nonce` holds no significance.
  This is ignored unless using the ec2 auth method.
- `iam_http_request_method` `(string: <required-iam>)` - HTTP method used in the
  signed request. Currently only POST is supported, but other methods may be
  supported in the future. This is required when using the iam auth method.
- `iam_request_url` `(string: <required-iam>)` - Base64-encoded HTTP URL used in
  the signed request. Most likely just `aHR0cHM6Ly9zdHMuYW1hem9uYXdzLmNvbS8=`
  (base64-encoding of `https://sts.amazonaws.com/`) as most requests will
  probably use POST with an empty URI. This is required when using the iam auth
  method.
- `iam_request_body` `(string: <required-iam>)` - Base64-encoded body of the
  signed request. Most likely
  `QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ==` which is the
  base64 encoding of `Action=GetCallerIdentity&Version=2011-06-15`. This is
  required when using the iam auth method.
- `iam_request_headers` `(string: <required-iam>)` - Base64-encoded,
  JSON-serialized representation of the sts:GetCallerIdentity HTTP request
  headers. The JSON serialization assumes that each header key maps to either a
  string value or an array of string values (though the length of that array
  will probably only be one). If the `iam_server_id_header_value` is configured
  in Vault for the aws auth mount, then the headers must include the
  X-Vault-AWS-IAM-Server-ID header, its value must match the value configured,
  and the header must be included in the signed headers.  This is required when
  using the iam auth method.


### Sample Payload

```json
{}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/aws/login
```

### Sample Response

```json
{
  "auth": {
    "renewable": true,
    "lease_duration": 1800000,
    "metadata": {
      "role_tag_max_ttl": "0",
      "instance_id": "i-de0f1344",
      "ami_id": "ami-fce36983",
      "role": "dev-role",
      "auth_type": "ec2"
    },
    "policies": [
      "default",
      "dev"
    ],
    "accessor": "20b89871-e6f2-1160-fb29-31c2f6d4645e",
    "client_token": "c9368254-3f21-aded-8a6f-7c818e81b17a"
  }
}
```

## Place Role Tags in Blacklist

Places a valid role tag in a blacklist. This ensures that the role tag
cannot be used by any instance to perform a login operation again.  Note
that if the role tag was previously used to perform a successful login,
placing the tag in the blacklist does not invalidate the already issued
token.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/roletag-blacklist/:role_tag`            | `204 (empty body)` |

### Parameters

- `role_tag` `(string: <required>)` - Role tag to be blacklisted. The tag can be
  supplied as-is. In order to avoid any encoding problems, it can be base64
  encoded.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/auth/aws/roletag-blacklist/djE6MDlWcDBxR3V5Qjg9OmE9YW1pLWZjZTNjNjk2OnA9ZGVmYXVsdCxwcm9kOmQ9ZmFsc2U6dD0zMDBoMG0wczp1UExLQ1F4cXNlZlJocnAxcW1WYTF3c1FWVVhYSkc4VVpQLwo=
```

### Read Role Tag Blacklist Information

Returns the blacklist entry of a previously blacklisted role tag.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/roletag-blacklist/:role_tag`            | `200 application/json` |

### Parameters

- `role_tag` `(string: <required>)` - Role tag to be blacklisted. The tag can be
  supplied as-is. In order to avoid any encoding problems, it can be base64
  encoded.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/roletag-blacklist/djE6MDlWcDBxR3V5Qjg9OmE9YW1pLWZjZTNjNjk2OnA9ZGVmYXVsdCxwcm9kOmQ9ZmFsc2U6dD0zMDBoMG0wczp1UExLQ1F4cXNlZlJocnAxcW1WYTF3c1FWVVhYSkc4VVpQLwo=
```


### Sample Response

```json
{
  "data": {
    "expiration_time": "2016-04-25T10:35:20.127058773-04:00",
    "creation_time": "2016-04-12T22:35:01.178348124-04:00"
  }
}
```

## List Blacklist Tags

Lists all the role tags that are blacklisted.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/aws/roletag-blacklist`       | `200 application/json` |
| `GET`   | `/auth/aws/roletag-blacklist?list=true`       | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/aws/roletag-blacklist
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "v1:09Vp0qGuyB8=:a=ami-fce3c696:p=default,prod:d=false:t=300h0m0s:uPLKCQxqsefRhrp1qmVa1wsQVUXXJG8UZP/"
    ]
  }
}
```

## Delete Blacklist Tags

Deletes a blacklisted role tag.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/aws/roletag-blacklist/:role_tag`       | `204 (empty body)`  |

### Parameters

- `role_tag` `(string: <required>)` - Role tag to be blacklisted. The tag can be
  supplied as-is. In order to avoid any encoding problems, it can be base64
  encoded.


### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/roletag-blacklist/djE6MDlWcDBxR3V5Qjg9OmE9YW1pLWZjZTNjNjk2OnA9ZGVmYXVsdCxwcm9kOmQ9ZmFsc2U6dD0zMDBoMG0wczp1UExLQ1F4cXNlZlJocnAxcW1WYTF3c1FWVVhYSkc4VVpQLwo=
```

## Tidy Blacklist Tags

Cleans up the entries in the blacklist based on expiration time on the entry and
`safety_buffer`.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/tidy/roletag-blacklist` | `204 (empty body)` |

### Parameters

- `safety_buffer` `(string: "72h")` - The amount of extra time that must have
  passed beyond the `roletag` expiration, before it is removed from the backend
  storage. Defaults to 72h.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/auth/aws/tidy/roletag-blacklist
```

### Read Identity Whitelist Information

Returns an entry in the whitelist. An entry will be created/updated by every
successful login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/aws/identity-whitelist/:instance_id`            | `200 application/json` |

### Parameters

- `instance_id` `(string: <required>)` - EC2 instance ID. A successful login
  operation from an EC2 instance gets cached in this whitelist, keyed off of
  instance ID.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/aws/identity-whitelist/i-aab47d37
```


### Sample Response

```json
{
  "data": {
    "pending_time": "2016-04-14T01:01:41Z",
    "expiration_time": "2016-05-05 10:09:16.67077232 +0000 UTC",
    "creation_time": "2016-04-14 14:09:16.67077232 +0000 UTC",
    "client_nonce": "5defbf9e-a8f9-3063-bdfc-54b7a42a1f95",
    "role": "dev-role"
  }
}
```

## List Identity Whitelist Entries

  Lists all the instance IDs that are in the whitelist of successful logins.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/aws/identity-whitelist`       | `200 application/json` |
| `GET`   | `/auth/aws/identity-whitelist?list=true`       | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/aws/roletag-blacklist
```

### Sample Response

```json
{
  "data": {
    "keys": [
      "i-aab47d37"
    ]
  }
}
```

## Delete Identity Whitelist Entries

Deletes a cache of the successful login from an instance.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/aws/identity-whitelist/:instance_id`       | `204 (empty body)`  |

### Parameters

- `instance_id` `(string: <required>)` - EC2 instance ID. A successful login
  operation from an EC2 instance gets cached in this whitelist, keyed off of
  instance ID.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/aws/identity-whitelist/i-aab47d37
```

## Tidy Identity Whitelist Entries

Cleans up the entries in the whitelist based on expiration time and
`safety_buffer`.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/aws/tidy/identity-whitelist` | `204 (empty body)` |

### Parameters

- `safety_buffer` `(string: "72h")` - The amount of extra time that must have
  passed beyond the `roletag` expiration, before it is removed from the backend
  storage. Defaults to 72h.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    https://vault.rocks/v1/auth/aws/tidy/identity-whitelist
```
