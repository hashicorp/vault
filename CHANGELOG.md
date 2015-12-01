## 0.4.0 (Unreleased)

DEPRECATIONS/BREAKING CHANGES:

 * Policy Name Casing: Policy names are now normalized to lower-case on write,
   helping prevent accidental case mismatches. For backwards compatibility,
   policy names are not currently normalized when reading or deleting. [GH-676]
 * Default etcd port number: the default connection string for the `etcd`
   physical store uses port 2379 instead of port 4001, which is the port used
   by the supported version 2.x of etcd. [GH-753]
 * As noted below in the FEATURES section, if your Vault installation contains
   a policy called `default`, new tokens created will inherit this policy
   automatically.
 * In the PKI backend there have been a few minor breaking changes:
   * The token display name is no longer a valid option for providing a base
   domain for issuance. Since this name is prepended with the name of the
   authentication backend that issued it, it provided a faulty use-case at best
   and a confusing experience at worst. We hope to figure out a better
   per-token value in a future release.
   * The `allowed_base_domain` parameter has been changed to `allowed_domains`,
   which accepts a comma-separated list of domains. This allows issuing
   certificates with DNS subjects across multiple domains. If you had a
   configured `allowed_base_domain` parameter, it will be migrated
   automatically when the role is read (either via a normal read, or via
   issuing a certificate).

FEATURES:

 * **Significantly Enhanced PKI Backend**: The `pki` backend can now generate
   and sign root CA certificates and intermediate CA CSRs. It can also now sign
   submitted client CSRs, as well as a significant number of other
   enhancements. See the updated documentation for the full API. [GH-666]
 * **CRL Checking for Certificate Authentication**: The `cert` backend now
   supports pushing CRLs into the mount and using the contained serial numbers
   for revocation checking. See the documentation for the `cert` backend for
   more info. [GH-330]
 * **Default Policy**: Vault now ensures that a policy named `default` is added
   to every token. This policy cannot be deleted, but it can be modified
   (including to an empty policy). There are three endpoints allowed in the
   default `default` policy, related to token self-management: `lookup-self`,
   which allows a token to retrieve its own information, and `revoke-self` and
   `renew-self`, which are self-explanatory. If your existing Vault
   installation contains a policy called `default`, it will not be overridden,
   but it will be added to each new token created. You can override this
   behavior when using manual token creation (i.e. not via an authentication
   backend) by setting the "no_default_policy" flag to true. [GH-732]

IMPROVEMENTS:

 * api: API client now uses a 60 second timeout instead of indefinite [GH-681]
 * api: Implement LookupSelf, RenewSelf, and RevokeSelf functions for auth
   tokens [GH-739]
 * api: Standardize environment variable reading logic inside the API; the CLI
   now uses this but can still override via command-line parameters [GH-618]
 * audit: HMAC-SHA256'd client tokens are now stored with each request entry.
   Previously they were only displayed at creation time; this allows much
   better traceability of client actions. [GH-713]
 * audit: There is now a `sys/audit-hash` endpoint that can be used to generate
   an HMAC-SHA256'd value from provided data using the given audit backend's
   salt [GH-784]
 * core: The physical storage read cache can now be disabled via
   "disable_cache" [GH-674]
 * core: The unsealing process can now be reset midway through (this feature
   was documented before, but not enabled) [GH-695]
 * core: Tokens can now renew themselves [GH-455]
 * core: Base64-encoded PGP keys can be used with the CLI for `init` and
   `rekey` operations [GH-653]
 * core: Print version on startup [GH-765]
 * core: Access to `sys/policy` and `sys/mounts` now uses the normal ACL system
   instead of requiring a root token [GH-769]
 * credential/token: Display whether or not a token is an orphan in the output
   of a lookup call [GH-766]
 * logical: Allow `.` in path-based variables in many more locations [GH-244]
 * logical: Responses now contain a "warnings" key containing a list of
   warnings returned from the server. These are conditions that did not require
   failing an operation, but of which the client should be aware. [GH-676]
 * physical/consul: Consul now uses a connection pool to limit the number of
   outstanding operations, improving behavior when a lot of operations must
   happen at once [GH-677]
 * physical/etcd: Support TLS-encrypted connections and use a connection pool
   to limit the number of outstanding operations [GH-780]
 * physical/s3: The S3 endpoint can now be configured, allowing using
   S3-API-compatible storage solutions [GH-750]
 * physical/s3: The S3 bucket can now be configured with the `AWS_S3_BUCKET`
   environment variable [GH-758]
 * secret/consul: Management tokens can now be created [GH-714]

BUG FIXES:

 * api: API client now checks for a 301 response for redirects. Vault doesn't
   generate these, but in certain conditions Go's internal HTTP handler can
generate them, leading to client errors.
 * cli: `token-create` now supports the `ttl` parameter in addition to the
   deprecated `lease` parameter. [GH-688]
 * core: Return data from `generic` backends on the last use of a limited-use
   token [GH-615]
 * core: Fix upgrade path for leases created in `generic` prior to 0.3 [GH-673]
 * core: Stale leader entries will now be reaped [GH-679]
 * core: Using `mount-tune` on the auth/token path did not take effect.
   [GH-688]
 * core: Fix a potential race condition when (un)sealing the vault with metrics
   enabled [GH-694]
 * core: Fix an error that could happen in some failure scenarios where Vault
   could fail to revert to a clean state [GH-733]
 * core: Ensure secondary indexes are removed when a lease is expired [GH-749]
 * core: Ensure rollback manager uses an up-to-date mounts table [GH-771]
 * everywhere: Don't use http.DefaultClient, as it shares state implicitly and
   is a source of hard-to-track-down bugs [GH-700]
 * credential/token: Allow creating orphan tokens via an API path [GH-748]
 * secret/generic: Validate given duration at write time, not just read time;
   if stored durations are not parseable, return a warning and the default
   duration rather than an error [GH-718]
 * secret/postgresql: Revoke permissions before dropping a user or revocation
   may fail [GH-699]

MISC:

 * Various documentation fixes and improvements [GH-685] [GH-688] [GH-697]
   [GH-710] [GH-715]

## 0.3.1 (October 6, 2015)

SECURITY:

 * core: In certain failure scenarios, the full values of requests and
   responses would be logged [GH-665]

FEATURES:

 * **Settable Maximum Open Connections**: The `mysql` and `postgresql` backends
   now allow setting the number of maximum open connections to the database,
which was previously capped to 2. [GH-661]
 * **Renewable Tokens for GitHub**: The `github` backend now supports
   specifying a TTL, enabling renewable tokens. [GH-664]

BUG FIXES:

 * dist: linux-amd64 distribution was dynamically linked [GH-656]
 * credential/github: Fix acceptance tests [GH-651]

MISC:

 * Various minor documentation fixes and improvements [GH-649] [GH-650]
   [GH-654] [GH-663]

## 0.3.0 (September 28, 2015)

DEPRECATIONS/BREAKING CHANGES:

Note: deprecations and breaking changes in upcoming releases are announced
ahead of time on the "vault-tool" mailing list.

 * **Cookie Authentication Removed**: As of 0.3 the only way to authenticate is
   via the X-Vault-Token header. Cookie authentication was hard to properly
test, could result in browsers/tools/applications saving tokens in plaintext on
disk, and other issues. [GH-564]
 * **Terminology/Field Names**: Vault is transitioning from overloading the
   term "lease" to mean both "a set of metadata" and "the amount of time the
metadata is valid". The latter is now being referred to as TTL (or
"lease_duration" for backwards-compatibility); some parts of Vault have already
switched to using "ttl" and others will follow in upcoming releases. In
particular, the "token", "generic", and "pki" backends accept both "ttl" and
"lease" but in 0.4 only "ttl" will be accepted. [GH-528]
 * **Downgrade Not Supported**: Due to enhancements in the storage subsytem,
   values written by Vault 0.3+ will not be able to be read by prior versions
of Vault. There are no expected upgrade issues, however, as with all critical
infrastructure it is recommended to back up Vault's physical storage before
upgrading.

FEATURES:

 * **SSH Backend**: Vault can now be used to delegate SSH access to machines,
   via a (recommended) One-Time Password approach or by issuing dynamic keys.
[GH-385]
 * **Cubbyhole Backend**: This backend works similarly to the "generic" backend
   but provides a per-token workspace. This enables some additional
authentication workflows (especially for containers) and can be useful to
applications to e.g. store local credentials while being restarted or upgraded,
rather than persisting to disk. [GH-612]
 * **Transit Backend Improvements**: The transit backend now allows key
   rotation and datakey generation. For rotation, data encrypted with previous
versions of the keys can still be decrypted, down to a (configurable) minimum
previous version; there is a rewrap function for manual upgrades of ciphertext
to newer versions. Additionally, the backend now allows generating and
returning high-entropy keys of a configurable bitsize suitable for AES and
other functions; this is returned wrapped by a named key, or optionally both
wrapped and plaintext for immediate use. [GH-626]
 * **Global and Per-Mount Default/Max TTL Support**: You can now set the
   default and maximum Time To Live for leases both globally and per-mount.
Per-mount settings override global settings. Not all backends honor these
settings yet, but the maximum is a hard limit enforced outside the backend. See
the documentation for "/sys/mounts/" for details on configuring per-mount TTLs.
[GH-469]
 * **PGP Encryption for Unseal Keys**: When initializing or rotating Vault's
   master key, PGP/GPG public keys can now be provided. The output keys will be
encrypted with the given keys, in order. [GH-570]
 * **Duo Multifactor Authentication Support**: Backends that support MFA can
   now use Duo as the mechanism. [GH-464]
 * **Performance Improvements**: Users of the "generic" backend will see a
   significant performance improvement as the backend no longer creates leases,
although it does return TTLs (global/mount default, or set per-item) as before.
[GH-631]
 * **Codebase Audit**: Vault's codebase was audited by iSEC. (The terms of the
   audit contract do not allow us to make the results public.) [GH-220]

IMPROVEMENTS:

 * audit: Log entries now contain a time field [GH-495]
 * audit: Obfuscated audit entries now use hmac-sha256 instead of sha1 [GH-627]
 * backends: Add ability for a cleanup function to be called on backend unmount
   [GH-608]
 * config: Allow specifying minimum acceptable TLS version [GH-447]
 * core: If trying to mount in a location that is already mounted, be more
   helpful about the error [GH-510]
 * core: Be more explicit on failure if the issue is invalid JSON [GH-553]
 * core: Tokens can now revoke themselves [GH-620]
 * credential/app-id: Give a more specific error when sending a duplicate POST
   to sys/auth/app-id [GH-392]
 * credential/github: Support custom API endpoints (e.g. for Github Enterprise)
   [GH-572]
 * credential/ldap: Add per-user policies and option to login with
   userPrincipalName [GH-420]
 * credential/token: Allow root tokens to specify the ID of a token being
   created from CLI [GH-502]
 * credential/userpass: Enable renewals for login tokens [GH-623]
 * scripts: Use /usr/bin/env to find Bash instead of hardcoding [GH-446]
 * scripts: Use godep for build scripts to use same environment as tests
   [GH-404]
 * secret/mysql: Allow reading configuration data [GH-529]
 * secret/pki: Split "allow_any_name" logic to that and "enforce_hostnames", to
   allow for non-hostname values (e.g. for client certificates) [GH-555]
 * storage/consul: Allow specifying certificates used to talk to Consul
   [GH-384]
 * storage/mysql: Allow SSL encrypted connections [GH-439]
 * storage/s3: Allow using temporary security credentials [GH-433]
 * telemetry: Put telemetry object in configuration to allow more flexibility
   [GH-419]
 * testing: Disable mlock for testing of logical backends so as not to require
   root [GH-479]

BUG FIXES:

 * audit/file: Do not enable auditing if file permissions are invalid [GH-550]
 * backends: Allow hyphens in endpoint patterns (fixes AWS and others) [GH-559]
 * cli: Fixed missing setup of client TLS certificates if no custom CA was
   provided
 * cli/read: Do not include a carriage return when using raw field output
   [GH-624]
 * core: Bad input data could lead to a panic for that session, rather than
   returning an error [GH-503]
 * core: Allow SHA2-384/SHA2-512 hashed certificates [GH-448]
 * core: Do not return a Secret if there are no uses left on a token (since it
   will be unable to be used) [GH-615]
 * core: Code paths that called lookup-self would decrement num_uses and
   potentially immediately revoke a token [GH-552]
 * core: Some /sys/ paths would not properly redirect from a standby to the
   leader [GH-499] [GH-551]
 * credential/aws: Translate spaces in a token's display name to avoid making
   IAM unhappy [GH-567]
 * credential/github: Integration failed if more than ten organizations or
   teams [GH-489]
 * credential/token: Tokens with sudo access to "auth/token/create" can now use
   root-only options [GH-629]
 * secret/cassandra: Work around backwards-incompatible change made in
   Cassandra 2.2 preventing Vault from properly setting/revoking leases
[GH-549]
 * secret/mysql: Use varbinary instead of varchar to avoid InnoDB/UTF-8 issues
   [GH-522]
 * secret/postgres: Explicitly set timezone in connections [GH-597]
 * storage/etcd: Renew semaphore periodically to prevent leadership flapping
   [GH-606]
 * storage/zk: Fix collisions in storage that could lead to data unavailability
   [GH-411]

MISC:

 * Various documentation fixes and improvements [GH-412] [GH-474] [GH-476]
   [GH-482] [GH-483] [GH-486] [GH-508] [GH-568] [GH-574] [GH-586] [GH-590]
[GH-591] [GH-592] [GH-595] [GH-613] [GH-637]
 * Less "armon" in stack traces [GH-453]
 * Sourcegraph integration [GH-456]

## 0.2.0 (July 13, 2015)

FEATURES:

 * **Key Rotation Support**: The `rotate` command can be used to rotate the
   master encryption key used to write data to the storage (physical) backend.
[GH-277]
 * **Rekey Support**: Rekey can be used to rotate the master key and change the
   configuration of the unseal keys (number of shares, threshold required).
[GH-277]
 * **New secret backend: `pki`**: Enable Vault to be a certificate authority
   and generate signed TLS certificates. [GH-310]
 * **New secret backend: `cassandra`**: Generate dynamic credentials for
   Cassandra [GH-363]
 * **New storage backend: `etcd`**: store physical data in etcd [GH-259]
   [GH-297]
 * **New storage backend: `s3`**: store physical data in S3. Does not support
   HA. [GH-242]
 * **New storage backend: `MySQL`**: store physical data in MySQL. Does not
   support HA. [GH-324]
 * `transit` secret backend supports derived keys for per-transaction unique
   keys [GH-399]

IMPROVEMENTS:

 * cli/auth: Enable `cert` method [GH-380]
 * cli/auth: read input from stdin [GH-250]
 * cli/read: Ability to read a single field from a secret [GH-257]
 * cli/write: Adding a force flag when no input required
 * core: allow time duration format in place of seconds for some inputs
 * core: audit log provides more useful information [GH-360]
 * core: graceful shutdown for faster HA failover
 * core: **change policy format** to use explicit globbing [GH-400] Any
   existing policy in Vault is automatically upgraded to avoid issues.  All
policy files must be updated for future writes. Adding the explicit glob
character `*` to the path specification is all that is required.
 * core: policy merging to give deny highest precedence [GH-400]
 * credential/app-id: Protect against timing attack on app-id
 * credential/cert: Record the common name in the metadata [GH-342]
 * credential/ldap: Allow TLS verification to be disabled [GH-372]
 * credential/ldap: More flexible names allowed [GH-245] [GH-379] [GH-367]
 * credential/userpass: Protect against timing attack on password
 * credential/userpass: Use bcrypt for password matching
 * http: response codes improved to reflect error [GH-366]
 * http: the `sys/health` endpoint supports `?standbyok` to return 200 on
   standby [GH-389]
 * secret/app-id: Support deleting AppID and UserIDs [GH-200]
 * secret/consul: Fine grained lease control [GH-261]
 * secret/transit: Decouple raw key from key management endpoint [GH-355]
 * secret/transit: Upsert named key when encrypt is used [GH-355]
 * storage/zk: Support for HA configuration [GH-252]
 * storage/zk: Changing node representation. **Backwards incompatible**.
   [GH-416]

BUG FIXES:

 * audit/file: file removing TLS connection state
 * audit/syslog: fix removing TLS connection state
 * command/*: commands accepting `k=v` allow blank values
 * core: Allow building on FreeBSD [GH-365]
 * core: Fixed various panics when audit logging enabled
 * core: Lease renewal does not create redundant lease
 * core: fixed leases with negative duration [GH-354]
 * core: token renewal does not create child token
 * core: fixing panic when lease increment is null [GH-408]
 * credential/app-id: Salt the paths in storage backend to avoid information
   leak
 * credential/cert: Fixing client certificate not being requested
 * credential/cert: Fixing panic when no certificate match found [GH-361]
 * http: Accept PUT as POST for sys/auth
 * http: Accept PUT as POST for sys/mounts [GH-349]
 * http: Return 503 when sealed [GH-225]
 * secret/postgres: Username length is capped to exceeding limit
 * server: Do not panic if backend not configured [GH-222]
 * server: Explicitly check value of tls_diable [GH-201]
 * storage/zk: Fixed issues with version conflicts [GH-190]

MISC:

 * cli/path-help: renamed from `help` to avoid confusion

## 0.1.2 (May 11, 2015)

FEATURES:

  * **New physical backend: `zookeeper`**: store physical data in Zookeeper.
    HA not supported yet.
  * **New credential backend: `ldap`**: authenticate using LDAP credentials.

IMPROVEMENTS:

  * core: Auth backends can store internal data about auth creds
  * audit: display name for auth is shown in logs [GH-176]
  * command/*: `-insecure` has been renamed to `-tls-skip-verify` [GH-130]
  * command/*: `VAULT_TOKEN` overrides local stored auth [GH-162]
  * command/server: environment variables are copy-pastable
  * credential/app-id: hash of app and user ID are in metadata [GH-176]
  * http: HTTP API accepts `X-Vault-Token` as auth header [GH-124]
  * logical/*: Generate help output even if no synopsis specified

BUG FIXES:

  * core: login endpoints should never return secrets
  * core: Internal data should never be returned from core endpoints
  * core: defer barrier initialization to as late as possible to avoid error
    cases during init that corrupt data (no data loss)
  * core: guard against invalid init config earlier
  * audit/file: create file if it doesn't exist [GH-148]
  * command/*: ignore directories when traversing CA paths [GH-181]
  * credential/*: all policy mapping keys are case insensitive [GH-163]
  * physical/consul: Fixing path for locking so HA works in every case

## 0.1.1 (May 2, 2015)

SECURITY CHANGES:

  * physical/file: create the storge with 0600 permissions [GH-102]
  * token/disk: write the token to disk with 0600 perms

IMPROVEMENTS:

  * core: Very verbose error if mlock fails [GH-59]
  * command/*: On error with TLS oversized record, show more human-friendly
    error message. [GH-123]
  * command/read: `lease_renewable` is now outputed along with the secret to
    show whether it is renewable or not
  * command/server: Add configuration option to disable mlock
  * command/server: Disable mlock for dev mode so it works on more systems

BUG FIXES:

  * core: if token helper isn't absolute, prepend with path to Vault
    executable, not "vault" (which requires PATH) [GH-60]
  * core: Any "mapping" routes allow hyphens in keys [GH-119]
  * core: Validate `advertise_addr` is a valid URL with scheme [GH-106]
  * command/auth: Using an invalid token won't crash [GH-75]
  * credential/app-id: app and user IDs can have hyphens in keys [GH-119]
  * helper/password: import proper DLL for Windows to ask password [GH-83]

## 0.1.0 (April 28, 2015)

  * Initial release
