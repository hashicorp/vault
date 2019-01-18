---
layout: "api"
page_title: "/sys/internal/specs/openapi - HTTP API"
sidebar_title: "<code>/sys/internal/specs/openapi</code>"
sidebar_current: "api-http-system-internal-specs-openapi"
description: |-
  The `/sys/internal/specs/openapi` endpoint is used to generate an OpenAPI document of the mounted backends.
---

# `/sys/internal/specs/openapi`

The `/sys/internal/specs/openapi` endpoint is used to generate an OpenAPI document of the mounted backends.
The response conforms to the [OpenAPI V3 specification](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md),
with path names matching the mount names used by the Vault server (i.e. customizations with `-path` will be reflected).
The set of included paths is based on the permissions of the request token.

The response may include Vault-specific [extensions](https://github.com/oai/openapi-specification/blob/master/versions/3.0.2.md#specification-extensions). Three are currently defined:

- `x-vault-sudo` - Endpoint requires [sudo](https://www.vaultproject.io/docs/concepts/policies.html#sudo) privileges.
- `x-vault-unauthenticated` - Endpoint is unauthenticated.
- `x-vault-create-supported` - Endpoint allows creation of new items, in addition to updating existing items.

Basic documentation will be generated for all paths, but a newer path definition structure now allows for
more detailed documentation to be added. At this time the `/sys` endpoints have been updated to use the new
structure, and other endpoints will be modified incrementally.

## Get OpenAPI Document

This endpoint returns a single OpenAPI document describing all paths visible to the requester.

| Method |           Path                |        Produces        |
| :----- | :------------------------     | :--------------------- |
| `GET`  | `/sys/internal/specs/openapi` | `200 application/json` |


### Sample Request

```
$ curl http://127.0.0.1:8200/v1/sys/internal/specs/openapi
```

### Sample Response

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "HashiCorp Vault API",
    "description": "HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.",
    "version": "1.0.0",
    "license": {
      "name": "Mozilla Public License 2.0",
      "url": "https://www.mozilla.org/en-US/MPL/2.0"
    }
  },
  "paths": {
    "/auth/token/create": {
      "description": "The token create path is used to create new tokens.",
      "post": {
        "summary": "The token create path is used to create new tokens.",
        "tags": [
          "auth"
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    },
    ...
```

