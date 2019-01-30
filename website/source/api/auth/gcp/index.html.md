---
layout: "api"
page_title: "Google Cloud - Auth Methods - HTTP API"
sidebar_title: "Google Cloud"
sidebar_current: "api-http-auth-gcp"
description: |-
  This is the API documentation for the Vault Google Cloud authentication
  method.
---

# Google Cloud Auth Method (API)

This is the API documentation for the Vault Google Cloud auth method. To learn
more about the usage and operation, see the
[Vault Google Cloud method documentation](/docs/auth/gcp.html).

This documentation assumes the plugin method is mounted at the
`/auth/gcp` path in Vault. Since it is possible to enable auth methods
at any location, please update your API calls accordingly.

## Configure

Configures the credentials required for the plugin to perform API calls
to Google Cloud. These credentials will be used to query the status of IAM
entities and get service account or other Google public certificates
to confirm signed JWTs passed in during login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/config`           | `204 (empty body)`     |

### Parameters

- `credentials` `(string: "")` - A JSON string containing the contents
  of a GCP credentials file. The credentials file must have the following
  [permissions](https://cloud.google.com/compute/docs/access/iam):

    ```
    iam.serviceAccounts.get
    iam.serviceAccountKeys.get
    ```

    If this value is empty, Vault will try to use [Application Default
    Credentials][gcp-adc] from the machine on which the Vault server is running.
    
    The project must have the `iam.googleapis.com` API [enabled](https://console.cloud.google.com/flows/enableapi?apiid=iam.googleapis.com).

### Sample Payload

```json
{
  "credentials": "{ \"type\": \"service_account\", \"project_id\": \"project-123456\", ...}"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/gcp/config
```

## Read Config

Returns the configuration, if any, including credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/gcp/config`           | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/gcp/config
```

### Sample Response

```json
{
  "data": {
    "client_email": "service-account@project-123456.iam.gserviceaccount.com",
    "client_id": "123456789101112131415",
    "private_key_id": "97fd7ba59a96e1f3830296aedb4f50879e4d5382",
    "project_id": "project-123456"
  },
}
```

## Delete Config

Deletes all GCP configuration data. This operation is idempotent.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/gcp/config`           | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/gcp/config
```

## Create Role

Registers a role in the method. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/role/:name`       | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` - The name of the role.

- `type` `(string: <required>)` - The type of this role. Certain fields
  correspond to specific roles and will be rejected otherwise. Please see below
  for more information.

- `ttl` `(string: "")` - The TTL period of tokens issued using this role. This
  can be specified as an integer number of seconds or as a duration value like
  "5m".

- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens issued in
  seconds using this role. This can be specified as an integer number of seconds
  or as a duration value like "5m".

- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter. This can be specified as an integer number of seconds
  or as a duration value like "5m".

- `policies` `(array: [default])` - The list of policies to be set on tokens
  issued using this role.

- `bound_service_accounts` `(array: <required for iam>)` - An array of 
   service account emails or IDs that login is restricted to,
   either directly or through an associated instance. If set to
  `*`, all service accounts are allowed (you can bind this further using
  `bound_projects`.)
  
- `bound_projects` `(array: [])` - An array of GCP project IDs. Only entities 
   belonging to this project can authenticate under the role.

- `add_group_aliases` `(bool: false)` - If true, any auth token
   generated under this token will have associated group aliases, namely
   `project-$PROJECT_ID`, `folder-$PROJECT_ID`, and `organization-$ORG_ID`
   for the entities project and all its folder or organization ancestors. This
   requires Vault to have IAM permission `resourcemanager.projects.get`.
    
#### `iam`-only Parameters

The following parameters are only valid when the role is of type `"iam"`:

- `max_jwt_exp` `(string: "15m")` - The number of seconds past the time of
  authentication that the login param JWT must expire within. For example, if a
  user attempts to login with a token that expires within an hour and this is
  set to 15 minutes, Vault will return an error prompting the user to create a
  new signed JWT with a shorter `exp`. The GCE metadata tokens currently do not
  allow the `exp` claim to be customized.

- `allow_gce_inference` `(bool: true)` - A flag to determine if this role should
   allow GCE instances to authenticate by inferring service accounts from the
   GCE identity metadata token.

#### `gce`-only Parameters

The following parameters are only valid when the role is of type `"gce"`:

- `bound_zones` `(array: [])`: The list of zones that a GCE instance must belong
  to in order to be authenticated. If `bound_instance_groups` is provided, it is
  assumed to be a zonal group and the group must belong to this zone.

- `bound_regions` `(array: [])`: The list of regions that a GCE instance must
  belong to in order to be authenticated. If `bound_instance_groups` is
  provided, it is assumed to be a regional group and the group must belong to
  this region. If `bound_zones` are provided, this attribute is ignored.

- `bound_instance_groups` `(array: [])`: The instance groups that an authorized
  instance must belong to in order to be authenticated. If specified, either
  `bound_zones` or `bound_regions` must be set too.

- `bound_labels` `(array: [])`: A comma-separated list of GCP labels formatted
  as "key:value" strings that must be set on authorized GCE instances. Because
  GCP labels are not currently ACL'd, we recommend that this be used in
  conjunction with other restrictions.

### Sample Payload

Example `iam` role:

```json
{
  "type": "iam",
  "project_id": "project-123456",
  "policies": ["prod"],
  "ttl": "30m",
  "max_ttl": "24h",
  "max_jwt_exp": "5m",
  "bound_service_accounts": [
    "dev-1@project-123456.iam.gserviceaccount.com"
  ]
}
```

Example `gce` role:

```json
{
  "type": "gce",
  "project_id": "project-123456",
  "policies": ["prod"],
  "bound_zones": ["us-east1-b", "eu-west2-a"],
  "ttl": "30m",
  "max_ttl": "24h",
  "bound_service_accounts": [
    "dev-1@project-123456.iam.gserviceaccount.com"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/gcp/role/my-role
```

## Edit Service Accounts on IAM Role

Edit service accounts for an existing IAM role in the method.
This allows you to add or remove service accounts from the list of
service accounts on the role.

| Method   | Path                                    | Produces           |
| :------- | :---------------------------------------| :------------------|
| `POST`   | `/auth/gcp/role/:name/service-accounts` | `204 (empty body)` |

### Parameters

- `name` `(string: <required>)` - The name of an existing `iam` type role. This
  will return an error if role is not an `iam` type role.

- `add` `(array: [])` - The list of service accounts to add to the role's
  service accounts.

- `remove` `(array: [])` - The list of service accounts to remove from the
  role's service accounts.

### Sample Payload

```json
{
  "add": [
    "dev-1@project-123456.iam.gserviceaccount.com",
    "123456789"
  ],
  "remove": [
    "dev-2@project-123456.iam.gserviceaccount.com"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/gcp/role/my-role
```

## Edit Labels on GCE Role

Edit labels for an existing GCE role in the backend. This allows you to add or
remove labels (keys, values, or both) from the list of keys on the role.

| Method   | Path                                    | Produces           |
| :------- | :---------------------------------------| :------------------|
| `POST`   | `/auth/gcp/role/:name/labels`           | `204 (empty body)` |

### Parameters

- `name` `(string: <required>)` - The name of an existing `gce` role. This will
  return an error if role is not a `gce` type role.

- `add` `(array: [])` - The list of `key:value` labels to add to the GCE role's
  bound labels.

- `remove` `(array: [])` - The list of label _keys_ to remove from the role's
  bound labels. If any of the specified keys do not exist, no error is returned
  (idempotent).

### Sample Payload

```json
{
  "add": [
    "foo:bar",
    "env:dev",
    "key:value"
  ],
  "remove": [
    "key1",
    "key2"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/gcp/role/my-role
```

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/gcp/role/:name`       | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - The name of the role to read.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/gcp/role/my-role
```

### Sample Response

```json
{
  "data": {
    "bound_labels": {
      "env": "dev",
      "foo": "bar",
      "key": "value"
    },
    "bound_service_accounts": [
      "dev-1@project-123456.iam.gserviceaccount.com"
    ],
    "bound_zones": [
      "eu-west2-a",
      "us-east1-b"
    ],
    "max_ttl": 86400,
    "policies": [
      "prod"
    ],
    "project_id": "project-123456",
    "type": "gce",
    "ttl": 1800
  }
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
    http://127.0.0.1:8200/v1/auth/gcp/roles
```

### Sample Response

```json  
{
  "data": {
    "keys": [
      "my-role",
      "my-other-role"
    ]
  }
}
```

## Delete Role

Deletes the previously registered role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/gcp/role/:role`       | `204 (empty body)`     |

### Parameters

- `role` `(string: <required>)` - The name of the role to delete.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    http://127.0.0.1:8200/v1/auth/gcp/role/my-role
```

## Login

Login to retrieve a Vault token. This endpoint takes a signed JSON Web Token
(JWT) and a role name for some entity. It verifies the JWT signature with Google
Cloud to authenticate that entity and then authorizes the entity for the given
role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/gcp/login`            | `200 application/json` |

### Sample Payload

- `role` `(string: <required>)` - The name of the role against which the login
  is being attempted.

- `jwt` `(string: <required>)` - A Signed [JSON Web Token][jwt].

  - For `iam` type roles, this is a JWT signed with the
  [`signJwt` method][signjwt-method] or a self-signed JWT.

  - For `gce` type roles, this is an [identity metadata token][instance-token].


### Sample Payload

```json
{
  "role": "my-role",
  "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/gcp/login
```

### Sample Response

```json
{
  "auth": {
    "client_token": "f33f8c72-924e-11f8-cb43-ac59d697597c",
    "accessor": "0e9e354a-520f-df04-6867-ee81cae3d42d",
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "metadata": {
      "project_id": "my-project",
      "role": "my-role",
      "service_account_email": "dev1@project-123456.iam.gserviceaccount.com",
      "service_account_id": "111111111111111111111"
    },
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

[gcp-adc]: https://developers.google.com/identity/protocols/application-default-credentials
[jwt]: https://tools.ietf.org/html/rfc7519
[signjwt-method]: https://cloud.google.com/iam/reference/rest/v1/projects.serviceAccounts/signJwt
[instance-token]: https://cloud.google.com/compute/docs/instances/verifying-instance-identity#request_signature
