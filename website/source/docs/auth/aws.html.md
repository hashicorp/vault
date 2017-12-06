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

Amazon EC2 instances have access to metadata which describes the instance. The
Vault EC2 authentication method leverages this components of this metadata to
authenticate and distribute an initial Vault token to an EC2 instance. The data
flow (which is also represented in the graphic below) is as follows:

[![Vault AWS EC2 Authentication Flow](/assets/images/vault-aws-ec2-auth-flow.png)](/assets/images/vault-aws-ec2-auth-flow.png)

1. An AWS EC2 instance fetches its [AWS Instance Identity Document][aws-iid]
from the [EC2 Metadata Service][aws-ec2-mds]. In addition to data itself, AWS
also provides the PKCS#7 signature of the data, and publishes the public keys
(by region) which can be used to verify the signature.

1. The AWS EC2 instance makes a request to Vault with the Instance Identity
Document and the PKCS#7 signature of the document.

1. Vault verifies the signature on the PKCS#7 document, ensuring the information
is certified accurate by AWS. This process validates both the validity and
integrity of the document data. As an added security measure, Vault verifies
that the instance is currently running using the public EC2 API endpoint.

1. Provided all steps are successful, Vault returns the initial Vault token to
the EC2 instance. This token is mapped to any configured policies based on the
instance metadata.

[aws-iid]: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html
[aws-ec2-mds]: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html

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
`X-Amz-Signature`, and `X-Amz-SignedHeaders` GET query parameters containing the
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

The iam authentication method allows you to specify a bound IAM principal ARN.
Clients authenticating to Vault must have an ARN that matches the ARN bound to
the role they are attempting to login to. The bound ARN allows specifying a
wildcard at the end of the bound ARN. For example, if the bound ARN were
`arn:aws:iam::123456789012:*` it would allow any principal in AWS account
123456789012 to login to it. Similarly, if it were
`arn:aws:iam::123456789012:role/*` it would allow any IAM role in the AWS
account to login to it. If you wish to specify a wildcard, you must give Vault
`iam:GetUser` and `iam:GetRole` permissions to properly resolve the full user
path.

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

## Recommended Vault IAM Policy

This specifies the recommended IAM policy needed by the AWS auth backend. Note
that if you are using the same credentials for the AWS auth and secret backends
(e.g., if you're running Vault on an EC2 instance in an IAM instance profile),
then you will need to add additional permissions as required by the AWS secret
backend.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "iam:GetInstanceProfile",
        "iam:GetUser",
        "iam:GetRole"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["sts:AssumeRole"],
      "Resource": [
        "arn:aws:iam:<AccountId>:role/<VaultRole>"
      ]
    }
  ]
}
```

Here are some of the scenarios in which Vault would need to use each of these
permissions. This isn't intended to be an exhaustive list of all the scenarios
in which Vault might make an AWS API call, but rather illustrative of why these
are needed.

* `ec2:DescribeInstances` is necessary when you are using the `ec2` auth method
  or when you are inferring an `ec2_instance` entity type to validate that the
  EC2 instance meets binding requirements of the role
* `iam:GetInstanceProfile` is used when you have a `bound_iam_role_arn` in the
  `ec2` auth method. Vault needs to determine which IAM role is attached to the
  instance profile.
* `iam:GetUser` and `iam:GetRole` are used when using the iam auth method and
  binding to an IAM user or role principal to determine the unique AWS user ID
  or when using a wildcard on the bound ARN to resolve the full ARN of the user
  or role.
* The `sts:AssumeRole` stanza is necessary when you are using [Cross Account
  Access](#cross-account-access). The `Resource`s specified should be a list of
  all the roles for which you have configured cross-account access, and each of
  those roles should have this IAM policy attached (except for the
  `sts:AssumeRole` statement).

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

The client nonce which is generated by the backend and which gets returned
along with the authentication response, will be audit logged in plaintext. If
this is undesired, clients can supply a custom nonce to the login endpoint
which will not be returned and hence will not be audit logged.

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
instance profile, or ECS task role, in that order). If you do not have IAM
credentials available at any of these locations, you can explicitly pass them
in on the command line (though this is not recommended), omitting
`aws_security_token` if not applicable.

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

The AWS authentication backend has a full HTTP API. Please see the
[AWS Auth API](/api/auth/aws/index.html) for more
details.
