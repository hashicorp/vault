---
layout: "api"
page_title: "AliCloud - Auth Methods - HTTP API"
sidebar_title: "AliCloud"
sidebar_current: "api-http-auth-alicloud"
description: |-
  This is the API documentation for the Vault AliCloud auth method.
---

# AliCloud Auth Method (API)

This is the API documentation for the Vault AliCloud auth method. For
general information about the usage and operation of the AliCloud method, please
see the [Vault AliCloud auth method documentation](/docs/auth/alicloud.html).

This documentation assumes the AliCloud auth method is mounted at the `/auth/alicloud`
path in Vault. Since it is possible to enable auth methods at any location,
please update your API calls accordingly.

## Create Role

Registers a role. Only entities using the role registered using this endpoint 
will be able to perform the login operation.

| Method   | Path                             |
| :------------------------------- | :--------------------- |
| `POST`   | `/auth/alicloud/role/:role`      |

### Parameters

- `role` `(string: <required>)` - Name of the role. Must correspond with the name of the role reflected in the arn.
- `arn` `(string: <required>)` - The role's arn.
- `policies` `(array: [])` - Policies to be set on tokens issued using this
  role.
- `ttl` `(string: "")` - The TTL period of tokens issued using this role,
  provided as "1h", where hour is the largest suffix.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens issued using
  this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.
- `bound_cidrs` `(string: "", or list: [])` – If set, restricts usage of the
  roles to client IPs falling within the range of the specified CIDR(s).

### Sample Payload

```json
{
  "arn": "acs:ram::5138828231865461:role/dev-role",
  "policies": [
    "dev",
    "prod"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/alicloud/role/dev-role
```

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/auth/alicloud/role/:role`  |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/alicloud/role/dev-role
```

### Sample Response

```json
{
  "data": {
    "arn": "acs:ram::5138828231865461:role/dev-role",
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "ttl": 1800000,
    "max_ttl": 1800000,
    "period": 0
  }
}
```

## List Roles

Lists all the roles that are registered with the method.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `LIST`   | `/auth/alicloud/roles`       |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    http://127.0.0.1:8200/v1/auth/alicloud/roles
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

| Method   | Path                             |
| :------------------------------- | :--------------------- |
| `DELETE` | `/auth/alicloud/role/:role`      |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/alicloud/role/dev-role
```

## Login

Fetch a token. This endpoint verifies the signature of the signed 
GetCallerIdentity request.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`   | `/auth/alicloud/login`       |

### Parameters

- `role` `(string: <required>)` - Name of the role.
- `identity_request_url` `(string: <required>)` - Base64-encoded HTTP URL used in
  the signed request.
- `identity_request_headers` `(string: <required>)` - Base64-encoded,
  JSON-serialized representation of the sts:GetCallerIdentity HTTP request
  headers. The JSON serialization assumes that each header key maps to either a
  string value or an array of string values (though the length of that array
  will probably only be one).


### Sample Payload

```json
{
  "role": "dev-role",
  "identity_request_url": "aWRlbnRpdHlabrVxdWVzdF91cmw=",
  "identity_request_headers": "aWRlimRpdHlfcmVxdWVzdF9oZWFkZXJz"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/alicloud/login
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
      "auth_type": "ec2",
      "account_id":    "5138828231865461",
      "user_id":       "216959339000654321",
      "role_id":       "4657-abcd",
      "arn":           "acs:ram::5138828231865461:assumed-role/dev-role/vm-ram-i-rj978rorvlg76urhqh7q",
      "identity_type": "assumed-role",
      "principal_id":  "vm-ram-i-rj978rorvlg76urhqh7q",
      "request_id":    "D6E46F10-F26C-4AA0-BB69-FE2743D9AE62",
      "role_name":     "dev-role"
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