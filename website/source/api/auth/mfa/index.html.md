---
layout: "api"
page_title: "MFA Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-mfa"
description: |-
  This is the API documentation for the Vault MFA authentication backend.
---

# MFA Auth Backend HTTP API

This is the API documentation for the Vault MFA authentication backend. For
general information about the usage and operation of the AppRole backend, please
see the [Vault MFA backend documentation](/docs/auth/mfa.html).

This documentation assumes the MFA backend is mounted at the `/auth/mfa`
path in Vault. Since it is possible to mount auth backends at any location,
please update your API calls accordingly.