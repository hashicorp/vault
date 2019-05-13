---
layout: "guides"
page_title: "Build Your Own Certificate Authority - Guides"
sidebar_title: "Build Your Own CA"
sidebar_current: "guides-secret-mgmt-pki"
description: |-
  The PKI secrets engine generates dynamic X.509 certificates. With this secrets
  engine, services can get certificates without going through the usual manual
  process of generating a private key and CSR, submitting to a CA, and waiting
  for a verification and signing process to complete. Vault's built-in
  authentication and authorization mechanisms provide the verification
  functionality.
---

# Build Your Own Certificate Authority (CA)

Vault's PKI secrets engine can dynamically generate X.509 certificates on
demand. This allows services to acquire certificates without going through the
usual manual process of generating a private key and Certificates Signing
Request (CSR), submitting to a CA, and then wait for the verification and
signing process to complete.


## Reference Material

- [PKI (Certificates) Secrets Engine](/docs/secrets/pki/index.html)
- [PKI Secrets Engine (API)](/api/secret/pki/index.html)
- [RFC 5280 Internet X.509 Public Key Infrastructure Certificate and Certificate
Revocation List (CRL) Profile](https://tools.ietf.org/html/rfc5280)
- [OpenSSL x509 Man Pages](https://www.openssl.org/docs/man1.1.0/apps/x509.html)


## Estimated Time to Complete

15 minutes


## Personas

The steps described in this guide are typically performed by **security
engineer**.


## Challenge

Organizations should protect their website; however, the Traditional PKI process
workflow takes a long time which motivates organizations to create certificates
which do not expire for a year or more.


## Solution

Use Vault to create X509 certificates for usage in MTLS or other arbitrary PKI
encryption.   While this can be used to create web server certificates.  If
users do not import the CA chains, the browser will complain about self-signed
certificates.

Creating PKI certificates is generally a cumbersome process using traditional
tools like `openssl` or even more advanced frameworks like CFSSL.   These tools
also require a human component to verify certificate distribution meets
organizational security policies.   

Vault PKI secrets engine makes this a lot simpler.  The PKI secrets engine can
be an Intermediate-Only certificate authority which potentially allows for
higher levels of security.

1. Store CA outside the Vault (air gapped)
1. Create CSRs for the intermediates
1. Sign CSR outside Vault and import intermediate
1. Issue leaf certificates from the Intermediate CA


## Prerequisites

To perform the tasks described in this guide, you need to have a Vault
environment.  Refer to the [Getting
Started](/intro/getting-started/install.html) guide to install Vault.

Alternatively, you can use the [Vault
Playground](https://www.katacoda.com/hashicorp/scenarios/vault-playground)
environment.


### <a name="policy"></a>Policy requirements

-> **NOTE:** For the purpose of this guide, you can use **`root`** token to work
with Vault. However, it is recommended that root tokens are only used for just
enough initial setup or in emergencies. As a best practice, use tokens with
appropriate set of policies based on your role in the organization.

To perform all tasks demonstrated in this guide, your policy must include the
following permissions:

```shell
# Enable secrets engine
path "sys/mounts/*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}

# List enabled secrets engine
path "sys/mounts" {
  capabilities = [ "read", "list" ]
}

# Work with pki secrets engine
path "pki*" {
  capabilities = [ "create", "read", "update", "delete", "list", "sudo" ]
}
```

If you are not familiar with policies, complete the
[policies](/guides/identity/policies.html) guide.


## Steps

In this guide, you are going to first generate a self-signed root certificate.
Then you are going to generate an intermediate certificate which is signed by
the root. Finally, you are going to generate a certificate for
`test.example.com` domain.

![Overview](/img/vault-pki-4.png)

In this guide, you perform the following:

1. [Generate Root CA](#step1)
1. [Generate Intermediate CA](#step2)
1. [Create a Role](#step3)
1. [Request Certificates](#step4)
1. [Revoke Certificates](#step5)
1. [Remove Expired Certificates](#step6)



### <a name="step1"></a>Step 1: Generate Root CA

In this step, you are going to generate a self-signed root certificate using
PKI secrets engine.

#### CLI command

1. First, enable the `pki` secrets engine at `pki` path:

    ```plaintext
    $ vault secrets enable pki
    ```

1. Tune the `pki` secrets engine to issue certificates with a maximum time-to-live (TTL) of 87600 hours.

    ```plaintext
    $ vault secrets tune -max-lease-ttl=87600h pki
    ```

1. Generate the ***root*** certificate and save the certificate in `CA_cert.crt`.

    ```plaintext
    $ vault write -field=certificate pki/root/generate/internal common_name="example.com" \
            ttl=87600h > CA_cert.crt
    ```

    This generates a new self-signed CA certificate and private key. Vault will
    _automatically_ revoke the generated root at the end of its lease period
    (TTL); the CA certificate will sign its own Certificate Revocation List
    (CRL).

1. Configure the CA and CRL URLs:

    ```plaintext
    $ vault write pki/config/urls \
            issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" \
            crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"
    ```

#### API call using cURL

1. First, enable the `pki` secrets engine at `pki` path using `/sys/mounts` endpoint:

    ```plaintext
    $ curl --header "X-Vault-Token: <TOKEN>" \
           --request POST \
           --data <PARAMETERS> \
           <VAULT_ADDRESS>/v1/sys/mounts/<PATH>
    ```

    Where `<TOKEN>` is your valid token, and `<PARAMETERS>` holds [configuration
    parameters](/api/system/mounts.html#enable-secrets-engine) of the secret engine.

    **Example:**

    The following example mounts `pki` secret engine.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"type":"pki"}' \
           https://127.0.0.1:8200/v1/sys/mounts/pki
    ```

1. Tune the `pki` secrets engine to issue certificates with a maximum
time-to-live (TTL) of 87600 hours.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"max_lease_ttl":"87600h"}' \
           https://127.0.0.1:8200/v1/sys/mounts/pki/tune
    ```

1. Generate the ***root*** certificate and extract the CA certificate and save
it as `CA_cert.crt`.

    > **NOTE:** The following command uses `jq` tool to parse the output JSON.
    You can install [`jq`](https://stedolan.github.io/jq/download/) or manually
    copy and paste the certificate in a file, `CA_cert.crt`.

    ```plaintext
    $ tee payload.json <<EOF
    {
      "common_name": "example.com",
      "ttl": "87600h"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload.json \
           https://127.0.0.1:8200/v1/pki/root/generate/internal \
           | jq -r ".data.certificate" > CA_cert.crt
    ```

    This generates a new self-signed CA certificate and private key. Vault will
    _automatically_ revoke the generated root at the end of its lease period
    (TTL); the CA certificate will sign its own Certificate Revocation List
    (CRL).

1. Configure the CA and CRL URLs:

    ```plaintext
    $ tee payload-url.json <<EOF
    {
      "issuing_certificates": "http://127.0.0.1:8200/v1/pki/ca",
      "crl_distribution_points": "http://127.0.0.1:8200/v1/pki/crl"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload-url.json \
           https://127.0.0.1:8200/v1/pki/config/urls
    ```


#### Web UI

Open a web browser and launch the Vault UI (e.g. http://127.0.0.1:8200/ui) and
then login.

1. Select **Enable new engine**.
1. Select **PKI** from the **Secrets engine type** drop-down list.
1. Click **More options** to expand and set the **Maximum lease TTL** to **`87600 hours`**.
1. Click **Enable Engine**.
1. Select **Configure**.
1. Click **Configure CA**.
1. Leave **CA Type** as **root**, and **Type** to be **internal**.  Enter
**`example.com`** in the **Common Name** field.
1. Select **Options** and then set **TTL** field to be **87600 hours**.
1. Click **Save**.
1. Click **Copy Certificate** and save it in a file named, **`CA_cert.crt`**.
1. Click the **URLs** tab, and then set:
    - Issuing certificates: `http://127.0.0.1:8200/v1/pki/ca`
    - CRL Distribution Points: `http://127.0.0.1:8200/v1/pki/crl`
    ![Configure URL](/img/vault-pki-1.png)
1. Click **Save**.

<br>

-> **NOTE:** To examine the generated root certificate, you can use
[OpenSSL](https://www.openssl.org/source/).

  ```shell
  # Print the certificate in text form
  $ openssl x509 -in CA_cert.crt -text

  # Print the validity dates
  $ openssl x509 -in CA_cert.crt -noout -dates
  ```


### <a name="step2"></a>Step 2: Generate Intermediate CA

Now, you are going to create an intermediate CA using the root CA you regenerate
in the previous step.

#### CLI Command

1. First, enable the `pki` secrets engine at **`pki_int`** path:

    ```plaintext
    $ vault secrets enable -path=pki_int pki
    ```

1. Tune the `pki_int` secrets engine to issue certificates with a maximum
time-to-live (TTL) of 43800 hours.

    ```plaintext
    $ vault secrets tune -max-lease-ttl=43800h pki_int
    ```

1. Execute the following command to generate an intermediate and save the CSR as
`pki_intermediate.csr`:

    ```plaintext
    $ vault write -format=json pki_int/intermediate/generate/internal \
            common_name="example.com Intermediate Authority" ttl="43800h" \
            | jq -r '.data.csr' > pki_intermediate.csr
    ```

1. Sign the intermediate certificate with the root certificate and save the
generated certificate as `intermediate.cert.pem`:

    ```plaintext
    $ vault write -format=json pki/root/sign-intermediate csr=@pki_intermediate.csr \
            format=pem_bundle \
            | jq -r '.data.certificate' > intermediate.cert.pem
    ```

1. Once the CSR is signed and the root CA returns a certificate, it can be
  imported back into Vault:

    ```plaintext
    $ vault write pki_int/intermediate/set-signed certificate=@intermediate.cert.pem
    ```


#### API call using cURL


1. First, enable the `pki` secrets engine at **`pki_int`** path:

    **Example:**

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"type":"pki"}' \
           https://127.0.0.1:8200/v1/sys/mounts/pki_int
    ```

1. Tune the `pki_int` secrets engine to issue certificates with a maximum
time-to-live (TTL) of 43800 hours.

    ```plaintext
    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data '{"max_lease_ttl":"43800h"}' \
           https://127.0.0.1:8200/v1/sys/mounts/pki_int/tune
    ```

1. Generate an intermediate using the `/pki_int/intermediate/generate/internal`
endpoint.

    ```plaintext
    $ tee payload-int.json <<EOF
    {
      "common_name": "example.com Intermediate Authority",
      "ttl": "43800h"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload-int.json \
           https://127.0.0.1:8200/v1/pki_int/intermediate/generate/internal | jq
    ```

    Copy the generated CSR.

1. Sign the intermediate certificate with the root certificate and save the
certificate as `intermediate.cert.pem`.  

    -> **NOTE:** The API request payload should contain the CSR you obtained.

    ```plaintext
    $ tee payload-int-cert.json <<EOF
    {
      "csr": "...",
      "format": "pem_bundle"
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload-int-cert.json \
           https://127.0.0.1:8200/v1/pki/root/sign-intermediate | jq
    ```

    > **NOTE:** The **`format`** in the payload specifies the format of the
    returned data.  When `pem_bundle`, the certificate field will contain the
    certificate.

    Copy the generated certificate.

1. Once the CSR is signed and the root CA returns a certificate, it can be
imported back into Vault using the `/pki_int/intermediate/set-signed` endpoint.

    -> **NOTE:** The API request payload should contain the certificate you
    obtained.

    ```plaintext
    $ tee payload-signed.json <<EOF
    {
      "certificate": "..."
    }
    EOF

    $ curl --header "X-Vault-Token: ..." \
           --request POST \
           --data @payload-signed.json \
           https://127.0.0.1:8200/v1/pki_int/intermediate/set-signed
    ```


#### Web UI

1. Select **Enable new engine** in the **Secrets** tab.
1. Select **PKI** from the **Secrets engine type** drop-down list.
1. Enter **`pki_int`** in the **Path** field.
1. Click **More options** to expand and set the **Maximum lease TTL** to **`43800
hours`**.
1. Click **Enable Engine**.
1. Select **Configure**.
1. Click **Configure CA**.
1. Select **intermediate** from **CA Type** drop-down list.  
1. Enter **`example.com Intermediate Authority`** in the **Common Name** field,
and then click **Save**.
1. Click **Copy CSR** and save it in a file, `pki_intermediate.csr`.
1. Select **pki** from the **Secrets** tab to return to the root CA.
1. Select **Configure** and then click **Sign intermediate**.
1. Paste in the CSR in the **Certificate Signing Request (CSR)** field.
1. Enter **`example.com`** in the **Common Name**.
1. Select **pem_bundle** from the **Format** drop-down list, and then click
**Save**.
1. Click **Copy Certificate** and save the generated certificate in a file, `intermediate.cert.pem`.
1. Select **pki_int** from the **Secrets** tab to return to the intermediate CA.
1. Select **Configure** and then click **Set signed intermediate**.
1. Paste in the certificate in the **Signed Intermediate Certificate** field and
then click **Save**.



### <a name="step3"></a>Step 3: Create a Role

A role is a logical name that maps to a policy used to generate those
credentials. It allows [configuration
parameters](/api/secret/pki/index.html#create-update-role) to control
certificate common names, alternate names, the key usages that they are valid
for, and more.

Calling out some of the parameters:

- **`allowed_domains`** - specifies the domains of the role (used with
  `allow_bare_domains` and `allow-subdomains` options)
- **`allow_bare_domains`** - specifies if clients can request certificates
matching the value of the actual domains themselves
- **`allow_subdomains`** - specifies if clients can request certificates with
CNs that are subdomains of the CNs allowed by the other role options (NOTE: This
  includes wildcard subdomains.)
- **`allow_glob_domains`** - allows names specified in allowed_domains to
contain glob patterns (e.g. ftp*.example.com)

In this step, you are going to create a role named, **`example-dot-com`**.

#### CLI Command

Create a role named **`example-dot-com`** which allows subdomains.

```plaintext
$ vault write pki_int/roles/example-dot-com \
        allowed_domains="example.com" \
        allow_subdomains=true \
        max_ttl="720h"
```

#### API call using cURL

Create a role named **`example-dot-com`** which allows subdomains.

```plaintext
$ tee payload-role.json <<EOF
{
  "allowed_domains": "example.com",
  "allow_subdomains": true,
  "max_ttl": "720h"
}

$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data @payload-role.json \
       https://127.0.0.1:8200/v1/pki_int/roles/example-dot-com
```

#### Web UI

Create a role named **`example-dot-com`** which allows subdomains.

1. Click **pki_int** and then select **Create role**.
1. Enter **`example-dot-com`** in the **Role name** field.
1. Select **Options** to expand, and then set the **Max TTL** to **`43800
hours`**.   Select **Hide Options**.
1. Select **Domain Handling** to expand, and then select the **Allow
subdomains** check-box. Enter **`example.com`** in the **Allowed domains**
field.
    ![Create Role](/img/vault-pki-2.png)
1. Click **Create role**.


### <a name="step4"></a>Step 4: Request Certificates

Keep certificate lifetimes to be short to align with Vault's philosophy of
short-lived secrets.

#### CLI Command

Execute the following command to request a new certificate for `test.example.com`
domain based on the `example-dot-com` role:

```plaintext
$ vault write pki_int/issue/example-dot-com common_name="test.example.com" ttl="24h"

Key                 Value
---                 -----
certificate         -----BEGIN CERTIFICATE-----
MIIDwzCCAqugAwIBAgIUTQABMCAsXjG6ExFTX8201xKVH4IwDQYJKoZIhvcNAQEL
BQAwGjEYMBYGA1UEAxMPd3d3LmV4YW1wbGUuY29tMB4XDTE4MDcyNDIxMTMxOVoX
             ...

-----END CERTIFICATE-----
issuing_ca          -----BEGIN CERTIFICATE-----
MIIDQTCCAimgAwIBAgIUbMYp39mdj7dKX033ZjK18rx05x8wDQYJKoZIhvcNAQEL
             ...

-----END CERTIFICATE-----
private_key         -----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAte1fqy2Ekj+EFqKV6N5QJlBgMo/U4IIxwLZI6a87yAC/rDhm
W58liadXrwjzRgWeqVOoCRr/B5JnRLbyIKBVp6MMFwZVkynEPzDmy0ynuomSfJkM
             ...

-----END RSA PRIVATE KEY-----
private_key_type    rsa
serial_number       4d:00:01:30:20:2c:5e:31:ba:13:11:53:5f:cd:b4:d7:12:95:1f:82
```

The response contains the PEM-encoded private key, key type and certificate
serial number.



#### API call using cURL

Invoke the **`/pki_int/issue/<role_name>`** endpoint to request a new certificate.

**Example:**

Request a certificate for `test.example.com` domain based on the
`example-dot-com` role:

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"common_name": "test.example.com", "ttl": "24h"}' \
       https://127.0.0.1:8200/v1/pki_int/issue/example-dot-com | jq
{
 "request_id": "6fa8d77d-0758-33ae-b5ea-8b3d15014fd1",
 "lease_id": "",
 "renewable": false,
 "lease_duration": 0,
 "data": {
   "certificate": "-----BEGIN CERTIFICATE-----\nMIIDvzCCAqegAwIBAgIUG7H0Pzpqm+...-----END CERTIFICATE-----",
   "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIDNTCCAh2gAwIBAgIUQhIX9D...-----END CERTIFICATE-----",
   "private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAr6IsROOW5...-----END RSA PRIVATE KEY-----",
   "private_key_type": "rsa",
   "serial_number": "1b:b1:f4:3f:3a:6a:9b:e8:33:af:f7:1b:b1:4d:57:7f:65:65:39:c1"
 },
 "wrap_info": null,
 "warnings": null,
 "auth": null
}
```

The response contains the PEM-encoded private key, key type and certificate
serial number.


#### Web UI

1. Select **Secrets**.
1. Select **pki_int** from the **Secrets Engines** list.
1. Select **example-dot-com** under **Roles**.
1. Enter **`test.example.com`** in the **Common Name** field.
1. Select **Options** to expand, and then set the **TTL** to **`24 hours`**.  
1. Select **Hide Options** and then click **Generate**.
    ![Issue Certificate](/img/vault-pki-3.png)

    > The response contains the PEM-encoded private key, key type and certificate
    serial number.

1. Click **Copy credentials** and save it in a file.

<br>

-> **NOTE:** A certificate can be rotated at any time by simply issuing a new
certificate with the same CN.


### <a name="step5"></a>Step 5: Revoke Certificates

If a certificate must be revoked, you can easily perform the revocation action
which will cause the CRL to be regenerated. When the CRL is regenerated, any
expired certificates are removed from the CRL.


#### CLI Command

In a certain circumstances, you may wish to revoke an issued certificate.

To revoke:

```plaintext
$ vault write pki_int/revoke serial_number=<serial_number>
```

**Example:**

```plaintext
$ vault write pki_int/revoke serial_number="48:97:82:dd:f0:d3:d9:7e:53:25:ba:fd:f6:77:3e:89:e5:65:cc:e7"
Key                        Value
---                        -----
revocation_time            1532539632
revocation_time_rfc3339    2018-07-25T17:27:12.165206399Z
```


#### API call using cURL

Invoke the **`/pki_int/revoke`** endpoint to invoke a certificate using its
serial number.

**Example:**

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"serial_number": "48:97:82:dd:f0:d3:d9:7e:53:25:ba:fd:f6:77:3e:89:e5:65:cc:e7"}' \
       https://127.0.0.1:8200/v1/pki_int/revoke
```

#### Web UI

1. Select **Secrets**.
1. Select **pki_int** from the **Secrets Engines** list.
1. Select the **eCertificates** tab.
1. Select the serial number for the certificate you wish to revoke.
1. Click **Revoke**.  At the confirmation, click **Revoke** again.


### <a name="step6"></a>Step 6: Remove Expired Certificates

Keep the storage backend and CRL by periodically removing certificates that have
expired and are past a certain buffer period beyond their expiration time.

#### CLI Command

To remove revoked certificate and clean the CRL.

```plaintext
$ vault write pki_int/tidy tidy_cert_store=true tidy_revoked_certs=true
```

#### API call using cURL

Invoke the **`/pki_int/tidy`** endpoint to remove revoked certificate and clean
the CRL.

**Example:**

```plaintext
$ curl --header "X-Vault-Token: ..." \
       --request POST \
       --data '{"tidy_cert_store": true, "tidy_revoked_certs": true}' \
       https://127.0.0.1:8200/v1/pki_int/tidy
```

#### Web UI

1. Select **Secrets**.
1. Select **pki_int** from the **Secrets Engines** list.
1. Select **Configure**.
1. Select the **Tidy** tab.
1. Select the check-box for **Tidy the Certificate Store** and **Tidy the Revocation List (CRL)**.
1. Click **Save**.


## Next steps

Check out the [Streamline Certificate Management with HashiCorp
Vault](https://www.hashicorp.com/resources/streamline-certificate-management-with-vault)
webinar recording.
