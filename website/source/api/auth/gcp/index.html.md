---
layout: "api"
page_title: "Google Cloud Platform Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-gcp"
description: |-
  This is the API documentation for the Vault GCP authentication
  backend.
---

# GCP Auth Backend HTTP API

This is the API documentation for the Vault GCP authentication backend.
To learn more about the usage and operation, see the
[Vault GCP backend documentation](/docs/auth/gcp.html).

This documentation assumes the backend is mounted at the `/auth/gcp`
path in Vault. Since it is possible to mount auth backends at any location,
please update your API calls accordingly.

## Configure Client

Configures the backend credentials required to perform API calls to GCP.
These credentials will be used to query the status of IAM entities and get
service account or other Google public certificates to confirm signed JWTs
passed in during login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/config/`          | `204 (empty body)`     |

### Parameters

- `credentials` `(string: "")` - A marshaled JSON string that is the content
  of a GCP credentials file. If not provided, the Vault server attempts
  to use [Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)


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
  "auth":null,
  "warnings":[
    "Read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords."
  ],

  "data":{
    "client_email":"serviceaccount1@project-123456.iam.gserviceaccount.com",
    "client_id":"...",
    "private_key":"-----BEGIN PRIVATE KEY-----...-----END PRIVATE KEY-----\n",
    "private_key_id":"...",
    "project_id":"project-123456"
  },
  "lease_duration":0,
  "lease_id":"",
  "renewable":false,
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
    https://vault.rocks/v1/auth/gcp/config/client
```

## Create Role

Registers a role in the backend. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/role/:role`       | `204 (empty body)`     |

### Parameters

- `role` `(string: <required>)` - Name of the role.
- `type` `(string: "iam")` - The type of this role. Only the  
  restrictions applicable to this role type will be allowed to
  be configured on the role (see below).

  Valid choices are: `iam`.  
- `project` `(string: "")` - Required. Only entities belonging to this
  project can login for this role.
- `ttl` `(string: "")` - The TTL period of tokens issued using this role in
  seconds.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens
  issued in seconds using this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.  The maximum allowed lifetime of tokens issued using
  this role.
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

`iam`:

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

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/gcp/role/:role`        | `200 application/json` |

### Parameters

- `role` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/gcp/role/dev-role
```

### Sample Response

```json
{
  "auth":null,
  "warnings":null,
  "data":{
    "max_jwt_exp": 900,
    "max_ttl": 0,
    "ttl":0,
    "period": 0,
    "policies":[
      "default",
      "dev",
      "prod"],
    "project_id":"project-123456",
    "role_type":"iam",
    "service_accounts": [
      "dev-1@project-123456.iam.gserviceaccount.com",
      "dev-2@project-123456.iam.gserviceaccount.com",
      "123456789",
    ]
  },
  "wrap_info":null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}

```

## List Roles

Lists all the roles that are registered with the backend.

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
  "auth": null,
  "warnings": null,
  "data": {
    "keys": [
      "dev-role",
      "prod-role"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
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
- `jwt` `(string: "")` - Signed JSON Web Token ([JWT](https://tools.ietf.org/html/rfc7519)).
  For `iam`, this is a JWT generated using the IAM API method [signJwt](https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt)
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
  "data": null,
  "lease_id":"",
  "renewable": false,
  "lease_duration": 0,
  "wrap_info": null,
  "warnings": null
}
```
