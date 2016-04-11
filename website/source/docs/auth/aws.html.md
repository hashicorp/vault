---
layout: "docs"
page_title: "Auth Backend: AWS EC2"
sidebar_current: "docs-auth-aws"
description: |-
  The AWS EC2 backend is a mechanism for AWS EC2 instances to authenticate with Vault.
---

# Auth Backend: AWS EC2

The AWS EC2 auth backend is a mechanism for AWS EC2 instances to authenticate
with Vault in an automated fashion. This solves the problem of secure introduction
of EC2 instances to Vault server and avoids the need to create and issue Vault
tokens to each instance manually. It works by using the dynamic metadata information
that uniquely represents each EC2 instance.

## Authentication workflow

EC2 instances will have access to its instance metadata. Details about EC2 instance
metadata can be found [here](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html).

Of all the "dynamic metadata" available to the EC2 instances, the instance identity
document and its PKCS#7 signature are of particular use in this backend. For details
on retrieving the PKCS#7 signature, see [here](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html).

Instance identity document contains enough information to uniquely identify an
EC2 instance. EC2 instance will have access to PKCS#7 signature of its identity
document. This signature contains the instance identity document, along with
the signer information that can establish the authenticity of the contents in
the signature. The signature can be verified using the public certificate provided
by AWS (public certificate varies by region).

During the login, to establish authenticity of the information provided by the
client (EC2 instance), the PKCS#7 signature is validated by the backend. Before
succeeding the login attempt and returning a Vault token, AWS API DescribeInstanceStatus
is invoked to check if the instance is healthy.

## Authorization workflow

The AMIs that are used by instances should be associated with Vault policies at
priori, which provides access control primitives on the resources. A successful
login returns a token. The policies of this token are the same policies that are
associated with the registered AMI. If `role_tag` option (refer API section) is
enabled on the AMI, then the policies of the token will be the subset of the
policies that are associated with the AMI.

## Client Nonce

If an unintended party gets access to the PKCS#7 signature of a particular
instance, it can impersonate that instance and fetch a Vault token. The design
of this backend addresses this problem by sharing the responsibility with the
clients of this backend. The backend will **NOT** be able to distinguish the
genuineness of the request, during the first login. But once an instance performs
a successful login, the backend can then thwart the replay-login attempts from
unintended parties, using a unique nonce that is supplied by the client, during
its first successful login. The login from an unintended party is detected when
the instance tries to login for the first time and it fails. A security alert
should be triggered in such cases.

The client should ensure that it generates unique nonces and makes sure that
it uses the same nonce for each login attempt. During the first login, the
backend caches the client nonce in a `whitelist`. For the subsequent login
requests to succeed, the presented client nonce should match the cached nonce.
Hence, if the nonce is lost/changed then a token cannot be refreshed (rotated).

## Advanced options and caveats

### Dynamic management of policies via role tags
If the instance is required to have customized set of policies based on the
role it plays, it can be achieved by setting `role_tag` option (refer API
section) on the registered AMI. When this option is set, during the login,
along with verification of PKCS#7 signature and instance health, the backend
will query for a specific tag that is attached to the instance. This tag will
hold information that represents a subset of capabilities that are set on the
AMI. Hence, a successful login when `role_tag` is enabled on AMI, returns a
token with the capabilities that are a subset of the capabilities configured
on the AMI. A `role_tag` can be created using `auth/aws/image/<ami_id>/roletag`
endpoint and is immutable. The information present in the tag is SHA256 hashed
and HMAC protected. The key to HMAC is only maintained in the backend.

### Handling lost client nonce
If an EC2 instance loses its client nonce when it migrates to a different host,
say after a stop and start action on the instance, the subsequent login attempts
will not succeed. If the client nonce is lost, 2 administrative actions can be
taken.One option is to delete the entry corresponding to the instance ID from
the identity `whitelist` in the backend. This can be done via `auth/aws/whitelist/identity/<instance_id>`
endpoint. This allows a new client nonce to be accepted by the backend during
the next login request. The other option is to relax the condition of matching
the client nonce through `allow_instance_migration`(refer API section). When
this option is enabled, only `pendingTime` in the instance identity document
will be checked to be newer than the `pendingTime` in the instance identity
document, that was used to login previously. This option should be used with
caution, since any entity that has access to instance PKCS#7 signature can imitate
the instance to get a new Vault token, and only the requirement of newer `pendingTime`,
will be the line of defense against such attacks.

### Disabling reauthentication
If a client chooses to fetch a long-lived Vault token and intends to not refresh
(rotate) the token, then it can disable all future logins. If the option
`disallow_reauthentication` is set, only one login will be allowed per instance.
If the instance successfully gets the token for the first time, it can use it
without worrying about its token getting hijacked by another entity. The client
will still need to raise a security alert if the first login fails, since the
backend will not be able to distinguish a genuine login attempt from an imitation,
for the first time.

When `disallow_reauthentication` option is enabled, the backend only allows a
single successful login from the client. In this case, the client nonce loses
its significance and hence the client can choose not to supply the nonce during
the login.

### Blacklisting role tags
It maybe difficult to track the created role tags and to get to know which instances
are indeed using specific role tags. In these cases, when a role tag needs to be
blocked from any further login attempts, it can be placed in a `blacklist` via the
endpoint `auth/aws/blacklist/roletag/<role_tag>`. Note that this will not invalidate
the tokens that were already issued. This only blocks any further login requests.

### Expiration times and tidying of `blacklist` and `whitelist` entries
The entries in both identity `whitelist` and role tag `blacklist` are not deleted
automatically. The entries in both of these lists will have an expiration time
which is dynamically determined by three factors: `max_ttl` set on the AMI,
`max_ttl` set on the role tag and `max_ttl` value of the backend mount. The
least of these three will be set as the expiration times of these entries.
Separate endpoints `aws/auth/whitelist/identity/tidy` and `aws/auth/blacklist/roletag/tidy`
are provided to cleanup the entries present in these lists.

### Varying public certificates
AWS public key which is used to verify the PKCS#7 signature varies by region.
To check if the default public certificate is applicable for the instances
or to get a different public certificate, refer [this](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html).
If the instances that are using this backend require more than one certificate,
then this backend needs to be mounted at as many paths as there are certificates.
The clients should then use appropriate mount of the backend which can verify its
PKCS#7 signature.

## Authentication

### Via the CLI

#### Enable AWS EC2 authentication in Vault.

```
$ vault auth-enable aws
```

#### Configure the credentials required to make AWS API calls.

```
$ vault write auth/aws/config/client secret_key=vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj access_key=VKIAJBRHKH6EVTTNXDHA region=us-east-1
```

#### Configure the policies on the AMI.

```
$ vault write auth/aws/image/ami-fce3c696 policies=prod,dev max_ttl=500h
```

#### Perform the login operation

```
$ vault write auth/aws/login pkcs7=MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGmewogICJkZXZwYXlQcm9kdWN0Q29kZXMiIDogbnVsbCwKICAicHJpdmF0ZUlwIiA6ICIxNzIuMzEuNjMuNjAiLAogICJhdmFpbGFiaWxpdHlab25lIiA6ICJ1cy1lYXN0LTFjIiwKICAidmVyc2lvbiIgOiAiMjAxMC0wOC0zMSIsCiAgImluc3RhbmNlSWQiIDogImktZGUwZjEzNDQiLAogICJiaWxsaW5nUHJvZHVjdHMiIDogbnVsbCwKICAiaW5zdGFuY2VUeXBlIiA6ICJ0Mi5taWNybyIsCiAgImFjY291bnRJZCIgOiAiMjQxNjU2NjE1ODU5IiwKICAiaW1hZ2VJZCIgOiAiYW1pLWZjZTNjNjk2IiwKICAicGVuZGluZ1RpbWUiIDogIjIwMTYtMDQtMDVUMTY6MjY6NTVaIiwKICAiYXJjaGl0ZWN0dXJlIiA6ICJ4ODZfNjQiLAogICJrZXJuZWxJZCIgOiBudWxsLAogICJyYW1kaXNrSWQiIDogbnVsbCwKICAicmVnaW9uIiA6ICJ1cy1lYXN0LTEiCn0AAAAAAAAxggEXMIIBEwIBATBpMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQwIJAJa6SNnlXhpnMAkGBSsOAwIaBQCgXTAYBgkqhkiG9w0BCQMxCwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0xNjA0MDUxNjI3MDBaMCMGCSqGSIb3DQEJBDEWBBRtiynzMTNfTw1TV/d8NvfgVw+XfTAJBgcqhkjOOAQDBC4wLAIUVfpVcNYoOKzN1c+h1Vsm/c5U0tQCFAK/K72idWrONIqMOVJ8Uen0wYg4AAAAAAAA nonce=vault-client-nonce
```


### Via the API

#### Enable AWS EC2 authentication in Vault.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/sys/auth/aws" -d '{"type":"aws"}'
```

#### Configure the credentials required to make AWS API calls.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/auth/aws/config/client" -d '{"access_key":"VKIAJBRHKH6EVTTNXDHA", "secret_key":"vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj", "region":"us-east-1"}'
```

#### Configure the policies on the AMI.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/auth/aws/image/ami-fce3c696" -d '{"policies":"prod,dev","max_ttl":"500h"}'
```

#### Perform the login operation

```
curl -X POST "http://127.0.0.1:8200/v1/auth/aws/login" -d '{"pkcs7":"MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGmewogICJkZXZwYXlQcm9kdWN0Q29kZXMiIDogbnVsbCwKICAicHJpdmF0ZUlwIiA6ICIxNzIuMzEuNjMuNjAiLAogICJhdmFpbGFiaWxpdHlab25lIiA6ICJ1cy1lYXN0LTFjIiwKICAidmVyc2lvbiIgOiAiMjAxMC0wOC0zMSIsCiAgImluc3RhbmNlSWQiIDogImktZGUwZjEzNDQiLAogICJiaWxsaW5nUHJvZHVjdHMiIDogbnVsbCwKICAiaW5zdGFuY2VUeXBlIiA6ICJ0Mi5taWNybyIsCiAgImFjY291bnRJZCIgOiAiMjQxNjU2NjE1ODU5IiwKICAiaW1hZ2VJZCIgOiAiYW1pLWZjZTNjNjk2IiwKICAicGVuZGluZ1RpbWUiIDogIjIwMTYtMDQtMDVUMTY6MjY6NTVaIiwKICAiYXJjaGl0ZWN0dXJlIiA6ICJ4ODZfNjQiLAogICJrZXJuZWxJZCIgOiBudWxsLAogICJyYW1kaXNrSWQiIDogbnVsbCwKICAicmVnaW9uIiA6ICJ1cy1lYXN0LTEiCn0AAAAAAAAxggEXMIIBEwIBATBpMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQwIJAJa6SNnlXhpnMAkGBSsOAwIaBQCgXTAYBgkqhkiG9w0BCQMxCwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0xNjA0MDUxNjI3MDBaMCMGCSqGSIb3DQEJBDEWBBRtiynzMTNfTw1TV/d8NvfgVw+XfTAJBgcqhkjOOAQDBC4wLAIUVfpVcNYoOKzN1c+h1Vsm/c5U0tQCFAK/K72idWrONIqMOVJ8Uen0wYg4AAAAAAAA","nonce":"ault-client-nonce"}'
```


The response will be in JSON. For example:

```javascript
{
  "auth": {
    "renewable": true,
    "lease_duration": 1800000,
    "metadata": {
      "role_tag_max_ttl": "0",
      "instance_id": "i-de0f1344"
    },
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "accessor": "20b89871-e6f2-1160-fb29-31c2f6d4645e",
    "client_token": "c9368254-3f21-aded-8a6f-7c818e81b17a"
  },
  "warnings": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

## API
### /auth/aws/config/client
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the credentials required to perform API calls to AWS.
    The instance identity document fetched from the PKCS#7 signature
    will provide the EC2 instance ID. The credentials configured using
    this endpoint will be used to query the status of the instances via
    DescribeInstanceStatus API. Also, if the login is performed using
    the role tag, then these credentials will also be used to fetch the
    tags that are set on the EC2 instance via DescribeTags API. If the
    static credentials are not provided using this endpoint, then the
    credentials will be retrieved from the environment variables
    `AWS_ACCESS_KEY`, `AWS_SECRET_KEY` and `AWS_REGION` respectively.
    If the credentials are still not found and if the backend is configured
    on an EC2 instance with metadata querying capabilities, the credentials
    are fetched automatically.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/client`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">access_key</span>
        <span class="param-flags">required</span>
        AWS Access key with permissions to query EC2 instance metadata.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">secret_key</span>
        <span class="param-flags">required</span>
        AWS Secret key with permissions to query EC2 instance metadata.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">region</span>
        <span class="param-flags">required</span>
        Region for API calls. Defaults to the value of the AWS_REGION env var.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

#### GET
<dl class="api">
  <dt>Description</dt>
    Returns the previously configured AWS access credentials.
  <dd>
    
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/client`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```
{
  "auth": null,
  "warnings": null,
  "data": {
    "secret_key": "vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj",
    "region": "us-east-1",
    "access_key": "VKIAJBRHKH6EVTTNXDHA"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the previously configured AWS access credentials.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/client`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/config/certificate
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Registers an AWS public key that is used to verify the PKCS#7 signature of the
    EC2 instance metadata.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/certificate`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">aws_public_key</span>
        <span class="param-flags">required</span>
        AWS Public key required to verify PKCS7 signature of the EC2 instance metadata.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns the previously configured AWS public key.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/certificate`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "aws_public_cert": "-----BEGIN CERTIFICATE-----\nMIIC7TCCAq0CCQCWukjZ5V4aZzAJBgcqhkjOOAQDMFwxCzAJBgNVBAYTAlVTMRkw\nFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYD\nVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQzAeFw0xMjAxMDUxMjU2MTJaFw0z\nODAxMDUxMjU2MTJaMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9u\nIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNl\ncnZpY2VzIExMQzCCAbcwggEsBgcqhkjOOAQBMIIBHwKBgQCjkvcS2bb1VQ4yt/5e\nih5OO6kK/n1Lzllr7D8ZwtQP8fOEpp5E2ng+D6Ud1Z1gYipr58Kj3nssSNpI6bX3\nVyIQzK7wLclnd/YozqNNmgIyZecN7EglK9ITHJLP+x8FtUpt3QbyYXJdmVMegN6P\nhviYt5JH/nYl4hh3Pa1HJdskgQIVALVJ3ER11+Ko4tP6nwvHwh6+ERYRAoGBAI1j\nk+tkqMVHuAFcvAGKocTgsjJem6/5qomzJuKDmbJNu9Qxw3rAotXau8Qe+MBcJl/U\nhhy1KHVpCGl9fueQ2s6IL0CaO/buycU1CiYQk40KNHCcHfNiZbdlx1E9rpUp7bnF\nlRa2v1ntMX3caRVDdbtPEWmdxSCYsYFDk4mZrOLBA4GEAAKBgEbmeve5f8LIE/Gf\nMNmP9CM5eovQOGx5ho8WqD+aTebs+k2tn92BBPqeZqpWRa5P/+jrdKml1qx4llHW\nMXrs3IgIb6+hUIB+S8dz8/mmO0bpr76RoZVCXYab2CZedFut7qc3WUH9+EUAH5mw\nvSeDCOUMYQR7R9LINYwouHIziqQYMAkGByqGSM44BAMDLwAwLAIUWXBlk40xTwSw\n7HX32MxXYruse9ACFBNGmdX2ZBrVNGrN9N2f6ROk0k9K\n-----END CERTIFICATE-----\n"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


### /auth/aws/image/<ami_id>
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Registers an AMI ID in the backend. Only those instances which are using the AMIs registered using this endpoint,
    will be able to perform login operation. If each EC2 instance is using unique AMI ID, then all those AMI IDs should
    be registered beforehand. In case the same AMI is shared among many EC2 instances, then that AMI should be registered
    using this endpoint with the option `role_tag` (refer API section), then a `roletag` should be created using
    `auth/aws/image/<ami_id>/roletag` endpoint, and this tag should be attached to the EC2 instance before the login operation
    is performed.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/image/<ami_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ami_id</span>
        <span class="param-flags">required</span>
        AMI ID to be mapped.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">role_tag</span>
        <span class="param-flags">optional</span>
        If set, enables the `roletag` login for this AMI, meaning that this AMI is shared among many EC2 instances. The value set for this field should be the `key` of the tag on the EC2 instance and the `tag_value` returned from `auth/aws/image/<ami_id>/roletag` should be the `value` of the tag on the instance. Defaults to empty string, meaning that this AMI is not shared among instances.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">optional</span>
        The maximum allowed lease duration.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Policies to be associated with the AMI.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">allow_instance_migration</span>
        <span class="param-flags">optional</span>
        If set, allows migration of the underlying instance where the client resides. This keys off of pendingTime in the metadata document, so essentially, this disables the client nonce check whenever the instance is migrated to a new host and pendingTime is newer than the previously-remembered time. Use with caution.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">disallow_reauthentication</span>
        <span class="param-flags">optional</span>
        If set, only allows a single token to be granted per instance ID. In order to perform a fresh login, the entry in whitelist for the instance ID needs to be cleared using 'auth/aws/whitelist/identity/<instance_id>' endpoint. Defaults to 'false'.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns the previously registered AMI ID configuration.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/image/<ami_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "role_tag": "",
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "max_ttl": 1800000,
    "disallow_reauthentication": false,
    "allow_instance_migration": false
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists all the AMI IDs that are registered with the backend.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/images?list=true`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "keys": [
      "ami-fce3c696",
      "ami-hei3d687"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes the previously registered AMI ID.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/image/<ami_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/image/<ami_id>/roletag
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Creates a `roletag` for the AMI_ID. Role tags provide an effective way to restrict the
    options that are set on the AMI ID. This is of use when AMI is shared by multiple instances
    and there is need to customize the options for specific instances.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/image/<ami_id>/roletag`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">ami_id</span>
        <span class="param-flags">required</span>
        AMI ID to create a tag for.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Policies to be associated with the tag.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">optional</span>
        The maximum allowed lease duration.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">disallow_reauthentication</span>
        <span class="param-flags">optional</span>
        If set, only allows a single token to be granted per instance ID. This can be cleared with the auth/aws/whitelist/identity endpoint. Defaults to 'false'.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "tag_value": "v1:09Vp0qGuyB8=:a=ami-fce3c696:p=default,prod:d=false:t=300h0m0s:uPLKCQxqsefRhrp1qmVa1wsQVUXXJG8UZP/pJIdVyOI=",
    "tag_key": "VaultRole"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


### /auth/aws/login
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
   Login and fetch a token. If the instance metadata signature is valid
   along with a few other conditions, a token will be issued. 
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/login`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">pkcs7</span>
        <span class="param-flags">required</span>
        PKCS7 signature of the identity document.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">required/optional, depends</span>
        The `nonce` created by a client of this backend. When `disallow_reauthentication`
        option is enabled on either the AMI or the role tag, then `nonce` parameter is
        optional. It is a required parameter otherwise.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": {
    "renewable": true,
    "lease_duration": 1800000,
    "metadata": {
      "role_tag_max_ttl": "0",
      "instance_id": "i-de0f1344"
    },
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "accessor": "20b89871-e6f2-1160-fb29-31c2f6d4645e",
    "client_token": "c9368254-3f21-aded-8a6f-7c818e81b17a"
  },
  "warnings": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


### /auth/aws/blacklist/roletag/<role_tag>
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Places a valid roletag in a blacklist. This ensures that the `roletag`
    cannot be used by any instance to perform a login operation again.
    Note that if this `roletag` was previousy used to perfom a successful
    login, placing the `roletag` in the blacklist does not invalidate the
    already issued token.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/blacklist/roletag/<role_tag>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role_tag</span>
        <span class="param-flags">required</span>
        Role tag that needs be blacklisted. The tag can be supplied as-is, or can be base64 encoded.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns the blacklist entry of a previously blacklisted `roletag`.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/blacklist/roletag/<role_tag>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "expiration_time": "2016-04-25T10:35:20.127058773-04:00",
    "creation_time": "2016-04-12T22:35:01.178348124-04:00"
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists all the `roletags` that are blacklisted.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/blacklist/roletag?list=true`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "keys": [
      "v1:09Vp0qGuyB8=:a=ami-fce3c696:p=default,prod:d=false:t=300h0m0s:uPLKCQxqsefRhrp1qmVa1wsQVUXXJG8UZP/"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a blacklisted `roletag`.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/blacklist/roletag/<role_tag>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/blacklist/roletag/tidy
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Cleans up the entries in the blacklist based on expiration time on the entry and `safety_buffer`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/blacklist/roletag/tidy`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">safety_buffer</span>
        <span class="param-flags">optional</span>
        The amount of extra time that must have passed beyond the `roletag` expiration, before it is removed from the backend storage. Defaults to 72h.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/whitelist/identity/<instance_id>
#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns an entry in the whitelist. An entry will be created/updated by every successful login.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/whitelist/identity/<instance_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">instance_id</span>
        <span class="param-flags">required</span>
        EC2 instance ID. A successful login operation from an EC2 instance gets cached in this whitelist, keyed off of instance ID.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists all the instance IDs that are in the whitelist of successful logins.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/whitelist/identity?list=true`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


#### DELETE
<dl class="api">
  <dt>Description</dt>
  <dd>
    Deletes a cache of the successful login from an instance.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/whitelist/identity/<instance_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/whitelist/identity/tidy
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Cleans up the entries in the whitelist based on expiration time and `safety_buffer`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/whitelist/identity/tidy`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">safety_buffer</span>
        <span class="param-flags">optional</span>
        The amount of extra time that must have passed beyond the identity expiration, before it is removed from the backend storage. Defaults to 72h.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
