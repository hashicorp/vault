---
layout: "api"
page_title: "Google Cloud Platform Auth Plugin Backend - HTTP API"
sidebar_current: "docs-http-auth-gcp"
description: |-
  This is the API documentation for the Vault GCP authentication
  backend plugin.
---

# GCP Auth Plugin HTTP API

This is the API documentation for the Vault GCP authentication backend
plugin. To learn more about the usage and operation, see the
[Vault GCP backend documentation](/docs/auth/gcp.html).

This documentation assumes the plugin backend is mounted at the
`/auth/gcp` path in Vault. Since it is possible to mount auth backends
at any location, please update your API calls accordingly.

## Configure

Configures the credentials required for the plugin to perform API calls
to GCP. These credentials will be used to query the status of IAM
entities and get service account or other Google public certificates
to confirm signed JWTs passed in during login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/config`          | `204 (empty body)`     |

### Parameters

- `credentials` `(string: "")` - A marshaled JSON string that is the content
  of a GCP credentials file. If you would rather specify a file, you can use
  `credentials="@path/to/creds.json`. The GCP permissions
  Vault currently requires are:
    - `iam.serviceAccounts.get`
    - `iam.serviceAccountKeys.get`

  If this value is not specified or if it is explicitly set to empty,
  Vault will attempt to use [Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
  for that server's machine.

### Sample Payload

```json
{
  "credentials": "{ \"type\": \"service_account\", \"project_id\": \"project-123456\",...}"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/gcp/config
```

## Read Config

Returns the previously configured config, including credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/gcp/config`           | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/gcp/config
```

### Sample Response

```json
{
  "data":{
    "client_email":"serviceaccount1@project-123456.iam.gserviceaccount.com",
    "client_id":"...",
    "private_key":"-----BEGIN PRIVATE KEY-----...-----END PRIVATE KEY-----\n",
    "private_key_id":"...",
    "project_id":"project-123456"
  },
  ...
}

```

## Delete Config

Deletes the previously configured GCP config and credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/gcp/config`           | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/gcp/config
```

## Create Role

Registers a role in the backend. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/role/:name`       | `204 (empty body)`     |

### Parameters
- `name` `(string: <required>)` - Name of the role.
- `type` `(string: <required>)` - The type of this role. Only the
  restrictions applicable to this role type will be allowed to
  be configured on the role (see below). Valid choices are: `iam`.
- `project_id` `(string: "")` - Required. Only entities belonging to this
  project can login for this role.
- `ttl` `(string: "")` - The TTL period of tokens issued using this role in
  seconds.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens
  issued in seconds using this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.
- `policies` `(array: [])` - Policies to be set on tokens issued using this
  role.
- `max_jwt_exp` `(string: "")` - Optional, defaults to 900 (15min).
  Number of seconds past the time of authentication that the login param JWT
  must expire within. For example, if a user attempts to login with a token
  that expires within an hour and this is set to 15 minutes, Vault will return
  an error prompting the user to create a new signed JWT with a shorter `exp`.

**`iam`-only params**:

- `service_accounts` `(array: [])` - Required for `iam` roles.
  A comma-separated list of service account emails or ids.
  Defines the service accounts that login is restricted to. If set to `\*`, all
  service accounts are allowed (role will still be bound by project).

### Sample Payload

`iam` Role:

```json
{
  "type": "iam",
  "project": "project-123456",
  "policies": [
    "default",
    "dev",
    "prod"
  ],
  "max_ttl": 1800000,
  "max_jwt_exp": 10000,
  "service_accounts": [
    "dev-1@project-123456.iam.gserviceaccount.com",
    "dev-2@project-123456.iam.gserviceaccount.com",
    "123456789",
  ],
  "allow_instance_migration": false
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/gcp/role/dev-role
```

## Edit Service Accounts For IAM Role

Edit service accounts for an existing IAM role in the backend.
This allows you to add or remove service accounts from the list of
service accounts on the role.

| Method   | Path                                    | Produces           |
| :------- | :---------------------------------------| :------------------|
| `POST`   | `/auth/gcp/role/:name/service-accounts` | `204 (empty body)` |

### Parameters
- `name` `(string: <required>)` - Name of an existing `iam` role.
    Returns error if role is not an `iam` role.
- `add` `(array: [])` - List of service accounts to add to the role's
    service accounts
- `remove` `(array: [])` - List of service accounts to remove from the
    role's service accounts

### Sample Payload

```json
{
  "add": [
      "dev-1@project-123456.iam.gserviceaccount.com",
      "123456789",
  ],
  "remove": [
      "dev-2@project-123456.iam.gserviceaccount.com",
  ],
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/gcp/role/dev-role
```

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/gcp/role/:name`        | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/gcp/role/dev-role
```

### Sample Response

```json
{
    "data":{
        "max_jwt_exp": 900,
        "max_ttl": 0,
        "ttl":0,
        "period": 0,
        "policies":[
            "default",
            "dev",
            "prod"
        ],
        "project_id":"project-123456",
        "role_type":"iam",
        "service_accounts": [
            "dev-1@project-123456.iam.gserviceaccount.com",
            "dev-2@project-123456.iam.gserviceaccount.com",
            "123456789",
        ]
    },
    ...
}

```

## List Roles

Lists all the roles that are registered with the plugin.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/gcp/roles`            | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/gcp/roles
```

### Sample Response

```json  
{
    "data": {
        "keys": [
            "dev-role",
            "prod-role"
        ]
    },
    ...
}
```

## Delete Role

Deletes the previously registered role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/gcp/role/:role`       | `204 (empty body)`     |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/gcp/role/dev-role
```

## Login

Fetch a token. This endpoint takes a signed JSON Web Token (JWT) and
a role name for some entity. It verifies the JWT signature to authenticate that
entity and then authorizes the entity for the given role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/login`            | `200 application/json` |

### Sample Payload

- `role` `(string: "")` - Name of the role against which the login is being
  attempted.
- `jwt` `(string: "")` - Signed [JSON Web Token](https://tools.ietf.org/html/rfc7519) (JWT).
  For `iam`, this is a JWT generated using the IAM API method
  [signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt)
  or a self-signed JWT.


### Sample Payload

```json
{
    "role": "dev-role",
    "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/gcp/login
```

### Sample Response

```json
{
    "auth":{
        "client_token":"f33f8c72-924e-11f8-cb43-ac59d697597c",
        "accessor":"0e9e354a-520f-df04-6867-ee81cae3d42d",
        "policies":[
            "default",
            "dev",
            "prod"
        ],
        "metadata":{
            "role": "dev-role",
            "service_account_email": "dev1@project-123456.iam.gserviceaccount.com",
            "service_account_id": "111111111111111111111"
        },
        "lease_duration":2764800,
        "renewable":true
    },
    ...
}
```
