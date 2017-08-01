---
layout: "api"
page_title: "TLS Certificate Auth Backend - HTTP API"
sidebar_current: "docs-http-auth-cert"
description: |-
  This is the API documentation for the Vault TLS Certificate authentication
  backend.
---

# TLS Certificate Auth Backend HTTP API

This is the API documentation for the Vault TLS Certificate authentication 
backend. For general information about the usage and operation of the TLS
Certificate backend, please see the [Vault TLS Certificate backend documentation](/docs/auth/cert.html).

This documentation assumes the TLS Certificate backend is mounted at the
`/auth/cert` path in Vault. Since it is possible to mount auth backends at any
location, please update your API calls accordingly.