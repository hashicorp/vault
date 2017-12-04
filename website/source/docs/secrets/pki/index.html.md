---
layout: "docs"
page_title: "PKI Secret Backend"
sidebar_current: "docs-secrets-pki"
description: |-
  The PKI secret backend for Vault generates TLS certificates.
---

# PKI Secret Backend

Name: `pki`

The PKI secret backend for Vault generates X.509 certificates dynamically based
on configured roles. This means services can get certificates needed for both
client and server authentication without going through the usual manual process
of generating a private key and CSR, submitting to a CA, and waiting for a
verification and signing process to complete. Vault's built-in authentication
and authorization mechanisms provide the verification functionality.

By keeping TTLs relatively short, revocations are less likely to be needed,
keeping CRLs short and helping the backend scale to large workloads. This in
turn allows each instance of a running application to have a unique
certificate, eliminating sharing and the accompanying pain of revocation and
rollover.

In addition, by allowing revocation to mostly be forgone, this backend allows
for ephemeral certificates; certificates can be fetched and stored in memory
upon application startup and discarded upon shutdown, without ever being
written to disk.

This page will show a quick start for this backend. For detailed documentation
on every path, use `vault path-help` after mounting the backend.

## Considerations

To successfully deploy this backend, there are a number of important
considerations to be aware of, as well as some preparatory steps that should be
undertaken. You should read all of these *before* using this backend or
generating the CA to use with this backend.

### Be Careful with Root CAs

Vault storage is secure, but not as secure as a piece of paper in a bank vault.
It is, after all, networked software. If your root CA is hosted outside of
Vault, don't put it in Vault as well; instead, issue a shorter-lived
intermediate CA certificate and put this into Vault. This aligns with industry
best practices.

Since 0.4, the backend supports generating self-signed root CAs and creating
and signing CSRs for intermediate CAs. In each instance, for security reasons,
the private key can *only* be exported at generation time, and the ability to
do so is part of the command path (so it can be put into ACL policies).

If you plan on using intermediate CAs with Vault, it is suggested that you let
Vault create CSRs and do not export the private key, then sign those with your
root CA (which may be a second mount of the `pki` backend).

### One CA Certificate, One Backend

In order to vastly simplify both the configuration and codebase of the PKI
backend, only one CA certificate is allowed per backend. If you want to issue
certificates from multiple CAs, mount the PKI backend at multiple mount points
with separate CA certificates in each.

This also provides a convenient method of switching to a new CA certificate
while keeping CRLs valid from the old CA certificate; simply mount a new
backend and issue from there.

A common pattern is to have one mount act as your root CA and to use this CA
only to sign intermediate CA CSRs from other PKI mounts.

### Keep certificate lifetimes short, for CRL's sake

This backend aligns with Vault's philosophy of short-lived secrets. As such it
is not expected that CRLs will grow large; the only place a private key is ever
returned is to the requesting client (this backend does *not* store generated
private keys, except for CA certificates). In most cases, if the key is lost,
the certificate can simply be ignored, as it will expire shortly.

If a certificate must truly be revoked, the normal Vault revocation function
can be used; alternately a root token can be used to revoke the certificate
using the certificate's serial number. Any revocation action will cause the CRL
to be regenerated. When the CRL is regenerated, any expired certificates are
removed from the CRL (and any revoked, expired certificate are removed from
backend storage).

This backend does not support multiple CRL endpoints with sliding date windows;
often such mechanisms will have the transition point a few days apart, but this
gets into the expected realm of the actual certificate validity periods issued
from this backend. A good rule of thumb for this backend would be to simply not
issue certificates with a validity period greater than your maximum comfortable
CRL lifetime. Alternately, you can control CRL caching behavior on the client
to ensure that checks happen more often.

Often multiple endpoints are used in case a single CRL endpoint is down so that
clients don't have to figure out what to do with a lack of response. Run Vault
in HA mode, and the CRL endpoint should be available even if a particular node
is down.

### You must configure issuing/CRL/OCSP information *in advance*

This backend serves CRLs from a predictable location, but it is not possible
for the backend to know where it is running. Therefore, you must configure
desired URLs for the issuing certificate, CRL distribution points, and OCSP
servers manually using the `config/urls` endpoint. It is supported to have more
than one of each of these by passing in the multiple URLs as a comma-separated
string parameter.

### Safe Minimums

Since its inception, this backend has enforced SHA256 for signature hashes
rather than SHA1. As of 0.5.1, a minimum of 2048 bits for RSA keys is also
enforced. Software that can handle SHA256 signatures should also be able to
handle 2048-bit keys, and 1024-bit keys are considered unsafe and are
disallowed in the Internet PKI.

### Token Lifetimes and Revocation

When a token expires, it revokes all leases associated with it. This means that
long-lived CA certs need correspondingly long-lived tokens, something that is
easy to forget. Starting with 0.6, root and intermediate CA certs no longer
have associated leases, to prevent unintended revocation when not using a token
with a long enough lifetime. To revoke these certificates, use the `pki/revoke`
endpoint.

## Quick Start

#### Mount the backend

The first step to using the PKI backend is to mount it. Unlike the `kv`
backend, the `pki` backend is not mounted by default.

```text
$ vault mount pki
Successfully mounted 'pki' at 'pki'!
```

#### Configure a CA certificate

Next, Vault must be configured with a CA certificate and associated private
key. We'll take advantage of the backend's self-signed root generation support,
but Vault also supports generating an intermediate CA (with a CSR for signing)
or setting a PEM-encoded certificate and private key bundle directly into the
backend.

Generally you'll want a root certificate to only be used to sign CA
intermediate certificates, but for this example we'll proceed as if you will
issue certificates directly from the root. As it's a root, we'll want to set a
long maximum life time for the certificate; since it honors the maximum mount
TTL, first we adjust that:

```text
$ vault mount-tune -max-lease-ttl=87600h pki
Successfully tuned mount 'pki'!
```

That sets the maximum TTL for secrets issued from the mount to 10 years. (Note
that roles can further restrict the maximum TTL.)

Now, we generate our root certificate:

```text
$ vault write pki/root/generate/internal common_name=myvault.com ttl=87600h
Key             Value
certificate     -----BEGIN CERTIFICATE-----
MIIDvTCCAqWgAwIBAgIUAsza+fvOw+Xh9ifYQ0gNN0ruuWcwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wHhcNMTUxMTE5MTYwNDU5WhcNMjUx
MTE2MTYwNDU5WjAWMRQwEgYDVQQDEwtteXZhdWx0LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMUhH4OLf/sa6GuJONGC/CWLY7nDbfH8jAaCKgqV
eJ81KrmcgP8WPhoFsYHFQEQXQZcrJagwYfm19jYn3CaqrYPbciv9bcWi+ECxZV3x
Hs/YdCFk7KgDGCci37w+cy6fSB943FKJqqVbvPv0odmq6LvgGGgneznvuvkIrOWG
qVDrDdvbEZ01XAyzUQJaaiJXExN+6xm1HcBoypCP8ZjjnXHcFQvw2QBItLRU7iUd
ESFgbrkrSPW3HA6KF0ov2qFMoHTiQ6aM4KaHPmXcFPicugYR9owZfZ4lwWJCqT7j
EkhokaMgHnvyRScuiRZhQm8ppHZoYsqrc3glfEuxGHkS+0cCAwEAAaOCAQEwgf4w
DgYDVR0PAQH/BAQDAgGuMBMGA1UdJQQMMAoGCCsGAQUFBwMJMA8GA1UdEwEB/wQF
MAMBAf8wHQYDVR0OBBYEFLvAbt0eUUOoo7hjKiQM2bRqDKrZMB8GA1UdIwQYMBaA
FLvAbt0eUUOoo7hjKiQM2bRqDKrZMDsGCCsGAQUFBwEBBC8wLTArBggrBgEFBQcw
AoYfaHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3BraS9jYTAWBgNVHREEDzANggtt
eXZhdWx0LmNvbTAxBgNVHR8EKjAoMCagJKAihiBodHRwOi8vMTI3LjAuMC4xOjgy
MDAvdjEvcGtpL2NybDANBgkqhkiG9w0BAQsFAAOCAQEAVSgIRl6XJs95D7iXGzeQ
Ab8OIei779k0pD7xxS/+knY3TM6733zL/LXs4BEL3wfcQWoDrMtCW0Ook455sAOE
PSnTaZYQSH/F74VawWhSee4ZyiWq+sTUI4IzqYG3IS36mCyb0t6RxEb3aoQ87WHs
BHIB6uWbj6WoGHYM8ESxY89aY9jnX3xSs1HuluVW1uPrpIoa/eudpyV40Y1+9RNM
6fCX5LHGM7vKYxqvudYe+7G1MdKVBQg17h6XuieiUswVt2/HvDlNr+9DHrUla9Ve
Ig43v+grirlG7DrAr6Aiu/MVWKJP6CvNwG/XzrGaqd6KqSsE+8oIGR9tCTuPxI6v
SQ==
-----END CERTIFICATE-----
expiration      1.763309099e+09
issuing_ca      -----BEGIN CERTIFICATE-----
MIIDvTCCAqWgAwIBAgIUAsza+fvOw+Xh9ifYQ0gNN0ruuWcwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wHhcNMTUxMTE5MTYwNDU5WhcNMjUx
MTE2MTYwNDU5WjAWMRQwEgYDVQQDEwtteXZhdWx0LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMUhH4OLf/sa6GuJONGC/CWLY7nDbfH8jAaCKgqV
eJ81KrmcgP8WPhoFsYHFQEQXQZcrJagwYfm19jYn3CaqrYPbciv9bcWi+ECxZV3x
Hs/YdCFk7KgDGCci37w+cy6fSB943FKJqqVbvPv0odmq6LvgGGgneznvuvkIrOWG
qVDrDdvbEZ01XAyzUQJaaiJXExN+6xm1HcBoypCP8ZjjnXHcFQvw2QBItLRU7iUd
ESFgbrkrSPW3HA6KF0ov2qFMoHTiQ6aM4KaHPmXcFPicugYR9owZfZ4lwWJCqT7j
EkhokaMgHnvyRScuiRZhQm8ppHZoYsqrc3glfEuxGHkS+0cCAwEAAaOCAQEwgf4w
DgYDVR0PAQH/BAQDAgGuMBMGA1UdJQQMMAoGCCsGAQUFBwMJMA8GA1UdEwEB/wQF
MAMBAf8wHQYDVR0OBBYEFLvAbt0eUUOoo7hjKiQM2bRqDKrZMB8GA1UdIwQYMBaA
FLvAbt0eUUOoo7hjKiQM2bRqDKrZMDsGCCsGAQUFBwEBBC8wLTArBggrBgEFBQcw
AoYfaHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3BraS9jYTAWBgNVHREEDzANggtt
eXZhdWx0LmNvbTAxBgNVHR8EKjAoMCagJKAihiBodHRwOi8vMTI3LjAuMC4xOjgy
MDAvdjEvcGtpL2NybDANBgkqhkiG9w0BAQsFAAOCAQEAVSgIRl6XJs95D7iXGzeQ
Ab8OIei779k0pD7xxS/+knY3TM6733zL/LXs4BEL3wfcQWoDrMtCW0Ook455sAOE
PSnTaZYQSH/F74VawWhSee4ZyiWq+sTUI4IzqYG3IS36mCyb0t6RxEb3aoQ87WHs
BHIB6uWbj6WoGHYM8ESxY89aY9jnX3xSs1HuluVW1uPrpIoa/eudpyV40Y1+9RNM
6fCX5LHGM7vKYxqvudYe+7G1MdKVBQg17h6XuieiUswVt2/HvDlNr+9DHrUla9Ve
Ig43v+grirlG7DrAr6Aiu/MVWKJP6CvNwG/XzrGaqd6KqSsE+8oIGR9tCTuPxI6v
SQ==
-----END CERTIFICATE-----
serial_number   02:cc:da:f9:fb:ce:c3:e5:e1:f6:27:d8:43:48:0d:37:4a:ee:b9:67
```

The returned certificate is purely informational; it and its private key are
safely stored in the backend mount.

#### Set URL configuration

Generated certificates can have the CRL location and the location of the
issuing certificate encoded. These values must be set manually, but can be
changed at any time.

```text
$ vault write pki/config/urls issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"
Success! Data written to: pki/ca/urls
```

#### Configure a role

The next step is to configure a role. A role is a logical name that maps to a
policy used to generate those credentials. For example, let's create an
"example-dot-com" role:

```text
$ vault write pki/roles/example-dot-com \
    allowed_domains="example.com" \
    allow_subdomains="true" max_ttl="72h"
Success! Data written to: pki/roles/example-dot-com
```

#### Generate credentials

By writing to the `roles/example-dot-com` path we are defining the
`example-dot-com` role. To generate a new set of credentials, we simply write
to the `issue` endpoint with that role name: Vault is now configured to create
and manage certificates!

```text
$ vault write pki/issue/example-dot-com \
    common_name=blah.example.com
Key             	Value
---             	-----
lease_id        	pki/issue/example-dot-com/6d8ab3e2-ce31-8821-81e4-740a498af51d
lease_duration  	259199
lease_renewable 	false
certificate     	-----BEGIN CERTIFICATE-----
MIIDbDCCAlSgAwIBAgIUPiAyxq+nIE6xlWf7hrzLkPQxtvMwDQYJKoZIhvcNAQEL
BQAwMzExMC8GA1UEAxMoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1
dGhvcml0eTAeFw0xNjA5MjcwMDA5MTNaFw0xNjA5MjcwMTA5NDNaMBsxGTAXBgNV
BAMTEGJsYWguZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQDJAYB04IVdmSC/TimaA6BbXlvgBTZHL5wBUTmO4iHhenL0eDEXVe2Fd7Yq
75LiBJmcC96hKbqh5rwS8KwN9ElZI52/mSMC+IvoNlYHAf7shwfsjrVx3q7/bTFg
lz6wECn1ugysxynmMvgQD/pliRkxTQ7RMh4Qlh75YG3R9BHy9ZddklZp0aNaitts
0uufHnN1UER/wxBCZdWTUu34KDL9I6yE7Br0slKKHPdEsGlFcMkbZhvjslZ7DGvO
974S0qtOdKiawJZbpNPg0foGZ3AxesDUlkHmmgzUNes/sjknDYTHEfeXM6Uap0j6
XvyhCxqdeahb/Vtibg0z9I0IusJbAgMBAAGjgY8wgYwwDgYDVR0PAQH/BAQDAgOo
MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQU/5oy0rL7
TT0wX7KZK7qcXqgayNwwHwYDVR0jBBgwFoAUgM37P8oXmA972ztLfw+b1eIY5now
GwYDVR0RBBQwEoIQYmxhaC5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEA
CT2vI6/taeLTw6ZulUhLXEXYXWZu1gF8n2COjZzbZXmHxQAoZ3GtnSNwacPHAyIj
f3cA9Moo552y39LUtWk+wgFtQokWGK7LXglLaveNUBowOHq/xk0waiIinJcgTG53
Z/qnbJnTjAOG7JwVJplWUIiS1avCksrHt7heE2EGRGJALqyLZ119+PW6ogtCLUv1
X8RCTw/UkIF/LT+sLF0bXWy4Hn38Gjwj1MVv1l76cEGOVSHyrYkN+6AMnAP58L5+
IWE9tN3oac4x7jhbuNpfxazIJ8Q6l/Up5U5Evfbh6N1DI0/gFCP20fMBkHwkuLfZ
2ekZoSeCgFRDlHGkr7Vv9w==
-----END CERTIFICATE-----
issuing_ca      	-----BEGIN CERTIFICATE-----
MIIDijCCAnKgAwIBAgIUB28DoGwgGFKL7fbOu9S4FalHLn0wDQYJKoZIhvcNAQEL
BQAwLzEtMCsGA1UEAxMkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9y
aXR5MB4XDTE2MDkyNzAwMDgyMVoXDTI2MDkxNjE2MDg1MVowMzExMC8GA1UEAxMo
VmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1dGhvcml0eTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAOSCiSij4wy1wiMwvZt+rtU3IaO6ZTn9
LfIPuGsR5/QSJk37pCZQco1LgoE/rTl+/xu3bDovyHDmgObghC6rzVOX2Tpi7kD+
DOZpqxOsaS8ebYgxB/XJTSxyEJuSAcpSNLqqAiZivuQXdaD0N7H3Or0awwmKE9mD
I0g8CF4fPDmuuOG0ASn9fMqXVVt5tXtEqZ9yJYfNOXx3FOPjRVOZf+kvSc31wCKe
i/KmR0AQOmToKMzq988nLqFPTi9KZB8sEU20cGFeTQFol+m3FTcIru94EPD+nLUn
xtlLELVspYb/PP3VpvRj9b+DY8FGJ5nfSJl7Rkje+CD4VxJpSadin3kCAwEAAaOB
mTCBljAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU
gM37P8oXmA972ztLfw+b1eIY5nowHwYDVR0jBBgwFoAUj4YAIxRwrBy0QMRKLnD0
kVidIuYwMwYDVR0RBCwwKoIoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3Vi
IEF1dGhvcml0eTANBgkqhkiG9w0BAQsFAAOCAQEAA4buJuPNJvA1kiATLw1dVU2J
HPubk2Kp26Mg+GwLn7Vz45Ub133JCYfF3/zXLFZZ5Yub9gWTtjScrvNfQTAbNGdQ
BdnUlMmIRmfB7bfckhryR2R9byumeHATgNKZF7h8liNHI7X8tTzZGs6wPdXOLlzR
TlM3m1RNK8pbSPOkfPb06w9cBRlD8OAbNtJmuypXA6tYyiiMYBhP0QLAO3i4m1ns
aAjAgEjtkB1rQxW5DxoTArZ0asiIdmIcIGmsVxfDQIjFlRxAkafMs74v+5U5gbBX
wsOledU0fLl8KLq8W3OXqJwhGLK65fscrP0/omPAcFgzXf+L4VUADM4XhW6Xyg==
-----END CERTIFICATE-----
ca_chain        	[-----BEGIN CERTIFICATE-----
MIIDijCCAnKgAwIBAgIUB28DoGwgGFKL7fbOu9S4FalHLn0wDQYJKoZIhvcNAQEL
BQAwLzEtMCsGA1UEAxMkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9y
aXR5MB4XDTE2MDkyNzAwMDgyMVoXDTI2MDkxNjE2MDg1MVowMzExMC8GA1UEAxMo
VmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1dGhvcml0eTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAOSCiSij4wy1wiMwvZt+rtU3IaO6ZTn9
LfIPuGsR5/QSJk37pCZQco1LgoE/rTl+/xu3bDovyHDmgObghC6rzVOX2Tpi7kD+
DOZpqxOsaS8ebYgxB/XJTSxyEJuSAcpSNLqqAiZivuQXdaD0N7H3Or0awwmKE9mD
I0g8CF4fPDmuuOG0ASn9fMqXVVt5tXtEqZ9yJYfNOXx3FOPjRVOZf+kvSc31wCKe
i/KmR0AQOmToKMzq988nLqFPTi9KZB8sEU20cGFeTQFol+m3FTcIru94EPD+nLUn
xtlLELVspYb/PP3VpvRj9b+DY8FGJ5nfSJl7Rkje+CD4VxJpSadin3kCAwEAAaOB
mTCBljAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU
gM37P8oXmA972ztLfw+b1eIY5nowHwYDVR0jBBgwFoAUj4YAIxRwrBy0QMRKLnD0
kVidIuYwMwYDVR0RBCwwKoIoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3Vi
IEF1dGhvcml0eTANBgkqhkiG9w0BAQsFAAOCAQEAA4buJuPNJvA1kiATLw1dVU2J
HPubk2Kp26Mg+GwLn7Vz45Ub133JCYfF3/zXLFZZ5Yub9gWTtjScrvNfQTAbNGdQ
BdnUlMmIRmfB7bfckhryR2R9byumeHATgNKZF7h8liNHI7X8tTzZGs6wPdXOLlzR
TlM3m1RNK8pbSPOkfPb06w9cBRlD8OAbNtJmuypXA6tYyiiMYBhP0QLAO3i4m1ns
aAjAgEjtkB1rQxW5DxoTArZ0asiIdmIcIGmsVxfDQIjFlRxAkafMs74v+5U5gbBX
wsOledU0fLl8KLq8W3OXqJwhGLK65fscrP0/omPAcFgzXf+L4VUADM4XhW6Xyg==
-----END CERTIFICATE----- -----BEGIN CERTIFICATE-----
MIIDejCCAmKgAwIBAgIUDXJyQ1uJPF5ridDOCvGtVF1F8HUwDQYJKoZIhvcNAQEL
BQAwJzElMCMGA1UEAxMcVmF1bHQgVGVzdGluZyBSb290IEF1dGhvcml0eTAeFw0x
NjA5MjcwMDA4MjBaFw0yNjA5MjAyMDA4NTBaMC8xLTArBgNVBAMTJFZhdWx0IFRl
c3RpbmcgSW50ZXJtZWRpYXRlIEF1dGhvcml0eTCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAKHsRTw3aShwDTbywK7AeXNvz7IrmdOLAsd+svDdIUn/4kWQ
lAy4uXYncQc/V9bqLjza3tflK7otXT+V5GjbK+WpW5WSp8LkVhKdLRWOnPWJEC+B
nOucmLR0mFQF1W4Bfx0fYYCLdN/YbjSevPmA0UzlIN/pdQQoxUIvraTHPNBar94K
zmlMu06qAvl27LXYUE3nAhQaRGq4M39WbAUtRsNKaTU72qTpMsstpnBB1QBT2m2U
44twFpXZAgfR/hSqcA4NegPWmB5l+E2GhYfihOhVcnFaH2tgXb4MOMUyRH1hNdgZ
28K5G1ILt2+Rp+NSosA0LI3pV490SJfAxuc0tsUCAwEAAaOBlTCBkjAOBgNVHQ8B
Af8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUj4YAIxRwrBy0QMRK
LnD0kVidIuYwHwYDVR0jBBgwFoAULNIU30rP+wMVelJMFNyDtxgtq04wLwYDVR0R
BCgwJoIkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9yaXR5MA0GCSqG
SIb3DQEBCwUAA4IBAQCOjH2n8H1Q5KpaWTm378FKd2YY1nzI/nCwjAQX96VcJUrZ
W1ofPsTcCASQKwo3HC2ayV46DMiKoJWI+xOux2N+S9uVd+SC4ZloFzSER8cCDRRk
huVra+cAaljnkJVb4Ojv6vHnXljx9NrcW6KzJzwMf1HzewyG+P1EjD4/kcA5r0Gw
vuzGXMXmjMATf0LZlklDOHkNLtvnLS8axbXI05TlHIj9y9Y+aQyFYebwip+ZXYju
pIJFswrsCk5e2G6+UmhV81JH29IvjBi4POgqm2+mrGz5xS/i6flcs/8pn01jlDpC
knj9MxY9j42z2BKkhHayyuOa0BQm0TTu4S2fhajl
-----END CERTIFICATE-----]
private_key     	-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAyQGAdOCFXZkgv04pmgOgW15b4AU2Ry+cAVE5juIh4Xpy9Hgx
F1XthXe2Ku+S4gSZnAveoSm6oea8EvCsDfRJWSOdv5kjAviL6DZWBwH+7IcH7I61
cd6u/20xYJc+sBAp9boMrMcp5jL4EA/6ZYkZMU0O0TIeEJYe+WBt0fQR8vWXXZJW
adGjWorbbNLrnx5zdVBEf8MQQmXVk1Lt+Cgy/SOshOwa9LJSihz3RLBpRXDJG2Yb
47JWewxrzve+EtKrTnSomsCWW6TT4NH6BmdwMXrA1JZB5poM1DXrP7I5Jw2ExxH3
lzOlGqdI+l78oQsanXmoW/1bYm4NM/SNCLrCWwIDAQABAoIBAQCCbHMJY1Wl8eIJ
v5HG2WuHXaaHqVoavo2fXTDXwWryfx1v+zz/Q0YnQBH3shPAi/OQCTOfpw/uVWTb
dUZul3+wUyfcVmUdXGCLgBY53dWna8Z8e+zHwhISsqtDXV/TpelUBDCNO324XIIR
Cg0TLO4nyzQ+ESLo6D+Y2DTp8lBjMEkmKTd8CLXR2ycEoVykN98qPZm8keiLGO91
I8K7aRd8uOyQ6HUfJRlzFHSuwaLReErxGTEPI4t/wVqh2nP2gGBsn3apiJ0ul6Jz
NlYO5PqiwpeDk4ibhQBpicnm1jnEcynH/WtGuKgMNB0M4SBRBsEguO7WoKx3o+qZ
iVIaPWDhAoGBAO05UBvyJpAcz/ZNQlaF0EAOhoxNQ3h6+6ZYUE52PgZ/DHftyJPI
Y+JJNclY91wn91Yk3ROrDi8gqhzA+2Lelxo1kuZDu+m+bpzhVUdJia7tZDNzRIhI
24eP2GdochooOZ0qjvrik4kuX43amBhQ4RHsBjmX5CnUlL5ZULs8v2xnAoGBANjq
VLAwiIIqJZEC6BuBvVYKaRWkBCAXvQ3j/OqxHRYu3P68PZ58Q7HrhrCuyQHTph2v
fzfmEMPbSCrFIrrMRmjUG8wopL7GjZjFl8HOBHFwzFiz+CT5DEC+IJIRkp4HM8F/
PAzjB2wCdRdSjLTD5ph0/xQIg5xfln7D+wqU0QHtAoGBAKkLF0/ivaIiNftw0J3x
WxXag/yErlizYpIGCqvuzII6lLr9YdoViT/eJYrmb9Zm0HS9biCu2zuwDijRSBIL
RieyF40opUaKoi3+0JMtDwTtO2MCd8qaCH3QfkgqAG0tTuj1Q8/6F2JA/myKYamq
MMhhpYny9+7rAlemM8ZJIqtvAoGBAKOI3zpKDNCdd98A4v7B7H2usZUIJ7gOTZDo
XqiNyRENWb2PK6GNq/e6SrxvuclvyKA+zFnXULJoYtsj7tAH69lieGaOCc5uoRgZ
eBU7/euMj/McE6vEO3GgJawaJYCQi3uJMjvA+bp7i81+hehOfU5ZfmmbFaZSBoMh
u+U5Vu3tAoGBANnBIbHfD3E7rqnqdpH1oRRHLA1VdghzEKgyUTPHNDzPJG87RY3c
rRqeXepblud3qFjD60xS9BzcBijOvZ4+KHk6VIMpkyqoeNVFCJbBVCw+JGMp88+v
e9t+2iwryh5+rnq+pg6anmgwHldptJc1XEFZA2UUQ89RP7kOGQF6IkIS
-----END RSA PRIVATE KEY-----
private_key_type	rsa
serial_number   	3e:20:32:c6:af:a7:20:4e:b1:95:67:fb:86:bc:cb:90:f4:31:b6:f3
```

Vault has now generated a new set of credentials using the `example-dot-com`
role configuration. Here we see the dynamically generated private key and
certificate. The issuing CA certificate and CA trust chain is returned as well.

Using ACLs, it is possible to restrict using the pki backend such that trusted
operators can manage the role definitions, and both users and applications are
restricted in the credentials they are allowed to read.

If you get stuck at any time, simply run `vault path-help pki` or with a
subpath for interactive help output.


## API

The PKI secret backend has a full HTTP API. Please see the
[PKI secret backend API](/api/secret/pki/index.html) for more
details.
