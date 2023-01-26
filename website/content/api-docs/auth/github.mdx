---
layout: api
page_title: GitHub - Auth Methods - HTTP API
description: This is the API documentation for the Vault GitHub auth method.
---

# GitHub Auth Method (API)

This is the API documentation for the Vault GitHub auth method. For
general information about the usage and operation of the GitHub method, please
see the [Vault GitHub method documentation](/vault/docs/auth/github).

This documentation assumes the GitHub method is enabled at the `/auth/github`
path in Vault. Since it is possible to enable auth methods at any location,
please update your API calls accordingly.

## Configure Method

Configures the connection parameters for GitHub. This path honors the
distinction between the `create` and `update` capabilities inside ACL policies.

| Method | Path                  |
| :----- | :-------------------- |
| `POST` | `/auth/github/config` |

### Parameters

- `organization` `(string: <required>)` - The organization users must be part
  of.
- `organization_id` `(int: 0)` - The ID of the organization users must be part
  of. Vault will attempt to fetch and set this value if it is not provided.
- `base_url` `(string: "")` - The API endpoint to use. Useful if you are running
  GitHub Enterprise or an API-compatible authentication server.

@include 'tokenfields.mdx'

### Sample Payload

```json
{
  "organization": "acme-org"
}
```

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/github/config
```

## Read Configuration

Reads the GitHub configuration.

| Method | Path                  |
| :----- | :-------------------- |
| `GET`  | `/auth/github/config` |

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/github/config
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "organization": "acme-org",
    "base_url": "",
    "ttl": "",
    "max_ttl": ""
  },
  "warnings": null
}
```

## Map GitHub Teams

Map a list of policies to a team that exists in the configured GitHub organization.

| Method | Path                                |
| :----- | :---------------------------------- |
| `POST` | `/auth/github/map/teams/:team_name` |

### Parameters

- `team_name` `(string)` - GitHub team name in "slugified" format
- `value` `(string)` - Comma separated list of policies to assign

### Sample Payload

```json
{
  "value": "dev-policy"
}
```

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/github/map/teams/dev
```

## Read Team Mapping

Reads the GitHub team policy mapping.

| Method | Path                                |
| :----- | :---------------------------------- |
| `GET`  | `/auth/github/map/teams/:team_name` |

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/github/map/teams/dev
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "key": "dev",
    "value": "dev-policy"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

## Map GitHub Users

Map a list of policies to a specific GitHub user exists in the configured
organization.

| Method | Path                                |
| :----- | :---------------------------------- |
| `POST` | `/auth/github/map/users/:user_name` |

### Parameters

- `user_name` `(string)` - GitHub user name
- `value` `(string)` - Comma separated list of policies to assign

### Sample Payload

```json
{
  "value": "sethvargo-policy"
}
```

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    http://127.0.0.1:8200/v1/auth/github/map/users/sethvargo
```

The user with username `sethvargo` will be assigned the `sethvargo-policy`
policy **in addition to** any team policies.

## Read User Mapping

Reads the GitHub user policy mapping.

| Method | Path                                |
| :----- | :---------------------------------- |
| `GET`  | `/auth/github/map/users/:user_name` |

### Sample Request

```shell-session
$ curl \
    --header "X-Vault-Token: ..." \
    http://127.0.0.1:8200/v1/auth/github/map/users/sethvargo
```

### Sample Response

```json
{
  "request_id": "764b6f88-efba-51bd-ed62-cf1c9e80e37a",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "key": "sethvargo",
    "value": "sethvargo-policy"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
```

## Login

Login using GitHub access token.

| Method | Path                 |
| :----- | :------------------- |
| `POST` | `/auth/github/login` |

### Parameters

- `token` `(string: <required>)` - GitHub personal API token.

### Sample Payload

```json
{
  "token": "ABC123..."
}
```

### Sample Request

```shell-session
$ curl \
    --request POST \
    http://127.0.0.1:8200/v1/auth/github/login
```

### Sample Response

```javascript
{
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "warnings": null,
  "auth": {
    "client_token": "64d2a8f2-2a2f-5688-102b-e6088b76e344",
    "accessor": "18bb8f89-826a-56ee-c65b-1736dc5ea27d",
    "policies": ["default"],
    "metadata": {
      "username": "fred",
      "org": "acme-org"
    },
  },
  "lease_duration": 7200,
  "renewable": true
}
```
