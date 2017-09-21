---
layout: "docs"
page_title: "PKI - Secrets Engines"
sidebar_current: "docs-secrets-pki"
description: |-
  The PKI secrets engine for Vault generates TLS certificates.
---

# PKI Secrets Engine

The PKI secrets engine generates dynamic X.509 certificates. With this secrets
engine, services can get certificates without going through the usual manual
process of generating a private key and CSR, submitting to a CA, and waiting for
a verification and signing process to complete. Vault's built-in authentication
and authorization mechanisms provide the verification functionality.

By keeping TTLs relatively short, revocations are less likely to be needed,
keeping CRLs short and helping the secrets engine scale to large workloads. This
in turn allows each instance of a running application to have a unique
certificate, eliminating sharing and the accompanying pain of revocation and
rollover.

In addition, by allowing revocation to mostly be forgone, this secrets engine
allows for ephemeral certificates. Certificates can be fetched and stored in
memory upon application startup and discarded upon shutdown, without ever being
written to disk.

## Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.

1. Enable the PKI secrets engine:

    ```text
    $ vault secrets enable pki
    Success! Enabled the pki secrets engine at: pki/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.


1. Increase the TTL by tuning the secrets engine. The default value of 30 days may be too short, so increase it to 1 year:

    ```text
    $ vault secrets tune -max-lease-ttl=8760h pki
    Success! Tuned the secrets engine at: pki/
    ```

    Note that individual roles can restrict this value to be shorter on a
    per-certificate basis. This just configures the global maximum for this
    secrets engine.

1. Configure a CA certificate and private key. Vault can accept an existing key
pair, or it can generate its own self-signed root. In general, we recommend
maintaining your root CA outside of Vault and providing Vault a signed
intermediate CA.

    ```text
    $ vault write pki/root/generate/internal \
        common_name=my-website.com \
        ttl=8760h

    Key              Value
    ---              -----
    certificate      -----BEGIN CERTIFICATE-----...
    expiration       1536807433
    issuing_ca       -----BEGIN CERTIFICATE-----...
    serial_number    7c:f1:fb:2c:6e:4d:99:0e:82:1b:08:0a:81:ed:61:3e:1d:fa:f5:29
    ```

    The returned certificate is purely informative. The private key is safely
    stored internally in Vault.

1. Update the CRL location and issuing certificates. These values can be updated
in the future.

    ```text
    $ vault write pki/config/urls \
        issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" \
        crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"
    Success! Data written to: pki/config/urls
    ```

1. Configure a role that maps a name in Vault to a procedure for generating a
certificate. When users or machines generate credentials, they are generated
against this role:

    ```text
    $ vault write pki/roles/my-role \
        allowed_domains=my-website.com \
        allow_subdomains=true \
        max_ttl=72h
    Success! Data written to: pki/roles/my-role
    ```

## Usage

After the secrets engine is configured and a user/machine has a Vault token with
the proper permission, it can generate credentials.

1. Generate a new credential by writing to the `/issue` endpoint with the name
of the role:

    ```text
    $ vault write pki/issue/my-role \
        common_name=www.my-website.com

    Key                 Value
    ---                 -----
    certificate         -----BEGIN CERTIFICATE-----...
    issuing_ca          -----BEGIN CERTIFICATE-----...
    private_key         -----BEGIN RSA PRIVATE KEY-----...
    private_key_type    rsa
    serial_number       1d:2e:c6:06:45:18:60:0e:23:d6:c5:17:43:c0:fe:46:ed:d1:50:be
    ```

    The output will include a dynamically generated private key and certificate
    which corresponds to the given role and expires in 72h (as dictated by our
    role definition). The issuing CA and trust chain is also returned for
    automation simplicity.

## Considerations

To successfully deploy this secrets engine, there are a number of important
considerations to be aware of, as well as some preparatory steps that should be
undertaken. You should read all of these *before* using this secrets engine or
generating the CA to use with this secrets engine.

### Be Careful with Root CAs

Vault storage is secure, but not as secure as a piece of paper in a bank vault.
It is, after all, networked software. If your root CA is hosted outside of
Vault, don't put it in Vault as well; instead, issue a shorter-lived
intermediate CA certificate and put this into Vault. This aligns with industry
best practices.

Since 0.4, the secrets engine supports generating self-signed root CAs and
creating and signing CSRs for intermediate CAs. In each instance, for security
reasons, the private key can *only* be exported at generation time, and the
ability to do so is part of the command path (so it can be put into ACL
policies).

If you plan on using intermediate CAs with Vault, it is suggested that you let
Vault create CSRs and do not export the private key, then sign those with your
root CA (which may be a second mount of the `pki` secrets engine).

### One CA Certificate, One Secrets Engine

In order to vastly simplify both the configuration and codebase of the PKI
secrets engine, only one CA certificate is allowed per secrets engine. If you
want to issue certificates from multiple CAs, mount the PKI secrets engine at
multiple mount points with separate CA certificates in each.

This also provides a convenient method of switching to a new CA certificate
while keeping CRLs valid from the old CA certificate; simply mount a new secrets
engine and issue from there.

A common pattern is to have one mount act as your root CA and to use this CA
only to sign intermediate CA CSRs from other PKI secrets engines.

### Keep certificate lifetimes short, for CRL's sake

This secrets engine aligns with Vault's philosophy of short-lived secrets. As
such it is not expected that CRLs will grow large; the only place a private key
is ever returned is to the requesting client (this secrets engine does *not*
store generated private keys, except for CA certificates). In most cases, if the
key is lost, the certificate can simply be ignored, as it will expire shortly.

If a certificate must truly be revoked, the normal Vault revocation function can
be used; alternately a root token can be used to revoke the certificate using
the certificate's serial number. Any revocation action will cause the CRL to be
regenerated. When the CRL is regenerated, any expired certificates are removed
from the CRL (and any revoked, expired certificate are removed from secrets
engine storage).

This secrets engine does not support multiple CRL endpoints with sliding date
windows; often such mechanisms will have the transition point a few days apart,
but this gets into the expected realm of the actual certificate validity periods
issued from this secrets engine. A good rule of thumb for this secrets engine
would be to simply not issue certificates with a validity period greater than
your maximum comfortable CRL lifetime. Alternately, you can control CRL caching
behavior on the client to ensure that checks happen more often.

Often multiple endpoints are used in case a single CRL endpoint is down so that
clients don't have to figure out what to do with a lack of response. Run Vault in HA mode, and the CRL endpoint should be available even if a particular node
is down.

### You must configure issuing/CRL/OCSP information *in advance*

This secrets engine serves CRLs from a predictable location, but it is not
possible for the secrets engine to know where it is running. Therefore, you must
configure desired URLs for the issuing certificate, CRL distribution points, and
OCSP servers manually using the `config/urls` endpoint. It is supported to have
more than one of each of these by passing in the multiple URLs as a
comma-separated string parameter.

### Safe Minimums

Since its inception, this secrets engine has enforced SHA256 for signature
hashes rather than SHA1. As of 0.5.1, a minimum of 2048 bits for RSA keys is
also enforced. Software that can handle SHA256 signatures should also be able to
handle 2048-bit keys, and 1024-bit keys are considered unsafe and are disallowed
in the Internet PKI.

### Token Lifetimes and Revocation

When a token expires, it revokes all leases associated with it. This means that
long-lived CA certs need correspondingly long-lived tokens, something that is
easy to forget. Starting with 0.6, root and intermediate CA certs no longer have
associated leases, to prevent unintended revocation when not using a token with
a long enough lifetime. To revoke these certificates, use the `pki/revoke`
endpoint.


## API

The PKI secrets engine has a full HTTP API. Please see the
[PKI secrets engine API](/api/secret/pki/index.html) for more
details.
