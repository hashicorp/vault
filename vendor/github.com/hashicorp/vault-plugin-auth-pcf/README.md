# vault-plugin-auth-pcf

This plugin leverages PCF's [App and Container Identity Assurance](https://content.pivotal.io/blog/new-in-pcf-2-1-app-container-identity-assurance-via-automatic-cert-rotation)
for authenticating to Vault. 

## Known Risks

This authentication engine uses PCF's instance identity service to authenticate users to Vault. Because PCF
makes its CA certificate and **private key** available to certain users at any time, it's possible for someone
with access to them to self-issue identity certificates that meet the criteria for a Vault role, allowing
them to gain unintended access to Vault.

For this reason, we recommend that if you choose this auth method, you **carefully guard access to
the private key** for your instance identity CA certificate. In CredHub, it can be obtained through the 
following call: `$ credhub get -n /cf/diego-instance-identity-root-ca`. 

Take extra steps to limit access to that path in CredHub, whether it be through use of CredHub's ACL system, 
or through carefully limiting the users who can access CredHub.

## Getting Started

### Obtaining Your Instance Identity CA Certificate

In most versions of PCF, instance identity is enabled out-of-the-box. Check by pulling your CA certificate,
which you'll need to configure this auth engine. There are undoubtedly multiple ways to do this, but this
is how we did it.

#### From CF Dev

```
$ bosh int --path /diego_instance_identity_ca ~/.cfdev/state/bosh/creds.yml
```

#### From CredHub

[Install and authenticate to the PCF command line tool](https://docs.pivotal.io/tiledev/2-2/pcf-command.html), 
and [install jq](https://stedolan.github.io/jq/). 

Get the credentials you'll use for CredHub:
```
$ pcf settings | jq '.products[0].director_credhub_client_credentials'
```

SSH into your Ops Manager VM:
```
ssh -i ops_mgr.pem ubuntu@$OPS_MGR_URL
```
Please note that the above `OPS_MGR_URL` shouldn't be prepended with `https://`.

Log in to Credhub using the credentials you obtained earlier:
```
$ credhub login --client-name=director_to_credhub --client-secret=CoJPkrsYi3c-Fx2QHEEDyaEEUuOfYMzw
```
6. Retrieve the CA information:
```
$ credhub get -n /cf/diego-instance-identity-root-ca
```
7. You'll receive a response like:
```
id: be2bd996-1d35-443b-b81c-90095024d5e7
name: /cf/diego-instance-identity-root-ca
type: certificate
value:
  ca: |
    -----BEGIN CERTIFICATE-----
    MIIDNDCCAhygAwIBAgITPqTy1qvfHNEVuxsl9l1glY85OTANBgkqhkiG9w0BAQsF
    ADAqMSgwJgYDVQQDEx9EaWVnbyBJbnN0YW5jZSBJZGVudGl0eSBSb290IENBMB4X
    DTE5MDYwNjA5MTIwMVoXDTIyMDYwNTA5MTIwMVowKjEoMCYGA1UEAxMfRGllZ28g
    SW5zdGFuY2UgSWRlbnRpdHkgUm9vdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEP
    ADCCAQoCggEBALa8xGDYT/q3UzEKAsLDajhuHxPpIPFlCXwp6u8U5Qrf427Xof7n
    rXRKzRu3g7E20U/OwzgBi3VZs8T29JGNWeA2k0HtX8oQ+Wc8Qngz9M8t1h9SZlx5
    fGfxPt3x7xozaIGJ8p4HKQH1ZlirL7dzun7Y+7m6Ey8cMVsepqUs64r8+KpCbxKJ
    rV04qtTNlr0LG3yOxSHlip+DDvUVL3jSFz/JDWxwCymiFBAh0QjG1LKp2FisURoX
    GY+HJbf2StpK3i4dYnxQXQlMDpipozK7WFxv3gH4Q6YMZvlmIPidAF8FxfDIsYcq
    TgQ5q0pr9mbu8oKbZ74vyZMqiy+r9vLhbu0CAwEAAaNTMFEwHQYDVR0OBBYEFAHf
    pwqBhZ8/A6ZAvU+p5JPz/omjMB8GA1UdIwQYMBaAFAHfpwqBhZ8/A6ZAvU+p5JPz
    /omjMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBADuDJev+6bOC
    v7t9SS4Nd/zeREuF9IKsHDHrYUZBIO1aBQbOO1iDtL4VA3LBEx6fOgN5fbxroUsz
    X9/6PtxLe+5U8i5MOztK+OxxPrtDfnblXVb6IW4EKhTnWesS7R2WnOWtzqRQXKFU
    voBn3QckLV1o9eqzYIE/aob4z0GaVanA9PSzzbVPsX79RCD1B7NmV0cKEQ7IrCrh
    L7ElDV/GlNrtVdHjY0mwz9iu+0YJvxvcHDTERi106b28KXzJz+P5/hyg2wqRXzdI
    faXAjW0kuq5nxyJUALwxD/8pz77uNt4w6WfJoSDM6XrAIhh15K3tZg9EzBmAZ/5D
    jK0RcmCyaXw=
    -----END CERTIFICATE-----
  certificate: |
    -----BEGIN CERTIFICATE-----
    MIIDNDCCAhygAwIBAgITPqTy1qvfHNEVuxsl9l1glY85OTANBgkqhkiG9w0BAQsF
    ADAqMSgwJgYDVQQDEx9EaWVnbyBJbnN0YW5jZSBJZGVudGl0eSBSb290IENBMB4X
    DTE5MDYwNjA5MTIwMVoXDTIyMDYwNTA5MTIwMVowKjEoMCYGA1UEAxMfRGllZ28g
    SW5zdGFuY2UgSWRlbnRpdHkgUm9vdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEP
    ADCCAQoCggEBALa8xGDYT/q3UzEKAsLDajhuHxPpIPFlCXwp6u8U5Qrf427Xof7n
    rXRKzRu3g7E20U/OwzgBi3VZs8T29JGNWeA2k0HtX8oQ+Wc8Qngz9M8t1h9SZlx5
    fGfxPt3x7xozaIGJ8p4HKQH1ZlirL7dzun7Y+7m6Ey8cMVsepqUs64r8+KpCbxKJ
    rV04qtTNlr0LG3yOxSHlip+DDvUVL3jSFz/JDWxwCymiFBAh0QjG1LKp2FisURoX
    GY+HJbf2StpK3i4dYnxQXQlMDpipozK7WFxv3gH4Q6YMZvlmIPidAF8FxfDIsYcq
    TgQ5q0pr9mbu8oKbZ74vyZMqiy+r9vLhbu0CAwEAAaNTMFEwHQYDVR0OBBYEFAHf
    pwqBhZ8/A6ZAvU+p5JPz/omjMB8GA1UdIwQYMBaAFAHfpwqBhZ8/A6ZAvU+p5JPz
    /omjMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBADuDJev+6bOC
    v7t9SS4Nd/zeREuF9IKsHDHrYUZBIO1aBQbOO1iDtL4VA3LBEx6fOgN5fbxroUsz
    X9/6PtxLe+5U8i5MOztK+OxxPrtDfnblXVb6IW4EKhTnWesS7R2WnOWtzqRQXKFU
    voBn3QckLV1o9eqzYIE/aob4z0GaVanA9PSzzbVPsX79RCD1B7NmV0cKEQ7IrCrh
    L7ElDV/GlNrtVdHjY0mwz9iu+0YJvxvcHDTERi106b28KXzJz+P5/hyg2wqRXzdI
    faXAjW0kuq5nxyJUALwxD/8pz77uNt4w6WfJoSDM6XrAIhh15K3tZg9EzBmAZ/5D
    jK0RcmCyaXw=
    -----END CERTIFICATE-----
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
    MIIEpAIBAAKCAQEAtrzEYNhP+rdTMQoCwsNqOG4fE+kg8WUJfCnq7xTlCt/jbteh
    /uetdErNG7eDsTbRT87DOAGLdVmzxPb0kY1Z4DaTQe1fyhD5ZzxCeDP0zy3WH1Jm
    XHl8Z/E+3fHvGjNogYnyngcpAfVmWKsvt3O6ftj7uboTLxwxWx6mpSzrivz4qkJv
    EomtXTiq1M2WvQsbfI7FIeWKn4MO9RUveNIXP8kNbHALKaIUECHRCMbUsqnYWKxR
    GhcZj4clt/ZK2kreLh1ifFBdCUwOmKmjMrtYXG/eAfhDpgxm+WYg+J0AXwXF8Mix
    hypOBDmrSmv2Zu7ygptnvi/JkyqLL6v28uFu7QIDAQABAoIBACGIPhjvWK3PGirz
    hVIr/b/hJT7IFs11Fup73qqEkQsPznI2i3l1FfUzDLQ7VqUcRAh7DoOmdOrRzRUl
    o/dZktZ77UW5w0wXFU0GV8Qq9I9X/+S7gCEUAeoo8LpVfOS37kNnBuhMtA+x8lfv
    AdCOIfjI5FhOdtq8N6pa04WX2pkkRzQkIpneRcLPqq1WwKZK1o8zxCnbP+4SI8yo
    dTB2ldDY+1vusMYoFch2IsPDMCxVUpAYOoO0jvyOM4cqm0m7P+Tb6Qy19GG50ZfC
    PlEK76YTOurdirGWdnGC+JEf4smrGgIvJGKEb55/qbqPow/o6Bf4XVXnOHaboW/W
    Mu4cGFkCgYEA2JaQxFF51yuz4W7/VFZCCix3woPH12fZ3fz7aVMUfn4WQwMbNLVa
    7G4gdRclReOs5FLM7jPqiTxDckXnoWRc5Ff9Xn4c2ioXrqzXr5V4qJ4GAWxAx0uM
    w1u5ZpVL2HO12pat74MnYw7EJ65oznQNFSC1FAGn5BJ9f5HFk4X3ngMCgYEA1/1Q
    XmAk1XJUQ0RP0EehwNipZQPhodGmrqHBGED+N/8+eRMo4shi//EcPXZ222j9kqAE
    inPA9qaDxhBjgMt+JBFkj/bmTO/Yz8XusBBa5YlN9Ev30zlO+dRlM41/piluPTzf
    vNQuzyNIzl2Gzd71R1TcuFWIDxn8BR0/cBA/5E8CgYAT7m8uEc1jlrr8AOnwSevT
    4dm3hccLNJxhCFnejG2zYkkMK6oCRLo0TcIg5Ftivhv3+wKu3Qo1TN1sE7DIMmM2
    BD7lxjdDgGIjifZjSx8KbVhiIyMm8/XlOHisTwrmxWcz0W/6PZiPThmRCUTN0vIt
    QpBHYgugOm9gIPsMo2RxHwKBgQDOUDjZvUrR3GCi1HjMwe+/bvX3+MopMULfYsE4
    srRittxs+KFAZxsx0ZUhHKySDurQiSttOP6kXBBZPERfvYFjYH3HipcX/K8EYNQL
    t8OrqAkfhwVV7VMEDx8QLGQ3SzHzKteo3qFL2S9teCcRNZzjoysmpQTPMAnstLBp
    EgyFvwKBgQDObNn/Kmfwi6TuGhIjLtBuUEhp5n4EUtysTUZs/15h02MWOfI8CCvm
    xWb6/vZrVggxGlZgZtKy9+COPVpEMFaVdwq9uq4lW77sSBwGIwfzHd1CIjce6mSg
    P5+wO3aTgvr4n8D5NyWcnYPJKRQzqWHHnfk+9TQA1l0g3/yQXfCx2A==
    -----END RSA PRIVATE KEY-----
version_created_at: "2019-06-06T09:12:01Z"
```

From that response, copy the first certificate (under `ca: |`) and place
it into its own separate file using a plain text editor like [Sublime](https://www.sublimetext.com/).
The following instructions assume you name the file `ca.crt`.
Remove any tabs before each line, and any trailing space or lines. When
complete, your CA certificate should look like this:

```
$ cat ca.crt
-----BEGIN CERTIFICATE-----
MIIDNDCCAhygAwIBAgITPqTy1qvfHNEVuxsl9l1glY85OTANBgkqhkiG9w0BAQsF
ADAqMSgwJgYDVQQDEx9EaWVnbyBJbnN0YW5jZSBJZGVudGl0eSBSb290IENBMB4X
DTE5MDYwNjA5MTIwMVoXDTIyMDYwNTA5MTIwMVowKjEoMCYGA1UEAxMfRGllZ28g
SW5zdGFuY2UgSWRlbnRpdHkgUm9vdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBALa8xGDYT/q3UzEKAsLDajhuHxPpIPFlCXwp6u8U5Qrf427Xof7n
rXRKzRu3g7E20U/OwzgBi3VZs8T29JGNWeA2k0HtX8oQ+Wc8Qngz9M8t1h9SZlx5
fGfxPt3x7xozaIGJ8p4HKQH1ZlirL7dzun7Y+7m6Ey8cMVsepqUs64r8+KpCbxKJ
rV04qtTNlr0LG3yOxSHlip+DDvUVL3jSFz/JDWxwCymiFBAh0QjG1LKp2FisURoX
GY+HJbf2StpK3i4dYnxQXQlMDpipozK7WFxv3gH4Q6YMZvlmIPidAF8FxfDIsYcq
TgQ5q0pr9mbu8oKbZ74vyZMqiy+r9vLhbu0CAwEAAaNTMFEwHQYDVR0OBBYEFAHf
pwqBhZ8/A6ZAvU+p5JPz/omjMB8GA1UdIwQYMBaAFAHfpwqBhZ8/A6ZAvU+p5JPz
/omjMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBADuDJev+6bOC
v7t9SS4Nd/zeREuF9IKsHDHrYUZBIO1aBQbOO1iDtL4VA3LBEx6fOgN5fbxroUsz
X9/6PtxLe+5U8i5MOztK+OxxPrtDfnblXVb6IW4EKhTnWesS7R2WnOWtzqRQXKFU
voBn3QckLV1o9eqzYIE/aob4z0GaVanA9PSzzbVPsX79RCD1B7NmV0cKEQ7IrCrh
L7ElDV/GlNrtVdHjY0mwz9iu+0YJvxvcHDTERi106b28KXzJz+P5/hyg2wqRXzdI
faXAjW0kuq5nxyJUALwxD/8pz77uNt4w6WfJoSDM6XrAIhh15K3tZg9EzBmAZ/5D
jK0RcmCyaXw=
-----END CERTIFICATE-----
```

Verify that this certificate can be properly parsed like so:

```
$ openssl x509 -in ca.crt -text -noout
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            3e:a4:f2:d6:ab:df:1c:d1:15:bb:1b:25:f6:5d:60:95:8f:39:39
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: CN=Diego Instance Identity Root CA
        Validity
            Not Before: Jun  6 09:12:01 2019 GMT
            Not After : Jun  5 09:12:01 2022 GMT
        Subject: CN=Diego Instance Identity Root CA
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:b6:bc:c4:60:d8:4f:fa:b7:53:31:0a:02:c2:c3:
                    6a:38:6e:1f:13:e9:20:f1:65:09:7c:29:ea:ef:14:
                    e5:0a:df:e3:6e:d7:a1:fe:e7:ad:74:4a:cd:1b:b7:
                    83:b1:36:d1:4f:ce:c3:38:01:8b:75:59:b3:c4:f6:
                    f4:91:8d:59:e0:36:93:41:ed:5f:ca:10:f9:67:3c:
                    42:78:33:f4:cf:2d:d6:1f:52:66:5c:79:7c:67:f1:
                    3e:dd:f1:ef:1a:33:68:81:89:f2:9e:07:29:01:f5:
                    66:58:ab:2f:b7:73:ba:7e:d8:fb:b9:ba:13:2f:1c:
                    31:5b:1e:a6:a5:2c:eb:8a:fc:f8:aa:42:6f:12:89:
                    ad:5d:38:aa:d4:cd:96:bd:0b:1b:7c:8e:c5:21:e5:
                    8a:9f:83:0e:f5:15:2f:78:d2:17:3f:c9:0d:6c:70:
                    0b:29:a2:14:10:21:d1:08:c6:d4:b2:a9:d8:58:ac:
                    51:1a:17:19:8f:87:25:b7:f6:4a:da:4a:de:2e:1d:
                    62:7c:50:5d:09:4c:0e:98:a9:a3:32:bb:58:5c:6f:
                    de:01:f8:43:a6:0c:66:f9:66:20:f8:9d:00:5f:05:
                    c5:f0:c8:b1:87:2a:4e:04:39:ab:4a:6b:f6:66:ee:
                    f2:82:9b:67:be:2f:c9:93:2a:8b:2f:ab:f6:f2:e1:
                    6e:ed
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Subject Key Identifier: 
                01:DF:A7:0A:81:85:9F:3F:03:A6:40:BD:4F:A9:E4:93:F3:FE:89:A3
            X509v3 Authority Key Identifier: 
                keyid:01:DF:A7:0A:81:85:9F:3F:03:A6:40:BD:4F:A9:E4:93:F3:FE:89:A3

            X509v3 Basic Constraints: critical
                CA:TRUE
    Signature Algorithm: sha256WithRSAEncryption
         3b:83:25:eb:fe:e9:b3:82:bf:bb:7d:49:2e:0d:77:fc:de:44:
         4b:85:f4:82:ac:1c:31:eb:61:46:41:20:ed:5a:05:06:ce:3b:
         58:83:b4:be:15:03:72:c1:13:1e:9f:3a:03:79:7d:bc:6b:a1:
         4b:33:5f:df:fa:3e:dc:4b:7b:ee:54:f2:2e:4c:3b:3b:4a:f8:
         ec:71:3e:bb:43:7e:76:e5:5d:56:fa:21:6e:04:2a:14:e7:59:
         eb:12:ed:1d:96:9c:e5:ad:ce:a4:50:5c:a1:54:be:80:67:dd:
         07:24:2d:5d:68:f5:ea:b3:60:81:3f:6a:86:f8:cf:41:9a:55:
         a9:c0:f4:f4:b3:cd:b5:4f:b1:7e:fd:44:20:f5:07:b3:66:57:
         47:0a:11:0e:c8:ac:2a:e1:2f:b1:25:0d:5f:c6:94:da:ed:55:
         d1:e3:63:49:b0:cf:d8:ae:fb:46:09:bf:1b:dc:1c:34:c4:46:
         2d:74:e9:bd:bc:29:7c:c9:cf:e3:f9:fe:1c:a0:db:0a:91:5f:
         37:48:7d:a5:c0:8d:6d:24:ba:ae:67:c7:22:54:00:bc:31:0f:
         ff:29:cf:be:ee:36:de:30:e9:67:c9:a1:20:cc:e9:7a:c0:22:
         18:75:e4:ad:ed:66:0f:44:cc:19:80:67:fe:43:8c:ad:11:72:
         60:b2:69:7c
```

Congratulations! You have obtained the CA certificate you'll use for configuring
this auth engine.

### Obtaining Your API Credentials

From the directory where you added `metadata` in the previous step to authenticate to the pcf command-line
tool, run the following commands:

```
$ pcf target
$ cf api
```

The api endpoint given will be used for configuring this Vault auth method. 

This plugin was tested with Org Manager level permissions, but lower level permissions may be usable.
```
$ cf create-user vault pa55word
$ cf orgs
$ cf org-users my-example-org
$ cf set-org-role Alice my-example-org OrgManager
```

Since the PCF API tends to use a self-signed certificate, you'll also need to configure
Vault to trust that certificate. You can obtain its API certificate via:

```
openssl s_client -showcerts -servername domain.com -connect domain.com:443
```

You'll see a certificate outputted as part of the response, which should be broken
out into a separate well-formatted file like the `ca.crt` above, and used for the
`pcf_api_trusted_certificates` field.

## Downloading the Plugin

- `$ git clone git@github.com:hashicorp/vault-plugin-auth-pcf.git`
- `$ cd vault-plugin-auth-pcf`
- `$ PCF_HOME=$(pwd)`

## Sample Usage

Please note that this example uses `generate-signature`, a tool installed through `$ make tools`.

First, enable the PCF auth engine.
```
$ vault auth enable pcf
```

Next, configure the plugin. In the `config` call below, `certificates` is intended to be the instance
identity CA certificate you pulled above.

In the CF Dev environment the default API address is `https://api.dev.cfdev.sh`. The default username and password
are `admin`, `admin`. In a production environment, these attributes will vary.
```
$ vault write auth/pcf/config \
      certificates=@ca.crt \
      pcf_api_addr=https://api.dev.cfdev.sh \
      pcf_username=admin \
      pcf_password=admin \
      pcf_api_trusted_certificates=@pcfapi.crt
```

Then, add a role that will be used to grant specific Vault policies to those logging in with it. When a constraint like
`bound_application_ids` is added, then the application ID on the cert used for logging in _must_ be one of the role's
application IDs. However, if `bound_application_ids` is omitted, then _any_ application ID will match. We recommend
configuring as many bound parameters as possible.

Also, by default, the IP address on the certificate presented at login must match that of the caller. However, if
your callers tend to be proxied, this may not work for you. If that's the case, set `disable_ip_matching` to true.
```
$ vault write auth/pcf/roles/test-role \
    bound_application_ids=2d3e834a-3a25-4591-974c-fa5626d5d0a1 \
    bound_space_ids=3d2eba6b-ef19-44d5-91dd-1975b0db5cc9 \
    bound_organization_ids=34a878d0-c2f9-4521-ba73-a9f664e82c7bf \
    policies=foo-policies
```

Logging in is intended to be performed using your `CF_INSTANCE_CERT` and `CF_INSTANCE_KEY`. This is an example of how
it can be done.
```
$ export CF_INSTANCE_CERT=$PCF_HOME/testdata/fake-certificates/instance.crt
$ export CF_INSTANCE_KEY=$PCF_HOME/testdata/fake-certificates/instance.key
$ vault login -method=pcf role=test-role
```

### Updating the CA Certificate

In PCF, most CA certificates expire after 4 years. However, it's possible to configure your own CA certificate for the
instance identity service, and its expiration date could vary. Either way, sometimes CA certificates expire and it may
be necessary to have multiple configured so the beginning date of once commences when another expires.

To configure multiple certificates, simply update the config to include the current one and future one.
```
$ CURRENT=$(cat /path/to/current-ca.crt)
$ FUTURE=$(cat /path/to/future-ca.crt)
$ vault write auth/vault-plugin-auth-pcf/config certificates="$CURRENT,$FUTURE"
```

All other configured values will remain untouched; however, the previous value for `certificates` will be overwritten
with the new one you've provided.

Providing a future CA certificate before the current one expires can protect you from having a downtime while the service
is switching over from the old to the new. If a client certificate was issued by _any_ CA certificate you've configured,
login will succeed.

## Troubleshooting

### Obtaining a Certificate Error from the PCF API

When configuring this plugin, you may encounter an error like:
```
Error writing data to auth/pcf/config: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/auth/pcf/config
Code: 500. Errors:

* 1 error occurred:
	* unable to establish an initial connection to the PCF API: Could not get api /v2/info: Get https://api.sys.lagunaniguel.cf-app.com/v2/info: x509: certificate signed by unknown authority
```

To resolve this error, review instructions above regarding setting the `pcf_api_trusted_certificates` field.

### verify-certs

This tool, installed by `make tools`, is for verifying that your CA certificate, client certificate, and client 
key are all properly related to each other and will pass verification if used by this auth engine. If you're 
debugging authentication problems that may be related to your certificates, it's a fantastic tool to use.

```
verify-certs -ca-cert=local/path/to/ca.crt -instance-cert=local/path/to/instance.crt -instance-key=local/path/to/instance.key
```
The `ca-cert` should be the cert that was used to issue the given client certificate.

The `instance-cert` given should be the value for the `CF_INSTANCE_CERT` variable in the PCF environment you're
using, and the `instance-key` should be the value for the `CF_INSTANCE_KEY`.

The tool does take the _local path to_ these certificates, so you'll need to gather them and place them on your
local machine to verify they all will work together.

### generate-signature

This tool, installed by `make tools`, is for generating a valid signature to be used for signing into Vault via PCF. 

It can be used as a standalone tool for generating a signature like so:
```
export CF_INSTANCE_CERT=path/to/instance.crt
export CF_INSTANCE_KEY=path/to/instance.key
export SIGNING_TIME=$(date -u)
export ROLE='test-role'
generate-signature
```

It can also be used for signing into Vault like so:
```
export CF_INSTANCE_CERT=path/to/instance.crt
export CF_INSTANCE_KEY=path/to/instance.key
export SIGNING_TIME=$(date -u)
export ROLE='test-role'

vault write auth/vault-plugin-auth-pcf/login \
    role=$ROLE \
    certificate=$CF_INSTANCE_CERT \
    signing-time=SIGNING_TIME \
    signature=$(generate-signature)
```
If the tool is being run in a PCF environment already containing the `CF_INSTANCE_CERT` and `CF_INSTANCE_KEY`, those
variables obviously won't need to be manually set before the tool is used and can just be pulled as they are.

## Developing

### mock-pcf-server

This tool, installed by `make tools`, is for use in development. It lets you run a mocked PCF server for use in local 
testing, with output that can be used as the `pcf_api_addr`, `pcf_username`, and `pcf_password` in your config.

Example use:
```
$ mock-pcf-server
running at http://127.0.0.1:33671
username is username
password is password
```

Simply hit CTRL+C to stop the test server.

### Implementing the Signature Algorithm in Other Languages

The signing algorithm used by this plugin is viewable in `signatures/version1.go`. There is also a test
called `TestSignature` in the same package that outputs a viewable signing string, hash of it, and
resulting signature. The signature will be different every time the test is run because some
of the input to the final signature includes cryptographically random material. This means that no matter
what you do, your final signature won't match any signatures shown; the important thing, however, is that 
it can be verified as having been signed by the private key that's associated with the given client
certificate.

To develop your own version of the signing algorithm in a different language, we recommend you duplicate
the inputs to `TestSignature`, duplicate its signing string and hash, and duplicate the signing algorithm used.

### Quick Start

```
# After cloning the repo, generate fake certs, a test binary, and install the tools.
make test
make dev
make tools

# In one shell window, run Vault with the plugin available in the catalog.
vault server -dev -dev-root-token-id=root -dev-plugin-dir=$PCF_HOME/bin -log-level=debug

# In another shell window, run a mock of the PCF API so the plugin's client calls won't fail.
mock-pcf-server

# In another shell window, execute the following commands to exercise each endpoint.
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=root
export MOCK_PCF_SERVER_ADDR='something' # ex. http://127.0.0.1:32937

vault auth enable vault-plugin-auth-pcf

vault write auth/vault-plugin-auth-pcf/config \
    certificates=@$PCF_HOME/testdata/fake-certificates/ca.crt \
    pcf_api_addr=$MOCK_PCF_SERVER_ADDR \
    pcf_username=username \
    pcf_password=password
    
vault write auth/vault-plugin-auth-pcf/roles/test-role \
    bound_application_ids=2d3e834a-3a25-4591-974c-fa5626d5d0a1 \
    bound_space_ids=3d2eba6b-ef19-44d5-91dd-1975b0db5cc9 \
    bound_organization_ids=34a878d0-c2f9-4521-ba73-a9f664e82c7bf \
    bound_instance_ids=1bf2e7f6-2d1d-41ec-501c-c70 \
    policies=foo,policies \
    disable_ip_matching=true \
    ttl=86400s \
    max_ttl=86400s \
    period=86400s
    
export CF_INSTANCE_CERT=$PCF_HOME/testdata/fake-certificates/instance.crt
export CF_INSTANCE_KEY=$PCF_HOME/testdata/fake-certificates/instance.key
export SIGNING_TIME=$(date -u)
export ROLE='test-role'
vault write auth/vault-plugin-auth-pcf/login \
    role=$ROLE \
    certificate=@$CF_INSTANCE_CERT \
    signing_time="$SIGNING_TIME" \
    signature=$(generate-signature)
    
vault token renew <token>

CURRENT=$(cat $PCF_HOME/testdata/fake-certificates/ca.crt)
FUTURE=$(cat $PCF_HOME/testdata/fake-certificates/ca.crt)
vault write auth/vault-plugin-auth-pcf/config certificates="$CURRENT,$FUTURE"
```