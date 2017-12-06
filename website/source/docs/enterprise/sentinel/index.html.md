---
layout: "docs"
page_title: "Vault Enterprise Sentinel Integration"
sidebar_current: "docs-vault-enterprise-sentinel"
description: |-
  An overview of how Sentinel interacts with Vault Enterprise.

---

# Overview

Vault Enterprise integrates HashiCorp Sentinel to provide a rich set of access
control functionality. Because Vault is a security-focused product trusted with
high-risk secrets and assets, and because of its default-deny stance,
integration with Vault is implemented in a defense-in-depth fashion. This takes
the form of multiple types of policies and a fixed evaluation order.

## Policy Types

Vault's policy system has been expanded to support three types of policies:

- `ACLs` - These are the [traditional Vault
  policies](/docs/concepts/policies.html) and remain unchanged.

- `Role Governing Policies (RGPs)` - RGPs are Sentinel policies that are tied
  to particular tokens, Identity entities, or Identity groups. They have access
  to a rich set of controls across various aspects of Vault.

- `Endpoint Governing Policies (EGPs)` - EGPs are Sentinel policies that are
  tied to particular paths instead of tokens. They have access to as much
  request information as possible, but they can take effect even on
  unauthenticated paths, such as login paths.

Not every unauthenticated path supports EGPs. For instance, the paths related
to root token generation cannot support EGPs because it's already the mechanism
of last resort if, for instance, all clients are locked out of Vault due to
misconfigured EGPs.

Like with ACLs, [root tokens](/docs/concepts/tokens.html#root-tokens) tokens
are not subject to Sentinel policy checks.

Sentinel execution should be considered to be significantly slower than normal
ACL policy checking. If high performance is needed, testing should be performed
appropriately when introducing Sentinel policies.

## Policy Evaluation

During evaluation, all policy types, if they exist, must grant access.
Evaluation uses the following logic:

1. If the request is unauthenticated, skip to step 3. Otherwise, evaluate the
   token's ACL policies. These must grant access; as always, a failure to be
   granted capabilities on a path via ACL policies denies the request.
2. RGPs attached to the token are evaluated. All policies must pass according
   to their enforcement level.
3. EGPs set on the requested path, and any prefix-matching EGPs set on
   less-specific paths, are evaluated. All policies must pass according to
   their enforcement level.

Any failure at any of these steps results in a denied request.

## Policy Overriding

Vault supports normal Sentinel overriding behavior. Requests to override can be
specified on the command line via the `policy-override` flag or in HTTP
requests by setting the `X-Vault-Policy-Override` header to `true`.

Override requests are visible in Vault's audit log; in addition, override
requests and their eventual status (whether they ended up being required) are
logged as warnings in Vault's server logs.

## MFA

Sentinel policies support the [Identity-based MFA
system](/docs/enterprise/mfa/index.html) in Vault Enterprise.  Within a single
request, multiple checks of any named MFA method will only trigger
authentication behavior for that method once, regardless of whether its
validity is checked via ACLs, RGPs, or EGPs.

EGPs can be used to require MFA on otherwise unauthenticated paths, such as
login paths. On such paths, the request data will perform a lookahead to try to
discover the appropriate Identity information to use for MFA. It may be
necessary to pre-populate Identity entries or supply additional parameters with
the request if you require more information to use MFA than the endpoint is
able to glean from the original request alone.

# Using Sentinel

## Configuration

Sentinel policies can be configured via the `sys/policies/rgp/` and
`sys/policies/egp/` endpoints; see [the
documentation](/api/system/policies.html) for more information.

Once set, RGPs can be assigned to Identity entities and groups or to tokens
just like ACL policies. As a result, they cannot share names with ACL policies.

When setting an EGP, a list of paths must be provided specifying on which paths
that EGP should take effect. Endpoints can have multiple distinct EGPs set on
them; all are evaluated for each request. Paths can use a glob character (`*`)
as the last character of the path to perform a prefix match; a path that
consists only of a `*` will apply to the root of the API. Since requests are
subject to an EGPs exactly matching the requested path and any glob EGPs
sitting further up the request path, an EGP with a path of `*` will thus take
effect on all requests.

## Properties and Examples

See the [Examples](/docs/enterprise/sentinel/examples.html) page for examples
of Sentinel in action, and the
[Properties](/docs/enterprise/sentinel/properties.html) page for detailed
property documentation.
