---
layout: "api"
page_title: "AppRole Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-approle"
description: |-
  This is the API documentation for the Vault AppRole authentication backend.
---

# AppRole Auth Backend HTTP API

This is the API documentation for the Vault AppRole authentication backend. For
general information about the usage and operation of the AppRole backend, please
see the [Vault AppRole backend documentation](/docs/auth/approle.html).

This documentation assumes the AppRole backend is mounted at the `/auth/approle`
path in Vault. Since it is possible to mount auth backends at any location,
please update your API calls accordingly.
