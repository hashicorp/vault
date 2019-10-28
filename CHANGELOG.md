## 1.3 (Unreleased)

FEATURES:

 * **Vault Debug**: A new top-level subcommand, `debug`, is added that allows 
   operators to retrieve debugging information related to a particular Vault
   node. Operators can use this simple workflow to capture triaging information,
   which can then be consumed programmatically or by support and engineering teams.
   It has the abilitity to probe for config, host, metrics, pprof, server status, 
   and replication status.
 * **Recovery Mode**: Vault server can be brought up in recovery mode to resolve
   outages caused due to data store being in bad state. This is a privileged mode
   that allows `sys/raw` API calls to perform surgical corrections to the data
   store. Bad storage state can be caused by bugs. However, this is usually
   observed when known (and fixed) bugs are hit by older versions of Vault.
 * **Entropy Augmentation (Enterprise)**: Vault now supports sourcing entropy from 
   external source for critical security parameters. Currently an HSM that
   supports PKCS#11 is the only supported source.
 * **Active Directory Secret Check-In/Check-Out**: In the Active Directory secrets
   engine, users or applications can check out a service account for use, and its
   password will be rotated when it's checked back in.
 * **Vault Agent Template** Vault Agent now supports rendering templates containing 
   Vault secrets to disk, similar to Consul Template [GH-7652]
 * **Transit Key Type Support**: Signing and verification is now supported with the P-384
   (secp384r1) and P-521 (secp521r1) ECDSA curves [GH-7551] and encryption and 
   decryption is now supported via AES128-GCM96 [GH-7555]
 * **SSRF Protection for Vault Agent**: Vault Agent has a configuration option to 
   require a specific header beffore allowing requests [GH-7627]
 * **AWS Auth Method Root Rotation**: The credential used by the AWS auth method can 
   now be rotated, to ensure that only Vault knows the credentials it is using [GH-7131]
 * **New UI Features** The UI now supports managing users and groups for the 
   Userpass, Cert, Okta, and Radius auth methods.
 * **Shamir with Stored Master Key** The on disk format for Shamir seals has changed,
   allowing for a secondary cluster using Shamir downstream from a primary cluster
   using Auto Unseal. [GH-7694]
 * **Stackdriver Metrics Sink**: Vault can now send metrics to
   [Stackdriver](https://cloud.google.com/stackdriver/). See the [configuration
   documentation](https://www.vaultproject.io/docs/config/index.html) for
   details. [GH-6957]

IMPROVEMENTS:

 * auth/jwt: The redirect callback host may now be specified for CLI logins
   [JWT-71]
 * auth/jwt: Bound claims may now contain boolean values [JWT-73]
 * auth/jwt: CLI logins can now open the browser when running in WSL [JWT-77]
 * core: Exit ScanView if context has been cancelled [GH-7419]
 * core: re-encrypt barrier and recovery keys if the unseal key is updated
   [GH-7493]
 * core (enterprise): Add background seal re-wrap
 * core/metrics: Add config parameter to allow unauthenticated sys/metrics 
   access. [GH-7550]  
 * replication (enterprise): Write-Ahead-Log entries will not duplicate the
   data belonging to the encompassing physical entries of the transaction,
   thereby improving the performance and storage capacity.
 * replication (enterprise): Added more replication metrics
 * replication (enterprise): Reindex process now compares subpages for a more
   accurate indexing process.
 * replication (enterprise): Reindex API now accepts a new `skip_flush` parameter
   indicating all the changes should not be flushed while the tree is locked.
 * secrets/aws: The root config can now be read [GH-7245]
 * secrets/aws: Role paths may now contain the '@' character [GH-7553]
 * secrets/database/cassandra: Add ability to skip verfication of connection [GH-7614]
 * storage/azure: Add config parameter to Azure storage backend to allow
   specifying the ARM endpoint [GH-7567]
 * storage/cassandra: Improve storage efficiency by eliminating unnecessary
   copies of value data [GH-7199]
 * storage/raft: Improve raft write performance by utilizing FSM Batching [GH-7527]
 * storage/raft: Add support for non-voter nodes [GH-7634]
 * sys: Add a new `sys/host-info` endpoint for querying information about 
   the host [GH-7330]
 * sys: Add a new set of endpoints under `sys/pprof/` that allows profiling
   information to be extracted [GH-7473]
 * sys: Add endpoint that counts the total number of active identity entities [GH-7541]
 * sys: `sys/seal-status` now has a `storage_type` field denoting what type of storage
   the cluster is configured to use
 * sys/config: Add  a new endpoint under `sys/config/state/sanitized` that
   returns the configuration state of the server. It excludes config values
   from `storage`, `ha_storage`, and `seal` stanzas and some values
   from `telemetry` due to potential sensitive entries in those fields.
 * ui: when using raft storage, you can now join a raft cluster, download a
   snapshot, and restore a snapshot from the UI [GH-7410]
 * ui: clarify when secret version is deleted in the secret version history dropdown [GH-7714]
 * sys: Add a new `sys/internal/counters/tokens` endpoint, that counts the
   total number of active service token accessors in the shared token storage.
   [GH-7541]

BUG FIXES:

 * agent: Fix a data race on the token value for inmemsink [GH-7707]
 * auth/gcp: Fix a bug where region information in instance groups names could
   cause an authorization attempt to fail [GCP-74]
 * cli: Fix a bug where a token of an unknown format (e.g. in ~/.vault-token)
   could cause confusing error messages during `vault login` [GH-7508]
 * cli: Fix a bug where the `namespace list` command with JSON formatting 
   always returned an empty object [GH-7705]
 * identity (enterprise): Fixed identity case sensitive loading in secondary
   cluster [GH-7327]
 * raft: Fixed VAULT_CLUSTER_ADDR env being ignored at startup [GH-7619]
 * ui: using the `wrapped_token` query param will work with `redirect_to` and
   will automatically log in as intended [GH-7398]
 * ui: fix an error when initializing from the UI using PGP keys [GH-7542]
 * ui: show all active kv v2 secret versions even when `delete_version_after` is configured [GH-7685] 
 * cli: Command timeouts are now always specified solely by the
   `VAULT_CLIENT_TIMEOUT` value. [GH-7469]
 
## 1.2.4 (Unreleased)

CHANGES: 

 * auth/aws: If a custom `sts_endpoint` is configured, Vault Agent and the CLI
   should provide the corresponding region via the `region` parameter (which
   already existed as a CLI parameter, and has now been added to Agent). The
   automatic region detection added to the CLI and Agent in 1.2 has been removed.

IMPROVEMENTS:
  * cli: Ignore existing token during CLI login [GH-7508]
  * core: Log proxy settings from environment on startup [GH-7528]
  * core: Cache whether we've been initialized to reduce load on storage [GH-7549]

BUG FIXES:

 * agent: Fix handling of gzipped responses [GH-7470]
 * cli: Fix panic when pgp keys list is empty [GH-7546]
 * core: add hook for initializing seals for migration [GH-7666]
 * core (enterprise): Migrating from one auto unseal method to another never
   worked on enterprise, now it does.
 * identity: Add required field `response_types_supported` to identity token
   `.well-known/openid-configuration` response [GH-7533]
 * identity: Fixed nil pointer panic when merging entities [GH-7712]
 * secrets/database: Fix bug in combined DB secrets engine that can result in
   writes to static-roles endpoints timing out [GH-7518]
 * secrets/pki: Improve tidy to continue when value is nil [GH-7589]
 * ui (Enterprise): Allow kv v2 secrets that are gated by Control Groups to be 
   viewed in the UI [GH-7504]
 * cli: Command timeouts are now always specified solely by the
   `VAULT_CLIENT_TIMEOUT` value. [GH-7469]
   
## 1.2.3 (September 12, 2019)

FEATURES:

* **Oracle Cloud (OCI) Integration**: Vault now support using Oracle Cloud for
  storage, auto unseal, and authentication.  

IMPROVEMENTS:

 * auth/jwt: Groups claim matching now treats a string response as a single
   element list [JWT-63]
 * auth/kubernetes: enable better support for projected tokens API by allowing
   user to specify issuer [GH-65]
 * auth/pcf: The PCF auth plugin was renamed to the CF auth plugin, maintaining
   full backwards compatibility [GH-7346]
 * replication: Premium packages now come with unlimited performance standby
   nodes

BUG FIXES:

 * agent: Allow batch tokens and other non-renewable tokens to be used for
   agent operations [GH-7441]
 * auth/jwt: Fix an error where newer (v1.2) token_* configuration parameters
   were not being applied to tokens generated using the OIDC login flow
   [JWT-67]
 * seal/transit: Allow using Vault Agent for transit seal operations [GH-7441]
 * storage/couchdb: Fix a file descriptor leak [GH-7345]
 * ui: Fix a bug where the status menu would disappear when trying to revoke a
   token [GH-7337]
 * ui: Fix a regression that prevented input of custom items in search-select
   [GH-7338]
 * ui: Fix an issue with the namespace picker being unable to render nested
   namespaces named with numbers and sorting of namespaces in the picker
   [GH-7333]

## 1.2.2 (August 15, 2019)

CHANGES:

 * auth/pcf: The signature format has been updated to use the standard Base64
   encoding instead of the URL-safe variant. Signatures created using the
   previous format will continue to be accepted [PCF-27]
 * core: The http response code returned when an identity token key is not found
   has been changed from 400 to 404

IMPROVEMENTS: 

 * identity: Remove 512 entity limit for groups [GH-7317]

BUG FIXES:

 * auth/approle: Fix an error where an empty `token_type` string was not being
   correctly handled as `TokenTypeDefault` [GH-7273]
 * auth/radius: Fix panic when logging in [GH-7286]
 * ui: the string-list widget will now honor multiline input [GH-7254]
 * ui: various visual bugs in the KV interface were addressed [GH-7307]
 * ui: fixed incorrect URL to access help in LDAP auth [GH-7299]

## 1.2.1 (August 6th, 2019)

BUG FIXES:

 * agent: Fix a panic on creds pulling in some error conditions in `aws` and
   `alicloud` auth methods [GH-7238]
 * auth/approle: Fix error reading role-id on a role created pre-1.2 [GH-7231]
 * auth/token: Fix sudo check in non-root namespaces on create [GH-7224]
 * core: Fix health checks with perfstandbyok=true returning the wrong status
   code [GH-7240]
 * ui: The web CLI will now parse input as a shell string, with special
   characters escaped [GH-7206]
 * ui: The UI will now redirect to a page after authentication [GH-7088]
 * ui (Enterprise): The list of namespaces is now cleared when logging
   out [GH-7186]

## 1.2.0 (July 30th, 2019)

CHANGES:

 * Token store roles use new, common token fields for the values
   that overlap with other auth backends. `period`, `explicit_max_ttl`, and
   `bound_cidrs` will continue to work, with priority being given to the
   `token_` prefixed versions of those parameters. They will also be returned
   when doing a read on the role if they were used to provide values initially;
   however, in Vault 1.4 if `period` or `explicit_max_ttl` is zero they will no
   longer be returned. (`explicit_max_ttl` was already not returned if empty.)
 * Due to underlying changes in Go version 1.12 and Go > 1.11.5, Vault is now
   stricter about what characters it will accept in path names. Whereas before
   it would filter out unprintable characters (and this could be turned off),
   control characters and other invalid characters are now rejected within Go's
   HTTP library before the request is passed to Vault, and this cannot be
   disabled. To continue using these (e.g. for already-written paths), they
   must be properly percent-encoded (e.g. `\r` becomes `%0D`, `\x00` becomes
   `%00`, and so on).
 * The user-configured regions on the AWSKMS seal stanza will now be preferred
   over regions set in the enclosing environment.  This is a _breaking_ change.
 * All values in audit logs now are omitted if they are empty.  This helps
   reduce the size of audit log entries by not reproducing keys in each entry
   that commonly don't contain any value, which can help in cases where audit
   log entries are above the maximum UDP packet size and others.
 * Both PeriodicFunc and WALRollback functions will be called if both are
   provided. Previously WALRollback would only be called if PeriodicFunc was
   not set. See [GH-6717](https://github.com/hashicorp/vault/pull/6717) for
   details.
 * Vault now uses Go's official dependency management system, Go Modules, to
   manage dependencies. As a result to both reduce transitive dependencies for
   API library users and plugin authors, and to work around various conflicts,
   we have moved various helpers around, mostly under an `sdk/` submodule. A
   couple of functions have also moved from plugin helper code to the `api/`
   submodule. If you are a plugin author, take a look at some of our official
   plugins and the paths they are importing for guidance.
 * AppRole uses new, common token fields for values that overlap
   with other auth backends. `period` and `policies` will continue to work,
   with priority being given to the `token_` prefixed versions of those
   parameters. They will also be returned when doing a read on the role if they
   were used to provide values initially.
 * In AppRole, `"default"` is no longer automatically added to the `policies`
   parameter. This was a no-op since it would always be added anyways by
   Vault's core; however, this can now be explicitly disabled with the new
   `token_no_default_policy` field.
 * In AppRole, `bound_cidr_list` is no longer returned when reading a role
 * rollback: Rollback will no longer display log messages when it runs; it will
   only display messages on error.
 * Database plugins will now default to 4 `max_open_connections`
   rather than 2.

FEATURES:

 * **Integrated Storage**: Vault 1.2 includes a _tech preview_ of a new way to 
   manage storage directly within a Vault cluster. This new integrated storage
   solution is based on the Raft protocol which is also used to back HashiCorp
   Consul and HashiCorp Nomad.
 * **Combined DB credential rotation**: Alternative mode for the Combined DB
   Secret Engine to automatically rotate existing database account credentials
   and set Vault as the source of truth for credentials.
 * **Identity Tokens**: Vault's Identity system can now generate OIDC-compliant
   ID tokens. These customizable tokens allow encapsulating a signed, verifiable
   snapshot of identity information and metadata. They can be use by other
   applications—even those without Vault authorization—as a way of establishing
   identity based on a Vault entity.
 * **Pivotal Cloud Foundry plugin**: New auth method using Pivotal Cloud
   Foundry certificates for Vault authentication.
 * **ElasticSearch database plugin**: New ElasticSearch database plugin issues
   unique, short-lived ElasticSearch credentials.
 * **New UI Features**: An HTTP Request Volume Page and new UI for editing LDAP
   Users and Groups have been added.
 * **HA support for Postgres**: PostgreSQL versions >= 9.5 may now but used as
   and HA storage backend.
 * **KMIP secrets engine (Enterprise)**: Allows Vault to operate as a KMIP
   Server, seamlessly brokering cryptographic operations for traditional
   infrastructure.
 * Common Token Fields: Auth methods now use common fields for controlling
   token behavior, making it easier to understand configuration across methods.
 * **Vault API explorer**: The Vault UI now includes an embedded API explorer 
   where you can browse the endpoints avaliable to you and make requests. To try
   it out, open the Web CLI and type `api`.

IMPROVEMENTS:

 * agent: Allow EC2 nonce to be passed in [GH-6953]
 * agent: Add optional `namespace` parameter, which sets the default namespace
   for the auto-auth functionality [GH-6988]
 * agent: Add cert auto-auth method [GH-6652]
 * api: Add support for passing data to delete operations via `DeleteWithData`
   [GH-7139]
 * audit/file: Dramatically speed up file operations by changing
   locking/marshaling order [GH-7024]
 * auth/jwt: A JWKS endpoint may now be configured for signature verification [JWT-43]
 * auth/jwt: A new `verbose_oidc_logging` role parameter has been added to help
   troubleshoot OIDC configuration [JWT-57]
 * auth/jwt: `bound_claims` will now match received claims that are lists if any element
   of the list is one of the expected values [JWT-50]
 * auth/jwt: Leeways for `nbf` and `exp` are now configurable, as is clock skew
   leeway [JWT-53]
 * auth/kubernetes: Allow service names/namespaces to be configured as globs
   [KUBEAUTH-58]
 * auth/token: Allow the support of the identity system for the token backend
   via token roles [GH-6267]
 * auth/token: Add a large set of token configuration options to token store
   roles [GH-6662]
 * cli: `path-help` now allows `-format=json` to be specified, which will
   output OpenAPI [GH-7006]
 * cli: Add support for passing parameters to `vault delete` operations
   [GH-7139]
 * cli: Add a log-format CLI flag that can specify either "standard" or "json"
   for the log format for the `vault server`command. [GH-6840]
 * cli: Add `-dev-no-store-token` to allow dev servers to not store the
   generated token at the tokenhelper location [GH-7104]
 * identity: Allow a group alias' canonical ID to be modified
 * namespaces: Namespaces can now be created and deleted from performance
   replication secondaries
 * plugins: Change the default for `max_open_connections` for DB plugins to 4
   [GH-7093]
 * replication: Client TLS authentication is now supported when enabling or
   updating a replication secondary
 * secrets/database: Cassandra operations will now cancel on client timeout
   [GH-6954]
 * secrets/kv: Add optional `delete_version_after` parameter, which takes a
   duration and can be set on the mount and/or the metadata for a specific key
   [GH-7005]
 * storage/postgres: LIST now performs better on large datasets [GH-6546]
 * storage/s3: A new `path` parameter allows selecting the path within a bucket
   for Vault data [GH-7157]
 * ui: KV v1 and v2 will now gracefully degrade allowing a write without read
   workflow in the UI [GH-6570]
 * ui: Many visual improvements with the addition of Toolbars [GH-6626], the restyling
   of the Confirm Action component [GH-6741], and using a new set of glyphs for our
   Icon component [GH-6736]
 * ui: Lazy loading parts of the application so that the total initial payload is
   smaller [GH-6718]
 * ui: Tabbing to auto-complete in filters will first complete a common prefix if there
   is one [GH-6759]
 * ui: Removing jQuery from the application makes the initial JS payload smaller [GH-6768]
 
BUG FIXES:

 * audit: Log requests and responses due to invalid wrapping token provided
   [GH-6541]
 * audit: Fix bug preventing request counter queries from working with auditing 
   enabled [GH-6767
 * auth/aws: AWS Roles are now upgraded and saved to the latest version just
   after the AWS credential plugin is mounted. [GH-7025]
 * auth/aws: Fix a case where a panic could stem from a malformed assumed-role ARN
   when parsing this value [GH-6917]
 * auth/aws: Fix an error complaining about a read-only view that could occur
   during updating of a role when on a performance replication secondary
   [GH-6926]
 * auth/jwt: Fix a regression introduced in 1.1.1 that disabled checking of client_id
   for OIDC logins [JWT-54]
 * auth/jwt: Fix a panic during OIDC CLI logins that could occur if the Vault server
   response is empty [JWT-55]
 * auth/jwt: Fix issue where OIDC logins might intermittently fail when using
   performance standbys [JWT-61]
 * identity: Fix a case where modifying aliases of an entity could end up
   moving the entity into the wrong namespace
 * namespaces: Fix a behavior (currently only known to be benign) where we
   wouldn't delete policies through the official functions before wiping the
   namespaces on deletion
 * secrets/database: Escape username/password before using in connection URL
   [GH-7089]
 * secrets/pki: Forward revocation requests to active node when on a
   performance standby [GH-7173]
 * ui: Fix timestamp on some transit keys [GH-6827]
 * ui: Show Entities and Groups in Side Navigation [GH-7138]
 * ui: Ensure dropdown updates selected item on HTTP Request Metrics page

## 1.1.4/1.1.5 (July 25th/30th, 2019)

NOTE:

Although 1.1.4 was tagged, we realized very soon after the tag was publicly
pushed that an intended fix was accidentally left out. As a result, 1.1.4 was
not officially announced and 1.1.5 should be used as the release after 1.1.3.

IMPROVEMENTS:

 * identity: Allow a group alias' canonical ID to be modified
 * namespaces: Improve namespace deletion performance [GH-6939]
 * namespaces: Namespaces can now be created and deleted from performance 
   replication secondaries

BUG FIXES:

 * api: Add backwards compat support for API env vars [GH-7135]
 * auth/aws: Fix a case where a panic could stem from a malformed assumed-role
   ARN when parsing this value [GH-6917]
 * auth/ldap: Add `use_pre111_group_cn_behavior` flag to allow recovering from
   a regression caused by a bug fix starting in 1.1.1 [GH-7208]
 * auth/aws: Use a role cache to avoid separate locking paths [GH-6926]
 * core: Fix a deadlock if a panic happens during request handling [GH-6920]
 * core: Fix an issue that may cause key upgrades to not be cleaned up properly
   [GH-6949]
 * core: Don't shutdown if key upgrades fail due to canceled context [GH-7070]
 * core: Fix panic caused by handling requests while vault is inactive
 * identity: Fix reading entity and groups that have spaces in their names 
   [GH-7055]
 * identity: Ensure entity alias operations properly verify namespace [GH-6886]
 * mfa: Fix a nil pointer panic that could occur if invalid Duo credentials
   were supplied
 * replication: Forward step-down on perf standbys to match HA behavior
 * replication: Fix various read only storage errors on performance standbys
 * replication: Stop forwarding before stopping replication to eliminate some
   possible bad states
 * secrets/database: Allow cassandra queries to be cancled [GH-6954]
 * storage/consul: Fix a regression causing vault to not connect to consul over
   unix sockets [GH-6859]
 * ui: Fix saving of TTL and string array fields generated by Open API [GH-7094]
 
## 1.1.3 (June 5th, 2019)

IMPROVEMENTS:

 * agent: Now supports proxying request query parameters [GH-6772]
 * core: Mount table output now includes a UUID indicating the storage path [GH-6633]
 * core: HTTP server timeout values are now configurable [GH-6666]
 * replication: Improve performance of the reindex operation on secondary clusters
   when mount filters are in use
 * replication: Replication status API now returns the state and progress of a reindex

BUG FIXES:

 * api: Return the Entity ID in the secret output [GH-6819]
 * auth/jwt: Consider bound claims when considering if there is at least one
   bound constraint [JWT-49]
 * auth/okta: Fix handling of group names containing slashes [GH-6665]
 * cli: Add deprecated stored-shares flag back to the init command [GH-6677]
 * cli: Fix a panic when the KV command would return no data [GH-6675]
 * cli: Fix issue causing CLI list operations to not return proper format when
   there is an empty response [GH-6776]
 * core: Correctly honor non-HMAC request keys when auditing requests [GH-6653]
 * core: Fix the `x-vault-unauthenticated` value in OpenAPI for a number of
   endpoints [GH-6654]
 * core: Fix issue where some OpenAPI parameters were incorrectly listed as
   being sent as a header [GH-6679]
 * core: Fix issue that would allow duplicate mount names to be used [GH-6771]
 * namespaces: Fix behavior when using `root` instead of `root/` as the
   namespace header value
 * pki: fix a panic when a client submits a null value [GH-5679]
 * replication: Properly update mount entry cache on a secondary to apply all
   new values after a tune
 * replication: Properly close connection on bootstrap error
 * replication: Fix an issue causing startup problems if a namespace policy
   wasn't replicated properly
 * replication: Fix longer than necessary WAL replay during an initial reindex
 * replication: Fix error during mount filter invalidation on DR secondary clusters
 * secrets/ad: Make time buffer configurable [AD-35]
 * secrets/gcp: Check for nil config when getting credentials [SGCP-35]
 * secrets/gcp: Fix error checking in some cases where the returned value could
   be 403 instead of 404 [SGCP-37]
 * secrets/gcpkms: Disable key rotation when deleting a key [GCPKMS-10]
 * storage/consul: recognize `https://` address even if schema not specified
   [GH-6602]
 * storage/dynamodb: Fix an issue where a deleted lock key in DynamoDB (HA)
   could cause constant switching of the active node [GH-6637]
 * storage/dynamodb: Eliminate a high-CPU condition that could occur if an
   error was received from the DynamoDB API [GH-6640]
 * storage/gcs: Correctly use configured chunk size values [GH-6655]
 * storage/mssql: Use the correct database when pre-created schemas exist
   [GH-6356]
 * ui: Fix issue with select arrows on drop down menus [GH-6627]
 * ui: Fix an issue where sensitive input values weren't being saved to the
   server [GH-6586]
 * ui: Fix web cli parsing when using quoted values [GH-6755]
 * ui: Fix a namespace workflow mapping identities from external namespaces by
   allowing arbitrary input in search-select component [GH-6728]

## 1.1.2 (April 18th, 2019)

This is a bug fix release containing the two items below. It is otherwise
unchanged from 1.1.1.

BUG FIXES:

 * auth/okta: Fix a potential dropped error [GH-6592]
 * secrets/kv: Fix a regression on upgrade where a KVv2 mount could fail to be
   mounted on unseal if it had previously been mounted but not written to
   [KV-31]

## 1.1.1 (April 11th, 2019)

SECURITY:

 * Given: (a) performance replication is enabled; (b) performance standbys are
   in use on the performance replication secondary cluster; and (c) mount
   filters are in use, if a mount that was previously available to a secondary
   is updated to be filtered out, although the data would be removed from the
   secondary cluster, the in-memory cache of the data would not be purged on
   the performance standby nodes. As a result, the previously-available data
   could still be read from memory if it was ever read from disk, and if this
   included mount configuration data this could result in token or lease
   issuance. The issue is fixed in this release; in prior releases either an
   active node changeover (such as a step-down) or a restart of the standby
   nodes is sufficient to cause the performance standby nodes to clear their
   cache. A CVE is in the process of being issued; the number is
   CVE-2019-11075.
 * Roles in the JWT Auth backend using the OIDC login flow (i.e. role_type of
   “oidc”) were not enforcing bound_cidrs restrictions, if any were configured
   for the role. This issue did not affect roles of type “jwt”.

CHANGES:

 * auth/jwt: Disallow logins of role_type "oidc" via the `/login` path [JWT-38]
 * core/acl:  New ordering defines which policy wins when there are multiple
   inexact matches and at least one path contains `+`. `+*` is now illegal in
   policy paths. The previous behavior simply selected any matching
   segment-wildcard path that matched. [GH-6532]
 * replication: Due to technical limitations, mounting and unmounting was not
   previously possible from a performance secondary. These have been resolved,
   and these operations may now be run from a performance secondary.

IMPROVEMENTS:

 * agent: Allow AppRole auto-auth without a secret-id [GH-6324]
 * auth/gcp: Cache clients to improve performance and reduce open file usage
 * auth/jwt: Bounds claims validiation will now allow matching the received
   claims against a list of expected values [JWT-41]
 * secret/gcp: Cache clients to improve performance and reduce open file usage
 * replication: Mounting/unmounting/remounting/mount-tuning is now supported
   from a performance secondary cluster
 * ui: Suport for authentication via the RADIUS auth method [GH-6488]
 * ui: Navigating away from secret list view will clear any page-specific
   filter that was applied [GH-6511]
 * ui: Improved the display when OIDC auth errors [GH-6553]

BUG FIXES:

 * agent: Allow auto-auth to be used with caching without having to define any
   sinks [GH-6468]
 * agent: Disallow some nonsensical config file combinations [GH-6471]
 * auth/ldap: Fix CN check not working if CN was not all in uppercase [GH-6518]
 * auth/jwt: The CLI helper for OIDC logins will now open the browser to the correct
   URL when running on Windows [JWT-37]
 * auth/jwt: Fix OIDC login issue where configured TLS certs weren't being used [JWT-40]
 * auth/jwt: Fix an issue where the `oidc_scopes` parameter was not being included in
   the response to a role read request [JWT-35]
 * core: Fix seal migration case when migrating to Shamir and a seal block
   wasn't explicitly specified [GH-6455]
 * core: Fix unwrapping when using namespaced wrapping tokens [GH-6536]
 * core: Fix incorrect representation of required properties in OpenAPI output
   [GH-6490]
 * core: Fix deadlock that could happen when using the UI [GH-6560]
 * identity: Fix updating groups removing existing members [GH-6527]
 * identity: Properly invalidate group alias in performance secondary [GH-6564]
 * identity: Use namespace context when loading entities and groups to ensure
   merging of duplicate entries works properly [GH-6563]
 * replication: Fix performance standby election failure [GH-6561]
 * replication: Fix mount filter invalidation on performance standby nodes
 * replication: Fix license reloading on performance standby nodes
 * replication: Fix handling of control groups on performance standby nodes
 * replication: Fix some forwarding scenarios with request bodies using
   performance standby nodes [GH-6538]
 * secret/gcp: Fix roleset binding when using JSON [GCP-27]
 * secret/pki: Use `uri_sans` param in when not using CSR parameters [GH-6505]
 * storage/dynamodb: Fix a race condition possible in HA configurations that could
   leave the cluster without a leader [GH-6512]
 * ui: Fix an issue where in production builds OpenAPI model generation was
   failing, causing any form using it to render labels with missing fields [GH-6474]
 * ui: Fix issue nav-hiding when moving between namespaces [GH-6473]
 * ui: Secrets will always show in the nav regardless of access to cubbyhole [GH-6477]
 * ui: fix SSH OTP generation [GH-6540]
 * ui: add polyfill to load UI in IE11 [GH-6567]
 * ui: Fix issue where some elements would fail to work properly if using ACLs
   with segment-wildcard paths (`/+/` segments) [GH-6525]

## 1.1.0 (March 18th, 2019)

CHANGES:

 * auth/jwt: The `groups_claim_delimiter_pattern` field has been removed. If the
   groups claim is not at the top level, it can now be specified as a
   [JSONPointer](https://tools.ietf.org/html/rfc6901).
 * auth/jwt: Roles now have a "role type" parameter with a default type of
   "oidc". To configure new JWT roles, a role type of "jwt" must be explicitly
   specified.
 * cli: CLI commands deprecated in 0.9.2 are now removed. Please see the CLI
   help/warning output in previous versions of Vault for updated commands.
 * core: Vault no longer automatically mounts a K/V backend at the "secret/"
   path when initializing Vault
 * core: Vault's cluster port will now be open at all times on HA standby nodes
 * plugins: Vault no longer supports running netRPC plugins. These were
   deprecated in favor of gRPC based plugins and any plugin built since 0.9.4
   defaults to gRPC. Older plugins may need to be recompiled against the latest
   Vault dependencies.

FEATURES:

 * **Vault Agent Caching**: Vault Agent can now be configured to act as a
   caching proxy to Vault. Clients can send requests to Vault Agent and the
   request will be proxied to the Vault server and cached locally in Agent.
   Currently Agent will cache generated leases and tokens and keep them
   renewed. The proxy can also use the Auto Auth feature so clients do not need
   to authenticate to Vault, but rather can make requests to Agent and have
   Agent fully manage token lifecycle.
 * **OIDC Redirect Flow Support**: The JWT auth backend now supports OIDC
   roles. These allow authentication via an OIDC-compliant provider via the
   user's browser. The login may be initiated from the Vault UI or through
   the `vault login` command.
 * **ACL Path Wildcard**: ACL paths can now use the `+` character to enable
   wild card matching for a single directory in the path definition.
 * **Transit Auto Unseal**: Vault can now be configured to use the Transit
   Secret Engine in another Vault cluster as an auto unseal provider.

IMPROVEMENTS:

 * auth/jwt: A default role can be set. It will be used during JWT/OIDC logins if
   a role is not specified.
 * auth/jwt: Arbitrary claims data can now be copied into token & alias metadata.
 * auth/jwt: An arbitrary set of bound claims can now be configured for a role.
 * auth/jwt: The name "oidc" has been added as an alias for the jwt backend. Either
   name may be specified in the `auth enable` command.
 * command/server: A warning will be printed when 'tls_cipher_suites' includes a
   blacklisted cipher suite or all cipher suites are blacklisted by the HTTP/2
   specification [GH-6300]
 * core/metrics: Prometheus pull support using a new sys/metrics endpoint. [GH-5308]
 * core: On non-windows platforms a SIGUSR2 will make the server log a dump of
   all running goroutines' stack traces for debugging purposes [GH-6240]
 * replication: The initial replication indexing process on newly initialized or upgraded
   clusters now runs asynchronously
 * sentinel: Add token namespace id and path, available in rules as
   token.namespace.id and token.namespace.path
 * ui: The UI is now leveraging OpenAPI definitions to pull in fields for various forms.
   This means, it will not be necessary to add fields on the go and JS sides in the future.
   [GH-6209]

BUG FIXES:

 * auth/jwt: Apply `bound_claims` validation across all login paths
 * auth/jwt: Update `bound_audiences` validation during non-OIDC logins to accept
   any matched audience, as documented and handled in OIDC logins [JWT-30]
 * auth/token: Fix issue where empty values for token role update call were
   ignored [GH-6314]
 * core: The `operator migrate` command will no longer hang on empty key names
   [GH-6371]
 * identity: Fix a panic at login when external group has a nil alias [GH-6230]
 * namespaces: Clear out identity store items upon namespace deletion
 * replication/perfstandby: Fixed a bug causing performance standbys to wait
   longer than necessary after forwarding a write to the active node
 * replication/mountfilter: Fix a deadlock that could occur when mount filters
   were updated [GH-6426]
 * secret/kv: Fix issue where a v1→v2 upgrade could run on a performance
   standby when using a local mount
 * secret/ssh: Fix for a bug where attempting to delete the last ssh role
   in the zeroaddress configuration could fail [GH-6390]
 * secret/totp: Uppercase provided keys so they don't fail base32 validation
   [GH-6400]
 * secret/transit: Multiple HMAC, Sign or Verify operations can now be
   performed with one API call using the new `batch_input` parameter [GH-5875]
 * sys: `sys/internal/ui/mounts` will no longer return secret or auth mounts
   that have been filtered. Similarly, `sys/internal/ui/mount/:path` will
   return a error response if a filtered mount path is requested. [GH-6412]
 * ui: Fix for a bug where you couldn't access the data tab after clicking on
   wrap details on the unwrap page [GH-6404]
 * ui: Fix an issue where the policies tab was erroneously hidden [GH-6301]
 * ui: Fix encoding issues with kv interfaces [GH-6294]

## 1.0.3.1 (March 14th, 2019) (Enterprise Only)

SECURITY:

 * A regression was fixed in replication mount filter code introduced in Vault
   1.0 that caused the underlying filtered data to be replicated to
   secondaries. This data was not accessible to users via Vault's API but via a
   combination of privileged configuration file changes/Vault commands it could
   be read.  Upgrading to this version or 1.1 will fix this issue and cause the
   replicated data to be deleted from filtered secondaries. More information
   was sent to customer contacts on file.

## 1.0.3 (February 12th, 2019)

CHANGES:

 * New AWS authentication plugin mounts will default to using the generated
   role ID as the Identity alias name. This applies to both EC2 and IAM auth.
   Existing mounts that explicitly set this value will not be affected but
   mounts that specified no preference will switch over on upgrade.
 * The default policy now allows a token to look up its associated identity
   entity either by name or by id [GH-6105]
 * The Vault UI's navigation and onboarding wizard now only displays items that
   are permitted in a users' policy [GH-5980, GH-6094]
 * An issue was fixed that caused recovery keys to not work on secondary
   clusters when using a different unseal mechanism/key than the primary. This
   would be hit if the cluster was rekeyed or initialized after 1.0. We recommend
   rekeying the recovery keys on the primary cluster if you meet the above
   requirements.

FEATURES:

 * **cURL Command Output**: CLI commands can now use the `-output-curl-string`
   flag to print out an equivalent cURL command.
 * **Response Headers From Plugins**: Plugins can now send back headers that
   will be included in the response to a client. The set of allowed headers can
   be managed by the operator.

IMPROVEMENTS:

 * auth/aws: AWS EC2 authentication can optionally create entity aliases by
   role ID [GH-6133]
 * auth/jwt: The supported set of signing algorithms is now configurable [JWT
   plugin GH-16]
 * core: When starting from an uninitialized state, HA nodes will now attempt
   to auto-unseal using a configured auto-unseal mechanism after the active
   node initializes Vault [GH-6039]
 * secret/database: Add socket keepalive option for Cassandra [GH-6201]
 * secret/ssh: Add signed key constraints, allowing enforcement of key types
   and minimum key sizes [GH-6030]
 * secret/transit: ECDSA signatures can now be marshaled in JWS-compatible
   fashion [GH-6077]
 * storage/etcd: Support SRV service names [GH-6087]
 * storage/aws: Support specifying a KMS key ID for server-side encryption
   [GH-5996]

BUG FIXES:

 * core: Fix a rare case where a standby whose connection is entirely torn down
   to the active node, then reconnects to the same active node, may not
   successfully resume operation [GH-6167]
 * cors: Don't duplicate headers when they're written [GH-6207]
 * identity: Persist merged entities only on the primary [GH-6075]
 * replication: Fix a potential race when a token is created and then used with
   a performance standby very quickly, before an associated entity has been
   replicated. If the entity is not found in this scenario, the request will
   forward to the active node.
 * replication: Fix issue where recovery keys would not work on secondary
   clusters if using a different unseal mechanism than the primary.
 * replication: Fix a "failed to register lease" error when using performance
   standbys
 * storage/postgresql: The `Get` method will now return an Entry object with
   the `Key` member correctly populated with the full path that was requested
   instead of just the last path element [GH-6044]

## 1.0.2 (January 15th, 2019)

SECURITY:

 * When creating a child token from a parent with `bound_cidrs`, the list of
   CIDRs would not be propagated to the child token, allowing the child token
   to be used from any address.

CHANGES:

 * secret/aws: Role now returns `credential_type` instead of `credential_types`
   to match role input. If a legacy role that can supply more than one
   credential type, they will be concatenated with a `,`.
 * physical/dynamodb, autoseal/aws: Instead of Vault performing environment
   variable handling, and overriding static (config file) values if found, we
   use the default AWS SDK env handling behavior, which also looks for
   deprecated values. If you were previously providing both config values and
   environment values, please ensure the config values are unset if you want to
   use environment values.
 * Namespaces (Enterprise): Providing "root" as the header value for
   `X-Vault-Namespace` will perform the request on the root namespace. This is
   equivalent to providing an empty value. Creating a namespace called "root" in
   the root namespace is disallowed.

FEATURES:

 * **InfluxDB Database Plugin**: Use Vault to dynamically create and manage InfluxDB
   users

IMPROVEMENTS:

 * auth/aws: AWS EC2 authentication can optionally create entity aliases by
   image ID [GH-5846]
 * autoseal/gcpckms: Reduce the required permissions for the GCPCKMS autounseal
   [GH-5999]
 * physical/foundationdb: TLS support added. [GH-5800]

BUG FIXES:

 * api: Fix a couple of places where we were using the `LIST` HTTP verb
   (necessary to get the right method into the wrapping lookup function) and
   not then modifying it to a `GET`; although this is officially the verb Vault
   uses for listing and it's fully legal to use custom verbs, since many WAFs
   and API gateways choke on anything outside of RFC-standardized verbs we fall
   back to `GET` [GH-6026]
 * autoseal/aws: Fix reading session tokens when AWS access key/secret key are
   also provided [GH-5965]
 * command/operator/rekey: Fix help output showing `-delete-backup` when it
   should show `-backup-delete` [GH-5981]
 * core: Fix bound_cidrs not being propagated to child tokens
 * replication: Correctly forward identity entity creation that originates from
   performance standby nodes (Enterprise)
 * secret/aws: Make input `credential_type` match the output type (string, not
   array) [GH-5972]
 * secret/cubbyhole: Properly cleanup cubbyhole after token revocation [GH-6006]
 * secret/pki: Fix reading certificates on windows with the file storage backend [GH-6013]
 * ui (enterprise): properly display perf-standby count on the license page [GH-5971]
 * ui: fix disappearing nested secrets and go to the nearest parent when deleting
   a secret - [GH-5976]
 * ui: fix error where deleting an item via the context menu would fail if the
   item name contained dots [GH-6018]
 * ui: allow saving of kv secret after an errored save attempt [GH-6022]
 * ui: fix display of kv-v1 secret containing a key named "keys" [GH-6023]

## 1.0.1 (December 14th, 2018)

SECURITY:

 * Update version of Go to 1.11.3 to fix Go bug
   https://github.com/golang/go/issues/29233 which corresponds to
   CVE-2018-16875
 * Database user revocation: If a client has configured custom revocation
   statements for a role with a value of `""`, that statement would be executed
   verbatim, resulting in a lack of actual revocation but success for the
   operation. Vault will now strip empty statements from any provided; as a
   result if an empty statement is provided, it will behave as if no statement
   is provided, falling back to the default revocation statement.

CHANGES:

 * secret/database: On role read, empty statements will be returned as empty
   slices instead of potentially being returned as JSON null values. This makes
   it more in line with other parts of Vault and makes it easier for statically
   typed languages to interpret the values.

IMPROVEMENTS:

 * cli: Strip iTerm extra characters from password manager input [GH-5837]
 * command/server: Setting default kv engine to v1 in -dev mode can now be
   specified via -dev-kv-v1 [GH-5919]
 * core: Add operationId field to OpenAPI output [GH-5876]
 * ui: Added ability to search for Group and Policy IDs when creating Groups
   and Entities instead of typing them in manually

BUG FIXES:

 * auth/azure: Cache azure authorizer [15]
 * auth/gcp: Remove explicit project for service account in GCE authorizer [58]
 * cli: Show correct stored keys/threshold for autoseals [GH-5910]
 * cli: Fix backwards compatibility fallback when listing plugins [GH-5913]
 * core: Fix upgrades when the seal config had been created on early versions
   of vault [GH-5956]
 * namespaces: Correctly reload the proper mount when tuning or reloading the
   mount [GH-5937]
 * secret/azure: Cache azure authorizer [19]
 * secret/database: Strip empty statements on user input [GH-5955]
 * secret/gcpkms: Add path for retrieving the public key [5]
 * secret/pki: Fix panic that could occur during tidy operation when malformed
   data was found [GH-5931]
 * secret/pki: Strip empty line in ca_chain output [GH-5779]
 * ui: Fixed a bug where the web CLI was not usable via the `fullscreen`
   command - [GH-5909]
 * ui: Fix a bug where you couldn't write a jwt auth method config [GH-5936]

## 0.11.6 (December 14th, 2018)

This release contains the three security fixes from 1.0.0 and 1.0.1 and the
following bug fixes from 1.0.0/1.0.1:

 * namespaces: Correctly reload the proper mount when tuning or reloading the
   mount [GH-5937]
 * replication/perfstandby: Fix audit table upgrade on standbys [GH-5811]
 * replication/perfstandby: Fix redirect on approle update [GH-5820]
 * secrets/kv: Fix issue where storage version would get incorrectly downgraded
   [GH-5809]

It is otherwise identical to 0.11.5.

## 1.0.0 (December 3rd, 2018)

SECURITY:

 * When debugging a customer incident we discovered that in the case of
   malformed data from an autoseal mechanism, Vault's master key could be
   logged in Vault's server log. For this to happen, the data would need to be
   modified by the autoseal mechanism after being submitted to it by Vault but
   prior to encryption, or after decryption, prior to it being returned to
   Vault. To put it another way, it requires the data that Vault submits for
   encryption to not match the data returned after decryption. It is not
   sufficient for the autoseal mechanism to return an error, and it cannot be
   triggered by an outside attacker changing the on-disk ciphertext as all
   autoseal mechanisms use authenticated encryption. We do not believe that
   this is generally a cause for concern; since it involves the autoseal
   mechanism returning bad data to Vault but with no error, in a working Vault
   configuration this code path should never be hit, and if hitting this issue
   Vault will not be unsealing properly anyways so it will be obvious what is
   happening and an immediate rekey of the master key can be performed after
   service is restored. We have filed for a CVE (CVE-2018-19786) and a CVSS V3
   score of 5.2 has been assigned.

CHANGES:

 * Tokens are now prefixed by a designation to indicate what type of token they
   are. Service tokens start with `s.` and batch tokens start with `b.`.
   Existing tokens will still work (they are all of service type and will be
   considered as such). Prefixing allows us to be more efficient when consuming
   a token, which keeps the critical path of requests faster.
 * Paths within `auth/token` that allow specifying a token or accessor in the
   URL have been removed. These have been deprecated since March 2016 and
   undocumented, but were retained for backwards compatibility. They shouldn't
   be used due to the possibility of those paths being logged, so at this point
   they are simply being removed.
 * Vault will no longer accept updates when the storage key has invalid UTF-8
   character encoding [GH-5819]
 * Mount/Auth tuning the `options` map on backends will now upsert any provided
   values, and keep any of the existing values in place if not provided. The
   options map itself cannot be unset once it's set, but the keypairs within the
   map can be unset if an empty value is provided, with the exception of the
   `version` keypair which is handled differently for KVv2 purposes.
 * Agent no longer automatically reauthenticates when new credentials are
   detected. It's not strictly necessary and in some cases was causing
   reauthentication much more often than intended.
 * HSM Regenerate Key Support Removed: Vault no longer supports destroying and
   regenerating encryption keys on an HSM; it only supports creating them.
   Although this has never been a source of a customer incident, it is simply a
   code path that is too trivial to activate, especially by mistyping
   `regenerate_key` instead of `generate_key`.
 * Barrier Config Upgrade (Enterprise): When upgrading from Vault 0.8.x, the
   seal type in the barrier config storage entry will be upgraded from
   "hsm-auto" to "awskms" or "pkcs11" upon unseal if using AWSKMS or HSM seals.
   If performing seal migration, the barrier config should first be upgraded
   prior to starting migration.
 * Go API client uses pooled HTTP client: The Go API client now uses a
   connection-pooling HTTP client by default. For CLI operations this makes no
   difference but it should provide significant performance benefits for those
   writing custom clients using the Go API library. As before, this can be
   changed to any custom HTTP client by the caller.
 * Builtin Secret Engines and Auth Methods are integrated deeper into the
   plugin system. The plugin catalog can now override builtin plugins with
   custom versions of the same name. Additionally the plugin system now
   requires a plugin `type` field when configuring plugins, this can be "auth",
   "database", or "secret".

FEATURES:

 * **Auto-Unseal in Open Source**: Cloud-based auto-unseal has been migrated
   from Enterprise to Open Source. We've created a migrator to allow migrating
   between Shamir seals and auto unseal methods.
 * **Batch Tokens**: Batch tokens trade off some features of service tokens for no
   storage overhead, and in most cases can be used across performance
   replication clusters.
 * **Replication Speed Improvements**: We've worked hard to speed up a lot of
   operations when using Vault Enterprise Replication.
 * **GCP KMS Secrets Engine**: This new secrets engine provides a Transit-like
   pattern to keys stored within GCP Cloud KMS.
 * **AppRole support in Vault Agent Auto-Auth**: You can now use AppRole
   credentials when having Agent automatically authenticate to Vault
 * **OpenAPI Support**: Descriptions of mounted backends can be served directly
   from Vault
 * **Kubernetes Projected Service Account Tokens**: Projected Service Account
   Tokens are now supported in Kubernetes auth
 * **Response Wrapping in UI**: Added ability to wrap secrets and easily copy
   the wrap token or secret JSON in the UI

IMPROVEMENTS:

 * agent: Support for configuring the location of the kubernetes service account
   [GH-5725]
 * auth/token: New tokens are indexed in storage HMAC-SHA256 instead of SHA1
 * secret/totp: Allow @ character to be part of key name [GH-5652]
 * secret/consul: Add support for new policy based tokens added in Consul 1.4
   [GH-5586]
 * ui: Improve the token auto-renew warning, and automatically begin renewal
   when a user becomes active again [GH-5662]
 * ui: The unbundled UI page now has some styling [GH-5665]
 * ui: Improved banner and popup design [GH-5672]
 * ui: Added token type to auth method mount config [GH-5723]
 * ui: Display additonal wrap info when unwrapping. [GH-5664]
 * ui: Empty states have updated styling and link to relevant actions and
   documentation [GH-5758]
 * ui: Allow editing of KV V2 data when a token doesn't have capabilities to
   read secret metadata [GH-5879]

BUG FIXES:

 * agent: Fix auth when multiple redirects [GH-5814]
 * cli: Restore the `-policy-override` flag [GH-5826]
 * core: Fix rekey progress reset which did not happen under certain
   circumstances. [GH-5743]
 * core: Migration from autounseal to shamir will clean up old keys [GH-5671]
 * identity: Update group memberships when entity is deleted [GH-5786]
 * replication/perfstandby: Fix audit table upgrade on standbys [GH-5811]
 * replication/perfstandby: Fix redirect on approle update [GH-5820]
 * secrets/azure: Fix valid roles being rejected for duplicate ids despite
   having distinct scopes
   [[GH-16]](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/16)
 * storage/gcs: Send md5 of values to GCS to avoid potential corruption
   [GH-5804]
 * secrets/kv: Fix issue where storage version would get incorrectly downgraded
   [GH-5809]
 * secrets/kv: Disallow empty paths on a `kv put` while accepting empty paths
   for all other operations for backwards compatibility
   [[GH-19]](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/19)
 * ui: Allow for secret creation in kv v2 when cas_required=true [GH-5823]
 * ui: Fix dr secondary operation token generation via the ui [GH-5818]
 * ui: Fix the PKI context menu so that items load [GH-5824]
 * ui: Update DR Secondary Token generation command [GH-5857]
 * ui: Fix pagination bug where controls would be rendered once for each
   item when viewing policies [GH-5866]
 * ui: Fix bug where `sys/leases/revoke` required 'sudo' capability to show
   the revoke button in the UI [GH-5647]
 * ui: Fix issue where certain pages wouldn't render in a namespace [GH-5692]

## 0.11.5 (November 13th, 2018)

BUG FIXES:

 * agent: Fix issue when specifying two file sinks [GH-5610]
 * auth/userpass: Fix minor timing issue that could leak the presence of a
   username [GH-5614]
 * autounseal/alicloud: Fix issue interacting with the API (Enterprise)
 * autounseal/azure: Fix key version tracking (Enterprise)
 * cli: Fix panic that could occur if parameters were not provided [GH-5603]
 * core: Fix buggy behavior if trying to remount into a namespace
 * identity: Fix duplication of entity alias entity during alias transfer
   between entities [GH-5733]
 * namespaces: Fix tuning of auth mounts in a namespace
 * ui: Fix bug where editing secrets as JSON doesn't save properly [GH-5660]
 * ui: Fix issue where IE 11 didn't render the UI and also had a broken form
   when trying to use tool/hash [GH-5714]

## 0.11.4 (October 23rd, 2018)

CHANGES:

 * core: HA lock file is no longer copied during `operator migrate` [GH-5503].
   We've categorized this as a change, but generally this can be considered
   just a bug fix, and no action is needed.

FEATURES:

 * **Transit Key Trimming**: Keys in transit secret engine can now be trimmed to
   remove older unused key versions
 * **Web UI support for KV Version 2**: Browse, delete, undelete and destroy
   individual secret versions in the UI
 * **Azure Existing Service Principal Support**: Credentials can now be generated
   against an existing service principal

IMPROVEMENTS:

 * core: Add last WAL in leader/health output for easier debugging [GH-5523]
 * identity: Identity names will now be handled case insensitively by default.
   This includes names of entities, aliases and groups [GH-5404]
 * secrets/aws: Added role-option max_sts_ttl to cap TTL for AWS STS
   credentials [GH-5500]
 * secret/database: Allow Cassandra user to be non-superuser so long as it has
   role creation permissions [GH-5402]
 * secret/radius: Allow setting the NAS Identifier value in the generated
   packet [GH-5465]
 * secret/ssh: Allow usage of JSON arrays when setting zero addresses [GH-5528]
 * secret/transit: Allow trimming unused keys [GH-5388]
 * ui: Support KVv2 [GH-5547], [GH-5563]
 * ui: Allow viewing and updating Vault license via the UI
 * ui: Onboarding will now display your progress through the chosen tutorials
 * ui: Dynamic secret backends obfuscate sensitive data by default and
   visibility is toggleable

BUG FIXES:

 * agent: Fix potential hang during agent shutdown [GH-5026]
 * auth/ldap: Fix listing of users/groups that contain slashes [GH-5537]
 * core: Fix memory leak during some expiration calls [GH-5505]
 * core: Fix generate-root operations requiring empty `otp` to be provided
   instead of an empty body [GH-5495]
 * identity: Remove lookup check during alias removal from entity [GH-5524]
 * secret/pki: Fix TTL/MaxTTL check when using `sign-verbatim` [GH-5549]
 * secret/pki: Fix regression in 0.11.2+ causing the NotBefore value of
   generated certificates to be set to the Unix epoch if the role value was not
   set, instead of using the default of 30 seconds [GH-5481]
 * storage/mysql: Use `varbinary` instead of `varchar` when creating HA tables
   [GH-5529]

## 0.11.3 (October 8th, 2018)

SECURITY:

 * Revocation: A regression in 0.11.2 (OSS) and 0.11.0 (Enterprise) caused
   lease IDs containing periods (`.`) to not be revoked properly. Upon startup
   when revocation is tried again these should now revoke successfully.

IMPROVEMENTS:

 * auth/ldap: Listing of users and groups return absolute paths [GH-5537]
 * secret/pki: OID SANs can now specify `*` to allow any value [GH-5459]

BUG FIXES:

 * auth/ldap: Fix panic if specific values were given to be escaped [GH-5471]
 * cli/auth: Fix panic if `vault auth` was given no parameters [GH-5473]
 * secret/database/mongodb: Fix panic that could occur at high load [GH-5463]
 * secret/pki: Fix CA generation not allowing OID SANs [GH-5459]

## 0.11.2 (October 2nd, 2018)

CHANGES:

 * `sys/seal-status` now includes an `initialized` boolean in the output. If
   Vault is not initialized, it will return a `200` with this value set `false`
   instead of a `400`.
 * `passthrough_request_headers` will now deny certain headers from being
   provided to backends based on a global denylist.
 * Token Format: Tokens are now represented as a base62 value; tokens in
   namespaces will have the namespace identifier appended. (This appeared in
   Enterprise in 0.11.0, but is only in OSS in 0.11.2.)

FEATURES:

 * **AWS Secret Engine Root Credential Rotation**: The credential used by the AWS
   secret engine can now be rotated, to ensure that only Vault knows the
   credentials it is using [GH-5140]
 * **Storage Backend Migrator**: A new `operator migrate` command allows offline
   migration of data between two storage backends
 * **AliCloud KMS Auto Unseal and Seal Wrap Support (Enterprise)**: AliCloud KMS can now be used a support seal for
   Auto Unseal and Seal Wrapping

BUG FIXES:

 * auth/okta: Fix reading deprecated `token` parameter if a token was
   previously set in the configuration [GH-5409]
 * core: Re-add deprecated capabilities information for now [GH-5360]
 * core: Fix handling of cyclic token relationships [GH-4803]
 * storage/mysql: Fix locking on MariaDB [GH-5343]
 * replication: Fix DR API when using a token [GH-5398]
 * identity: Ensure old group alias is removed when a new one is written [GH-5350]
 * storage/alicloud: Don't call uname on package init [GH-5358]
 * secrets/jwt: Fix issue where request context would be canceled too early
 * ui: fix need to have update for aws iam creds generation [GF-5294]
 * ui: fix calculation of token expiry [GH-5435]

IMPROVEMENTS:

 * auth/aws: The identity alias name can now configured to be either IAM unique
   ID of the IAM Principal, or ARN of the caller identity [GH-5247]
 * auth/cert: Add allowed_organizational_units support [GH-5252]
 * cli: Format TTLs for non-secret responses [GH-5367]
 * identity: Support operating on entities and groups by their names [GH-5355]
 * plugins: Add `env` parameter when registering plugins to the catalog to allow
   operators to include environment variables during plugin execution. [GH-5359]
 * secrets/aws: WAL Rollback improvements [GH-5202]
 * secrets/aws: Allow specifying STS role-default TTLs [GH-5138]
 * secrets/pki: Add configuration support for setting NotBefore [GH-5325]
 * core: Support for passing the Vault token via an Authorization Bearer header [GH-5397]
 * replication: Reindex process now runs in the background and does not block other
   vault operations
 * storage/zookeeper: Enable TLS based communication with Zookeeper [GH-4856]
 * ui: you can now init a cluster with a seal config [GH-5428]
 * ui: added the option to force promote replication clusters [GH-5438]
 * replication: Allow promotion of a secondary when data is syncing with a "force" flag

## 0.11.1.1 (September 17th, 2018) (Enterprise Only)

BUG FIXES:

 * agent: Fix auth handler-based wrapping of output tokens [GH-5316]
 * core: Properly store the replication checkpoint file if it's larger than the
   storage engine's per-item limit
 * core: Improve WAL deletion rate
 * core: Fix token creation on performance standby nodes
 * core: Fix unwrapping inside a namespace
 * core: Always forward tidy operations from performance standby nodes

IMPROVEMENTS:

 * auth/aws: add support for key/value pairs or JSON values for
   `iam_request_headers` with IAM auth method [GH-5320]
 * auth/aws, secret/aws: Throttling errors from the AWS API will now be
   reported as 502 errors by Vault, along with the original error [GH-5270]
 * replication: Start fetching during a sync from where it previously errored

## 0.11.1 (September 6th, 2018)

SECURITY:

 * Random Byte Reading in Barrier: Prior to this release, Vault was not
   properly checking the error code when reading random bytes for the IV for
   AES operations in its cryptographic barrier. Specifically, this means that
   such an IV could potentially be zero multiple times, causing nonce re-use
   and weakening the security of the key. On most platforms this should never
   happen because reading from kernel random sources is non-blocking and always
   successful, but there may be platform-specific behavior that has not been
   accounted for. (Vault has tests to check exactly this, and the tests have
   never seen nonce re-use.)

FEATURES:

 * AliCloud Agent Support: Vault Agent can now authenticate against the
   AliCloud auth method.
 * UI: Enable AliCloud auth method and Azure secrets engine via the UI.

IMPROVEMENTS:

 * core: Logging level for most logs (not including secrets/auth plugins) can
   now be changed on-the-fly via `SIGHUP`, reading the desired value from
   Vault's config file [GH-5280]

BUG FIXES:

 * core: Ensure we use a background context when stepping down [GH-5290]
 * core: Properly check error return from random byte reading [GH-5277]
 * core: Re-add `sys/` top-route injection for now [GH-5241]
 * core: Policies stored in minified JSON would return an error [GH-5229]
 * core: Evaluate templated policies in capabilities check [GH-5250]
 * identity: Update MemDB with identity group alias while loading groups [GH-5289]
 * secrets/database: Fix nil pointer when revoking some leases [GH-5262]
 * secrets/pki: Fix sign-verbatim losing extra Subject attributes [GH-5245]
 * secrets/pki: Remove certificates from store when tidying revoked
   certificates and simplify API [GH-5231]
 * ui: JSON editor will not coerce input to an object, and will now show an
   error about Vault expecting an object [GH-5271]
 * ui: authentication form will now default to any methods that have been tuned
   to show up for unauthenticated users [GH-5281]


## 0.11.0 (August 28th, 2018)

DEPRECATIONS/CHANGES:

 * Request Timeouts: A default request timeout of 90s is now enforced. This
   setting can be overwritten in the config file. If you anticipate requests
   taking longer than 90s this setting should be updated before upgrading.
 * (NOTE: will be re-added into 0.11.1 as it broke more than anticipated. There
   will be some further guidelines around when this will be removed again.)
   * `sys/` Top Level Injection: For the last two years for backwards
   compatibility data for various `sys/` routes has been injected into both the
   Secret's Data map and into the top level of the JSON response object.
   However, this has some subtle issues that pop up from time to time and is
   becoming increasingly complicated to maintain, so it's finally being
   removed.
 * Path Fallback for List Operations: For a very long time Vault has
   automatically adjusted `list` operations to always end in a `/`, as list
   operations operates on prefixes, so all list operations by definition end
   with `/`. This was done server-side so affects all clients. However, this
   has also led to a lot of confusion for users writing policies that assume
   that the path that they use in the CLI is the path used internally. Starting
   in 0.11, ACL policies gain a new fallback rule for listing: they will use a
   matching path ending in `/` if available, but if not found, they will look
   for the same path without a trailing `/`. This allows putting `list`
   capabilities in the same path block as most other capabilities for that
   path, while not providing any extra access if `list` wasn't actually
   provided there.
 * Performance Standbys On By Default: If you flavor/license of Vault
   Enterprise supports Performance Standbys, they are on by default. You can
   disable this behavior per-node with the `disable_performance_standby`
   configuration flag.
 * AWS Secret Engine Roles: The AWS Secret Engine roles are now explicit about
   the type of AWS credential they are generating; this reduces reduce
   ambiguity that existed previously as well as enables new features for
   specific credential types. Writing role data and generating credentials
   remain backwards compatible; however, the data returned when reading a
   role's configuration has changed in backwards-incompatible ways. Anything
   that depended on reading role data from the AWS secret engine will break
   until it is updated to work with the new format.
 * Token Format (Enterprise): Tokens are now represented as a base62 value;
   tokens in namespaces will have the namespace identifier appended.

FEATURES:

 * **Namespaces (Enterprise)**: A set of features within Vault Enterprise
   that allows Vault environments to support *Secure Multi-tenancy* within a
   single Vault Enterprise infrastructure. Through namespaces, Vault
   administrators can support tenant isolation for teams and individuals as
   well as empower those individuals to self-manage their own tenant
   environment.
 * **Performance Standbys (Enterprise)**: Standby nodes can now service
   requests that do not modify storage. This provides near-horizontal scaling
   of a cluster in some workloads, and is the intra-cluster analogue of
   the existing Performance Replication feature, which replicates to distinct
   clusters in other datacenters, geos, etc.
 * **AliCloud OSS Storage**: AliCloud OSS can now be used for Vault storage.
 * **AliCloud Auth Plugin**: AliCloud's identity services can now be used to
   grant access to Vault. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-auth-alicloud) for
   more information.
 * **Azure Secrets Plugin**: There is now a plugin (pulled in to Vault) that
   allows generating credentials to allow access to Azure. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-secrets-azure) for
   more information.
 * **HA Support for MySQL Storage**: MySQL storage now supports HA.
 * **ACL Templating**: ACL policies can now be templated using identity Entity,
   Groups, and Metadata.
 * **UI Onboarding wizards**: The Vault UI can provide contextual help and
   guidance, linking out to relevant 
