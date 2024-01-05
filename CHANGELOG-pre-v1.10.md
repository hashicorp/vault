## 1.9.10

### September 30, 2022

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]

## 1.9.9

### August 31, 2022

SECURITY:

* core: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. This vulnerability, CVE-2022-40186, is fixed in 1.11.3, 1.10.6, and 1.9.9. [[HSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

CHANGES:

* core: Bump Go version to 1.17.13.

BUG FIXES:

* core (enterprise): Fix some races in merkle index flushing code found in testing
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]

SECURITY:

* identity/entity: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. [[HCSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

## 1.9.8

### July 21, 2022

SECURITY:

* storage/raft: Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

CHANGES:

* core: Bump Go version to 1.17.12.

IMPROVEMENTS:

* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]

BUG FIXES:

* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* transform (enterprise): Fix a bug in the handling of nested or unmatched capture groups in FPE transformations.
* ui: Fix issue where metadata tab is hidden even though policy grants access [[GH-15824](https://github.com/hashicorp/vault/pull/15824)]
* ui: Updated `leasId` to `leaseId` in the "Copy Credentials" section of "Generate AWS Credentials" [[GH-15685](https://github.com/hashicorp/vault/pull/15685)]

## 1.9.7

### June 10, 2022

CHANGES:

* core: Bump Go version to 1.17.11. [[GH-go-ver-197](https://github.com/hashicorp/vault/pull/go-ver-197)]

IMPROVEMENTS:

* ui: Allow namespace param to be parsed from state queryParam [[GH-15378](https://github.com/hashicorp/vault/pull/15378)]

BUG FIXES:

* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* auth/ldap: The logic for setting the entity alias when `username_as_alias` is set
has been fixed. The previous behavior would make a request to the LDAP server to
get `user_attr` before discarding it and using the username instead. This would
make it impossible for a user to connect if this attribute was missing or had
multiple values, even though it would not be used anyway. This has been fixed
and the username is now used without making superfluous LDAP searches. [[GH-15525](https://github.com/hashicorp/vault/pull/15525)]
* core (enterprise): Fix overcounting of lease count quota usage at startup.
* core/config: Only ask the system about network interfaces when address configs contain a template having the format: {{ ... }} [[GH-15224](https://github.com/hashicorp/vault/pull/15224)]
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: Prevent metrics generation from causing deadlocks. [[GH-15693](https://github.com/hashicorp/vault/pull/15693)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* transform (enterprise): Fix non-overridable column default value causing tokenization tokens to expire prematurely when using the MySQL storage backend.
* ui: Fixes client count timezone bug [[GH-15743](https://github.com/hashicorp/vault/pull/15743)]
* ui: Fixes issue logging in with OIDC from a listed auth mounts tab [[GH-15666](https://github.com/hashicorp/vault/pull/15666)]

## 1.9.6

### April 29, 2022

BUG FIXES:

* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]

## 1.9.5

### April 22, 2022

CHANGES:

* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.17.9. [[GH-15045](https://github.com/hashicorp/vault/pull/15045)]

IMPROVEMENTS:

* auth/ldap: Add username_as_alias configurable to change how aliases are named [[GH-14324](https://github.com/hashicorp/vault/pull/14324)]
* core: Systemd unit file included with the Linux packages now sets the service type to notify. [[GH-14385](https://github.com/hashicorp/vault/pull/14385)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer
* website/docs: added a link to an Enigma secret plugin. [[GH-14389](https://github.com/hashicorp/vault/pull/14389)]

BUG FIXES:

* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* api: Respect increment value in grace period calculations in LifetimeWatcher [[GH-14836](https://github.com/hashicorp/vault/pull/14836)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* core (enterprise): Allow local alias create RPCs to persist alias metadata
* core/metrics: Fix incorrect table size metric for local mounts [[GH-14755](https://github.com/hashicorp/vault/pull/14755)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* identity/token: Fixes a bug where duplicate public keys could appear in the .well-known JWKS [[GH-14543](https://github.com/hashicorp/vault/pull/14543)]
* metrics/autosnapshots (enterprise) : Fix bug that could cause
vault.autosnapshots.save.errors to not be incremented when there is an
autosnapshot save error.
* replication (enterprise): fix panic due to missing entity during invalidation of local aliases. [[GH-14622](https://github.com/hashicorp/vault/pull/14622)]
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not excepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fix issue where UI incorrectly handled API errors when mounting backends [[GH-14551](https://github.com/hashicorp/vault/pull/14551)]
* ui: Fixes caching issue on kv new version create [[GH-14489](https://github.com/hashicorp/vault/pull/14489)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging out with wrapped token query parameter [[GH-14329](https://github.com/hashicorp/vault/pull/14329)]
* ui: Fixes issue with correct auth method not selected when logging out from OIDC or JWT methods [[GH-14545](https://github.com/hashicorp/vault/pull/14545)]
* ui: Redirects to managed namespace if incorrect namespace in URL param [[GH-14422](https://github.com/hashicorp/vault/pull/14422)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]

## 1.9.4

### March 3, 2022

SECURITY:

* secrets/pki: Vault and Vault Enterprise (“Vault”) allowed the PKI secrets engine under certain configurations to issue wildcard certificates to authorized users for a specified domain, even if the PKI role policy attribute allow_subdomains is set to false. This vulnerability, CVE-2022-25243, was fixed in Vault 1.8.9 and 1.9.4.
* transform (enterprise): Vault Enterprise (“Vault”) clusters using the tokenization transform feature can expose the tokenization key through the tokenization key configuration endpoint to authorized operators with read permissions on this endpoint. This vulnerability, CVE-2022-25244, was fixed in Vault Enterprise 1.7.10, 1.8.9, and 1.9.4.

CHANGES:

* secrets/azure: Changes the configuration parameter `use_microsoft_graph_api` to use the Microsoft
Graph API by default. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]

IMPROVEMENTS:

* core: Bump Go version to 1.17.7. [[GH-14232](https://github.com/hashicorp/vault/pull/14232)]
* secrets/pki: Restrict issuance of wildcard certificates via role parameter (`allow_wildcard_certificates`) [[GH-14238](https://github.com/hashicorp/vault/pull/14238)]

BUG FIXES:

* Fixed bug where auth method only considers system-identity when multiple identities are available. [#50](https://github.com/hashicorp/vault-plugin-auth-azure/pull/50) [[GH-14138](https://github.com/hashicorp/vault/pull/14138)]
* auth/kubernetes: Properly handle the migration of role storage entries containing an empty `alias_name_source` [[GH-13925](https://github.com/hashicorp/vault/pull/13925)]
* auth/kubernetes: ensure valid entity alias names created for projected volume tokens [[GH-14144](https://github.com/hashicorp/vault/pull/14144)]
* identity/oidc: Adds support for port-agnostic validation of loopback IP redirect URIs. [[GH-13871](https://github.com/hashicorp/vault/pull/13871)]
* identity/oidc: Fixes inherited group membership when evaluating client assignments [[GH-14013](https://github.com/hashicorp/vault/pull/14013)]
* secrets/azure: Fixed bug where Azure environment did not change Graph URL [[GH-13973](https://github.com/hashicorp/vault/pull/13973)]
* secrets/azure: Fixes the [rotate root](https://www.vaultproject.io/api-docs/secret/azure#rotate-root)
operation for upgraded configurations with a `root_password_ttl` of zero. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]
* secrets/gcp: Fixed bug where error was not reported for invalid bindings [[GH-13974](https://github.com/hashicorp/vault/pull/13974)]
* secrets/openldap: Fix panic from nil logger in backend [[GH-14171](https://github.com/hashicorp/vault/pull/14171)]
* secrets/pki: Fix issuance of wildcard certificates matching glob patterns [[GH-14235](https://github.com/hashicorp/vault/pull/14235)]
* storage/raft: Fix issues allowing invalid nodes to become leadership candidates. [[GH-13703](https://github.com/hashicorp/vault/pull/13703)]
* ui: Fix default TTL display and set on database role [[GH-14224](https://github.com/hashicorp/vault/pull/14224)]
* ui: Fix incorrect validity message on transit secrets engine [[GH-14233](https://github.com/hashicorp/vault/pull/14233)]
* ui: Fix kv engine access bug [[GH-13872](https://github.com/hashicorp/vault/pull/13872)]
* ui: Fix issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Trigger background token self-renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.9.3

### January 27, 2022

IMPROVEMENTS:

* auth/kubernetes: Added support for dynamically reloading short-lived tokens for better Kubernetes 1.21+ compatibility [[GH-13698](https://github.com/hashicorp/vault/pull/13698)]
* auth/ldap: Add username to alias metadata [[GH-13669](https://github.com/hashicorp/vault/pull/13669)]
* core/identity: Support updating an alias' `custom_metadata` to be empty. [[GH-13395](https://github.com/hashicorp/vault/pull/13395)]
* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]
* http (enterprise): Serve /sys/license/status endpoint within namespaces

BUG FIXES:

* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13492](https://github.com/hashicorp/vault/pull/13492)]
* cli: Fix using kv patch with older server versions that don't support HTTP PATCH. [[GH-13615](https://github.com/hashicorp/vault/pull/13615)]
* core (enterprise): Workaround AWS CloudHSM v5 SDK issue not allowing read-only sessions
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13476](https://github.com/hashicorp/vault/pull/13476)]
* core: add support for go-sockaddr templates in the top-level cluster_addr field [[GH-13678](https://github.com/hashicorp/vault/pull/13678)]
* identity/oidc: Check for a nil signing key on rotation to prevent panics. [[GH-13716](https://github.com/hashicorp/vault/pull/13716)]
* kmip (enterprise): Fix locate by name operations fail to find key after a rekey operation.
* secrets/database/mssql: Accept a boolean for `contained_db`, rather than just a string. [[GH-13469](https://github.com/hashicorp/vault/pull/13469)]
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13548](https://github.com/hashicorp/vault/pull/13548)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-13759](https://github.com/hashicorp/vault/pull/13759)]
* storage/raft: On linux, use map_populate for bolt files to improve startup time. [[GH-13573](https://github.com/hashicorp/vault/pull/13573)]
* storage/raft: Units for bolt metrics now given in milliseconds instead of nanoseconds [[GH-13749](https://github.com/hashicorp/vault/pull/13749)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]

## 1.9.2

### December 21, 2021

CHANGES:

* go: Update go version to 1.17.5 [[GH-13408](https://github.com/hashicorp/vault/pull/13408)]

IMPROVEMENTS:

* auth/jwt: The Authorization Code flow makes use of the Proof Key for Code Exchange (PKCE) extension. [[GH-13365](https://github.com/hashicorp/vault/pull/13365)]

BUG FIXES:

* ui: Fix client count current month data not showing unless monthly history data exists [[GH-13396](https://github.com/hashicorp/vault/pull/13396)]

## 1.9.1

### December 9, 2021

SECURITY:

* storage/raft: Integrated Storage backend could be caused to crash by an authenticated user with write permissions to the KV secrets engine. This vulnerability, CVE-2021-45042, was fixed in Vault 1.7.7, 1.8.6, and 1.9.1.

IMPROVEMENTS:

* storage/aerospike: Upgrade `aerospike-client-go` to v5.6.0. [[GH-12165](https://github.com/hashicorp/vault/pull/12165)]

BUG FIXES:

* auth/approle: Fix regression where unset cidrlist is returned as nil instead of zero-length array. [[GH-13235](https://github.com/hashicorp/vault/pull/13235)]
* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* http:Fix /sys/monitor endpoint returning streaming not supported [[GH-13200](https://github.com/hashicorp/vault/pull/13200)]
* identity/oidc: Make the `nonce` parameter optional for the Authorization Endpoint of OIDC providers. [[GH-13231](https://github.com/hashicorp/vault/pull/13231)]
* identity: Fixes a panic in the OIDC key rotation due to a missing nil check. [[GH-13298](https://github.com/hashicorp/vault/pull/13298)]
* sdk/queue: move lock before length check to prevent panics. [[GH-13146](https://github.com/hashicorp/vault/pull/13146)]
* secrets/azure: Fixes service principal generation when assigning roles that have [DataActions](https://docs.microsoft.com/en-us/azure/role-based-access-control/role-definitions#dataactions). [[GH-13277](https://github.com/hashicorp/vault/pull/13277)]
* secrets/pki: Recognize ed25519 when requesting a response in PKCS8 format [[GH-13257](https://github.com/hashicorp/vault/pull/13257)]
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* ui: Do not show verify connection value on database connection config page [[GH-13152](https://github.com/hashicorp/vault/pull/13152)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with automate secret deletion value not displaying initially if set in secret metadata edit view [[GH-13177](https://github.com/hashicorp/vault/pull/13177)]
* ui: Fixes issue with placeholder not displaying for automatically deleted secrets when deletion time has passed [[GH-13166](https://github.com/hashicorp/vault/pull/13166)]
* ui: Fixes node-forge error when parsing EC (elliptical curve) certs [[GH-13238](https://github.com/hashicorp/vault/pull/13238)]

## 1.9.0

### November 17, 2021

CHANGES:

* auth/kubernetes: `disable_iss_validation` defaults to true. [#127](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/127) [[GH-12975](https://github.com/hashicorp/vault/pull/12975)]
* expiration: VAULT_16_REVOKE_PERMITPOOL environment variable has been removed. [[GH-12888](https://github.com/hashicorp/vault/pull/12888)]
* expiration: VAULT_LEASE_USE_LEGACY_REVOCATION_STRATEGY environment variable has
been removed. [[GH-12888](https://github.com/hashicorp/vault/pull/12888)]
* go: Update go version to 1.17.2
* secrets/ssh: Roles with empty allowed_extensions will now forbid end-users
specifying extensions when requesting ssh key signing. Update roles setting
allowed_extensions to `*` to permit any extension to be specified by an end-user. [[GH-12847](https://github.com/hashicorp/vault/pull/12847)]

FEATURES:

* **Customizable HTTP Headers**: Add support to define custom HTTP headers for root path (`/`) and also on API endpoints (`/v1/*`) [[GH-12485](https://github.com/hashicorp/vault/pull/12485)]
* **Deduplicate Token With Entities in Activity Log**: Vault tokens without entities are now tracked with client IDs and deduplicated in the Activity Log [[GH-12820](https://github.com/hashicorp/vault/pull/12820)]
* **Elasticsearch Database UI**: The UI now supports adding and editing Elasticsearch connections in the database secret engine. [[GH-12672](https://github.com/hashicorp/vault/pull/12672)]
* **KV Custom Metadata**: Add ability in kv-v2 to specify version-agnostic custom key metadata via the
metadata endpoint. The data will be present in responses made to the data endpoint independent of the
calling token's `read` access to the metadata endpoint. [[GH-12907](https://github.com/hashicorp/vault/pull/12907)]
* **KV patch (Tech Preview)**: Add partial update support for the `/<mount>/data/:path` kv-v2
endpoint through HTTP `PATCH`.  A new `patch` ACL capability has been added and
is required to make such requests. [[GH-12687](https://github.com/hashicorp/vault/pull/12687)]
* **Key Management Secrets Engine (Enterprise)**: Adds support for distributing and managing keys in GCP Cloud KMS.
* **Local Auth Mount Entities (enterprise)**: Logins on `local` auth mounts will
generate identity entities for the tokens issued. The aliases of the entity
resulting from local auth mounts (local-aliases), will be scoped by the cluster.
This means that the local-aliases will never leave the geographical boundary of
the cluster where they were issued. This is something to be mindful about for
those who have implemented local auth mounts for complying with GDPR guidelines.
* **Namespaces (Enterprise)**: Adds support for locking Vault API for particular namespaces.
* **OIDC Identity Provider (Tech Preview)**: Adds support for Vault to be an OpenID Connect (OIDC) provider. [[GH-12932](https://github.com/hashicorp/vault/pull/12932)]
* **Oracle Database UI**: The UI now supports adding and editing Oracle connections in the database secret engine. [[GH-12752](https://github.com/hashicorp/vault/pull/12752)]
* **Postgres Database UI**: The UI now supports adding and editing Postgres connections in the database secret engine. [[GH-12945](https://github.com/hashicorp/vault/pull/12945)]

SECURITY:

* core/identity: A Vault user with write permission to an entity alias ID sharing a mount accessor with another user may acquire this other user’s policies by merging their identities. This vulnerability, CVE-2021-41802, was fixed in Vault and Vault Enterprise 1.7.5, 1.8.4, and 1.9.0.
* core/identity: Templated ACL policies would always match the first-created entity alias if multiple entity aliases existed for a specified entity and mount combination, potentially resulting in incorrect policy enforcement. This vulnerability, CVE-2021-43998, was fixed in Vault and Vault Enterprise 1.7.6, 1.8.5, and 1.9.0.

IMPROVEMENTS:

* agent/cache: Process persistent cache leases in dependency order during restore to ensure child leases are always correctly restored [[GH-12843](https://github.com/hashicorp/vault/pull/12843)]
* agent/cache: Use an in-process listener between consul-template and vault-agent when caching is enabled and either templates or a listener is defined [[GH-12762](https://github.com/hashicorp/vault/pull/12762)]
* agent/cache: tolerate partial restore failure from persistent cache [[GH-12718](https://github.com/hashicorp/vault/pull/12718)]
* agent/template: add support for new 'writeToFile' template function [[GH-12505](https://github.com/hashicorp/vault/pull/12505)]
* api: Add configuration option for ensuring isolated read-after-write semantics for all Client requests. [[GH-12814](https://github.com/hashicorp/vault/pull/12814)]
* api: adds native Login method to Go client module with different auth method interfaces to support easier authentication [[GH-12796](https://github.com/hashicorp/vault/pull/12796)]
* api: Move mergeStates and other required utils from agent to api module [[GH-12731](https://github.com/hashicorp/vault/pull/12731)]
* api: Support VAULT_HTTP_PROXY environment variable to allow overriding the Vault client's HTTP proxy [[GH-12582](https://github.com/hashicorp/vault/pull/12582)]
* auth/approle: The `role/:name/secret-id-accessor/lookup` endpoint now returns a 404 status code when the `secret_id_accessor` cannot be found [[GH-12788](https://github.com/hashicorp/vault/pull/12788)]
* auth/approle: expose secret_id_accessor as WrappedAccessor when creating wrapped secret-id. [[GH-12425](https://github.com/hashicorp/vault/pull/12425)]
* auth/aws: add profile support for AWS credentials when using the AWS auth method [[GH-12621](https://github.com/hashicorp/vault/pull/12621)]
* auth/kubernetes: validate JWT against the provided role on alias look ahead operations [[GH-12688](https://github.com/hashicorp/vault/pull/12688)]
* auth/kubernetes: Add ability to configure entity alias names based on the serviceaccount's namespace and name. [#110](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/110) [#112](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/112) [[GH-12633](https://github.com/hashicorp/vault/pull/12633)]
* auth/ldap: include support for an optional user filter field when searching for users [[GH-11000](https://github.com/hashicorp/vault/pull/11000)]
* auth/oidc: Adds the `skip_browser` CLI option to allow users to skip opening the default browser during the authentication flow. [[GH-12876](https://github.com/hashicorp/vault/pull/12876)]
* auth/okta: Send x-forwarded-for in Okta Push Factor request [[GH-12320](https://github.com/hashicorp/vault/pull/12320)]
* auth/token: Add `allowed_policies_glob` and `disallowed_policies_glob` fields to token roles to allow glob matching of policies [[GH-7277](https://github.com/hashicorp/vault/pull/7277)]
* cli: Operator diagnose now tests for missing or partial telemetry configurations. [[GH-12802](https://github.com/hashicorp/vault/pull/12802)]
* cli: add new http option : -header which enable sending arbitrary headers with the cli [[GH-12508](https://github.com/hashicorp/vault/pull/12508)]
* command: operator generate-root -decode: allow passing encoded token via stdin [[GH-12881](https://github.com/hashicorp/vault/pull/12881)]
* core/token: Return the token_no_default_policy config on token role read if set [[GH-12565](https://github.com/hashicorp/vault/pull/12565)]
* core: Add support for go-sockaddr templated addresses in config. [[GH-9109](https://github.com/hashicorp/vault/pull/9109)]
* core: adds custom_metadata field for aliases [[GH-12502](https://github.com/hashicorp/vault/pull/12502)]
* core: Update Oracle Cloud library to enable seal integration with the uk-gov-london-1 region [[GH-12724](https://github.com/hashicorp/vault/pull/12724)]
* core: Update github.com/ulikunitz/xz to fix security vulnerability GHSA-25xm-hr59-7c27. [[GH-12253](https://github.com/hashicorp/vault/pull/12253)]
* core: Upgrade github.com/gogo/protobuf [[GH-12255](https://github.com/hashicorp/vault/pull/12255)]
* core: build with Go 1.17, and mitigate a breaking change they made that could impact how approle and ssh interpret IPs/CIDRs [[GH-12868](https://github.com/hashicorp/vault/pull/12868)]
* core: observe the client counts broken down by namespace for partial month client count [[GH-12393](https://github.com/hashicorp/vault/pull/12393)]
* core: Artifact builds will now only run on merges to the release branches or to `main`
* core: The [dockerfile](https://github.com/hashicorp/vault/blob/main/Dockerfile) that is used to build the vault docker image available at [hashicorp/vault](https://hub.docker.com/repository/docker/hashicorp/vault) now lives in the root of this repo, and the entrypoint is available under [.release/docker/docker-entrypoint.sh](https://github.com/hashicorp/vault/blob/main/.release/docker/docker-entrypoint.sh)
* core: The vault linux packaging service configs and pre/post install scripts are now available under [.release/linux](https://github.com/hashicorp/vault/blob/main/.release/linux)
* core: Vault linux packages are now available for all supported linux architectures including arm, arm64, 386, and amd64
* db/cassandra: make the connect_timeout config option actually apply to connection timeouts, in addition to non-connection operations [[GH-12903](https://github.com/hashicorp/vault/pull/12903)]
* identity/token: Only return keys from the `.well-known/keys` endpoint that are being used by roles to sign/verify tokens. [[GH-12780](https://github.com/hashicorp/vault/pull/12780)]
* identity: fix issue where Cache-Control header causes stampede of requests for JWKS keys [[GH-12414](https://github.com/hashicorp/vault/pull/12414)]
* physical/etcd: Upgrade etcd3 client to v3.5.0 and etcd2 to v2.305.0. [[GH-11980](https://github.com/hashicorp/vault/pull/11980)]
* pki: adds signature_bits field to customize signature algorithm on CAs and certs signed by Vault [[GH-11245](https://github.com/hashicorp/vault/pull/11245)]
* plugin: update the couchbase gocb version in the couchbase plugin [[GH-12483](https://github.com/hashicorp/vault/pull/12483)]
* replication (enterprise): Add merkle.flushDirty.num_pages_outstanding metric which specifies number of
outstanding dirty pages that were not flushed. [[GH-2093](https://github.com/hashicorp/vault/pull/2093)]
* sdk/framework: The '+' wildcard is now supported for parameterizing unauthenticated paths. [[GH-12668](https://github.com/hashicorp/vault/pull/12668)]
* secrets/aws: Add conditional template that allows custom usernames for both STS and IAM cases [[GH-12185](https://github.com/hashicorp/vault/pull/12185)]
* secrets/azure: Adds support for rotate-root. [#70](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/70) [[GH-13034](https://github.com/hashicorp/vault/pull/13034)]
* secrets/azure: Adds support for using Microsoft Graph API since Azure Active Directory API is being removed in 2022. [#67](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/67) [[GH-12629](https://github.com/hashicorp/vault/pull/12629)]
* secrets/database: Update MSSQL dependency github.com/denisenkom/go-mssqldb to v0.11.0 and include support for contained databases in MSSQL plugin [[GH-12839](https://github.com/hashicorp/vault/pull/12839)]
* secrets/pki: Allow signing of self-issued certs with a different signature algorithm. [[GH-12514](https://github.com/hashicorp/vault/pull/12514)]
* secrets/pki: Use entropy augmentation when available when generating root and intermediate CA key material. [[GH-12559](https://github.com/hashicorp/vault/pull/12559)]
* secrets/pki: select appropriate signature algorithm for ECDSA signature on certificates. [[GH-11216](https://github.com/hashicorp/vault/pull/11216)]
* secrets/pki: Support ed25519 as a key for the pki backend [[GH-11780](https://github.com/hashicorp/vault/pull/11780)]
* secrets/rabbitmq: Update dependency github.com/michaelklishin/rabbit-hole to v2 and resolve UserInfo.tags regression from RabbitMQ v3.9 [[GH-12877](https://github.com/hashicorp/vault/pull/12877)]
* secrets/ssh: Let allowed_users template mix templated and non-templated parts. [[GH-10886](https://github.com/hashicorp/vault/pull/10886)]
* secrets/ssh: Use entropy augmentation when available for generation of the signing key. [[GH-12560](https://github.com/hashicorp/vault/pull/12560)]
* serviceregistration: add `external-source: "vault"` metadata value for Consul registration. [[GH-12163](https://github.com/hashicorp/vault/pull/12163)]
* storage/raft: Best-effort handling of cancelled contexts. [[GH-12162](https://github.com/hashicorp/vault/pull/12162)]
* transform (enterprise): Add advanced features for encoding and decoding for Transform FPE
* transform (enterprise): Add a `reference` field to batch items, and propogate it to the response
* ui: Add KV secret search box when no metadata list access. [[GH-12626](https://github.com/hashicorp/vault/pull/12626)]
* ui: Add custom metadata to KV secret engine and metadata to config [[GH-12169](https://github.com/hashicorp/vault/pull/12169)]
* ui: Creates new StatText component [[GH-12295](https://github.com/hashicorp/vault/pull/12295)]
* ui: client count monthly view [[GH-12554](https://github.com/hashicorp/vault/pull/12554)]
* ui: creates bar chart component for displaying client count data by namespace [[GH-12437](https://github.com/hashicorp/vault/pull/12437)]
* ui: Add creation time to KV 2 version history and version view [[GH-12663](https://github.com/hashicorp/vault/pull/12663)]
* ui: Added resize for JSON editor [[GH-12906](https://github.com/hashicorp/vault/pull/12906)] [[GH-12906](https://github.com/hashicorp/vault/pull/12906)]
* ui: Adds warning about white space in KV secret engine. [[GH-12921](https://github.com/hashicorp/vault/pull/12921)]
* ui: Click to copy database static role last rotation value in tooltip [[GH-12890](https://github.com/hashicorp/vault/pull/12890)]
* ui: Filter DB connection attributes so only relevant attrs POST to backend [[GH-12770](https://github.com/hashicorp/vault/pull/12770)]
* ui: Removes empty rows from DB config views [[GH-12819](https://github.com/hashicorp/vault/pull/12819)]
* ui: Standardizes toolbar presentation of destructive actions [[GH-12895](https://github.com/hashicorp/vault/pull/12895)]
* ui: Updates font for table row value fields [[GH-12908](https://github.com/hashicorp/vault/pull/12908)]
* ui: namespace search in client count views [[GH-12577](https://github.com/hashicorp/vault/pull/12577)]
* ui: parse and display pki cert metadata [[GH-12541](https://github.com/hashicorp/vault/pull/12541)]
* ui: replaces Vault's use of elazarl/go-bindata-assetfs in building the UI with Go's native Embed package [[GH-11208](https://github.com/hashicorp/vault/pull/11208)]
* ui: updated client tracking config view [[GH-12422](https://github.com/hashicorp/vault/pull/12422)]

DEPRECATIONS:

* auth/kubernetes: deprecate `disable_iss_validation` and `issuer` configuration fields [#127](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/127) [[GH-12975](https://github.com/hashicorp/vault/pull/12975)]

BUG FIXES:

* activity log (enterprise): allow partial monthly client count to be accessed from namespaces [[GH-13086](https://github.com/hashicorp/vault/pull/13086)]
* agent: Avoid possible `unexpected fault address` panic when using persistent cache. [[GH-12534](https://github.com/hashicorp/vault/pull/12534)]
* api: Fixes storage APIs returning incorrect error when parsing responses [[GH-12338](https://github.com/hashicorp/vault/pull/12338)]
* auth/aws: Fix ec2 auth on instances that have a cert in their PKCS7 signature [[GH-12519](https://github.com/hashicorp/vault/pull/12519)]
* auth/aws: Fixes ec2 login no longer supporting DSA signature verification [[GH-12340](https://github.com/hashicorp/vault/pull/12340)]
* auth/aws: fix config/rotate-root to store new key [[GH-12715](https://github.com/hashicorp/vault/pull/12715)]
* auth/jwt: Fixes OIDC auth from the Vault UI when using `form_post` as the `oidc_response_mode`. [[GH-12265](https://github.com/hashicorp/vault/pull/12265)]
* cli/api: Providing consistency for the use of comma separated parameters in auth/secret enable/tune [[GH-12126](https://github.com/hashicorp/vault/pull/12126)]
* cli: fixes CLI requests when namespace is both provided as argument and part of the path [[GH-12720](https://github.com/hashicorp/vault/pull/12720)]
* cli: fixes CLI requests when namespace is both provided as argument and part of the path [[GH-12911](https://github.com/hashicorp/vault/pull/12911)]
* cli: vault debug now puts newlines after every captured log line. [[GH-12175](https://github.com/hashicorp/vault/pull/12175)]
* core (enterprise): Allow deletion of stored licenses on DR secondary nodes
* core (enterprise): Disallow autogenerated licenses to be used in diagnose even when config is specified
* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Fix data race during perf standby sealing
* core (enterprise): Fixes reading raft auto-snapshot configuration from performance standby node [[GH-12317](https://github.com/hashicorp/vault/pull/12317)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* core (enterprise): namespace header included in responses, Go client uses it when displaying error messages [[GH-12196](https://github.com/hashicorp/vault/pull/12196)]
* core/api: Fix an arm64 bug converting a negative int to an unsigned int [[GH-12372](https://github.com/hashicorp/vault/pull/12372)]
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13093](https://github.com/hashicorp/vault/pull/13093)]
* core/identity: Cleanup alias in the in-memory entity after an alias deletion by ID [[GH-12834](https://github.com/hashicorp/vault/pull/12834)]
* core/identity: Disallow entity alias creation/update if a conflicting alias exists for the target entity and mount combination [[GH-12747](https://github.com/hashicorp/vault/pull/12747)]
* core: Fix a deadlock on HA leadership transfer [[GH-12691](https://github.com/hashicorp/vault/pull/12691)]
* core: Fix warnings logged on perf standbys re stored versions [[GH-13042](https://github.com/hashicorp/vault/pull/13042)]
* core: fix byte printing for diagnose disk checks [[GH-12229](https://github.com/hashicorp/vault/pull/12229)]
* core: revert some unintentionally downgraded dependencies from 1.9.0-rc1 [[GH-13168](https://github.com/hashicorp/vault/pull/13168)]
* database/couchbase: change default template to truncate username at 128 characters [[GH-12301](https://github.com/hashicorp/vault/pull/12301)]
* database/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* http: removed unpublished true from logical_system path, making openapi spec consistent with documentation [[GH-12713](https://github.com/hashicorp/vault/pull/12713)]
* identity/token: Adds missing call to unlock mutex in key deletion error handling [[GH-12916](https://github.com/hashicorp/vault/pull/12916)]
* identity: Fail alias rename if the resulting (name,accessor) exists already [[GH-12473](https://github.com/hashicorp/vault/pull/12473)]
* identity: Fix a panic on arm64 platform when doing identity I/O. [[GH-12371](https://github.com/hashicorp/vault/pull/12371)]
* identity: Fix regression preventing startup when aliases were created pre-1.9. [[GH-13169](https://github.com/hashicorp/vault/pull/13169)]
* identity: dedup from_entity_ids when merging two entities [[GH-10101](https://github.com/hashicorp/vault/pull/10101)]
* identity: disallow creation of role without a key parameter [[GH-12208](https://github.com/hashicorp/vault/pull/12208)]
* identity: do not allow a role's token_ttl to be longer than the signing key's verification_ttl [[GH-12151](https://github.com/hashicorp/vault/pull/12151)]
* identity: merge associated entity groups when merging entities [[GH-10085](https://github.com/hashicorp/vault/pull/10085)]
* identity: suppress duplicate policies on entities [[GH-12812](https://github.com/hashicorp/vault/pull/12812)]
* kmip (enterprise): Fix handling of custom attributes when servicing GetAttributes requests
* kmip (enterprise): Fix handling of invalid role parameters within various vault api calls
* kmip (enterprise): Forward KMIP register operations to the active node
* license: ignore stored terminated license while autoloading is enabled [[GH-2104](https://github.com/hashicorp/vault/pull/2104)]
* licensing (enterprise): Revert accidental inclusion of the TDE feature from the `prem` build.
* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* pki: Fix regression preventing email addresses being used as a common name within certificates [[GH-12716](https://github.com/hashicorp/vault/pull/12716)]
* plugin/couchbase: Fix an issue in which the locking patterns did not allow parallel requests. [[GH-13033](https://github.com/hashicorp/vault/pull/13033)]
* plugin/snowflake: Fixed bug where plugin would crash on 32 bit systems [[GH-12378](https://github.com/hashicorp/vault/pull/12378)]
* raft (enterprise): Fix panic when updating auto-snapshot config
* replication (enterprise): Fix issue where merkle.flushDirty.num_pages metric is not emitted if number
of dirty pages is 0. [[GH-2093](https://github.com/hashicorp/vault/pull/2093)]
* replication (enterprise): Fix merkle.saveCheckpoint.num_dirty metric to accurately specify the number
of dirty pages in the merkle tree at time of checkpoint creation. [[GH-2093](https://github.com/hashicorp/vault/pull/2093)]
* sdk/database: Fix a DeleteUser error message on the gRPC client. [[GH-12351](https://github.com/hashicorp/vault/pull/12351)]
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/gcp: Fixes a potential panic in the service account policy rollback for rolesets. [[GH-12379](https://github.com/hashicorp/vault/pull/12379)]
* secrets/keymgmt (enterprise): Fix support for Azure Managed HSM Key Vault instances. [[GH-12934](https://github.com/hashicorp/vault/pull/12934)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12600](https://github.com/hashicorp/vault/pull/12600)]
* secrets/transit: Enforce minimum cache size for transit backend and init cache size on transit backend without restart. [[GH-12418](https://github.com/hashicorp/vault/pull/12418)]
* storage/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* storage/raft (enterprise): Ensure that raft autosnapshot backoff retry duration never hits 0s
* storage/raft: Detect incomplete raft snapshots in api.RaftSnapshot(), and thereby in `vault operator raft snapshot save`. [[GH-12388](https://github.com/hashicorp/vault/pull/12388)]
* storage/raft: Fix regression in 1.9.0-rc1 that changed how time is represented in Raft logs; this prevented using a raft db created pre-1.9. [[GH-13165](https://github.com/hashicorp/vault/pull/13165)]
* storage/raft: Support `addr_type=public_v6` in auto-join [[GH-12366](https://github.com/hashicorp/vault/pull/12366)]
* transform (enterprise): Enforce minimum cache size for Transform backend and reset cache size without a restart
* transform (enterprise): Fix an error where the decode response of an expired token is an empty result rather than an error.
* ui: Adds pagination to auth methods list view [[GH-13054](https://github.com/hashicorp/vault/pull/13054)]
* ui: Fix bug where capabilities check on secret-delete-menu was encoding the forward slashes. [[GH-12550](https://github.com/hashicorp/vault/pull/12550)]
* ui: Fix bug where edit role form on auth method is invalid by default [[GH-12646](https://github.com/hashicorp/vault/pull/12646)]
* ui: Fixed api explorer routing bug [[GH-12354](https://github.com/hashicorp/vault/pull/12354)]
* ui: Fixed text overflow in flash messages [[GH-12357](https://github.com/hashicorp/vault/pull/12357)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: Remove spinner after token renew [[GH-12887](https://github.com/hashicorp/vault/pull/12887)]
* ui: Removes ability to tune token_type for token auth methods [[GH-12904](https://github.com/hashicorp/vault/pull/12904)]
* ui: Show day of month instead of day of year in the expiration warning dialog [[GH-11984](https://github.com/hashicorp/vault/pull/11984)]
* ui: fix issue where on MaskedInput on auth methods if tab it would clear the value. [[GH-12409](https://github.com/hashicorp/vault/pull/12409)]
* ui: fix missing navbar items on login to namespace [[GH-12478](https://github.com/hashicorp/vault/pull/12478)]
* ui: update bar chart when model changes [[GH-12622](https://github.com/hashicorp/vault/pull/12622)]
* ui: updating database TTL picker help text. [[GH-12212](https://github.com/hashicorp/vault/pull/12212)]

## 1.8.12

### June 10, 2022

BUG FIXES:

* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* transform (enterprise): Fix non-overridable column default value causing tokenization tokens to expire prematurely when using the MySQL storage backend.

## 1.8.11

### April 29, 2022

BUG FIXES:

* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]

## 1.8.10

### April 22, 2022

CHANGES:

* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.16.15. [[GH-go-ver-1810](https://github.com/hashicorp/vault/pull/go-ver-1810)]

IMPROVEMENTS:

* auth/ldap: Add username_as_alias configurable to change how aliases are named [[GH-14324](https://github.com/hashicorp/vault/pull/14324)]
* core: Systemd unit file included with the Linux packages now sets the service type to notify. [[GH-14385](https://github.com/hashicorp/vault/pull/14385)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer

BUG FIXES:

* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* metrics/autosnapshots (enterprise) : Fix bug that could cause
vault.autosnapshots.save.errors to not be incremented when there is an
autosnapshot save error.
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not excepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging out with wrapped token query parameter [[GH-14329](https://github.com/hashicorp/vault/pull/14329)]
* ui: Fixes issue with correct auth method not selected when logging out from OIDC or JWT methods [[GH-14545](https://github.com/hashicorp/vault/pull/14545)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]

## 1.8.9

### March 3, 2022

* secrets/pki: Vault and Vault Enterprise (“Vault”) allowed the PKI secrets engine under certain configurations to issue wildcard certificates to authorized users for a specified domain, even if the PKI role policy attribute allow_subdomains is set to false. This vulnerability, CVE-2022-25243, was fixed in Vault 1.8.9 and 1.9.4.
* transform (enterprise): Vault Enterprise (“Vault”) clusters using the tokenization transform feature can expose the tokenization key through the tokenization key configuration endpoint to authorized operators with read permissions on this endpoint. This vulnerability, CVE-2022-25244, was fixed in Vault Enterprise 1.7.10, 1.8.9, and 1.9.4.

IMPROVEMENTS:

* secrets/pki: Restrict issuance of wildcard certificates via role parameter (`allow_wildcard_certificates`) [[GH-14238](https://github.com/hashicorp/vault/pull/14238)]

BUG FIXES:

* auth/aws: Fix ec2 auth on instances that have a cert in their PKCS7 signature [[GH-12519](https://github.com/hashicorp/vault/pull/12519)]
* database/mssql: Removed string interpolation on internal queries and replaced them with inline queries using named parameters. [[GH-13799](https://github.com/hashicorp/vault/pull/13799)]
* secrets/openldap: Fix panic from nil logger in backend [[GH-14170](https://github.com/hashicorp/vault/pull/14170)]
* secrets/pki: Fix issuance of wildcard certificates matching glob patterns [[GH-14235](https://github.com/hashicorp/vault/pull/14235)]
* ui: Fix issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Trigger background token self-renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.8.8

### January 27, 2022

IMPROVEMENTS:

* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]

BUG FIXES:

* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13494](https://github.com/hashicorp/vault/pull/13494)]
* core (enterprise): Workaround AWS CloudHSM v5 SDK issue not allowing read-only sessions
* kmip (enterprise): Fix locate by name operations fail to find key after a rekey operation.
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13549](https://github.com/hashicorp/vault/pull/13549)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-13759](https://github.com/hashicorp/vault/pull/13759)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-2456](https://github.com/hashicorp/vault/pull/2456)]
* storage/raft: Fix issues allowing invalid nodes to become leadership candidates. [[GH-13703](https://github.com/hashicorp/vault/pull/13703)]
* storage/raft: On linux, use map_populate for bolt files to improve startup time. [[GH-13573](https://github.com/hashicorp/vault/pull/13573)]
* storage/raft: Units for bolt metrics now given in milliseconds instead of nanoseconds [[GH-13749](https://github.com/hashicorp/vault/pull/13749)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]

## 1.8.7

### December 21, 2021

CHANGES:

* go: Update go version to 1.16.12 [[GH-13422](https://github.com/hashicorp/vault/pull/13422)]

## 1.8.6

### December 9, 2021

CHANGES:

* go: Update go version to 1.16.9 [[GH-13029](https://github.com/hashicorp/vault/pull/13029)]

SECURITY:

* storage/raft: Integrated Storage backend could be caused to crash by an authenticated user with write permissions to the KV secrets engine. This vulnerability, CVE-2021-45042, was fixed in Vault 1.7.7, 1.8.6, and 1.9.1.

BUG FIXES:

* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* ui: Adds pagination to auth methods list view [[GH-13054](https://github.com/hashicorp/vault/pull/13054)]
* ui: Do not show verify connection value on database connection config page [[GH-13152](https://github.com/hashicorp/vault/pull/13152)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]

## 1.8.5

### November 4, 2021

SECURITY:

* core/identity: Templated ACL policies would always match the first-created entity alias if multiple entity aliases existed for a specified entity and mount combination, potentially resulting in incorrect policy enforcement. This vulnerability, CVE-2021-43998, was fixed in Vault and Vault Enterprise 1.7.6, 1.8.5, and 1.9.0.

BUG FIXES:

* auth/aws: fix config/rotate-root to store new key [[GH-12715](https://github.com/hashicorp/vault/pull/12715)]
* core/identity: Cleanup alias in the in-memory entity after an alias deletion by ID [[GH-12834](https://github.com/hashicorp/vault/pull/12834)]
* core/identity: Disallow entity alias creation/update if a conflicting alias exists for the target entity and mount combination [[GH-12747](https://github.com/hashicorp/vault/pull/12747)]
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* identity/token: Adds missing call to unlock mutex in key deletion error handling [[GH-12916](https://github.com/hashicorp/vault/pull/12916)]
* kmip (enterprise): Fix handling of custom attributes when servicing GetAttributes requests
* kmip (enterprise): Fix handling of invalid role parameters within various vault api calls
* kmip (enterprise): Forward KMIP register operations to the active node
* secrets/keymgmt (enterprise): Fix support for Azure Managed HSM Key Vault instances. [[GH-12952](https://github.com/hashicorp/vault/pull/12952)]
* transform (enterprise): Fix an error where the decode response of an expired token is an empty result rather than an error.

## 1.8.4

### 6 October 2021

SECURITY:

* core/identity: A Vault user with write permission to an entity alias ID sharing a mount accessor with another user may acquire this other user’s policies by merging their identities. This vulnerability, CVE-2021-41802, was fixed in Vault and Vault Enterprise 1.7.5 and 1.8.4.

IMPROVEMENTS:

* core: Update Oracle Cloud library to enable seal integration with the uk-gov-london-1 region [[GH-12724](https://github.com/hashicorp/vault/pull/12724)]

BUG FIXES:

* core: Fix a deadlock on HA leadership transfer [[GH-12691](https://github.com/hashicorp/vault/pull/12691)]
* database/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* pki: Fix regression preventing email addresses being used as a common name within certificates [[GH-12716](https://github.com/hashicorp/vault/pull/12716)]
* storage/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* ui: Fix bug where edit role form on auth method is invalid by default [[GH-12646](https://github.com/hashicorp/vault/pull/12646)]

## 1.8.3

### 29 September 2021

IMPROVEMENTS:

* secrets/pki: Allow signing of self-issued certs with a different signature algorithm. [[GH-12514](https://github.com/hashicorp/vault/pull/12514)]

BUG FIXES:

* agent: Avoid possible `unexpected fault address` panic when using persistent cache. [[GH-12534](https://github.com/hashicorp/vault/pull/12534)]
* core (enterprise): Allow deletion of stored licenses on DR secondary nodes
* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* identity: Fail alias rename if the resulting (name,accessor) exists already [[GH-12473](https://github.com/hashicorp/vault/pull/12473)]
* raft (enterprise): Fix panic when updating auto-snapshot config
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12599](https://github.com/hashicorp/vault/pull/12599)]
* secrets/transit: Enforce minimum cache size for transit backend and init cache size on transit backend without restart. [[GH-12418](https://github.com/hashicorp/vault/pull/12418)]
* storage/raft: Detect incomplete raft snapshots in api.RaftSnapshot(), and thereby in `vault operator raft snapshot save`. [[GH-12388](https://github.com/hashicorp/vault/pull/12388)]
* ui: Fix bug where capabilities check on secret-delete-menu was encoding the forward slashes. [[GH-12550](https://github.com/hashicorp/vault/pull/12550)]
* ui: Show day of month instead of day of year in the expiration warning dialog [[GH-11984](https://github.com/hashicorp/vault/pull/11984)]

## 1.8.2

### 26 August 2021

CHANGES:

* Alpine: Docker images for Vault 1.6.6+, 1.7.4+, and 1.8.2+ are built with Alpine 3.14, due to CVE-2021-36159
* go: Update go version to 1.16.7 [[GH-12408](https://github.com/hashicorp/vault/pull/12408)]

BUG FIXES:

* auth/aws: Fixes ec2 login no longer supporting DSA signature verification [[GH-12340](https://github.com/hashicorp/vault/pull/12340)]
* cli: vault debug now puts newlines after every captured log line. [[GH-12175](https://github.com/hashicorp/vault/pull/12175)]
* database/couchbase: change default template to truncate username at 128 characters [[GH-12300](https://github.com/hashicorp/vault/pull/12300)]
* identity: Fix a panic on arm64 platform when doing identity I/O. [[GH-12371](https://github.com/hashicorp/vault/pull/12371)]
* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* plugin/snowflake: Fixed bug where plugin would crash on 32 bit systems [[GH-12378](https://github.com/hashicorp/vault/pull/12378)]
* sdk/database: Fix a DeleteUser error message on the gRPC client. [[GH-12351](https://github.com/hashicorp/vault/pull/12351)]
* secrets/gcp: Fixes a potential panic in the service account policy rollback for rolesets. [[GH-12379](https://github.com/hashicorp/vault/pull/12379)]
* ui: Fixed api explorer routing bug [[GH-12354](https://github.com/hashicorp/vault/pull/12354)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: fix issue where on MaskedInput on auth methods if tab it would clear the value. [[GH-12409](https://github.com/hashicorp/vault/pull/12409)]

## 1.8.1

### August 5th, 2021

CHANGES:

* go: Update go version to 1.16.6 [[GH-12245](https://github.com/hashicorp/vault/pull/12245)]

IMPROVEMENTS:

* serviceregistration: add `external-source: "vault"` metadata value for Consul registration. [[GH-12163](https://github.com/hashicorp/vault/pull/12163)]

BUG FIXES:

* auth/aws: Remove warning stating AWS Token TTL will be capped by the Default Lease TTL. [[GH-12026](https://github.com/hashicorp/vault/pull/12026)]
* auth/jwt: Fixes OIDC auth from the Vault UI when using `form_post` as the `oidc_response_mode`. [[GH-12258](https://github.com/hashicorp/vault/pull/12258)]
* core (enterprise): Disallow autogenerated licenses to be used in diagnose even when config is specified
* core: fix byte printing for diagnose disk checks [[GH-12229](https://github.com/hashicorp/vault/pull/12229)]
* identity: do not allow a role's token_ttl to be longer than the signing key's verification_ttl [[GH-12151](https://github.com/hashicorp/vault/pull/12151)]

## 1.8.0

### July 28th, 2021

CHANGES:

* agent: Errors in the template engine will no longer cause agent to exit unless
explicitly defined to do so. A new configuration parameter,
`exit_on_retry_failure`, within the new top-level stanza, `template_config`, can
be set to `true` in order to cause agent to exit. Note that for agent to exit if
`template.error_on_missing_key` is set to `true`, `exit_on_retry_failure` must
be also set to `true`. Otherwise, the template engine will log an error but then
restart its internal runner. [[GH-11775](https://github.com/hashicorp/vault/pull/11775)]
* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* core (enterprise): License/EULA changes that ensure the presence of a valid HashiCorp license to
start Vault. More information is available in the [Vault License FAQ](https://www.vaultproject.io/docs/enterprise/license/faqs)

FEATURES:

* **GCP Secrets Engine Static Accounts**: Adds ability to use existing service accounts for generation
  of service account keys and access tokens. [[GH-12023](https://github.com/hashicorp/vault/pull/12023)]
* **Key Management Secrets Engine (Enterprise)**: Adds general availability for distributing and managing keys in AWS KMS. [[GH-11958](https://github.com/hashicorp/vault/pull/11958)]
* **License Autoloading (Enterprise)**: Licenses may now be automatically loaded from the environment or disk.
* **MySQL Database UI**: The UI now supports adding and editing MySQL connections in the database secret engine [[GH-11532](https://github.com/hashicorp/vault/pull/11532)]
* **Vault Diagnose**: A new `vault operator` command to detect common issues with vault server setups.

SECURITY:

* storage/raft: When initializing Vault’s Integrated Storage backend, excessively broad filesystem permissions may be set for the underlying Bolt database used by Vault’s Raft implementation. This vulnerability, CVE-2021-38553, was fixed in Vault 1.8.0.
* ui: The Vault UI erroneously cached and exposed user-viewed secrets between authenticated sessions in a single shared browser, if the browser window / tab was not refreshed or closed between logout and a subsequent login. This vulnerability, CVE-2021-38554, was fixed in Vault 1.8.0 and will be addressed in pending 1.7.4 / 1.6.6 releases.

IMPROVEMENTS:

* agent/template: Added static_secret_render_interval to specify how often to fetch non-leased secrets [[GH-11934](https://github.com/hashicorp/vault/pull/11934)]
* agent: Allow Agent auto auth to read symlinked JWT files [[GH-11502](https://github.com/hashicorp/vault/pull/11502)]
* api: Allow a leveled logger to be provided to `api.Client` through `SetLogger`. [[GH-11696](https://github.com/hashicorp/vault/pull/11696)]
* auth/aws: Underlying error included in validation failure message. [[GH-11638](https://github.com/hashicorp/vault/pull/11638)]
* cli/api: Add lease lookup command [[GH-11129](https://github.com/hashicorp/vault/pull/11129)]
* core: Add `prefix_filter` to telemetry config [[GH-12025](https://github.com/hashicorp/vault/pull/12025)]
* core: Add a darwin/arm64 binary release supporting the Apple M1 CPU [[GH-12071](https://github.com/hashicorp/vault/pull/12071)]
* core: Add a small (<1s) exponential backoff to failed TCP listener Accept failures. [[GH-11588](https://github.com/hashicorp/vault/pull/11588)]
* core (enterprise): Add controlled capabilities to control group policy stanza
* core: Add metrics for standby node forwarding. [[GH-11366](https://github.com/hashicorp/vault/pull/11366)]
* core: Add metrics to report if a node is a perf standby, if a node is a dr secondary or primary, and if a node is a perf secondary or primary. [[GH-11472](https://github.com/hashicorp/vault/pull/11472)]
* core: Send notifications to systemd on start, stop, and configuration reload. [[GH-11517](https://github.com/hashicorp/vault/pull/11517)]
* core: add irrevocable lease list and count apis [[GH-11607](https://github.com/hashicorp/vault/pull/11607)]
* core: allow arbitrary length stack traces upon receiving SIGUSR2 (was 32MB) [[GH-11364](https://github.com/hashicorp/vault/pull/11364)]
* core: Improve renew/revoke performance using per-lease locks [[GH-11122](https://github.com/hashicorp/vault/pull/11122)]
* db/cassandra: Added tls_server_name to specify server name for TLS validation [[GH-11820](https://github.com/hashicorp/vault/pull/11820)]
* go: Update to Go 1.16.5 [[GH-11802](https://github.com/hashicorp/vault/pull/11802)]
* replication: Delay evaluation of X-Vault-Index headers until merkle sync completes.
* secrets/rabbitmq: Add ability to customize dynamic usernames [[GH-11899](https://github.com/hashicorp/vault/pull/11899)]
* secrets/ad: Add `rotate-role` endpoint to allow rotations of service accounts. [[GH-11942](https://github.com/hashicorp/vault/pull/11942)]
* secrets/aws: add IAM tagging support for iam_user roles [[GH-10953](https://github.com/hashicorp/vault/pull/10953)]
* secrets/aws: add ability to provide a role session name when generating STS credentials [[GH-11345](https://github.com/hashicorp/vault/pull/11345)]
* secrets/database/elasticsearch: Add ability to customize dynamic usernames [[GH-11957](https://github.com/hashicorp/vault/pull/11957)]
* secrets/database/influxdb: Add ability to customize dynamic usernames [[GH-11796](https://github.com/hashicorp/vault/pull/11796)]
* secrets/database/mongodb: Add ability to customize `SocketTimeout`, `ConnectTimeout`, and `ServerSelectionTimeout` [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database/mongodb: Increased throughput by allowing for multiple request threads to simultaneously update users in MongoDB [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database/mongodbatlas: Adds the ability to customize username generation for dynamic users in MongoDB Atlas. [[GH-11956](https://github.com/hashicorp/vault/pull/11956)]
* secrets/database/redshift: Add ability to customize dynamic usernames [[GH-12016](https://github.com/hashicorp/vault/pull/12016)]
* secrets/database/snowflake: Add ability to customize dynamic usernames [[GH-11997](https://github.com/hashicorp/vault/pull/11997)]
* ssh: add support for templated values in SSH CA DefaultExtensions [[GH-11495](https://github.com/hashicorp/vault/pull/11495)]
* storage/raft: Improve raft batch size selection [[GH-11907](https://github.com/hashicorp/vault/pull/11907)]
* storage/raft: change freelist type to map and set nofreelistsync to true [[GH-11895](https://github.com/hashicorp/vault/pull/11895)]
* storage/raft: Switch to shared raft-boltdb library and add boltdb metrics [[GH-11269](https://github.com/hashicorp/vault/pull/11269)]
* storage/raft: Support autopilot for HA only raft storage. [[GH-11260](https://github.com/hashicorp/vault/pull/11260)]
* storage/raft (enterprise): Enable Autopilot on DR secondary clusters
* ui: Add Validation to KV secret engine [[GH-11785](https://github.com/hashicorp/vault/pull/11785)]
* ui: Add database secret engine support for MSSQL [[GH-11231](https://github.com/hashicorp/vault/pull/11231)]
* ui: Add push notification message when selecting okta auth. [[GH-11442](https://github.com/hashicorp/vault/pull/11442)]
* ui: Add regex validation to Transform Template pattern input [[GH-11586](https://github.com/hashicorp/vault/pull/11586)]
* ui: Add specific error message if unseal fails due to license [[GH-11705](https://github.com/hashicorp/vault/pull/11705)]
* ui: Add validation support for open api form fields [[GH-11963](https://github.com/hashicorp/vault/pull/11963)]
* ui: Added auth method descriptions to UI login page [[GH-11795](https://github.com/hashicorp/vault/pull/11795)]
* ui: JSON fields on database can be cleared on edit [[GH-11708](https://github.com/hashicorp/vault/pull/11708)]
* ui: Obscure secret values on input and displayOnly fields like certificates. [[GH-11284](https://github.com/hashicorp/vault/pull/11284)]
* ui: Redesign of KV 2 Delete toolbar. [[GH-11530](https://github.com/hashicorp/vault/pull/11530)]
* ui: Replace tool partials with components. [[GH-11672](https://github.com/hashicorp/vault/pull/11672)]
* ui: Show description on secret engine list [[GH-11995](https://github.com/hashicorp/vault/pull/11995)]
* ui: Update ember to latest LTS and upgrade UI dependencies [[GH-11447](https://github.com/hashicorp/vault/pull/11447)]
* ui: Update partials to components [[GH-11680](https://github.com/hashicorp/vault/pull/11680)]
* ui: Updated ivy code mirror component for consistency [[GH-11500](https://github.com/hashicorp/vault/pull/11500)]
* ui: Updated node to v14, latest stable build [[GH-12049](https://github.com/hashicorp/vault/pull/12049)]
* ui: Updated search select component styling [[GH-11360](https://github.com/hashicorp/vault/pull/11360)]
* ui: add transform secrets engine to features list [[GH-12003](https://github.com/hashicorp/vault/pull/12003)]
* ui: add validations for duplicate path kv engine [[GH-11878](https://github.com/hashicorp/vault/pull/11878)]
* ui: show site-wide banners for license warnings if applicable [[GH-11759](https://github.com/hashicorp/vault/pull/11759)]
* ui: update license page with relevant autoload info [[GH-11778](https://github.com/hashicorp/vault/pull/11778)]

DEPRECATIONS:

* secrets/gcp: Deprecated the `/gcp/token/:roleset` and `/gcp/key/:roleset` paths for generating
  secrets for rolesets. Use `/gcp/roleset/:roleset/token` and `/gcp/roleset/:roleset/key` instead. [[GH-12023](https://github.com/hashicorp/vault/pull/12023)]

BUG FIXES:

* activity: Omit wrapping tokens and control groups from client counts [[GH-11826](https://github.com/hashicorp/vault/pull/11826)]
* agent/cert: Fix issue where the API client on agent was not honoring certificate
  information from the auto-auth config map on renewals or retries. [[GH-11576](https://github.com/hashicorp/vault/pull/11576)]
* agent/template: fix command shell quoting issue [[GH-11838](https://github.com/hashicorp/vault/pull/11838)]
* agent: Fixed agent templating to use configured tls servername values [[GH-11288](https://github.com/hashicorp/vault/pull/11288)]
* agent: fix timestamp format in log messages from the templating engine [[GH-11838](https://github.com/hashicorp/vault/pull/11838)]
* auth/approle: fixing dereference of nil pointer [[GH-11864](https://github.com/hashicorp/vault/pull/11864)]
* auth/jwt: Updates the [hashicorp/cap](https://github.com/hashicorp/cap) library to `v0.1.0` to
  bring in a verification key caching fix. [[GH-11784](https://github.com/hashicorp/vault/pull/11784)]
* auth/kubernetes: Fix AliasLookahead to correctly extract ServiceAccount UID when using ephemeral JWTs [[GH-12073](https://github.com/hashicorp/vault/pull/12073)]
* auth/ldap: Fix a bug where the LDAP auth method does not return the request_timeout configuration parameter on config read. [[GH-11975](https://github.com/hashicorp/vault/pull/11975)]
* cli: Add support for response wrapping in `vault list` and `vault kv list` with output format other than `table`. [[GH-12031](https://github.com/hashicorp/vault/pull/12031)]
* cli: vault delete and vault kv delete should support the same output options (e.g. -format) as vault write. [[GH-11992](https://github.com/hashicorp/vault/pull/11992)]
* core (enterprise): Fix orphan return value from auth methods executed on performance standby nodes.
* core (enterprise): Fix plugins mounted in namespaces being unable to use password policies [[GH-11596](https://github.com/hashicorp/vault/pull/11596)]
* core (enterprise): serialize access to HSM entropy generation to avoid errors in concurrent key generation.
* core/metrics: Add generic KV mount support for vault.kv.secret.count telemetry metric [[GH-12020](https://github.com/hashicorp/vault/pull/12020)]
* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* core: Fix edge cases in the configuration endpoint for barrier key autorotation. [[GH-11541](https://github.com/hashicorp/vault/pull/11541)]
* core: Fix goroutine leak when updating rate limit quota [[GH-11371](https://github.com/hashicorp/vault/pull/11371)]
* core (enterprise): Fix panic on DR secondary when there are lease count quotas [[GH-11742](https://github.com/hashicorp/vault/pull/11742)]
* core: Fix race that allowed remounting on path used by another mount [[GH-11453](https://github.com/hashicorp/vault/pull/11453)]
* core: Fix storage entry leak when revoking leases created with non-orphan batch tokens. [[GH-11377](https://github.com/hashicorp/vault/pull/11377)]
* core: Fixed double counting of http requests after operator stepdown [[GH-11970](https://github.com/hashicorp/vault/pull/11970)]
* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]
* identity: Use correct mount accessor when refreshing external group memberships. [[GH-11506](https://github.com/hashicorp/vault/pull/11506)]
* mongo-db: default username template now strips invalid '.' characters [[GH-11872](https://github.com/hashicorp/vault/pull/11872)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* replication: Fix panic trying to update walState during identity group invalidation.
* replication: Fix: mounts created within a namespace that was part of an Allow
  filtering rule would not appear on performance secondary if created after rule
  was defined.
* secret/pki: use case insensitive domain name comparison as per RFC1035 section 2.3.3
* secret: fix the bug where transit encrypt batch doesn't work with key_version [[GH-11628](https://github.com/hashicorp/vault/pull/11628)]
* secrets/ad: Forward all creds requests to active node [[GH-76](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/76)] [[GH-11836](https://github.com/hashicorp/vault/pull/11836)]
* secrets/database/cassandra: Fixed issue where hostnames were not being validated when using TLS [[GH-11365](https://github.com/hashicorp/vault/pull/11365)]
* secrets/database/cassandra: Fixed issue where the PEM parsing logic of `pem_bundle` and `pem_json` didn't work for CA-only configurations [[GH-11861](https://github.com/hashicorp/vault/pull/11861)]
* secrets/database/cassandra: Updated default statement for password rotation to allow for special characters. This applies to root and static credentials. [[GH-11262](https://github.com/hashicorp/vault/pull/11262)]
* secrets/database: Fix marshalling to allow providing numeric arguments to external database plugins. [[GH-11451](https://github.com/hashicorp/vault/pull/11451)]
* secrets/database: Fixed an issue that prevented external database plugin processes from restarting after a shutdown. [[GH-12087](https://github.com/hashicorp/vault/pull/12087)]
* secrets/database: Fixed minor race condition when rotate-root is called [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database: Fixes issue for V4 database interface where `SetCredentials` wasn't falling back to using `RotateRootCredentials` if `SetCredentials` is `Unimplemented` [[GH-11585](https://github.com/hashicorp/vault/pull/11585)]
* secrets/openldap: Fix bug where schema was not compatible with rotate-root [#24](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/24) [[GH-12019](https://github.com/hashicorp/vault/pull/12019)]
* storage/dynamodb: Handle throttled batch write requests by retrying, without which writes could be lost. [[GH-10181](https://github.com/hashicorp/vault/pull/10181)]
* storage/raft: Support cluster address change for nodes in a cluster managed by autopilot [[GH-11247](https://github.com/hashicorp/vault/pull/11247)]
* storage/raft: Tweak creation of vault.db file [[GH-12034](https://github.com/hashicorp/vault/pull/12034)]
* storage/raft: leader_tls_servername wasn't used unless leader_ca_cert_file and/or mTLS were configured. [[GH-11252](https://github.com/hashicorp/vault/pull/11252)]
* tokenutil: Perform the num uses check before token type. [[GH-11647](https://github.com/hashicorp/vault/pull/11647)]
* transform (enterprise): Fix an issue with malformed transform configuration
  storage when upgrading from 1.5 to 1.6.  See Upgrade Notes for 1.6.x.
* ui: Add role from database connection automatically populates the database for new role [[GH-11119](https://github.com/hashicorp/vault/pull/11119)]
* ui: Add root rotation statements support to appropriate database secret engine plugins [[GH-11404](https://github.com/hashicorp/vault/pull/11404)]
* ui: Automatically refresh the page when user logs out [[GH-12035](https://github.com/hashicorp/vault/pull/12035)]
* ui: Fix Version History queryParams on LinkedBlock [[GH-12079](https://github.com/hashicorp/vault/pull/12079)]
* ui: Fix bug where database secret engines with custom names cannot delete connections [[GH-11127](https://github.com/hashicorp/vault/pull/11127)]
* ui: Fix bug where the UI does not recognize version 2 KV until refresh, and fix [object Object] error message [[GH-11258](https://github.com/hashicorp/vault/pull/11258)]
* ui: Fix database role CG access [[GH-12111](https://github.com/hashicorp/vault/pull/12111)]
* ui: Fix date display on expired token notice [[GH-11142](https://github.com/hashicorp/vault/pull/11142)]
* ui: Fix entity group membership and metadata not showing [[GH-11641](https://github.com/hashicorp/vault/pull/11641)]
* ui: Fix error message caused by control group [[GH-11143](https://github.com/hashicorp/vault/pull/11143)]
* ui: Fix footer URL linking to the correct version changelog. [[GH-11283](https://github.com/hashicorp/vault/pull/11283)]
* ui: Fix issue where logging in without namespace input causes error [[GH-11094](https://github.com/hashicorp/vault/pull/11094)]
* ui: Fix namespace-bug on login [[GH-11182](https://github.com/hashicorp/vault/pull/11182)]
* ui: Fix status menu no showing on login [[GH-11213](https://github.com/hashicorp/vault/pull/11213)]
* ui: Fix text link URL on database roles list [[GH-11597](https://github.com/hashicorp/vault/pull/11597)]
* ui: Fixed and updated lease renewal picker [[GH-11256](https://github.com/hashicorp/vault/pull/11256)]
* ui: fix control group access for database credential [[GH-12024](https://github.com/hashicorp/vault/pull/12024)]
* ui: fix issue where select-one option was not showing in secrets database role creation [[GH-11294](https://github.com/hashicorp/vault/pull/11294)]
* ui: fix oidc login with Safari [[GH-11884](https://github.com/hashicorp/vault/pull/11884)]

## 1.7.10

### March 3, 2022

SECURITY:

* transform (enterprise): Vault Enterprise (“Vault”) clusters using the tokenization transform feature can expose the tokenization key through the tokenization key configuration endpoint to authorized operators with read permissions on this endpoint. This vulnerability, CVE-2022-25244, was fixed in Vault Enterprise 1.7.10, 1.8.9, and 1.9.4.

BUG FIXES:

* database/mssql: Removed string interpolation on internal queries and replaced them with inline queries using named parameters. [[GH-13799](https://github.com/hashicorp/vault/pull/13799)]
* ui: Fix issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Trigger background token self-renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.7.9

### January 27, 2022

IMPROVEMENTS:

* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]

BUG FIXES:

* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13493](https://github.com/hashicorp/vault/pull/13493)]
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13735](https://github.com/hashicorp/vault/pull/13735)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]

## 1.7.8

### December 21, 2021

CHANGES:

* go: Update go version to 1.16.12 [[GH-13422](https://github.com/hashicorp/vault/pull/13422)]

BUG FIXES:

* auth/aws: Fixes ec2 login no longer supporting DSA signature verification [[GH-12340](https://github.com/hashicorp/vault/pull/12340)]
* identity: Fix a panic on arm64 platform when doing identity I/O. [[GH-12371](https://github.com/hashicorp/vault/pull/12371)]

## 1.7.7

### December 9, 2021

SECURITY:

* storage/raft: Integrated Storage backend could be caused to crash by an authenticated user with write permissions to the KV secrets engine. This vulnerability, CVE-2021-45042, was fixed in Vault 1.7.7, 1.8.6, and 1.9.1.

BUG FIXES:

* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]

## 1.7.6

### November 4, 2021

SECURITY:

* core/identity: Templated ACL policies would always match the first-created entity alias if multiple entity aliases existed for a specified entity and mount combination, potentially resulting in incorrect policy enforcement. This vulnerability, CVE-2021-43998, was fixed in Vault and Vault Enterprise 1.7.6, 1.8.5, and 1.9.0.

BUG FIXES:

* auth/aws: fix config/rotate-root to store new key [[GH-12715](https://github.com/hashicorp/vault/pull/12715)]
* core/identity: Cleanup alias in the in-memory entity after an alias deletion by ID [[GH-12834](https://github.com/hashicorp/vault/pull/12834)]
* core/identity: Disallow entity alias creation/update if a conflicting alias exists for the target entity and mount combination [[GH-12747](https://github.com/hashicorp/vault/pull/12747)]
* core: Fix a deadlock on HA leadership transfer [[GH-12691](https://github.com/hashicorp/vault/pull/12691)]
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* kmip (enterprise): Fix handling of custom attributes when servicing GetAttributes requests
* kmip (enterprise): Fix handling of invalid role parameters within various vault api calls
* kmip (enterprise): Forward KMIP register operations to the active node
* secrets/keymgmt (enterprise): Fix support for Azure Managed HSM Key Vault instances. [[GH-12957](https://github.com/hashicorp/vault/pull/12957)]
* storage/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* database/postgres: Update postgres library (github.com/lib/pq) to properly remove terminated TLS connections from the connection pool. [[GH-12413](https://github.com/hashicorp/vault/pull/12413)]
* transform (enterprise): Fix an error where the decode response of an expired token is an empty result rather than an error.

## 1.7.5

### 29 September 2021

SECURITY:

* core/identity: A Vault user with write permission to an entity alias ID sharing a mount accessor with another user may acquire this other user’s policies by merging their identities. This vulnerability, CVE-2021-41802, was fixed in Vault and Vault Enterprise 1.7.5 and 1.8.4.

IMPROVEMENTS:

* secrets/pki: Allow signing of self-issued certs with a different signature algorithm. [[GH-12514](https://github.com/hashicorp/vault/pull/12514)]

BUG FIXES:

* agent: Avoid possible `unexpected fault address` panic when using persistent cache. [[GH-12534](https://github.com/hashicorp/vault/pull/12534)]
* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* identity: Fail alias rename if the resulting (name,accessor) exists already [[GH-12473](https://github.com/hashicorp/vault/pull/12473)]
* raft (enterprise): Fix panic when updating auto-snapshot config
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12598](https://github.com/hashicorp/vault/pull/12598)]
* storage/raft: Detect incomplete raft snapshots in api.RaftSnapshot(), and thereby in `vault operator raft snapshot save`. [[GH-12388](https://github.com/hashicorp/vault/pull/12388)]
* ui: Fixed api explorer routing bug [[GH-12354](https://github.com/hashicorp/vault/pull/12354)]

## 1.7.4

### 26 August 2021

SECURITY:

* _UI Secret Caching_: The Vault UI erroneously cached and exposed user-viewed secrets between authenticated sessions in a single shared browser, if the browser window / tab was not refreshed or closed between logout and a subsequent login. This vulnerability, CVE-2021-38554, was fixed in Vault 1.8.0 and will be addressed in pending 1.7.4 / 1.6.6 releases.

CHANGES:

* Alpine: Docker images for Vault 1.6.6+, 1.7.4+, and 1.8.2+ are built with Alpine 3.14, due to CVE-2021-36159
* go: Update go version to 1.15.15 [[GH-12411](https://github.com/hashicorp/vault/pull/12411)]

IMPROVEMENTS:

* ui: Updated node to v14, latest stable build [[GH-12049](https://github.com/hashicorp/vault/pull/12049)]

BUG FIXES:

* replication (enterprise): Fix a panic that could occur when checking the last wal and the log shipper buffer is empty.
* cli: vault debug now puts newlines after every captured log line. [[GH-12175](https://github.com/hashicorp/vault/pull/12175)]
* database/couchbase: change default template to truncate username at 128 characters [[GH-12299](https://github.com/hashicorp/vault/pull/12299)]
* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* secrets/database/cassandra: Fixed issue where the PEM parsing logic of `pem_bundle` and `pem_json` didn't work for CA-only configurations [[GH-11861](https://github.com/hashicorp/vault/pull/11861)]
* secrets/database: Fixed an issue that prevented external database plugin processes from restarting after a shutdown. [[GH-12087](https://github.com/hashicorp/vault/pull/12087)]
* ui: Automatically refresh the page when user logs out [[GH-12035](https://github.com/hashicorp/vault/pull/12035)]
* ui: Fix database role CG access [[GH-12111](https://github.com/hashicorp/vault/pull/12111)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: fix control group access for database credential [[GH-12024](https://github.com/hashicorp/vault/pull/12024)]
* ui: fix oidc login with Safari [[GH-11884](https://github.com/hashicorp/vault/pull/11884)]

## 1.7.3

### June 16th, 2021

CHANGES:

* go: Update go version to 1.15.13 [[GH-11857](https://github.com/hashicorp/vault/pull/11857)]

IMPROVEMENTS:

* db/cassandra: Added tls_server_name to specify server name for TLS validation [[GH-11820](https://github.com/hashicorp/vault/pull/11820)]
* ui: Add specific error message if unseal fails due to license [[GH-11705](https://github.com/hashicorp/vault/pull/11705)]

BUG FIXES:

* auth/jwt: Updates the [hashicorp/cap](https://github.com/hashicorp/cap) library to `v0.1.0` to
bring in a verification key caching fix. [[GH-11784](https://github.com/hashicorp/vault/pull/11784)]
* core (enterprise): serialize access to HSM entropy generation to avoid errors in concurrent key generation.
* secret: fix the bug where transit encrypt batch doesn't work with key_version [[GH-11628](https://github.com/hashicorp/vault/pull/11628)]
* secrets/ad: Forward all creds requests to active node [[GH-76](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/76)] [[GH-11836](https://github.com/hashicorp/vault/pull/11836)]
* tokenutil: Perform the num uses check before token type. [[GH-11647](https://github.com/hashicorp/vault/pull/11647)]

## 1.7.2

### May 20th, 2021

SECURITY:

* Non-Expiring Leases: Vault and Vault Enterprise renewed nearly-expiring token
leases and dynamic secret leases with a zero-second TTL, causing them to be
treated as non-expiring, and never revoked. This issue affects Vault and Vault
Enterprise versions 0.10.0 through 1.7.1, and is fixed in 1.5.9, 1.6.5, and
1.7.2 (CVE-2021-32923).

CHANGES:

* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* auth/gcp: Update to v0.9.1 to use IAM Service Account Credentials API for
signing JWTs [[GH-11494](https://github.com/hashicorp/vault/pull/11494)]

IMPROVEMENTS:

* api, agent: LifetimeWatcher now does more retries when renewal failures occur.  This also impacts Agent auto-auth and leases managed via Agent caching. [[GH-11445](https://github.com/hashicorp/vault/pull/11445)]
* auth/aws: Underlying error included in validation failure message. [[GH-11638](https://github.com/hashicorp/vault/pull/11638)]
* http: Add optional HTTP response headers for hostname and raft node ID [[GH-11289](https://github.com/hashicorp/vault/pull/11289)]
* secrets/aws: add ability to provide a role session name when generating STS credentials [[GH-11345](https://github.com/hashicorp/vault/pull/11345)]
* secrets/database/mongodb: Add ability to customize `SocketTimeout`, `ConnectTimeout`, and `ServerSelectionTimeout` [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database/mongodb: Increased throughput by allowing for multiple request threads to simultaneously update users in MongoDB [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]

BUG FIXES:

* agent/cert: Fix issue where the API client on agent was not honoring certificate
information from the auto-auth config map on renewals or retries. [[GH-11576](https://github.com/hashicorp/vault/pull/11576)]
* agent: Fixed agent templating to use configured tls servername values [[GH-11288](https://github.com/hashicorp/vault/pull/11288)]
* core (enterprise): Fix plugins mounted in namespaces being unable to use password policies [[GH-11596](https://github.com/hashicorp/vault/pull/11596)]
* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]
* identity: Use correct mount accessor when refreshing external group memberships. [[GH-11506](https://github.com/hashicorp/vault/pull/11506)]
* replication: Fix panic trying to update walState during identity group invalidation. [[GH-1865](https://github.com/hashicorp/vault/pull/1865)]
* secrets/database: Fix marshalling to allow providing numeric arguments to external database plugins. [[GH-11451](https://github.com/hashicorp/vault/pull/11451)]
* secrets/database: Fixed minor race condition when rotate-root is called [[GH-11600](https://github.com/hashicorp/vault/pull/11600)]
* secrets/database: Fixes issue for V4 database interface where `SetCredentials` wasn't falling back to using `RotateRootCredentials` if `SetCredentials` is `Unimplemented` [[GH-11585](https://github.com/hashicorp/vault/pull/11585)]
* secrets/keymgmt (enterprise): Fixes audit logging for the read key response.
* storage/raft: Support cluster address change for nodes in a cluster managed by autopilot [[GH-11247](https://github.com/hashicorp/vault/pull/11247)]
* ui: Fix entity group membership and metadata not showing [[GH-11641](https://github.com/hashicorp/vault/pull/11641)]
* ui: Fix text link URL on database roles list [[GH-11597](https://github.com/hashicorp/vault/pull/11597)]

## 1.7.1

### 21 April 2021

SECURITY:

* The PKI Secrets Engine tidy functionality may cause Vault to exclude revoked-but-unexpired certificates from the
  Vault CRL. This vulnerability affects Vault and Vault Enterprise 1.5.1 and newer and was fixed in versions
  1.5.8, 1.6.4, and 1.7.1. (CVE-2021-27668)
* The Cassandra Database and Storage backends were not correctly verifying TLS certificates. This issue affects all
  versions of Vault and Vault Enterprise and was fixed in versions 1.6.4, and 1.7.1. (CVE-2021-27400)

CHANGES:

* go: Update to Go 1.15.11 [[GH-11395](https://github.com/hashicorp/vault/pull/11395)]

IMPROVEMENTS:

* auth/jwt: Adds ability to directly provide service account JSON in G Suite provider config. [[GH-11388](https://github.com/hashicorp/vault/pull/11388)]
* core: Add tls_max_version listener config option. [[GH-11226](https://github.com/hashicorp/vault/pull/11226)]
* core: Add metrics for standby node forwarding. [[GH-11366](https://github.com/hashicorp/vault/pull/11366)]
* core: allow arbitrary length stack traces upon receiving SIGUSR2 (was 32MB) [[GH-11364](https://github.com/hashicorp/vault/pull/11364)]
* storage/raft: Support autopilot for HA only raft storage. [[GH-11260](https://github.com/hashicorp/vault/pull/11260)]

BUG FIXES:

* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* core: Fix goroutine leak when updating rate limit quota [[GH-11371](https://github.com/hashicorp/vault/pull/11371)]
* core: Fix storage entry leak when revoking leases created with non-orphan batch tokens. [[GH-11377](https://github.com/hashicorp/vault/pull/11377)]
* core: requests forwarded by standby weren't always timed out. [[GH-11322](https://github.com/hashicorp/vault/pull/11322)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* replication: Fix: mounts created within a namespace that was part of an Allow
  filtering rule would not appear on performance secondary if created after rule
  was defined.
* replication: Perf standby nodes on newly enabled DR secondary sometimes couldn't connect to active node with TLS errors. [[GH-1823](https://github.com/hashicorp/vault/pull/1823)]
* secrets/database/cassandra: Fixed issue where hostnames were not being validated when using TLS [[GH-11365](https://github.com/hashicorp/vault/pull/11365)]
* secrets/database/cassandra: Updated default statement for password rotation to allow for special characters. This applies to root and static credentials. [[GH-11262](https://github.com/hashicorp/vault/pull/11262)]
* storage/dynamodb: Handle throttled batch write requests by retrying, without which writes could be lost. [[GH-10181](https://github.com/hashicorp/vault/pull/10181)]
* storage/raft: leader_tls_servername wasn't used unless leader_ca_cert_file and/or mTLS were configured. [[GH-11252](https://github.com/hashicorp/vault/pull/11252)]
* storage/raft: using raft for ha_storage with a different storage backend was broken in 1.7.0, now fixed. [[GH-11340](https://github.com/hashicorp/vault/pull/11340)]
* ui: Add root rotation statements support to appropriate database secret engine plugins [[GH-11404](https://github.com/hashicorp/vault/pull/11404)]
* ui: Fix bug where the UI does not recognize version 2 KV until refresh, and fix [object Object] error message [[GH-11258](https://github.com/hashicorp/vault/pull/11258)]
* ui: Fix OIDC bug seen when running on HCP [[GH-11283](https://github.com/hashicorp/vault/pull/11283)]
* ui: Fix namespace-bug on login [[GH-11182](https://github.com/hashicorp/vault/pull/11182)]
* ui: Fix status menu no showing on login [[GH-11213](https://github.com/hashicorp/vault/pull/11213)]
* ui: fix issue where select-one option was not showing in secrets database role creation [[GH-11294](https://github.com/hashicorp/vault/pull/11294)]

## 1.7.0

### 24 March 2021

CHANGES:

* agent: Failed auto-auth attempts are now throttled by an exponential backoff instead of the
~2 second retry delay. The maximum backoff may be configured with the new `max_backoff` parameter,
which defaults to 5 minutes. [[GH-10964](https://github.com/hashicorp/vault/pull/10964)]
* aws/auth: AWS Auth concepts and endpoints that use the "whitelist" and "blacklist" terms
have been updated to more inclusive language (e.g. `/auth/aws/identity-whitelist` has been
updated to`/auth/aws/identity-accesslist`). The old and new endpoints are aliases,
sharing the same underlying data. The legacy endpoint names are considered **deprecated**
and will be removed in a future release (not before Vault 1.9). The complete list of
endpoint changes is available in the [AWS Auth API docs](/api-docs/auth/aws#deprecations-effective-in-vault-1-7).
* go: Update Go version to 1.15.10 [[GH-11114](https://github.com/hashicorp/vault/pull/11114)] [[GH-11173](https://github.com/hashicorp/vault/pull/11173)]

FEATURES:

* **Aerospike Storage Backend**: Add support for using Aerospike as a storage backend [[GH-10131](https://github.com/hashicorp/vault/pull/10131)]
* **Autopilot for Integrated Storage**: A set of features has been added to allow for automatic operator-friendly management of Vault servers. This is only applicable when integrated storage is in use.
  * **Dead Server Cleanup**: Dead servers will periodically be cleaned up and removed from the Raft peer set, to prevent them from interfering with the quorum size and leader elections.
  * **Server Health Checking**: An API has been added to track the state of servers, including their health.
  * **New Server Stabilization**: When a new server is added to the cluster, there will be a waiting period where it must be healthy and stable for a certain amount of time before being promoted to a full, voting member.
* **Tokenization Secrets Engine (Enterprise)**: The Tokenization Secrets Engine is now generally available. We have added support for MySQL, key rotation, and snapshot/restore.
* replication (enterprise): The log shipper is now memory as well as length bound, and length and size can be separately configured.
* agent: Support for persisting the agent cache to disk [[GH-10938](https://github.com/hashicorp/vault/pull/10938)]
* auth/jwt: Adds `max_age` role parameter and `auth_time` claim validation. [[GH-10919](https://github.com/hashicorp/vault/pull/10919)]
* core (enterprise): X-Vault-Index and related headers can be used by clients to manage eventual consistency.
* kmip (enterprise): Use entropy augmentation to generate kmip certificates
* sdk: Private key generation in the certutil package now allows custom io.Readers to be used. [[GH-10653](https://github.com/hashicorp/vault/pull/10653)]
* secrets/aws: add IAM tagging support for iam_user roles [[GH-10953](https://github.com/hashicorp/vault/pull/10953)]
* secrets/database/cassandra: Add ability to customize dynamic usernames [[GH-10906](https://github.com/hashicorp/vault/pull/10906)]
* secrets/database/couchbase: Add ability to customize dynamic usernames [[GH-10995](https://github.com/hashicorp/vault/pull/10995)]
* secrets/database/mongodb: Add ability to customize dynamic usernames [[GH-10858](https://github.com/hashicorp/vault/pull/10858)]
* secrets/database/mssql: Add ability to customize dynamic usernames [[GH-10767](https://github.com/hashicorp/vault/pull/10767)]
* secrets/database/mysql: Add ability to customize dynamic usernames [[GH-10834](https://github.com/hashicorp/vault/pull/10834)]
* secrets/database/postgresql: Add ability to customize dynamic usernames [[GH-10766](https://github.com/hashicorp/vault/pull/10766)]
* secrets/db/snowflake: Added support for Snowflake to the Database Secret Engine [[GH-10603](https://github.com/hashicorp/vault/pull/10603)]
* secrets/keymgmt (enterprise): Adds beta support for distributing and managing keys in AWS KMS.
* secrets/keymgmt (enterprise): Adds general availability for distributing and managing keys in Azure Key Vault.
* secrets/openldap: Added dynamic roles to OpenLDAP similar to the combined database engine [[GH-10996](https://github.com/hashicorp/vault/pull/10996)]
* secrets/terraform: New secret engine for managing Terraform Cloud API tokens [[GH-10931](https://github.com/hashicorp/vault/pull/10931)]
* ui: Adds check for feature flag on application, and updates namespace toolbar on login if present [[GH-10588](https://github.com/hashicorp/vault/pull/10588)]
* ui: Adds the wizard to the Database Secret Engine [[GH-10982](https://github.com/hashicorp/vault/pull/10982)]
* ui: Database secrets engine, supporting MongoDB only [[GH-10655](https://github.com/hashicorp/vault/pull/10655)]

IMPROVEMENTS:

* agent: Add a `vault.retry` stanza that allows specifying number of retries on failure; this applies both to templating and proxied requests. [[GH-11113](https://github.com/hashicorp/vault/pull/11113)]
* agent: Agent can now run as a Windows service. [[GH-10231](https://github.com/hashicorp/vault/pull/10231)]
* agent: Better concurrent request handling on identical requests proxied through Agent. [[GH-10705](https://github.com/hashicorp/vault/pull/10705)]
* agent: Route templating server through cache when persistent cache is enabled. [[GH-10927](https://github.com/hashicorp/vault/pull/10927)]
* agent: change auto-auth to preload an existing token on start [[GH-10850](https://github.com/hashicorp/vault/pull/10850)]
* auth/approle: Secrets ID generation endpoint now returns `secret_id_ttl` as part of its response. [[GH-10826](https://github.com/hashicorp/vault/pull/10826)]
* auth/ldap: Improve consistency in error messages [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]
* auth/okta: Adds support for Okta Verify TOTP MFA. [[GH-10942](https://github.com/hashicorp/vault/pull/10942)]
* changelog: Add dependencies listed in dependencies/2-25-21 [[GH-11015](https://github.com/hashicorp/vault/pull/11015)]
* command/debug: Now collects logs (at level `trace`) as a periodic output. [[GH-10609](https://github.com/hashicorp/vault/pull/10609)]
* core (enterprise): "vault status" command works when a namespace is set. [[GH-10725](https://github.com/hashicorp/vault/pull/10725)]
* core (enterprise): Update Trial Enterprise license from 30 minutes to 6 hours
* core/metrics: Added "vault operator usage" command. [[GH-10365](https://github.com/hashicorp/vault/pull/10365)]
* core/metrics: New telemetry metrics reporting lease expirations by time interval and namespace [[GH-10375](https://github.com/hashicorp/vault/pull/10375)]
* core: Added active since timestamp to the status output of active nodes. [[GH-10489](https://github.com/hashicorp/vault/pull/10489)]
* core: Check audit device with a test message before adding it. [[GH-10520](https://github.com/hashicorp/vault/pull/10520)]
* core: Track barrier encryption count and automatically rotate after a large number of operations or on a schedule [[GH-10774](https://github.com/hashicorp/vault/pull/10774)]
* core: add metrics for active entity count [[GH-10514](https://github.com/hashicorp/vault/pull/10514)]
* core: add partial month client count api [[GH-11022](https://github.com/hashicorp/vault/pull/11022)]
* core: dev mode listener allows unauthenticated sys/metrics requests [[GH-10992](https://github.com/hashicorp/vault/pull/10992)]
* core: reduce memory used by leases [[GH-10726](https://github.com/hashicorp/vault/pull/10726)]
* secrets/gcp: Truncate ServiceAccount display names longer than 100 characters. [[GH-10558](https://github.com/hashicorp/vault/pull/10558)]
* storage/raft (enterprise): Listing of peers is now allowed on DR secondary
cluster nodes, as an update operation that takes in DR operation token for
authenticating the request.
* transform (enterprise): Improve FPE transformation performance
* transform (enterprise): Use transactions with batch tokenization operations for improved performance
* ui: Clarify language on usage metrics page empty state [[GH-10951](https://github.com/hashicorp/vault/pull/10951)]
* ui: Customize MongoDB input fields on Database Secrets Engine [[GH-10949](https://github.com/hashicorp/vault/pull/10949)]
* ui: Upgrade Ember-cli from 3.8 to 3.22. [[GH-9972](https://github.com/hashicorp/vault/pull/9972)]
* ui: Upgrade Storybook from 5.3.19 to 6.1.17. [[GH-10904](https://github.com/hashicorp/vault/pull/10904)]
* ui: Upgrade date-fns from 1.3.0 to 2.16.1. [[GH-10848](https://github.com/hashicorp/vault/pull/10848)]
* ui: Upgrade dependencies to resolve potential JS vulnerabilities [[GH-10677](https://github.com/hashicorp/vault/pull/10677)]
* ui: better errors on Database secrets engine role create [[GH-10980](https://github.com/hashicorp/vault/pull/10980)]

BUG FIXES:

* agent: Only set the namespace if the VAULT_NAMESPACE env var isn't present [[GH-10556](https://github.com/hashicorp/vault/pull/10556)]
* agent: Set TokenParent correctly in the Index to be cached. [[GH-10833](https://github.com/hashicorp/vault/pull/10833)]
* agent: Set namespace for template server in agent. [[GH-10757](https://github.com/hashicorp/vault/pull/10757)]
* api/sys/config/ui: Fixes issue where multiple UI custom header values are ignored and only the first given value is used [[GH-10490](https://github.com/hashicorp/vault/pull/10490)]
* api: Fixes CORS API methods that were outdated and invalid [[GH-10444](https://github.com/hashicorp/vault/pull/10444)]
* auth/jwt: Fixes `bound_claims` validation for provider-specific group and user info fetching. [[GH-10546](https://github.com/hashicorp/vault/pull/10546)]
* auth/jwt: Fixes an issue where JWT verification keys weren't updated after a `jwks_url` change. [[GH-10919](https://github.com/hashicorp/vault/pull/10919)]
* auth/jwt: Fixes an issue where `jwt_supported_algs` were not being validated for JWT auth using
`jwks_url` and `jwt_validation_pubkeys`. [[GH-10919](https://github.com/hashicorp/vault/pull/10919)]
* auth/oci: Fixes alias name to use the role name, and not the literal string `name` [[GH-10](https://github.com/hashicorp/vault-plugin-auth-oci/pull/10)] [[GH-10952](https://github.com/hashicorp/vault/pull/10952)]
* consul-template: Update consul-template vendor version and associated dependencies to master,
pulling in <https://github.com/hashicorp/consul-template/pull/1447> [[GH-10756](https://github.com/hashicorp/vault/pull/10756)]
* core (enterprise): Limit entropy augmentation during token generation to root tokens. [[GH-10487](https://github.com/hashicorp/vault/pull/10487)]
* core (enterprise): Vault EGP policies attached to path * were not correctly scoped to the namespace.
* core/identity: Fix deadlock in entity merge endpoint. [[GH-10877](https://github.com/hashicorp/vault/pull/10877)]
* core: Avoid disclosing IP addresses in the errors of unauthenticated requests [[GH-10579](https://github.com/hashicorp/vault/pull/10579)]
* core: Fix client.Clone() to include the address [[GH-10077](https://github.com/hashicorp/vault/pull/10077)]
* core: Fix duplicate quotas on performance standby nodes. [[GH-10855](https://github.com/hashicorp/vault/pull/10855)]
* core: Fix rate limit resource quota migration from 1.5.x to 1.6.x by ensuring `purgeInterval` and
`staleAge` are set appropriately. [[GH-10536](https://github.com/hashicorp/vault/pull/10536)]
* core: Make all APIs that report init status consistent, and make them report
initialized=true when a Raft join is in progress. [[GH-10498](https://github.com/hashicorp/vault/pull/10498)]
* core: Make the response to an unauthenticated request to sys/internal endpoints consistent regardless of mount existence. [[GH-10650](https://github.com/hashicorp/vault/pull/10650)]
* core: Turn off case sensitivity for allowed entity alias check during token create operation. [[GH-10743](https://github.com/hashicorp/vault/pull/10743)]
* http: change max_request_size to be unlimited when the config value is less than 0 [[GH-10072](https://github.com/hashicorp/vault/pull/10072)]
* license: Fix license caching issue that prevents new licenses to get picked up by the license manager [[GH-10424](https://github.com/hashicorp/vault/pull/10424)]
* metrics: Protect emitMetrics from panicking during post-seal [[GH-10708](https://github.com/hashicorp/vault/pull/10708)]
* quotas/rate-limit: Fix quotas enforcing old rate limit quota paths [[GH-10689](https://github.com/hashicorp/vault/pull/10689)]
* replication (enterprise): Fix bug with not starting merkle sync while requests are in progress
* secrets/database/influxdb: Fix issue where not all errors from InfluxDB were being handled [[GH-10384](https://github.com/hashicorp/vault/pull/10384)]
* secrets/database/mysql: Fixes issue where the DisplayName within generated usernames was the incorrect length [[GH-10433](https://github.com/hashicorp/vault/pull/10433)]
* secrets/database: Sanitize `private_key` field when reading database plugin config [[GH-10416](https://github.com/hashicorp/vault/pull/10416)]
* secrets/gcp: Fix issue with account and iam_policy roleset WALs not being removed after attempts when GCP project no longer exists [[GH-10759](https://github.com/hashicorp/vault/pull/10759)]
* secrets/transit: allow for null string to be used for optional parameters in encrypt and decrypt [[GH-10386](https://github.com/hashicorp/vault/pull/10386)]
* serviceregistration: Fix race during shutdown of Consul service registration. [[GH-10901](https://github.com/hashicorp/vault/pull/10901)]
* storage/raft (enterprise): Automated snapshots with Azure required specifying
`azure_blob_environment`, which should have had as a default `AZUREPUBLICCLOUD`.
* storage/raft (enterprise): Reading a non-existent auto snapshot config now returns 404.
* storage/raft (enterprise): The parameter aws_s3_server_kms_key was misnamed and
didn't work.  Renamed to aws_s3_kms_key, and make it work so that when provided
the given key will be used to encrypt the snapshot using AWS KMS.
* transform (enterprise): Fix bug tokenization handling metadata on exportable stores
* transform (enterprise): Fix bug where tokenization store changes are persisted but don't take effect
* transform (enterprise): Fix transform configuration not handling `stores` parameter on the legacy path
* transform (enterprise): Make expiration timestamps human readable
* transform (enterprise): Return false for invalid tokens on the validate endpoint rather than returning an HTTP error
* ui: Add role from database connection automatically populates the database for new role [[GH-11119](https://github.com/hashicorp/vault/pull/11119)]
* ui: Fix bug in Transform secret engine when a new role is added and then removed from a transformation [[GH-10417](https://github.com/hashicorp/vault/pull/10417)]
* ui: Fix bug that double encodes secret route when there are spaces in the path and makes you unable to view the version history. [[GH-10596](https://github.com/hashicorp/vault/pull/10596)]
* ui: Fix expected response from feature-flags endpoint [[GH-10684](https://github.com/hashicorp/vault/pull/10684)]
* ui: Fix footer URL linking to the correct version changelog. [[GH-10491](https://github.com/hashicorp/vault/pull/10491)]

DEPRECATIONS:

* aws/auth: AWS Auth endpoints that use the "whitelist" and "blacklist" terms have been deprecated.
Refer to the CHANGES section for additional details.

## 1.6.7

### 29 September 2021

BUG FIXES:

* core (enterprise): Fix bug where password generation through password policies do not work on namespaces if performed outside a request callback or from an external plugin. [[GH-12635](https://github.com/hashicorp/vault/pull/12635)]
* core (enterprise): Only delete quotas on primary cluster. [[GH-12339](https://github.com/hashicorp/vault/pull/12339)]
* secrets/db: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [[GH-12563](https://github.com/hashicorp/vault/pull/12563)]
* secrets/openldap: Fix bug where Vault can rotate static role passwords early during start up under certain conditions. [#28](https://github.com/hashicorp/vault-plugin-secrets-openldap/pull/28) [[GH-12597](https://github.com/hashicorp/vault/pull/12597)]

## 1.6.6

### 26 August 2021

SECURITY:

* _UI Secret Caching_: The Vault UI erroneously cached and exposed user-viewed secrets between authenticated sessions in a single shared browser, if the browser window / tab was not refreshed or closed between logout and a subsequent login. This vulnerability, CVE-2021-38554, was fixed in Vault 1.8.0 and will be addressed in pending 1.7.4 / 1.6.6 releases.

CHANGES:

* Alpine: Docker images for Vault 1.6.6+, 1.7.4+, and 1.8.2+ are built with Alpine 3.14, due to CVE-2021-36159
* go: Update go version to 1.15.15 [[GH-12423](https://github.com/hashicorp/vault/pull/12423)]

IMPROVEMENTS:

* db/cassandra: Added tls_server_name to specify server name for TLS validation [[GH-11820](https://github.com/hashicorp/vault/pull/11820)]

BUG FIXES:

* physical/raft: Fix safeio.Rename error when restoring snapshots on windows [[GH-12377](https://github.com/hashicorp/vault/pull/12377)]
* secret: fix the bug where transit encrypt batch doesn't work with key_version [[GH-11628](https://github.com/hashicorp/vault/pull/11628)]
* secrets/database: Fixed an issue that prevented external database plugin processes from restarting after a shutdown. [[GH-12087](https://github.com/hashicorp/vault/pull/12087)]
* ui: Automatically refresh the page when user logs out [[GH-12035](https://github.com/hashicorp/vault/pull/12035)]
* ui: Fixes metrics page when read on counter config not allowed [[GH-12348](https://github.com/hashicorp/vault/pull/12348)]
* ui: fix oidc login with Safari [[GH-11884](https://github.com/hashicorp/vault/pull/11884)]

## 1.6.5

### May 20th, 2021

SECURITY:

* Non-Expiring Leases: Vault and Vault Enterprise renewed nearly-expiring token
leases and dynamic secret leases with a zero-second TTL, causing them to be
treated as non-expiring, and never revoked. This issue affects Vault and Vault
Enterprise versions 0.10.0 through 1.7.1, and is fixed in 1.5.9, 1.6.5, and
1.7.2 (CVE-2021-32923).

CHANGES:

* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* auth/gcp: Update to v0.8.1 to use IAM Service Account Credentials API for
signing JWTs [[GH-11498](https://github.com/hashicorp/vault/pull/11498)]

BUG FIXES:

* core (enterprise): Fix plugins mounted in namespaces being unable to use password policies [[GH-11596](https://github.com/hashicorp/vault/pull/11596)]
* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]
* secrets/database: Fix marshalling to allow providing numeric arguments to external database plugins. [[GH-11451](https://github.com/hashicorp/vault/pull/11451)]
* secrets/database: Fixes issue for V4 database interface where `SetCredentials` wasn't falling back to using `RotateRootCredentials` if `SetCredentials` is `Unimplemented` [[GH-11585](https://github.com/hashicorp/vault/pull/11585)]
* ui: Fix namespace-bug on login [[GH-11182](https://github.com/hashicorp/vault/pull/11182)]

## 1.6.4

### 21 April 2021

SECURITY:

* The PKI Secrets Engine tidy functionality may cause Vault to exclude revoked-but-unexpired certificates from the
  Vault CRL. This vulnerability affects Vault and Vault Enterprise 1.5.1 and newer and was fixed in versions
  1.5.8, 1.6.4, and 1.7.1. (CVE-2021-27668)
* The Cassandra Database and Storage backends were not correctly verifying TLS certificates. This issue affects all
  versions of Vault and Vault Enterprise and was fixed in versions 1.6.4, and 1.7.1. (CVE-2021-27400)

CHANGES:

* go: Update to Go 1.15.11 [[GH-11396](https://github.com/hashicorp/vault/pull/11396)]

IMPROVEMENTS:

* command/debug: Now collects logs (at level `trace`) as a periodic output. [[GH-10609](https://github.com/hashicorp/vault/pull/10609)]
* core: Add tls_max_version listener config option. [[GH-11226](https://github.com/hashicorp/vault/pull/11226)]
* core: allow arbitrary length stack traces upon receiving SIGUSR2 (was 32MB) [[GH-11364](https://github.com/hashicorp/vault/pull/11364)]

BUG FIXES:

* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* core: Fix goroutine leak when updating rate limit quota [[GH-11371](https://github.com/hashicorp/vault/pull/11371)]
* core: Fix storage entry leak when revoking leases created with non-orphan batch tokens. [[GH-11377](https://github.com/hashicorp/vault/pull/11377)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* pki: Preserve ordering of all DN attribute values when issuing certificates [[GH-11259](https://github.com/hashicorp/vault/pull/11259)]
* replication: Fix: mounts created within a namespace that was part of an Allow
  filtering rule would not appear on performance secondary if created after rule
  was defined.
* secrets/database/cassandra: Fixed issue where hostnames were not being validated when using TLS [[GH-11365](https://github.com/hashicorp/vault/pull/11365)]
* storage/raft: leader_tls_servername wasn't used unless leader_ca_cert_file and/or mTLS were configured. [[GH-11252](https://github.com/hashicorp/vault/pull/11252)]

## 1.6.3

### February 25, 2021

SECURITY:

* Limited Unauthenticated License Metadata Read: We addressed a security vulnerability that allowed for the unauthenticated
reading of Vault license metadata from DR Secondaries. This vulnerability affects Vault Enterprise and is
fixed in 1.6.3 (CVE-2021-27668).

CHANGES:

* secrets/mongodbatlas: Move from whitelist to access list API [[GH-10966](https://github.com/hashicorp/vault/pull/10966)]

IMPROVEMENTS:

* ui: Clarify language on usage metrics page empty state [[GH-10951](https://github.com/hashicorp/vault/pull/10951)]

BUG FIXES:

* auth/kubernetes: Cancel API calls to TokenReview endpoint when request context
is closed [[GH-10930](https://github.com/hashicorp/vault/pull/10930)]
* core/identity: Fix deadlock in entity merge endpoint. [[GH-10877](https://github.com/hashicorp/vault/pull/10877)]
* quotas: Fix duplicate quotas on performance standby nodes. [[GH-10855](https://github.com/hashicorp/vault/pull/10855)]
* quotas/rate-limit: Fix quotas enforcing old rate limit quota paths [[GH-10689](https://github.com/hashicorp/vault/pull/10689)]
* replication (enterprise): Don't write request count data on DR Secondaries.
Fixes DR Secondaries becoming out of sync approximately every 30s. [[GH-10970](https://github.com/hashicorp/vault/pull/10970)]
* secrets/azure (enterprise): Forward service principal credential creation to the
primary cluster if called on a performance standby or performance secondary. [[GH-10902](https://github.com/hashicorp/vault/pull/10902)]

## 1.6.2

### January 29, 2021

SECURITY:

* IP Address Disclosure: We fixed a vulnerability where, under some error
conditions, Vault would return an error message disclosing internal IP
addresses. This vulnerability affects Vault and Vault Enterprise and is fixed in
1.6.2 (CVE-2021-3024).
* Limited Unauthenticated Remove Peer: As of Vault 1.6, the remove-peer command
on DR secondaries did not require authentication. This issue impacts the
stability of HA architecture, as a bad actor could remove all standby
nodes from a DR
secondary. This issue affects Vault Enterprise 1.6.0 and 1.6.1, and is fixed in
1.6.2 (CVE-2021-3282).
* Mount Path Disclosure: Vault previously returned different HTTP status codes for
existent and non-existent mount paths. This behavior would allow unauthenticated
brute force attacks to reveal which paths had valid mounts. This issue affects
Vault and Vault Enterprise and is fixed in 1.6.2 (CVE-2020-25594).

CHANGES:

* go: Update go version to 1.15.7 [[GH-10730](https://github.com/hashicorp/vault/pull/10730)]

FEATURES:

* ui: Adds check for feature flag on application, and updates namespace toolbar on login if present [[GH-10588](https://github.com/hashicorp/vault/pull/10588)]

IMPROVEMENTS:

* core (enterprise): "vault status" command works when a namespace is set. [[GH-10725](https://github.com/hashicorp/vault/pull/10725)]
* core: reduce memory used by leases [[GH-10726](https://github.com/hashicorp/vault/pull/10726)]
* storage/raft (enterprise): Listing of peers is now allowed on DR secondary
cluster nodes, as an update operation that takes in DR operation token for
authenticating the request.
* core: allow setting tls_servername for raft retry/auto-join [[GH-10698](https://github.com/hashicorp/vault/pull/10698)]

BUG FIXES:

* agent: Set namespace for template server in agent. [[GH-10757](https://github.com/hashicorp/vault/pull/10757)]
* core: Make the response to an unauthenticated request to sys/internal endpoints consistent regardless of mount existence. [[GH-10650](https://github.com/hashicorp/vault/pull/10650)]
* metrics: Protect emitMetrics from panicking during post-seal [[GH-10708](https://github.com/hashicorp/vault/pull/10708)]
* secrets/gcp: Fix issue with account and iam_policy roleset WALs not being removed after attempts when GCP project no longer exists [[GH-10759](https://github.com/hashicorp/vault/pull/10759)]
* storage/raft (enterprise): Automated snapshots with Azure required specifying
`azure_blob_environment`, which should have had as a default `AZUREPUBLICCLOUD`.
* storage/raft (enterprise): Autosnapshots config and storage weren't excluded from
performance replication, causing conflicts and errors.
* ui: Fix bug that double encodes secret route when there are spaces in the path and makes you unable to view the version history. [[GH-10596](https://github.com/hashicorp/vault/pull/10596)]
* ui: Fix expected response from feature-flags endpoint [[GH-10684](https://github.com/hashicorp/vault/pull/10684)]

## 1.6.1

### December 16, 2020

SECURITY:

* LDAP Auth Method: We addressed an issue where error messages returned by the
  LDAP auth method allowed user enumeration [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]. This vulnerability affects Vault OSS and Vault
  Enterprise and is fixed in 1.5.6 and 1.6.1 (CVE-2020-35177).
* Sentinel EGP: We've fixed incorrect handling of namespace paths to prevent
  users within namespaces from applying Sentinel EGP policies to paths above
  their namespace. This vulnerability affects Vault Enterprise and is fixed in
  1.5.6 and 1.6.1 (CVE-2020-35453).

IMPROVEMENTS:

* auth/ldap: Improve consistency in error messages [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]
* core/metrics: Added "vault operator usage" command. [[GH-10365](https://github.com/hashicorp/vault/pull/10365)]
* secrets/gcp: Truncate ServiceAccount display names longer than 100 characters. [[GH-10558](https://github.com/hashicorp/vault/pull/10558)]

BUG FIXES:

* agent: Only set the namespace if the VAULT_NAMESPACE env var isn't present [[GH-10556](https://github.com/hashicorp/vault/pull/10556)]
* auth/jwt: Fixes `bound_claims` validation for provider-specific group and user info fetching. [[GH-10546](https://github.com/hashicorp/vault/pull/10546)]
* core (enterprise): Vault EGP policies attached to path * were not correctly scoped to the namespace.
* core: Avoid deadlocks by ensuring that if grabLockOrStop returns stopped=true, the lock will not be held. [[GH-10456](https://github.com/hashicorp/vault/pull/10456)]
* core: Fix client.Clone() to include the address [[GH-10077](https://github.com/hashicorp/vault/pull/10077)]
* core: Fix rate limit resource quota migration from 1.5.x to 1.6.x by ensuring `purgeInterval` and
`staleAge` are set appropriately. [[GH-10536](https://github.com/hashicorp/vault/pull/10536)]
* core: Make all APIs that report init status consistent, and make them report
initialized=true when a Raft join is in progress. [[GH-10498](https://github.com/hashicorp/vault/pull/10498)]
* secrets/database/influxdb: Fix issue where not all errors from InfluxDB were being handled [[GH-10384](https://github.com/hashicorp/vault/pull/10384)]
* secrets/database/mysql: Fixes issue where the DisplayName within generated usernames was the incorrect length [[GH-10433](https://github.com/hashicorp/vault/pull/10433)]
* secrets/database: Sanitize `private_key` field when reading database plugin config [[GH-10416](https://github.com/hashicorp/vault/pull/10416)]
* secrets/transit: allow for null string to be used for optional parameters in encrypt and decrypt [[GH-10386](https://github.com/hashicorp/vault/pull/10386)]
* storage/raft (enterprise): The parameter aws_s3_server_kms_key was misnamed and didn't work.  Renamed to aws_s3_kms_key, and make it work so that when provided the given key will be used to encrypt the snapshot using AWS KMS.
* transform (enterprise): Fix bug tokenization handling metadata on exportable stores
* transform (enterprise): Fix transform configuration not handling `stores` parameter on the legacy path
* transform (enterprise): Make expiration timestamps human readable
* transform (enterprise): Return false for invalid tokens on the validate endpoint rather than returning an HTTP error
* transform (enterprise): Fix bug where tokenization store changes are persisted but don't take effect
* ui: Fix bug in Transform secret engine when a new role is added and then removed from a transformation [[GH-10417](https://github.com/hashicorp/vault/pull/10417)]
* ui: Fix footer URL linking to the correct version changelog. [[GH-10491](https://github.com/hashicorp/vault/pull/10491)]
* ui: Fox radio click on secrets and auth list pages. [[GH-10586](https://github.com/hashicorp/vault/pull/10586)]

## 1.6.0

### November 11th, 2020

NOTE:

Binaries for 32-bit macOS (i.e. the `darwin_386` build) will no longer be published. This target was dropped in the latest version of the Go compiler.

CHANGES:

* agent: Agent now properly returns a non-zero exit code on error, such as one due to template rendering failure. Using `error_on_missing_key` in the template config will cause agent to immediately exit on failure. In order to make agent properly exit due to continuous failure from template rendering errors, the old behavior of indefinitely restarting the template server is now changed to exit once the default retry attempt of 12 times (with exponential backoff) gets exhausted. [[GH-9670](https://github.com/hashicorp/vault/pull/9670)]
* token: Periodic tokens generated by auth methods will have the period value stored in its token entry. [[GH-7885](https://github.com/hashicorp/vault/pull/7885)]
* core: New telemetry metrics reporting mount table size and number of entries [[GH-10201](https://github.com/hashicorp/vault/pull/10201)]
* go: Updated Go version to 1.15.4 [[GH-10366](https://github.com/hashicorp/vault/pull/10366)]

FEATURES:

* **Couchbase Secrets**: Vault can now manage static and dynamic credentials for Couchbase. [[GH-9664](https://github.com/hashicorp/vault/pull/9664)]
* **Expanded Password Policy Support**: Custom password policies are now supported for all database engines.
* **Integrated Storage Auto Snapshots (Enterprise)**: This feature enables an operator to schedule snapshots of the integrated storage backend and ensure those snapshots are persisted elsewhere.
* **Integrated Storage Cloud Auto Join**: This feature for integrated storage enables Vault nodes running in the cloud to automatically discover and join a Vault cluster via operator-supplied metadata.
* **Key Management Secrets Engine (Enterprise; Tech Preview)**: This new secret engine allows securely distributing and managing keys to Azure cloud KMS services.
* **Seal Migration**: With Vault 1.6, we will support migrating from an auto unseal mechanism to a different mechanism of the same type. For example, if you were using an AWS KMS key to automatically unseal, you can now migrate to a different AWS KMS key.
* **Tokenization (Enterprise; Tech Preview)**: Tokenization supports creating irreversible “tokens” from sensitive data. Tokens can be used in less secure environments, protecting the original data.
* **Vault Client Count**: Vault now counts the number of active entities (and non-entity tokens) per month and makes this information available via the "Metrics" section of the UI.

IMPROVEMENTS:

* auth/approle: Role names can now be referenced in templated policies through the `approle.metadata.role_name` property [[GH-9529](https://github.com/hashicorp/vault/pull/9529)]
* auth/aws: Improve logic check on wildcard `BoundIamPrincipalARNs` and include role name on error messages on check failure [[GH-10036](https://github.com/hashicorp/vault/pull/10036)]
* auth/jwt: Add support for fetching groups and user information from G Suite during authentication. [[GH-123](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/123)]
* auth/jwt: Adding EdDSA (ed25519) to supported algorithms [[GH-129](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/129)]
* auth/jwt: Improve cli authorization error [[GH-137](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/137)]
* auth/jwt: Add OIDC namespace_in_state option [[GH-140](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/140)]
* secrets/transit: fix missing plaintext in bulk decrypt response [[GH-9991](https://github.com/hashicorp/vault/pull/9991)]
* command/server: Delay informational messages in -dev mode until logs have settled. [[GH-9702](https://github.com/hashicorp/vault/pull/9702)]
* command/server: Add environment variable support for `disable_mlock`. [[GH-9931](https://github.com/hashicorp/vault/pull/9931)]
* core/metrics: Add metrics for storage cache [[GH_10079](https://github.com/hashicorp/vault/pull/10079)]
* core/metrics: Add metrics for leader status [[GH 10147](https://github.com/hashicorp/vault/pull/10147)]
* physical/azure: Add the ability to use Azure Instance Metadata Service to set the credentials for Azure Blob storage on the backend. [[GH-10189](https://github.com/hashicorp/vault/pull/10189)]
* sdk/framework: Add a time type for API fields. [[GH-9911](https://github.com/hashicorp/vault/pull/9911)]
* secrets/database: Added support for password policies to all databases [[GH-9641](https://github.com/hashicorp/vault/pull/9641),
  [and more](https://github.com/hashicorp/vault/pulls?q=is%3Apr+is%3Amerged+dbpw)]
* secrets/database/cassandra: Added support for static credential rotation [[GH-10051](https://github.com/hashicorp/vault/pull/10051)]
* secrets/database/elasticsearch: Added support for static credential rotation [[GH-19](https://github.com/hashicorp/vault-plugin-database-elasticsearch/pull/19)]
* secrets/database/hanadb: Added support for root credential & static credential rotation [[GH-10142](https://github.com/hashicorp/vault/pull/10142)]
* secrets/database/hanadb: Default password generation now includes dashes. Custom statements may need to be updated
  to include quotes around the password field [[GH-10142](https://github.com/hashicorp/vault/pull/10142)]
* secrets/database/influxdb: Added support for static credential rotation [[GH-10118](https://github.com/hashicorp/vault/pull/10118)]
* secrets/database/mongodbatlas: Added support for root credential rotation [[GH-14](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/14)]
* secrets/database/mongodbatlas: Support scopes field in creations statements for MongoDB Atlas database plugin [[GH-15](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/15)]
* seal/awskms: Add logging during awskms auto-unseal [[GH-9794](https://github.com/hashicorp/vault/pull/9794)]
* storage/azure: Update SDK library to use [azure-storage-blob-go](https://github.com/Azure/azure-storage-blob-go) since previous library has been deprecated. [[GH-9577](https://github.com/hashicorp/vault/pull/9577/)]
* secrets/ad: `rotate-root` now supports POST requests like other secret engines [[GH-70](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/70)]
* ui: Add ui functionality for the Transform Secret Engine [[GH-9665](https://github.com/hashicorp/vault/pull/9665)]
* ui: Pricing metrics dashboard [[GH-10049](https://github.com/hashicorp/vault/pull/10049)]

BUG FIXES:

* auth/jwt: Fix bug preventing config edit UI from rendering [[GH-141](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/141)]
* cli: Don't open or overwrite a raft snapshot file on an unsuccessful `vault operator raft snapshot` [[GH-9894](https://github.com/hashicorp/vault/pull/9894)]
* core: Implement constant time version of shamir GF(2^8) math [[GH-9932](https://github.com/hashicorp/vault/pull/9932)]
* core: Fix resource leak in plugin API (plugin-dependent, not all plugins impacted) [[GH-9557](https://github.com/hashicorp/vault/pull/9557)]
* core: Fix race involved in enabling certain features via a license change
* core: Fix error handling in HCL parsing of objects with invalid syntax [[GH-410](https://github.com/hashicorp/hcl/pull/410)]
* identity: Check for timeouts in entity API [[GH-9925](https://github.com/hashicorp/vault/pull/9925)]
* secrets/database: Fix handling of TLS options in mongodb connection strings [[GH-9519](https://github.com/hashicorp/vault/pull/9519)]
* secrets/gcp: Ensure that the IAM policy version is appropriately set after a roleset's bindings have changed. [[GH-93](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/93)]
* ui: Mask LDAP bindpass while typing [[GH-10087](https://github.com/hashicorp/vault/pull/10087)]
* ui: Update language in promote dr modal flow [[GH-10155](https://github.com/hashicorp/vault/pull/10155)]
* ui: Update language on replication primary dashboard for clarity [[GH-10205](https://github.com/hashicorp/vault/pull/10217)]
* core: Fix bug where updating an existing path quota could introduce a conflict. [[GH-10285](https://github.com/hashicorp/vault/pull/10285)]

## 1.5.9

### May 20th, 2021

SECURITY:

* Non-Expiring Leases: Vault and Vault Enterprise renewed nearly-expiring token
leases and dynamic secret leases with a zero-second TTL, causing them to be
treated as non-expiring, and never revoked. This issue affects Vault and Vault
Enterprise versions 0.10.0 through 1.7.1, and is fixed in 1.5.9, 1.6.5, and
1.7.2 (CVE-2021-32923).

CHANGES:

* agent: Update to use IAM Service Account Credentials endpoint for signing JWTs
when using GCP Auto-Auth method [[GH-11473](https://github.com/hashicorp/vault/pull/11473)]
* auth/gcp: Update to v0.7.2 to use IAM Service Account Credentials API for
signing JWTs [[GH-11499](https://github.com/hashicorp/vault/pull/11499)]

BUG FIXES:

* core: correct logic for renewal of leases nearing their expiration time. [[GH-11650](https://github.com/hashicorp/vault/pull/11650)]

## 1.5.8

### 21 April 2021

SECURITY:

* The PKI Secrets Engine tidy functionality may cause Vault to exclude revoked-but-unexpired certificates from the
  Vault CRL. This vulnerability affects Vault and Vault Enterprise 1.5.1 and newer and was fixed in versions
  1.5.8, 1.6.4, and 1.7.1. (CVE-2021-27668)

CHANGES:

* go: Update to Go 1.14.15 [[GH-11397](https://github.com/hashicorp/vault/pull/11397)]

IMPROVEMENTS:

* core: Add tls_max_version listener config option. [[GH-11226](https://github.com/hashicorp/vault/pull/11226)]

BUG FIXES:

* core/identity: Fix deadlock in entity merge endpoint. [[GH-10877](https://github.com/hashicorp/vault/pull/10877)]
* core: Fix cleanup of storage entries from cubbyholes within namespaces. [[GH-11408](https://github.com/hashicorp/vault/pull/11408)]
* pki: Only remove revoked entry for certificates during tidy if they are past their NotAfter value [[GH-11367](https://github.com/hashicorp/vault/pull/11367)]
* core: Avoid deadlocks by ensuring that if grabLockOrStop returns stopped=true, the lock will not be held. [[GH-10456](https://github.com/hashicorp/vault/pull/10456)]

## 1.5.7

### January 29, 2021

SECURITY:

* IP Address Disclosure: We fixed a vulnerability where, under some error
conditions, Vault would return an error message disclosing internal IP
addresses. This vulnerability affects Vault and Vault Enterprise and is fixed in
1.6.2 and 1.5.7 (CVE-2021-3024).
* Mount Path Disclosure: Vault previously returned different HTTP status codes for
existent and non-existent mount paths. This behavior would allow unauthenticated
brute force attacks to reveal which paths had valid mounts. This issue affects
Vault and Vault Enterprise and is fixed in 1.6.2 and 1.5.7 (CVE-2020-25594).

IMPROVEMENTS:

* storage/raft (enterprise): Listing of peers is now allowed on DR secondary
cluster nodes, as an update operation that takes in DR operation token for
authenticating the request.

BUG FIXES:

* core: Avoid disclosing IP addresses in the errors of unauthenticated requests [[GH-10579](https://github.com/hashicorp/vault/pull/10579)]
* core: Make the response to an unauthenticated request to sys/internal endpoints consistent regardless of mount existence. [[GH-10650](https://github.com/hashicorp/vault/pull/10650)]

## 1.5.6

### December 16, 2020

SECURITY:

* LDAP Auth Method: We addressed an issue where error messages returned by the
  LDAP auth method allowed user enumeration [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]. This vulnerability affects Vault OSS and Vault
  Enterprise and is fixed in 1.5.6 and 1.6.1 (CVE-2020-35177).
* Sentinel EGP: We've fixed incorrect handling of namespace paths to prevent
  users within namespaces from applying Sentinel EGP policies to paths above
  their namespace. This vulnerability affects Vault Enterprise and is fixed in
  1.5.6 and 1.6.1.

IMPROVEMENTS:

* auth/ldap: Improve consistency in error messages [[GH-10537](https://github.com/hashicorp/vault/pull/10537)]

BUG FIXES:

* core (enterprise): Vault EGP policies attached to path * were not correctly scoped to the namespace.
* core: Fix bug where updating an existing path quota could introduce a conflict [[GH-10285](https://github.com/hashicorp/vault/pull/10285)]
* core: Fix client.Clone() to include the address [[GH-10077](https://github.com/hashicorp/vault/pull/10077)]
* quotas (enterprise): Reset cache before loading quotas in the db during startup
* secrets/transit: allow for null string to be used for optional parameters in encrypt and decrypt [[GH-10386](https://github.com/hashicorp/vault/pull/10386)]

## 1.5.5

### October 21, 2020

IMPROVEMENTS:

* auth/aws, core/seal, secret/aws: Set default IMDS timeouts to match AWS SDK [[GH-10133](https://github.com/hashicorp/vault/pull/10133)]

BUG FIXES:

* auth/aws: Restrict region selection when in the aws-us-gov partition to avoid IAM errors [[GH-9947](https://github.com/hashicorp/vault/pull/9947)]
* core (enterprise): Allow operators to add and remove (Raft) peers in a DR secondary cluster using Integrated Storage.
* core (enterprise): Add DR operation token to the remove peer API and CLI command (when DR secondary).
* core (enterprise): Fix deadlock in handling EGP policies
* core (enterprise): Fix extraneous error messages in DR Cluster
* secrets/mysql: Conditionally overwrite TLS parameters for MySQL secrets engine [[GH-9729](https://github.com/hashicorp/vault/pull/9729)]
* secrets/ad: Fix bug where `password_policy` setting was not using correct key when `ad/config` was read [[GH-71](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/71)]
* ui: Fix issue with listing roles and methods on the same auth methods with different names [[GH-10122](https://github.com/hashicorp/vault/pull/10122)]

## 1.5.4

### September 24th, 2020

SECURITY:

* Batch Token Expiry: We addressed an issue where batch token leases could outlive their TTL because we were not scheduling the expiration time correctly. This vulnerability affects Vault OSS and Vault Enterprise 1.0 and newer and is fixed in 1.4.7 and 1.5.4 (CVE-2020-25816).

IMPROVEMENTS:

* secrets/pki: Handle expiration of a cert not in storage as a success [[GH-9880](https://github.com/hashicorp/vault/pull/9880)]
* auth/kubernetes: Add an option to disable defaulting to the local CA cert and service account JWT when running in a Kubernetes pod [[GH-97]](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/97)
* secrets/gcp: Add check for 403 during rollback to prevent repeated deletion calls [[GH-97](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/97)]
* core: Disable usage metrics collection on performance standby nodes. [[GH-9966](https://github.com/hashicorp/vault/pull/9966)]
* credential/aws: Added X-Amz-Content-Sha256 as a default STS request header [[GH-10009](https://github.com/hashicorp/vault/pull/10009)]

BUG FIXES:

* agent: Fix `disable_fast_negotiation` not being set on the auth method when configured by user. [[GH-9892](https://github.com/hashicorp/vault/pull/9892)]
* core (enterprise): Fix hang when cluster-wide plugin reload cleanup is slow on unseal
* core (enterprise): Fix an error in cluster-wide plugin reload cleanup following such a reload
* core: Fix crash when metrics collection encounters zero-length keys in KV store [[GH-9811](https://github.com/hashicorp/vault/pull/9881)]
* mfa (enterprise): Fix incorrect handling of PingID responses that could result in auth requests failing
* replication (enterprise): Improve race condition when using a newly created token on a performance standby node
* replication (enterprise): Only write failover cluster addresses if they've changed
* ui: fix bug where dropdown for identity/entity management is not reflective of actual policy [[GH-9958](https://github.com/hashicorp/vault/pull/9958)]

## 1.5.3

### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust
* secrets/ssh: Fixed a bug with role option for SSH signing algorithm to allow more than RSA signing

## 1.5.2.1

### August 21st, 2020

### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.5.2

### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.5.2 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.5.2)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.

## 1.5.1

### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)
* When using Vault Agent with cert auto-auth and caching enabled, under certain circumstances, clients without permission to access agent's token may retrieve the token without login credentials. This vulnerability affects Vault Agent 1.1.0 and newer and is fixed in 1.5.1 (CVE-2020-17455)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.5.1 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.5.1)

CHANGES:

* pki: The tidy operation will now remove revoked certificates if the parameter `tidy_revoked_certs` is set to `true`. This will result in certificate entries being immediately removed, as opposed to awaiting until its NotAfter time. Note that this only affects certificates that have been already revoked. [[GH-9609](https://github.com/hashicorp/vault/pull/9609)]
* go: Updated Go version to 1.14.7

IMPROVEMENTS:

* auth/jwt: Add support for fetching groups and user information from G Suite during authentication. [[GH-9574](https://github.com/hashicorp/vault/pull/9574)]
* auth/jwt: Add EdDSA to supported algorithms. [[GH-129](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/129)]
* secrets/openldap: Add "ad" schema that allows the engine to correctly rotate AD passwords. [[GH-9740](https://github.com/hashicorp/vault/pull/9740)]
* pki: Add a `allowed_domains_template` parameter that enables the use of identity templating within the `allowed_domains` parameter. [[GH-8509](https://github.com/hashicorp/vault/pull/8509)]
* secret/azure: Use write-ahead-logs to cleanup any orphaned Service Principals [[GH-9773](https://github.com/hashicorp/vault/pull/9773)]
* ui: Wrap TTL option on transit engine export action is updated to a new component. [[GH-9632](https://github.com/hashicorp/vault/pull/9632)]
* ui: Wrap Tool uses newest version of TTL Picker component. [[GH-9691](https://github.com/hashicorp/vault/pull/9691)]

BUG FIXES:

* secrets/gcp: Ensure that the IAM policy version is appropriately set after a roleset's bindings have changed. [[GH-9603](https://github.com/hashicorp/vault/pull/9603)]
* replication (enterprise): Fix status API output incorrectly stating replication is in `idle` state.
* replication (enterprise): Use PrimaryClusterAddr if it's been set
* core: Fix panic when printing over-long info fields at startup [[GH-9681](https://github.com/hashicorp/vault/pull/9681)]
* core: Seal migration using the new minimal-downtime strategy didn't work properly with performance standbys. [[GH-9690](https://github.com/hashicorp/vault/pull/9690)]
* core: Vault failed to start when there were non-string values in seal configuration [[GH-9555](https://github.com/hashicorp/vault/pull/9555)]
* core: Handle a trailing slash in the API address used for enabling replication

## 1.5.0

### July 21st, 2020

CHANGES:

* audit: Token TTL and issue time are now provided in the auth portion of audit logs. [[GH-9091](https://github.com/hashicorp/vault/pull/9091)]
* auth/gcp: Changes the default name of the entity alias that gets created to be the role ID for both IAM and GCE authentication. [[GH-99](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/99)]
* core: Remove the addition of newlines to parsed configuration when using integer/boolean values [[GH-8928](https://github.com/hashicorp/vault/pull/8928)]
* cubbyhole: Reject reads and writes to an empty ("") path. [[GH-8971](https://github.com/hashicorp/vault/pull/8971)]
* secrets/azure: Default password generation changed from uuid to cryptographically secure randomized string [[GH-40](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/40)]
* storage/gcs: The `credentials_file` config option has been removed. The `GOOGLE_APPLICATION_CREDENTIALS` environment variable
  or default credentials may be used instead [[GH-9424](https://github.com/hashicorp/vault/pull/9424)]
* storage/raft: The storage configuration now accepts a new `max_entry_size` config that will limit
  the total size in bytes of any entry committed via raft. It defaults to `"1048576"` (1MiB). [[GH-9027](https://github.com/hashicorp/vault/pull/9027)]
* token: Token creation with custom token ID via `id` will no longer allow periods (`.`) as part of the input string.
  The final generated token value may contain periods, such as the `s.` prefix for service token
  indication. [[GH-8646](https://github.com/hashicorp/vault/pull/8646/files)]
* token: Token renewals will now return token policies within the `token_policies` , identity policies within `identity_policies`, and the full policy set within `policies`. [[GH-8535](https://github.com/hashicorp/vault/pull/8535)]
* go: Updated Go version to 1.14.4

FEATURES:

* **Monitoring**: We have released a Splunk App [9] for Enterprise customers. The app is accompanied by an updated monitoring guide and a few new metrics to enable OSS users to effectively monitor Vault.
* **Password Policies**: Allows operators to customize how passwords are generated for select secret engines (OpenLDAP, Active Directory, Azure, and RabbitMQ).
* **Replication UI Improvements**: We have redesigned the replication UI to highlight the state and relationship between primaries and secondaries and improved management workflows, enabling a more holistic understanding of multiple Vault clusters.
* **Resource Quotas**: As of 1.5, Vault supports specifying a quota to rate limit requests on OSS and Enterprise. Enterprise customers also have access to set quotas on the number of leases that can be generated on a path.
* **OpenShift Support**: We have updated the Helm charts to allow users to install Vault onto their OpenShift clusters.
* **Seal Migration**: We have made updates to allow migrations from auto unseal to Shamir unseal on Enterprise.
* **AWS Auth Web Identity Support**: We've added support for AWS Web Identities, which will be used in the credentials chain if present.
* **Vault Monitor**: Similar to the monitor command for Consul and Nomad, we have added the ability for Vault to stream logs from other Vault servers at varying log levels.
* **AWS Secrets Groups Support**: IAM users generated by Vault may now be added to IAM Groups.
* **Integrated Storage as HA Storage**: In Vault 1.5, it is possible to use Integrated Storage as HA Storage with a different storage backend as regular storage.
* **OIDC Auth Provider Extensions**: We've added support to OIDC Auth to incorporate IdP-specific extensions. Currently this includes expanded Azure AD groups support.
* **GCP Secrets**: Support BigQuery dataset ACLs in absence of IAM endpoints.
* **KMIP**: Add support for signing client certificates requests (CSRs) rather than having them be generated entirely within Vault.

IMPROVEMENTS:

* audit: Replication status requests are no longer audited. [[GH-8877](https://github.com/hashicorp/vault/pull/8877)]
* audit: Added mount_type field to requests and responses. [[GH-9167](https://github.com/hashicorp/vault/pull/9167)]
* auth/aws: Add support for Web Identity credentials [[GH-7738](https://github.com/hashicorp/vault/pull/7738)]
* auth/jwt: Support users that are members of more than 200 groups on Azure [[GH-120](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/120)]
* auth/kerberos: Support identities without userPrincipalName [[GH-44](https://github.com/hashicorp/vault-plugin-auth-kerberos/issues/44)]
* auth/kubernetes: Allow disabling `iss` validation [[GH-91](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/91)]
* auth/kubernetes: Try reading the ca.crt and TokenReviewer JWT from the default service account [[GH-83](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/83)]
* cli: Support reading TLS parameters from file for the `vault operator raft join` command. [[GH-9060](https://github.com/hashicorp/vault/pull/9060)]
* cli: Add a new subcommand, `vault monitor`, for tailing server logs in the console. [[GH-8477](https://github.com/hashicorp/vault/pull/8477)]
* core: Add the Go version used to build a Vault binary to the server message output. [[GH-9078](https://github.com/hashicorp/vault/pull/9078)]
* core: Added Password Policies for user-configurable password generation [[GH-8637](https://github.com/hashicorp/vault/pull/8637)]
* core: New telemetry metrics covering token counts, token creation, KV secret counts, lease creation. [[GH-9239](https://github.com/hashicorp/vault/pull/9239)] [[GH-9250](https://github.com/hashicorp/vault/pull/9250)] [[GH-9244](https://github.com/hashicorp/vault/pull/9244)] [[GH-9052](https://github.com/hashicorp/vault/pull/9052)]
* physical/gcs: The storage backend now uses a dedicated client for HA lock updates to prevent lock table update failures when flooded by other client requests. [[GH-9424](https://github.com/hashicorp/vault/pull/9424)]
* physical/spanner: The storage backend now uses a dedicated client for HA lock updates to prevent lock table update failures when flooded by other client requests. [[GH-9423](https://github.com/hashicorp/vault/pull/9423)]
* plugin: Add SDK method, `Sys.ReloadPlugin`, and CLI command, `vault plugin reload`, for reloading plugins. [[GH-8777](https://github.com/hashicorp/vault/pull/8777)]
* plugin (enterprise): Add a scope field to plugin reload, which when global, reloads the plugin anywhere in a cluster. [[GH-9347](https://github.com/hashicorp/vault/pull/9347)]
* sdk/framework: Support accepting TypeFloat parameters over the API [[GH-8923](https://github.com/hashicorp/vault/pull/8923)]
* secrets/aws: Add iam_groups parameter to role create/update [[GH-8811](https://github.com/hashicorp/vault/pull/8811)]
* secrets/database: Add static role rotation for MongoDB Atlas database plugin [[GH-11](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/11)]
* secrets/database: Add static role rotation for MSSQL database plugin [[GH-9062](https://github.com/hashicorp/vault/pull/9062)]
* secrets/database: Allow InfluxDB to use insecure TLS without cert bundle [[GH-8778](https://github.com/hashicorp/vault/pull/8778)]
* secrets/gcp: Support BigQuery dataset ACLs in absence of IAM endpoints [[GH-78](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/78)]
* secrets/pki: Allow 3072-bit RSA keys [[GH-8343](https://github.com/hashicorp/vault/pull/8343)]
* secrets/ssh: Add a CA-mode role option to specify signing algorithm [[GH-9096](https://github.com/hashicorp/vault/pull/9096)]
* secrets/ssh: The [Vault SSH Helper](https://github.com/hashicorp/vault-ssh-helper) can now be configured to reference a mount in a namespace [[GH-44](https://github.com/hashicorp/vault-ssh-helper/pull/44)]
* secrets/transit: Transit requests that make use of keys now include a new field  `key_version` in their responses [[GH-9100](https://github.com/hashicorp/vault/pull/9100)]
* secrets/transit: Improving transit batch encrypt and decrypt latencies [[GH-8775](https://github.com/hashicorp/vault/pull/8775)]
* sentinel: Add a sentinel config section, and "additional_enabled_modules", a list of Sentinel modules that may be imported in addition to the defaults.
* ui: Update TTL picker styling on SSH secret engine [[GH-8891](https://github.com/hashicorp/vault/pull/8891)]
* ui: Only render the JWT input field of the Vault login form on mounts configured for JWT auth [[GH-8952](https://github.com/hashicorp/vault/pull/8952)]
* ui: Add replication dashboards.  Improve replication management workflows. [[GH-8705]](https://github.com/hashicorp/vault/pull/8705).
* ui: Update alert banners to match design systems black text. [[GH-9463]](https://github.com/hashicorp/vault/pull/9463).

BUG FIXES:

* auth/oci: Fix issue where users of the Oracle Cloud Infrastructure (OCI) auth method could not authenticate when the plugin backend was mounted at a non-default path. [[GH-7](https://github.com/hashicorp/vault-plugin-auth-oci/pull/7)]
* core: Extend replicated cubbyhole fix in 1.4.0 to cover case where a performance primary is also a DR primary [[GH-9148](https://github.com/hashicorp/vault/pull/9148)]
* replication (enterprise): Use the PrimaryClusterAddr if it's been set
* seal/awskms: fix AWS KMS auto-unseal when AWS_ROLE_SESSION_NAME not set [[GH-9416](https://github.com/hashicorp/vault/pull/9416)]
* sentinel: fix panic due to concurrent map access when rules iterate over metadata maps
* secrets/aws: Fix issue where performance standbys weren't able to generate STS credentials after an IAM access key rotation in AWS and root IAM credential update in Vault [[GH-9186](https://github.com/hashicorp/vault/pull/9186)]
* secrets/database: Fix issue where rotating root database credentials while Vault's storage backend is unavailable causes Vault to lose access to the database [[GH-8782](https://github.com/hashicorp/vault/pull/8782)]
* secrets/database: Fix issue that prevents performance standbys from connecting to databases after a root credential rotation [[GH-9129](https://github.com/hashicorp/vault/pull/9129)]
* secrets/database: Fix parsing of multi-line PostgreSQL statements [[GH-8512](https://github.com/hashicorp/vault/pull/8512)]
* secrets/gcp: Fix issue were updates were not being applied to the `token_scopes` of a roleset. [[GH-90](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/90)]
* secrets/kv: Return the value of delete_version_after when reading kv/config, even if it is set to the default. [[GH-42](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/42)]
* ui: Add Toggle component into core addon so it is available in KMIP and other Ember Engines.[[GH-8913]](https://github.com/hashicorp/vault/pull/8913)
* ui: Disallow max versions value of large than 9999999999999999 on kv2 secrets engine. [[GH-9242](https://github.com/hashicorp/vault/pull/9242)]
* ui: Add and upgrade missing dependencies to resolve a failure with `make static-dist`. [[GH-9277](https://github.com/hashicorp/vault/pull/9371)]

## 1.4.7.1

### October 15th, 2020

### Enterprise Only

BUG FIXES:

* replication (enterprise): Fix panic when old filter path evaluation fails

## 1.4.7

### September 24th, 2020

SECURITY:

* Batch Token Expiry: We addressed an issue where batch token leases could outlive their TTL because we were not scheduling the expiration time correctly. This vulnerability affects Vault OSS and Vault Enterprise 1.0 and newer and is fixed in 1.4.7 and 1.5.4 (CVE-2020-25816).

IMPROVEMENTS:

* secret/azure: Use write-ahead-logs to cleanup any orphaned Service Principals [[GH-9773](https://github.com/hashicorp/vault/pull/9773)]

BUG FIXES:

* replication (enterprise): Don't stop replication if old filter path evaluation fails

## 1.4.6

### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust
* secrets/ssh: Fixed a bug with role option for SSH signing algorithm to allow more than RSA signing [[GH-9824](https://github.com/hashicorp/vault/pull/9824)]

## 1.4.5.1

### August 21st, 2020

### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.4.5

### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.4.5 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.4.5)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.

## 1.4.4

### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.4.4 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.4.4)

BUG FIXES:

* auth/okta: fix bug introduced in 1.4.0: only 200 external groups were fetched even if user belonged to more [[GH-9580](https://github.com/hashicorp/vault/pull/9580)]
* seal/awskms: fix AWS KMS auto-unseal when AWS_ROLE_SESSION_NAME not set [[GH-9416](https://github.com/hashicorp/vault/pull/9416)]
* secrets/aws: Fix possible issue creating access keys when using Performance Standbys  [[GH-9606](https://github.com/hashicorp/vault/pull/9606)]

IMPROVEMENTS:

* auth/aws: Retry on transient failures during AWS IAM auth login attempts [[GH-8727](https://github.com/hashicorp/vault/pull/8727)]
* ui: Add transit key algorithms aes128-gcm96, ecdsa-p384, ecdsa-p521 to the UI. [[GH-9070](https://github.com/hashicorp/vault/pull/9070)] & [[GH-9520](https://github.com/hashicorp/vault/pull/9520)]

## 1.4.3

### July 2nd, 2020

IMPROVEMENTS:

* auth/aws: Add support for Web Identity credentials [[GH-9251](https://github.com/hashicorp/vault/pull/9251)]
* auth/kerberos: Support identities without userPrincipalName [[GH-44](https://github.com/hashicorp/vault-plugin-auth-kerberos/issues/44)]
* core: Add the Go version used to build a Vault binary to the server message output. [[GH-9078](https://github.com/hashicorp/vault/pull/9078)]
* secrets/database: Add static role rotation for MongoDB Atlas database plugin [[GH-9311](https://github.com/hashicorp/vault/pull/9311)]
* physical/mysql: Require TLS or plaintext flagging in MySQL configuration [[GH-9012](https://github.com/hashicorp/vault/pull/9012)]
* ui: Link to the Vault Changelog in the UI footer [[GH-9216](https://github.com/hashicorp/vault/pull/9216)]

BUG FIXES:

* agent: Restart template server when it shuts down [[GH-9200](https://github.com/hashicorp/vault/pull/9200)]
* auth/oci: Fix issue where users of the Oracle Cloud Infrastructure (OCI) auth method could not authenticate when the plugin backend was mounted at a non-default path. [[GH-9278](https://github.com/hashicorp/vault/pull/9278)]
* replication: The issue causing cubbyholes in namespaces on performance secondaries to not work, which was fixed in 1.4.0, was still an issue when the primary was both a performance primary and DR primary.
* seal: (enterprise) Fix issue causing stored seal and recovery keys to be mistaken as sealwrapped values
* secrets/aws: Fix issue where performance standbys weren't able to generate STS credentials after an IAM access key rotation in AWS and root IAM credential update in Vault [[GH-9207](https://github.com/hashicorp/vault/pull/9207)]
* secrets/database: Fix issue that prevents performance standbys from connecting to databases after a root credential rotation [[GH-9208](https://github.com/hashicorp/vault/pull/9208)]
* secrets/gcp: Fix issue were updates were not being applied to the `token_scopes` of a roleset. [[GH-9277](https://github.com/hashicorp/vault/pull/9277)]

## 1.4.2 (May 21st, 2020)

SECURITY:

* core: Proxy environment variables are now redacted before being logged, in case the URLs include a username:password. This vulnerability, CVE-2020-13223, is fixed in 1.3.6 and 1.4.2, but affects 1.4.0 and 1.4.1, as well as older versions of Vault [[GH-9022](https://github.com/hashicorp/vault/pull/9022)]
* secrets/gcp: Fix a regression in 1.4.0 where the system TTLs were being used instead of the configured backend TTLs for dynamic service accounts. This vulnerability is CVE-2020-12757. [[GH-85](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/85)]

IMPROVEMENTS:

* storage/raft: The storage stanza now accepts `leader_ca_cert_file`, `leader_client_cert_file`, and
  `leader_client_key_file` parameters to read and parse TLS certificate information from paths on disk.
  Existing non-path based parameters will continue to work, but their values will need to be provided as a
  single-line string with newlines delimited by `\n`.  [[GH-8894](https://github.com/hashicorp/vault/pull/8894)]
* storage/raft: The `vault status` CLI command and the `sys/leader` API now contain the committed and applied
  raft indexes. [[GH-9011](https://github.com/hashicorp/vault/pull/9011)]

BUG FIXES:

* auth/aws: Fix token renewal issues caused by the metadata changes in 1.4.1 [[GH-8991](https://github.com/hashicorp/vault/pull/8991)]
* auth/ldap: Fix 1.4.0 regression that could result in auth failures when LDAP auth config includes upndomain. [[GH-9041](https://github.com/hashicorp/vault/pull/9041)]
* secrets/ad: Forward rotation requests from standbys to active clusters [[GH-66](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/66)]
* secrets/database: Prevent generation of usernames that are not allowed by the MongoDB Atlas API [[GH-9](https://github.com/hashicorp/vault-plugin-database-mongodbatlas/pull/9)]
* secrets/database: Return an error if a manual rotation of static account credentials fails [[GH-9035](https://github.com/hashicorp/vault/pull/9035)]
* secrets/openldap: Forward all rotation requests from standbys to active clusters [[GH-9028](https://github.com/hashicorp/vault/pull/9028)]
* secrets/transform (enterprise): Fix panic that could occur when accessing cached template entries, such as a requests
  that accessed templates directly or indirectly from a performance standby node.
* serviceregistration: Fix a regression for Consul service registration that ignored using the listener address as
  the redirect address unless api_addr was provided. It now properly uses the same redirect address as the one
  used by Vault's Core object. [[GH-8976](https://github.com/hashicorp/vault/pull/8976)]
* storage/raft: Advertise the configured cluster address to the rest of the nodes in the raft cluster. This fixes
  an issue where a node advertising 0.0.0.0 is not using a unique hostname. [[GH-9008](https://github.com/hashicorp/vault/pull/9008)]
* storage/raft: Fix panic when multiple nodes attempt to join the cluster at once. [[GH-9008](https://github.com/hashicorp/vault/pull/9008)]
* sys: The path provided in `sys/internal/ui/mounts/:path` is now namespace-aware. This fixes an issue
  with `vault kv` subcommands that had namespaces provided in the path returning permission denied all the time.
  [[GH-8962](https://github.com/hashicorp/vault/pull/8962)]
* ui: Fix snowman that appears when namespaces have more than one period [[GH-8910](https://github.com/hashicorp/vault/pull/8910)]

## 1.4.1 (April 30th, 2020)

CHANGES:

* auth/aws: The default set of metadata fields added in 1.4.1 has been changed to `account_id` and `auth_type` [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]
* storage/raft: Disallow `ha_storage` to be specified if `raft` is set as the `storage` type. [[GH-8707](https://github.com/hashicorp/vault/pull/8707)]

IMPROVEMENTS:

* auth/aws: The set of metadata stored during login is now configurable [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]
* auth/aws: Improve region selection to avoid errors seen if the account hasn't enabled some newer AWS regions [[GH-8679](https://github.com/hashicorp/vault/pull/8679)]
* auth/azure: Enable login from Azure VMs with user-assigned identities [[GH-33](https://github.com/hashicorp/vault-plugin-auth-azure/pull/33)]
* auth/gcp: The set of metadata stored during login is now configurable [[GH-92](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/92)]
* auth/gcp: The type of alias name used during login is now configurable [[GH-95](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/95)]
* auth/ldap: Improve error messages during LDAP operation failures [[GH-8740](https://github.com/hashicorp/vault/pull/8740)]
* identity: Add a batch delete API for identity entities [[GH-8785]](https://github.com/hashicorp/vault/pull/8785)
* identity: Improve performance of logins when no group updates are needed [[GH-8795]](https://github.com/hashicorp/vault/pull/8795)
* metrics: Add `vault.identity.num_entities` metric [[GH-8816]](https://github.com/hashicorp/vault/pull/8816)
* secrets/kv: Allow `delete-version-after` to be reset to 0 via the CLI [[GH-8635](https://github.com/hashicorp/vault/pull/8635)]
* secrets/rabbitmq: Improve error handling and reporting [[GH-8619](https://github.com/hashicorp/vault/pull/8619)]
* ui: Provide One Time Password during Operation Token generation process [[GH-8630]](https://github.com/hashicorp/vault/pull/8630)

BUG FIXES:

* auth/okta: Fix MFA regression (introduced in [GH-8143](https://github.com/hashicorp/vault/pull/8143)) from 1.4.0 [[GH-8807](https://github.com/hashicorp/vault/pull/8807)]
* auth/userpass: Fix upgrade value for `token_bound_cidrs` being ignored due to incorrect key provided [[GH-8826](https://github.com/hashicorp/vault/pull/8826/files)]
* config/seal: Fix segfault when seal block is removed [[GH-8517](https://github.com/hashicorp/vault/pull/8517)]
* core: Fix an issue where users attempting to build Vault could receive Go module checksum errors [[GH-8770](https://github.com/hashicorp/vault/pull/8770)]
* core: Fix blocked requests if a SIGHUP is issued during a long-running request has the state lock held.
  Also fixes deadlock that can happen if `vault debug` with the config target is ran during this time.
  [[GH-8755](https://github.com/hashicorp/vault/pull/8755)]
* core: Always rewrite the .vault-token file as part of a `vault login` to ensure permissions and ownership are set correctly [[GH-8867](https://github.com/hashicorp/vault/pull/8867)]
* database/mongodb: Fix context deadline error that may result due to retry attempts on failed commands
  [[GH-8863](https://github.com/hashicorp/vault/pull/8863)]
* http: Fix superflous call messages from the http package on logs caused by missing returns after
  `respondError` calls [[GH-8796](https://github.com/hashicorp/vault/pull/8796)]
* namespace (enterprise): Fix namespace listing to return `key_info` when a scoping namespace is also provided.
* seal/gcpkms: Fix panic that could occur if all seal parameters were provided via environment
  variables [[GH-8840](https://github.com/hashicorp/vault/pull/8840)]
* storage/raft: Fix memory allocation and incorrect metadata tracking issues with snapshots [[GH-8793](https://github.com/hashicorp/vault/pull/8793)]
* storage/raft: Fix panic that could occur if `disable_clustering` was set to true on Raft storage cluster [[GH-8784](https://github.com/hashicorp/vault/pull/8784)]
* storage/raft: Handle errors returned from the API during snapshot operations [[GH-8861](https://github.com/hashicorp/vault/pull/8861)]
* sys/wrapping: Allow unwrapping of wrapping tokens which contain nil data [[GH-8714](https://github.com/hashicorp/vault/pull/8714)]

## 1.4.0 (April 7th, 2020)

CHANGES:

* cli: The raft configuration command has been renamed to list-peers to avoid
  confusion.

FEATURES:

* **Kerberos Authentication**: Vault now supports Kerberos authentication using a SPNEGO token.
   Login can be performed using the Vault CLI, API, or agent.
* **Kubernetes Service Discovery**: A new Kubernetes service discovery feature where, if
   configured, Vault will tag Vault pods with their current health status. For more, see [#8249](https://github.com/hashicorp/vault/pull/8249).
* **MongoDB Atlas Secrets**: Vault can now generate dynamic credentials for both MongoDB Atlas databases
  as well as the [Atlas programmatic interface](https://docs.atlas.mongodb.com/tutorial/manage-programmatic-access/).
* **OpenLDAP Secrets Engine**: We now support password management of existing OpenLDAP user entries. For more, see [#8360](https://github.com/hashicorp/vault/pull/8360/).
* **Redshift Database Secrets Engine**: The database secrets engine now supports static and dynamic secrets for the Amazon Web Services (AWS) Redshift service.
* **Service Registration Config**: A newly introduced `service_registration` configuration stanza, that allows for service registration to be configured separately from the storage backend. For more, see [#7887](https://github.com/hashicorp/vault/pull/7887/).
* **Transform Secrets Engine (Enterprise)**: A new secrets engine that handles secure data transformations against provided input values.
* **Integrated Storage**: Promoted out of beta and into general availability for both open-source and enterprise workloads.

IMPROVEMENTS:

* agent: add option to force the use of the auth-auth token, and ignore the Vault token in the request [[GH-8101](https://github.com/hashicorp/vault/pull/8101)]
* api: Restore and fix DNS SRV Lookup [[GH-8520](https://github.com/hashicorp/vault/pull/8520)]
* audit: HMAC http_raw_body in audit log; this ensures that large authenticated Prometheus metrics responses get
  replaced with short HMAC values [[GH-8130](https://github.com/hashicorp/vault/pull/8130)]
* audit: Generate-root, generate-recovery-token, and generate-dr-operation-token requests and responses are now audited. [[GH-8301](https://github.com/hashicorp/vault/pull/8301)]
* auth/aws: Reduce the number of simultaneous STS client credentials needed  [[GH-8161](https://github.com/hashicorp/vault/pull/8161)]
* auth/azure: subscription ID, resource group, vm and vmss names are now stored in alias metadata [[GH-30](https://github.com/hashicorp/vault-plugin-auth-azure/pull/30)]
* auth/jwt: Additional OIDC callback parameters available for CLI logins [[GH-80](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/80) & [GH-86](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/86)]
* auth/jwt: Bound claims may be optionally configured using globs [[GH-89](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/89)]
* auth/jwt: Timeout during OIDC CLI login if process doesn't complete within 2 minutes [[GH-97](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/97)]
* auth/jwt: Add support for the `form_post` response mode [[GH-98](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/98)]
* auth/jwt: add optional client_nonce to authorization flow [[GH-104](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/104)]
* auth/okta: Upgrade okta sdk lib, which should improve handling of groups [[GH-8143](https://github.com/hashicorp/vault/pull/8143)]
* aws: Add support for v2 of the instance metadata service (see [issue 7924](https://github.com/hashicorp/vault/issues/7924) for all linked PRs)
* core: Separate out service discovery interface from storage interface to allow
  new types of service discovery not coupled to storage [[GH-7887](https://github.com/hashicorp/vault/pull/7887)]
* core: Add support for telemetry option `metrics_prefix` [[GH-8340](https://github.com/hashicorp/vault/pull/8340)]
* core: Entropy Augmentation can now be used with AWS KMS and Vault Transit seals
* core: Allow tls_min_version to be set to TLS 1.3 [[GH-8305](https://github.com/hashicorp/vault/pull/8305)]
* cli: Incorrect TLS configuration will now correctly fail [[GH-8025](https://github.com/hashicorp/vault/pull/8025)]
* identity: Allow specifying a custom `client_id` for identity tokens [[GH-8165](https://github.com/hashicorp/vault/pull/8165)]
* metrics/prometheus: improve performance with high volume of metrics updates [[GH-8507](https://github.com/hashicorp/vault/pull/8507)]
* replication (enterprise): Fix race condition causing clusters with high throughput writes to sometimes
  fail to enter streaming-wal mode
* replication (enterprise): Secondary clusters can now perform an extra gRPC call to all nodes in a primary
  cluster in an attempt to resolve the active node's address
* replication (enterprise): The replication status API now outputs `last_performance_wal`, `last_dr_wal`,
  and `connection_state` values
* replication (enterprise): DR secondary clusters can now be recovered by the `replication/dr/secondary/recover`
  API
* replication (enterprise): We now allow for an alternate means to create a Disaster Recovery token, by using a batch
  token that is created with an ACL that allows for access to one or more of the DR endpoints.
* secrets/database/mongodb: Switched internal MongoDB driver to mongo-driver [[GH-8140](https://github.com/hashicorp/vault/pull/8140)]
* secrets/database/mongodb: Add support for x509 client authorization to MongoDB [[GH-8329](https://github.com/hashicorp/vault/pull/8329)]
* secrets/database/oracle: Add support for static credential rotation [[GH-26](https://github.com/hashicorp/vault-plugin-database-oracle/pull/26)]
* secrets/consul: Add support to specify TLS options per Consul backend [[GH-4800](https://github.com/hashicorp/vault/pull/4800)]
* secrets/gcp: Allow specifying the TTL for a service key [[GH-54](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/54)]
* secrets/gcp: Add support for rotating root keys [[GH-53](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/53)]
* secrets/gcp: Handle version 3 policies for Resource Manager IAM requests [[GH-77](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/77)]
* secrets/nomad: Add support to specify TLS options per Nomad backend [[GH-8083](https://github.com/hashicorp/vault/pull/8083)]
* secrets/ssh: Allowed users can now be templated with identity information [[GH-7548](https://github.com/hashicorp/vault/pull/7548)]
* secrets/transit: Adding RSA3072 key support [[GH-8151](https://github.com/hashicorp/vault/pull/8151)]
* storage/consul: Vault returns now a more descriptive error message when only a client cert or
  a client key has been provided [[GH-4930]](https://github.com/hashicorp/vault/pull/8084)
* storage/raft: Nodes in the raft cluster can all be given possible leader
  addresses for them to continuously try and join one of them, thus automating
  the process of join to a greater extent [[GH-7856](https://github.com/hashicorp/vault/pull/7856)]
* storage/raft: Fix a potential deadlock that could occur on leadership transition [[GH-8547](https://github.com/hashicorp/vault/pull/8547)]
* storage/raft: Refresh TLS keyring on snapshot restore [[GH-8546](https://github.com/hashicorp/vault/pull/8546)]
* storage/etcd: Bumped etcd client API SDK [[GH-7931](https://github.com/hashicorp/vault/pull/7931) & [GH-4961](https://github.com/hashicorp/vault/pull/4961) & [GH-4349](https://github.com/hashicorp/vault/pull/4349) & [GH-7582](https://github.com/hashicorp/vault/pull/7582)]
* ui: Make Transit Key actions more prominent [[GH-8304](https://github.com/hashicorp/vault/pull/8304)]
* ui: Add Core Usage Metrics [[GH-8347](https://github.com/hashicorp/vault/pull/8347)]
* ui: Add refresh Namespace list on the Namespace dropdown, and redesign of Namespace dropdown menu [[GH-8442](https://github.com/hashicorp/vault/pull/8442)]
* ui: Update transit actions to codeblocks & automatically encode plaintext unless indicated [[GH-8462](https://github.com/hashicorp/vault/pull/8462)]
* ui: Display the results of transit key actions in a modal window [[GH-8462](https://github.com/hashicorp/vault/pull/8575)]
* ui: Transit key version styling updates & ability to copy key from dropdown [[GH-8480](https://github.com/hashicorp/vault/pull/8480)]

BUG FIXES:

* agent: Fix issue where TLS options are ignored for agent template feature [[GH-7889](https://github.com/hashicorp/vault/pull/7889)]
* auth/jwt: Use lower case role names for `default_role` to match the `role` case convention [[GH-100](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/100)]
* auth/ldap: Fix a bug where the UPNDOMAIN parameter was wrongly used to lookup the group
  membership of the given user [[GH-6325]](https://github.com/hashicorp/vault/pull/8333)
* cli: Support autocompletion for nested mounts [[GH-8303](https://github.com/hashicorp/vault/pull/8303)]
* cli: Fix CLI namespace autocompletion [[GH-8315](https://github.com/hashicorp/vault/pull/8315)]
* identity: Fix incorrect caching of identity token JWKS responses [[GH-8412](https://github.com/hashicorp/vault/pull/8412)]
* metrics/stackdriver: Fix issue that prevents the stackdriver metrics library to create unnecessary stackdriver descriptors [[GH-8073](https://github.com/hashicorp/vault/pull/8073)]
* replication (enterprise): Fix issue causing cubbyholes in namespaces on performance secondaries to not work.
* replication (enterprise): Unmounting a dynamic secrets backend could sometimes lead to replication errors.  Change the order of operations to prevent that.
* seal (enterprise): Fix seal migration when transactional seal wrap backend is in use.
* secrets/database/influxdb: Fix potential panic if connection to the InfluxDB database cannot be established [[GH-8282](https://github.com/hashicorp/vault/pull/8282)]
* secrets/database/mysql: Ensures default static credential rotation statements are used [[GH-8240](https://github.com/hashicorp/vault/pull/8240)]
* secrets/database/mysql: Fix inconsistent query parameter names: {{name}} or {{username}} for
  different queries. Now it allows for either for backwards compatibility [[GH-8240](https://github.com/hashicorp/vault/pull/8240)]
* secrets/database/postgres: Fix inconsistent query parameter names: {{name}} or {{username}} for
  different queries. Now it allows for either for backwards compatibility [[GH-8240](https://github.com/hashicorp/vault/pull/8240)]
* secrets/pki: Support FQDNs in DNS Name [[GH-8288](https://github.com/hashicorp/vault/pull/8288)]
* storage/raft: Allow seal migration to be performed on Vault clusters using raft storage [[GH-8103](https://github.com/hashicorp/vault/pull/8103)]
* telemetry: Prometheus requests on standby nodes will now return an error instead of forwarding
  the request to the active node [[GH-8280](https://github.com/hashicorp/vault/pull/8280)]
* ui: Fix broken popup menu on the transit secrets list page [[GH-8348](https://github.com/hashicorp/vault/pull/8348)]
* ui: Update headless Chrome flag to fix `yarn run test:oss` [[GH-8035](https://github.com/hashicorp/vault/pull/8035)]
* ui: Update CLI to accept empty strings as param value to reset previously-set values
* ui: Fix bug where error states don't clear when moving between action tabs on Transit [[GH-8354](https://github.com/hashicorp/vault/pull/8354)]

## 1.3.10

### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust

## 1.3.9.1

### August 21st, 2020

### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.3.9

### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.3.9 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.3.9)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.

## 1.3.8

### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.3.8 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.3.8)

## 1.3.7

### July 2nd, 2020

BUG FIXES:

* seal: (enterprise) Fix issue causing stored seal and recovery keys to be mistaken as sealwrapped values
* secrets/aws: Fix issue where performance standbys weren't able to generate STS credentials after an IAM access key rotation in AWS and root IAM credential update in Vault [[GH-9363](https://github.com/hashicorp/vault/pull/9363)]

## 1.3.6 (May 21st, 2020)

SECURITY:

* core: proxy environment variables are now redacted before being logged, in case the URLs include a username:password. This vulnerability, CVE-2020-13223, is fixed in 1.3.6 and 1.4.2, but affects 1.4 and 1.4.1, as well as older versions of Vault [[GH-9022](https://github.com/hashicorp/vault/pull/9022)]

BUG FIXES:

* auth/aws: Fix token renewal issues caused by the metadata changes in 1.3.5 [[GH-8991](https://github.com/hashicorp/vault/pull/8991)]
* replication: Fix mount filter bug that allowed replication filters to hide local mounts on a performance secondary

## 1.3.5 (April 28th, 2020)

CHANGES:

* auth/aws: The default set of metadata fields added in 1.3.2 has been changed to `account_id` and `auth_type` [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]

IMPROVEMENTS:

* auth/aws: The set of metadata stored during login is now configurable [[GH-8783](https://github.com/hashicorp/vault/pull/8783)]

## 1.3.4 (March 19th, 2020)

SECURITY:

* A vulnerability was identified in Vault and Vault Enterprise such that, under certain circumstances,  an Entity's Group membership may inadvertently include Groups the Entity no longer has permissions to. This vulnerability, CVE-2020-10660, affects Vault and Vault Enterprise versions 0.9.0 and newer, and is fixed in 1.3.4. [[GH-8606](https://github.com/hashicorp/vault/pull/8606)]
* A vulnerability was identified in Vault Enterprise such that, under certain circumstances, existing nested-path policies may give access to Namespaces created after-the-fact. This vulnerability, CVE-2020-10661, affects Vault Enterprise versions 0.11 and newer, and is fixed in 1.3.4.

## 1.3.3 (March 5th, 2020)

BUG FIXES:

* approle: Fix excessive locking during tidy, which could potentially block new approle logins for long enough to cause an outage [[GH-8418](https://github.com/hashicorp/vault/pull/8418)]
* cli: Fix issue where Raft snapshots from standby nodes created an empty backup file [[GH-8097](https://github.com/hashicorp/vault/pull/8097)]
* identity: Fix incorrect caching of identity token JWKS responses [[GH-8412](https://github.com/hashicorp/vault/pull/8412)]
* kmip: role read now returns tls_client_ttl
* kmip: fix panic when templateattr not provided in rekey request
* secrets/database/influxdb: Fix potential panic if connection to the InfluxDB database cannot be established [[GH-8282](https://github.com/hashicorp/vault/pull/8282)]
* storage/mysql: Fix potential crash when using MySQL as coordination for high availability [[GH-8300](https://github.com/hashicorp/vault/pull/8300)]
* storage/raft: Fix potential crash when using Raft as coordination for high availability [[GH-8356](https://github.com/hashicorp/vault/pull/8356)]
* ui: Fix missing License menu item [[GH-8230](https://github.com/hashicorp/vault/pull/8230)]
* ui: Fix bug where default auth method on login is defaulted to auth method that is listing-visibility=unauth instead of “other” [[GH-8218](https://github.com/hashicorp/vault/pull/8218)]
* ui: Fix bug where KMIP details were not shown in the UI Wizard [[GH-8255](https://github.com/hashicorp/vault/pull/8255)]
* ui: Show Error messages on Auth Configuration page when you hit permission errors [[GH-8500](https://github.com/hashicorp/vault/pull/8500)]
* ui: Remove duplicate form inputs for the GitHub config [[GH-8519](https://github.com/hashicorp/vault/pull/8519)]
* ui: Correct HMAC capitalization [[GH-8528](https://github.com/hashicorp/vault/pull/8528)]
* ui: Fix danger message in DR [[GH-8555](https://github.com/hashicorp/vault/pull/8555)]
* ui: Fix certificate field for LDAP config [[GH-8573](https://github.com/hashicorp/vault/pull/8573)]

## 1.3.2 (January 22nd, 2020)

SECURITY:

* When deleting a namespace on Vault Enterprise, in certain circumstances, the deletion
   process will fail to revoke dynamic secrets for a mount in that namespace. This will
   leave any dynamic secrets in remote systems alive and will fail to clean them up. This
   vulnerability, CVE-2020-7220, affects Vault Enterprise 0.11.0 and newer.

IMPROVEMENTS:

* auth/aws: Add aws metadata to identity alias [[GH-7985](https://github.com/hashicorp/vault/pull/7985)]
* auth/kubernetes: Allow both names and namespaces to be set to "*" [[GH-78](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/78)]

BUG FIXES:

* auth/azure: Fix Azure compute client to use correct base URL [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* auth/ldap: Fix renewal of tokens without configured policies that are
  generated by an LDAP login [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* auth/okta: Fix renewal of tokens without configured policies that are
  generated by an Okta login [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* core: Fix seal migration error when attempting to migrate from auto unseal to shamir [[GH-8172](https://github.com/hashicorp/vault/pull/8172)]
* core: Fix seal migration config issue when migrating from auto unseal to auto unseal [[GH-8172](https://github.com/hashicorp/vault/pull/8172)]
* plugin: Fix issue where a plugin unwrap request potentially used an expired token [[GH-8058](https://github.com/hashicorp/vault/pull/8058)]
* replication: Fix issue where a forwarded request from a performance/standby node could run into
  a timeout
* secrets/database: Fix issue where a manual static role rotation could potentially panic [[GH-8098](https://github.com/hashicorp/vault/pull/8098)]
* secrets/database: Fix issue where a manual root credential rotation request is not forwarded
  to the primary node [[GH-8125](https://github.com/hashicorp/vault/pull/8125)]
* secrets/database: Fix issue where a manual static role rotation request is not forwarded
  to the primary node [[GH-8126](https://github.com/hashicorp/vault/pull/8126)]
* secrets/database/mysql: Fix issue where special characters for a MySQL password were encoded [[GH-8040](https://github.com/hashicorp/vault/pull/8040)]
* ui: Fix deleting namespaces [[GH-8132](https://github.com/hashicorp/vault/pull/8132)]
* ui: Fix Error handler on kv-secret edit and kv-secret view pages [[GH-8133](https://github.com/hashicorp/vault/pull/8133)]
* ui: Fix OIDC callback to check storage [[GH-7929](https://github.com/hashicorp/vault/pull/7929)].
* ui: Change `.box-radio` height to min-height to prevent overflow issues [[GH-8065](https://github.com/hashicorp/vault/pull/8065)]

## 1.3.1 (December 18th, 2019)

IMPROVEMENTS:

* agent: Add ability to set `exit-after-auth` via the CLI [[GH-7920](https://github.com/hashicorp/vault/pull/7920)]
* auth/ldap: Add a `request_timeout` configuration option to prevent connection
  requests from hanging [[GH-7909](https://github.com/hashicorp/vault/pull/7909)]
* auth/kubernetes: Add audience to tokenreview API request for Kube deployments where issuer
  is not Kube. [[GH-74](https://github.com/hashicorp/vault/pull/74)]
* secrets/ad: Add a `request_timeout` configuration option to prevent connection
  requests from hanging [[GH-59](https://github.com/hashicorp/vault-plugin-secrets-ad/pull/59)]
* storage/postgresql: Add support for setting `connection_url` from enviornment
  variable `VAULT_PG_CONNECTION_URL` [[GH-7937](https://github.com/hashicorp/vault/pull/7937)]
* telemetry: Add `enable_hostname_label` option to telemetry stanza [[GH-7902](https://github.com/hashicorp/vault/pull/7902)]
* telemetry: Add accept header check for prometheus mime type [[GH-7958](https://github.com/hashicorp/vault/pull/7958)]

BUG FIXES:

* agent: Fix issue where Agent exits before all templates are rendered when
  using and `exit_after_auth` [[GH-7899](https://github.com/hashicorp/vault/pull/7899)]
* auth/aws: Fixes region-related issues when using a custom `sts_endpoint` by adding
  a `sts_region` parameter [[GH-7922](https://github.com/hashicorp/vault/pull/7922)]
* auth/token: Fix panic when getting batch tokens on a performance standby from a role
  that does not exist [[GH-8027](https://github.com/hashicorp/vault/pull/8027)]
* core: Improve warning message for lease TTLs [[GH-7901](https://github.com/hashicorp/vault/pull/7901)]
* identity: Fix identity token panic during invalidation [[GH-8043](https://github.com/hashicorp/vault/pull/8043)]
* plugin: Fix a panic that could occur if a mount/auth entry was unable to
  mount the plugin backend and a request that required the system view to be
  retrieved was made [[GH-7991](https://github.com/hashicorp/vault/pull/7991)]
* replication: Add `generate-public-key` endpoint to list of allowed endpoints
  for existing DR secondaries
* secrets/gcp: Fix panic if bindings aren't provided in roleset create/update. [[GH-56](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/56)]
* secrets/pki: Prevent generating certificate on performance standby when storing
  [[GH-7904](https://github.com/hashicorp/vault/pull/7904)]
* secrets/transit: Prevent restoring keys to new names that are sub paths [[GH-7998](https://github.com/hashicorp/vault/pull/7998)]
* storage/s3: Fix a bug in configurable S3 paths that was preventing use of S3 as
  a source during `operator migrate` operations [[GH-7966](https://github.com/hashicorp/vault/pull/7966)]
* ui: Ensure secrets with a period in their key can be viewed and copied [[GH-7926](https://github.com/hashicorp/vault/pull/7926)]
* ui: Fix status menu after demotion [[GH-7997](https://github.com/hashicorp/vault/pull/7997)]
* ui: Fix select dropdowns in Safari when running Mojave [[GH-8023](https://github.com/hashicorp/vault/pull/8023)]

## 1.3 (November 14th, 2019)

CHANGES:

* Secondary cluster activation: There has been a change to the way that activating
   performance and DR secondary clusters works when using public keys for
   encryption of the parameters rather than a wrapping token. This flow was
   experimental and never documented. It is now officially supported and
   documented but is not backwards compatible with older Vault releases.
* Cluster cipher suites: On its cluster port, Vault will no longer advertise
   the full TLS 1.2 cipher suite list by default. Although this port is only
   used for Vault-to-Vault communication and would always pick a strong cipher,
   it could cause false flags on port scanners and other security utilities
   that assumed insecure ciphers were being used. The previous behavior can be
   achieved by setting the value of the (undocumented) `cluster_cipher_suites`
   config flag to `tls12`.
* API/Agent Renewal behavior: The API now allows multiple options for how it
   deals with renewals. The legacy behavior in the Agent/API is for the renewer
   (now called the lifetime watcher) to exit on a renew error, leading to a
   reauthentication. The new default behavior is for the lifetime watcher to
   ignore 5XX errors and simply retry as scheduled, using the existing lease
   duration. It is also possible, within custom code, to disable renewals
   entirely, which allows the lifetime watcher to simply return when it
   believes it is time for your code to renew or reauthenticate.

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
* **Vault Agent Template**: Vault Agent now supports rendering templates containing
   Vault secrets to disk, similar to Consul Template [[GH-7652](https://github.com/hashicorp/vault/pull/7652)]
* **Transit Key Type Support**: Signing and verification is now supported with the P-384
   (secp384r1) and P-521 (secp521r1) ECDSA curves [[GH-7551](https://github.com/hashicorp/vault/pull/7551)] and encryption and
   decryption is now supported via AES128-GCM96 [[GH-7555](https://github.com/hashicorp/vault/pull/7555)]
* **SSRF Protection for Vault Agent**: Vault Agent has a configuration option to
   require a specific header before allowing requests [[GH-7627](https://github.com/hashicorp/vault/pull/7627)]
* **AWS Auth Method Root Rotation**: The credential used by the AWS auth method can
   now be rotated, to ensure that only Vault knows the credentials it is using [[GH-7131](https://github.com/hashicorp/vault/pull/7131)]
* **New UI Features**: The UI now supports managing users and groups for the
   Userpass, Cert, Okta, and Radius auth methods.
* **Shamir with Stored Master Key**: The on disk format for Shamir seals has changed,
   allowing for a secondary cluster using Shamir downstream from a primary cluster
   using Auto Unseal. [[GH-7694](https://github.com/hashicorp/vault/pull/7694)]
* **Stackdriver Metrics Sink**: Vault can now send metrics to
   [Stackdriver](https://cloud.google.com/stackdriver/). See the [configuration
   documentation](https://www.vaultproject.io/docs/config/index.html) for
   details. [[GH-6957](https://github.com/hashicorp/vault/pull/6957)]
* **Filtered Paths Replication (Enterprise)**: Based on the predecessor Filtered Mount Replication,
   Filtered Paths Replication allows now filtering of namespaces in addition to mounts.
   With this feature, Filtered Mount Replication should be considered deprecated.
* **Token Renewal via Accessor**: Tokens can now be renewed via the accessor value through
   the new `auth/token/renew-accessor` endpoint if the caller's token has
   permission to access that endpoint.
* **Improved Integrated Storage (Beta)**: Improved raft write performance, added support for
   non-voter nodes, along with UI support for: using raft storage, joining a raft cluster,
   and downloading and restoring a snapshot.

IMPROVEMENTS:

* agent: Add ability to set the TLS SNI name used by Agent [[GH-7519](https://github.com/hashicorp/vault/pull/7519)]
* agent & api: Change default renewer behavior to ignore 5XX errors [[GH-7733](https://github.com/hashicorp/vault/pull/7733)]
* auth/jwt: The redirect callback host may now be specified for CLI logins
   [[GH-71](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/71)]
* auth/jwt: Bound claims may now contain boolean values [[GH-73](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/73)]
* auth/jwt: CLI logins can now open the browser when running in WSL [[GH-77](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/77)]
* core: Exit ScanView if context has been cancelled [[GH-7419](https://github.com/hashicorp/vault/pull/7419)]
* core: re-encrypt barrier and recovery keys if the unseal key is updated
   [[GH-7493](https://github.com/hashicorp/vault/pull/7493)]
* core: Don't advertise the full set of TLS 1.2 cipher suites on the cluster
   port, even though only strong ciphers were used [[GH-7487](https://github.com/hashicorp/vault/pull/7487)]
* core (enterprise): Add background seal re-wrap
* core/metrics: Add config parameter to allow unauthenticated sys/metrics
   access. [[GH-7550](https://github.com/hashicorp/vault/pull/7550)]
* metrics: Upgrade DataDog library to improve performance [[GH-7794](https://github.com/hashicorp/vault/pull/7794)]
* replication (enterprise): Write-Ahead-Log entries will not duplicate the
   data belonging to the encompassing physical entries of the transaction,
   thereby improving the performance and storage capacity.
* replication (enterprise): Added more replication metrics
* replication (enterprise): Reindex process now compares subpages for a more
   accurate indexing process.
* replication (enterprise): Reindex API now accepts a new `skip_flush`
   parameter indicating all the changes should not be flushed while the tree is
   locked.
* secrets/aws: The root config can now be read [[GH-7245](https://github.com/hashicorp/vault/pull/7245)]
* secrets/aws: Role paths may now contain the '@' character [[GH-7553](https://github.com/hashicorp/vault/pull/7553)]
* secrets/database/cassandra: Add ability to skip verfication of connection
   [[GH-7614](https://github.com/hashicorp/vault/pull/7614)]
* secrets/gcp: Fix panic during rollback if the roleset has been deleted
   [[GH-52](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/52)]
* storage/azure: Add config parameter to Azure storage backend to allow
   specifying the ARM endpoint [[GH-7567](https://github.com/hashicorp/vault/pull/7567)]
* storage/cassandra: Improve storage efficiency by eliminating unnecessary
   copies of value data [[GH-7199](https://github.com/hashicorp/vault/pull/7199)]
* storage/raft: Improve raft write performance by utilizing FSM Batching
   [[GH-7527](https://github.com/hashicorp/vault/pull/7527)]
* storage/raft: Add support for non-voter nodes [[GH-7634](https://github.com/hashicorp/vault/pull/7634)]
* sys: Add a new `sys/host-info` endpoint for querying information about
   the host [[GH-7330](https://github.com/hashicorp/vault/pull/7330)]
* sys: Add a new set of endpoints under `sys/pprof/` that allows profiling
   information to be extracted [[GH-7473](https://github.com/hashicorp/vault/pull/7473)]
* sys: Add endpoint that counts the total number of active identity entities
   [[GH-7541](https://github.com/hashicorp/vault/pull/7541)]
* sys: `sys/seal-status` now has a `storage_type` field denoting what type of
   storage
   the cluster is configured to use
* sys: Add a new `sys/internal/counters/tokens` endpoint, that counts the
   total number of active service token accessors in the shared token storage.
   [[GH-7541](https://github.com/hashicorp/vault/pull/7541)]
* sys/config: Add  a new endpoint under `sys/config/state/sanitized` that
   returns the configuration state of the server. It excludes config values
   from `storage`, `ha_storage`, and `seal` stanzas and some values
   from `telemetry` due to potential sensitive entries in those fields.
* ui: when using raft storage, you can now join a raft cluster, download a
   snapshot, and restore a snapshot from the UI [[GH-7410](https://github.com/hashicorp/vault/pull/7410)]
* ui: clarify when secret version is deleted in the secret version history
   dropdown [[GH-7714](https://github.com/hashicorp/vault/pull/7714)]

BUG FIXES:

* agent: Fix a data race on the token value for inmemsink [[GH-7707](https://github.com/hashicorp/vault/pull/7707)]
* api: Fix Go API using lease revocation via URL instead of body [[GH-7777](https://github.com/hashicorp/vault/pull/7777)]
* api: Allow setting a function to control retry behavior [[GH-7331](https://github.com/hashicorp/vault/pull/7331)]
* auth/gcp: Fix a bug where region information in instance groups names could
   cause an authorization attempt to fail [[GH-74](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/74)]
* cli: Fix a bug where a token of an unknown format (e.g. in ~/.vault-token)
   could cause confusing error messages during `vault login` [[GH-7508](https://github.com/hashicorp/vault/pull/7508)]
* cli: Fix a bug where the `namespace list` command with JSON formatting
   always returned an empty object [[GH-7705](https://github.com/hashicorp/vault/pull/7705)]
* cli: Command timeouts are now always specified solely by the
   `VAULT_CLIENT_TIMEOUT` value. [[GH-7469](https://github.com/hashicorp/vault/pull/7469)]
* core: Don't allow registering a non-root zero TTL token lease. This is purely
   defense in depth as the lease would be revoked immediately anyways, but
   there's no real reason to allow registration. [[GH-7524](https://github.com/hashicorp/vault/pull/7524)]
* core: Correctly revoke the token that's present in the response auth from a
   auth/token/ request if there's partial failure during the process. [[GH-7835](https://github.com/hashicorp/vault/pull/7835)]
* identity (enterprise): Fixed identity case sensitive loading in secondary
   cluster [[GH-7327](https://github.com/hashicorp/vault/pull/7327)]
* identity: Ensure only replication primary stores the identity case sensitivity state [[GH-7820](https://github.com/hashicorp/vault/pull/7820)]
* raft: Fixed VAULT_CLUSTER_ADDR env being ignored at startup [[GH-7619](https://github.com/hashicorp/vault/pull/7619)]
* secrets/pki: Don't allow duplicate SAN names in issued certs [[GH-7605](https://github.com/hashicorp/vault/pull/7605)]
* sys/health: Pay attention to the values provided for `standbyok` and
   `perfstandbyok` rather than simply using their presence as a key to flip on
   that behavior [[GH-7323](https://github.com/hashicorp/vault/pull/7323)]
* ui: using the `wrapped_token` query param will work with `redirect_to` and
   will automatically log in as intended [[GH-7398](https://github.com/hashicorp/vault/pull/7398)]
* ui: fix an error when initializing from the UI using PGP keys [[GH-7542](https://github.com/hashicorp/vault/pull/7542)]
* ui: show all active kv v2 secret versions even when `delete_version_after` is configured [[GH-7685](https://github.com/hashicorp/vault/pull/7685)]
* ui: Ensure that items in the top navigation link to pages that users have access to [[GH-7590](https://github.com/hashicorp/vault/pull/7590)]

## 1.2.7

### August 27th, 2020

NOTE:

All security content from 1.5.2, 1.5.1, 1.4.5, 1.4.4, 1.3.9, 1.3.8, 1.2.6, and 1.2.5 has been made fully open source, and the git tags for 1.5.3, 1.4.6, 1.3.10, and 1.2.7 will build correctly for open source users.

BUG FIXES:

* auth/aws: Made header handling for IAM authentication more robust

## 1.2.6.1

### August 21st, 2020

### Enterprise Only

NOTE:

Includes correct license in the HSM binary.

## 1.2.6

### August 20th, 2020

NOTE:

OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.

KNOWN ISSUES:

* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.2.6 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.2.6)
* In versions 1.2.6, 1.3.9, 1.4.5, and 1.5.2, enterprise licenses on the HSM build were not incorporated correctly - enterprise
  customers should use 1.2.6.1, 1.3.9.1, 1.4.5.1, and 1.5.2.1.

## 1.2.5

### August 20th, 2020

SECURITY:

* When using the IAM AWS Auth Method, under certain circumstances, values Vault uses to validate identities and roles can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.7.1 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16250) (Discovered by Felix Wilhelm of Google Project Zero)
* When using the GCP GCE Auth Method, under certain circumstances, values Vault uses to validate GCE VMs can be manipulated and bypassed. This vulnerability affects Vault and Vault Enterprise 0.8.3 and newer and is fixed in 1.2.5, 1.3.8, 1.4.4, and 1.5.1 (CVE-2020-16251) (Discovered by Felix Wilhelm of Google Project Zero)

KNOWN ISSUES:

* OSS binaries of 1.5.1, 1.4.4, 1.3.8, and 1.2.5 were built without the Vault UI. Enterprise binaries are not affected.
* AWS IAM logins may return an error depending on the headers sent with the request.
  For more details and a workaround, see the [1.2.5 Upgrade Guide](https://www.vaultproject.io/docs/upgrading/upgrade-to-1.2.5)

BUG FIXES:

* seal: (enterprise) Fix issue causing stored seal and recovery keys to be mistaken as sealwrapped values

## 1.2.4 (November 7th, 2019)

SECURITY:

* In a non-root namespace, revocation of a token scoped to a non-root
   namespace did not trigger the expected revocation of dynamic secret leases
   associated with that token. As a result, dynamic secret leases in non-root
   namespaces may outlive the token that created them.  This vulnerability,
   CVE-2019-18616, affects Vault Enterprise 0.11.0 and newer.
* Disaster Recovery secondary clusters did not delete already-replicated data
   after a mount filter has been created on an upstream Performance secondary
   cluster. As a result, encrypted secrets may remain replicated on a Disaster
   Recovery secondary cluster after application of a mount filter excluding
   those secrets from replication. This vulnerability, CVE-2019-18617, affects
   Vault Enterprise 0.8 and newer.
* Update version of Go to 1.12.12 to fix Go bug golang.org/issue/34960 which
   corresponds to CVE-2019-17596.

CHANGES:

* auth/aws: If a custom `sts_endpoint` is configured, Vault Agent and the CLI
   should provide the corresponding region via the `region` parameter (which
   already existed as a CLI parameter, and has now been added to Agent). The
   automatic region detection added to the CLI and Agent in 1.2 has been removed.

IMPROVEMENTS:

* cli: Ignore existing token during CLI login [[GH-7508](https://github.com/hashicorp/vault/pull/7508)]
* core: Log proxy settings from environment on startup [[GH-7528](https://github.com/hashicorp/vault/pull/7528)]
* core: Cache whether we've been initialized to reduce load on storage [[GH-7549](https://github.com/hashicorp/vault/pull/7549)]

BUG FIXES:

* agent: Fix handling of gzipped responses [[GH-7470](https://github.com/hashicorp/vault/pull/7470)]
* cli: Fix panic when pgp keys list is empty [[GH-7546](https://github.com/hashicorp/vault/pull/7546)]
* cli: Command timeouts are now always specified solely by the
   `VAULT_CLIENT_TIMEOUT` value. [[GH-7469](https://github.com/hashicorp/vault/pull/7469)]
* core: add hook for initializing seals for migration [[GH-7666](https://github.com/hashicorp/vault/pull/7666)]
* core (enterprise): Migrating from one auto unseal method to another never
   worked on enterprise, now it does.
* identity: Add required field `response_types_supported` to identity token
   `.well-known/openid-configuration` response [[GH-7533](https://github.com/hashicorp/vault/pull/7533)]
* identity: Fixed nil pointer panic when merging entities [[GH-7712](https://github.com/hashicorp/vault/pull/7712)]
* replication (Enterprise): Fix issue causing performance standbys nodes
   disconnecting when under high loads.
* secrets/azure: Fix panic that could occur if client retries timeout [[GH-7793](https://github.com/hashicorp/vault/pull/7793)]
* secrets/database: Fix bug in combined DB secrets engine that can result in
   writes to static-roles endpoints timing out [[GH-7518](https://github.com/hashicorp/vault/pull/7518)]
* secrets/pki: Improve tidy to continue when value is nil [[GH-7589](https://github.com/hashicorp/vault/pull/7589)]
* ui (Enterprise): Allow kv v2 secrets that are gated by Control Groups to be
   viewed in the UI [[GH-7504](https://github.com/hashicorp/vault/pull/7504)]

## 1.2.3 (September 12, 2019)

FEATURES:

* **Oracle Cloud (OCI) Integration**: Vault now support using Oracle Cloud for
  storage, auto unseal, and authentication.

IMPROVEMENTS:

* auth/jwt: Groups claim matching now treats a string response as a single
   element list [[GH-63](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/63)]
* auth/kubernetes: enable better support for projected tokens API by allowing
   user to specify issuer [[GH-65](https://github.com/hashicorp/vault/pull/65)]
* auth/pcf: The PCF auth plugin was renamed to the CF auth plugin, maintaining
   full backwards compatibility [[GH-7346](https://github.com/hashicorp/vault/pull/7346)]
* replication: Premium packages now come with unlimited performance standby
   nodes

BUG FIXES:

* agent: Allow batch tokens and other non-renewable tokens to be used for
   agent operations [[GH-7441](https://github.com/hashicorp/vault/pull/7441)]
* auth/jwt: Fix an error where newer (v1.2) token_* configuration parameters
   were not being applied to tokens generated using the OIDC login flow
   [[GH-67](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/67)]
* raft: Fix an incorrect JSON tag on `leader_ca_cert` in the join request [[GH-7393](https://github.com/hashicorp/vault/pull/7393)]
* seal/transit: Allow using Vault Agent for transit seal operations [[GH-7441](https://github.com/hashicorp/vault/pull/7441)]
* storage/couchdb: Fix a file descriptor leak [[GH-7345](https://github.com/hashicorp/vault/pull/7345)]
* ui: Fix a bug where the status menu would disappear when trying to revoke a
   token [[GH-7337](https://github.com/hashicorp/vault/pull/7337)]
* ui: Fix a regression that prevented input of custom items in search-select
   [[GH-7338](https://github.com/hashicorp/vault/pull/7338)]
* ui: Fix an issue with the namespace picker being unable to render nested
   namespaces named with numbers and sorting of namespaces in the picker
   [[GH-7333](https://github.com/hashicorp/vault/pull/7333)]

## 1.2.2 (August 15, 2019)

CHANGES:

* auth/pcf: The signature format has been updated to use the standard Base64
   encoding instead of the URL-safe variant. Signatures created using the
   previous format will continue to be accepted [PCF-27]
* core: The http response code returned when an identity token key is not found
   has been changed from 400 to 404

IMPROVEMENTS:

* identity: Remove 512 entity limit for groups [[GH-7317](https://github.com/hashicorp/vault/pull/7317)]

BUG FIXES:

* auth/approle: Fix an error where an empty `token_type` string was not being
   correctly handled as `TokenTypeDefault` [[GH-7273](https://github.com/hashicorp/vault/pull/7273)]
* auth/radius: Fix panic when logging in [[GH-7286](https://github.com/hashicorp/vault/pull/7286)]
* ui: the string-list widget will now honor multiline input [[GH-7254](https://github.com/hashicorp/vault/pull/7254)]
* ui: various visual bugs in the KV interface were addressed [[GH-7307](https://github.com/hashicorp/vault/pull/7307)]
* ui: fixed incorrect URL to access help in LDAP auth [[GH-7299](https://github.com/hashicorp/vault/pull/7299)]

## 1.2.1 (August 6th, 2019)

BUG FIXES:

* agent: Fix a panic on creds pulling in some error conditions in `aws` and
   `alicloud` auth methods [[GH-7238](https://github.com/hashicorp/vault/pull/7238)]
* auth/approle: Fix error reading role-id on a role created pre-1.2 [[GH-7231](https://github.com/hashicorp/vault/pull/7231)]
* auth/token: Fix sudo check in non-root namespaces on create [[GH-7224](https://github.com/hashicorp/vault/pull/7224)]
* core: Fix health checks with perfstandbyok=true returning the wrong status
   code [[GH-7240](https://github.com/hashicorp/vault/pull/7240)]
* ui: The web CLI will now parse input as a shell string, with special
   characters escaped [[GH-7206](https://github.com/hashicorp/vault/pull/7206)]
* ui: The UI will now redirect to a page after authentication [[GH-7088](https://github.com/hashicorp/vault/pull/7088)]
* ui (Enterprise): The list of namespaces is now cleared when logging
   out [[GH-7186](https://github.com/hashicorp/vault/pull/7186)]

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
   not set. See [[GH-6717](https://github.com/hashicorp/vault/pull/6717)] for
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

* agent: Allow EC2 nonce to be passed in [[GH-6953](https://github.com/hashicorp/vault/pull/6953)]
* agent: Add optional `namespace` parameter, which sets the default namespace
   for the auto-auth functionality [[GH-6988](https://github.com/hashicorp/vault/pull/6988)]
* agent: Add cert auto-auth method [[GH-6652](https://github.com/hashicorp/vault/pull/6652)]
* api: Add support for passing data to delete operations via `DeleteWithData`
   [[GH-7139](https://github.com/hashicorp/vault/pull/7139)]
* audit/file: Dramatically speed up file operations by changing
   locking/marshaling order [[GH-7024](https://github.com/hashicorp/vault/pull/7024)]
* auth/jwt: A JWKS endpoint may now be configured for signature verification [[GH-43](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/43)]
* auth/jwt: A new `verbose_oidc_logging` role parameter has been added to help
   troubleshoot OIDC configuration [[GH-57](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/57)]
* auth/jwt: `bound_claims` will now match received claims that are lists if any element
   of the list is one of the expected values [[GH-50](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/50)]
* auth/jwt: Leeways for `nbf` and `exp` are now configurable, as is clock skew
   leeway [[GH-53](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/53)]
* auth/kubernetes: Allow service names/namespaces to be configured as globs
   [[GH-58](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/58)]
* auth/token: Allow the support of the identity system for the token backend
   via token roles [[GH-6267](https://github.com/hashicorp/vault/pull/6267)]
* auth/token: Add a large set of token configuration options to token store
   roles [[GH-6662](https://github.com/hashicorp/vault/pull/6662)]
* cli: `path-help` now allows `-format=json` to be specified, which will
   output OpenAPI [[GH-7006](https://github.com/hashicorp/vault/pull/7006)]
* cli: Add support for passing parameters to `vault delete` operations
   [[GH-7139](https://github.com/hashicorp/vault/pull/7139)]
* cli: Add a log-format CLI flag that can specify either "standard" or "json"
   for the log format for the `vault server`command. [[GH-6840](https://github.com/hashicorp/vault/pull/6840)]
* cli: Add `-dev-no-store-token` to allow dev servers to not store the
   generated token at the tokenhelper location [[GH-7104](https://github.com/hashicorp/vault/pull/7104)]
* identity: Allow a group alias' canonical ID to be modified
* namespaces: Namespaces can now be created and deleted from performance
   replication secondaries
* plugins: Change the default for `max_open_connections` for DB plugins to 4
   [[GH-7093](https://github.com/hashicorp/vault/pull/7093)]
* replication: Client TLS authentication is now supported when enabling or
   updating a replication secondary
* secrets/database: Cassandra operations will now cancel on client timeout
   [[GH-6954](https://github.com/hashicorp/vault/pull/6954)]
* secrets/kv: Add optional `delete_version_after` parameter, which takes a
   duration and can be set on the mount and/or the metadata for a specific key
   [[GH-7005](https://github.com/hashicorp/vault/pull/7005)]
* storage/postgres: LIST now performs better on large datasets [[GH-6546](https://github.com/hashicorp/vault/pull/6546)]
* storage/s3: A new `path` parameter allows selecting the path within a bucket
   for Vault data [[GH-7157](https://github.com/hashicorp/vault/pull/7157)]
* ui: KV v1 and v2 will now gracefully degrade allowing a write without read
   workflow in the UI [[GH-6570](https://github.com/hashicorp/vault/pull/6570)]
* ui: Many visual improvements with the addition of Toolbars [[GH-6626](https://github.com/hashicorp/vault/pull/6626)], the restyling
   of the Confirm Action component [[GH-6741](https://github.com/hashicorp/vault/pull/6741)], and using a new set of glyphs for our
   Icon component [[GH-6736](https://github.com/hashicorp/vault/pull/6736)]
* ui: Lazy loading parts of the application so that the total initial payload is
   smaller [[GH-6718](https://github.com/hashicorp/vault/pull/6718)]
* ui: Tabbing to auto-complete in filters will first complete a common prefix if there
   is one [[GH-6759](https://github.com/hashicorp/vault/pull/6759)]
* ui: Removing jQuery from the application makes the initial JS payload smaller [[GH-6768](https://github.com/hashicorp/vault/pull/6768)]

BUG FIXES:

* audit: Log requests and responses due to invalid wrapping token provided
   [[GH-6541](https://github.com/hashicorp/vault/pull/6541)]
* audit: Fix bug preventing request counter queries from working with auditing
   enabled [[GH-6767](https://github.com/hashicorp/vault/pull/6767)
* auth/aws: AWS Roles are now upgraded and saved to the latest version just
   after the AWS credential plugin is mounted. [[GH-7025](https://github.com/hashicorp/vault/pull/7025)]
* auth/aws: Fix a case where a panic could stem from a malformed assumed-role ARN
   when parsing this value [[GH-6917](https://github.com/hashicorp/vault/pull/6917)]
* auth/aws: Fix an error complaining about a read-only view that could occur
   during updating of a role when on a performance replication secondary
   [[GH-6926](https://github.com/hashicorp/vault/pull/6926)]
* auth/jwt: Fix a regression introduced in 1.1.1 that disabled checking of client_id
   for OIDC logins [[GH-54](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/54)]
* auth/jwt: Fix a panic during OIDC CLI logins that could occur if the Vault server
   response is empty [[GH-55](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/55)]
* auth/jwt: Fix issue where OIDC logins might intermittently fail when using
   performance standbys [[GH-61](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/61)]
* identity: Fix a case where modifying aliases of an entity could end up
   moving the entity into the wrong namespace
* namespaces: Fix a behavior (currently only known to be benign) where we
   wouldn't delete policies through the official functions before wiping the
   namespaces on deletion
* secrets/database: Escape username/password before using in connection URL
   [[GH-7089](https://github.com/hashicorp/vault/pull/7089)]
* secrets/pki: Forward revocation requests to active node when on a
   performance standby [[GH-7173](https://github.com/hashicorp/vault/pull/7173)]
* ui: Fix timestamp on some transit keys [[GH-6827](https://github.com/hashicorp/vault/pull/6827)]
* ui: Show Entities and Groups in Side Navigation [[GH-7138](https://github.com/hashicorp/vault/pull/7138)]
* ui: Ensure dropdown updates selected item on HTTP Request Metrics page

## 1.1.4/1.1.5 (July 25th/30th, 2019)

NOTE:

Although 1.1.4 was tagged, we realized very soon after the tag was publicly
pushed that an intended fix was accidentally left out. As a result, 1.1.4 was
not officially announced and 1.1.5 should be used as the release after 1.1.3.

IMPROVEMENTS:

* identity: Allow a group alias' canonical ID to be modified
* namespaces: Improve namespace deletion performance [[GH-6939](https://github.com/hashicorp/vault/pull/6939)]
* namespaces: Namespaces can now be created and deleted from performance
   replication secondaries

BUG FIXES:

* api: Add backwards compat support for API env vars [[GH-7135](https://github.com/hashicorp/vault/pull/7135)]
* auth/aws: Fix a case where a panic could stem from a malformed assumed-role
   ARN when parsing this value [[GH-6917](https://github.com/hashicorp/vault/pull/6917)]
* auth/ldap: Add `use_pre111_group_cn_behavior` flag to allow recovering from
   a regression caused by a bug fix starting in 1.1.1 [[GH-7208](https://github.com/hashicorp/vault/pull/7208)]
* auth/aws: Use a role cache to avoid separate locking paths [[GH-6926](https://github.com/hashicorp/vault/pull/6926)]
* core: Fix a deadlock if a panic happens during request handling [[GH-6920](https://github.com/hashicorp/vault/pull/6920)]
* core: Fix an issue that may cause key upgrades to not be cleaned up properly
   [[GH-6949](https://github.com/hashicorp/vault/pull/6949)]
* core: Don't shutdown if key upgrades fail due to canceled context [[GH-7070](https://github.com/hashicorp/vault/pull/7070)]
* core: Fix panic caused by handling requests while vault is inactive
* identity: Fix reading entity and groups that have spaces in their names
   [[GH-7055](https://github.com/hashicorp/vault/pull/7055)]
* identity: Ensure entity alias operations properly verify namespace [[GH-6886](https://github.com/hashicorp/vault/pull/6886)]
* mfa: Fix a nil pointer panic that could occur if invalid Duo credentials
   were supplied
* replication: Forward step-down on perf standbys to match HA behavior
* replication: Fix various read only storage errors on performance standbys
* replication: Stop forwarding before stopping replication to eliminate some
   possible bad states
* secrets/database: Allow cassandra queries to be cancled [[GH-6954](https://github.com/hashicorp/vault/pull/6954)]
* storage/consul: Fix a regression causing vault to not connect to consul over
   unix sockets [[GH-6859](https://github.com/hashicorp/vault/pull/6859)]
* ui: Fix saving of TTL and string array fields generated by Open API [[GH-7094](https://github.com/hashicorp/vault/pull/7094)]

## 1.1.3 (June 5th, 2019)

IMPROVEMENTS:

* agent: Now supports proxying request query parameters [[GH-6772](https://github.com/hashicorp/vault/pull/6772)]
* core: Mount table output now includes a UUID indicating the storage path [[GH-6633](https://github.com/hashicorp/vault/pull/6633)]
* core: HTTP server timeout values are now configurable [[GH-6666](https://github.com/hashicorp/vault/pull/6666)]
* replication: Improve performance of the reindex operation on secondary clusters
   when mount filters are in use
* replication: Replication status API now returns the state and progress of a reindex

BUG FIXES:

* api: Return the Entity ID in the secret output [[GH-6819](https://github.com/hashicorp/vault/pull/6819)]
* auth/jwt: Consider bound claims when considering if there is at least one
   bound constraint [[GH-49](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/49)]
* auth/okta: Fix handling of group names containing slashes [[GH-6665](https://github.com/hashicorp/vault/pull/6665)]
* cli: Add deprecated stored-shares flag back to the init command [[GH-6677](https://github.com/hashicorp/vault/pull/6677)]
* cli: Fix a panic when the KV command would return no data [[GH-6675](https://github.com/hashicorp/vault/pull/6675)]
* cli: Fix issue causing CLI list operations to not return proper format when
   there is an empty response [[GH-6776](https://github.com/hashicorp/vault/pull/6776)]
* core: Correctly honor non-HMAC request keys when auditing requests [[GH-6653](https://github.com/hashicorp/vault/pull/6653)]
* core: Fix the `x-vault-unauthenticated` value in OpenAPI for a number of
   endpoints [[GH-6654](https://github.com/hashicorp/vault/pull/6654)]
* core: Fix issue where some OpenAPI parameters were incorrectly listed as
   being sent as a header [[GH-6679](https://github.com/hashicorp/vault/pull/6679)]
* core: Fix issue that would allow duplicate mount names to be used [[GH-6771](https://github.com/hashicorp/vault/pull/6771)]
* namespaces: Fix behavior when using `root` instead of `root/` as the
   namespace header value
* pki: fix a panic when a client submits a null value [[GH-5679](https://github.com/hashicorp/vault/pull/5679)]
* replication: Properly update mount entry cache on a secondary to apply all
   new values after a tune
* replication: Properly close connection on bootstrap error
* replication: Fix an issue causing startup problems if a namespace policy
   wasn't replicated properly
* replication: Fix longer than necessary WAL replay during an initial reindex
* replication: Fix error during mount filter invalidation on DR secondary clusters
* secrets/ad: Make time buffer configurable [AD-35]
* secrets/gcp: Check for nil config when getting credentials [[GH-35](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/35)]
* secrets/gcp: Fix error checking in some cases where the returned value could
   be 403 instead of 404 [[GH-37](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/37)]
* secrets/gcpkms: Disable key rotation when deleting a key [[GH-10](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/10)]
* storage/consul: recognize `https://` address even if schema not specified
   [[GH-6602](https://github.com/hashicorp/vault/pull/6602)]
* storage/dynamodb: Fix an issue where a deleted lock key in DynamoDB (HA)
   could cause constant switching of the active node [[GH-6637](https://github.com/hashicorp/vault/pull/6637)]
* storage/dynamodb: Eliminate a high-CPU condition that could occur if an
   error was received from the DynamoDB API [[GH-6640](https://github.com/hashicorp/vault/pull/6640)]
* storage/gcs: Correctly use configured chunk size values [[GH-6655](https://github.com/hashicorp/vault/pull/6655)]
* storage/mssql: Use the correct database when pre-created schemas exist
   [[GH-6356](https://github.com/hashicorp/vault/pull/6356)]
* ui: Fix issue with select arrows on drop down menus [[GH-6627](https://github.com/hashicorp/vault/pull/6627)]
* ui: Fix an issue where sensitive input values weren't being saved to the
   server [[GH-6586](https://github.com/hashicorp/vault/pull/6586)]
* ui: Fix web cli parsing when using quoted values [[GH-6755](https://github.com/hashicorp/vault/pull/6755)]
* ui: Fix a namespace workflow mapping identities from external namespaces by
   allowing arbitrary input in search-select component [[GH-6728](https://github.com/hashicorp/vault/pull/6728)]

## 1.1.2 (April 18th, 2019)

This is a bug fix release containing the two items below. It is otherwise
unchanged from 1.1.1.

BUG FIXES:

* auth/okta: Fix a potential dropped error [[GH-6592](https://github.com/hashicorp/vault/pull/6592)]
* secrets/kv: Fix a regression on upgrade where a KVv2 mount could fail to be
   mounted on unseal if it had previously been mounted but not written to
   [[GH-31](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/31)]

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

* auth/jwt: Disallow logins of role_type "oidc" via the `/login` path [[GH-38](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/38)]
* core/acl:  New ordering defines which policy wins when there are multiple
   inexact matches and at least one path contains `+`. `+*` is now illegal in
   policy paths. The previous behavior simply selected any matching
   segment-wildcard path that matched. [[GH-6532](https://github.com/hashicorp/vault/pull/6532)]
* replication: Due to technical limitations, mounting and unmounting was not
   previously possible from a performance secondary. These have been resolved,
   and these operations may now be run from a performance secondary.

IMPROVEMENTS:

* agent: Allow AppRole auto-auth without a secret-id [[GH-6324](https://github.com/hashicorp/vault/pull/6324)]
* auth/gcp: Cache clients to improve performance and reduce open file usage
* auth/jwt: Bounds claims validiation will now allow matching the received
   claims against a list of expected values [[GH-41](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/41)]
* secret/gcp: Cache clients to improve performance and reduce open file usage
* replication: Mounting/unmounting/remounting/mount-tuning is now supported
   from a performance secondary cluster
* ui: Suport for authentication via the RADIUS auth method [[GH-6488](https://github.com/hashicorp/vault/pull/6488)]
* ui: Navigating away from secret list view will clear any page-specific
   filter that was applied [[GH-6511](https://github.com/hashicorp/vault/pull/6511)]
* ui: Improved the display when OIDC auth errors [[GH-6553](https://github.com/hashicorp/vault/pull/6553)]

BUG FIXES:

* agent: Allow auto-auth to be used with caching without having to define any
   sinks [[GH-6468](https://github.com/hashicorp/vault/pull/6468)]
* agent: Disallow some nonsensical config file combinations [[GH-6471](https://github.com/hashicorp/vault/pull/6471)]
* auth/ldap: Fix CN check not working if CN was not all in uppercase [[GH-6518](https://github.com/hashicorp/vault/pull/6518)]
* auth/jwt: The CLI helper for OIDC logins will now open the browser to the correct
   URL when running on Windows [[GH-37](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/37)]
* auth/jwt: Fix OIDC login issue where configured TLS certs weren't being used [[GH-40](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/40)]
* auth/jwt: Fix an issue where the `oidc_scopes` parameter was not being included in
   the response to a role read request [[GH-35](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/35)]
* core: Fix seal migration case when migrating to Shamir and a seal block
   wasn't explicitly specified [[GH-6455](https://github.com/hashicorp/vault/pull/6455)]
* core: Fix unwrapping when using namespaced wrapping tokens [[GH-6536](https://github.com/hashicorp/vault/pull/6536)]
* core: Fix incorrect representation of required properties in OpenAPI output
   [[GH-6490](https://github.com/hashicorp/vault/pull/6490)]
* core: Fix deadlock that could happen when using the UI [[GH-6560](https://github.com/hashicorp/vault/pull/6560)]
* identity: Fix updating groups removing existing members [[GH-6527](https://github.com/hashicorp/vault/pull/6527)]
* identity: Properly invalidate group alias in performance secondary [[GH-6564](https://github.com/hashicorp/vault/pull/6564)]
* identity: Use namespace context when loading entities and groups to ensure
   merging of duplicate entries works properly [[GH-6563](https://github.com/hashicorp/vault/pull/6563)]
* replication: Fix performance standby election failure [[GH-6561](https://github.com/hashicorp/vault/pull/6561)]
* replication: Fix mount filter invalidation on performance standby nodes
* replication: Fix license reloading on performance standby nodes
* replication: Fix handling of control groups on performance standby nodes
* replication: Fix some forwarding scenarios with request bodies using
   performance standby nodes [[GH-6538](https://github.com/hashicorp/vault/pull/6538)]
* secret/gcp: Fix roleset binding when using JSON [[GH-27](https://github.com/hashicorp/vault-plugin-secrets-gcp/pull/27)]
* secret/pki: Use `uri_sans` param in when not using CSR parameters [[GH-6505](https://github.com/hashicorp/vault/pull/6505)]
* storage/dynamodb: Fix a race condition possible in HA configurations that could
   leave the cluster without a leader [[GH-6512](https://github.com/hashicorp/vault/pull/6512)]
* ui: Fix an issue where in production builds OpenAPI model generation was
   failing, causing any form using it to render labels with missing fields [[GH-6474](https://github.com/hashicorp/vault/pull/6474)]
* ui: Fix issue nav-hiding when moving between namespaces [[GH-6473](https://github.com/hashicorp/vault/pull/6473)]
* ui: Secrets will always show in the nav regardless of access to cubbyhole [[GH-6477](https://github.com/hashicorp/vault/pull/6477)]
* ui: fix SSH OTP generation [[GH-6540](https://github.com/hashicorp/vault/pull/6540)]
* ui: add polyfill to load UI in IE11 [[GH-6567](https://github.com/hashicorp/vault/pull/6567)]
* ui: Fix issue where some elements would fail to work properly if using ACLs
   with segment-wildcard paths (`/+/` segments) [[GH-6525](https://github.com/hashicorp/vault/pull/6525)]

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
   specification [[GH-6300](https://github.com/hashicorp/vault/pull/6300)]
* core/metrics: Prometheus pull support using a new sys/metrics endpoint. [[GH-5308](https://github.com/hashicorp/vault/pull/5308)]
* core: On non-windows platforms a SIGUSR2 will make the server log a dump of
   all running goroutines' stack traces for debugging purposes [[GH-6240](https://github.com/hashicorp/vault/pull/6240)]
* replication: The initial replication indexing process on newly initialized or upgraded
   clusters now runs asynchronously
* sentinel: Add token namespace id and path, available in rules as
   token.namespace.id and token.namespace.path
* ui: The UI is now leveraging OpenAPI definitions to pull in fields for various forms.
   This means, it will not be necessary to add fields on the go and JS sides in the future.
   [[GH-6209](https://github.com/hashicorp/vault/pull/6209)]

BUG FIXES:

* auth/jwt: Apply `bound_claims` validation across all login paths
* auth/jwt: Update `bound_audiences` validation during non-OIDC logins to accept
   any matched audience, as documented and handled in OIDC logins [[GH-30](https://github.com/hashicorp/vault-plugin-auth-jwt/pull/30)]
* auth/token: Fix issue where empty values for token role update call were
   ignored [[GH-6314](https://github.com/hashicorp/vault/pull/6314)]
* core: The `operator migrate` command will no longer hang on empty key names
   [[GH-6371](https://github.com/hashicorp/vault/pull/6371)]
* identity: Fix a panic at login when external group has a nil alias [[GH-6230](https://github.com/hashicorp/vault/pull/6230)]
* namespaces: Clear out identity store items upon namespace deletion
* replication/perfstandby: Fixed a bug causing performance standbys to wait
   longer than necessary after forwarding a write to the active node
* replication/mountfilter: Fix a deadlock that could occur when mount filters
   were updated [[GH-6426](https://github.com/hashicorp/vault/pull/6426)]
* secret/kv: Fix issue where a v1→v2 upgrade could run on a performance
   standby when using a local mount
* secret/ssh: Fix for a bug where attempting to delete the last ssh role
   in the zeroaddress configuration could fail [[GH-6390](https://github.com/hashicorp/vault/pull/6390)]
* secret/totp: Uppercase provided keys so they don't fail base32 validation
   [[GH-6400](https://github.com/hashicorp/vault/pull/6400)]
* secret/transit: Multiple HMAC, Sign or Verify operations can now be
   performed with one API call using the new `batch_input` parameter [[GH-5875](https://github.com/hashicorp/vault/pull/5875)]
* sys: `sys/internal/ui/mounts` will no longer return secret or auth mounts
   that have been filtered. Similarly, `sys/internal/ui/mount/:path` will
   return a error response if a filtered mount path is requested. [[GH-6412](https://github.com/hashicorp/vault/pull/6412)]
* ui: Fix for a bug where you couldn't access the data tab after clicking on
   wrap details on the unwrap page [[GH-6404](https://github.com/hashicorp/vault/pull/6404)]
* ui: Fix an issue where the policies tab was erroneously hidden [[GH-6301](https://github.com/hashicorp/vault/pull/6301)]
* ui: Fix encoding issues with kv interfaces [[GH-6294](https://github.com/hashicorp/vault/pull/6294)]

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
   entity either by name or by id [[GH-6105](https://github.com/hashicorp/vault/pull/6105)]
* The Vault UI's navigation and onboarding wizard now only displays items that
   are permitted in a users' policy [[GH-5980](https://github.com/hashicorp/vault/pull/5980), [GH-6094](https://github.com/hashicorp/vault/pull/6094)]
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
   role ID [[GH-6133](https://github.com/hashicorp/vault/pull/6133)]
* auth/jwt: The supported set of signing algorithms is now configurable [JWT
   plugin [GH-16](https://github.com/hashicorp/vault/pull/16)]
* core: When starting from an uninitialized state, HA nodes will now attempt
   to auto-unseal using a configured auto-unseal mechanism after the active
   node initializes Vault [[GH-6039](https://github.com/hashicorp/vault/pull/6039)]
* secret/database: Add socket keepalive option for Cassandra [[GH-6201](https://github.com/hashicorp/vault/pull/6201)]
* secret/ssh: Add signed key constraints, allowing enforcement of key types
   and minimum key sizes [[GH-6030](https://github.com/hashicorp/vault/pull/6030)]
* secret/transit: ECDSA signatures can now be marshaled in JWS-compatible
   fashion [[GH-6077](https://github.com/hashicorp/vault/pull/6077)]
* storage/etcd: Support SRV service names [[GH-6087](https://github.com/hashicorp/vault/pull/6087)]
* storage/aws: Support specifying a KMS key ID for server-side encryption
   [[GH-5996](https://github.com/hashicorp/vault/pull/5996)]

BUG FIXES:

* core: Fix a rare case where a standby whose connection is entirely torn down
   to the active node, then reconnects to the same active node, may not
   successfully resume operation [[GH-6167](https://github.com/hashicorp/vault/pull/6167)]
* cors: Don't duplicate headers when they're written [[GH-6207](https://github.com/hashicorp/vault/pull/6207)]
* identity: Persist merged entities only on the primary [[GH-6075](https://github.com/hashicorp/vault/pull/6075)]
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
   instead of just the last path element [[GH-6044](https://github.com/hashicorp/vault/pull/6044)]

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
   image ID [[GH-5846](https://github.com/hashicorp/vault/pull/5846)]
* autoseal/gcpckms: Reduce the required permissions for the GCPCKMS autounseal
   [[GH-5999](https://github.com/hashicorp/vault/pull/5999)]
* physical/foundationdb: TLS support added. [[GH-5800](https://github.com/hashicorp/vault/pull/5800)]

BUG FIXES:

* api: Fix a couple of places where we were using the `LIST` HTTP verb
   (necessary to get the right method into the wrapping lookup function) and
   not then modifying it to a `GET`; although this is officially the verb Vault
   uses for listing and it's fully legal to use custom verbs, since many WAFs
   and API gateways choke on anything outside of RFC-standardized verbs we fall
   back to `GET` [[GH-6026](https://github.com/hashicorp/vault/pull/6026)]
* autoseal/aws: Fix reading session tokens when AWS access key/secret key are
   also provided [[GH-5965](https://github.com/hashicorp/vault/pull/5965)]
* command/operator/rekey: Fix help output showing `-delete-backup` when it
   should show `-backup-delete` [[GH-5981](https://github.com/hashicorp/vault/pull/5981)]
* core: Fix bound_cidrs not being propagated to child tokens
* replication: Correctly forward identity entity creation that originates from
   performance standby nodes (Enterprise)
* secret/aws: Make input `credential_type` match the output type (string, not
   array) [[GH-5972](https://github.com/hashicorp/vault/pull/5972)]
* secret/cubbyhole: Properly cleanup cubbyhole after token revocation [[GH-6006](https://github.com/hashicorp/vault/pull/6006)]
* secret/pki: Fix reading certificates on windows with the file storage backend [[GH-6013](https://github.com/hashicorp/vault/pull/6013)]
* ui (enterprise): properly display perf-standby count on the license page [[GH-5971](https://github.com/hashicorp/vault/pull/5971)]
* ui: fix disappearing nested secrets and go to the nearest parent when deleting
   a secret - [[GH-5976](https://github.com/hashicorp/vault/pull/5976)]
* ui: fix error where deleting an item via the context menu would fail if the
   item name contained dots [[GH-6018](https://github.com/hashicorp/vault/pull/6018)]
* ui: allow saving of kv secret after an errored save attempt [[GH-6022](https://github.com/hashicorp/vault/pull/6022)]
* ui: fix display of kv-v1 secret containing a key named "keys" [[GH-6023](https://github.com/hashicorp/vault/pull/6023)]

## 1.0.1 (December 14th, 2018)

SECURITY:

* Update version of Go to 1.11.3 to fix Go bug
   <https://github.com/golang/go/issues/29233> which corresponds to
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

* cli: Strip iTerm extra characters from password manager input [[GH-5837](https://github.com/hashicorp/vault/pull/5837)]
* command/server: Setting default kv engine to v1 in -dev mode can now be
   specified via -dev-kv-v1 [[GH-5919](https://github.com/hashicorp/vault/pull/5919)]
* core: Add operationId field to OpenAPI output [[GH-5876](https://github.com/hashicorp/vault/pull/5876)]
* ui: Added ability to search for Group and Policy IDs when creating Groups
   and Entities instead of typing them in manually

BUG FIXES:

* auth/azure: Cache azure authorizer [15]
* auth/gcp: Remove explicit project for service account in GCE authorizer [[GH-58](https://github.com/hashicorp/vault-plugin-auth-gcp/pull/58)]
* cli: Show correct stored keys/threshold for autoseals [[GH-5910](https://github.com/hashicorp/vault/pull/5910)]
* cli: Fix backwards compatibility fallback when listing plugins [[GH-5913](https://github.com/hashicorp/vault/pull/5913)]
* core: Fix upgrades when the seal config had been created on early versions
   of vault [[GH-5956](https://github.com/hashicorp/vault/pull/5956)]
* namespaces: Correctly reload the proper mount when tuning or reloading the
   mount [[GH-5937](https://github.com/hashicorp/vault/pull/5937)]
* secret/azure: Cache azure authorizer [19]
* secret/database: Strip empty statements on user input [[GH-5955](https://github.com/hashicorp/vault/pull/5955)]
* secret/gcpkms: Add path for retrieving the public key [[GH-5](https://github.com/hashicorp/vault-plugin-secrets-gcpkms/pull/5)]
* secret/pki: Fix panic that could occur during tidy operation when malformed
   data was found [[GH-5931](https://github.com/hashicorp/vault/pull/5931)]
* secret/pki: Strip empty line in ca_chain output [[GH-5779](https://github.com/hashicorp/vault/pull/5779)]
* ui: Fixed a bug where the web CLI was not usable via the `fullscreen`
   command - [[GH-5909](https://github.com/hashicorp/vault/pull/5909)]
* ui: Fix a bug where you couldn't write a jwt auth method config [[GH-5936](https://github.com/hashicorp/vault/pull/5936)]

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
   character encoding [[GH-5819](https://github.com/hashicorp/vault/pull/5819)]
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
   [[GH-5725](https://github.com/hashicorp/vault/pull/5725)]
* auth/token: New tokens are indexed in storage HMAC-SHA256 instead of SHA1
* secret/totp: Allow @ character to be part of key name [[GH-5652](https://github.com/hashicorp/vault/pull/5652)]
* secret/consul: Add support for new policy based tokens added in Consul 1.4
   [[GH-5586](https://github.com/hashicorp/vault/pull/5586)]
* ui: Improve the token auto-renew warning, and automatically begin renewal
   when a user becomes active again [[GH-5662](https://github.com/hashicorp/vault/pull/5662)]
* ui: The unbundled UI page now has some styling [[GH-5665](https://github.com/hashicorp/vault/pull/5665)]
* ui: Improved banner and popup design [[GH-5672](https://github.com/hashicorp/vault/pull/5672)]
* ui: Added token type to auth method mount config [[GH-5723](https://github.com/hashicorp/vault/pull/5723)]
* ui: Display additonal wrap info when unwrapping. [[GH-5664](https://github.com/hashicorp/vault/pull/5664)]
* ui: Empty states have updated styling and link to relevant actions and
   documentation [[GH-5758](https://github.com/hashicorp/vault/pull/5758)]
* ui: Allow editing of KV V2 data when a token doesn't have capabilities to
   read secret metadata [[GH-5879](https://github.com/hashicorp/vault/pull/5879)]

BUG FIXES:

* agent: Fix auth when multiple redirects [[GH-5814](https://github.com/hashicorp/vault/pull/5814)]
* cli: Restore the `-policy-override` flag [[GH-5826](https://github.com/hashicorp/vault/pull/5826)]
* core: Fix rekey progress reset which did not happen under certain
   circumstances. [[GH-5743](https://github.com/hashicorp/vault/pull/5743)]
* core: Migration from autounseal to shamir will clean up old keys [[GH-5671](https://github.com/hashicorp/vault/pull/5671)]
* identity: Update group memberships when entity is deleted [[GH-5786](https://github.com/hashicorp/vault/pull/5786)]
* replication/perfstandby: Fix audit table upgrade on standbys [[GH-5811](https://github.com/hashicorp/vault/pull/5811)]
* replication/perfstandby: Fix redirect on approle update [[GH-5820](https://github.com/hashicorp/vault/pull/5820)]
* secrets/azure: Fix valid roles being rejected for duplicate ids despite
   having distinct scopes
   [[GH-16](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/16)]
* storage/gcs: Send md5 of values to GCS to avoid potential corruption
   [[GH-5804](https://github.com/hashicorp/vault/pull/5804)]
* secrets/kv: Fix issue where storage version would get incorrectly downgraded
   [[GH-5809](https://github.com/hashicorp/vault/pull/5809)]
* secrets/kv: Disallow empty paths on a `kv put` while accepting empty paths
   for all other operations for backwards compatibility
   [[GH-19](https://github.com/hashicorp/vault-plugin-secrets-kv/pull/19)]
* ui: Allow for secret creation in kv v2 when cas_required=true [[GH-5823](https://github.com/hashicorp/vault/pull/5823)]
* ui: Fix dr secondary operation token generation via the ui [[GH-5818](https://github.com/hashicorp/vault/pull/5818)]
* ui: Fix the PKI context menu so that items load [[GH-5824](https://github.com/hashicorp/vault/pull/5824)]
* ui: Update DR Secondary Token generation command [[GH-5857](https://github.com/hashicorp/vault/pull/5857)]
* ui: Fix pagination bug where controls would be rendered once for each
   item when viewing policies [[GH-5866](https://github.com/hashicorp/vault/pull/5866)]
* ui: Fix bug where `sys/leases/revoke` required 'sudo' capability to show
   the revoke button in the UI [[GH-5647](https://github.com/hashicorp/vault/pull/5647)]
* ui: Fix issue where certain pages wouldn't render in a namespace [[GH-5692](https://github.com/hashicorp/vault/pull/5692)]
