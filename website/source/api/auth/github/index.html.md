---
layout: "api"
page_title: "Github Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-github"
description: |-
  This is the API documentation for the Vault Github authentication backend.
---

# Github Auth Backend HTTP API

This is the API documentation for the Vault Github authentication backend. For
general information about the usage and operation of the Github backend, please
see the [Vault Github backend documentation](/docs/auth/github.html).

This documentation assumes the Github backend is mounted at the `/auth/github`
path in Vault. Since it is possible to mount auth backends at any location,
please update your API calls accordingly.

## Configure Backend

Configures the connection parameters for Github. This path honors the 
distinction between the `create` and `update` capabilities inside ACL policies.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/github/config`          | `204 (empty body)`     |

### Parameters

- `organization` `(string: <required>)` - The organization users must be part 
  of.
- `base_url` `(string: "")` - The API endpoint to use. Useful if you are running
  GitHub Enterprise or an API-compatible authentication server.
- `ttl` `(string: "")` - Duration after which authentication will be expired.
- `max_ttl` `(string: "")` - Maximum duration after which authentication will 
  be expired.

### Sample Payload

```json
{
  "organization": "acme-org"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/github/config
```

## Read Configuration

Reads the Github configuration.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`    | `/auth/github/config`        | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/github/config
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

## Login

Login using GitHub access token.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/github/login`         | `200 application/json` |

### Parameters

- `token` `(string: <required>)` - GitHub personal API token.

### Sample Payload

```json
{
  "token": "ABC123..."
}
```

### Sample Request

```
$ curl \
    --request POST \
    https://vault.rocks/v1/auth/github/login
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
