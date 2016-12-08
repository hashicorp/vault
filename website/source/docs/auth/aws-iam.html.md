---
layout: "docs"
page_title: "Auth Backend: AWS-IAM"
sidebar_current: "docs-auth-aws-iam"
description: |-
  The aws-iam backend allows automated authentication AWS IAM principals.
---

# Auth Backend: aws-iam

The aws-iam auth backend provides a secure introduction mechanism for AWS IAM
principals. This allows you to reduce the problem of securely introducing a
Vault token to the problem of securely introducing AWS IAM credentials, which
AWS has already solved in a number of use cases, such as via EC2 Instance
Profiles and IAM roles attached to Lambda functions. It also allows you to have
a consistent workflow between developers working on a local laptop with tools
such as Hologram.

## Comparison with AWS-EC2 Authentication Backend

The AWS-IAM and AWS-EC2 authentication backends serve similar purposes. Both
look to authenticate some type of AWS entity to Vault. However, in many ways,
the similarity ends there. The following is a comparison of the two entities:

* What type of entity is authenticated:
  * The AWS-EC2 backend authenticates AWS EC2 instances only.
  * The AWS-IAM backend authenticates AWS IAM principals. This can include
    IAM users, IAM roles assumed from other accounts, AWS Lambdas that are
    launched in an IAM role, or even EC2 instances that are launched in an
    EC2 instance profile.
* How the entities are authenticated
  * The AWS-EC2 backend authenticates instances by making use of the EC2 instance
    identity document, which is a cryptographically signed document containing
    metadata about the instance. This document changes relatively infrequently,
    so Vault adds a number of other constructs to mitigate against replay
    attacks, such as client nonces, role tags, instance migrations, etc.
  * The AWS-IAM backend authenticates by having clients provide a specially
    signed AWS API request which the backend then passes on to AWS to validate
    the signature and tell Vault who created it. The actual secret (i.e.,
    the AWS secret access key) is never transmitted over the wire, and the
    AWS signature algorithm automatically expires requests after 15 minutes,
    providing simple and robust protection against replay attacks.
* Specific use cases
  * If you have a long-lived EC2 instance which you are unable to relaunch into
    an IAM instance profile, then the AWS-EC2 backend is probably the best
    solution for you. (While you could store long-lived AWS IAM user
    credentials on disk and use those to authenticate to Vault, that would not
    be recommended.)
  * If you have non-EC2 instance entities, such as IAM users, Lambdas in IAM
    roles, or developer laptops using [AdRoll's Hologram](https://github.com/AdRoll/hologram)
    then you would need to use the AWS-IAM backend.
  * If you have EC2 instances which are already in an IAM instance profile, then
    you could use either backend.

## Authentication Workflow

The AWS STS API includes a method,
[`sts:GetCallerIdentity`](http://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html),
which allow you to validate the identity of a caller. The client signs
a `GetCallerIdentity` query using the [AWS Signature v4
algorithm](http://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html) and
submits that signed query to the Vault server. The Vault server then forwards
it on to the AWS STS service and validates the result back. Clients don't even
need network-level access to talk to the AWS STS API endpoint; they merely need
to sign the credentials. However, it means that the Vault server DOES need
network-level access to send requests to the STS endpoint.

Each signed AWS request includes the current timestamp to mitigate the risk of
replay attacks. In addition, Vault allows you to require an additional header,
`X-Vault-AWSIAM-Server-ID`, to be present to mitigate against different types of replay
attacks (such as a signed `GetCallerIdentity` request stolen from a dev Vault
instance and used to authenticate to a prod Vault instance). Vault further
requires that this header be one of the headers included in the AWS signature
and relies upon AWS to authenticate that signature.

While the AWS API endpoints support both signed GET and POST requests, for
simplicity, the aws-iam backend supports only POST requests.

## Authorization Workflow

Roles are mapped from AWS IAM principals back to Vault roles. At present, two
types of principals are supported: AWS IAM Users and AWS IAM roles. For users,
you simply bind the full ARN of the user, e.g., `arn:aws:iam::123456789012:user/MyUserName`
(where `123456789012` is your AWS account number nad `MyUserName` is your IAM
username). Assumed roles are a little tricky: when authenticating with IAM roles,
you are simultaneously _two_ different principals: `arn:aws:iam::123456789012:role/MyRoleName`
*and* `arn:aws:sts::123456789012:assumed-role/MyRoleName/RoleSessionName`. In
this case, you should bind the former ARN to your Vault roles.

## Authentication

### Via the CLI

#### Enable AWS IAM authentication in Vault.

```
$ vault auth-enable aws-iam
```

#### Configure the policies on the role.

```
$ vault write auth/aws-iam/role/dev-role bound_iam_principal=arn:aws:iam::123456789012:role/my_role policies=prod,dev max_ttl=500h
```

#### Configure a required X-Vault-AWSIAM-Server-ID Header (recommended)

```
$ vault write auth/aws-iam/client/config vault_header_vaule=vault.example.xom
```

#### Perform the login operation

Generating the signed request is a non-standard operation. The Vault cli supports generating this for you:

```
$ vault auth -method=aws-iam header_value=vault.example.com role=dev-role
```

This assumes you have AWS credentials configured in the standard locations AWS SDKs search for credentials
(environment variables, ~/.aws/credentials, EC2 instance profile in that order). If you do not have
IAM credentials available at any of these locations, you can explicitly pass them in on the command line
(though this is not recommended):
```
$ vault auth -method aws-iam header_value=vault.example.com role=dev-role \
        aws_access_key_id=<access_key> \
        aws_secret_access_key=<secret_key> \
        aws_security_token=<security_token>
```

For reference, the following Go program also demonstrates how to generate the
required parameters (assuming you are using a default AWS credential provider):

```
package main

import (
        "encoding/base64"
        "encoding/json"
        "fmt"
        "io/ioutil"

        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/sts"
)

func transformHeaders(input map[string][]string) map[string]string {
        retval := map[string]string{}
        for k, v := range input {
                retval[k] = v[0]
        }
        return retval
}

func main() {
        sess, err := session.NewSession()
        if err != nil {
                fmt.Println("failed to create session,", err)
                return
        }

        svc := sts.New(sess)
        var params *sts.GetCallerIdentityInput
        stsRequest, _ := svc.GetCallerIdentityRequest(params)
        stsRequest.HTTPRequest.Header.Add("X-Vault-AWSIAM-Server-ID", "vault.example.com")
        stsRequest.Sign()

        headersJson, err := json.Marshal(transformHeaders(stsRequest.HTTPRequest.Header))
        if err != nil {
                fmt.Println(fmt.Errorf("Error:", err))
                return
        }
        requestBody, err := ioutil.ReadAll(stsRequest.HTTPRequest.Body)
        if err != nil {
                fmt.Println(fmt.Errorf("Error:", err))
                return
        }
        fmt.Println("method=" + stsRequest.HTTPRequest.Method)
        fmt.Println("url=" + stsRequest.HTTPRequest.URL.String())
        fmt.Println("headers=" + base64.StdEncoding.EncodeToString(headersJson))
        fmt.Println("body=" + base64.StdEncoding.EncodeToString(requestBody))
}
```
Using this, we can get the values to pass in to the `vault write` operation:

```
$ vault write auth/aws-iam/login role=dev method=POST url=https://sts.amazonaws.com/ headers=eyJBdXRob3JpemF0aW9uIjoiQVdTNC1ITUFDLVNIQTI1NiBDcmVkZW50aWFsPWZvby8yMDE2MDkzMC91cy1lYXN0LTEvc3RzL2F3czRfcmVxdWVzdCwgU2lnbmVkSGVhZGVycz1jb250ZW50LWxlbmd0aDtjb250ZW50LXR5cGU7aG9zdDt4LWFtei1kYXRlO3gtdmF1bHQtc2VydmVyLCBTaWduYXR1cmU9YTY5ZmQ3NTBhMzQ0NWM0ZTU1M2UxYjNlNzlkM2RhOTBlZWY1NDA0N2YxZWI0ZWZlOGZmYmM5YzQyOGMyNjU1YiIsIkNvbnRlbnQtTGVuZ3RoIjoiNDMiLCJDb250ZW50LVR5cGUiOiJhcHBsaWNhdGlvbi94LXd3dy1mb3JtLXVybGVuY29kZWQ7IGNoYXJzZXQ9dXRmLTgiLCJVc2VyLUFnZW50IjoiYXdzLXNkay1nby8xLjQuMTIgKGdvMS43LjE7IGxpbnV4OyBhbWQ2NCkiLCJYLUFtei1EYXRlIjoiMjAxNjA5MzBUMDQzMTIxWiIsIlgtVmF1bHQtU2VydmVyIjoidmF1bHQuZGV2Lmp0aG9tcHNvbi5pbyJ9 body=QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ==
```

### Via the API

#### Enable AWS IAM authentication in Vault.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/sys/auth/aws-iam" -d '{"type":"aws-iam"}'
```

#### Configure the policies on the role.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/auth/aws-iam/role/dev-role -d '{"bound_iam_principal":"arn:aws:iam::123456789012:role/my_role","policies":"prod,dev","max_ttl":"500h"}'
```

#### Perform the login operation

```
curl -X POST "http://127.0.0.1:8200/v1/auth/aws-iam/login" -d '{"role":"dev", "method": "POST", "url": "https://sts.amazonaws.com/", "body": "QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ==", "headers": "eyJBdXRob3JpemF0aW9uIjoiQVdTNC1ITUFDLVNIQTI1NiBDcmVkZW50aWFsPWZvby8yMDE2MDkzMC91cy1lYXN0LTEvc3RzL2F3czRfcmVxdWVzdCwgU2lnbmVkSGVhZGVycz1jb250ZW50LWxlbmd0aDtjb250ZW50LXR5cGU7aG9zdDt4LWFtei1kYXRlO3gtdmF1bHQtc2VydmVyLCBTaWduYXR1cmU9YTY5ZmQ3NTBhMzQ0NWM0ZTU1M2UxYjNlNzlkM2RhOTBlZWY1NDA0N2YxZWI0ZWZlOGZmYmM5YzQyOGMyNjU1YiIsIkNvbnRlbnQtTGVuZ3RoIjoiNDMiLCJDb250ZW50LVR5cGUiOiJhcHBsaWNhdGlvbi94LXd3dy1mb3JtLXVybGVuY29kZWQ7IGNoYXJzZXQ9dXRmLTgiLCJVc2VyLUFnZW50IjoiYXdzLXNkay1nby8xLjQuMTIgKGdvMS43LjE7IGxpbnV4OyBhbWQ2NCkiLCJYLUFtei1EYXRlIjoiMjAxNjA5MzBUMDQzMTIxWiIsIlgtVmF1bHQtU2VydmVyIjoidmF1bHQuZGV2Lmp0aG9tcHNvbi5pbyJ9" }'
```


The response will be in JSON. For example:

```javascript
{
    "auth": {
        "accessor": "52033390-a416-9aaf-e8d9-68b45947ce59",
        "client_token": "8b1d634a-6d02-7452-e2b9-9cec9e12cbf7",
        "lease_duration": 2592000,
        "metadata": {
            "canonical_arn": "arn:aws:iam::123456789012:role/MyRole",
            "client_arn": "arn:aws:sts::123456789012:assumed-role/MyRole/RoleSessionName"
        },
        "policies": [
            "default"
        ],
        "renewable": true
    },
    "data": null,
    "lease_duration": 0,
    "lease_id": "",
    "renewable": false,
    "request_id": "84915669-8386-5c99-2d7b-131cac038d52",
    "warnings": null,
    "wrap_info": null
}
```

## API
### /auth/aws-iam/config/client
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the behavior of Vault's STS client. 
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/config/client`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">endpoint</span>
        <span class="param-flags">optional</span>
        URL to override the default generated endpoint for making AWS STS API calls.
        This is useful if you want to use STS in a [region other than
us-east-1](http://docs.aws.amazon.com/general/latest/gr/rande.html#sts_region).
      </li>
      <li>
        <span class="param">vault_header_value</span>
        <span class="param-flags">optional</span>
        Value to require be present in the X-Vault-AWSIAM-Server-ID header. If not set, then
        no value is required or validated. If set, clients must include an X-VaultServer
        header in the headers of login requests, and further the X-VaultServer must be
        among the signed headers validated by AWS. This is to protect against different types
        of replay attacks, for example a signed request sent to a dev server being
        resent to a production server. Consider setting this to the Vault server's DNS name.
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
    Returns the previously configured AWS STS client.
  <dd>

  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/config/client`</dd>

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
    "endpoint": "",
    "vault_header_value": "vault.example.com"
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
    Deletes the previously configured AWS STS client configuration.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/config/client`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>



### /auth/aws-iam/role/[role]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Registers a role in the backend. Only those IAM principals which have been
mapped to a role will be able to perform the login operation.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/role/<role>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role</span>
        <span class="param-flags">required</span>
        Name of the role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_iam_principal</span>
        <span class="param-flags">required</span>
        Defines the ARN of the IAM principal to map to the role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The TTL period of tokens issued using this role, provided as "1h", where hour is
        the largest suffix.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">optional</span>
        The maximum allowed lifetime of tokens issued using this role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Policies to be set on tokens issued using this role.
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
    Returns the previously registered role configuration.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/role/<role>`</dd>

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
    "bound_iam_principal": "arn:aws:iam::123456789012:role/my_role"
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "max_ttl": 1800000
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
    Lists all the roles that are registered with the backend.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/roles?list=true`</dd>

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
      "dev-role",
      "prod-role"
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
    Deletes the previously registered role.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/role/<role>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>



### /auth/aws-iam/login
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
   Fetch a token. This endpoint verifies the query signature with AWS and
ensure the query was signed by a valid IAM principal. Note that the HTTP request
mentioned in the parameter descriptions below refers to the signed GetCallerIdentity
STS request that is embedded in the request to the Vault server, NOT the request to
the Vault server itself.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws-iam/login`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role</span>
        <span class="param-flags">optional</span>
        Name of the Vault role against which the login is being attempted.
        If `role` is not specified, then the login endpoint looks for a role
        bearing the name of the AMI ID of the EC2 instance that is trying to login.
        If a matching role is not found, login fails.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">method</span>
        <span class="param-flags">required</span>
        HTTP method used in the signed request. Currently only POST is supported, but
        other methods may be supported in the future.
      </li>
      <li>
        <span class="param">url</span>
        <span class="param-flags">required</span>
        HTTP URL used in the signed request. Most likely just https://sts.amazonaws.com/
        as most requests will probably use POST with an empty URI.
      </li>
      <li>
        <span class="param">body</span>
        <span class="param-flags">required</span>
        Base64-encoded body of the signed request. Most likely
        <em>QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ==</em>
        which is the base64 encoding of <em>Action=GetCallerIdentity&Version=2011-06-15</em>
      </li>
      <li>
        <span class="param">headers</span>
        <span class="param-flags">required</span>
        Base64-encoded, JSON-serialized representation of the HTTP request headers. The JSON
        serialization assumes that each header key only has a single string value. While the
        HTTP spec allows for a given header to have multiple values, that is extremely unlikely
        to occur with any request sent to STS.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>

```javascript
{
    "auth": {
        "accessor": "52033390-a416-9aaf-e8d9-68b45947ce59",
        "client_token": "8b1d634a-6d02-7452-e2b9-9cec9e12cbf7",
        "lease_duration": 2592000,
        "metadata": {
            "canonical_arn": "arn:aws:iam::854766835649:user/joelt49",
            "client_arn": "arn:aws:iam::854766835649:user/joelt49"
        },
        "policies": [
            "default"
        ],
        "renewable": true
    },
    "data": null,
    "lease_duration": 0,
    "lease_id": "",
    "renewable": false,
    "request_id": "84915669-8386-5c99-2d7b-131cac038d52",
    "warnings": null,
    "wrap_info": null
}
```

  </dd>
</dl>
