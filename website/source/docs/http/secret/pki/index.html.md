---
layout: "http"
page_title: "PKI Secret Backend - HTTP API"
sidebar_current: "docs-http-secret-pki"
description: |-
  This is the API documentation for the Vault PKI secret backend.
---

# PKI Secret Backend HTTP API

This is the API documentation for the Vault PKI secret backend. For general
information about the usage and operation of the PKI backend, please see the
[Vault PKI backend documentation](/docs/secrets/pki/index.html).

This documentation assumes the PKI backend is mounted at the `/pki` path in
Vault. Since it is possible to mount secret backends at any location, please
update your API calls accordingly.


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
    <binary DER-encoded certificate>
    ```

  </dd>
</dl>

### /pki/ca_chain
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves the CA certificate chain, including the CA *in PEM format*. This
    is a bare endpoint that does not return a standard Vault data structure.
    <br /><br />This is an unauthenticated endpoint.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/pki/ca_chain`</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```
    <PEM-encoded certificate chain>
    ```

  </dd>
</dl>

### /pki/cert/
#### GET

<dl class="api">
  <dt>Description</dt>
  <dd>
    Retrieves one of a selection of certificates. Valid values: `ca` for the CA
    certificate, `crl` for the current CRL, `ca_chain` for the CA trust chain
    or a serial number in either hyphen-separated or colon-separated octal
    format. This endpoint returns the certificate in PEM formatting in the
    `certificate` key of the JSON object. <br /><br />This is an unauthenticated endpoint.
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

### /pki/certs/
#### LIST

<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns a list of the current certificates by serial number only.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/pki/certs` (LIST) or `/pki/certs?list=true` (GET)</dd>

  <dt>Parameters</dt>
  <dd>
     None
  </dd>

  <dt>Returns</dt>
  <dd>

    ```javascript
    {
      "lease_id":"",
      "renewable":false,
      "lease_duration":0,
      "data":{
        "keys":[
          "17:67:16:b0:b9:45:58:c0:3a:29:e3:cb:d6:98:33:7a:a6:3b:66:c1",
          "26:0f:76:93:73:cb:3f:a0:7a:ff:97:85:42:48:3a:aa:e5:96:03:21"
        ]
      },
      "wrap_info":null,
      "warnings":null,
      "auth":null
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
        `pem_bundle`, the `certificate` field will contain the private key and
        certificate, concatenated; if the issuing CA is not a Vault-derived
        self-signed root, this will be included as well.
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
        "ca_chain": ["-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n"],
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
        <span class="param">key_usage</span>
        <span class="param-flags">optional</span>
        This sets the allowed key usage constraint on issued certificates. This
        is a comma-separated string; valid values can be found at
        https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop the
        `KeyUsage` part of the value. Values are not case-sensitive. To specify
        no key usage constraints, set this to an empty string. Defaults to
        `DigitalSignature,KeyAgreement,KeyEncipherment`.
      </li>
      <li>
        <span class="param">use_csr_common_name</span>
        <span class="param-flags">optional</span>
        If set, when used with the CSR signing endpoint, the common name in the
        CSR will be used instead of taken from the JSON data. This does `not`
        include any requested SANs in the CSR; use `use_csr_sans` for that.
        Defaults to `true`.
      </li>
      <li>
        <span class="param">use_csr_sans</span>
        <span class="param-flags">optional</span>
        If set, when used with the CSR signing endpoint, the subject alternate
        names in the CSR will be used instead of taken from the JSON data. This
        does `not` include the common name in the CSR; use
        `use_csr_common_name` for that. Defaults to `true`.
      </li>
      <li>
        <span class="param">allow_token_displayname</span>
        <span class="param-flags">optional</span>
        If set, the display name of the token used when requesting a
        certificate will be considered to be a valid host name by the role.
        Normal verification behavior applies with respect to subdomains and
        wildcards.
      </li>
      <li>
        <span class="param">ou</span>
        <span class="param-flags">optional</span>
        This sets the OU (OrganizationalUnit) values in the subject field of
        issued certificates. This is a comma-separated string.
      </li>
      <li>
        <span class="param">organization</span>
        <span class="param-flags">optional</span>
        This sets the O (Organization) values in the subject field of issued
        certificates. This is a comma-separated string.
      </li>
      <li>
        <span class="param">generate_lease</span>
        <span class="param-flags">optional</span>
        If set, certificates issued/signed against this role will have Vault
        leases attached to them. Defaults to "false". Certificates can be added
        to the CRL by `vault revoke <lease_id>` when certificates are
        associated with leases.  It can also be done using the `pki/revoke`
        endpoint. However, when lease generation is disabled, invoking
        `pki/revoke` would be the only way to add the certificates to the CRL.
        When large number of certificates are generated with long lifetimes, it
        is recommended that lease generation be disabled, as large amount of
        leases adversely affect the startup time of Vault.
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
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/pki/roles` (LIST) or `/pki/roles?list=true` (GET)</dd>

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
    "lease_duration": 2764800,
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
        `pem_bundle`, the `certificate` field will contain the private key (if
        exported) and certificate, concatenated; if the issuing CA is not a
        Vault-derived self-signed root, this will be included as well.
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
      "lease_id": "",
      "lease_duration": 0,
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
        `pem_bundle`, the `certificate` field will contain the certificate and,
        if the issuing CA is not a Vault-derived self-signed root, it will be
        concatenated with the certificate.
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
      "lease_id": "",
      "renewable": false,
      "lease_duration": 0,
      "data": {
        "certificate": "-----BEGIN CERTIFICATE-----\nMIIDzDCCAragAwIBAgIUOd0ukLcjH43TfTHFG9qE0FtlMVgwCwYJKoZIhvcNAQEL\n...\numkqeYeO30g1uYvDuWLXVA==\n-----END CERTIFICATE-----\n",
        "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n",
        "ca_chain": ["-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n"],
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
        `pem_bundle`, the `certificate` field will contain the certificate and,
        if the issuing CA is not a Vault-derived self-signed root, it will be
        concatenated with the certificate.
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
        "ca_chain": ["-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n"],
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
        `pem_bundle`, the `certificate` field will contain the certificate and,
        if the issuing CA is not a Vault-derived self-signed root, it will be
        concatenated with the certificate.
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
        "ca_chain": ["-----BEGIN CERTIFICATE-----\nMIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\n...\nG/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==\n-----END CERTIFICATE-----\n"]
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
=======
The PKI secret backend has a full HTTP API. Please see the
[PKI secret backend API](/docs/http/secret/pki/index.html) for more
details.
>>>>>>> e54ffcd1... Break out API documentation for secret backends
