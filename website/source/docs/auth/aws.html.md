---
layout: "docs"
page_title: "Auth Backend: AWS"
sidebar_current: "docs-auth-aws"
description: |-
  The aws backend allows automated authentication of AWS entities.
---

# Auth Backend: aws

The aws auth backend provides an automated mechanism to retrieve
a Vault token for AWS EC2 instances and IAM principals.  Unlike most Vault
authentication backends, this backend does not require manual first-deploying, or
provisioning security-sensitive credentials (tokens, username/password, client
certificates, etc), by operators under many circumstances. It treats
AWS as a Trusted Third Party and uses either
the cryptographically signed dynamic metadata information that uniquely
represents each EC2 instance or a special AWS request signed with AWS IAM
credentials. The metadata information is automatically supplied by AWS to all
EC2 instances, and IAM credentials are automatically supplied to AWS instances
in IAM instance profiles, Lambda functions, and others, and it is this
information already provided by AWS which Vault can use to authenticate
clients.

## Authentication Workflow

There are two authentication types present in the aws backend: `ec2` and `iam`.
Based on how you attempt to authenticate, Vault will determine if you are
attempting to use the `ec2` or `iam` type.  Each has a different authentication
workflow, and each can solve different use cases.  See the section on comparing
the two auth methods below to help determine which method is more appropriate
for your use cases.

### EC2 Authentication Method

EC2 instances have access to metadata describing the instance. (For those not
familiar with instance metadata, details can be found
[here](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html).)

One piece of "dynamic metadata" available to the EC2 instance, is the instance
identity document, a JSON representation of a collection of instance metadata.
AWS also provides PKCS#7 signature of the instance metadata document, and
publishes the public keys (grouped by region) which can be used to verify the
signature. Details on the instance identity document and the signature can be
found
[here](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html).

During login, the backend verifies the signature on the PKCS#7 document,
ensuring that the information contained within, is certified accurate by AWS.
Before succeeding the login attempt and returning a Vault token, the backend
verifies the current running status of the instance via the EC2 API.

There are various modifications to this workflow that provide more or less
security, as detailed later in this documentation.

### IAM Authentication Method

The AWS STS API includes a method,
[`sts:GetCallerIdentity`](http://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html),
which allows you to validate the identity of a client. The client signs
a `GetCallerIdentity` query using the [AWS Signature v4
algorithm](http://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html) and
submits 4 pieces of information to the Vault server to recreate a valid signed
request: the request URL, the request body, the request headers, and the request
method, as the AWS signature is computed over those fields. The Vault server
then reconstructs the query and forwards it on to the AWS STS service and
validates the result back. Clients don't need network-level access to talk to
the AWS STS API endpoint; they merely need access to the credentials to sign the
request. However, it means that the Vault server does need network-level access
to send requests to the STS endpoint.

Importantly, the credentials used to sign the GetCallerIdentity request can come
from the EC2 instance metadata service for an EC2 instance, or from the AWS
environment variables in an AWS Lambda function execution, which obviates the
need for an operator to manually provision some sort of identity material first.
However, the credentials can, in principle, come from anywhere, not just from
the locations AWS has provided for you.

Each signed AWS request includes the current timestamp to mitigate the risk of
replay attacks. In addition, Vault allows you to require an additional header,
`X-Vault-AWS-IAM-Server-ID`, to be present to mitigate against different types of replay
attacks (such as a signed `GetCallerIdentity` request stolen from a dev Vault
instance and used to authenticate to a prod Vault instance). Vault further
requires that this header be one of the headers included in the AWS signature
and relies upon AWS to authenticate that signature.

While AWS API endpoints support both signed GET and POST requests, for
simplicity, the aws backend supports only POST requests. It also does not
support `presigned` requests, i.e., requests with `X-Amz-Credential`,
`X-Amz-signature`, and `X-Amz-SignedHeaders` GET query parameters containing the
authenticating information.

It's also important to note that Amazon does NOT appear to include any sort
of authorization around calls to `GetCallerIdentity`. For example, if you have
an IAM policy on your credential that requires all access to be MFA authenticated,
non-MFA authenticated credentials (i.e., raw credentials, not those retrieved
by calling `GetSessionToken` and supplying an MFA code) will still be able to
authenticate to Vault using this backend. It does not appear possible to enforce
an IAM principal to be MFA authenticated while authenticating to Vault.

## Authorization Workflow

The basic mechanism of operation is per-role. Roles are registered in the
backend and associated with a specific authentication type that cannot be
changed once the role has been created. Roles can also be associated with
various optional restrictions, such as the set of allowed policies and max TTLs
on the generated tokens. Each role can be specified with the constraints that
are to be met during the login. For example, one such constraint that is
supported is to bind against AMI ID. A role which is bound to a specific AMI,
can only be used for login by EC2 instances that are deployed on the same AMI.

In general, role bindings that are specific to an EC2 instance are only checked
when the ec2 auth method is used to login, while bindings specific to IAM
principals are only checked when the iam auth method is used to login. However,
the iam method includes the ability for you to "infer" an EC2 instance ID from
the authenticated client and apply many of the bindings that would otherwise
only apply specifically to EC2 instances.

In many cases, an organization will use a "seed AMI" that is specialized after
bootup by configuration management or similar processes. For this reason, a
role entry in the backend can also be associated with a "role tag" when using
the ec2 auth type. These tags
are generated by the backend and are placed as the value of a tag with the
given key on the EC2 instance. The role tag can be used to further restrict the
parameters set on the role, but cannot be used to grant additional privileges.
If a role with an AMI bind constraint has "role tag" enabled on the role, and
the EC2 instance performing login does not have an expected tag on it, or if the
tag on the instance is deleted for some reason, authentication fails.

The role tags can be generated at will by an operator with appropriate API
access. They are HMAC-signed by a per-role key stored within the backend, allowing
the backend to verify the authenticity of a found role tag and ensure that it has
not been tampered with. There is also a mechanism to blacklist role tags if one
has been found to be distributed outside of its intended set of machines.

## IAM Authentication Inferences

With the iam auth method, normally Vault will see the IAM principal that
authenticated, either the IAM user or role. However, when you have an EC2
instance in an IAM instance profile, Vault can actually see the instance ID of
the instance and can "infer" that it's an EC2 instance. However, there are
important security caveats to be aware of before configuring Vault to make that
inference.

Each AWS IAM role has a "trust policy" which specifies which entities are
trusted to call
[`sts:AssumeRole`](http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html)
on the role and retrieve credentials that can be used to authenticate with that
role. When AssumeRole is called, a parameter called RoleSessionName is passed
in, which is chosen arbitrarily by the entity which calls AssumeRole. If you
have a role with an ARN `arn:aws:iam::123456789012:role/MyRole`, then the
credentials returned by calling AssumeRole on that role will be
`arn:aws:sts::123456789012:assumed-role/MyRole/RoleSessionName` where
RoleSessionName is the session name in the AssumeRole API call. It is this
latter value which Vault actually sees.

When you have an EC2 instance in an instance profile, the corresponding role's
trust policy specifies that the principal `"Service": "ec2.amazonaws.com"` is
trusted to call AssumeRole. When this is configured, EC2 calls AssumeRole on
behalf of your instance, with a RoleSessionName corresponding to the
instance's instance ID. Thus, it is possible for Vault to extract the instance
ID out of the value it sees when an EC2 instance in an instance profile
authenticates to Vault with the iam authentication method. This is known as
"inferencing." Vault can be configured, on a role-by-role basis, to infer that a
caller is an EC2 instance and, if so, apply further bindings that apply
specifically to EC2 instances -- most of the bindings available to the ec2
authentication backend.

However, it is very important to note that if any entity other than an AWS
service is permitted to call AssumeRole on your role, then that entity can
simply pass in your instance's instance ID and spoof your instance to Vault.
This also means that anybody who is able to modify your role's trust policy
(e.g., via
[`iam:UpdateAssumeRolePolicy`](http://docs.aws.amazon.com/IAM/latest/APIReference/API_UpdateAssumeRolePolicy.html),
then that person could also spoof your instances. If this is a concern but you
would like to take advantage of inferencing, then you should tightly restrict
who is able to call AssumeRole on the role, tightly restrict who is able to call
UpdateAssumeRolePolicy on the role, and monitor CloudTrail logs for calls to
AssumeRole and UpdateAssumeRolePolicy. All of these caveats apply equally to
using the iam authentication method without inferencing; the point is merely
that Vault cannot offer an iron-clad guarantee about the inference and it is up
to operators to determine, based on their own AWS controls and use cases,
whether or not it's appropriate to configure inferencing.

## Mixing Authentication Types

Vault allows you to configure using either the ec2 auth method or the iam auth
method, but not both auth methods. Further, Vault will prevent you from
enforcing restrictions that it cannot enforce given the chosen auth type for a
role. Some examples of how this works in practice:

1. You configure a role with the ec2 auth type, with a bound AMI ID. A
   client would not be able to login using the iam auth type.
2. You configure a role with the iam auth type, with a bound IAM
   principal ARN. A client would not be able to login with the ec2 auth method.
3. You configure a role with the iam auth type and further configure
   inferencing. You have a bound AMI ID and a bound IAM principal ARN. A client
   must login using the iam method; the RoleSessionName must be a valid instance
   ID viewable by Vault, and the instance must have come from the bound AMI ID.

## Comparison of the EC2 and IAM Methods

The iam and ec2 authentication methods serve similar and somewhat overlapping
functionality, in that both authenticate some type of AWS entity to Vault. To
help you determine which method is more appropriate for your use case, here is a
comparison of the two authentication methods.

* What type of entity is authenticated:
  * The ec2 auth method authenticates only AWS EC2 instances and is specialized
    to handle EC2 instances, such as restricting access to EC2 instances from
    a particular AMI, EC2 instances in a particular instance profile, or EC2
    instances with a specialized tag value (via the role_tag feature).
  * The iam auth method authenticates generic AWS IAM principals. This can
    include IAM users, IAM roles assumed from other accounts, AWS Lambdas that
    are launched in an IAM role, or even EC2 instances that are launched in an
    IAM instance profile. However, because it authenticates more generalized IAM
    principals, this backend doesn't offer more granular controls beyond binding
    to a given IAM principal without the use of inferencing.
* How the entities are authenticated
  * The ec2 auth method authenticates instances by making use of the EC2
    instance identity document, which is a cryptographically signed document
    containing metadata about the instance. This document changes relatively
    infrequently, so Vault adds a number of other constructs to mitigate against
    replay attacks, such as client nonces, role tags, instance migrations, etc.
    Because the instance identity document is signed by AWS, you have a strong
    guarantee that it came from an EC2 instance.
  * The iam auth method authenticates by having clients provide a specially
    signed AWS API request which the backend then passes on to AWS to validate
    the signature and tell Vault who created it. The actual secret (i.e.,
    the AWS secret access key) is never transmitted over the wire, and the
    AWS signature algorithm automatically expires requests after 15 minutes,
    providing simple and robust protection against replay attacks. The use of
    inferencing, however, provides a weaker guarantee that the credentials came
    from an EC2 instance in an IAM instance profile compared to the ec2
    authentication mechanism.
  * The instance identity document used in the ec2 auth method is more likely to
    be stolen given its relatively static nature, but it's harder to spoof. On
    the other hand, the credentials of an EC2 instance in an IAM instance
    profile are less likely to be stolen given their dynamic and short-lived
    nature, but it's easier to spoof credentials that might have come from an
    EC2 instance.
* Specific use cases
  * If you have non-EC2 instance entities, such as IAM users, Lambdas in IAM
    roles, or developer laptops using [AdRoll's
    Hologram](https://github.com/AdRoll/hologram) then you would need to use the
    iam auth method.
  * If you have EC2 instances, then you could use either auth method. If you
    need more granular filtering beyond just the instance profile of given EC2
    instances (such as filtering based off the AMI the instance was launched
    from), then you would need to use the ec2 auth method, change the instance
    profile associated with your EC2 instances so they have unique IAM roles
    for each different Vault role you would want them to authenticate
    to, or make use of inferencing. If you need to make use of role tags, then
    you will need to use the ec2 auth method.

## Client Nonce

Note: this only applies to the ec2 authentication method.

If an unintended party gains access to the PKCS#7 signature of the identity
document (which by default is available to every process and user that gains
access to an EC2 instance), it can impersonate that instance and fetch a Vault
token. The backend addresses this problem by using a Trust On First Use (TOFU)
mechanism that allows the first client to present the PKCS#7 signature of the
document to be authenticated and denying the rest. An important property of
this design is detection of unauthorized access: if an unintended party authenticates,
the intended client will be unable to authenticate and can raise an alert for
investigation.

During the first login, the backend stores the instance ID that authenticated
in a `whitelist`. One method of operation of the backend is to disallow any
authentication attempt for an instance ID contained in the whitelist, using the
`disallow_reauthentication` option on the role, meaning that an instance is
allowed to login only once. However, this has consequences for token rotation,
as it means that once a token has expired, subsequent authentication attempts
would fail. By default, reauthentication is enabled in this backend, and can be
turned off using `disallow_reauthentication` parameter on the registered role.

In the default method of operation, the backend will return a unique nonce
during the first authentication attempt, as part of auth `metadata`. Clients
should present this `nonce` for subsequent login attempts and it should match
the `nonce` cached at the identity-whitelist entry at the backend. Since only
the original client knows the `nonce`, only the original client is allowed to
reauthenticate. (This is the reason that this is a whitelist rather than a
blacklist; by default, it's keeping track of clients allowed to reauthenticate,
rather than those that are not.). Clients can choose to provide a `nonce` even
for the first login attempt, in which case the provided `nonce` will be tied to
the cached identity-whitelist entry. It is recommended to use a strong `nonce`
value in this case.

It is up to the client to behave correctly with respect to the nonce; if the
client stores the nonce on disk it can survive reboots, but could also give
access to other users or applications on the instance. It is also up to the
operator to ensure that client nonces are in fact unique; sharing nonces allows
a compromise of the nonce value to enable an attacker that gains access to any
EC2 instance to imitate the legitimate client on that instance. This is why
nonces can be disabled on the backend side in favor of only a single
authentication per instance; in some cases, such as when using ASGs, instances
are immutable and single-boot anyways, and in conjunction with a high max TTL,
reauthentication may not be needed (and if it is, the instance can simply be
shut down and allow ASG to start a new one).

In both cases, entries can be removed from the whitelist by instance ID,
allowing reauthentication by a client if the nonce is lost (or not used) and an
operator approves the process.

One other point: if available by the OS/distribution being used with the EC2
instance, it is not a bad idea to firewall access to the signed PKCS#7 metadata
to ensure that it is accessible only to the matching user(s) that require
access.

## Advanced Options and Caveats

### Dynamic Management of Policies Via Role Tags

Note: This only applies to the ec2 auth method or the iam auth method when
inferencing is used.

If the instance is required to have customized set of policies based on the
role it plays, the `role_tag` option can be used to provide a tag to set on
instances, for a given role. When this option is set, during login, along with
verification of PKCS#7 signature and instance health, the backend will query
for the value of a specific tag with the configured key that is attached to the
instance. The tag holds information that represents a *subset* of privileges that
are set on the role and are used to further restrict the set of the role's
privileges for that particular instance.

A `role_tag` can be created using `auth/aws/role/<role>/tag` endpoint
and is immutable. The information present in the tag is SHA256 hashed and HMAC
protected. The per-role key to HMAC is only maintained in the backend. This prevents
an adversarial operator from modifying the tag when setting it on the EC2 instance
in order to escalate privileges.

When 'role_tag' option is enabled on a role, the instances are required to have a
role tag. If the tag is not found on the EC2 instance, authentication will fail.
This is to ensure that privileges of an instance are never escalated for not
having the tag on it or for getting the tag removed. If the role tag creation does
not specify the policy component, the client will inherit the allowed policies set
on the role. If the role tag creation specifies the policy component but it contains
no policies, the token will contain only the `default` policy; by default, this policy
allows only manipulation (revocation, renewal, lookup) of the existing token, plus
access to its [cubbyhole](/docs/secrets/cubbyhole/index.html).
This can be useful to allow instances access to a secure "scratch space" for
storing data (via the token's cubbyhole) but without granting any access to
other resources provided by or resident in Vault.

### Handling Lost Client Nonces

Note: This only applies to the ec2 auth method.

If an EC2 instance loses its client nonce (due to a reboot, a stop/start of the
client, etc.), subsequent login attempts will not succeed. If the client nonce
is lost, normally the only option is to delete the entry corresponding to the
instance ID from the identity `whitelist` in the backend. This can be done via
the `auth/aws/identity-whitelist/<instance_id>` endpoint. This allows a new
client nonce to be accepted by the backend during the next login request.

Under certain circumstances there is another useful setting. When the instance
is placed onto a host upon creation, it is given a `pendingTime` value in the
instance identity document (documentation from AWS does not cover this option,
unfortunately). If an instance is stopped and started, the `pendingTime` value
is updated (this does not apply to reboots, however).

The backend can take advantage of this via the `allow_instance_migration`
option, which is set per-role. When this option is enabled, if the client nonce
does not match the saved nonce, the `pendingTime` value in the instance
identity document will be checked; if it is newer than the stored `pendingTime`
value, the backend assumes that the client was stopped/started and allows the
client to log in successfully, storing the new nonce as the valid nonce for
that client. This essentially re-starts the TOFU mechanism any time the
instance is stopped and started, so should be used with caution. Just like with
initial authentication, the legitimate client should have a way to alert (or an
alert should trigger based on its logs) if it is denied authentication.

Unfortunately, the `allow_instance_migration` only helps during stop/start
actions; the current metadata does not provide for a way to allow this
automatic behavior during reboots. The backend will be updated if this needed
metadata becomes available.

The `allow_instance_migration` option is set per-role, and can also be
specified in a role tag. Since role tags can only restrict behavior, if the
option is set to `false` on the role, a value of `true` in the role tag takes
effect; however, if the option is set to `true` on the role, a value set in the
role tag has no effect.

### Disabling Reauthentication

Note: this only applies to the ec2 authentication method.

If in a given organization's architecture, a client fetches a long-lived Vault
token and has no need to rotate the token, all future logins for that instance
ID can be disabled. If the option `disallow_reauthentication` is set, only one
login will be allowed per instance.  If the intended client successfully
retrieves a token during login, it can be sure that its token will not be
hijacked by another entity.

When `disallow_reauthentication` option is enabled, the client can choose not
to supply a nonce during login, although it is not an error to do so (the nonce
is simply ignored). Note that reauthentication is enabled by default. If only
a single login is desired, `disallow_reauthentication` should be set explicitly
on the role or on the role tag.

The `disallow_reauthentication` option is set per-role, and can also be
specified in a role tag. Since role tags can only restrict behavior, if the
option is set to `false` on the role, a value of `true` in the role tag takes
effect; however, if the option is set to `true` on the role, a value set in the
role tag has no effect.

### Blacklisting Role Tags

Note: this only applies to the ec2 authentication method or the iam auth method
when inferencing is used.

Role tags are tied to a specific role, but the backend has no control over, which
instances using that role, should have any particular role tag; that is purely up
to the operator. Although role tags are only restrictive (a tag cannot escalate
privileges above what is set on its role), if a role tag is found to have been
used incorrectly, and the administrator wants to ensure that the role tag has no
further effect, the role tag can be placed on a `blacklist` via the endpoint
`auth/aws/roletag-blacklist/<role_tag>`. Note that this will not invalidate the
tokens that were already issued; this only blocks any further login requests from
those instances that have the blacklisted tag attached to them.

### Expiration Times and Tidying of `blacklist` and `whitelist` Entries

The expired entries in both identity `whitelist` and role tag `blacklist` are
deleted automatically.  The entries in both of these lists contain an expiration
time which is dynamically determined by three factors: `max_ttl` set on the role,
`max_ttl` set on the role tag, and `max_ttl` value of the backend mount. The
least of these three dictates the maximum TTL of the issued token, and
correspondingly will be set as the expiration times of these entries.

The endpoints `aws/auth/tidy/identity-whitelist` and `aws/auth/tidy/roletag-blacklist` are
provided to clean up the entries present in these lists. These endpoints allow
defining a safety buffer, such that an entry must not only be expired, but be
past expiration by the amount of time dictated by the safety buffer in order
to actually remove the entry.

Automatic deletion of expired entries is performed by the periodic function
of the backend. This function does the tidying of both blacklist role tags
and whitelist identities. Periodic tidying is activated by default and will
have a safety buffer of 72 hours, meaning only those entries are deleted which
were expired before 72 hours from when the tidy operation is being performed.
This can be configured via `config/tidy/roletag-blacklist` and `config/tidy/identity-whitelist`
endpoints.

### Varying Public Certificates

Note: this only applies to the ec2 authentication method.

The AWS public certificate, which contains the public key used to verify the
PKCS#7 signature, varies for different AWS regions. The primary AWS public
certificate, which covers most AWS regions, is already included in Vault and
does not need to be added. Instances whose PKCS#7 signatures cannot be
verified by the default public certificate included in Vault can register a
different public certificate which can be found [here]
(http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html),
via the `auth/aws/config/certificate/<cert_name>` endpoint.

### Dangling Tokens

An EC2 instance, after authenticating itself with the backend, gets a Vault token.
After that, if the instance terminates or goes down for any reason, the backend
will not be aware of such events. The token issued will still be valid, until
it expires. The token will likely be expired sooner than its lifetime when the
instance fails to renew the token on time.

### Cross Account Access

To allow Vault to authenticate EC2 instances running in other accounts, AWS STS
(Security Token Service) can be used to retrieve temporary credentials by
assuming an IAM Role in those accounts. All these accounts should be configured
at the backend using the `auth/aws-ec2/config/sts/<account_id>` endpoint.

The account in which Vault is running (i.e. the master account) must be listed as
a trusted entity in the IAM Role being assumed on the remote account. The Role itself
must allow the `ec2:DescribeInstances` action, and `iam:GetInstanceProfile` if IAM Role
binding is used (see below).

Furthermore, in the master account, Vault must be granted the action `sts:AssumeRole`
for the IAM Role to be assumed.

## Authentication

### Via the CLI

#### Enable AWS EC2 authentication in Vault.

```
$ vault auth-enable aws
```

#### Configure the credentials required to make AWS API calls

If not specified, Vault will attempt to use standard environment variables
(`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) or IAM EC2 instance role
credentials if available.

The IAM account or role to which the credentials map must allow the
`ec2:DescribeInstances` action.  In addition, if IAM Role binding is used (see
`bound_iam_role_arn` below), `iam:GetInstanceProfile` must also be allowed.

```
$ vault write auth/aws/config/client secret_key=vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj access_key=VKIAJBRHKH6EVTTNXDHA
```

#### Configure the policies on the role.

```
$ vault write auth/aws/role/dev-role auth_type=ec2 bound_ami_id=ami-fce3c696 policies=prod,dev max_ttl=500h

$ vault write auth/aws/role/dev-role-iam auth_type=iam \
              bound_iam_principal_arn=arn:aws:iam::123456789012:role/MyRole policies=prod,dev max_ttl=500h
```

#### Configure a required X-Vault-AWS-IAM-Server-ID Header (recommended)

```
$ vault write auth/aws/config/client iam_server_id_header_value=vault.example.com
```


#### Perform the login operation

```
$ vault write auth/aws/login role=dev-role \
pkcs7=MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGmewogICJkZXZwYXlQcm9kdWN0Q29kZXMiIDogbnVsbCwKICAicHJpdmF0ZUlwIiA6ICIxNzIuMzEuNjMuNjAiLAogICJhdmFpbGFiaWxpdHlab25lIiA6ICJ1cy1lYXN0LTFjIiwKICAidmVyc2lvbiIgOiAiMjAxMC0wOC0zMSIsCiAgImluc3RhbmNlSWQiIDogImktZGUwZjEzNDQiLAogICJiaWxsaW5nUHJvZHVjdHMiIDogbnVsbCwKICAiaW5zdGFuY2VUeXBlIiA6ICJ0Mi5taWNybyIsCiAgImFjY291bnRJZCIgOiAiMjQxNjU2NjE1ODU5IiwKICAiaW1hZ2VJZCIgOiAiYW1pLWZjZTNjNjk2IiwKICAicGVuZGluZ1RpbWUiIDogIjIwMTYtMDQtMDVUMTY6MjY6NTVaIiwKICAiYXJjaGl0ZWN0dXJlIiA6ICJ4ODZfNjQiLAogICJrZXJuZWxJZCIgOiBudWxsLAogICJyYW1kaXNrSWQiIDogbnVsbCwKICAicmVnaW9uIiA6ICJ1cy1lYXN0LTEiCn0AAAAAAAAxggEXMIIBEwIBATBpMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQwIJAJa6SNnlXhpnMAkGBSsOAwIaBQCgXTAYBgkqhkiG9w0BCQMxCwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0xNjA0MDUxNjI3MDBaMCMGCSqGSIb3DQEJBDEWBBRtiynzMTNfTw1TV/d8NvfgVw+XfTAJBgcqhkjOOAQDBC4wLAIUVfpVcNYoOKzN1c+h1Vsm/c5U0tQCFAK/K72idWrONIqMOVJ8Uen0wYg4AAAAAAAA nonce=5defbf9e-a8f9-3063-bdfc-54b7a42a1f95
```

For the iam auth method, generating the signed request is a non-standard
operation. The Vault cli supports generating this for you:

```
$ vault auth -method=aws header_value=vault.example.com role=dev-role-iam
```

This assumes you have AWS credentials configured in the standard locations AWS
SDKs search for credentials (environment variables, ~/.aws/credentials, IAM
instance profile in that order). If you do not have IAM credentials available at
any of these locations, you can explicitly pass them in on the command line
(though this is not recommended), omitting `aws_security_token` if not
applicable .

```
$ vault auth -method=aws header_value=vault.example.com role=dev-role-iam \
        aws_access_key_id=<access_key> \
        aws_secret_access_key=<secret_key> \
        aws_security_token=<security_token>
```

An example of how to generate the required request values for the `login` method
can be found found in the [vault cli
source code](https://github.com/hashicorp/vault/blob/master/builtin/credential/aws/cli.go).
Using an approach such as this, the request parameters can be generated and
passed to the `login` method:

```
$ vault write auth/aws/login role=dev-role-iam \
        iam_http_request_method=POST \
        iam_request_url=aHR0cHM6Ly9zdHMuYW1hem9uYXdzLmNvbS8= \
        iam_request_body=QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ== \
        iam_request_headers=eyJDb250ZW50LUxlbmd0aCI6IFsiNDMiXSwgIlVzZXItQWdlbnQiOiBbImF3cy1zZGstZ28vMS40LjEyIChnbzEuNy4xOyBsaW51eDsgYW1kNjQpIl0sICJYLVZhdWx0LUFXU0lBTS1TZXJ2ZXItSWQiOiBbInZhdWx0LmV4YW1wbGUuY29tIl0sICJYLUFtei1EYXRlIjogWyIyMDE2MDkzMFQwNDMxMjFaIl0sICJDb250ZW50LVR5cGUiOiBbImFwcGxpY2F0aW9uL3gtd3d3LWZvcm0tdXJsZW5jb2RlZDsgY2hhcnNldD11dGYtOCJdLCAiQXV0aG9yaXphdGlvbiI6IFsiQVdTNC1ITUFDLVNIQTI1NiBDcmVkZW50aWFsPWZvby8yMDE2MDkzMC91cy1lYXN0LTEvc3RzL2F3czRfcmVxdWVzdCwgU2lnbmVkSGVhZGVycz1jb250ZW50LWxlbmd0aDtjb250ZW50LXR5cGU7aG9zdDt4LWFtei1kYXRlO3gtdmF1bHQtc2VydmVyLCBTaWduYXR1cmU9YTY5ZmQ3NTBhMzQ0NWM0ZTU1M2UxYjNlNzlkM2RhOTBlZWY1NDA0N2YxZWI0ZWZlOGZmYmM5YzQyOGMyNjU1YiJdfQ==
```

### Via the API

#### Enable AWS authentication in Vault.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/sys/auth/aws" -d '{"type":"aws"}'
```

#### Configure the credentials required to make AWS API calls.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/auth/aws/config/client" -d '{"access_key":"VKIAJBRHKH6EVTTNXDHA", "secret_key":"vCtSM8ZUEQ3mOFVlYPBQkf2sO6F/W7a5TVzrl3Oj"}'
```

#### Configure the policies on the role.

```
curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/auth/aws/role/dev-role -d '{"bound_ami_id":"ami-fce3c696","policies":"prod,dev","max_ttl":"500h"}'

curl -X POST -H "x-vault-token:123" "http://127.0.0.1:8200/v1/auth/aws/role/dev-role-iam -d '{"auth_type":"iam","policies":"prod,dev","max_ttl":"500h","bound_iam_principal_arn":"arn:aws:iam::123456789012:role/MyRole"}'
```

#### Perform the login operation

```
curl -X POST "http://127.0.0.1:8200/v1/auth/aws/login" -d '{"role":"dev-role","pkcs7":"'$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/pkcs7 | tr -d '\n')'","nonce":"5defbf9e-a8f9-3063-bdfc-54b7a42a1f95"}'

curl -X POST "http://127.0.0.1:8200/v1/auth/aws/login" -d '{"role":"dev", "iam_http_request_method": "POST", "iam_request_url": "aHR0cHM6Ly9zdHMuYW1hem9uYXdzLmNvbS8=", "iam_request_body": "QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ==", "iam_request_headers": "eyJDb250ZW50LUxlbmd0aCI6IFsiNDMiXSwgIlVzZXItQWdlbnQiOiBbImF3cy1zZGstZ28vMS40LjEyIChnbzEuNy4xOyBsaW51eDsgYW1kNjQpIl0sICJYLVZhdWx0LUFXU0lBTS1TZXJ2ZXItSWQiOiBbInZhdWx0LmV4YW1wbGUuY29tIl0sICJYLUFtei1EYXRlIjogWyIyMDE2MDkzMFQwNDMxMjFaIl0sICJDb250ZW50LVR5cGUiOiBbImFwcGxpY2F0aW9uL3gtd3d3LWZvcm0tdXJsZW5jb2RlZDsgY2hhcnNldD11dGYtOCJdLCAiQXV0aG9yaXphdGlvbiI6IFsiQVdTNC1ITUFDLVNIQTI1NiBDcmVkZW50aWFsPWZvby8yMDE2MDkzMC91cy1lYXN0LTEvc3RzL2F3czRfcmVxdWVzdCwgU2lnbmVkSGVhZGVycz1jb250ZW50LWxlbmd0aDtjb250ZW50LXR5cGU7aG9zdDt4LWFtei1kYXRlO3gtdmF1bHQtc2VydmVyLCBTaWduYXR1cmU9YTY5ZmQ3NTBhMzQ0NWM0ZTU1M2UxYjNlNzlkM2RhOTBlZWY1NDA0N2YxZWI0ZWZlOGZmYmM5YzQyOGMyNjU1YiJdfQ==" }'
```

The response will be in JSON. For example:

```javascript
{
  "auth": {
    "renewable": true,
    "lease_duration": 72000,
    "metadata": {
      "role_tag_max_ttl": "0s",
      "role": "ami-f083709d",
      "region": "us-east-1",
      "nonce": "5defbf9e-a8f9-3063-bdfc-54b7a42a1f95",
      "instance_id": "i-a832f734",
      "ami_id": "ami-f083709d"
    },
    "policies": [
      "default",
      "dev",
      "prod"
    ],
    "accessor": "5cd96cd1-58b7-2904-5519-75ddf957ec06",
    "client_token": "150fc858-2402-49c9-56a5-f4b57f2c8ff1"
  },
  "warnings": null,
  "wrap_info": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": "",
  "request_id": "d7d50c06-56b8-37f4-606c-ccdc87a1ee4c"
}
```

## API
### /auth/aws/config/client
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the credentials required to perform API calls to AWS as well as
    custom endpoints to talk to AWS APIs. The instance identity document
    fetched from the PKCS#7 signature will provide the EC2 instance ID. The
    credentials configured using this endpoint will be used to query the status
    of the instances via DescribeInstances API. If static credentials are not
    provided using this endpoint, then the credentials will be retrieved from
    the environment variables `AWS_ACCESS_KEY`, `AWS_SECRET_KEY` and
    `AWS_REGION` respectively. If the credentials are still not found and if the
    backend is configured on an EC2 instance with metadata querying
    capabilities, the credentials are fetched automatically.
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
        <span class="param-flags">optional</span>
        AWS Access key with permissions to query AWS APIs. The permissions
        required depend on the specific configurations. If using the `iam` auth
        method without inferencing, then no credentials are necessary. If using
        the `ec2` auth method or using the `iam` auth method with inferencing,
        then these credentials need access to `ec2:DescribeInstances`. If
        additionally a `bound_iam_role` is specified, then these credentials
        also need access to `iam:GetInstanceProfile`. If, however, an alterate
        sts configuration is set for the target account, then the credentials
        must be permissioned to call `sts:AssumeRole` on the configured role,
        and that role must have the permissions described here.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">secret_key</span>
        <span class="param-flags">optional</span>
        AWS Secret key with permissions to query AWS APIs.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">endpoint</span>
        <span class="param-flags">optional</span>
        URL to override the default generated endpoint for making AWS EC2 API calls.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">iam_endpoint</span>
        <span class="param-flags">optional</span>
        URL to override the default generated endpoint for making AWS IAM API calls.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">sts_endpoint</span>
        <span class="param-flags">optional</span>
        URL to override the default generated endpoint for making AWS STS API calls.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">iam_server_id_header_value</span>
        <span class="param-flags">optional</span>
        The value to require in the `X-Vault-AWS-IAM-Server-ID` header as part of
        GetCallerIdentity requests that are used in the iam auth method. If not
        set, then no value is required or validated. If set, clients must
        include an X-Vault-AWS-IAM-Server-ID header in the headers of login
        requests, and further this header must be among the signed headers
        validated by AWS. This is to protect against different types of replay
        attacks, for example a signed request sent to a dev server being resent
        to a production server. Consider setting this to the Vault server's DNS
        name.
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
    "access_key": "VKIAJBRHKH6EVTTNXDHA"
    "endpoint" "",
    "iam_endpoint" "",
    "sts_endpoint" "",
    "iam_server_id_header_value" "",
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


### /auth/aws/config/certificate/<cert_name>
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Registers an AWS public key to be used to verify the instance identity
    documents. While the PKCS#7 signature of the identity documents have DSA
    digest, the identity signature will have RSA digest, and hence the public
    keys for each type varies respectively. Indicate the type of the public key
    using the "type" parameter.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/certificate/<cert_name>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">cert_name</span>
        <span class="param-flags">required</span>
        Name of the certificate.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">aws_public_cert</span>
        <span class="param-flags">required</span>
        AWS Public key required to verify PKCS7 signature of the EC2 instance metadata.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">type</span>
        <span class="param-flags">optional</span>
        Takes the value of either "pkcs7" or "identity", indicating the type of
        document which can be verified using the given certificate. The PKCS#7
        document will have a DSA digest and the identity signature will have an
        RSA signature, and accordingly the public certificates to verify those
        also vary. Defaults to "pkcs7".
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
  <dd>`/auth/aws/config/certificate/<cert_name>`</dd>

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

#### LIST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Lists all the AWS public certificates that are registered with the backend.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/certificates` (LIST) or `/auth/aws/config/certificates?list=true` (GET)</dd>

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
      "cert1"
    ]
  },
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}
```

  </dd>
</dl>


### /auth/aws/config/sts/<account_id>
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Allows the explicit association of STS roles to satellite AWS accounts
    (i.e. those which are not the account in which the Vault server is
    running.) Login attempts from EC2 instances running in these accounts will
    be verified using credentials obtained by assumption of these STS roles.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/sts/<account_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">account_id</span>
        <span class="param-flags">required</span>
        AWS account ID to be associated with STS role. If set, Vault will use
        assumed credentials to verify any login attempts from EC2 instances in
        this account.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">sts_role</span>
        <span class="param-flags">required</span>
        AWS ARN for STS role to be assumed when interacting with the account
        specified.  The Vault server must have permissions to assume this role.
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
    Returns the previously configured STS role. 
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/sts/<account_id>`</dd>

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
    "sts_role ": "arn:aws:iam:<account_id>:role/myRole"
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
    Lists all the AWS Account IDs for which an STS role is registered 
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/sts` (LIST) or `/auth/aws/config/sts?list=true` (GET)</dd>

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
      "<account_id_1>",
      "<account_id_2>"
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
    Deletes a previously configured AWS account/STS role association  
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/sts/<account_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>

### /auth/aws/config/tidy/identity-whitelist
##### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the periodic tidying operation of the whitelisted identity entries.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/tidy/identity-whitelist`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">safety_buffer</span>
        <span class="param-flags">optional</span>
        The amount of extra time that must have passed beyond the `roletag`
        expiration, before it is removed from the backend storage. Defaults to
        72h.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">disable_periodic_tidy</span>
        <span class="param-flags">optional</span>
        If set to 'true', disables the periodic tidying of the
        'identity-whitelist/<instance_id>' entries.
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
    Returns the previously configured periodic whitelist tidying settings.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/tidy/identity-whitelist`</dd>

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
    "safety_buffer": 60,
    "disable_periodic_tidy": false
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
    Deletes the previously configured periodic whitelist tidying settings.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/tidy/identity-whitelist`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>



### /auth/aws/config/tidy/roletag-blacklist
##### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Configures the periodic tidying operation of the blacklisted role tag entries.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/tidy/roletag-blacklist`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">safety_buffer</span>
        <span class="param-flags">optional</span>
        The amount of extra time that must have passed beyond the `roletag`
        expiration, before it is removed from the backend storage. Defaults to
        72h.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">disable_periodic_tidy</span>
        <span class="param-flags">optional</span>
        If set to 'true', disables the periodic tidying of the
        'roletag-blacklist/<role_tag>' entries.
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
    Returns the previously configured periodic blacklist tidying settings.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/tidy/roletag-blacklist`</dd>

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
    "safety_buffer": 60,
    "disable_periodic_tidy": false
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
    Deletes the previously configured periodic blacklist tidying settings.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/config/tidy/roletag-blacklist`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>



### /auth/aws/role/[role]
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Registers a role in the backend. Only those instances or principals which
    are using the role registered using this endpoint, will be able to perform
    the login operation. Contraints can be specified on the role, that are
    applied on the instances or principals attempting to login. At least one
    constraint should be specified on the role. The available constraints you
    can choose are dependent on the `auth_type` of the role and, if the
    `auth_type` is `iam`, then whether inferencing is enabled. A role will not
    let you configure a constraint if it is not checked by the `auth_type` and
    inferencing configuration of that role.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/role/<role>`</dd>

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
        <span class="param">auth_type</span>
        <span class="param-flags">optional</span>
        The auth type permitted for this role. Valid choices are "ec2" or "iam".
        If no value is specified, then it will default to "iam" (except for
        legacy `aws-ec2` auth types, for which it will default to "ec2"). Only
        those bindings applicable to the auth type chosen will be allowed to be
        configured on the role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_ami_id</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the EC2 instances that they should be
        using the AMI ID specified by this parameter. This constraint is checked
        during ec2 auth as well as the iam auth method only when inferring an
        EC2 instance.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_account_id</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the EC2 instances that the account ID in
        its identity document to match the one specified by this parameter. This
        constraint is checked during ec2 auth as well as the iam auth method
        only when inferring an EC2 instance.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_region</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the EC2 instances that the region in
        its identity document must match the one specified by this parameter. This
        constraint is only checked by the ec2 auth method as well as the iam
        auth method only when inferring an ec2 instance..
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_vpc_id</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the EC2 instance to be associated with
        the VPC ID that matches the value specified by this parameter. This
        constraint is only checked by the ec2 auth method as well as the iam
        auth method only when inferring an ec2 instance.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_subnet_id</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the EC2 instance to be associated with
        the subnet ID that matches the value specified by this parameter. This
        constraint is only checked by the ec2 auth method as well as the iam
        auth method only when inferring an ec2 instance.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_iam_role_arn</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the authenticating EC2 instance that it must
        match the IAM role ARN specified by this parameter.  The value is
        prefix-matched (as though it were a glob ending in `*`).  The configured IAM
        user or EC2 instance role must be allowed to execute the
        `iam:GetInstanceProfile` action if this is specified. This constraint is
        checked by the ec2 auth method as well as the iam auth method only when
        inferring an EC2 instance.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_iam_instance_profile_arn</span>
        <span class="param-flags">optional</span>
        If set, defines a constraint on the EC2 instances to be associated with
        an IAM instance profile ARN which has a prefix that matches the value
        specified by this parameter. The value is prefix-matched (as though it
        were a glob ending in `*`). This constraint is checked by the ec2 auth
        method as well as the iam auth method only when inferring an ec2
        instance.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">role_tag</span>
        <span class="param-flags">optional</span>
        If set, enables the role tags for this role. The value set for this
        field should be the 'key' of the tag on the EC2 instance. The 'value'
        of the tag should be generated using `role/<role>/tag` endpoint.
        Defaults to an empty string, meaning that role tags are disabled. This
        constraint is valid only with the ec2 auth method and is not allowed
        when an auth_type is iam.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">bound_iam_principal_arn</span>
        <span class="param-flags">optional</span>
        Defines the IAM principal that must be authenticated using the iam
        auth method. It should look like
        "arn:aws:iam::123456789012:user/MyUserName" or
        "arn:aws:iam::123456789012:role/MyRoleName". This constraint is only
        checked by the iam auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">inferred_entity_type</span>
        <span class="param-flags">optional</span>
        When set, instructs Vault to turn on inferencing. The only current valid
        value is "ec2_instance" instructing Vault to infer that the role comes
        from an EC2 instance in an IAM instance profile. This only applies to
        the iam auth method. If you set this on an existing role where it had
        not previously been set, tokens that had been created prior will not be
        renewable; clients will need to get a new token.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">inferred_aws_region</span>
        <span class="param-flags">optional</span>
        When role inferencing is activated, the region to search for the
        inferred entities (e.g., EC2 instances). Required if role inferencing is
        activated. This only applies to the iam auth method.
      </li>
    </ul>
    <ul>
    <li>
        <span class="param">resolve_aws_unique_ids</span>
        <span class="param-flags">optional</span>
        When set, resolves the `bound_iam_principal_arn` to the [AWS Unique
        ID](http://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-unique-ids).
        This requires Vault to be able to call `iam:GetUser` or `iam:GetRole` on
        the `bound_iam_principal_arn` that is being bound. Resolving to
        internal AWS IDs more closely mimics the behavior of AWS services in
        that if an IAM user or role is deleted and a new one is recreated with
        the same name, those new users or roles won't get access to roles in
        Vault that were permissioned to the prior principals of the same name.
        The default value for new roles is true, while the default value for
        roles that existed prior to this option existing is false (you can
        check the value for a given role using the GET method on the role). Any
        authentication tokens created prior to this being supported won't
        verify the unique ID upon token renewal.  When this is changed from
        false to true on an existing role, Vault will attempt to resolve the
        role's bound IAM ARN to the unique ID and, if unable to do so, will
        fail to enable this option.  Changing this from `true` to `false` is
        not supported; if absolutely necessary, you would need to delete the
        role and recreate it explicitly setting it to `false`. However; the
        instances in which you would want to do this should be rare. If the
        role creation (or upgrading to use this) succeed, then Vault has
        already been able to resolve internal IDs, and it doesn't need any
        further IAM permissions to authenticate users. If a role has been
        deleted and recreated, and Vault has cached the old unique ID, you
        should just call this endpoint specifying the same
        `bound_iam_principal_arn` and, as long as Vault still has the necessary
        IAM permissions to resolve the unique ID, Vault will update the unique
        ID. (If it does not have the necessary permissions to resolve the
        unique ID, then it will fail to update.) If this option is set to
        false, then you MUST leave out the path component in
        bound_iam_principal_arn for **roles** only, but not IAM users. That is,
        if your IAM role ARN is of the form
        `arn:aws:iam::123456789012:role/some/path/to/MyRoleName`, you **must**
        specify a bound_iam_principal_arn of
        `arn:aws:iam::123456789012:role/MyRoleName` for authentication to
        work.
      </li>
    </ul>

    <ul>
      <li>
        <span class="param">ttl</span>
        <span class="param-flags">optional</span>
        The TTL period of tokens issued using this role, provided as "1h",
        where hour is the largest suffix.
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
        <span class="param">period</span>
        <span class="param-flags">optional</span>
        If set, indicates that the token generated using this role should never
        expire. The token should be renewed within the duration specified by
        this value. At each renewal, the token's TTL will be set to the value
        of this parameter.  The maximum allowed lifetime of tokens issued using
        this role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Policies to be set on tokens issued using this role.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">allow_instance_migration</span>
        <span class="param-flags">optional</span>
        If set, allows migration of the underlying instance where the client
        resides. This keys off of pendingTime in the metadata document, so
        essentially, this disables the client nonce check whenever the instance
        is migrated to a new host and pendingTime is newer than the
        previously-remembered time. Use with caution. This only applies to
        authentications via the ec2 auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">disallow_reauthentication</span>
        <span class="param-flags">optional</span>
        If set, only allows a single token to be granted per instance ID. In
        order to perform a fresh login, the entry in whitelist for the instance
        ID needs to be cleared using
        'auth/aws/identity-whitelist/<instance_id>' endpoint. Defaults to
        'false'. This only applies to authentications via the ec2 auth method.
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
  <dd>`/auth/aws/role/<role>`</dd>

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
    "bound_ami_id": "ami-fce36987",
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
    Lists all the roles that are registered with the backend.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/roles` (LIST) or `/auth/aws/roles?list=true` (GET)</dd>

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
  <dd>`/auth/aws/role/<role>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/role/[role]/tag
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
     Creates a role tag on the role, which help in restricting the capabilities
     that are set on the role. Role tags are not tied to any specific ec2
     instance unless specified explicitly using the `instance_id` parameter. By
     default, role tags are designed to be used across all instances that
     satisfies the constraints on the role. Regardless of which instances have
     role tags on them, capabilities defined in a role tag must be a strict
     subset of the given role's capabilities.  Note that, since adding and
     removing a tag is often a widely distributed privilege, care needs to be
     taken to ensure that the instances are attached with correct tags to not
     let them gain more privileges than what were intended.  If a role tag is
     changed, the capabilities inherited by the instance will be those defined
     on the new role tag. Since those must be a subset of the role
     capabilities, the role should never provide more capabilities than any
     given instance can be allowed to gain in a worst-case scenario.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/role/<role>/tag`</dd>

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
        <span class="param">policies</span>
        <span class="param-flags">optional</span>
        Policies to be associated with the tag. If set, must be a subset of the
        role's policies. If set, but set to an empty value, only the 'default'
        policy will be given to issued tokens.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">max_ttl</span>
        <span class="param-flags">optional</span>
        If set, specifies the maximum allowed token lifetime.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">instance_id</span>
        <span class="param-flags">optional</span>
        Instance ID for which this tag is intended for. If set, the created tag
        can only be used by the instance with the given ID.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">disallow_reauthentication</span>
        <span class="param-flags">optional</span>
        If set, only allows a single token to be granted per instance ID. This
        can be cleared with the auth/aws/identity-whitelist endpoint.
        Defaults to 'false'.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">allow_instance_migration</span>
        <span class="param-flags">optional</span>
        If set, allows migration of the underlying instance where the client
        resides. This keys off of pendingTime in the metadata document, so
        essentially, this disables the client nonce check whenever the instance
        is migrated to a new host and pendingTime is newer than the
        previously-remembered time. Use with caution. Defaults to 'false'.
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
    "tag_value": "v1:09Vp0qGuyB8=:r=dev-role:p=default,prod:d=false:t=300h0m0s:uPLKCQxqsefRhrp1qmVa1wsQVUXXJG8UZP/pJIdVyOI=",
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
    Fetch a token. This endpoint verifies the pkcs7 signature of the instance
    identity document or the signature of the signed GetCallerIdentity request.
    With the ec2 auth method, or when inferring an EC2 instance, verifies that
    the instance is actually in a running state.  Cross checks the constraints
    defined on the role with which the login is being performed. With the ec2
    auth method, as an alternative to pkcs7 signature, the identity document
    along with its RSA digest can be supplied to this endpoint.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/login`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role</span>
        <span class="param-flags">optional</span>
        Name of the role against which the login is being attempted.
        If `role` is not specified, then the login endpoint looks for a role
        bearing the name of the AMI ID of the EC2 instance that is trying to
        login if using the ec2 auth method, or the "friendly name" (i.e., role
        name or username) of the IAM principal authenticated.
        If a matching role is not found, login fails.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">identity</span>
        <span class="param-flags">required</span>
        Base64 encoded EC2 instance identity document. This needs to be
        supplied along with the `signature` parameter. If using `curl` for
        fetching the identity document, consider using the option `-w 0` while
        piping the output to `base64` binary.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">signature</span>
        <span class="param-flags">required</span>
        Base64 encoded SHA256 RSA signature of the instance identity document.
        This needs to be supplied along with `identity` parameter when using the
        ec2 auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">pkcs7</span>
        <span class="param-flags">required</span>
        PKCS7 signature of the identity document with all `\n` characters
        removed.  Either this needs to be set *OR* both `identity` and
        `signature` need to be set when using the ec2 auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">nonce</span>
        <span class="param-flags">optional</span>
        The nonce to be used for subsequent login requests. If this parameter
        is not specified at all and if reauthentication is allowed, then the
        backend will generate a random nonce, attaches it to the instance's
        identity-whitelist entry and returns the nonce back as part of auth
        metadata. This value should be used with further login requests, to
        establish client authenticity. Clients can choose to set a custom nonce
        if preferred, in which case, it is recommended that clients provide a
        strong nonce.  If a nonce is provided but with an empty value, it
        indicates intent to disable reauthentication. Note that, when
        `disallow_reauthentication` option is enabled on either the role or the
        role tag, the `nonce` holds no significance. This is ignored unless
        using the ec2 auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">iam_http_request_method</span>
        <span class="param-flags">required</span>
        HTTP method used in the signed request. Currently only POST is
        supported, but other methods may be supported in the future. This is
        required when using the iam auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">iam_request_url</span>
        <span class="param-flags">required</span>
        Base64-encoded HTTP URL used in the signed request. Most likely just
        `aHR0cHM6Ly9zdHMuYW1hem9uYXdzLmNvbS8=` (base64-encoding of
        `https://sts.amazonaws.com/`) as most requests will probably use POST
        with an empty URI. This is required when using the iam auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">iam_request_body</span>
        <span class="param-flags">required</span>
        Base64-encoded body of the signed request. Most likely
        `QWN0aW9uPUdldENhbGxlcklkZW50aXR5JlZlcnNpb249MjAxMS0wNi0xNQ==`
        which is the base64 encoding of
        `Action=GetCallerIdentity&Version=2011-06-15`. This is required
        when using the iam auth method.
      </li>
    </ul>
    <ul>
      <li>
        <span class="param">iam_request_headers</span>
        <span class="param-flags">required</span>
        Base64-encoded, JSON-serialized representation of the HTTP request
        headers. The JSON serialization assumes that each header key maps to an
        array of string values (though the length of that array will probably
        only be one). If the `iam_server_id_header_value` is configured in Vault
        for the aws auth mount, then the headers must include the
        X-Vault-AWS-IAM-Server-ID header, its value must match the value
        configured, and the header must be included in the signed headers.  This
        is required when using the iam auth method.
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
      "ami_id": "ami-fce36983"
      "role": "dev-role",
      "auth_type": "ec2"
    },
    "policies": [
      "default",
      "dev",
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


### /auth/aws/roletag-blacklist/<role_tag>
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Places a valid role tag in a blacklist. This ensures that the role tag
    cannot be used by any instance to perform a login operation again.  Note
    that if the role tag was previously used to perform a successful login,
    placing the tag in the blacklist does not invalidate the already issued
    token.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/roletag-blacklist/<role_tag>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">role_tag</span>
        <span class="param-flags">required</span>
        Role tag to be blacklisted. The tag can be supplied as-is. In order to
        avoid any encoding problems, it can be base64 encoded.
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
    Returns the blacklist entry of a previously blacklisted role tag.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/broletag-blacklist/<role_tag>`</dd>

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
    Lists all the role tags that are blacklisted.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/roletag-blacklist` (LIST) or `/auth/aws/roletag-blacklist?list=true` (GET)</dd>

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
    Deletes a blacklisted role tag.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/roletag-blacklist/<role_tag>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/tidy/roletag-blacklist
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Cleans up the entries in the blacklist based on expiration time on the
    entry and `safety_buffer`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/tidy/roletag-blacklist`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">safety_buffer</span>
        <span class="param-flags">optional</span>
        The amount of extra time that must have passed beyond the `roletag`
        expiration, before it is removed from the backend storage. Defaults to
        72h.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/identity-whitelist/<instance_id>
#### GET
<dl class="api">
  <dt>Description</dt>
  <dd>
    Returns an entry in the whitelist. An entry will be created/updated by
    every successful login.
  </dd>

  <dt>Method</dt>
  <dd>GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/identity-whitelist/<instance_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">instance_id</span>
        <span class="param-flags">required</span>
        EC2 instance ID. A successful login operation from an EC2 instance gets
        cached in this whitelist, keyed off of instance ID.
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
    "pending_time": "2016-04-14T01:01:41Z",
    "expiration_time": "2016-05-05 10:09:16.67077232 +0000 UTC",
    "creation_time": "2016-04-14 14:09:16.67077232 +0000 UTC",
    "client_nonce": "5defbf9e-a8f9-3063-bdfc-54b7a42a1f95",
    "role": "dev-role"
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
    Lists all the instance IDs that are in the whitelist of successful logins.
  </dd>

  <dt>Method</dt>
  <dd>LIST/GET</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/identity-whitelist` (LIST) or `/auth/aws/identity-whitelist?list=true` (GET)</dd>
  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.

```javascript
{
  "auth": null,
  "warnings": null,
  "data": {
    "keys": [
      "i-aab47d37"
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
    Deletes a cache of the successful login from an instance.
  </dd>

  <dt>Method</dt>
  <dd>DELETE</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/identity-whitelist/<instance_id>`</dd>

  <dt>Parameters</dt>
  <dd>
    None.
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>


### /auth/aws/tidy/identity-whitelist
#### POST
<dl class="api">
  <dt>Description</dt>
  <dd>
    Cleans up the entries in the whitelist based on expiration time and `safety_buffer`.
  </dd>

  <dt>Method</dt>
  <dd>POST</dd>

  <dt>URL</dt>
  <dd>`/auth/aws/tidy/identity-whitelist`</dd>

  <dt>Parameters</dt>
  <dd>
    <ul>
      <li>
        <span class="param">safety_buffer</span>
        <span class="param-flags">optional</span>
        The amount of extra time that must have passed beyond the identity
        expiration, before it is removed from the backend storage. Defaults to
        72h.
      </li>
    </ul>
  </dd>

  <dt>Returns</dt>
  <dd>`204` response code.
  </dd>
</dl>
