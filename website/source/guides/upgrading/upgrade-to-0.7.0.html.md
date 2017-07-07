---
layout: "guides"
page_title: "Upgrading to Vault 0.7.0 - Guides"
sidebar_current: "guides-upgrading-to-0.7.0"
description: |-
  This page contains the list of deprecations and important or breaking changes
  for Vault 0.7.0. Please read it carefully.
---

# Overview

This page contains the list of deprecations and important or breaking changes
for Vault 0.7.0 compared to the most recent release. Please read it carefully.

## List Operations Always Use Trailing Slash

 Any list operation, whether via the `GET` or `LIST` HTTP verb, will now
 internally canonicalize the path to have a trailing slash. This makes policy
 writing more predictable, as it means clients will no longer work or fail
 based on which client they're using or which HTTP verb they're using. However,
 it also means that policies allowing `list` capability must be carefully
 checked to ensure that they contain a trailing slash; some policies may need
 to be split into multiple stanzas to accommodate.

## PKI Defaults to Unleased Certificates

When issuing certificates from the PKI backend, by default, no leases will be
issued. If you want to manually revoke a certificate, its serial number can be
used with the `pki/revoke` endpoint. Issuing leases is still possible by
enabling the `generate_lease` toggle in PKI role entries (this will default to
`true` for upgrades, to keep existing behavior), which will allow using lease
IDs to revoke certificates. For installations issuing large numbers of
certificates (tens to hundreds of thousands, or millions), this will
significantly improve Vault startup time since leases associated with these
certificates will not have to be loaded; however note that it also means that
revocation of a token used to issue certificates will no longer add these
certificates to a CRL. If this behavior is desired or needed, consider keeping
leases enabled and ensuring lifetimes are reasonable, and issue long-lived
certificates via a different role with leases disabled.
