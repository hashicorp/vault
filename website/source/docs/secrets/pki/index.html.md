---
layout: "docs"
page_title: "Secret Backend: PKI"
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

A common pattern is to have one mount act as your root CA, and which is only
used for signing intermediate CA CSRs mounted at other locations.

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

### No OCSP support, yet

Vault's architecture does not currently allow for a binary protocol such as
OCSP to be supported by a backend. As such, you should configure your software
to use CRLs for revocation information, with a caching lifetime that feels good
to you. Since you are following the advice above about keeping lifetimes short
(right?), CRLs should not grow too large, however, you can configure alternate
CRL and/or OCSP servers using `config/urls` if you wish.

If you are using issued certificates for client authentication to Vault, note
that as of 0.4, the `cert` authentication endpoint supports being pushed CRLs,
but it cannot read CRLs directly from this backend.

### Safe Minimums

Since its inception, this backend has enforced SHA256 for signature hashes
rather than SHA1. As of 0.5.1, a minimum of 2048 bits for RSA keys is also
enforced. Software that can handle SHA256 signatures should also be able to
handle 2048-bit keys, and 1024-bit keys are considered unsafe and are
disallowed in the Internet PKI.

## Quick Start

#### Mount the backend

The first step to using the PKI backend is to mount it. Unlike the `generic`
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
lease_id        pki/root/generate/internal/aa959dd4-467e-e5ff-642b-371add518b40
lease_duration  315359999
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
Key                     Value
lease_id                pki/issue/example-dot-com/32db49a9-61dd-f9ca-f4a6-aaefafe53739
lease_duration          259199
lease_renewable         false
certificate             -----BEGIN CERTIFICATE-----
MIIDVDCCAjygAwIBAgIUFMne7ro1DyvpZV+URHMPtz3vv9MwDQYJKoZIhvcNAQEL
BQAwGzEZMBcGA1UEAxMQaW50ZXJtZWRpYXRlLmNvbTAeFw0xNTExMjAxODEzMDNa
Fw0xNTExMjMxODEzMDNaMBsxGTAXBgNVBAMTEGJsYWguZXhhbXBsZS5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDEfl5iimLtEkyZiKYs+PfoaAbe
VhtZxTYWFA+B7IIlr8iPTEZDuf72OOf0Nu/8TNrV6Zkoq0EuHvBtNi1ut3usYt6V
i9lrcofD/Qjn0EH6aOZnE5J7c+gmsODhlxflLfz8uEytonKxofwOw26J9+subuSI
migMbpbTP/BL2n29K4NJoKdhh8VMxNCFPQxgu5ACEdJ9GsiPO5wBb7ifdIq0HcSU
0ONe6uZqDXeKiqfrTg6eap4EaALogkJhuk8BcAJv9aSbJswOSXTGROa4XChtCXEu
D3yVOoZOOm7JSm60y7ntf/dxZF5xcZXjRe6GkXJAIADOL9E5dOlgTFlojYjpAgMB
AAGjgY8wgYwwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggr
BgEFBQcDAjAdBgNVHQ4EFgQUXSgxVNG7knU9Od1iKh3eCjE/P9UwHwYDVR0jBBgw
FoAU83Ovrk0BzVapM+C/fA25dv/Ju34wGwYDVR0RBBQwEoIQYmxhaC5leGFtcGxl
LmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAJwwnK+GH+X2CK18qzv3fuPzbsR4pWKTB
aaSweW83+QUiRgsR6sgdNqYH0bE1BO/nysDC8hj2IH77KtOfJJgcqG8w709022rh
DtZVB33lU17oZC4LUIhq/Ym0JEwYryKck8ClxJpWYKy/kcNwt/WcoAY+aX07c+a4
0ACCpzTX8vUAxFxcp6ZCwebbSTbv56KHDxMRUSiWiwcDaLlqECsqfETN6eCN/M6A
GlPoswzIHjXSrAhI8KADOQf4oHI2cOj7ecJX9EqTq5snxFKblS8B12q1javiQGJS
7eni6Irw6x/enuPxp2VdPJOxPkMSf/+BcADDQ4mOrFtYg7u7+AvWBw==
-----END CERTIFICATE-----
issuing_ca              -----BEGIN CERTIFICATE-----
MIIDUzCCAjugAwIBAgIUaJJpBnXW+GkzJ6k6fQ9mqVKkEmQwDQYJKoZIhvcNAQEL
BQAwEzERMA8GA1UEAxMIcm9vdC5jb20wHhcNMTUxMTIwMTgxMDU2WhcNMzUxMTE1
MTcxMDU2WjAbMRkwFwYDVQQDExBpbnRlcm1lZGlhdGUuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzTR8agj+uOYCvOSDsNYYb186aDo9QLPHAr96
tDPjQFyp5Yjr60I6+say7zRPPGLgd3BealJ9EFeO7qlXeqi5z3H0wLzrMYqcGlok
2vSzdbs03+0/QXOys2R1Bzb+FOd4VoUSTiZ+a8wf07tmusNJCDE5/kI32+etQa8l
5a9ZQlwpgIWZSQmQjCA5B/0T6WQVwnELHOAGv+mJk7bAY/LVQkjUHzvimySsUmSb
sB20BPKhammJUDEcObwuJxA+f7NbkXzEypnR0pGULM32TR5Bmzij/iX33XDt3JKi
I9e0gJIT7bk91heWXLpEv2/+7g/lcXm/Hl9KvtX5WUtVPNENUQIDAQABo4GWMIGT
MA4GA1UdDwEB/wQEAwIBrjATBgNVHSUEDDAKBggrBgEFBQcDCTAPBgNVHRMBAf8E
BTADAQH/MB0GA1UdDgQWBBTzc6+uTQHNVqkz4L98Dbl2/8m7fjAfBgNVHSMEGDAW
gBSmqzdLJm2brQl3mpyIDRTohpQrMDAbBgNVHREEFDASghBpbnRlcm1lZGlhdGUu
Y29tMA0GCSqGSIb3DQEBCwUAA4IBAQAubC0cuwbitp2Fq5FgH8Mu/Fzhf5qWftxE
a7VagVExs2uxP5yD57bWck6vZrks03SVk4GFR9yyIVbOIUAVEm1Rw1/PK77l9/2c
fYhy0OQVZweO+olOgEfC8gYLaBT5Vo3D1CjV/Vb2VGCct3dmMsXuD04HOy1mTz2p
3yPx1wPoUYNaEu+7gzvUxh+8AM3JmCcrsaa1R9AsayAXtLuCJm9Fy6bU4I3wbxBp
zTOT7fmkjpjCV4acfgcPF2F90TfcesHl9oUgNsu4tChABiPENA4h2A4yVku9onaQ
JrNqv2SnJaYH4OTgtguC0cLB7hvr/Sc73pU55OSs2KZWhLZRWAJv
-----END CERTIFICATE-----
private_key             -----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAxH5eYopi7RJMmYimLPj36GgG3lYbWcU2FhQPgeyCJa/Ij0xG
Q7n+9jjn9Dbv/Eza1emZKKtBLh7wbTYtbrd7rGLelYvZa3KHw/0I59BB+mjmZxOS
e3PoJrDg4ZcX5S38/LhMraJysaH8DsNuiffrLm7kiJooDG6W0z/wS9p9vSuDSaCn
YYfFTMTQhT0MYLuQAhHSfRrIjzucAW+4n3SKtB3ElNDjXurmag13ioqn604Onmqe
BGgC6IJCYbpPAXACb/WkmybMDkl0xkTmuFwobQlxLg98lTqGTjpuyUputMu57X/3
cWRecXGV40XuhpFyQCAAzi/ROXTpYExZaI2I6QIDAQABAoIBAQC11Rc060kmh6OZ
BOp6fZ60U+ffQiGnTieCAOhk29+IToYzjWsMa4d0hS6pQVmNyfVMtRJFn0z/CCSH
e/ZJGcR5vzipfTQjCWZ3yKXAF2mm+AIW6vbIBXeUrmQ8fpzfOVJ+73IN0GGA3hyp
8NJPHLxnSLl1a+qZrpEmRmnxV+y563Uporr7KriwUTDO+F/4YoJ92deDgW6FTR1/
vB6QIeSoTa1bXrkz+jvNyFh/Z+c3DKAyndzjCRfMTOygwbxJQdrZv1yNfkNSym6/
lGjZ6GnG/3hkh9nDldHLNJeJ4FbKyN5dbGnXrj2BI156KDE0jnO88TciX0khlRkp
efj9LcOBAoGBANB2HvtNp0eKDfGtey8Qxl6HP48mwg3pYLEATZtRowBFkyThfOoJ
Kf8ANNZHCPI/EilEYT6Xin3gdZNpC19gh3K86bCG8Uvq5TRtahO2oM4CenNOJAhe
utGfv5TfBBRodSRk2939masuavPLEqORKsBt3GSTv6Z1+Pfa+SjLm5F5AoGBAPFN
kvwXBvNkJOx594C9RVYmznD8z+IxVgqSIInpvu47SXJOWoS/hEwxfTCzCldl2ejK
eZ+mhJFud+wMbryF7AsW58+JfXr3lNJj7RquNx5cN+DSDV+1fESm9Me6e8XTxEqH
+ZoVPe/TRG6zJ+k7S99IDR1rE6D/oEjNpB/NxsbxAoGBAJ0Sc9O9Ni8UWd9hbTEQ
fbfaRszxUkSzNZUI+nDuuVhKFE40zS93CjrHCAjw60/EsEWB7ZgBDWw9hbo160jJ
biXJLHhDpWsjqeKwEr6Z3F59xZA+L65S2od6zBs7U1KhRqrOiFCjdndieVoLCJdQ
mZr27JqoLT8bIyZ2y0iu6iBZAoGADyTQMbP8QrApRRIOf2zhehurXxnurgJspPMw
yZb63Zao8FyMf8JJOkLs2W6TGpMQzvROF7/ql/n32r+Y/4nkG3oPiE3Xqyz4kQ+m
ZMNEQEqHUzu7jSMlrmVP/WztsaetrQPFnW7x2ShIJi5mNdP72gJ6mDsNG1CPraIC
R+CxNfECgYEAkqjs3j8Div/XnjPs3TuYtef7B/vgamzQujVkRR6+WkTvnZ1Hp9ge
vWrnOH3LbV8kfzzq0nbLK0iF7q4gl0czMJioVusaWrga2xkIfiI4yuHrIGatdTos
nPDHwoPZeRqBv/9OXSfQkYu+FiJnLEoztMb6f1Z1cPjvbuou2FB1p18=
-----END RSA PRIVATE KEY-----
private_key_type        rsa
serial_number           14:c9:de:ee:ba:35:0f:2b:e9:65:5f:94:44:73:0f:b7:3d:ef:bf:d3
```

Vault has now generated a new set of credentials using the `example-dot-com`
role configuration. Here we see the dynamically generated private key and
certificate. The issuing CA certificate is returned as well.

Using ACLs, it is possible to restrict using the pki backend such that trusted
operators can manage the role definitions, and both users and applications are
restricted in the credentials they are allowed to read.

If you get stuck at any time, simply run `vault path-help pki` or with a
subpath for interactive help output.


## API

### /pki/ca(/pem)
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves the CA certificate *in raw DER-encoded form*. This is a bare
    endpoint that does not return a standard Vault data structure. If `/pem` is
    added to the endpoint, the CA certificate is returned in PEM format. <br />
    <br />This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/ca(/pem)`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```
    <binary DER-encoded certficiate>
    ```

  </dd>
</dl>

### /pki/cert/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves one of a selection of certificates. Valid values: `ca` for the CA
    certificate, `crl` for the current CRL, or a serial number in either
    hyphen-separated or colon-separated octal format. This endpoint returns
    the certificate in PEM formatting in the `certificate` key of the JSON
    object. <br /><br />This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/cert/<serial>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIGmDCCBYCgAwIBAgIHBzEB3fTzhTANBgkqhkiG9w0BAQsFADCBjDELMAkGA1UE\n..."
      }
    }
    ...
    ```

  </dd>
</dl>

### /pki/config/ca
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows submitting the CA information for the backend via a PEM file
    containing the CA certificate and its private key, concatenated. Not needed
    if you are generating a self-signed root certificate, and not used if you
    have a signed intermediate CA certificate with a generated key (use the
    `/pki/intermediate/set-signed` endpoint for that). _If you have already set
    a certificate and key, they will be overridden._<br /><br />The information
    can be provided from a file via a `curl` command similar to the
    following:<br/>

    ```text
    $ curl \
        -H "X-Vault-Token:06b9d..." \
        -X POST \
        --data "@cabundle.json" \
        http://127.0.0.1:8200/v1/pki/config/ca
    ```

    Note that if you provide the data through the HTTP API it must be
    JSON-formatted, with newlines replaced with `\n`, like so:

    ```javascript
    {
      "pem_bundle": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END CERTIFICATE-----"
    }
    ```
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/config/ca`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">pem_bundle</span>
        <span class="param-flags">required</span>
        The key and certificate concatenated in PEM format.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>


### /pki/config/crl
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows getting the duration for which the generated CRL should be marked
    valid.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/config/crl`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
          "expiry": "72h"
        },
      "auth": null
    }
    ```

  </dd>
</dl>

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows setting the duration for which the generated CRL should be marked
    valid.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/config/crl`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
      <li>
        <span class="param">expiry</span>
        <span class="param-flags">required</span>
        The time until expiration. Defaults to `72h`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /pki/config/urls

#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Fetch the URLs to be encoded in generated certificates.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/config/urls`</dd>

  <dt>Parameters</dt>
  <dd>
    None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
          "issuing_certificates": [<url1>, <url2>],
          "crl_distribution_points": [<url1>, <url2>],
          "ocsp_servers": [<url1>, <url2>],
        },
      "auth": null
    }
    ```

  </dd>
</dl>

#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows setting the issuing certificate endpoints, CRL distribution points,
    and OCSP server endpoints that will be encoded into issued certificates.
    You can update any of the values at any time without affecting the other
    existing values. To remove the values, simply use a blank string as the
    parameter.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/config/urls`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">issuing_certificates</span>
        <span class="param-flags">optional</span>
        The URL values for the Issuing Certificate field.
      </li>
      <li>
        <span class="param">crl_distribution_points</span>
        <span class="param-flags">optional</span>
        The URL values for the CRL Distribution Points field.
      </li>
      <li>
        <span class="param">ocsp_servers</span>
        <span class="param-flags">optional</span>
        The URL values for the OCSP Servers field.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /pki/crl(/pem)
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves the current CRL *in raw DER-encoded form*. This endpoint
    is suitable for usage in the CRL Distribution Points extension in a
    CA certificate. This is a bare endpoint that does not return a
    standard Vault data structure. If `/pem` is added to the endpoint,
    the CRL is returned in PEM format.
    <br /><br />This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/crl(/pem)`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```
    <binary DER-encoded CRL>
    ```

  </dd>
</dl>

### /pki/crl/rotate
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
  This endpoint forces a rotation of the CRL. This can be used
  by administrators to cut the size of the CRL if it contains
  a number of certificates that have now expired, but has
  not been rotated due to no further certificates being revoked.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/crl/rotate`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "success": true
      }
    }
    ```

  </dd>
</dl>

### /pki/intermediate/generate
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new private key and a CSR for signing. If using Vault as a
    root, and for many other CAs, the various parameters on the final
    certificate are set at signing time and may or may not honor the parameters
    set here. _This will overwrite any previously existing CA private key._ If
    the path ends with `exported`, the private key will be returned in the
    response; if it is `internal` the private key will not be returned and
    *cannot be retrieved later*. <br /><br />This is mostly meant as a helper
    function, and not all possible parameters that can be set in a CSR are
    supported.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/intermediate/generate/[exported|internal]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">common_name</span>
        <span class="param-flags">required</span>
        The requested CN for the certificate.
      </li>
      <li>
        <span class="param">alt_names</span>
        <span class="param-flags">optional</span>
        Requested Subject Alternative Names, in a comma-delimited list. These
        can be host names or email addresses; they will be parsed into their
        respective fields.
      </li>
      <li>
        <span class="param">ip_sans</span>
        <span class="param-flags">optional</span>
        Requested IP Subject Alternative Names, in a comma-delimited list.
      </li>
      <li>
      <span class="param">format</span>
      <span class="param-flags">optional</span>
        Format for returned data. Can be `pem`, `der`, or `pem_bundle`;
        defaults to `pem`. If `der`, the output is base64 encoded. If
        `pem_bundle`, the `csr` field will contain the private key (if
        exported) and CSR, concatenated.
      </li>
      <li>
        <span class="param">key_type</span>
        <span class="param-flags">optional</span>
        Desired key type; must be `rsa` or `ec`. Defaults to `rsa`.
      </li>
      <li>
        <span class="param">key_bits</span>
        <span class="param-flags">optional</span>
        The number of bits to use. Defaults to `2048`. Must be changed to a
        valid value if the `key_type` is `ec`.
      </li>
      <li>
        <span class="param">exclude_cn_from_sans</span>
        <span class="param-flags">optional</span>
        If set, the given `common_name` will not be included in DNS or Email
        Subject Alternate Names (as appropriate). Useful if the CN is not a
        hostname or email address, but is instead some human-readable
        identifier.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "",
      "renewable": false,
      "lease_duration": 21600,
      "data": {
        "csr": "-----BEGIN CERTIFICATE REQUEST-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE REQUEST-----\n",
        "private_key": "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpAIBAAKCAQEAwsANtGz9gS3o5SwTSlOG1l-----END RSA PRIVATE KEY-----",
        "private_key_type": "rsa"
        },
      "warnings": null,
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/intermediate/set-signed
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows submitting the signed CA certificate corresponding to a private key generated via `/pki/intermediate/generate`. The certificate should be submitted in PEM format; see the documentation for `/pki/config/ca` for some hints on submitting.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/intermediate/set-signed`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">certificate</span>
        <span class="param-flags">required</span>
        The certificate in PEM format.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /pki/issue/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new set of credentials (private key and certificate) based on
    the role named in the endpoint. The issuing CA certificate is returned as
    well, so that only the root CA need be in a client's trust store.  <br
    /><br />*The private key is _not_ stored.  If you do not save the private
    key, you will need to request a new certificate.*
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/issue/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">common_name</span>
        <span class="param-flags">required</span>
        The requested CN for the certificate. If the CN is allowed by role
        policy, it will be issued.
      </li>
      <li>
        <span class="param">alt_names</span>
        <span class="param-flags">optional</span>
        Requested Subject Alternative Names, in a comma-delimited list. These
        can be host names or email addresses; they will be parsed into their
        respective fields. If any requested names do not match role policy, the
        entire request will be denied.
      </li>
      <li>
        <span class="param">ip_sans</span>
        <span class="param-flags">optional</span>
        Requested IP Subject Alternative Names, in a comma-delimited list. Only
        valid if the role allows IP SANs (which is the default).
      </li>
      <li>
      <span class="param">ttl</span>
      <span class="param-flags">optional</span>
        Requested Time To Live. Cannot be greater than the role's `max_ttl`
        value. If not provided, the role's `ttl` value will be used. Note that
        the role values default to system values if not explicitly set.
      </li>
      <li>
      <span class="param">format</span>
      <span class="param-flags">optional</span>
        Format for returned data. Can be `pem`, `der`, or `pem_bundle`;
        defaults to `pem`. If `der`, the output is base64 encoded. If
        `pem_bundle`, the `certificate` field will contain the private key,
        certificate, and issuing CA, concatenated.
      </li>
      <li>
        <span class="param">exclude_cn_from_sans</span>
        <span class="param-flags">optional</span>
        If set, the given `common_name` will not be included in DNS or Email
        Subject Alternate Names (as appropriate). Useful if the CN is not a
        hostname or email address, but is instead some human-readable
        identifier.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "pki/issue/test/7ad6cfa5-f04f-c62a-d477-f33210475d05",
      "renewable": false,
      "lease_duration": 21600,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n",
        "private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAnVHfwoKsUG1GDVyWB1AFroaKl2ImMBO8EnvGLRrmobIkQvh+\n...\nQN351pgTphi6nlCkGPzkDuwvtxSxiCWXQcaxrHAL7MiJpPzkIBq1\n-----END RSA PRIVATE KEY-----\n",
        "private_key_type": "rsa",
        "serial_number": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
      "warnings": "",  
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/revoke
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Revokes a certificate using its serial number. This is an
    alternative option to the standard method of revoking
    using Vault lease IDs. A successful revocation will
    rotate the CRL.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/revoke`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">serial_number</span>
        <span class="param-flags">required</span>
        The serial number of the certificate to revoke, in
        hyphen-separated or colon-separated octal.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "revocation_time": 1433269787
      }
    }
    ```
  </dd>
</dl>

### /pki/roles/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates or updates the role definition. Note that the
    `allowed_domains`, `allow_subdomains`, and
    `allow_any_name` attributes are additive; between them nearly and across
    multiple roles nearly any issuing policy can be accommodated.
    `server_flag`, `client_flag`, and `code_signing_flag` are additive as well.
    If a client requests a certificate that is not allowed by the CN policy in
    the role, the request is denied.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/roles/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The Time To Live value provided as a string duration with time suffix.
        Hour is the largest suffix.  If not set, uses the system default value
        or the value of `max_ttl`, whichever is shorter.
      </li>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">optional</span>
        The maximum Time To Live provided as a string duration with time
        suffix. Hour is the largest suffix. If not set, defaults to the system
        maximum lease TTL.
      </li>
      <li>
        <span class="param">allow_localhost</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates for `localhost` as one of the
        requested common names. This is useful for testing and to allow clients
        on a single host to talk securely. Defaults to true.
      </li>
      <li>
        <span class="param">allowed_domains</span>
        <span class="param-flags">optional</span>
        Designates the domains of the role, provided as a comma-separated list.
        This is used with the `allow_bare_domains` and `allow_subdomains`
        options. There is no default.
      </li>
      <li>
        <span class="param">allow_bare_domains</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates matching the value of the
        actual domains themselves; e.g. if a configured domain set with
        `allowed_domains` is `example.com`, this allows clients to actually
        request a certificate containing the name `example.com` as one of the
        DNS values on the final certificate. In some scenarios, this can be
        considered a security risk. Defaults to false.
      </li>
      <li>
        <span class="param">allow_subdomains</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates with CNs that are subdomains
        of the CNs allowed by the other role options. _This includes wildcard
        subdomains._ For example, an `allowed_domains` value of
        `example.com` with this option set to true will allow `foo.example.com`
        and `bar.example.com` as well as `*.example.com`. This is redundant
        when using the `allow_any_name` option.  Defaults to `false`.
      </li>
      <li>
        <span class="param">allow_any_name</span>
        <span class="param-flags">optional</span>
        If set, clients can request any CN. Useful in some circumstances, but
        make sure you understand whether it is appropriate for your
        installation before enabling it.  Defaults to `false`.
      </li>
      <li>
        <span class="param">enforce_hostnames</span>
        <span class="param-flags">optional</span>
        If set, only valid host names are allowed for CNs, DNS SANs, and the
        host part of email addresses. Defaults to `true`.
      </li>
      <li>
        <span class="param">allow_ip_sans</span>
        <span class="param-flags">optional</span>
        If set, clients can request IP Subject Alternative Names. No
        authorization checking is performed except to verify that the given
        values are valid IP addresses. Defaults to `true`.
      <li>
        <span class="param">server_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for server use.  Defaults to `true`.
      </li>
      <li>
        <span class="param">client_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for client use.  Defaults to `true`.
      </li>
      <li>
        <span class="param">code_signing_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for code signing use. Defaults to
        `false`.
      </li>
      <li>
        <span class="param">email_protection_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for email protection use. Defaults to
        `false`.
      </li>
      <li>
        <span class="param">key_type</span>
        <span class="param-flags">optional</span>
        The type of key to generate for generated private
        keys. Currently, `rsa` and `ec` are supported.
        Defaults to `rsa`.
      </li>
      <li>
        <span class="param">key_bits</span>
        <span class="param-flags">optional</span>
        The number of bits to use for the generated keys.
        Defaults to `2048`; this will need to be changed for
        `ec` keys. See https://golang.org/pkg/crypto/elliptic/#Curve
        for an overview of allowed bit lengths for `ec`.
      </li>
      <li>
        <span class="param">use_csr_common_name</span>
        <span class="param-flags">optional</span>
        If set, when used with the CSR signing endpoint, the common name in the
        CSR will be used instead of taken from the JSON data. This does `not`
        include any requested SANs in the CSR. Defaults to `false`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Queries the role definition.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/roles/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "data": {
        "allow_any_name": false,
        "allow_ip_sans": true,
        "allow_localhost": true,
        "allow_subdomains": false,
        "allowed_domains": "example.com,foobar.com",
        "client_flag": true,
        "code_signing_flag": false,
        "key_bits": 2048,
        "key_type": "rsa",
        "ttl": "6h",
        "max_ttl": "12h",
        "server_flag": true
      }
    }
    ```

  </dd>
</dl>

#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns a list of available roles. Only the role names are returned, not
    any values.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/roles/?list=true`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

  ```javascript
  {
    "auth": null,
    "data": {
      "keys": ["dev", "prod"]
    },
    "lease_duration": 2592000,
    "lease_id": "",
    "renewable": false
  }
  ```

  </dd>
</dl>

#### DELETE

<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the role definition. Deleting a role does <b>not</b> revoke
    certificates previously issued under this role.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/pki/roles/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
  </dd>
</dl>

### /pki/root/generate
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new self-signed CA certificate and private key. _This will
    overwrite any previously-existing private key and certificate._ If the path
    ends with `exported`, the private key will be returned in the response; if
    it is `internal` the private key will not be returned and *cannot be
    retrieved later*. Distribution points use the values set via `config/urls`.
    <br /><br />As with other issued certificates, Vault will automatically
    revoke the generated root at the end of its lease period; the CA
    certificate will sign its own CRL.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/root/generate/[exported|internal]`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">common_name</span>
        <span class="param-flags">required</span>
        The requested CN for the certificate.
      </li>
      <li>
        <span class="param">alt_names</span>
        <span class="param-flags">optional</span>
        Requested Subject Alternative Names, in a comma-delimited list. These
        can be host names or email addresses; they will be parsed into their
        respective fields.
      </li>
      <li>
        <span class="param">ip_sans</span>
        <span class="param-flags">optional</span>
        Requested IP Subject Alternative Names, in a comma-delimited list.
      </li>
      <li>
      <span class="param">ttl</span>
      <span class="param-flags">optional</span>
        Requested Time To Live (after which the certificate will be expired).
        This cannot be larger than the mount max (or, if not set, the system
        max).
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
        Format for returned data. Can be `pem`, `der`, or `pem_bundle`;
        defaults to `pem`. If `der`, the output is base64 encoded. If
        `pem_bundle`, the `certificate` field will contain the private key (if exported),
        certificate, and issuing CA, concatenated.
      </li>
      <li>
        <span class="param">key_type</span>
        <span class="param-flags">optional</span>
        Desired key type; must be `rsa` or `ec`. Defaults to `rsa`.
      </li>
      <li>
        <span class="param">key_bits</span>
        <span class="param-flags">optional</span>
        The number of bits to use. Defaults to `2048`. Must be changed to a
        valid value if the `key_type` is `ec`.
      </li>
      <li>
        <span class="param">max_path_length</span>
        <span class="param-flags">optional</span>
        If set, the maximum path length to encode in the generated certificate.
        Defaults to `-1`, which means no limit.  unless the signing certificate
        has a maximum path length set, in which case the path length is set to
        one less than that of the signing certificate.  A limit of `0` means a
        literal path length of zero.
      </li>
      <li>
        <span class="param">exclude_cn_from_sans</span>
        <span class="param-flags">optional</span>
        If set, the given `common_name` will not be included in DNS or Email
        Subject Alternate Names (as appropriate). Useful if the CN is not a
        hostname or email address, but is instead some human-readable
        identifier.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "pki/root/generate/internal/aa959dd4-467e-e5ff-642b-371add518b40",
      "lease_duration": 315359999,
      "renewable": false,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "serial": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/root/sign-intermediate
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Uses the configured CA certificate to issue a certificate with appropriate
    values for acting as an intermediate CA. Distribution points use the values
    set via `config/urls`. Values set in the CSR are ignored unless
    `use_csr_values` is set to true, in which case the values from the CSR are
    used verbatim.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/root/sign-intermediate`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">csr</span>
        <span class="param-flags">required</span>
        The PEM-encoded CSR.
      </li>
      <li>
        <span class="param">common_name</span>
        <span class="param-flags">required</span>
        The requested CN for the certificate.
      </li>
      <li>
        <span class="param">alt_names</span>
        <span class="param-flags">optional</span>
        Requested Subject Alternative Names, in a comma-delimited list. These
        can be host names or email addresses; they will be parsed into their
        respective fields.
      </li>
      <li>
        <span class="param">ip_sans</span>
        <span class="param-flags">optional</span>
        Requested IP Subject Alternative Names, in a comma-delimited list.
      </li>
      <li>
      <span class="param">ttl</span>
      <span class="param-flags">optional</span>
        Requested Time To Live (after which the certificate will be expired).
        This cannot be larger than the mount max (or, if not set, the system
        max).
      </li>
      <li>
      <span class="param">format</span>
      <span class="param-flags">optional</span>
        Format for returned data. Can be `pem`, `der`, or `pem_bundle`;
        defaults to `pem`. If `der`, the output is base64 encoded. If
        `pem_bundle`, the `certificate` field will contain the certificate and
        issuing CA, concatenated.
      </li>
      <li>
        <span class="param">max_path_length</span>
        <span class="param-flags">optional</span>
        If set, the maximum path length to encode in the generated certificate.
        Defaults to `-1`, which means no limit.  unless the signing certificate
        has a maximum path length set, in which case the path length is set to
        one less than that of the signing certificate.  A limit of `0` means a
        literal path length of zero.
      </li>
      <li>
        <span class="param">exclude_cn_from_sans</span>
        <span class="param-flags">optional</span>
        If set, the given `common_name` will not be included in DNS or Email
        Subject Alternate Names (as appropriate). Useful if the CN is not a
        hostname or email address, but is instead some human-readable
        identifier.
      </li>
      <li>
        <span class="param">use_csr_values</span>
        <span class="param-flags">optional</span>
        If set to `true`, then: 1) Subject information, including names and
        alternate names, will be preserved from the CSR rather than using the
        values provided in the other parameters to this path; 2) Any key usages
        (for instance, non-repudiation) requested in the CSR will be added to
        the basic set of key usages used for CA certs signed by this path; 3)
        Extensions requested in the CSR will be copied into the issued
        certificate.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "pki/root/sign-intermediate/bc23e3c6-8dcd-48c6-f3af-dd2db7f815c2",
      "renewable": false,
      "lease_duration": 21600,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n",
        "serial": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/sign/
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Signs a new certificate based upon the provided CSR and the supplied
    parameters, subject to the restrictions contained in the role named in the
    endpoint. The issuing CA certificate is returned as well, so that only the
    root CA need be in a client's trust store.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/sign/<role name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">csr</span>
        <span class="param-flags">required</span>
        The PEM-encoded CSR.
      </li>
      <li>
        <span class="param">common_name</span>
        <span class="param-flags">required</span>
        The requested CN for the certificate. If the CN is allowed by role
        policy, it will be issued.
      </li>
      <li>
        <span class="param">alt_names</span>
        <span class="param-flags">optional</span>
        Requested Subject Alternative Names, in a comma-delimited list. These
        can be host names or email addresses; they will be parsed into their
        respective fields. If any requested names do not match role policy, the
        entire request will be denied.
      </li>
      <li>
        <span class="param">ip_sans</span>
        <span class="param-flags">optional</span>
        Requested IP Subject Alternative Names, in a comma-delimited list. Only
        valid if the role allows IP SANs (which is the default).
      </li>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        Requested Time To Live. Cannot be greater than the role's `max_ttl`
        value. If not provided, the role's `ttl` value will be used. Note that
        the role values default to system values if not explicitly set.
      </li>
      <li>
        <span class="param">format</span>
        <span class="param-flags">optional</span>
        Format for returned data. Can be `pem`, `der`, or `pem_bundle`;
        defaults to `pem`. If `der`, the output is base64 encoded. If
        `pem_bundle`, the `certificate` field will contain the certificate and
        issuing CA, concatenated.
      </li>
      <li>
        <span class="param">exclude_cn_from_sans</span>
        <span class="param-flags">optional</span>
        If set, the given `common_name` will not be included in DNS or Email
        Subject Alternate Names (as appropriate). Useful if the CN is not a
        hostname or email address, but is instead some human-readable
        identifier.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "pki/sign/test/7ad6cfa5-f04f-c62a-d477-f33210475d05",
      "renewable": false,
      "lease_duration": 21600,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n",
        "serial": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/sign-verbatim
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Signs a new certificate based upon the provided CSR. Values are taken
    verbatim from the CSR; the _only_ restriction is that this endpoint will
    refuse to issue an intermediate CA certificate (see the
    `/pki/root/sign-intermediate` endpoint for that functionality.) _This is a
    potentially dangerous endpoint and only highly trusted users should
    have access._
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/sign-verbatim`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">csr</span>
        <span class="param-flags">required</span>
        The PEM-encoded CSR.
      </li>
      <li>
      <span class="param">ttl</span>
      <span class="param-flags">optional</span>
        Requested Time To Live. Cannot be greater than the mount's `max_ttl`
        value. If not provided, the mount's `ttl` value will be used, which
        defaults to system values if not explicitly set.
      </li>
      <li>
      <span class="param">format</span>
      <span class="param-flags">optional</span>
        Format for returned data. Can be `pem`, `der`, or `pem_bundle`;
        defaults to `pem`. If `der`, the output is base64 encoded. If
        `pem_bundle`, the `certificate` field will contain the certificate and
        issuing CA, concatenated.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "pki/sign-verbatim/7ad6cfa5-f04f-c62a-d477-f33210475d05",
      "renewable": false,
      "lease_duration": 21600,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n",
        "serial": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/tidy
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows tidying up the backend storage and/or CRL by removing certificates
    that have expired and are past a certain buffer period beyond their
    expiration time.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/tidy`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">tidy_cert_store</span>
        <span class="param-flags">optional</span>
        Whether to tidy up the certificate store. Defaults to `false`.
      </li>
      <li>
      <span class="param">tidy_revocation_list</span>
      <span class="param-flags">optional</span>
        Whether to tidy up the revocation list (CRL). Defaults to `false`.
      </li>
      <li>
      <span class="param">safety_buffer</span>
      <span class="param-flags">optional</span>
        A duration (given as an integer number of seconds or a string; defaults
        to `72h`) used as a safety buffer to ensure certificates are not
        expunged prematurely; as an example, this can keep certificates from
        being removed from the CRL that, due to clock skew, might still be
        considered valid on other hosts. For a certificate to be expunged, the
        time must be after the expiration time of the certificate (according to
        the local clock) plus the duration of `safety_buffer`.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` status code.
  </dd>
</dl>
