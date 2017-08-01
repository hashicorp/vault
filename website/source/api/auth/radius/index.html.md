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