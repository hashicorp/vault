---
layout: docs
page_title: OpenLDAP - Secrets Engine
description: >-
  The OpenLDAP secret engine manages OpenLDAP entry passwords.
---

# OpenLDAP Secrets Engine

@include 'x509-sha1-deprecation.mdx'

The OpenLDAP secret engine allows management of LDAP entry passwords as well as dynamic creation of credentials.
This engine supports interacting with Active Directory which is compatible with LDAP v3.

This plugin currently supports LDAP v3.

## Quick Setup

1. Enable the OpenLDAP secret engine:

   ```sh
   $ vault secrets enable openldap
   ```

   By default, the secrets engine will mount at the name of the engine. To
   enable the secrets engine at a different path, use the `-path` argument.

2. Configure the credentials that Vault uses to communicate with OpenLDAP
   to generate passwords:

   ```sh
   $ vault write openldap/config \
       binddn=$USERNAME \
       bindpass=$PASSWORD \
       url=ldaps://138.91.247.105
   ```

   Note: it's recommended a dedicated entry management account be created specifically for Vault.

3. Rotate the root password so only Vault knows the credentials:

   ```sh
   $ vault write -f openldap/rotate-root
   ```

   Note: it's not possible to retrieve the generated password once rotated by Vault.
   It's recommended a dedicated entry management account be created specifically for Vault.

### Password Generation

This engine previously allowed configuration of the length of the password that is generated
when rotating credentials. This mechanism was deprecated in Vault 1.5 in favor of
[password policies](/docs/concepts/password-policies). This means the `length` field should
no longer be used. The following password policy can be used to mirror the same behavior
that the `length` field provides:

```hcl
length=<length>
rule "charset" {
  charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
}
```

## Static Roles

### Setup

1. Configure a static role that maps a name in Vault to an entry in OpenLDAP.
   Password rotation settings will be managed by this role.

   ```sh
   $ vault write openldap/static-role/hashicorp \
       dn='uid=hashicorp,ou=users,dc=hashicorp,dc=com' \
       username='hashicorp' \
       rotation_period="24h"
   ```

2. Request credentials for the "hashicorp" role:

   ```sh
   $ vault read openldap/static-cred/hashicorp
   ```

### LDAP Password Policy

The OpenLDAP secret engine does not hash or encrypt passwords prior to modifying
values in LDAP. This behavior can cause plaintext passwords to be stored in LDAP.

To avoid having plaintext passwords stored, the LDAP server should be configured
with an LDAP password policy (ppolicy, not to be confused with a Vault password
policy). A ppolicy can enforce rules such as hashing plaintext passwords by default.

The following is an example of an LDAP password policy to enforce hashing on the
data information tree (DIT) `dc=hashicorp,dc=com`:

```
dn: cn=module{0},cn=config
changetype: modify
add: olcModuleLoad
olcModuleLoad: ppolicy

dn: olcOverlay={2}ppolicy,olcDatabase={1}mdb,cn=config
changetype: add
objectClass: olcPPolicyConfig
objectClass: olcOverlayConfig
olcOverlay: {2}ppolicy
olcPPolicyDefault: cn=default,ou=pwpolicies,dc=hashicorp,dc=com
olcPPolicyForwardUpdates: FALSE
olcPPolicyHashCleartext: TRUE
olcPPolicyUseLockout: TRUE
```

### Schema

The OpenLDAP Secret Engine supports three different schemas: `openldap` (default),
`racf` and `ad`.

#### OpenLDAP

By default the OpenLDAP Secret Engine assumes the entry password is stored in `userPassword`.
The following object classes provide `userPassword`:

- `organization`
- `organizationalUnit`
- `person`
- `posixAccount`

#### Resource Access Control Facility (RACF)

For managing IBM's Resource Access Control Facility (RACF) security system, the secret
engine must be configured to use the schema `racf`.

Generated passwords must be 8 characters or less to support RACF. The length of the
password can be configured using a [password policy](/docs/concepts/password-policies):

```bash
$ vault write openldap/config \
	binddn=$USERNAME \
	bindpass=$PASSWORD \
	url=ldaps://138.91.247.105 \
	schema=racf \
	password_policy=racf_password_policy
```

#### Active Directory (AD)

For managing Active Directory instances, the secret engine must be configured to use the
schema `ad`.

```bash
$ vault write openldap/config \
	binddn=$USERNAME \
	bindpass=$PASSWORD \
	url=ldaps://138.91.247.105 \
	schema=ad
```

### Password Rotation

Passwords can be managed in two ways:

- automatic time based rotation, and
- manual rotation.

### Auto Password Rotation

Passwords will automatically be rotated based on the `rotation_period` configured
in the static role (minimum of 5 seconds). When requesting credentials for a static
role, the response will include the time before the next rotation (`ttl`).

Auto-rotation is currently only supported for static roles. The `binddn` account used
by Vault should be rotated using the `rotate-root` endpoint to generate a password
only Vault will know.

### Manual Rotation

Static roles can be manually rotated using the `rotate-role` endpoint. When manually
rotated the rotation period will start over.

### Deleting Static Roles

Passwords are not rotated upon deletion of a static role. The password should be manually
rotated prior to deleting the role or revoking access to the static role.

## Dynamic Credentials

### Setup

Dynamic credentials can be configured by calling the `/role/:role_name` endpoint:

```bash
$ vault write openldap/role/dynamic-role \
  creation_ldif=@/path/to/creation.ldif \
  deletion_ldif=@/path/to/deletion.ldif \
  rollback_ldif=@/path/to/rollback.ldif \
  default_ttl=1h \
  max_ttl=24h
```

-> Note: The `rollback_ldif` argument is optional, but recommended. The statements within `rollback_ldif` will be
executed if the creation fails for any reason. This ensures any entities are removed in the event of a failure.

To generate credentials:

```bash
$ vault read openldap/creds/dynamic-role
Key                    Value
---                    -----
lease_id               openldap/creds/dynamic-role/HFgd6uKaDomVMvJpYbn9q4q5
lease_duration         1h
lease_renewable        true
distinguished_names    [cn=v_token_dynamic-role_FfH2i1c4dO_1611952635,ou=users,dc=learn,dc=example]
password               xWMjkIFMerYttEbzfnBVZvhRQGmhpAA0yeTya8fdmDB3LXDzGrjNEPV2bCPE9CW6
username               v_token_testrole_FfH2i1c4dO_1611952635
```

The `distinguished_names` field is an array of DNs that are created from the `creation_ldif` statements. If more than
one LDIF entry is included, the DN from each statement will be included in this field. Each entry in this field
corresponds to a single LDIF statement. No de-duplication occurs and order is maintained.

### LDIF Entries

User account management is provided through LDIF entries. The LDIF entries may be a base64-encoded version of the
LDIF string. The string will be parsed and validated to ensure that it adheres to LDIF syntax. A good reference
for proper LDIF syntax can be found [here](https://ldap.com/ldif-the-ldap-data-interchange-format/).

Some important things to remember when crafting your LDIF entries:

- There should not be any trailing spaces on any line, including empty lines
- Each `modify` block needs to be preceded with an empty line
- Multiple modifications for a `dn` can be defined in a single `modify` block. Each modification needs to close
  with a single dash (`-`)

### Active Directory (AD)

For Active Directory, there are a few additional details that are important to remember:

To create a user programmatically in AD, you first `add` a user object and then `modify` that user to provide a
password and enable the account.

- Passwords in AD are set using the `unicodePwd` field. This must be proceeded by two (2) colons (`::`).
- When setting a password programmatically in AD, the following criteria must be met:

  - The password must be enclosed in double quotes (`" "`)
  - The password must be in [`UTF16LE` format](https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-adts/6e803168-f140-4d23-b2d3-c3a8ab5917d2)
  - The password must be `base64`-encoded
  - Additional details can be found [here](https://docs.microsoft.com/en-us/troubleshoot/windows-server/identity/set-user-password-with-ldifde)

- Once a user's password has been set, it can be enabled. AD uses the `userAccountControl` field for this purpose:
  - To enable the account, set `userAccountControl` to `512`
  - You will likely also want to disable AD's password expiration for this dynamic user account. The
    `userAccountControl` value for this is: `65536`
  - `userAccountControl` flags are cumulative, so to set both of the above two flags, add up the two values
    (`512 + 65536 = 66048`): set `userAccountControl` to `66048`
  - See [here](https://docs.microsoft.com/en-us/troubleshoot/windows-server/identity/useraccountcontrol-manipulate-account-properties#property-flag-descriptions)
    for details on `userAccountControl` flags

`sAMAccountName` is a common field when working with AD users. It is used to provide compatibility with legacy
Windows NT systems and has a limit of 20 characters. Keep this in mind when defining your `username_template`.
See [here](https://docs.microsoft.com/en-us/windows/win32/adschema/a-samaccountname) for additional details.

With regard to adding dynamic users to groups, AD doesn't let you directly modify a user's `memberOf` attribute.
The `member` attribute of a group and `memberOf` attribute of a user are
[linked attributes](https://docs.microsoft.com/en-us/windows/win32/ad/linked-attributes). Linked attributes are
forward link/back link pairs, with the forward link able to be modified. In the case of AD group membership, the
group's `member` attribute is the forward link. In order to add a newly-created dynamic user to a group, we also
need to issue a `modify` request to the desired group and update the group membership with the new user.

#### Active Directory LDIF Example

The various `*_ldif` parameters are templates that use the [go template](https://golang.org/pkg/text/template/)
language. A complete LDIF example for creating an Active Directory user account is provided here for reference:

```ldif
dn: CN={{.Username}},OU=HashiVault,DC=adtesting,DC=lab
changetype: add
objectClass: top
objectClass: person
objectClass: organizationalPerson
objectClass: user
userPrincipalName: {{.Username}}@adtesting.lab
sAMAccountName: {{.Username}}

dn: CN={{.Username}},OU=HashiVault,DC=adtesting,DC=lab
changetype: modify
replace: unicodePwd
unicodePwd::{{ printf "%q" .Password | utf16le | base64 }}
-
replace: userAccountControl
userAccountControl: 66048
-

dn: CN=test-group,OU=HashiVault,DC=adtesting,DC=lab
changetype: modify
add: member
member: CN={{.Username}},OU=HashiVault,DC=adtesting,DC=lab
-
```

## API

The OpenLDAP secrets engine has a full HTTP API. Please see the [OpenLDAP secrets engine API docs](/api-docs/secret/openldap)
for more details.
