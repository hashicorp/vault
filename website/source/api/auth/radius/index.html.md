---
layout: "api"
page_title: "RADIUS Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-radius"
description: |-
  This is the API documentation for the Vault RADIUS authentication backend.
---

# RADIUS Auth Backend HTTP API

This is the API documentation for the Vault RADIUS authentication backend. For
general information about the usage and operation of the RADIUS backend, please
see the [Vault RADIUS backend documentation](/docs/auth/radius.html).

This documentation assumes the RADIUS backend is mounted at the `/auth/radius`
path in Vault. Since it is possible to mount auth backends at any location,
please update your API calls accordingly.

## Configure RADIUS

Configures the connection parameters and shared secret used to communicate with 
RADIUS.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/radius/config`        | `204 (empty body)`     |

### Parameters

- `host` `(string: <required>)` - The RADIUS server to connect to. Examples: 
  `radius.myorg.com`, `127.0.0.1`
- `port` `(integer: 1812)` - The UDP port where the RADIUS server is listening
   on. Defaults is 1812.
- `secret` `(string: <required>)` - The RADIUS shared secret.
- `unregistered_user_policies` `(string: "")` - A comma-separated list of 
  policies to be granted to unregistered users.
- `dial_timeout` `(integer: 10)` - Number of second to wait for a backend 
  connection before timing out. Default is 10.
- `nas_port` `(integer: 10)` - The NAS-Port attribute of the RADIUS request. 
  Defaults is 10.

### Sample Payload

```json
{
  "host": "radius.myorg.com",
  "port": 1812,
  "secret": "mySecret"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/radius/config
```

## Register User

Registers a new user and maps a set of policies to it.  This path honors the 
distinction between the `create` and `update` capabilities inside ACL policies.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/radius/users/:username` | `204 (empty body)`   |

### Parameters

- `username` `(string: <required>)` - Username for this user.
- `policies` `(string: "")` - Comma-separated list of policies.  If set to 
  empty string, only the `default` policy will be applicable to the user.

```json
{
  "policies": "dev,prod",
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/radius/users/test-user
```

## Read User

Reads the properties of an existing username.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `GET`   | `/auth/radius/users/:username` | `200 application/json`   |

### Parameters

- `username` `(string: <required>)` - Username for this user.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    https://vault.rocks/v1/auth/radius/users/test-user
```

### Sample Response

```json
{
  "request_id": "812229d7-a82e-0b20-c35b-81ce8c1b9fa6",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "policies": "default,dev"
  },
  "warnings": null
}
```

## Delete User

Deletes an existing username from the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE`   | `/auth/radius/users/:username` | `204 (empty body)`   |

### Parameters

- `username` `(string: <required>)` - Username for this user.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/auth/radius/users/test-user
```

## List Users

List the users registered with the backend.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/auth/radius/users`         | `200 application/json` |
| `GET`   | `/auth/radius/users?list=true` | `200 application/json` |

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request LIST \
    https://vault.rocks/v1/auth/radius/users
```

### Sample Response

```json
{
  "auth": null,
  "warnings": null,
  "wrap_info": null,
  "data": {
    "keys": [
      "devuser",
	    "produser"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## Login

Login with the username and password.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `POST`   | `/auth/radius/login`         | `200 application/json` |
| `POST`   | `/auth/radius/login/:username` | `200 application/json` |

### Parameters

- `username` `(string: <required>)` - Username for this user.
- `password` `(string: <required>)` - Password for the autheticating user.

### Sample Payload

```json
{
  "password": "Password!"
}
```

### Sample Request

```
$ curl \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/auth/radius/login/test-user
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
      "username": "vishal"
    },
  },
  "lease_duration": 7200,
  "renewable": true
}
 ```
