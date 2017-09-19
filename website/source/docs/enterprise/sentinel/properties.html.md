---
layout: "docs"
page_title: "Sentinel Properties"
sidebar_current: "docs-vault-enterprise-sentinel-properties"
description: |-
  An overview of how Sentinel interacts with Vault Enterprise.

---

# Overview

Vault injects a rich set of data into the running Sentinel environment,
allowing for very fine-grained controls. The set of available properties are
enumerated on this page.

# Examples

Following are some examples that help to introduce concepts. If you are
unfamiliar with writing Sentinel policies in Vault, please read through to
understand some best practices.

## MFA and CIDR Check on Login

The following Sentinel policy requires the incoming user to successfully
validate with an Okta MFA push request before authenticating with LDAP.
Additionally, it ensures that only users on the 10.20.0.0/16 subnet are able to
authenticate using LDAP.

```python
import "sockaddr"

# We expect logins to come only from our private IP range
cidrcheck = rule {
    sockaddr.is_contained(request.connection.remote_addr, "10.20.0.0/16")
}

# Require Ping MFA validation to succeed
ping_valid = rule {
    mfa.methods.ping.valid
}

main = rule when request.path is "auth/ldap/login" {
    ping_valid and cidrcheck
}
```

Note the `rule when` construct on the `main` rule. This scopes the policy to
the given condition.

Vault takes a default-deny approach to security. Without such scoping, because
active Sentinel policies must all pass successfully, the user would be forced
to start with a passing status and then define the conditions under which
access is denied, breaking the default-deny concept.

By instead indicating the conditions under which the `main` rule (and thus, in
this example, the entire policy) should be evaluated, the policy instead
describes the conditions under which a matching request is successful. This
keeps the default-deny feeling of Vault; if the evaluation condition isn't met,
the policy is simply a no-op.

## Allow Only Specific Identity Entities or Groups

```python
main = rule {
    identity.entity.name is "jeff" or
    identity.entity.id is "fe2a5bfd-c483-9263-b0d4-f9d345efdf9f" or
    "sysops" in identity.groups.names or
    "14c0940a-5c07-4b97-81ec-0d423accb8e0" in keys(identity.groups.by_id)
}
```

This example shows accessing Identity properties to make decisions, showing
that for Identity values IDs or names can be used for reference.

In general, it is more secure to use IDs. While convenient, entity names and
group names can be switched from one entity to another, because their only
constraint is that they must be unique. Using IDs guarantees that only that
specific entity or group is sufficient; if the group or entity are deleted and
recreated with the same name, the match will fail.

## Instantly Disallow All Previously-Generated Tokens

Imagine a break-glass scenario where it is discovered that there have been
compromises of some unknown number of previously-generated tokens.

In such a situation it would be possible to revoke all previous tokens, but
this may take a while for a number of reasons, from requiring revocation of
generated secrets to the simple delay required to remove many entries from
storage. In addition, it could revoke tokens and generated secrets that later
forensic analysis shows were not compromised, unnecessarily widening the impact
of the mass revocation.

In Vault's ACL system a simple deny could be put into place, but this is a very
coarse-grained control and would require forethought to ensure that a policy
that can be modified in such a way is attached to every token. It also would
not prevent access to login paths or other unauthenticated paths.

Sentinel offers much more fine-grained control:

```python
import "time"

main = rule when not request.unauthenticated {
    time.load(token.creation_time).unix >
      time.load("2017-09-17T13:25:29Z").unix
}
```

Created as an EGP on `*`, this will block all access to any path Sentinel
operates on with a token created before the given time. Tokens created after
this time, since they were not a part of the compromise, will not be subject to
this restriction.

## Delegate EGP Policy Management Under a Path

The following policy gives token holders with this policy (via their tokens or
their Identity entities/groups) the ability to write EGP policies that can only
take effect at Vault paths below certain prefixes. This effectively delegates
policy management to the team for their own key-value spaces.

```python
import "strings"

data_match = func() {
    # Make sure there is request data
    if length(request.data else 0) is 0 {
        return false
    }

    # Make sure request data includes paths
    if length(request.data.paths else 0) is 0 {
        return false
    }

    # For each path, verify that it is in the allowed list
    for strings.split(request.data.paths, ",") as path {
        # Make it easier for users who might be used to starting paths with
        # slashes
        sanitizedPath = strings.trim_prefix(path, "/")
        if not strings.has_prefix(sanitizedPath, "dev-kv/teama/") and
           not strings.has_prefix(sanitizedPath, "prod-kv/teama/") {
            return false
        }
    }

    return true
}

# Only care about writing; reading can be allowed by normal ACLs
precond = rule {
    request.operation in ["create", "update"] and
    strings.has_prefix(request.path, "sys/policies/egp/")
}

main = rule when precond {
    strings.has_prefix(request.path, "sys/policies/egp/teama-") and data_match()
}
```

# Properties

The following properties are available for use in Sentinel policies.

## Request Properties

The following properties are available in the `request` namespace.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `connection.remote_addr` | `string` | TCP/IP source address/port of the client |
| `data` | `map (string -> any)` | Raw request data |
| `operation` | `string` | Operation type, e.g. "read" or "update" |
| `path` | `string` | Path, with any leading `/` trimmed |
| `policy_override` | `bool` | `true` if a `soft-mandatory` policy override was requested |
| `unauthenticated` | `bool` | `true` if the requested path is an unauthenticated path |
| `wrapping.ttl` | `duration` | The requested response-wrapping TTL in nanoseconds, suitable for use with the `time` import|
| `wrapping.ttl_seconds` | `int` | The requested response-wrapping TTL in seconds |

## Token Properties

The following properties, if available, are in the `token` namespace. The
namespace will not exist if there is no token information attached to a
request, e.g. when logging in.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `creation_time` | `string` | The timestamp of the token's creation, in RFC3339 format |
| `creation_time_unix` | `int` | The timestamp of the token's creation, in seconds since Unix epoch UTC |
| `creation_ttl` | `duration` | The TTL the token was first created with in nanoseconds, suitable for use with the `time` import |
| `creation_ttl_seconds` | `int` | The TTL the token was first created with in seconds |
| `display_name` | `string` | The display name set on the token, if any |
| `entity_id` | `string` | The Identity entity ID attached to the token, if any |
| `explicit_max_ttl` | `duration` | If the token has an explicit max TTL, the duration of the explicit max TTL in nanoseconds, suitable for use with the `time` import |
| `explicit_max_ttl_seconds` | `int` | If the token has an explicit max TTL, the duration of the explicit max TTL in seconds |
| `metadata` | `map (string -> string)` | Metadata set on the token |
| `num_uses` | `int` | The number of uses remaining on a use-count-limited token; 0 if the token has no use-count limit |
| `path` | `string` | The request path that resulted in creation of this token |
| `period` | `duration` | If the token has a period, the duration of the period in nanoseconds, suitable for use with the `time` import |
| `period_seconds` | `int` | If the token has a period, the duration of the period in seconds |
| `policies` | `list (string)` | Policies directly attached to the token |
| `role` | `string` | If created via a token role, the role that created the token |

## Identity Properties

The following properties, if available, are in the `identity` namespace. The
namespace may not exist if there is no token information attached to the
request; however, at login time the user's request data will be used to attempt
to find any existing Identity information, or create some information to pass
to MFA functions.

### Entity Properties

These exist at the `identity.entity` namespace.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `creation_time` | `string` | The entity's creation time in RFC3339 format |
| `id` | `string` | The entity's ID |
| `last_update_time` | `string` | The entity's last update (modify) time in RFC3339 format |
| `metadata` | `map (string -> string)` | Metadata associated with the entity |
| `name` | `string` | The entity's name |
| `merged_entity_ids` | `list (string)` | A list of IDs of entities that have been merged into this one |
| `personas` | `list (persona)` | List of personas associated with this entity |
| `policies` | `list (string)` | List of the policies set on this entity |

### Persona Properties

These can be retrieved from `identity.entity.personas`.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `creation_time` | `string` | The persona's creation time in RFC3339 format |
| `id` | `string` | The persona's ID |
| `last_update_time` | `string` | The persona's last update (modify) time in RFC3339 format |
| `metadata` | `map (string -> string)` | Metadata associated with the persona|
| `merged_from_entity_ids` | `list (string)` | If this persona was attached to the current entity via one or more merges, the original entity/entities will be in this list |
| `mount_accessor` | `string` | The immutable accessor of the mount that created this persona |
| `mount_path` | `string` | The path of the mount that created this persona; unlike the accessor, there is no guarantee that the current path represents the original mount |
| `mount_type` | `string` | The type of the mount that created this persona |
| `name` | `string` | The persona's name |

### Groups Properties

These exist at the `identity.groups` namespace.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `by_id` | `map (string -> group)` | A map of group ID to group information |
| `by_name` | `map (string -> group)` | A map of group name to group information; unlike the group ID, there is no guarantee that the current name will always represent the same group |

### Group Properties

These can be retrieved from the `identity.groups` maps.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `creation_time` | `string` | The group's creation time in RFC3339 format |
| `id` | `string` | The group's ID |
| `last_update_time` | `string` | The group's last update (modify) time in RFC3339 format |
| `metadata` | `map (string -> string)` | Metadata associated with the group |
| `name` | `string` | The group's name |
| `member_entity_ids` | `list (string)` | A list of IDs of entities that are directly assigned to this group |
| `parent_group_ids` | `list (string)` | A list of IDs of groups that are parents of this group |
| `policies` | `list (string)` | List of the policies set on this group |

## MFA Properties

These properties exist at the `mfa` namespace.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `methods` | `map (string -> method)` | A map of method name to method properties |

### MFA Method Properties

These properties can be accessed via the `mfa.methods` selector.

| Name | Type | Description |
| :------- | :--------------------------- | :--------------------- |
| `valid` | `bool` | Whether the method has successfully been validated; if validation has not been attempted, this will trigger the validation attempt. The result of the validation attempt will be used for this method for all policies for the given request. |
