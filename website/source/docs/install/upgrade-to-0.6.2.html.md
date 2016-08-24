---
layout: "install"
page_title: "Upgrading to Vault 0.6.2"
sidebar_current: "docs-install-upgrade-to-0.6.2"
description: |-
  Learn how to upgrade to Vault 0.6.2
---

# Overview

This page contains the list of breaking changes for Vault 0.6.2. Please read it
carefully.

## Paths in auth/token May No Longer Contain a Token; Change to api.RevokeSelf

Several of the lookup, renew, and revoke operations on tokens via `auth/token`
allowed the token to be specified as part of the URL. Although never
recommended, we have now removed this capability altogether, as accidental
misuse can lead to problematic outcomes (for instance, request paths are logged
in cleartext in audit logs).

As part of this change, the `auth/token/lookup` endpoint no longer accepts
`GET` requests; all requests must be `PUT` or `POST` and contain a `token`
parameter in the JSON body data.

As part of this change, we have also removed the `token` parameter from the
`RevokeSelf` call in the Go API. This parameter was unused and confusing.
