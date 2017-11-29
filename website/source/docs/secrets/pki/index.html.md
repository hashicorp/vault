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
---             -----
certificate     -----BEGIN CERTIFICATE-----
MIIDNTCCAh2gAwIBAgIUE7JOMnCYSNeHmEkXIV38pbgQB+cwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wHhcNMTcxMTI5MTQ1MjE1WhcNMjcx
MTI3MTQ1MjQ1WjAWMRQwEgYDVQQDEwtteXZhdWx0LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMzROgEADqYx16LwaSwG+GAbAOIGRsYcillLbgAb
7IbkTP7y95o1XiGotYCUipkSdXelRcXPuV3LjaFjSLbJk8Zcka3I5edSNT5ILcMt
R5QkMsTKtByvpyCNZTF+FHFl9NRZFCGT88KcC2MhUQ8Y/tT0AA4yw0WgzRKNLBj+
Knuxu8xz5pbZeCXTq2gUhjSym8S4yx8I8DNtDBSJTjJXZb+tRMTQx+hQcWQqRx4b
wCZW6EOfbqTk5IP6frdnd5E+GpX9hHMF+C4ULI3sAULqgUTjokiUr3cIoDKrw8+d
I5NV/TKrDX415A6zlMoUcg5quWLNkPxQPy7n5ds8nht6LMUCAwEAAaN7MHkwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHIkFpNXoTG5
Ub/wp0VDdCih4yoDMB8GA1UdIwQYMBaAFHIkFpNXoTG5Ub/wp0VDdCih4yoDMBYG
A1UdEQQPMA2CC215dmF1bHQuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQBeUAyqhicD
VnlfWvKXaADjN3iiboiR/iKO/tKni7pc5QfNdah6IxT0gJuTNXQtNlhVDfY2xD1y
MDhtq7QO82J5NyWA75rQBzBIKFhMyL3s7GYY4kKiGmA5i975vpfh2BdHkAIdXWpZ
1CmdW+D6+9uCzK/2SgTKulJH0zEF9qokQXVOCEzQDv1GsmiapWJqeObin0H9Ff0r
RwKcWPNYwdnmzL9WtMd+VjFIEly0NtEoX6lwWyDlLRTxe9etoML3y9COxEKkgOxB
0SkP2MH511ka9tX1WjQHHCmdwLHLJ6kXMqPQnI70DyZkIzgQWmqhiKQUJiQA8jlI
Az3toI2X3ipS
-----END CERTIFICATE-----
expiration      1827327165
issuing_ca      -----BEGIN CERTIFICATE-----
MIIDNTCCAh2gAwIBAgIUE7JOMnCYSNeHmEkXIV38pbgQB+cwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wHhcNMTcxMTI5MTQ1MjE1WhcNMjcx
MTI3MTQ1MjQ1WjAWMRQwEgYDVQQDEwtteXZhdWx0LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMzROgEADqYx16LwaSwG+GAbAOIGRsYcillLbgAb
7IbkTP7y95o1XiGotYCUipkSdXelRcXPuV3LjaFjSLbJk8Zcka3I5edSNT5ILcMt
R5QkMsTKtByvpyCNZTF+FHFl9NRZFCGT88KcC2MhUQ8Y/tT0AA4yw0WgzRKNLBj+
Knuxu8xz5pbZeCXTq2gUhjSym8S4yx8I8DNtDBSJTjJXZb+tRMTQx+hQcWQqRx4b
wCZW6EOfbqTk5IP6frdnd5E+GpX9hHMF+C4ULI3sAULqgUTjokiUr3cIoDKrw8+d
I5NV/TKrDX415A6zlMoUcg5quWLNkPxQPy7n5ds8nht6LMUCAwEAAaN7MHkwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHIkFpNXoTG5
Ub/wp0VDdCih4yoDMB8GA1UdIwQYMBaAFHIkFpNXoTG5Ub/wp0VDdCih4yoDMBYG
A1UdEQQPMA2CC215dmF1bHQuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQBeUAyqhicD
VnlfWvKXaADjN3iiboiR/iKO/tKni7pc5QfNdah6IxT0gJuTNXQtNlhVDfY2xD1y
MDhtq7QO82J5NyWA75rQBzBIKFhMyL3s7GYY4kKiGmA5i975vpfh2BdHkAIdXWpZ
1CmdW+D6+9uCzK/2SgTKulJH0zEF9qokQXVOCEzQDv1GsmiapWJqeObin0H9Ff0r
RwKcWPNYwdnmzL9WtMd+VjFIEly0NtEoX6lwWyDlLRTxe9etoML3y9COxEKkgOxB
0SkP2MH511ka9tX1WjQHHCmdwLHLJ6kXMqPQnI70DyZkIzgQWmqhiKQUJiQA8jlI
Az3toI2X3ipS
-----END CERTIFICATE-----
serial_number   13:b2:4e:32:70:98:48:d7:87:98:49:17:21:5d:fc:a5:b8:10:07:e7
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
Key                 Value
---                 -----
certificate         -----BEGIN CERTIFICATE-----
MIIDvzCCAqegAwIBAgIUeLWwG/YiWlT3HCOlFsSJv89KswUwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wHhcNMTcxMTI5MTQ1MzI1WhcNMTcx
MjAyMTQ1MzU1WjAbMRkwFwYDVQQDExBibGFoLmV4YW1wbGUuY29tMIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwpNktB/nBtwZrEus2tWcMISd4KqTgiXC
XZMxSG7MSSDU872TD1pbBasxLmCyl+ZgKORA/tH6xqKqJdfujsTdBBqAF/7zsMDw
wMR7uViZhtg5zAR1qrpmYbal4o92INBS2HDQ8irFD6uswpZqsJKpfHTl/BI6VChK
d+eLeT89I89XmE7LKaP5RN/W4AuPvQy5wB5TtU4vqQPmQpoHWc8qjxStSBg9tsEp
k1VcnA0drr+enjyzHfOV+CAwuPHPsrsAoEqiwxXVoo8TK+LQn35U5w1I/SjwrvS0
bWJf9Fh4bnC5fonSKmohY/yVvKJuCDM2R9Tg86NCBMSnhK2Nn1tw6QIDAQABo4H/
MIH8MA4GA1UdDwEB/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
AwIwHQYDVR0OBBYEFDNWOkqGCzBpV2TU5EqzYSVhzyhhMB8GA1UdIwQYMBaAFHIk
FpNXoTG5Ub/wp0VDdCih4yoDMDsGCCsGAQUFBwEBBC8wLTArBggrBgEFBQcwAoYf
aHR0cDovLzEyNy4wLjAuMTo4MjAwL3YxL3BraS9jYTAbBgNVHREEFDASghBibGFo
LmV4YW1wbGUuY29tMDEGA1UdHwQqMCgwJqAkoCKGIGh0dHA6Ly8xMjcuMC4wLjE6
ODIwMC92MS9wa2kvY3JsMA0GCSqGSIb3DQEBCwUAA4IBAQBBKJnvvInuLtB5N8nH
G9F+uVkTFAgJw4FeaV+zD0KG2Kbu3En4C6+pucxDvSTJITfoAm/EyomlkiB4G+Ql
B8RSHehaL8fG6X4Cz7xcoab8Y38TY93h4k2w4aVi7Rpu2PPpwrms941F1CmKcfqH
4gjPi7t35jZ3jbke45qiF02Xorzn8l1ifmcf7asOdlUymaoxwQb/StfdqI8tNpSL
yt+XXhwZKgdnLf9hiBKa5RtXwOLMRp6o7MPZiNrzLO/pR4YQlk45v4h6gzgRA8VX
HHEz7P/xMjRz9cRalody/khrc0pjZ/K7Tz1V5MwEWrlm6GoNpy4/SwAZgxZCB4Jt
aADC
-----END CERTIFICATE-----
issuing_ca          -----BEGIN CERTIFICATE-----
MIIDNTCCAh2gAwIBAgIUE7JOMnCYSNeHmEkXIV38pbgQB+cwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wHhcNMTcxMTI5MTQ1MjE1WhcNMjcx
MTI3MTQ1MjQ1WjAWMRQwEgYDVQQDEwtteXZhdWx0LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAMzROgEADqYx16LwaSwG+GAbAOIGRsYcillLbgAb
7IbkTP7y95o1XiGotYCUipkSdXelRcXPuV3LjaFjSLbJk8Zcka3I5edSNT5ILcMt
R5QkMsTKtByvpyCNZTF+FHFl9NRZFCGT88KcC2MhUQ8Y/tT0AA4yw0WgzRKNLBj+
Knuxu8xz5pbZeCXTq2gUhjSym8S4yx8I8DNtDBSJTjJXZb+tRMTQx+hQcWQqRx4b
wCZW6EOfbqTk5IP6frdnd5E+GpX9hHMF+C4ULI3sAULqgUTjokiUr3cIoDKrw8+d
I5NV/TKrDX415A6zlMoUcg5quWLNkPxQPy7n5ds8nht6LMUCAwEAAaN7MHkwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHIkFpNXoTG5
Ub/wp0VDdCih4yoDMB8GA1UdIwQYMBaAFHIkFpNXoTG5Ub/wp0VDdCih4yoDMBYG
A1UdEQQPMA2CC215dmF1bHQuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQBeUAyqhicD
VnlfWvKXaADjN3iiboiR/iKO/tKni7pc5QfNdah6IxT0gJuTNXQtNlhVDfY2xD1y
MDhtq7QO82J5NyWA75rQBzBIKFhMyL3s7GYY4kKiGmA5i975vpfh2BdHkAIdXWpZ
1CmdW+D6+9uCzK/2SgTKulJH0zEF9qokQXVOCEzQDv1GsmiapWJqeObin0H9Ff0r
RwKcWPNYwdnmzL9WtMd+VjFIEly0NtEoX6lwWyDlLRTxe9etoML3y9COxEKkgOxB
0SkP2MH511ka9tX1WjQHHCmdwLHLJ6kXMqPQnI70DyZkIzgQWmqhiKQUJiQA8jlI
Az3toI2X3ipS
-----END CERTIFICATE-----
private_key         -----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAwpNktB/nBtwZrEus2tWcMISd4KqTgiXCXZMxSG7MSSDU872T
D1pbBasxLmCyl+ZgKORA/tH6xqKqJdfujsTdBBqAF/7zsMDwwMR7uViZhtg5zAR1
qrpmYbal4o92INBS2HDQ8irFD6uswpZqsJKpfHTl/BI6VChKd+eLeT89I89XmE7L
KaP5RN/W4AuPvQy5wB5TtU4vqQPmQpoHWc8qjxStSBg9tsEpk1VcnA0drr+enjyz
HfOV+CAwuPHPsrsAoEqiwxXVoo8TK+LQn35U5w1I/SjwrvS0bWJf9Fh4bnC5fonS
KmohY/yVvKJuCDM2R9Tg86NCBMSnhK2Nn1tw6QIDAQABAoIBAQCiS17L/3AsAJpJ
ZDWhslq8/WDSbHTtMaGVe5i32fL7bC8zvwRS4gLaD5jOHknY+YdrlDfCVFLgK/3P
4vRQkVPogFInsbiBze0CXOD2EDi+iMxsp6ud5CxRhI+JEjBt3lW7wx4FYDdOgttu
0xfaf/punPLX9jxAxfqXWMK1N1I/8tVYWnt/EA7ZTazp+z3+kDxSLsnxF2Y+c3Bq
OB8Gw397P6jsMraoFpEdkymUYy+/GqeJ8sKZ5l4DDfcH+++oMfdrFE0Sr90fQf1x
aD1cf+GbvjaCYxrSTxymCtFh7Xk++axSvFfkmMtHgFDTrgSzVVGzio4m/zJl/Nq+
c+lTkvIhAoGBAOoHEuOviC5WVz6cUneZNFFiM31MosOQ1HEAcIvR7OvU5HII426h
vSmx8Gki0snc6gcyQ9jmEM+wwkf0Q83t74vU0jE9ll8kzADu7jC2blllCvfT7lKs
ZpNZ3W67Si8w8GP6BBaDa4PCMLBmv3fBYiSZJvtk85uvQIAmhgtf190rAoGBANTY
FOKGOwxqJc6qah4gQcitvb9eSRQ/oN7iSTWzaBQ/YaDo1PWe5pVsG6xKrHbzirRv
jteJVKtVqPE+ukJ9RBibT2sVCqWqAzoo/1hxJiu1aYXyX08LnNuhByyVQ1wG5plF
ikqdgz5IFnzTg0+8v2/UFk3toq4j1z+OCoyXRmg7AoGAHk9cOvD5CkdUdV95rtPA
2umFEa1jR0Dyws/zw6gkr0abb8mG60U3YrcRFAzWkB50kQoJj4X8l2mlP/x666jt
ZYbi0k3Ps/LoGRbY8qYuFJXpnb9tFngNsPfqnfTT3tjPyaMP9HqA6ke0VqR4F+KL
+4F6cwTYKEnCaNaUddSr+JECgYEAh/2bsnQTLEZx646kiKURgvfHQYsrZB2XWnD4
V7BOMomghh/dWSXyq8vMDpQTh1jp6YlRmdLr3yC29ZSfizXgGVy6LG/gQqLStwlU
xJxeyBR73JJUZPvFd+p12/1ucVETayCsUCo9ncCPZaf6wSqWogu/SIEprNvHfprx
kIxi9tsCgYEAg8aBSCEBgaUrW6yvHb6PzJJTbwzrtcveInrKfelESlr2mu7wQmew
JizJZLZKjg8dV0iLQ0fpy9PHaLYSiYiH0utZ5PVKPfuWBb3jYDDFGuyzYIpnlCVq
psMmsg8wwknFBVLNkFUofPt1BJH2TAi46CWK2Qy6yHmeHwSEiz4857I=
-----END RSA PRIVATE KEY-----
private_key_type    rsa
serial_number       78:b5:b0:1b:f6:22:5a:54:f7:1c:23:a5:16:c4:89:bf:cf:4a:b3:05
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
