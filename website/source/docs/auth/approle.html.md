---
layout: "docs"
page_title: "Auth Backend: AppRole"
sidebar_current: "docs-auth-approle"
description: |-
  The AppRole backend allows machines and services to authenticate with Vault.
---

# Auth Backend: AppRole

This backend allows machines and services (_apps_) to authenticate with Vault
via a series of administratively defined _roles_: AppRoles. The open design of
`AppRole` enables a varied set of workflows and configurations to handle large
numbers of apps and their needs. This backend is oriented to automated
workflows, and is the successor to the `App-ID` backend.

## AppRoles

An AppRole represents a set of Vault policies and login constraints that must
be met to receive a token with those policies. The scope can be as narrow or
broad as desired -- an AppRole can be created for a particular machine, or even
a particular user on that machine, or a service spread across machines. The
credentials required for successful login depend upon on the constraints set on
the AppRole associated with the credentials.

## Credentials/Constraints

### RoleID

RoleID is an identifier that selects the AppRole against which the other
credentials are evaluated. When authenticating against this backend's login
endpoint, the RoleID is a required argument (via `role_id`) at all times. By
default, RoleIDs are unique UUIDs, which allow them to serve as secondary
secrets to the other credential information. However, they can be set to
particular values to match introspected information by the client (for
instance, the client's domain name).

### SecretID

SecretID is a credential that is required by default for any login (via
`secret_id`) and is intended to always be secret. (For advanced usage,
requiring a SecretID can be disabled via an AppRole's `bind_secret_id`
parameter, allowing machines with only knowledge of the RoleID, or matching
other set constraints, to fetch a token). SecretIDs can be created against an
AppRole either via generation of a 128-bit purely random UUID by the role
itself (`Pull` mode) or via specific, custom values (`Push` mode). Similarly to
tokens, SecretIDs have properties like usage-limit, TTLs and expirations.

#### Pull And Push SecretID Modes

If the SecretID used for login is fetched from an AppRole, this is operating in
Pull mode. If a "custom" SecretID is set against an AppRole by the client, it
is referred to as a Push mode. Push mode mimics the behavior of the deprecated
App-ID backend; however, in most cases Pull mode is the better approach. The
reason is that Push mode requires some other system to have knowledge of the
full set of client credentials (RoleID and SecretID) in order to create the
entry, even if these are then distributed via different paths. However, in Pull
mode, even though the RoleID must be known in order to distribute it to the
client, the SecretID can be kept confidential from all parties except for the
final authenticating client by using [Response
Wrapping](/docs/concepts/response-wrapping.html).

Push mode is available for App-ID workflow compatibility, which in some
specific cases is preferable, but in most cases Pull mode is more secure and
should be preferred.

### Further Constraints

`role_id` is a required credential at the login endpoint. AppRole pointed to by
the `role_id` will have constraints set on it. This dictates other `required`
credentials for login. The `bind_secret_id` constraint requires `secret_id` to
be presented at the login endpoint.  Going forward, this backend can support
more constraint parameters to support varied set of Apps. Some constraints will
not require a credential, but still enforce constraints for login.  For
example, `bound_cidr_list` will only allow requests coming from IP addresses
belonging to configured CIDR blocks on the AppRole.

## Comparison to Tokens

## Authentication

### Via the CLI

#### Enable AppRole authentication

```shell
$ vault auth-enable approle
```

#### Create a role

```shell
$ vault write auth/approle/role/testrole secret_id_ttl=10m token_num_uses=10 token_ttl=20m token_max_ttl=30m secret_id_num_uses=40
```

#### Fetch the RoleID of the AppRole

```shell
$ vault read auth/approle/role/testrole/role-id
```

```shell
role_id     db02de05-fa39-4855-059b-67221c5c2f63
```

#### Get a SecretID issued against the AppRole

```shell
$ vault write -f auth/approle/role/testrole/secret-id
```

```shell
secret_id               6a174c20-f6de-a53c-74d2-6018fcceff64
secret_id_accessor      c454f7e5-996e-7230-6074-6ef26b7bcf86
```


#### Login to get a Vault Token

```shell
$ vault write auth/approle/login role_id=db02de05-fa39-4855-059b-67221c5c2f63 secret_id=6a174c20-f6de-a53c-74d2-6018fcceff64
```

```shell
token           65b74ffd-842c-fd43-1386-f7d7006e520a
token_accessor  3c29bc22-5c72-11a6-f778-2bc8f48cea0e
token_duration  20m0s
token_renewable true
token_policies  [default]
```

### Via the API

#### Enable the AppRole authentication.

```javascript
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"type":"approle"}' http://127.0.0.1:8200/v1/sys/auth/approle
```

#### Create an AppRole with desired set of policies.

```javascript
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" -d '{"policies":"dev-policy,test-policy"}' http://127.0.0.1:8200/v1/auth/approle/role/testrole
```

#### Fetch the identifier of the role.

```javascript
$ curl -X GET -H "X-Vault-Token:$VAULT_TOKEN" http://127.0.0.1:8200/v1/auth/approle/role/testrole/role-id | jq .
```

```javascript
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "role_id": "988a9dfd-ea69-4a53-6cb6-9d6b86474bba"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "ef5c9b3f-e15e-0527-5457-79b4ecfe7b60"
}
```

#### Create a new secret identifier under the role.

```javascript
$ curl -X POST -H "X-Vault-Token:$VAULT_TOKEN" http://127.0.0.1:8200/v1/auth/approle/role/testrole/secret-id | jq .
```

```javascript
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "secret_id_accessor": "45946873-1d96-a9d4-678c-9229f74386a5",
    "secret_id": "37b74931-c4cd-d49a-9246-ccc62d682a25"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "c98fa1c2-7565-fd45-d9de-0b43c307f2aa"
}
```

#### Perform the login operation to fetch a new Vault token.

```javascript
$ curl -X POST \
     -d '{"role_id":"988a9dfd-ea69-4a53-6cb6-9d6b86474bba","secret_id":"37b74931-c4cd-d49a-9246-ccc62d682a25"}' \
     http://127.0.0.1:8200/v1/auth/approle/login | jq .
```

```javascript
{
  "auth": {
    "renewable": true,
    "lease_duration": 2764800,
    "metadata": {},
    "policies": [
      "default",
      "dev-policy",
      "test-policy"
    ],
    "accessor": "5d7fb475-07cb-4060-c2de-1ca3fcbf0c56",
    "client_token": "98a4c7ab-b1fe-361b-ba0b-e307aacfd587"
  },
  "warnings": null,
  "wrap_info": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "988fb8db-ce3b-0167-0ac7-1a568b902d75"
}
```

## API
### /auth/approle/role
#### LIST/GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Lists the existing AppRoles in the backend
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role` (LIST) or `/auth/approle/role?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
  None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "keys": [
          "dev",
          "prod",
          "test"
        ]
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>


### /auth/approle/role/[role_name]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Creates a new AppRole or updates an existing AppRole. This endpoint
  supports both `create` and `update` capabilities. There can be one or more
  constraints enabled on the role. It is required to have at least one of them
  enabled while creating or updating a role.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role_name</span>
        <span class="param-flags">required</span>
        Name of the AppRole.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bind_secret_id</span>
        <span class="param-flags">optional</span>
        Require `secret_id` to be presented when logging in using this AppRole.
        Defaults to 'true'.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_cidr_list</span>
        <span class="param-flags">optional</span>
        Comma-separated list of CIDR blocks; if set, specifies blocks of IP
        addresses which can perform the login operation.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Comma-separated list of policies set on tokens issued via this AppRole.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">secret_id_num_uses</span>
        <span class="param-flags">optional</span>
        Number of times any particular SecretID can be used to fetch a token
        from this AppRole, after which the SecretID will expire.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">secret_id_ttl</span>
        <span class="param-flags">optional</span>
        Duration in either an integer number of seconds (`3600`) or an integer
        time unit (`60m`) after which any SecretID expires.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">token_num_uses</span>
        <span class="param-flags">optional</span>
        Number of times issued tokens can be used.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">token_ttl</span>
        <span class="param-flags">optional</span>
        Duration in either an integer number of seconds (`3600`) or an integer
        time unit (`60m`) to set as the TTL for issued tokens and at renewal
        time.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">token_max_ttl</span>
        <span class="param-flags">optional</span>
        Duration in either an integer number of seconds (`3600`) or an integer
        time unit (`60m`) after which the issued token can no longer be
        renewed.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">period</span>
        <span class="param-flags">optional</span>
        Duration in either an integer number of seconds (`3600`) or an integer
        time unit (`60m`). If set, the token generated using this AppRole is a
        _periodic_ token; so long as it is renewed it never expires, but the
        TTL set on the token at each renewal is fixed to the value specified
        here. If this value is modified, the token will pick up the new value
        at its next renewal.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads the properties of an existing AppRole.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "token_ttl": 1200,
        "token_max_ttl": 1800,
        "secret_id_ttl": 600,
        "secret_id_num_uses": 40,
        "policies": [
          "default"
        ],
        "period": 0,
        "bind_secret_id": true,
        "bound_cidr_list": ""
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
  Deletes an existing AppRole from the backend.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/approle/role/[role_name]/role-id
#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads the RoleID of an existing AppRole.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/role-id`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "role_id": "e5a7b66e-5d08-da9c-7075-71984634b882"
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Updates the RoleID of an existing AppRole to a custom value.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/role-id`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role_id</span>
        <span class="param-flags">required</span>
        Value to be set as RoleID.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
  `204` response code.
  </dd>
</dl>



### /auth/approle/role/[role_name]/secret-id
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Generates and issues a new SecretID on an existing AppRole. Similar to
  tokens, the response will also contain a `secret_id_accessor` value which can
  be used to read the properties of the SecretID without divulging the SecretID
  itself, and also to delete the SecretID from the AppRole.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/secret-id`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">metadata</span>
        <span class="param-flags">optional</span>
        Metadata to be tied to the SecretID. This should be a JSON-formatted
        string containing the metadata in key-value pairs. This metadata will
        be set on tokens issued with this SecretID, and is logged in audit logs
        _in plaintext_.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">cidr_list</span>
        <span class="param-flags">optional</span>
Comma separated list of CIDR blocks enforcing secret IDs to be used from
specific set of IP addresses. If 'bound_cidr_list' is set on the role, then the
list of CIDR blocks listed here should be a subset of the CIDR blocks listed on
the role.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "secret_id_accessor": "84896a0c-1347-aa90-a4f6-aca8b7558780",
        "secret_id": "841771dc-11c9-bbc7-bcac-6a3945a69cd9"
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Lists the accessors of all the SecretIDs issued against the AppRole.
  This includes the accessors for "custom" SecretIDs as well.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/secret-id` (LIST) or `/auth/approle/role/[role_name]/secret-id?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
  None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "keys": [
          "ce102d2a-8253-c437-bf9a-aceed4241491",
          "a1c8dee4-b869-e68d-3520-2040c1a0849a",
          "be83b7e2-044c-7244-07e1-47560ca1c787",
          "84896a0c-1347-aa90-a4f6-aca8b7558780",
          "239b1328-6523-15e7-403a-a48038cdc45a"
        ]
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

### /auth/approle/role/[role_name]/secret-id/lookup
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads out the properties of a SecretID.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/secret-id/lookup`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">secret_id</span>
        <span class="param-flags">required</span>
Secret ID attached to the role
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "request_id": "0d25d8ec-0d16-2842-1dda-c28c25aefd4b",
      "lease_id": "",
      "lease_duration": 0,
      "renewable": false,
      "data": {
        "cidr_list": null,
        "creation_time": "2016-09-28T21:00:46.760570318-04:00",
        "expiration_time": "0001-01-01T00:00:00Z",
        "last_updated_time": "2016-09-28T21:00:46.760570318-04:00",
        "metadata": {},
        "secret_id_accessor": "b4bea6b2-0214-9f7f-33cf-e732155feadb",
        "secret_id_num_uses": 10,
        "secret_id_ttl": 0
      }
    }
    ```

  </dd>
</dl>

### /auth/approle/role/[role_name]/secret-id/destroy
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Deletes a SecretID.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/secret-id/destroy`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">secret_id</span>
        <span class="param-flags">required</span>
Secret ID attached to the role
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
  `204` response code.
  </dd>
</dl>

### /auth/approle/role/[role_name]/secret-id-accessor/lookup
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Reads out the properties of the SecretID associated with the supplied
  accessor.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/secret-id-accessor/lookup`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">secret_id_accessor</span>
        <span class="param-flags">required</span>
Accessor of the secret ID
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "request_id": "2132237e-d1b6-d298-6117-b54a2d938d00",
      "lease_id": "",
      "lease_duration": 0,
      "renewable": false,
      "data": {
        "cidr_list": null,
        "creation_time": "2016-09-28T22:09:02.834238344-04:00",
        "expiration_time": "0001-01-01T00:00:00Z",
        "last_updated_time": "2016-09-28T22:09:02.834238344-04:00",
        "metadata": {},
        "secret_id_accessor": "54ba219d-b539-ac4f-e3cf-763c02f351fb",
        "secret_id_num_uses": 10,
        "secret_id_ttl": 0
      }
    }
    ```

  </dd>
</dl>

### /auth/approle/role/[role_name]/secret-id-accessor/destroy
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Deletes the SecretID associated with the given accessor.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/secret-id-accessor/destroy`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">secret_id_accessor</span>
        <span class="param-flags">required</span>
Accessor of the secret ID
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
  `204` response code.
  </dd>
</dl>


### /auth/approle/role/[role_name]/custom-secret-id
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Assigns a "custom" SecretID against an existing AppRole. This is used in the
  "Push" model of operation.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/custom-secret-id`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">secret_id</span>
        <span class="param-flags">required</span>
        SecretID to be attached to the Role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">metadata</span>
        <span class="param-flags">optional</span>
        Metadata to be tied to the SecretID. This should be a JSON-formatted
        string containing the metadata in key-value pairs. This metadata will
        be set on tokens issued with this SecretID, and is logged in audit logs
        _in plaintext_.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">cidr_list</span>
        <span class="param-flags">optional</span>
Comma separated list of CIDR blocks enforcing secret IDs to be used from
specific set of IP addresses. If 'bound_cidr_list' is set on the role, then the
list of CIDR blocks listed here should be a subset of the CIDR blocks listed on
the role.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": null,
      "warnings": null,
      "wrap_info": null,
      "data": {
        "secret_id_accessor": "a109dc4a-1fd3-6df6-feda-0ca28b2d4a81",
        "secret_id": "testsecretid"
      },
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>


### /auth/approle/login
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
  Issues a Vault token based on the presented credentials. `role_id` is always
  required; if `bind_secret_id` is enabled (the default) on the AppRole,
  `secret_id` is required too. Any other bound authentication values on the
  AppRole (such as client IP CIDR) are also evaluated.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/login`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role_id</span>
        <span class="param-flags">required</span>
        RoleID of the AppRole.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">secret_id</span>
        <span class="param-flags">required when `bind_secret_id` is enabled</span>
        SecretID belonging to AppRole.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "auth": {
        "renewable": true,
        "lease_duration": 1200,
        "metadata": null,
        "policies": [
          "default"
        ],
        "accessor": "fd6c9a00-d2dc-3b11-0be5-af7ae0e1d374",
        "client_token": "5b1a0318-679c-9c45-e5c6-d1b9a9035d49"
      },
      "warnings": null,
      "wrap_info": null,
      "data": null,
      "lease_duration": 0,
      "renewable": false,
      "lease_id": ""
    }
    ```

  </dd>
</dl>

### /auth/approle/role/[role_name]/policies
### /auth/approle/role/[role_name]/secret-id-num-uses
### /auth/approle/role/[role_name]/secret-id-ttl
### /auth/approle/role/[role_name]/token-ttl
### /auth/approle/role/[role_name]/token-max-ttl
### /auth/approle/role/[role_name]/bind-secret-id
### /auth/approle/role/[role_name]/bound-cidr-list
### /auth/approle/role/[role_name]/period
#### POST/GET/DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
  Updates the respective property in the existing AppRole. All of these
  parameters of the AppRole can be updated using the `/auth/approle/role/[role_name]`
  endpoint directly. The endpoints for each field is provided separately
  to be able to delegate specific endpoints using Vault's ACL system.
  </dd>

  <dt>Method</dt>
  <dd>POST/GET/DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/approle/role/[role_name]/[field_name]`</dd>

  <dt>Parameters</dt>
  <dd>
  Refer to `/auth/approle/role/[role_name]` endpoint.
  </dd>

  <dt>Returns</dt>
  <dd>
  Refer to `/auth/approle/role/[role_name]` endpoint.
  </dd>
</dl>
