---
layout: "docs"
page_title: "User Interface - Configuration"
sidebar_current: "docs-configuration-ui"
description: |-
  Vault features a web based user interface (UI) for interacting with popular features. You can use the UI  to easily create, read, update, and delete secrets, authenticate, unseal, and more.
---

# Vault UI

Vault features a web based user interface (UI) for interacting with popular features. You can use the UI  to easily create, read, update, and delete secrets, authenticate, unseal, and more.

## Activating the Vault UI

The Vault UI is not activated by default. To activate the UI, set the `ui` configuration option in the Vault server configuration. Vault clients do not need to set this option, since they do not be serve the UI.

Here is an example configuration snippet:

```hcl
ui = true

listener "tcp" {
  address = "10.0.1.35:8200"
}

storage "file" {
  path = "/tmp/vault"
}
```

For more information on configuration file options please see [Vault Configuration](/docs/configuration/index.html).

## Accessing the Vault UI

The UI listens on the same port as the Vault API listener. As such, you must configure at least one `listener` stanza to expose the UI. Building on the example above, the listener is configured and some comments describing an alternative configuration are shown as a reference:

```hcl
listener "tcp" {
  address = "10.0.1.35:8200"

  # If bound to localhost, the Vault UI is only
  # accessible from the local machine!
  # address = "127.0.0.1:8200"
}
```

In this case, the UI is accessible the following URL from any machine on the subnet (provided no network firewalls are in place):

```text
https://10.0.1.35:8200/ui
```

It is also accessible at any DNS entry that resolves to that IP address, such as the Consul service address (if using Consul) as well:

```text
https://vault.service.consul:8200/ui
```

### Note on TLS

When using TLS (recommended), the certificate must be valid for all DNS entries you will be accessing the Vault UI on, and any IP addresses on the Subject Alternate Name. If you are running Vault with a self-signed certificate, any browsers that access the Vault UI will need to have the root CA installed. Failure to do so may result in the browser displaying a warning that the site is "untrusted". It is highly recommended that client browsers accessing the Vault UI install the proper CA root certificate into the OS trust store for validation to reduce the chance of a MITM attack.
