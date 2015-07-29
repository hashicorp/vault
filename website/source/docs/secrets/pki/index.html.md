---
layout: "docs"
page_title: "Secret Backend: PKI"
sidebar_current: "docs-secrets-pki"
description: |-
  The PKI secret backend for Vault generates TLS certificates.
---

# PKI Secret Backend

Name: `pki`

The PKI secret backend for Vault generates X.509 certificates dynamically based on configured roles. This means services can get certificates needed for both client and server authentication without going through the usual manual process of generating a private key and CSR, submitting to a CA, and waiting for a verification and signing process to complete. Vault's built-in authentication and authorization mechanisms provide the verification functionality.

By keeping leases relatively short, revocations are less likely to be needed, keeping CRLs short and helping the backend scale to large workloads. This in turn allows each instance of a running application to have a unique certificate, eliminating sharing and the accompanying pain of revocation and rollover.

In addition, by allowing revocation to mostly be forgone, this backend allows for ephemeral certificates; certificates can be fetched and stored in memory upon application startup and discarded upon shutdown, without ever being written to disk.

This page will show a quick start for this backend. For detailed documentation on every path, use `vault path-help` after mounting the backend.

## Considerations

To successfully deploy this backend, there are a number of important considerations to be aware of, as well as some preparatory steps that should be undertaken. You should read all of these *before* using this backend or generating the CA to use with this backend.

### Never use root CAs

Vault storage is secure, but not as secure as a piece of paper in a bank vault. It is, after all, networked software. Your long-lived self-signed root CA's private key should instead be used to issue a shorter-lived intermediate CA certificate, and this is what you should put into Vault. This aligns with industry best practices.

### One CA Certificate, One Backend

In order to vastly simplify both the configuration and codebase of the PKI backend, only one CA certificate is allowed per backend. If you want to issue certificates from multiple CAs, mount the PKI backend at multiple mount points with separate CA certificates in each.

This also provides a convenient method of switching to a new CA certificate while keeping CRLs valid from the old CA certificate; simply mount a new backend and issue from there.

### Keep certificate lifetimes short, for CRL's sake

This backend aligns with Vault's philosophy of short-lived secrets. As such it is not expected that CRLs will grow large; the only place a private key is ever returned is to the requesting client (this backend does *not* store generated private keys). In most cases, if the key is lost, the certificate can simply be ignored, as it will expire shortly.

If a certificate must truly be revoked, the normal Vault revocation function can be used; alternately a root token can be used to revoke the certificate using the certificate's serial number. Any revocation action will cause the CRL to be regenerated. When the CRL is regenerated, any expired certificates are removed from the CRL (and any revoked, expired certificate are removed from backend storage).

This backend does not support multiple CRL endpoints with sliding date windows; often such mechanisms will have the transition point a few days apart, but this gets into the expected realm of the actual certificate validity periods issued from this backend. A good rule of thumb for this backend would be to simply not issue certificates with a validity period greater than your maximum comfortable CRL lifetime. Alternately, you can control CRL caching behavior on the client to ensure that checks happen more often.

Often multiple endpoints are used in case a single CRL endpoint is down so that clients don't have to figure out what to do with a lack of response. Run Vault in HA mode, and the CRL endpoint should be available even if a particular node is down.

### You must configure CRL information *in advance*

This backend serves CRLs from a predictable location. That location must be encoded into your CA certificate if you want to allow applications to use the CRL endpoint encoded in certificates to find the CRL. Instructions for doing so are below. If you need to adjust this later, you will have to generate a new CA certificate using the same private key if you want to keep validity for already-issued certificates.

### No OCSP support, yet

Vault's architecture does not currently allow for a binary protocol such as OCSP to be supported by a backend. As such, you should configure your software to use CRLs for revocation information, with a caching lifetime that feels good to you. Since you are following the advice above about keeping lifetimes short (right?), CRLs should not grow too large.

## Quick Start

### CA certificate

In order for this backend to serve CRL information at the expected location, you will need to generate your CA certificate with this information. For OpenSSL, this means putting a value in the CA section with the appropriate URL; in this example the PKI backend is mounted at `pki`:

```text
crlDistributionPoints = URI:https://vault.example.com:8200/v1/pki/crl
```

Adjust the URI as appropriate.

### Vault

The first step to using the PKI backend is to mount it. Unlike the `generic` backend, the `pki` backend is not mounted by default.

```text
$ vault mount pki
Successfully mounted 'pki' at 'pki'!
```

Next, Vault must be configured with a root certificate and associated private key. This is done by writing the contents of a file or *stdin*:

```text
$ vault write pki/config/ca pem_bundle="@ca_bundle.pem"
Success! Data written to: pki/config/ca
```

or

```
$ cat bundle.pem | vault write pki/config/ca pem_bundle="-"
Success! Data written to: pki/config/ca
```

Although in this example the value being piped into *stdin* could be passed directly into the Vault CLI command, a more complex usage might be to use [Ansible](http://www.ansible.com) to securely store the certificate and private key in an `ansible-vault` file, then have an `ansible-playbook` command decrypt this value and pass it in to Vault.

The next step is to configure a role. A role is a logical name that maps to a policy used to generated those credentials. For example, let's create an "example-dot-com" role:

```text
$ vault write pki/roles/example-dot-com \
    allowed_base_domain="example.com" \
    allow_subdomains="true" lease_max="72h"
Success! Data written to: pki/roles/example-dot-com
```

By writing to the `roles/example-dot-com` path we are defining the `example-dot-com` role. To generate a new set of credentials, we simply write to the `issue` endpoint with that role name: Vault is now configured to create and manage certificates!

```text
$ vault write pki/issue/example-dot-com common_name=blah.example.com
Key            	Value
lease_id       	pki/issue/example-dot-com/819393b5-e1a1-9efd-b72f-4dc3a1972e31
lease_duration 	259200
lease_renewable	false
certificate    	-----BEGIN CERTIFICATE-----
MIIECDCCAvKgAwIBAgIUXmLrLkTdBIOOIYg2/BXO7docKfUwCwYJKoZIhvcNAQEL
...
az3gfwlOqVTdgi/ZVAtIzhSEJ0OY136bq4NOaw==
-----END CERTIFICATE-----
issuing_ca      -----BEGIN CERTIFICATE-----
MIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV
...
-----END CERTIFICATE-----
private_key    	-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0cczc7Y2yIu7aD/IaDi23Io+tvvDS9XaXXDUFW1kqd58P83r
...
3xhCNnZ3CMQaM2I48sloVK/XoikMLb5MZwOUQn/V+TrhWP4Lu7qD
-----END RSA PRIVATE KEY-----
serial         	5e:62:eb:2e:44:dd:04:83:8e:21:88:36:fc:15:ce:ed:da:1c:29:f5
```

Note that this is a write, not a read, to allow values to be passed in at request time.

Vault has now generated a new set of credentials using the `example-dot-com` role configuration. Here we see the dynamically generated private key and certificate. The issuing CA certificate is returned as well.

Using ACLs, it is possible to restrict using the pki backend such that trusted operators can manage the role definitions, and both users and applications are restricted in the credentials they are allowed to read.

If you get stuck at any time, simply run `vault path-help pki` or with a subpath for interactive help output.

## API

### /pki/ca(/pem)
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves the CA certificate *in raw DER-encoded form*.
    This is a bare endpoint that does not return a
    standard Vault data structure. If `/pem` is added to the
    endpoint, the CA certificate is returned in PEM format.
    <br /><br />This is an unauthenticated endpoint.
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
    Retrieves one of a selection of certificates. Valid values: `ca`
    for the CA certificate, `crl` for the current CRL, or a serial
    number in either hyphen-separated or colon-separated octal format.
    This endpoint returns the certificate in PEM formatting in the
    `certificate` key of the JSON object.
    <br /><br />This is an unauthenticated endpoint.
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
    A PEM file containing the issuing CA certificate
    and its private key, concatenated.
    <br /><br />This is a root-protected endpoint.
    <br /><br />The information can be provided from a file via a `curl`
    command similar to the following:<br/>

    ```text
    curl -X POST --data "@cabundle.json" http://127.0.0.1:8200/v1/pki/config/ca -H X-Vault-Token:06b9d...
    ```

    Note that if you provide the data through the HTTP API it must be
    JSON-formatted, with newlines replaced with `\n`, like so:

    ```text
    { "pem_bundle": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END CERTIFICATE-----" }
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
  <br /><br />This is a root-protected endpoint.
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
    Generates a new set of credentials (private key and
    certificate) based on the named role. The issuing CA
    certificate is returned as well, so that only the root CA
    need be in a client's trust store.
    <br /><br />*The private key is _not_ stored.
    If you do not save the private key, you will need to
    request a new certificate.*
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
        The requested CN for the certificate. If the CN is allowed
        by role policy, it will be issued.
      </li>
      <li>
        <span class="param">alt_names</span>
        <span class="param-flags">optional</span>
        Requested Subject Alternative Names, in a comma-delimited
        list. If any requested names do not match role policy,
        the entire request will be denied.
      </li>
      <li>
        <span class="param">ip_sans</span>
        <span class="param-flags">optional</span>
        Requested IP Subject Alternative Names, in a comma-delimited
        list. Only valid if the role allows IP SANs (which is the
        default).
      </li>
      <li>
      <span class="param">lease</span>
      <span class="param-flags">optional</span>
        Requested lease time. Cannot be greater than the role's
        `lease_max` parameter. If not provided, the role's `lease`
        value will be used.
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
    <br /><br />This is a root-protected endpoint.
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
    Creates or updates the role definition. Note that
    the `allowed_base_domain`, `allow_token_displayname`,
    `allow_subdomains`, and `allow_any_name` attributes
    are additive; between them nearly and across multiple
    roles nearly any issuing policy can be accommodated.
    `server_flag`, `client_flag`, and `code_signing_flag`
    are additive as well. If a client requests a
    certificate that is not allowed by the CN policy in
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
        <span class="param">lease</span>
        <span class="param-flags">optional</span>
        The lease value provided as a string duration
        with time suffix. Hour is the largest suffix.
        If not set, uses the value of `lease_max`.
      </li>
      <li>
        <span class="param">lease_max</span>
        <span class="param-flags">required</span>
        The maximum lease value provided as a string duration
        with time suffix. Hour is the largest suffix.
      </li>
      <li>
        <span class="param">allow_localhost</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates for `localhost`
        as one of the requested common names. This is useful
        for testing and to allow clients on a single host to
        talk securely.
        Defaults to true.
      </li>
      <li>
        <span class="param">allowed_base_domain</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates for subdomains
        directly off of this base domain. _This includes the
        wildcard subdomain._ For instance, a base_domain of
        `example.com` allows clients to request certificates for
        `foo.example.com` and `*.example.com`. To allow further
        levels of subdomains, enable the `allow_subdomains` option.
        There is no default.
      </li>
      <li>
        <span class="param">allow_token_displayname</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates matching
        the value of Display Name from the requesting token.
        Remember, this stacks with the other CN options,
        including `allowed_base_domain`. Defaults to `false`.
      </li>
      <li>
        <span class="param">allow_subdomains</span>
        <span class="param-flags">optional</span>
        If set, clients can request certificates with CNs that
        are subdomains of the CNs allowed by the other role
        options. _This includes wildcard subdomains._ This is
        redundant when using the `allow_any_name` option.
        Defaults to `false`.
      </li>
      <li>
        <span class="param">allow_any_name</span>
        <span class="param-flags">optional</span>
        If set, clients can request any CN. Useful in some
        circumstances, but make sure you understand whether it
        is appropriate for your installation before enabling it.
        Defaults to `false`.
      </li>
      <li>
        <span class="param">allow_ip_sans</span>
        <span class="param-flags">optional</span>
        If set, clients can request IP Subject Alternative
        Names. Unlike CNs, no authorization checking is
        performed except to verify that the given values
        are valid IP addresses. Defaults to `true`.
      <li>
        <span class="param">server_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for server use.
        Defaults to `true`.
      </li>
      <li>
        <span class="param">client_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for client use.
        Defaults to `true`.
      </li>
      <li>
        <span class="param">code_signing_flag</span>
        <span class="param-flags">optional</span>
        If set, certificates are flagged for code signing
        use. Defaults to `false`.
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
            "lease": "6h",
            "lease_max": "12h",
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
