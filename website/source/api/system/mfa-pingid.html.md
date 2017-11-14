---
layout: "api"
page_title: "/sys/mfa/method/pingid - HTTP API"
sidebar_current: "docs-http-system-mfa-pingid"
description: |-
  The '/sys/mfa/method/pingid' endpoint focuses on managing PingID MFA behaviors in Vault Enterprise.
---

## Configure PingID MFA Method

This endpoint defines a MFA method of type PingID.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `POST`   | `/sys/mfa/method/pingid/:name`   | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

- `mount_accessor` `(string: <required>)` - The mount to tie this method to for use in automatic mappings. The mapping will use the Name field of Aliases associated with this mount as the username in the mapping.

- `username_format` `(string)` - A format string for mapping Identity names to MFA method names. Values to substitute should be placed in `{{}}`. For example, `"{{alias.name}}@example.com"`. If blank, the Alias's Name field will be used as-is. Currently-supported mappings:
  - alias.name: The name returned by the mount configured via the `mount_accessor` parameter
  - entity.name: The name configured for the Entity
  - alias.metadata.`<key>`: The value of the Alias's metadata parameter
  - entity.metadata.`<key>`: The value of the Entity's metadata paramater

- `settings_file_base64` `(string)` - A base64-encoded third-party settings file retrieved from PingID's configuration page.

### Sample Payload

```json
{
  "mount_accessor": "auth_userpass_1793464a",
  "settings_file_base64": "AA8owj3..."
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/mfa/method/pingid/ping
```

## Read PingiD MFA Method

This endpoint queries the MFA configuration of PingID type for a given method
name.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `GET`    | `/sys/mfa/method/pingid/:name`   | `200 application/json`   |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://vault.rocks/v1/sys/mfa/method/pingid/ping

```

### Sample Response

```json
{
        "data": {
                "use_signature": true,
                "idp_url": "https://idpxnyl3m.pingidentity.com/pingid",
                "admin_url": "https://idpxnyl3m.pingidentity.com/pingid",
                "authenticator_url": "https://authenticator.pingone.com/pingid/ppm",
                "mount_accessor": "auth_userpass_1793464a",
                "name": "ping",
                "org_alias": "181459b0-9fb1-4938-8c86...",
                "type": "pingid",
                "username_format": ""
        }
}
```
## Delete PingID MFA Method

This endpoint deletes a PingID MFA method.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `DELETE` | `/sys/mfa/method/pingid/:name`   | `204 (empty body)`       |


### Parameters

- `name` `(string: <required>)` - Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/mfa/method/pingid/ping

```
