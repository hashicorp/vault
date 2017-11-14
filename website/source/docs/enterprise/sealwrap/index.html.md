---
layout: "docs"
page_title: "Vault Enterprise Seal Wrap"
sidebar_current: "docs-vault-enterprise-sealwrap"
description: |-
  Vault Enterprise features a mechanism to wrap values with an extra layer of
  encryption for supporting seals
---

# Seal Wrap

Vault Enterprise features a mechanism to wrap values with an extra layer of
encryption for supporting [seals](/docs/configuration/seal.html). This adds an
extra layer of protection and is useful in some compliance and regulatory
environments, including FIPS 140-2 environments.

To use this feature, you must have an active or trial license for Vault
Enterprise (HSMs) or Vault Pro (AWS KMS). To start a trial, contact [HashiCorp
sales](mailto:sales@hashicorp.com).

## FIPS 140-2 Compliance

Vault's Seal Wrap feature has been evaluated by Leidos for compliance with
FIPS 140-2 requirements. When used with a FIPS 140-2-compliant HSM, Vault will
store Critical Security Parameters (CSPs) in a manner that is compliant with
KeyStorage and KeyTransit requirements. This is on by default for many parts of
Vault and opt-in for each individual mount; see the Activating Seal Wrapping
section below for details.

[Download the current certification letter](/docs/enterprise/sealwrap/Vault_Compliance_Letter_signed.pdf)

### Updates Since The Latest FIPS Compliance Audit

The following are values that take advantage of seal wrapping in the current
release of Vault that have not yet been certified by Leidos. The mechanism for
seal wrapping is the same, they simply were not specifically evaluated by the
auditors.

* Root tokens
* Replication secondary activation tokens

## Activating Seal Wrapping

For some values, seal wrapping is always enabled with a supporting seal. This
includes the recovery key, any stored key shares, the master key, the keyring,
and more; essentially, any Critical Security Parameter (CSP) within Vault's
core. If upgrading from a version of Vault that did not support seal wrapping,
the next time these values are read they will be seal-wrapped and stored.

Backend mounts within Vault can also take advantage of seal wrapping. Seal
wrapping can be activated at mount time for a given mount by mounting the
backend with the `seal_wrap` configuration value set to `true`. (This value
cannot currently be changed later.)

A given backend's author can specify which values should be seal-wrapped by
identifying where CSPs are stored. If no specific CSPs are identifiable, all
data for the backend may be seal-wrapped.

To see the current list of seal-wrapped data per backend type, see the latest
audit letter and updates in the FIPS 140-2 Compliance section above.

Note that it is often an order of magnitude or two slower to write to and read
from HSMs or remote seals. However, values will be cached in memory
un-seal-wrapped (but still encrypted by Vault's built-in cryptographic barrier)
in Vault, which will mitigate this for read-heavy workloads.

## Seal Wrap and Replication

Seal wrapping takes place below the replication logic. As a result, it is
transparent to replication. Replication will convey which values should be
seal-wrapped, but it is up to the seal on the local cluster to implement it.
In practice, this means that seal wrapping can be used without needing to have
the replicated keys on both ends of the connection; each cluster can have
distinct keys in an HSM or in KMS.

In addition, it is possible to replicate from a Shamir-protected primary
cluster to clusters that use HSMs when seal wrapping is required in downstream
datacenters but not in the primary.

Because of the level of flexibility targeted for replication, values sent over
replication connections do not currently meet KeyTransit requirements for FIPS
140-2. Vault's clustering implementation does support best practices guidance
given in FIPS 140-2, but the cryptographic implementation of TLS is not FIPS
140-2 certified. We may look into providing certified TLS in the future for
replication traffic; in the meantime, a transparent TCP proxy that supports
certified FIPS 140-2 TLS (such as
[stunnel](https://www.stunnel.org/index.html)) can be used for replication
traffic if meeting KeyTransit requirements for replication is necessary.
