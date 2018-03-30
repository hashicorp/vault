---
layout: "api"
page_title: "Azure - Auth Methods - HTTP API"
sidebar_current: "docs-http-auth-azure"
description: |-
  This is the API documentation for the Vault Azure authentication
  method plugin.
---

# Azure Auth Method (API)

This is the API documentation for the Vault Azure auth method
plugin. To learn more about the usage and operation, see the
[Vault Azure method documentation](/docs/auth/azure.html).

This documentation assumes the plugin method is mounted at the
`/auth/azure` path in Vault. Since it is possible to enable auth methods
at any location, please update your API calls accordingly.

## Configure

Configures the credentials required for the plugin to perform API calls
to Azure. These credentials will be used to query the metadata about the
virtual machine.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/azure/config`         | `204 (empty body)`     |

### Parameters

- `tenant_id` `(string: <required>)` - The tenant id for the Azure Active Directory organization.
- `resource` `(string: <required>)` - The configured URL for the application registered in Azure Active Directory.
- `client_id` `(string: '')` - The client id for credentials to query the Azure APIs.  Currently read permissions to query compute resources are required.
- `client_secret` `(string: '')` - The client secret for credentials to query the Azure APIs. 

### Sample Payload

```json
{
  "tenant_id": "kd83...",
  "resource": "https://vault.hashicorp.com/",
  "client_id": "12ud...",
  "client_secret": "DUJDS3..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/azure/config
```

# Read Config

Returns the previously configured config, including credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/azure/config`           | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/azure/config
```

### Sample Response

```json
{
  "data":{
    "tenant_id": "kd83...",
    "resource": "https://vault.hashicorp.com/",
    "client_id": "12ud...",
    "client_secret": "DUJDS3..."
  },
  ...
}

```

## Delete Config

Deletes the previously configured Azure config and credentials.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/auth/azure/config`         | `204 (empty body)`     |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/azure/config
```

## Create Role

Registers a role in the method. Role types have specific entities
that can perform login operations against this endpoint. Constraints specific
to the role type must be set on the role. These are applied to the authenticated
entities attempting to login.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/azure/role/:name`       | `204 (empty body)`     |

### Parameters
- `name` `(string: <required>)` - Name of the role.
- `policies` `(array: [])` - Policies to be set on tokens issued using this
  role.
- `ttl` `(string: "")` - The TTL period of tokens issued using this role in
  seconds.
- `max_ttl` `(string: "")` - The maximum allowed lifetime of tokens
  issued in seconds using this role.
- `period` `(string: "")` - If set, indicates that the token generated using
  this role should never expire. The token should be renewed within the duration
  specified by this value. At each renewal, the token's TTL will be set to the
  value of this parameter.
- `bound_service_principal_ids` `(array: [])` - The list of Service Principal IDs 
  that login is restricted to.
- `bound_group_ids` `(array: [])` - The list of group ids that login is restricted 
  to.
- `bound_location` `(array: [])` - The list of locations that login is restricted to.
- `bound_subscription_ids` `(array: [])` - The list of subscription IDs that login 
  is restricted to.
- `bound_resource_group_names` `(array: [])` - The list of resource groups that 
  login is restricted to. 

### Sample Payload

```json
{
  "policies": [
    "default",
    "dev",
    "prod"
  ],
  "max_ttl": 1800000,
  "max_jwt_exp": 10000,
  "bound_resource_groups": [
    "vault-dev",
    "vault-staging",
    "vault-prod"
  ]
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/azure/role/dev-role
```

## Read Role

Returns the previously registered role configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/azure/role/:name`      | `200 application/json` |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/azure/role/dev-role
```

### Sample Response

```json
{
  "data":{
    "policies": [
        "default",
        "dev",
        "prod"
    ],
    "max_ttl": 1800000,
    "max_jwt_exp": 10000,
    "bound_resource_groups": [
        "vault-dev",
        "vault-staging",
        "vault-prod"
    ]
  },
  ...
}

```

## List Roles

Lists all the roles that are registered with the plugin.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/azure/roles`          | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/azure/roles
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
| `DELETE` | `/auth/azure/role/:name`     | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` - Name of the role.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/azure/role/dev-role
```

## Login

Fetch a token. This endpoint takes a signed JSON Web Token (JWT) and
a role name for some entity. It verifies the JWT signature to authenticate that
entity and then authorizes the entity for the given role.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/azure/login`          | `200 application/json` |

### Sample Payload

- `role` `(string: <required>)` - Name of the role against which the login is being
  attempted.
- `jwt` `(string: <required>)` - Signed [JSON Web Token](https://tools.ietf.org/html/rfc7519) (JWT) from Azure MSI.
- `subscription_id` `(string: "")` - The subscription ID for the machine that
  generated the MSI token.  This information can be obtained through instance
  metadata.
- `resource_group_name` `(string: "")` - The resource group for the machine that
  generated the MSI token.  This information can be obtained through instance
  metadata.
- `vm_name` `(string: "")` - The virtual machine name for the machine that
  generated the MSI token.  This information can be obtained through instance
  metadata.

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
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/azure/login
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
        "lease_duration":2764800,
        "renewable":true
    },
    ...
}
```
