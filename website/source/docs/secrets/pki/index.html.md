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

**N.B.**: This information is valid for 0.4+.

This backend serves CRLs from a predictable location, but it is not possible
for the backend to know where it is running. Therefore, you must configure
desired URLs for the issuing certificate, CRL distribution points, and OCSP
servers manually using the `config/urls` endpoint. It is supported to have more
than one of each of these.

### No OCSP support, yet

Vault's architecture does not currently allow for a binary protocol such as
OCSP to be supported by a backend. As such, you should configure your software
to use CRLs for revocation information, with a caching lifetime that feels good
to you. Since you are following the advice above about keeping lifetimes short
(right?), CRLs should not grow too large.

If OCSP is important to you, you can configure alternate CRL and/or OCSP
servers using `config/urls`.

If you are using issued certificates for client authentication to Vault, note
that as of 0.4, the `cert` authentication endpoint supports being pushed CRLs.

## Quick Start

**N.B.**: This quick start assumes that you are pushing a CA certificate and
private key in. Since 0.4 there are many different methods of getting CA
information into the backend according to your needs, so please note that this
is simply one example.

The first step to using the PKI backend is to mount it. Unlike the `generic`
backend, the `pki` backend is not mounted by default.

```text
$ vault mount pki
Successfully mounted 'pki' at 'pki'!
```

Next, Vault must be configured with a CA certificate and associated private key. This is done by writing the contents of a file or *stdin*:

```text
$ vault write pki/config/ca \
    pem_bundle="@ca_bundle.pem"
Success! Data written to: pki/config/ca
```

or

```
$ cat bundle.pem | vault write pki/config/ca pem_bundle="-"
Success! Data written to: pki/config/ca
```

Although in this example the value being piped into *stdin* could be passed
directly into the Vault CLI command, a more complex usage might be to use
[Ansible](http://www.ansible.com) to securely store the certificate and private
key in an `ansible-vault` file, then have an `ansible-playbook` command decrypt
this value and pass it in to Vault.

The next step is to configure a role. A role is a logical name that maps to a
policy used to generated those credentials. For example, let's create an
"example-dot-com" role:

```text
$ vault write pki/roles/example-dot-com \
    allowed_base_domain="example.com" \
    allow_subdomains="true" max_ttl="72h"
Success! Data written to: pki/roles/example-dot-com
```

By writing to the `roles/example-dot-com` path we are defining the
`example-dot-com` role. To generate a new set of credentials, we simply write
to the `issue` endpoint with that role name: Vault is now configured to create
and manage certificates!

```text
$ vault write pki/issue/example-dot-com \
    common_name=blah.example.com
Key             Value
lease_id        pki/issue/example-dot-com/819393b5-e1a1-9efd-b72f-4dc3a1972e31
lease_duration  259200
lease_renewable false
certificate     -----BEGIN CERTIFICATE-----
MIIECDCCAvKgAwIBAgIUXmLrLkTdBIOOIYg2/BXO7docKfUwCwYJKoZIhvcNAQEL
...
az3gfwlOqVTdgi/ZVAtIzhSEJ0OY136bq4NOaw==
-----END CERTIFICATE-----
issuing_ca      -----BEGIN CERTIFICATE-----
MIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV
...
-----END CERTIFICATE-----
private_key     -----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0cczc7Y2yIu7aD/IaDi23Io+tvvDS9XaXXDUFW1kqd58P83r
...
3xhCNnZ3CMQaM2I48sloVK/XoikMLb5MZwOUQn/V+TrhWP4Lu7qD
-----END RSA PRIVATE KEY-----
serial_number   5e:62:eb:2e:44:dd:04:83:8e:21:88:36:fc:15:ce:ed:da:1c:29:f5
```

Note that this is a write, not a read, to allow values to be passed in at
request time.

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

### /pki/config/ca/set
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows submitting the CA information via a PEM file containing the CA
    certificate and its private key, concatenated.  If you generated an
    intermediate CA CSR and received a signed certificate, you do not need to
    include the private key in the PEM file. <br /><br />The information can
    be provided from a file via a `curl` command similar to the following:<br/>

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
  <dd>`/pki/config/ca/set`</dd>

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

### /pki/config/ca/generate/intermediate
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new private key and a CSR for signing. If using Vault as a
    root, and for many other CAs, the various parameters on the final
    certificate are set at signing time and may or may not honor the parameters
    set here. If the path ends with `external`, the private key will be
    returned in the response; if it is `internal` the private key will not be
    returned and *cannot be retrieved later*. <br /><br />This is mostly meant
    as a helper function, and not all possible parameters that can be set in a
    CSR are supported.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/config/ca/generate/intermediate/[exported|internal]`</dd>

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
        Format for returned data. Can be `pem` or `der`; defaults to `pem`. If
        `der`, the output is base64 encoded.
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
        },
      "auth": null
    }
    ```

  </dd>
</dl>


### /pki/config/ca/generate/root
#### POST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Generates a new self-signed CA certificate and private key. If the path
    ends with `external`, the private key will be returned in the response; if
    it is `internal` the private key will not be returned and *cannot be
    retrieved later*. Distribution points use the values set via `config/urls`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/config/ca/generate/root/[exported|internal]`</dd>

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
        Format for returned data. Can be `pem` or `der`; defaults to `pem`. If
        `der`, the output is base64 encoded.
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
        Defaults to `-1`, which means no limit. A limit of `0` means a literal
        path length of zero.
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
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "serial": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
      "auth": null
    }
    ```

  </dd>
</dl>

### /pki/config/ca/sign
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
  <dd>`/pki/config/ca/sign`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
      <li>
        <span class="param">csr</span>
        <span class="param-flags">required</span>
        The PEM-encoded CSR.
      </li>
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
        Format for returned data. Can be `pem` or `der`; defaults to `pem`. If
        `der`, the output is base64 encoded.
      </li>
      <li>
        <span class="param">max_path_length</span>
        <span class="param-flags">optional</span>
        If set, the maximum path length to encode in the generated certificate.
        Defaults to `-1`, which means no limit. A limit of `0` means a literal
        path length of zero.
      </li>
      <li>
        <span class="param">use_csr_values</span>
        <span class="param-flags">optional</span>
        If set to `true`, then: 1) Subject information, including names and
        alternate names, will be preserved from the CSR rather than using the
        values provided in the other parameters to this path; 2) Any key usages
        (for instance, non-repudiation) requested in the CSR will be added to
        the basic set of key usages used for CA certs signed by this path.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id": "pki/config/ca/sign/bc23e3c6-8dcd-48c6-f3af-dd2db7f815c2",
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

### /pki/config/urls
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Fetch the URLs to be encoded in generated certificates.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

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

### /pki/config/crl
#### GET

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
  <dd>`/pki/issue/<name>`</dd>

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
      <span class="param">format</span>
      <span class="param-flags">optional</span>
        Format for returned data. Can be `pem` or `der`; defaults to `pem`. If
        `der`, the output is base64 encoded.
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
        "serial": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
        },
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
    `allowed_base_domain`, `allow_token_displayname`, `allow_subdomains`, and
    `allow_any_name` attributes are additive; between them nearly and across
    multiple roles nearly any issuing policy can be accommodated.
    `server_flag`, `client_flag`, and `code_signing_flag` are additive as well.
    If a client requests a certificate that is not allowed by the CN policy in
    the role, the request is denied.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/pki/roles/<name>`</dd>

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
        <span class="param">allowed_base_domain</span>
        <span class="param-flags">optional</span>
        **N.B.**: In 0.4+, the meaning of this value has changed, although the
        name is kept for backwards compatibility.<br /><br />Designates the
        base domain of the role. This is used with the `allow_base_domain` and
        `allow_subdomains` options. There is no default.
      </li>
      <li>
        <span class="param">allow_base_domain</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates matching the value of the
        actual base domain. Defaults to false.
      </li>
      <li>
        <span class="param">allow_token_displayname</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates matching the value of Display
        Name from the requesting token. Remember, this stacks with the other CN
        options, including `allow_subdomains`. Defaults to `false`.
      </li>
      <li>
        <span class="param">allow_subdomains</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates with CNs that are subdomains
        of the CNs allowed by the other role options. _This includes wildcard
        subdomains._ For example, an `allowed_base_domain` value of
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
        host part of email addresses. Defaults to `false`.
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
  <dd>`/pki/roles/<name>`</dd>

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
        "allow_token_displayname": false,
        "allowed_base_domain": "example.com",
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
  <dd>`/pki/roles/<name>`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>
    A `204` response code.
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
  <dd>`/pki/sign/<name>`</dd>

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
      <span class="param">format</span>
      <span class="param-flags">optional</span>
        Format for returned data. Can be `pem` or `der`; defaults to `pem`. If
        `der`, the output is base64 encoded.
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
