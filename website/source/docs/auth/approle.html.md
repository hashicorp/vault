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
credentials required for successful login depend upon the constraints set on
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

The AppRole authentication backend has a full HTTP API. Please see the
[AppRole API](/api/auth/approle/index.html) for more
details.