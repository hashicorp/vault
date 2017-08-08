---
layout: "api"
page_title: "/sys/mfa/method/okta - HTTP API"
sidebar_current: "docs-http-system-mfa-okta"
description: |-
  The '/sys/mfa/method/okta' endpoint focuses on managing Okta MFA behaviors in Vault Enterprise.
---

## Configure Okta MFA Method

This endpoint defines a MFA method of type Okta.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `POST`   | `/sys/mfa/method/okta/:name`   | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

- `mount_accessor` `(string: <required>)` - The mount to tie this method to for use in automatic mappings. The mapping will use the Name field of Personas associated with this mount as the username in the mapping.

- `username_format` `(string)` - A format string for mapping Identity names to MFA method names. Values to substitute should be placed in `{{}}`. For example, `"{{persona.name}}@example.com"`. If blank, the Persona's Name field will be used as-is. Currently-supported mappings:
  - persona.name: The name returned by the mount configured via the `mount_accessor` parameter

- `org_name` `(string)` - Name of the organization to be used in the Okta API.

- `api_token` `(string)` - Okta API key.

- `production` `(string)` - If set, production API URL prefix will be used to communicate with Okta and if not set, a preview production API URL prefix will be used. Defaults to true.

### Sample Payload

```json
{
  "mount_accessor": "auth_userpass_1793464a",
  "org_name": "dev-262778",
  "api_token": "0081u7KrReNkzmABZJAP2oDyIXccveqx9vIOEyCZDC"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/mfa/method/okta/my_okta
```

## Read Okta MFA Method

This endpoint queries the MFA configuration of Okta type for a given method
name.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `GET`    | `/sys/mfa/method/okta/:name`   | `200 application/json`   |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://vault.rocks/v1/sys/mfa/method/okta/my_okta

```

### Sample Response

```json
{
        "data": {
                "api_token": "0081u7KrReNkzmABZJAP2oDyIXccveqx9vIOEyCZDC",
                "id": "e39f08a1-a42d-143d-5b87-15c61d89c15a",
                "mount_accessor": "auth_userpass_1793464a",
                "name": "my_okta",
                "org_name": "dev-262778",
                "production": true,
                "type": "okta",
                "username_format": ""
        }
}
```
## Delete Okta MFA Method

This endpoint deletes a Okta MFA method.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `DELETE` | `/sys/mfa/method/okta/:name`   | `204 (empty body)`       |


### Parameters

- `name` `(string: <required>)` - Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/mfa/method/okta/my_okta

```
