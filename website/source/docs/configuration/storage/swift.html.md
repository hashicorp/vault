---
layout: "docs"
page_title: "Swift - Storage Backends - Configuration"
sidebar_current: "docs-configuration-storage-swift"
description: |-
  The Swift storage backend is used to persist Vault's data in an OpenStack
  Swift Container.
---

# Swift Storage Backend

The Swift storage backend is used to persist Vault's data in an
[OpenStack Swift Container][swift].


- **No High Availability** – the Swift storage backend does not support high
  availability.

- **Community Supported** – the Swift storage backend is supported by the
  community. While it has undergone review by HashiCorp employees, they may not
  be as knowledgeable about the technology. If you encounter problems with them,
  you may be referred to the original author.

```hcl
storage "swift" {
  auth_url  = "https://..."
  username  = "admin"
  password  = "secret123!"
  container = "my-storage-container"
}
```

## `swift` Parameters

- `auth_url` `(string: <required>)` – Specifies the OpenStack authentication
  endpoint. Currently only v1.0 authentication endpoints are supported. This can
  also be provided via the environment variable `OS_AUTH_URL`.

- `container` `(string: <required>)` – Specifies the name of the Swift
  container. This can also be provided via the environment variable
  `OS_CONTAINER`.

- `max_parallel` `(string: "128")` – The maximum number of concurrent requests.

- `password` `(string: <required>)` – Specifies the OpenStack password. This can
  also be provided via the environment variable `OS_PASSWORD`.

- `tenant` `(string: "")` – Specifies the name of the tenant. If left blank,
  this will default to the default tenant of the username. This can also be
  provided via the environment variable `OS_TENANT_NAME`.

- `username` `(string: <required>)` – Specifies the OpenStack account/username.
  This can also be provided via the environment variable `OS_USERNAME`.

## `swift` Examples

### Default Example

This example shows a default configuration for Swift.

```hcl
storage "swift" {
  auth_url  = "https://os.initernal/v1/auth"
  container = "container-239"

  username  = "user1234"
  password  = "pass5678"
}
```

[swift]: http://docs.openstack.org/developer/swift/
