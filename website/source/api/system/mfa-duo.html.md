---
layout: "api"
page_title: "/sys/mfa/method/duo - HTTP API"
sidebar_current: "docs-http-system-mfa-duo"
description: |-
  The '/sys/mfa/method/duo' endpoint focuses on managing Duo MFA behaviors in Vault Enterprise.
---

## Configure Duo MFA Method

This endpoint defines a MFA method of type Duo.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `POST`   | `/sys/mfa/method/duo/:name`   | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

- `mount_accessor` `(string: <required>)` - The mount to tie this method to for use in automatic mappings. The mapping will use the Name field of Aliases associated with this mount as the username in the mapping.

- `username_format` `(string)` - A format string for mapping Identity names to MFA method names. Values to substitute should be placed in `{{}}`. For example, `"{{alias.name}}@example.com"`. If blank, the Alias's Name field will be used as-is. Currently-supported mappings:
  - alias.name: The name returned by the mount configured via the `mount_accessor` parameter
  - entity.name: The name configured for the Entity
  - alias.metadata.`<key>`: The value of the Alias's metadata parameter
  - entity.metadata.`<key>`: The value of the Entity's metadata paramater

- `secret_key` `(string)` - Secret key for Duo.

- `integration_key` `(string)` - Integration key for Duo.

- `api_hostname` `(string)` - API hostname for Duo.

- `push_info` `(string)` - Push information for Duo.

### Sample Payload

```json
{
  "mount_accessor": "auth_userpass_1793464a",
  "secret_key": "BIACEUEAXI20BNWTEYXT",
  "integration_key":"8C7THtrIigh2rPZQMbguugt8IUftWhMRCOBzbuyz",
  "api_hostname":"api-2b5c39f5.duosecurity.com"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/mfa/method/duo/my_duo
```

## Read Duo MFA Method

This endpoint queries the MFA configuration of Duo type for a given method
name.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `GET`    | `/sys/mfa/method/duo/:name`   | `200 application/json`   |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://vault.rocks/v1/sys/mfa/method/duo/my_duo

```

### Sample Response

```json
{
        "data": {
                "api_hostname": "api-2b5c39f5.duosecurity.com",
                "id": "0ad21b78-e9bb-64fa-88b8-1e38db217bde",
                "integration_key": "BIACEUEAXI20BNWTEYXT",
                "mount_accessor": "auth_userpass_1793464a",
                "name": "my_duo",
                "pushinfo": "",
                "secret_key": "8C7THtrIigh2rPZQMbguugt8IUftWhMRCOBzbuyz",
                "type": "duo",
                "username_format": ""
        }
}
```
## Delete Duo MFA Method

This endpoint deletes a Duo MFA method.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `DELETE` | `/sys/mfa/method/duo/:name`   | `204 (empty body)`       |


### Parameters

- `name` `(string: <required>)` - Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/mfa/method/duo/my_duo

```
