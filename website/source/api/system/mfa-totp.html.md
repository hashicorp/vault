---
layout: "api"
page_title: "/sys/mfa/method/totp - HTTP API"
sidebar_current: "docs-http-system-mfa-totp"
description: |-
  The '/sys/mfa/method/totp' endpoint focuses on managing TOTP MFA behaviors in Vault Enterprise.
---

## Configure TOTP MFA Method

This endpoint defines a MFA method of type TOTP.

| Method   | Path                           | Produces               |
| :------- | :----------------------------- | :--------------------- |
| `POST`   | `/sys/mfa/method/totp/:name`   | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

- `issuer` `(string: <required>)` - The name of the key's issuing organization.

- `period` `(int or duration format string: 30)` - The length of time used to generate a counter for the TOTP token calculation.

- `key_size` `(int: 20)` – Specifies the size in bytes of the generated key.

- `qr_size` `(int: 200)` - The pixel size of the generated square QR code.

- `algorithm` `(string: "SHA1")` – Specifies the hashing algorithm used to generate the TOTP code. Options include "SHA1", "SHA256" and "SHA512".

- `digits` `(int: 6)` - The number of digits in the generated TOTP token. This value can either be 6 or 8.

- `skew` `(int: 1)` - The number of delay periods that are allowed when validating a TOTP token. This value can either be 0 or 1.


### Sample Payload

```json
{
  "issuer": "vault"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json \
    https://vault.rocks/v1/sys/mfa/method/totp/my_totp
```

## Read TOTP MFA Method

This endpoint queries the MFA configuration of TOTP type for a given method
name.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `GET`    | `/sys/mfa/method/totp/:name`   | `200 application/json`   |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://vault.rocks/v1/sys/mfa/method/totp/my_totp

```

### Sample Response

```json
{
        "data": {
                "algorithm": "SHA1",
                "digits": 6,
                "id": "865587ba-6229-7f2a-6da0-609d5370af70",
                "issuer": "vault",
                "key_size": 20,
                "name": "my_totp",
                "period": 30,
                "qr_size": 200,
                "skew": 1,
                "type": "totp"
        }
}
```

## Delete TOTP MFA Method

This endpoint deletes a TOTP MFA method.

| Method   | Path                           | Produces                 |
| :------- | :----------------------------- | :----------------------- |
| `DELETE` | `/sys/mfa/method/totp/:name`   | `204 (empty body)`       |


### Parameters

- `name` `(string: <required>)` - Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request DELETE \
    https://vault.rocks/v1/sys/mfa/method/totp/my_totp

```

## Generate a TOTP MFA Secret

This endpoint generates an MFA secret in the entity of the calling token, if it
doesn't exist already, using the configuration stored under the given MFA
method name.

| Method   | Path                                  | Produces                 |
| :------- | :------------------------------------ | :----------------------- |
| `GET`    | `/sys/mfa/method/totp/:name/generate` | `200 application/json`   |

### Parameters

- `name` `(string: <required>)` - Name of the MFA method.

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request GET \
    https://vault.rocks/v1/sys/mfa/method/totp/my_totp/generate
```

### Sample Response

```json
{
  "data": {
    "barcode": "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 rkJggg==",
    "url": "otpauth://totp/vault:4746fb81-028c-cd4e-026b-7dd18fe4c2f4?algorithm=SHA1&digits=6&issuer=vault&period=30&secret=XVE7TOZWJVEWQOATOD7 U53IEAJG72Z2I"
  }
}
```

## Administratively Generate a TOTP MFA Secret

This endpoint can be used to generate a TOTP MFA secret. Unlike the `generate`
API which stores the generated secret on the entity ID of the calling token,
the `admin-generate` API stores the generated secret on the given entity ID.

| Method   | Path                                         | Produces                 |
| :------- | :------------------------------------------- | :----------------------- |
| `POST`   | `/sys/mfa/method/totp/:name/admin-generate`  | `200 application/json`   |

### Parameters

- `name` `(string: <required>)` - Name of the MFA method.

- `entity_id` `(string: <required>)` - Entity ID on which the generated secret
  needs to get stored.

### Sample Payload

```json
{
  "entity_id":"4746fb81-028c-cd4e-026b-7dd18fe4c2f4"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json
    https://vault.rocks/v1/sys/mfa/method/totp/my_totp/admin-generate
```

### Sample Response

```json
{
  "data": {
    "barcode": "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",
    "url": "otpauth://totp/vault:4746fb81-028c-cd4e-026b-7dd18fe4c2f4?algorithm=SHA1&digits=6&issuer=vault&period=30&secret=6HQ4RZ7GM6MMLRKVDCI23LXNZF7UDZ2U"
  }
}
```

### Administratively Destroy TOTP MFA Secret

This endpoint deletes a TOTP MFA secret from the given entity ID.

Note that in order to overwrite a secret on the entity, it is required to
explicitly delete the secret first. This API can be used to delete the secret
and the `generate` or `admin-generate` APIs should be used to regenerate a new
secret.

| Method   | Path                                    | Produces               |
| :------- | :-------------------------------------- | :--------------------- |
| `POST`   | `/sys/mfa/method/:name/admin-destroy`   | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Name of the MFA method.

- `entity_id` `(string: <required>)` - Entity ID from which the MFA secret
  should be removed.

### Sample Payload

```json
{
  "entity_id": "4746fb81-028c-cd4e-026b-7dd18fe4c2f4"
}
```

### Sample Request

```
$ curl \
    --header "X-Vault-Token: ..." \
    --request POST \
    --data @payload.json
    https://vault.rocks/v1/sys/mfa/method/totp/my_totp/admin-destroy
```
