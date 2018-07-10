---
layout: "docs"
page_title: "Active Directory - Secrets Engines"
sidebar_current: "docs-secrets-active-directory"
description: |-
  The Active Directory secrets engine for Vault generates passwords dynamically based on
  roles.
---

# Active Directory Secrets Engine

The Active Directory (AD) secrets engine is a plugin residing [here](https://github.com/hashicorp/vault-plugin-secrets-active-directory).

The AD secrets engine rotates AD passwords dynamically,
and is designed for a high-load environment where many instances may be accessing
a shared password simultaneously. With a simple set up and a simple creds API,
it doesn't require instances to be manually registered in advance to gain access. 
As long as access has been granted to the creds path via a method like 
[AppRole](https://www.vaultproject.io/api/auth/approle/index.html), they're available.

Passwords are lazily rotated based on preset TTLs and can have a length configured to meet 
your needs.

## A Note on Lazy Rotation

To drive home the point that passwords are rotated "lazily", consider this scenario:

- A password is configured with a TTL of 1 hour.
- All instances of a service using this password are off for 12 hours.
- Then they wake up and again request the password.

In this scenario, although the password TTL was set to 1 hour, the password wouldn't be rotated for 12 hours when it
was next requested. "Lazy" rotation means passwords are rotated when all of the following conditions are true:

- They are over their TTL
- They are requested

Therefore, the AD TTL can be considered a soft contract. It's fulfilled when the given password is next requested. 

To ensure your passwords are rotated as expected, we'd recommend you configure services to request each password at least
twice as often as its TTL.

## A Note on Escaping

**It is up to the administrator** to provide properly escaped DNs. This
includes the user DN, bind DN for search, and so on.

The only DN escaping performed by this method is on usernames given at login
time when they are inserted into the final bind DN, and uses escaping rules
defined in RFC 4514.

Additionally, Active Directory has escaping rules that differ slightly from the
RFC; in particular it requires escaping of '#' regardless of position in the DN
(the RFC only requires it to be escaped when it is the first character), and
'=', which the RFC indicates can be escaped with a backslash, but does not
contain in its set of required escapes. If you are using Active Directory and
these appear in your usernames, please ensure that they are escaped, in
addition to being properly escaped in your configured DNs.

For reference, see [RFC 4514](https://www.ietf.org/rfc/rfc4514.txt) and this
[TechNet post on characters to escape in Active
Directory](http://social.technet.microsoft.com/wiki/contents/articles/5312.active-directory-characters-to-escape.aspx).

## Quick Setup

Most secrets engines must be configured in advance before they can perform their
functions. These steps are usually completed by an operator or configuration
management tool.
    
1. Enable the Active Directory secrets engine:

    ```text
    $ vault secrets enable ad
    Success! Enabled the ad secrets engine at: ad/
    ```

    By default, the secrets engine will mount at the name of the engine. To
    enable the secrets engine at a different path, use the `-path` argument.

2. Configure the credentials that Vault uses to communicate with Active Directory 
to generate passwords:

    ```text
    $ vault write ad/config \
        binddn=$USERNAME \
        bindpass=$PASSWORD \
        url=ldap://138.91.247.105 \
        userdn='dc=example,dc=com'
    ```

    The `$USERNAME` and `$PASSWORD` given must have access to modify passwords
    for the given account. It is possible to delegate access to change
    passwords for these accounts to the one Vault is in control of, and this is
    usually the highest-security solution.
    
    If you'd like to do a quick, insecure evaluation, also set `insecure_tls` to true. However, this is NOT RECOMMENDED
    in a production environment. In production, we recommend `insecure_tls` is false (its default) and is used with a valid 
    `certificate`.

3. Configure a role that maps a name in Vault to an account in Active Directory.
When applications request passwords, password rotation settings will be managed by
this role.

    ```text
    $ vault write ad/roles/my-application \
        service_account_name="my-application@example.com"
    ```

4. Grant "my-application" access to its creds at `ad/creds/my-application` using an 
auth method like [AppRole](https://www.vaultproject.io/api/auth/approle/index.html).

## FAQ

### What if someone directly rotates an Active Directory password that Vault is managing?

If an administrator at your company rotates a password that Vault is managing,
the next time an application asks _Vault_ for that password, Vault won't know
it. 

To maintain that application's up-time, Vault will need to return to a state of
knowing the password. Vault will generate a new password, update it, and return
it to the application(s) asking for it. This all occurs automatically, without
human intervention.

Thus, we wouldn't recommend that administrators directly rotate the passwords
for accounts that Vault is managing. This may lead to behavior the
administrator wouldn't expect, like finding very quickly afterwards that their
new password has already been changed. 

The password `ttl` on a role can be updated at any time to ensure that the
responsibility of updating passwords can be left to Vault, rather than
requiring manual administrator updates.

### Why does Vault return the last password in addition to the current one?

Active Directory promises _eventual consistency_, which means that new
passwords may not be propagated to all instances immediately. To deal with
this, Vault returns the current password with the last password if it's known.
That way, if a new password isn't fully operational, the last password can also
be used.

## API

The Active Directory secrets engine has a full HTTP API. Please see the
[Active Directory secrets engine API](/api/secret/ad/index.html) for more
details.
