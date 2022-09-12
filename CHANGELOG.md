## 1.12.0
### Unreleased

CHANGES:

* core/entities: Fixed stranding of aliases upon entity merge, and require explicit selection of which aliases should be kept when some must be deleted [[GH-16539](https://github.com/hashicorp/vault/pull/16539)]
* core: Bump Go version to 1.18.5.
* core: Validate input parameters for vault operator init command. Vault 1.12 CLI version is needed to run operator init now. [[GH-16379](https://github.com/hashicorp/vault/pull/16379)]
* identity: a request to `/identity/group` that includes `member_group_ids` that contains a cycle will now be responded to with a 400 rather than 500 [[GH-15912](https://github.com/hashicorp/vault/pull/15912)]
* licensing (enterprise): Terminated licenses will no longer result in shutdown. Instead, upgrades
will not be allowed if the license termination time is before the build date of the binary.
* plugins: `GET /sys/plugins/catalog/:type/:name` endpoint now returns an additional `version` field in the response data. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* plugins: `GET /sys/plugins/catalog` endpoint now returns an additional `detailed` field in the response data with a list of additional plugin metadata. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]

FEATURES:

* **Secrets/auth plugin multiplexing**: manage multiple plugin configurations with a single plugin process [[GH-14946](https://github.com/hashicorp/vault/pull/14946)]
* secrets/database/hana: Add ability to customize dynamic usernames [[GH-16631](https://github.com/hashicorp/vault/pull/16631)]
* secrets/pki: Add an OCSP responder that implements a subset of RFC6960, answering single serial number OCSP requests for
a specific cluster's revoked certificates in a mount. [[GH-16723](https://github.com/hashicorp/vault/pull/16723)]
* ui: UI support for Okta Number Challenge. [[GH-15998](https://github.com/hashicorp/vault/pull/15998)]

IMPROVEMENTS:

* activity (enterprise): Added new clients unit tests to test accuracy of estimates
* agent: Added `disable_idle_connections` configuration to disable leaving idle connections open in auto-auth, caching and templating. [[GH-15986](https://github.com/hashicorp/vault/pull/15986)]
* agent: Added `disable_keep_alives` configuration to disable keep alives in auto-auth, caching and templating. [[GH-16479](https://github.com/hashicorp/vault/pull/16479)]
* agent: JWT auto auth now supports a `remove_jwt_after_reading` config option which defaults to true. [[GH-11969](https://github.com/hashicorp/vault/pull/11969)]
* agent: Send notifications to systemd on start and stop. [[GH-9802](https://github.com/hashicorp/vault/pull/9802)]
* api/mfa: Add namespace path to the MFA read/list endpoint [[GH-16911](https://github.com/hashicorp/vault/pull/16911)]
* api: Add a sentinel error for missing KV secrets [[GH-16699](https://github.com/hashicorp/vault/pull/16699)]
* auth/aws: PKCS7 signatures will now use SHA256 by default in prep for Go 1.18 [[GH-16455](https://github.com/hashicorp/vault/pull/16455)]
* auth/cert: Add metadata to identity-alias [[GH-14751](https://github.com/hashicorp/vault/pull/14751)]
* auth/gcp: Add support for GCE regional instance groups [[GH-16435](https://github.com/hashicorp/vault/pull/16435)]
* auth/jwt: Adds support for Microsoft US Gov L4 to the Azure provider for groups fetching. [[GH-16525](https://github.com/hashicorp/vault/pull/16525)]
* auth/jwt: Improves detection of Windows Subsystem for Linux (WSL) for CLI-based logins. [[GH-16525](https://github.com/hashicorp/vault/pull/16525)]
* auth/kerberos: add `add_group_aliases` config to include LDAP groups in Vault group aliases [[GH-16890](https://github.com/hashicorp/vault/pull/16890)]
* auth/kerberos: add `remove_instance_name` parameter to the login CLI and the 
Kerberos config in Vault. This removes any instance names found in the keytab 
service principal name. [[GH-16594](https://github.com/hashicorp/vault/pull/16594)]
* auth/oidc: Adds support for group membership parsing when using SecureAuth as an OIDC provider. [[GH-16274](https://github.com/hashicorp/vault/pull/16274)]
* cli: CLI commands will print a warning if flags will be ignored because they are passed after positional arguments. [[GH-16441](https://github.com/hashicorp/vault/pull/16441)]
* command/audit: Improve missing type error message [[GH-16409](https://github.com/hashicorp/vault/pull/16409)]
* command/server: add `-dev-tls` and `-dev-tls-cert-dir` subcommands to create a Vault dev server with generated certificates and private key. [[GH-16421](https://github.com/hashicorp/vault/pull/16421)]
* core (enterprise): Add HTTP PATCH support for namespaces with an associated `namespace patch` CLI command
* core (enterprise): Add check to `vault server` command to ensure configured storage backend is supported.
* core (enterprise): Add custom metadata support for namespaces
* core/activity: generate hyperloglogs containing clientIds for each month during precomputation [[GH-16146](https://github.com/hashicorp/vault/pull/16146)]
* core/activity: refactor activity log api to reuse partial api functions in activity endpoint when current month is specified [[GH-16162](https://github.com/hashicorp/vault/pull/16162)]
* core/activity: use monthly hyperloglogs to calculate new clients approximation for current month [[GH-16184](https://github.com/hashicorp/vault/pull/16184)]
* core/quotas (enterprise): Added ability to add path suffixes for lease-count resource quotas
* core/quotas (enterprise): Added ability to add role information for lease-count resource quotas, to limit login requests on auth mounts made using that role
* core/quotas: Added ability to add path suffixes for rate-limit resource quotas [[GH-15989](https://github.com/hashicorp/vault/pull/15989)]
* core/quotas: Added ability to add role information for rate-limit resource quotas, to limit login requests on auth mounts made using that role [[GH-16115](https://github.com/hashicorp/vault/pull/16115)]
* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* core: Limit activity log client count usage by namespaces [[GH-16000](https://github.com/hashicorp/vault/pull/16000)]
* core: Upgrade github.com/hashicorp/raft [[GH-16609](https://github.com/hashicorp/vault/pull/16609)]
* core: remove gox [[GH-16353](https://github.com/hashicorp/vault/pull/16353)]
* docs: Clarify the behaviour of local mounts in the context of DR replication [[GH-16218](https://github.com/hashicorp/vault/pull/16218)]
* identity/oidc: Adds support for detailed listing of clients and providers. [[GH-16567](https://github.com/hashicorp/vault/pull/16567)]
* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]
* identity/oidc: allows filtering the list providers response by an allowed_client_id [[GH-16181](https://github.com/hashicorp/vault/pull/16181)]
* identity: Prevent possibility of data races on entity creation. [[GH-16487](https://github.com/hashicorp/vault/pull/16487)]
* physical/postgresql: pass context to queries to propagate timeouts and cancellations on requests. [[GH-15866](https://github.com/hashicorp/vault/pull/15866)]
* plugins: Add Deprecation Status method to builtinregistry. [[GH-16846](https://github.com/hashicorp/vault/pull/16846)]
* plugins: Plugin catalog supports registering and managing plugins with semantic version information. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* secret/nomad: allow reading CA and client auth certificate from /nomad/config/access [[GH-15809](https://github.com/hashicorp/vault/pull/15809)]
* secret/pki: Add RSA PSS signature support for issuing certificates, signing CRLs [[GH-16519](https://github.com/hashicorp/vault/pull/16519)]
* secret/pki: Add signature_bits to sign-intermediate, sign-verbatim endpoints [[GH-16124](https://github.com/hashicorp/vault/pull/16124)]
* secret/pki: Allow issuing certificates with non-domain, non-email Common Names from roles, sign-verbatim, and as issuers (`cn_validations`). [[GH-15996](https://github.com/hashicorp/vault/pull/15996)]
* secret/pki: Allow specifying SKID for cross-signed issuance from older Vault versions. [[GH-16494](https://github.com/hashicorp/vault/pull/16494)]
* secret/transit: Allow importing Ed25519 keys from PKCS#8 with inner RFC 5915 ECPrivateKey blobs (NSS-wrapped keys). [[GH-15742](https://github.com/hashicorp/vault/pull/15742)]
* secrets/ad: set config default length only if password_policy is missing [[GH-16140](https://github.com/hashicorp/vault/pull/16140)]
* secrets/kubernetes: Add allowed_kubernetes_namespace_selector to allow selecting Kubernetes namespaces with a label selector when configuring roles. [[GH-16240](https://github.com/hashicorp/vault/pull/16240)]
* secrets/pki/tidy: Add another pair of metrics counting certificates not deleted by the tidy operation. [[GH-16702](https://github.com/hashicorp/vault/pull/16702)]
* secrets/pki: Add ability to periodically rebuild CRL before expiry [[GH-16762](https://github.com/hashicorp/vault/pull/16762)]
* secrets/pki: Add ability to periodically run tidy operations to remove expired certificates. [[GH-16900](https://github.com/hashicorp/vault/pull/16900)]
* secrets/pki: Add support for per-issuer Authority Information Access (AIA) URLs [[GH-16563](https://github.com/hashicorp/vault/pull/16563)]
* secrets/pki: Allow revocation of certificates with explicitly provided certificate (bring your own certificate / BYOC). [[GH-16564](https://github.com/hashicorp/vault/pull/16564)]
* secrets/pki: Allow revocation via proving possession of certificate's private key [[GH-16566](https://github.com/hashicorp/vault/pull/16566)]
* secrets/pki: Allow tidy to associate revoked certs with their issuers for OCSP performance [[GH-16871](https://github.com/hashicorp/vault/pull/16871)]
* secrets/pki: Honor If-Modified-Since header on CA, CRL fetch; requires passthrough_request_headers modification on the mount point. [[GH-16249](https://github.com/hashicorp/vault/pull/16249)]
* secrets/pki: Improve stability of association of revoked cert with its parent issuer; when an issuer loses crl-signing usage, do not place certs on default issuer's CRL. [[GH-16874](https://github.com/hashicorp/vault/pull/16874)]
* secrets/pki: Support generating delta CRLs for up-to-date CRLs when auto-building is enabled. [[GH-16773](https://github.com/hashicorp/vault/pull/16773)]
* secrets/ssh: Add allowed_domains_template to allow templating of allowed_domains. [[GH-16056](https://github.com/hashicorp/vault/pull/16056)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]
* secrets/ssh: Allow the use of Identity templates in the `default_user` field [[GH-16351](https://github.com/hashicorp/vault/pull/16351)]
* ssh: Addition of an endpoint `ssh/issue/:role` to allow the creation of signed key pairs [[GH-15561](https://github.com/hashicorp/vault/pull/15561)]
* storage/cassandra: tuning parameters for clustered environments `connection_timeout`, `initial_connection_timeout`, `simple_retry_policy_retries`. [[GH-10467](https://github.com/hashicorp/vault/pull/10467)]
* storage/gcs: Add documentation explaining how to configure the gcs backend using environment variables instead of options in the configuration stanza [[GH-14455](https://github.com/hashicorp/vault/pull/14455)]
* ui: Changed the tokenBoundCidrs tooltip content to clarify that comma separated values are not accepted in this field. [[GH-15852](https://github.com/hashicorp/vault/pull/15852)]
* ui: Removed deprecated version of core-js 2.6.11 [[GH-15898](https://github.com/hashicorp/vault/pull/15898)]
* ui: Renamed labels under Tools for wrap, lookup, rewrap and unwrap with description. [[GH-16489](https://github.com/hashicorp/vault/pull/16489)]
* ui: redirect_to param forwards from auth route when authenticated [[GH-16821](https://github.com/hashicorp/vault/pull/16821)]
* website/docs: API generate-recovery-token documentation. [[GH-16213](https://github.com/hashicorp/vault/pull/16213)]
* website/docs: Update replication docs to mention Integrated Storage [[GH-16063](https://github.com/hashicorp/vault/pull/16063)]
* website/docs: changed to echo for all string examples instead of (<<<) here-string. [[GH-9081](https://github.com/hashicorp/vault/pull/9081)]

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* agent: Update consul-template for pkiCert bug fixes [[GH-16087](https://github.com/hashicorp/vault/pull/16087)]
* api/sys/internal/specs/openapi: support a new "dynamic" query parameter to generate generic mountpaths [[GH-15835](https://github.com/hashicorp/vault/pull/15835)]
* api: Fixed erroneous warnings of unrecognized parameters when unwrapping data. [[GH-16794](https://github.com/hashicorp/vault/pull/16794)]
* api: Fixed issue with internal/ui/mounts and internal/ui/mounts/(?P<path>.+) endpoints where it was not properly handling /auth/ [[GH-15552](https://github.com/hashicorp/vault/pull/15552)]
* api: properly handle switching to/from unix domain socket when changing client address [[GH-11904](https://github.com/hashicorp/vault/pull/11904)]
* auth/kerberos: Maintain headers set by the client [[GH-16636](https://github.com/hashicorp/vault/pull/16636)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core (enterprise): Fix bug where wrapping token lookup does not work within namespaces. [[GH-15583](https://github.com/hashicorp/vault/pull/15583)]
* core (enterprise): Fix creation of duplicate entities via alias metadata changes on local auth mounts.
* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/license (enterprise): Always remove stored license and allow unseal to complete when license cleanup fails
* core/quotas (enterprise): Fixed issue with improper counting of leases if lease count quota created after leases
* core/quotas: Added globbing functionality on the end of path suffix quota paths [[GH-16386](https://github.com/hashicorp/vault/pull/16386)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* core: Fixes parsing boolean values for ha_storage backends in config [[GH-15900](https://github.com/hashicorp/vault/pull/15900)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* debug: Fix panic when capturing debug bundle on Windows [[GH-14399](https://github.com/hashicorp/vault/pull/14399)]
* debug: Remove extra empty lines from vault.log when debug command is run [[GH-16714](https://github.com/hashicorp/vault/pull/16714)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the 
Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* openapi: Fixed issue where information about /auth/token endpoints was not present with explicit policy permissions [[GH-15552](https://github.com/hashicorp/vault/pull/15552)]
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* plugin/secrets/auth: Fix a bug with aliased backends such as aws-ec2 or generic [[GH-16673](https://github.com/hashicorp/vault/pull/16673)]
* quotas/lease-count: Fix lease-count quotas on mounts not properly being enforced when the lease generating request is a read [[GH-15735](https://github.com/hashicorp/vault/pull/15735)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* replication (enterprise): Fix data race in saveCheckpoint.
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* secrets/kv: Fix `kv get` issue preventing the ability to read a secret when providing a leading slash [[GH-16443](https://github.com/hashicorp/vault/pull/16443)]
* secrets/pki: Allow import of issuers without CRLSign KeyUsage; prohibit setting crl-signing usage on such issuers [[GH-16865](https://github.com/hashicorp/vault/pull/16865)]
* secrets/pki: Fix migration to properly handle mounts that contain only keys, no certificates [[GH-16813](https://github.com/hashicorp/vault/pull/16813)]
* secrets/pki: Ignore EC PARAMETER PEM blocks during issuer import (/config/ca, /issuers/import/*, and /intermediate/set-signed) [[GH-16721](https://github.com/hashicorp/vault/pull/16721)]
* secrets/pki: LIST issuers endpoint is now unauthenticated. [[GH-16830](https://github.com/hashicorp/vault/pull/16830)]
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix info tooltip submitting form [[GH-16659](https://github.com/hashicorp/vault/pull/16659)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]
* ui: Fixed bug where red spellcheck underline appears in sensitive/secret kv values when it should not appear [[GH-15681](https://github.com/hashicorp/vault/pull/15681)]
* ui: OIDC login type uses localStorage instead of sessionStorage [[GH-16170](https://github.com/hashicorp/vault/pull/16170)]
* vault: Fix a bug where duplicate policies could be added to an identity group. [[GH-15638](https://github.com/hashicorp/vault/pull/15638)]

## 1.11.3
### August 31, 2022

CHANGES:

* core: Bump Go version to 1.17.13.

IMPROVEMENTS:

* auth/kerberos: add `add_group_aliases` config to include LDAP groups in Vault group aliases [[GH-16890](https://github.com/hashicorp/vault/pull/16890)]
* auth/kerberos: add `remove_instance_name` parameter to the login CLI and the 
Kerberos config in Vault. This removes any instance names found in the keytab 
service principal name. [[GH-16594](https://github.com/hashicorp/vault/pull/16594)]
* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]
* storage/gcs: Add documentation explaining how to configure the gcs backend using environment variables instead of options in the configuration stanza [[GH-14455](https://github.com/hashicorp/vault/pull/14455)]

BUG FIXES:

* api: Fixed erroneous warnings of unrecognized parameters when unwrapping data. [[GH-16794](https://github.com/hashicorp/vault/pull/16794)]
* auth/gcp: Fixes the ability to reset the configuration's credentials to use application default credentials. [[GH-16523](https://github.com/hashicorp/vault/pull/16523)]
* auth/kerberos: Maintain headers set by the client [[GH-16636](https://github.com/hashicorp/vault/pull/16636)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core/license (enterprise): Always remove stored license and allow unseal to complete when license cleanup fails
* database/elasticsearch: Fixes a bug in boolean parsing for initialize [[GH-16526](https://github.com/hashicorp/vault/pull/16526)]
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the 
Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* plugin/secrets/auth: Fix a bug with aliased backends such as aws-ec2 or generic [[GH-16673](https://github.com/hashicorp/vault/pull/16673)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* secrets/pki: Fix migration to properly handle mounts that contain only keys, no certificates [[GH-16813](https://github.com/hashicorp/vault/pull/16813)]
* secrets/pki: Ignore EC PARAMETER PEM blocks during issuer import (/config/ca, /issuers/import/*, and /intermediate/set-signed) [[GH-16721](https://github.com/hashicorp/vault/pull/16721)]
* secrets/pki: LIST issuers endpoint is now unauthenticated. [[GH-16830](https://github.com/hashicorp/vault/pull/16830)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix info tooltip submitting form [[GH-16659](https://github.com/hashicorp/vault/pull/16659)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]
  
## 1.11.2
### August 2, 2022

IMPROVEMENTS:

* agent: Added `disable_keep_alives` configuration to disable keep alives in auto-auth, caching and templating. [[GH-16479](https://github.com/hashicorp/vault/pull/16479)]

BUG FIXES:

* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* secrets/kv: Fix `kv get` issue preventing the ability to read a secret when providing a leading slash [[GH-16443](https://github.com/hashicorp/vault/pull/16443)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]
  
## 1.11.1
### July 21, 2022

CHANGES:

* core: Bump Go version to 1.17.12.

IMPROVEMENTS:

* agent: Added `disable_idle_connections` configuration to disable leaving idle connections open in auto-auth, caching and templating. [[GH-15986](https://github.com/hashicorp/vault/pull/15986)]
* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* agent: Update consul-template for pkiCert bug fixes [[GH-16087](https://github.com/hashicorp/vault/pull/16087)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* kmip (enterprise): Return SecretData as supported Object Type.
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* transform (enterprise): Fix a bug in the handling of nested or unmatched capture groups in FPE transformations.
* ui: OIDC login type uses localStorage instead of sessionStorage [[GH-16170](https://github.com/hashicorp/vault/pull/16170)]
  
SECURITY:

* storage/raft (enterprise): Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HCSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]
  
## 1.11.0
### June 20, 2022

CHANGES:

* auth/aws: Add RoleSession to DisplayName when using assumeRole for authentication [[GH-14954](https://github.com/hashicorp/vault/pull/14954)]
* auth/kubernetes: If `kubernetes_ca_cert` is unset, and there is no pod-local CA available, an error will be surfaced when writing config instead of waiting for login. [[GH-15584](https://github.com/hashicorp/vault/pull/15584)]
* auth: Remove support for legacy MFA
(https://www.vaultproject.io/docs/v1.10.x/auth/mfa) [[GH-14869](https://github.com/hashicorp/vault/pull/14869)]
* core/fips: Disable and warn about entropy augmentation in FIPS 140-2 Inside mode [[GH-15858](https://github.com/hashicorp/vault/pull/15858)]
* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.17.11. [[GH-go-ver-1110](https://github.com/hashicorp/vault/pull/go-ver-1110)]
* database & storage: Change underlying driver library from [lib/pq](https://github.com/lib/pq) to [pgx](https://github.com/jackc/pgx). This change affects Redshift & Postgres database secrets engines, and CockroachDB & Postgres storage engines [[GH-15343](https://github.com/hashicorp/vault/pull/15343)]
* licensing (enterprise): Remove support for stored licenses and associated `sys/license` and `sys/license/signed`
endpoints in favor of [autoloaded licenses](https://www.vaultproject.io/docs/enterprise/license/autoloading).
* replication (enterprise): The `/sys/replication/performance/primary/mount-filter` endpoint has been removed. Please use [Paths Filter](https://www.vaultproject.io/api-docs/system/replication/replication-performance#create-paths-filter) instead.
* secret/pki: Remove unused signature_bits parameter from intermediate CSR generation; this parameter doesn't control the final certificate's signature algorithm selection as that is up to the signing CA [[GH-15478](https://github.com/hashicorp/vault/pull/15478)]
* secrets/kubernetes: Split `additional_metadata` into `extra_annotations` and `extra_labels` parameters [[GH-15655](https://github.com/hashicorp/vault/pull/15655)]
* secrets/pki: A new aliased api path (/pki/issuer/:issuer_ref/sign-self-issued)
providing the same functionality as the existing API(/pki/root/sign-self-issued)
does not require sudo capabilities but the latter still requires it in an
effort to maintain backwards compatibility. [[GH-15211](https://github.com/hashicorp/vault/pull/15211)]
* secrets/pki: Err on unknown role during sign-verbatim. [[GH-15543](https://github.com/hashicorp/vault/pull/15543)]
* secrets/pki: Existing CRL API (/pki/crl) now returns an X.509 v2 CRL instead
of a v1 CRL. [[GH-15100](https://github.com/hashicorp/vault/pull/15100)]
* secrets/pki: The `ca_chain` response field within issuing (/pki/issue/:role)
and signing APIs will now include the root CA certificate if the mount is
aware of it. [[GH-15155](https://github.com/hashicorp/vault/pull/15155)]
* secrets/pki: existing Delete Root API (pki/root) will now delete all issuers
and keys within the mount path. [[GH-15004](https://github.com/hashicorp/vault/pull/15004)]
* secrets/pki: existing Generate Root (pki/root/generate/:type),
Set Signed Intermediate (/pki/intermediate/set-signed) APIs will
add new issuers/keys to a mount instead of warning that an existing CA exists [[GH-14975](https://github.com/hashicorp/vault/pull/14975)]
* secrets/pki: the signed CA certificate from the sign-intermediate api will now appear within the ca_chain
response field along with the issuer's ca chain. [[GH-15524](https://github.com/hashicorp/vault/pull/15524)]
* ui: Upgrade Ember to version 3.28 [[GH-14763](https://github.com/hashicorp/vault/pull/14763)]

FEATURES:

* **Autopilot Improvements (Enterprise)**: Autopilot on Vault Enterprise now supports automated upgrades and redundancy zones when using integrated storage.
* **KeyMgmt UI**: Add UI support for managing the Key Management Secrets Engine [[GH-15523](https://github.com/hashicorp/vault/pull/15523)]
* **Kubernetes Secrets Engine**: This new secrets engine generates Kubernetes service account tokens, service accounts, role bindings, and roles dynamically. [[GH-15551](https://github.com/hashicorp/vault/pull/15551)]
* **Non-Disruptive Intermediate/Root Certificate Rotation**: This allows
import, generation and configuration of any number of keys and/or issuers
within a PKI mount, providing operators the ability to rotate certificates
in place without affecting existing client configurations. [[GH-15277](https://github.com/hashicorp/vault/pull/15277)]
* **Print minimum required policy for any command**: The global CLI flag `-output-policy` can now be used with any command to print out the minimum required policy HCL for that operation, including whether the given path requires the "sudo" capability. [[GH-14899](https://github.com/hashicorp/vault/pull/14899)]
* **Snowflake Database Plugin**: Adds ability to manage RSA key pair credentials for dynamic and static Snowflake users. [[GH-15376](https://github.com/hashicorp/vault/pull/15376)]
* **Transit BYOK**: Allow import of externally-generated keys into the Transit secrets engine. [[GH-15414](https://github.com/hashicorp/vault/pull/15414)]
* nomad: Bootstrap Nomad ACL system if no token is provided [[GH-12451](https://github.com/hashicorp/vault/pull/12451)]
* storage/dynamodb: Added `AWS_DYNAMODB_REGION` environment variable. [[GH-15054](https://github.com/hashicorp/vault/pull/15054)]

IMPROVEMENTS:

* activity: return nil response months in activity log API when no month data exists [[GH-15420](https://github.com/hashicorp/vault/pull/15420)]
* agent/auto-auth: Add `min_backoff` to the method stanza for configuring initial backoff duration. [[GH-15204](https://github.com/hashicorp/vault/pull/15204)]
* agent: Update consul-template to v0.29.0 [[GH-15293](https://github.com/hashicorp/vault/pull/15293)]
* agent: Upgrade hashicorp/consul-template version for sprig template functions and improved writeTo function [[GH-15092](https://github.com/hashicorp/vault/pull/15092)]
* api/monitor: Add log_format option to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* api: Add ability to pass certificate as PEM bytes to api.Client. [[GH-14753](https://github.com/hashicorp/vault/pull/14753)]
* api: Add context-aware functions to vault/api for each API wrapper function. [[GH-14388](https://github.com/hashicorp/vault/pull/14388)]
* api: Added MFALogin() for handling MFA flow when using login helpers. [[GH-14900](https://github.com/hashicorp/vault/pull/14900)]
* api: If the parameters supplied over the API payload are ignored due to not
being what the endpoints were expecting, or if the parameters supplied get
replaced by the values in the endpoint's path itself, warnings will be added to
the non-empty responses listing all the ignored and replaced parameters. [[GH-14962](https://github.com/hashicorp/vault/pull/14962)]
* api: KV helper methods to simplify the common use case of reading and writing KV secrets [[GH-15305](https://github.com/hashicorp/vault/pull/15305)]
* api: Provide a helper method WithNamespace to create a cloned client with a new NS [[GH-14963](https://github.com/hashicorp/vault/pull/14963)]
* api: Support VAULT_PROXY_ADDR environment variable to allow overriding the Vault client's HTTP proxy. [[GH-15377](https://github.com/hashicorp/vault/pull/15377)]
* api: Use the context passed to the api/auth Login helpers. [[GH-14775](https://github.com/hashicorp/vault/pull/14775)]
* api: make ListPlugins parse only known plugin types [[GH-15434](https://github.com/hashicorp/vault/pull/15434)]
* audit: Add a policy_results block into the audit log that contains the set of
policies that granted this request access. [[GH-15457](https://github.com/hashicorp/vault/pull/15457)]
* audit: Include mount_accessor in audit request and response logs [[GH-15342](https://github.com/hashicorp/vault/pull/15342)]
* audit: added entity_created boolean to audit log, set when login operations create an entity [[GH-15487](https://github.com/hashicorp/vault/pull/15487)]
* auth/aws: Add rsa2048 signature type to API [[GH-15719](https://github.com/hashicorp/vault/pull/15719)]
* auth/gcp: Enable the Google service endpoints used by the underlying client to be customized [[GH-15592](https://github.com/hashicorp/vault/pull/15592)]
* auth/gcp: Vault CLI now infers the service account email when running on Google Cloud [[GH-15592](https://github.com/hashicorp/vault/pull/15592)]
* auth/jwt: Adds ability to use JSON pointer syntax for the `user_claim` value. [[GH-15593](https://github.com/hashicorp/vault/pull/15593)]
* auth/okta: Add support for Google provider TOTP type in the Okta auth method [[GH-14985](https://github.com/hashicorp/vault/pull/14985)]
* auth/okta: Add support for performing [the number
challenge](https://help.okta.com/en-us/Content/Topics/Mobile/ov-admin-config.htm?cshid=csh-okta-verify-number-challenge-v1#enable-number-challenge)
during an Okta Verify push challenge [[GH-15361](https://github.com/hashicorp/vault/pull/15361)]
* auth: Globally scoped Login MFA method Get/List endpoints [[GH-15248](https://github.com/hashicorp/vault/pull/15248)]
* auth: enforce a rate limit for TOTP passcode validation attempts [[GH-14864](https://github.com/hashicorp/vault/pull/14864)]
* auth: forward cached MFA auth response to the leader using RPC instead of forwarding all login requests [[GH-15469](https://github.com/hashicorp/vault/pull/15469)]
* cli/debug: added support for retrieving metrics from DR clusters if `unauthenticated_metrics_access` is enabled [[GH-15316](https://github.com/hashicorp/vault/pull/15316)]
* cli/vault: warn when policy name contains upper-case letter [[GH-14670](https://github.com/hashicorp/vault/pull/14670)]
* cli: Alternative flag-based syntax for KV to mitigate confusion from automatically appended /data [[GH-14807](https://github.com/hashicorp/vault/pull/14807)]
* cockroachdb: add high-availability support [[GH-12965](https://github.com/hashicorp/vault/pull/12965)]
* command/debug: Add log_format flag to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* command: Support optional '-log-level' flag to be passed to 'operator migrate' command (defaults to info). Also support VAULT_LOG_LEVEL env var. [[GH-15405](https://github.com/hashicorp/vault/pull/15405)]
* command: Support the optional '-detailed' flag to be passed to 'vault list' command to show ListResponseWithInfo data. Also supports the VAULT_DETAILED env var. [[GH-15417](https://github.com/hashicorp/vault/pull/15417)]
* core (enterprise): Include `termination_time` in `sys/license/status` response
* core (enterprise): Include termination time in `license inspect` command output
* core,transit: Allow callers to choose random byte source including entropy augmentation sources for the sys/tools/random and transit/random endpoints. [[GH-15213](https://github.com/hashicorp/vault/pull/15213)]
* core/activity: Order month data in ascending order of timestamps [[GH-15259](https://github.com/hashicorp/vault/pull/15259)]
* core/activity: allow client counts to be precomputed and queried on non-contiguous chunks of data [[GH-15352](https://github.com/hashicorp/vault/pull/15352)]
* core/managed-keys (enterprise): Allow configuring the number of parallel operations to PKCS#11 managed keys.
* core: Add an export API for historical activity log data [[GH-15586](https://github.com/hashicorp/vault/pull/15586)]
* core: Add new DB methods that do not prepare statements. [[GH-15166](https://github.com/hashicorp/vault/pull/15166)]
* core: check uid and permissions of config dir, config file, plugin dir and plugin binaries [[GH-14817](https://github.com/hashicorp/vault/pull/14817)]
* core: Fix some identity data races found by Go race detector (no known impact yet). [[GH-15123](https://github.com/hashicorp/vault/pull/15123)]
* core: Include build date in `sys/seal-status` and `sys/version-history` endpoints. [[GH-14957](https://github.com/hashicorp/vault/pull/14957)]
* core: Upgrade github.org/x/crypto/ssh [[GH-15125](https://github.com/hashicorp/vault/pull/15125)]
* kmip (enterprise): Implement operations Query, Import, Encrypt and Decrypt. Improve operations Locate, Add Attribute, Get Attributes and Get Attribute List to handle most supported attributes.
* mfa/okta: migrate to use official Okta SDK [[GH-15355](https://github.com/hashicorp/vault/pull/15355)]
* sdk: Change OpenAPI code generator to extract request objects into /components/schemas and reference them by name. [[GH-14217](https://github.com/hashicorp/vault/pull/14217)]
* secrets/consul: Add support for Consul node-identities and service-identities [[GH-15295](https://github.com/hashicorp/vault/pull/15295)]
* secrets/consul: Vault is now able to automatically bootstrap the Consul ACL system. [[GH-10751](https://github.com/hashicorp/vault/pull/10751)]
* secrets/database/elasticsearch: Use the new /_security base API path instead of /_xpack/security when managing elasticsearch. [[GH-15614](https://github.com/hashicorp/vault/pull/15614)]
* secrets/pki: Add not_before_duration to root CA generation, intermediate CA signing paths. [[GH-14178](https://github.com/hashicorp/vault/pull/14178)]
* secrets/pki: Add support for CPS URLs and User Notice to Policy Information [[GH-15751](https://github.com/hashicorp/vault/pull/15751)]
* secrets/pki: Allow operators to control the issuing certificate behavior when
the requested TTL is beyond the NotAfter value of the signing certificate [[GH-15152](https://github.com/hashicorp/vault/pull/15152)]
* secrets/pki: Always return CRLs, URLs configurations, even if using the default value. [[GH-15470](https://github.com/hashicorp/vault/pull/15470)]
* secrets/pki: Enable Patch Functionality for Roles and Issuers (API only) [[GH-15510](https://github.com/hashicorp/vault/pull/15510)]
* secrets/pki: Have pki/sign-verbatim use the not_before_duration field defined in the role [[GH-15429](https://github.com/hashicorp/vault/pull/15429)]
* secrets/pki: Warn on empty Subject field during issuer generation (root/generate and root/sign-intermediate). [[GH-15494](https://github.com/hashicorp/vault/pull/15494)]
* secrets/pki: Warn on missing AIA access information when generating issuers (config/urls). [[GH-15509](https://github.com/hashicorp/vault/pull/15509)]
* secrets/pki: Warn when `generate_lease` and `no_store` are both set to `true` on requests. [[GH-14292](https://github.com/hashicorp/vault/pull/14292)]
* secrets/ssh: Add connection timeout of 1 minute for outbound SSH connection in deprecated Dynamic SSH Keys mode. [[GH-15440](https://github.com/hashicorp/vault/pull/15440)]
* secrets/ssh: Support for `add_before_duration` in SSH [[GH-15250](https://github.com/hashicorp/vault/pull/15250)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer
* storage/raft: Use larger timeouts at startup to reduce likelihood of inducing elections. [[GH-15042](https://github.com/hashicorp/vault/pull/15042)]
* ui: Allow namespace param to be parsed from state queryParam [[GH-15378](https://github.com/hashicorp/vault/pull/15378)]
* ui: Default auto-rotation period in transit is 30 days [[GH-15474](https://github.com/hashicorp/vault/pull/15474)]
* ui: Parse schema refs from OpenAPI [[GH-14508](https://github.com/hashicorp/vault/pull/14508)]
* ui: Remove stored license references [[GH-15513](https://github.com/hashicorp/vault/pull/15513)]
* ui: Remove storybook. [[GH-15074](https://github.com/hashicorp/vault/pull/15074)]
* ui: Replaces the IvyCodemirror wrapper with a custom ember modifier. [[GH-14659](https://github.com/hashicorp/vault/pull/14659)]
* website/docs: Add usage documentation for Kubernetes Secrets Engine [[GH-15527](https://github.com/hashicorp/vault/pull/15527)]
* website/docs: added a link to an Enigma secret plugin. [[GH-14389](https://github.com/hashicorp/vault/pull/14389)]

DEPRECATIONS:

* docs: Document removal of X.509 certificates with signatures who use SHA-1 in Vault 1.12 [[GH-15581](https://github.com/hashicorp/vault/pull/15581)]
* secrets/consul: Deprecate old parameters "token_type" and "policy" [[GH-15550](https://github.com/hashicorp/vault/pull/15550)]
* secrets/consul: Deprecate parameter "policies" in favor of "consul_policies" for consistency [[GH-15400](https://github.com/hashicorp/vault/pull/15400)]

BUG FIXES:

* Fixed panic when adding or modifying a Duo MFA Method in Enterprise
* agent: Fix log level mismatch between ERR and ERROR [[GH-14424](https://github.com/hashicorp/vault/pull/14424)]
* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* api: Fixes bug where OutputCurlString field was unintentionally being copied over during client cloning [[GH-14968](https://github.com/hashicorp/vault/pull/14968)]
* api: Respect increment value in grace period calculations in LifetimeWatcher [[GH-14836](https://github.com/hashicorp/vault/pull/14836)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* auth/kubernetes: Fix error code when using the wrong service account [[GH-15584](https://github.com/hashicorp/vault/pull/15584)]
* auth/ldap: The logic for setting the entity alias when `username_as_alias` is set
has been fixed. The previous behavior would make a request to the LDAP server to
get `user_attr` before discarding it and using the username instead. This would
make it impossible for a user to connect if this attribute was missing or had
multiple values, even though it would not be used anyway. This has been fixed
and the username is now used without making superfluous LDAP searches. [[GH-15525](https://github.com/hashicorp/vault/pull/15525)]
* auth: Fixed erroneous success message when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed erroneous token information being displayed when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed two-phase MFA information missing from table format when using vault login [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Prevent deleting a valid MFA method ID using the endpoint for a different MFA method type [[GH-15482](https://github.com/hashicorp/vault/pull/15482)]
* auth: forward requests subject to login MFA from perfStandby to Active node [[GH-15009](https://github.com/hashicorp/vault/pull/15009)]
* auth: load login MFA configuration upon restart [[GH-15261](https://github.com/hashicorp/vault/pull/15261)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* cli: kv get command now honors trailing spaces to retrieve secrets [[GH-15188](https://github.com/hashicorp/vault/pull/15188)]
* command: do not report listener and storage types as key not found warnings [[GH-15383](https://github.com/hashicorp/vault/pull/15383)]
* core (enterprise): Allow local alias create RPCs to persist alias metadata
* core (enterprise): Fix overcounting of lease count quota usage at startup.
* core (enterprise): Fix some races in merkle index flushing code found in testing
* core (enterprise): Handle additional edge cases reinitializing PKCS#11 libraries after login errors.
* core/config: Only ask the system about network interfaces when address configs contain a template having the format: {{ ... }} [[GH-15224](https://github.com/hashicorp/vault/pull/15224)]
* core/managed-keys (enterprise): Allow PKCS#11 managed keys to use 0 as a slot number
* core/metrics: Fix incorrect table size metric for local mounts [[GH-14755](https://github.com/hashicorp/vault/pull/14755)]
* core: Fix double counting for "route" metrics [[GH-12763](https://github.com/hashicorp/vault/pull/12763)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: Limit SSCT WAL checks on perf standbys to raft backends only [[GH-15879](https://github.com/hashicorp/vault/pull/15879)]
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: Prevent metrics generation from causing deadlocks. [[GH-15693](https://github.com/hashicorp/vault/pull/15693)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]
* core: renaming the environment variable VAULT_DISABLE_FILE_PERMISSIONS_CHECK to VAULT_ENABLE_FILE_PERMISSIONS_CHECK and adjusting the logic [[GH-15452](https://github.com/hashicorp/vault/pull/15452)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* identity: deduplicate policies when creating/updating identity groups [[GH-15055](https://github.com/hashicorp/vault/pull/15055)]
* mfa/okta: disable client side rate limiting causing delays in push notifications [[GH-15369](https://github.com/hashicorp/vault/pull/15369)]
* plugin: Fix a bug where plugin reload would falsely report success in certain scenarios. [[GH-15579](https://github.com/hashicorp/vault/pull/15579)]
* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* raft: Ensure initialMmapSize is set to 0 on Windows [[GH-14977](https://github.com/hashicorp/vault/pull/14977)]
* replication (enterprise): fix panic due to missing entity during invalidation of local aliases. [[GH-14622](https://github.com/hashicorp/vault/pull/14622)]
* sdk/cidrutil: Only check if cidr contains remote address for IP addresses [[GH-14487](https://github.com/hashicorp/vault/pull/14487)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]
* sdk: Fix OpenApi spec generator to remove duplicate sha_256 parameter [[GH-15163](https://github.com/hashicorp/vault/pull/15163)]
* secrets/database: Ensure that a `connection_url` password is redacted in all cases. [[GH-14744](https://github.com/hashicorp/vault/pull/14744)]
* secrets/kv: Fix issue preventing the ability to reset the `delete_version_after` key metadata field to 0s via HTTP `PATCH`. [[GH-15792](https://github.com/hashicorp/vault/pull/15792)]
* secrets/pki: CRLs on performance secondary clusters are now automatically
rebuilt upon changes to the list of issuers. [[GH-15179](https://github.com/hashicorp/vault/pull/15179)]
* secrets/pki: Fix handling of "any" key type with default zero signature bits value. [[GH-14875](https://github.com/hashicorp/vault/pull/14875)]
* secrets/pki: Fixed bug where larger SHA-2 hashes were truncated with shorter ECDSA CA certificates [[GH-14943](https://github.com/hashicorp/vault/pull/14943)]
* secrets/ssh: Convert role field not_before_duration to seconds before returning it [[GH-15559](https://github.com/hashicorp/vault/pull/15559)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* storage/raft: Forward autopilot state requests on perf standbys to active node. [[GH-15493](https://github.com/hashicorp/vault/pull/15493)]
* storage/raft: joining a node to a cluster now ignores any VAULT_NAMESPACE environment variable set on the server process [[GH-15519](https://github.com/hashicorp/vault/pull/15519)]
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not accepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fix KV secret showing in the edit form after a user creates a new version but doesn't have read capabilities [[GH-14794](https://github.com/hashicorp/vault/pull/14794)]
* ui: Fix inconsistent behavior in client count calendar widget [[GH-15789](https://github.com/hashicorp/vault/pull/15789)]
* ui: Fix issue where metadata tab is hidden even though policy grants access [[GH-15824](https://github.com/hashicorp/vault/pull/15824)]
* ui: Fix issue with KV not recomputing model when you changed versions. [[GH-14941](https://github.com/hashicorp/vault/pull/14941)]
* ui: Fixed client count timezone for start and end months [[GH-15167](https://github.com/hashicorp/vault/pull/15167)]
* ui: Fixed unsupported revocation statements field for DB roles [[GH-15573](https://github.com/hashicorp/vault/pull/15573)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging in with OIDC from a listed auth mounts tab [[GH-14916](https://github.com/hashicorp/vault/pull/14916)]
* ui: Revert using localStorage in favor of sessionStorage [[GH-15769](https://github.com/hashicorp/vault/pull/15769)]
* ui: Updated `leasId` to `leaseId` in the "Copy Credentials" section of "Generate AWS Credentials" [[GH-15685](https://github.com/hashicorp/vault/pull/15685)]
* ui: fix firefox inability to recognize file format of client count csv export [[GH-15364](https://github.com/hashicorp/vault/pull/15364)]
* ui: fix form validations ignoring default values and disabling submit button [[GH-15560](https://github.com/hashicorp/vault/pull/15560)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]

## 1.10.6
### August 31, 2022

CHANGES:

* core: Bump Go version to 1.17.13.

IMPROVEMENTS:

* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]

BUG FIXES:

* auth/gcp: Fixes the ability to reset the configuration's credentials to use application default credentials. [[GH-16524](https://github.com/hashicorp/vault/pull/16524)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the 
Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]
  
## 1.10.5
### July 21, 2022

CHANGES:

* core/fips: Disable and warn about entropy augmentation in FIPS 140-2 Inside mode [[GH-15858](https://github.com/hashicorp/vault/pull/15858)]
* core: Bump Go version to 1.17.12.

IMPROVEMENTS:

* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* core: Limit SSCT WAL checks on perf standbys to raft backends only [[GH-15879](https://github.com/hashicorp/vault/pull/15879)]
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* storage/raft (enterprise): Prevent unauthenticated voter status with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* transform (enterprise): Fix a bug in the handling of nested or unmatched capture groups in FPE transformations.
* ui: Fix issue where metadata tab is hidden even though policy grants access [[GH-15824](https://github.com/hashicorp/vault/pull/15824)]
* ui: Revert using localStorage in favor of sessionStorage [[GH-16169](https://github.com/hashicorp/vault/pull/16169)]
* ui: Updated `leasId` to `leaseId` in the "Copy Credentials" section of "Generate AWS Credentials" [[GH-15685](https://github.com/hashicorp/vault/pull/15685)]

## 1.10.4
### June 10, 2022

CHANGES:

* core: Bump Go version to 1.17.11. [[GH-go-ver-1104](https://github.com/hashicorp/vault/pull/go-ver-1104)]

IMPROVEMENTS:

* api/monitor: Add log_format option to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* auth: Globally scoped Login MFA method Get/List endpoints [[GH-15248](https://github.com/hashicorp/vault/pull/15248)]
* auth: forward cached MFA auth response to the leader using RPC instead of forwarding all login requests [[GH-15469](https://github.com/hashicorp/vault/pull/15469)]
* cli/debug: added support for retrieving metrics from DR clusters if `unauthenticated_metrics_access` is enabled [[GH-15316](https://github.com/hashicorp/vault/pull/15316)]
* command/debug: Add log_format flag to allow for logs to be emitted in JSON format [[GH-15536](https://github.com/hashicorp/vault/pull/15536)]
* core: Fix some identity data races found by Go race detector (no known impact yet). [[GH-15123](https://github.com/hashicorp/vault/pull/15123)]
* storage/raft: Use larger timeouts at startup to reduce likelihood of inducing elections. [[GH-15042](https://github.com/hashicorp/vault/pull/15042)]
* ui: Allow namespace param to be parsed from state queryParam [[GH-15378](https://github.com/hashicorp/vault/pull/15378)]

BUG FIXES:

* agent: Redact auto auth token from renew endpoints [[GH-15380](https://github.com/hashicorp/vault/pull/15380)]
* auth/kubernetes: Fix error code when using the wrong service account [[GH-15585](https://github.com/hashicorp/vault/pull/15585)]
* auth/ldap: The logic for setting the entity alias when `username_as_alias` is set
has been fixed. The previous behavior would make a request to the LDAP server to
get `user_attr` before discarding it and using the username instead. This would
make it impossible for a user to connect if this attribute was missing or had
multiple values, even though it would not be used anyway. This has been fixed
and the username is now used without making superfluous LDAP searches. [[GH-15525](https://github.com/hashicorp/vault/pull/15525)]
* auth: Fixed erroneous success message when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed erroneous token information being displayed when using vault login in case of two-phase MFA [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Fixed two-phase MFA information missing from table format when using vault login [[GH-15428](https://github.com/hashicorp/vault/pull/15428)]
* auth: Prevent deleting a valid MFA method ID using the endpoint for a different MFA method type [[GH-15482](https://github.com/hashicorp/vault/pull/15482)]
* core (enterprise): Fix overcounting of lease count quota usage at startup.
* core: Prevent changing file permissions of audit logs when mode 0000 is used. [[GH-15759](https://github.com/hashicorp/vault/pull/15759)]
* core: Prevent metrics generation from causing deadlocks. [[GH-15693](https://github.com/hashicorp/vault/pull/15693)]
* core: fixed systemd reloading notification [[GH-15041](https://github.com/hashicorp/vault/pull/15041)]
* mfa/okta: disable client side rate limiting causing delays in push notifications [[GH-15369](https://github.com/hashicorp/vault/pull/15369)]
* storage/raft (enterprise):  Auto-snapshot configuration now forbids slashes in file prefixes for all types, and "/" in path prefix for local storage type.  Strip leading prefix in path prefix for AWS.  Improve error handling/reporting.
* transform (enterprise): Fix non-overridable column default value causing tokenization tokens to expire prematurely when using the MySQL storage backend.
* ui: Fix inconsistent behavior in client count calendar widget [[GH-15789](https://github.com/hashicorp/vault/pull/15789)]
* ui: Fixed client count timezone for start and end months [[GH-15167](https://github.com/hashicorp/vault/pull/15167)]
* ui: fix firefox inability to recognize file format of client count csv export [[GH-15364](https://github.com/hashicorp/vault/pull/15364)]

## 1.10.3
### May 11, 2022

SECURITY:
* auth: A vulnerability was identified in Vault and Vault Enterprise (“Vault”) from 1.10.0 to 1.10.2 where MFA may not be enforced on user logins after a server restart. This vulnerability, CVE-2022-30689, was fixed in Vault 1.10.3.

BUG FIXES:

* auth: load login MFA configuration upon restart [[GH-15261](https://github.com/hashicorp/vault/pull/15261)]
* core/config: Only ask the system about network interfaces when address configs contain a template having the format: {{ ... }} [[GH-15224](https://github.com/hashicorp/vault/pull/15224)]
* core: pre-calculate namespace specific paths when tainting a route during postUnseal [[GH-15067](https://github.com/hashicorp/vault/pull/15067)]

## 1.10.2
### April 29, 2022

BUG FIXES:

* raft: fix Raft TLS key rotation panic that occurs if active key is more than 24 hours old [[GH-15156](https://github.com/hashicorp/vault/pull/15156)]
* sdk: Fix OpenApi spec generator to properly convert TypeInt64 to OAS supported int64 [[GH-15104](https://github.com/hashicorp/vault/pull/15104)]

## 1.10.1
### April 22, 2022

CHANGES:

* core: A request that fails path validation due to relative path check will now be responded to with a 400 rather than 500. [[GH-14328](https://github.com/hashicorp/vault/pull/14328)]
* core: Bump Go version to 1.17.9. [[GH-15044](https://github.com/hashicorp/vault/pull/15044)]

IMPROVEMENTS:

* agent: Upgrade hashicorp/consul-template version for sprig template functions and improved writeTo function [[GH-15092](https://github.com/hashicorp/vault/pull/15092)]
* auth: enforce a rate limit for TOTP passcode validation attempts [[GH-14864](https://github.com/hashicorp/vault/pull/14864)]
* cli/vault: warn when policy name contains upper-case letter [[GH-14670](https://github.com/hashicorp/vault/pull/14670)]
* cockroachdb: add high-availability support [[GH-12965](https://github.com/hashicorp/vault/pull/12965)]
* sentinel (enterprise): Upgrade sentinel to [v0.18.5](https://docs.hashicorp.com/sentinel/changelog#0-18-5-january-14-2022) to avoid potential naming collisions in the remote installer

BUG FIXES:

* Fixed panic when adding or modifying a Duo MFA Method in Enterprise
* agent: Fix log level mismatch between ERR and ERROR [[GH-14424](https://github.com/hashicorp/vault/pull/14424)]
* api/sys/raft: Update RaftSnapshotRestore to use net/http client allowing bodies larger than allocated memory to be streamed [[GH-14269](https://github.com/hashicorp/vault/pull/14269)]
* api: Respect increment value in grace period calculations in LifetimeWatcher [[GH-14836](https://github.com/hashicorp/vault/pull/14836)]
* auth/approle: Add maximum length for input values that result in SHA56 HMAC calculation [[GH-14746](https://github.com/hashicorp/vault/pull/14746)]
* auth: forward requests subject to login MFA from perfStandby to Active node [[GH-15009](https://github.com/hashicorp/vault/pull/15009)]
* cassandra: Update gocql Cassandra client to fix "no hosts available in the pool" error [[GH-14973](https://github.com/hashicorp/vault/pull/14973)]
* cli: Fix panic caused by parsing key=value fields whose value is a single backslash [[GH-14523](https://github.com/hashicorp/vault/pull/14523)]
* core (enterprise): Allow local alias create RPCs to persist alias metadata [[GH-changelog:_2747](https://github.com/hashicorp/vault/pull/changelog:_2747)]
* core/managed-keys (enterprise): Allow PKCS#11 managed keys to use 0 as a slot number
* core/metrics: Fix incorrect table size metric for local mounts [[GH-14755](https://github.com/hashicorp/vault/pull/14755)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited integers [[GH-15072](https://github.com/hashicorp/vault/pull/15072)]
* core: Fix panic caused by parsing JSON integers for fields defined as comma-delimited strings [[GH-14522](https://github.com/hashicorp/vault/pull/14522)]
* core: Fix panic caused by parsing policies with empty slice values. [[GH-14501](https://github.com/hashicorp/vault/pull/14501)]
* core: Fix panic for help request URL paths without /v1/ prefix [[GH-14704](https://github.com/hashicorp/vault/pull/14704)]
* core: fixing excessive unix file permissions [[GH-14791](https://github.com/hashicorp/vault/pull/14791)]
* core: fixing excessive unix file permissions on dir, files and archive created by vault debug command [[GH-14846](https://github.com/hashicorp/vault/pull/14846)]
* core: report unused or redundant keys in server configuration [[GH-14752](https://github.com/hashicorp/vault/pull/14752)]
* core: time.After() used in a select statement can lead to memory leak [[GH-14814](https://github.com/hashicorp/vault/pull/14814)]
* raft: Ensure initialMmapSize is set to 0 on Windows [[GH-14977](https://github.com/hashicorp/vault/pull/14977)]
* replication (enterprise): fix panic due to missing entity during invalidation of local aliases. [[GH-14622](https://github.com/hashicorp/vault/pull/14622)]
* secrets/database: Ensure that a `connection_url` password is redacted in all cases. [[GH-14744](https://github.com/hashicorp/vault/pull/14744)]
* secrets/pki: Fix handling of "any" key type with default zero signature bits value. [[GH-14875](https://github.com/hashicorp/vault/pull/14875)]
* secrets/pki: Fixed bug where larger SHA-2 hashes were truncated with shorter ECDSA CA certificates [[GH-14943](https://github.com/hashicorp/vault/pull/14943)]
* ui: Fix Generated Token's Policies helpText to clarify that comma separated values are not excepted in this field. [[GH-15046](https://github.com/hashicorp/vault/pull/15046)]
* ui: Fixes edit auth method capabilities issue [[GH-14966](https://github.com/hashicorp/vault/pull/14966)]
* ui: Fixes issue logging in with OIDC from a listed auth mounts tab [[GH-14916](https://github.com/hashicorp/vault/pull/14916)]
* ui: fix search-select component showing blank selections when editing group member entity [[GH-15058](https://github.com/hashicorp/vault/pull/15058)]
* ui: masked values no longer give away length or location of special characters [[GH-15025](https://github.com/hashicorp/vault/pull/15025)]

## 1.10.0
### March 23, 2022

CHANGES:

* core (enterprise): requests with newly generated tokens to perf standbys which are lagging behind the active node return http 412 instead of 400/403/50x.
* core: Changes the unit of `default_lease_ttl` and `max_lease_ttl` values returned by
the `/sys/config/state/sanitized` endpoint from nanoseconds to seconds. [[GH-14206](https://github.com/hashicorp/vault/pull/14206)]
* core: Bump Go version to 1.17.7. [[GH-14232](https://github.com/hashicorp/vault/pull/14232)]
* plugin/database: The return value from `POST /database/config/:name` has been updated to "204 No Content" [[GH-14033](https://github.com/hashicorp/vault/pull/14033)]
* secrets/azure: Changes the configuration parameter `use_microsoft_graph_api` to use the Microsoft 
Graph API by default. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]
* storage/etcd: Remove support for v2. [[GH-14193](https://github.com/hashicorp/vault/pull/14193)]
* ui: Upgrade Ember to version 3.24 [[GH-13443](https://github.com/hashicorp/vault/pull/13443)]

FEATURES:

* **Database plugin multiplexing**: manage multiple database connections with a single plugin process [[GH-14033](https://github.com/hashicorp/vault/pull/14033)]
* **Login MFA**: Single and two phase MFA is now available when authenticating to Vault. [[GH-14025](https://github.com/hashicorp/vault/pull/14025)]
* **Mount Migration**: Vault supports moving secrets and auth mounts both within and across namespaces.
* **Postgres in the UI**: Postgres DB is now supported by the UI [[GH-12945](https://github.com/hashicorp/vault/pull/12945)]
* **Report in-flight requests**: Adding a trace capability to show in-flight requests, and a new gauge metric to show the total number of in-flight requests [[GH-13024](https://github.com/hashicorp/vault/pull/13024)]
* **Server Side Consistent Tokens**: Service tokens have been updated to be longer (a minimum of 95 bytes) and token prefixes for all token types are updated from s., b., and r. to hvs., hvb., and hvr. for service, batch, and recovery tokens respectively. Vault clusters with integrated storage will now have read-after-write consistency by default. [[GH-14109](https://github.com/hashicorp/vault/pull/14109)]
* **Transit SHA-3 Support**: Add support for SHA-3 in the Transit backend. [[GH-13367](https://github.com/hashicorp/vault/pull/13367)]
* **Transit Time-Based Key Autorotation**: Add support for automatic, time-based key rotation to transit secrets engine, including in the UI. [[GH-13691](https://github.com/hashicorp/vault/pull/13691)]
* **UI Client Count Improvements**: Restructures client count dashboard, making use of billing start date to improve accuracy. Adds mount-level distribution and filtering. [[GH-client-counts](https://github.com/hashicorp/vault/pull/client-counts)]
* **Agent Telemetry**: The Vault Agent can now collect and return telemetry information at the `/agent/v1/metrics` endpoint.

IMPROVEMENTS:

* agent: Adds ability to configure specific user-assigned managed identities for Azure auto-auth. [[GH-14214](https://github.com/hashicorp/vault/pull/14214)]
* agent: The `agent/v1/quit` endpoint can now be used to stop the Vault Agent remotely [[GH-14223](https://github.com/hashicorp/vault/pull/14223)]
* api: Allow cloning `api.Client` tokens via `api.Config.CloneToken` or `api.Client.SetCloneToken()`. [[GH-13515](https://github.com/hashicorp/vault/pull/13515)]
* api: Define constants for X-Vault-Forward and X-Vault-Inconsistent headers [[GH-14067](https://github.com/hashicorp/vault/pull/14067)]
* api: Implements Login method in Go client libraries for GCP and Azure auth methods [[GH-13022](https://github.com/hashicorp/vault/pull/13022)]
* api: Implements Login method in Go client libraries for LDAP auth methods [[GH-13841](https://github.com/hashicorp/vault/pull/13841)]
* api: Trim newline character from wrapping token in logical.Unwrap from the api package [[GH-13044](https://github.com/hashicorp/vault/pull/13044)]
* api: add api method for modifying raft autopilot configuration [[GH-12428](https://github.com/hashicorp/vault/pull/12428)]
* api: respect WithWrappingToken() option during AppRole login authentication when used with secret ID specified from environment or from string [[GH-13241](https://github.com/hashicorp/vault/pull/13241)]
* audit: The audit logs now contain the port used by the client [[GH-12790](https://github.com/hashicorp/vault/pull/12790)]
* auth/aws: Enable region detection in the CLI by specifying the region as `auto` [[GH-14051](https://github.com/hashicorp/vault/pull/14051)]
* auth/cert: Add certificate extensions as metadata [[GH-13348](https://github.com/hashicorp/vault/pull/13348)]
* auth/jwt: The Authorization Code flow makes use of the Proof Key for Code Exchange (PKCE) extension. [[GH-13365](https://github.com/hashicorp/vault/pull/13365)]
* auth/kubernetes: Added support for dynamically reloading short-lived tokens for better Kubernetes 1.21+ compatibility [[GH-13595](https://github.com/hashicorp/vault/pull/13595)]
* auth/ldap: Add a response warning and server log whenever the config is accessed
if `userfilter` doesn't consider `userattr` [[GH-14095](https://github.com/hashicorp/vault/pull/14095)]
* auth/ldap: Add username to alias metadata [[GH-13669](https://github.com/hashicorp/vault/pull/13669)]
* auth/ldap: Add username_as_alias configurable to change how aliases are named [[GH-14324](https://github.com/hashicorp/vault/pull/14324)]
* auth/okta: Update [okta-sdk-golang](https://github.com/okta/okta-sdk-golang) dependency to version v2.9.1 for improved request backoff handling [[GH-13439](https://github.com/hashicorp/vault/pull/13439)]
* auth/token: The `auth/token/revoke-accessor` endpoint is now idempotent and will
not error out if the token has already been revoked. [[GH-13661](https://github.com/hashicorp/vault/pull/13661)]
* auth: reading `sys/auth/:path` now returns the configuration for the auth engine mounted at the given path [[GH-12793](https://github.com/hashicorp/vault/pull/12793)]
* cli: interactive CLI for login mfa [[GH-14131](https://github.com/hashicorp/vault/pull/14131)]
* command (enterprise): "vault license get" now uses non-deprecated endpoint /sys/license/status
* core/ha: Add new mechanism for keeping track of peers talking to active node, and new 'operator members' command to view them. [[GH-13292](https://github.com/hashicorp/vault/pull/13292)]
* core/identity: Support updating an alias' `custom_metadata` to be empty. [[GH-13395](https://github.com/hashicorp/vault/pull/13395)]
* core/pki: Support Y10K value in notAfter field to be compliant with IEEE 802.1AR-2018 standard [[GH-12795](https://github.com/hashicorp/vault/pull/12795)]
* core/pki: Support Y10K value in notAfter field when signing non-CA certificates [[GH-13736](https://github.com/hashicorp/vault/pull/13736)]
* core: Add duration and start_time to completed requests log entries [[GH-13682](https://github.com/hashicorp/vault/pull/13682)]
* core: Add support to list password policies at `sys/policies/password` [[GH-12787](https://github.com/hashicorp/vault/pull/12787)]
* core: Add support to list version history via API at `sys/version-history` and via CLI with `vault version-history` [[GH-13766](https://github.com/hashicorp/vault/pull/13766)]
* core: Fixes code scanning alerts [[GH-13667](https://github.com/hashicorp/vault/pull/13667)]
* core: Periodically test the health of connectivity to auto-seal backends [[GH-13078](https://github.com/hashicorp/vault/pull/13078)]
* core: Reading `sys/mounts/:path` now returns the configuration for the secret engine at the given path [[GH-12792](https://github.com/hashicorp/vault/pull/12792)]
* core: Replace "master key" terminology with "root key" [[GH-13324](https://github.com/hashicorp/vault/pull/13324)]
* core: Small changes to ensure goroutines terminate in tests [[GH-14197](https://github.com/hashicorp/vault/pull/14197)]
* core: Systemd unit file included with the Linux packages now sets the service type to notify. [[GH-14385](https://github.com/hashicorp/vault/pull/14385)]
* core: Update github.com/prometheus/client_golang to fix security vulnerability CVE-2022-21698. [[GH-14190](https://github.com/hashicorp/vault/pull/14190)]
* core: Vault now supports the PROXY protocol v2. Support for UNKNOWN connections
has also been added to the PROXY protocol v1. [[GH-13540](https://github.com/hashicorp/vault/pull/13540)]
* http (enterprise): Serve /sys/license/status endpoint within namespaces
* identity/oidc: Adds a default OIDC provider [[GH-14119](https://github.com/hashicorp/vault/pull/14119)]
* identity/oidc: Adds a default key for OIDC clients [[GH-14119](https://github.com/hashicorp/vault/pull/14119)]
* identity/oidc: Adds an `allow_all` assignment that permits all entities to authenticate via an OIDC client [[GH-14119](https://github.com/hashicorp/vault/pull/14119)]
* identity/oidc: Adds proof key for code exchange (PKCE) support to OIDC providers. [[GH-13917](https://github.com/hashicorp/vault/pull/13917)]
* sdk: Add helper for decoding root tokens [[GH-10505](https://github.com/hashicorp/vault/pull/10505)]
* secrets/azure: Adds support for rotate-root. [#70](https://github.com/hashicorp/vault-plugin-secrets-azure/pull/70) [[GH-13034](https://github.com/hashicorp/vault/pull/13034)]
* secrets/consul: Add support for consul enterprise namespaces and admin partitions. [[GH-13850](https://github.com/hashicorp/vault/pull/13850)]
* secrets/consul: Add support for consul roles. [[GH-14014](https://github.com/hashicorp/vault/pull/14014)]
* secrets/database/influxdb: Switch/upgrade to the `influxdb1-client` module [[GH-12262](https://github.com/hashicorp/vault/pull/12262)]
* secrets/database: Add database configuration parameter 'disable_escaping' for username and password when connecting to a database. [[GH-13414](https://github.com/hashicorp/vault/pull/13414)]
* secrets/kv: add full secret path output to table-formatted responses [[GH-14301](https://github.com/hashicorp/vault/pull/14301)]
* secrets/kv: add patch support for KVv2 key metadata [[GH-13215](https://github.com/hashicorp/vault/pull/13215)]
* secrets/kv: add subkeys endpoint to retrieve a secret's stucture without its values [[GH-13893](https://github.com/hashicorp/vault/pull/13893)]
* secrets/pki: Add ability to fetch individual certificate as DER or PEM [[GH-10948](https://github.com/hashicorp/vault/pull/10948)]
* secrets/pki: Add count and duration metrics to PKI issue and revoke calls. [[GH-13889](https://github.com/hashicorp/vault/pull/13889)]
* secrets/pki: Add error handling for error types other than UserError or InternalError [[GH-14195](https://github.com/hashicorp/vault/pull/14195)]
* secrets/pki: Allow URI SAN templates in allowed_uri_sans when allowed_uri_sans_template is set to true. [[GH-10249](https://github.com/hashicorp/vault/pull/10249)]
* secrets/pki: Allow other_sans in sign-intermediate and sign-verbatim [[GH-13958](https://github.com/hashicorp/vault/pull/13958)]
* secrets/pki: Calculate the Subject Key Identifier as suggested in [RFC 5280, Section 4.2.1.2](https://datatracker.ietf.org/doc/html/rfc5280#section-4.2.1.2). [[GH-11218](https://github.com/hashicorp/vault/pull/11218)]
* secrets/pki: Restrict issuance of wildcard certificates via role parameter (`allow_wildcard_certificates`) [[GH-14238](https://github.com/hashicorp/vault/pull/14238)]
* secrets/pki: Return complete chain (in `ca_chain` field) on calls to `pki/cert/ca_chain` [[GH-13935](https://github.com/hashicorp/vault/pull/13935)]
* secrets/pki: Use application/pem-certificate-chain for PEM certificates, application/x-pem-file for PEM CRLs [[GH-13927](https://github.com/hashicorp/vault/pull/13927)]
* secrets/pki: select appropriate signature algorithm for ECDSA signature on certificates. [[GH-11216](https://github.com/hashicorp/vault/pull/11216)]
* secrets/ssh: Add support for generating non-RSA SSH CAs [[GH-14008](https://github.com/hashicorp/vault/pull/14008)]
* secrets/ssh: Allow specifying multiple approved key lengths for a single algorithm [[GH-13991](https://github.com/hashicorp/vault/pull/13991)]
* secrets/ssh: Use secure default for algorithm signer (rsa-sha2-256) with RSA SSH CA keys on new roles [[GH-14006](https://github.com/hashicorp/vault/pull/14006)]
* secrets/transit: Don't abort transit encrypt or decrypt batches on single item failure. [[GH-13111](https://github.com/hashicorp/vault/pull/13111)]
* storage/aerospike: Upgrade `aerospike-client-go` to v5.6.0. [[GH-12165](https://github.com/hashicorp/vault/pull/12165)]
* storage/raft: Set InitialMmapSize to 100GB on 64bit architectures [[GH-13178](https://github.com/hashicorp/vault/pull/13178)]
* storage/raft: When using retry_join stanzas, join against all of them in parallel. [[GH-13606](https://github.com/hashicorp/vault/pull/13606)]
* sys/raw: Enhance sys/raw to read and write values that cannot be encoded in json. [[GH-13537](https://github.com/hashicorp/vault/pull/13537)]
* ui: Add support for ECDSA and Ed25519 certificate views [[GH-13894](https://github.com/hashicorp/vault/pull/13894)]
* ui: Add version diff view for KV V2 [[GH-13000](https://github.com/hashicorp/vault/pull/13000)]
* ui: Added client side paging for namespace list view [[GH-13195](https://github.com/hashicorp/vault/pull/13195)]
* ui: Adds flight icons to UI [[GH-12976](https://github.com/hashicorp/vault/pull/12976)]
* ui: Adds multi-factor authentication support [[GH-14049](https://github.com/hashicorp/vault/pull/14049)]
* ui: Allow static role credential rotation in Database secrets engines [[GH-14268](https://github.com/hashicorp/vault/pull/14268)]
* ui: Display badge for all versions in secrets engine header [[GH-13015](https://github.com/hashicorp/vault/pull/13015)]
* ui: Swap browser localStorage in favor of sessionStorage [[GH-14054](https://github.com/hashicorp/vault/pull/14054)]
* ui: The integrated web terminal now accepts both `-f` and `--force` as aliases
for `-force` for the `write` command. [[GH-13683](https://github.com/hashicorp/vault/pull/13683)]
* ui: Transform advanced templating with encode/decode format support [[GH-13908](https://github.com/hashicorp/vault/pull/13908)]
* ui: Updates ember blueprints to glimmer components [[GH-13149](https://github.com/hashicorp/vault/pull/13149)]
* ui: customizes empty state messages for transit and transform [[GH-13090](https://github.com/hashicorp/vault/pull/13090)]

BUG FIXES:

* Fixed bug where auth method only considers system-identity when multiple identities are available. [#50](https://github.com/hashicorp/vault-plugin-auth-azure/pull/50) [[GH-14138](https://github.com/hashicorp/vault/pull/14138)]
* activity log (enterprise): allow partial monthly client count to be accessed from namespaces [[GH-13086](https://github.com/hashicorp/vault/pull/13086)]
* agent: Fixes bug where vault agent is unaware of the namespace in the config when wrapping token
* api/client: Fixes an issue where the `replicateStateStore` was being set to `nil` upon consecutive calls to `client.SetReadYourWrites(true)`. [[GH-13486](https://github.com/hashicorp/vault/pull/13486)]
* auth/approle: Fix regression where unset cidrlist is returned as nil instead of zero-length array. [[GH-13235](https://github.com/hashicorp/vault/pull/13235)]
* auth/approle: Fix wrapping of nil errors in `login` endpoint [[GH-14107](https://github.com/hashicorp/vault/pull/14107)]
* auth/github: Use the Organization ID instead of the Organization name to verify the org membership. [[GH-13332](https://github.com/hashicorp/vault/pull/13332)]
* auth/kubernetes: Properly handle the migration of role storage entries containing an empty `alias_name_source` [[GH-13925](https://github.com/hashicorp/vault/pull/13925)]
* auth/kubernetes: ensure valid entity alias names created for projected volume tokens [[GH-14144](https://github.com/hashicorp/vault/pull/14144)]
* auth/oidc: Fixes OIDC auth from the Vault UI when using the implicit flow and `form_post` response mode. [[GH-13492](https://github.com/hashicorp/vault/pull/13492)]
* cli: Fix using kv patch with older server versions that don't support HTTP PATCH. [[GH-13615](https://github.com/hashicorp/vault/pull/13615)]
* core (enterprise): Fix a data race in logshipper.
* core (enterprise): Workaround AWS CloudHSM v5 SDK issue not allowing read-only sessions
* core/api: Fix overwriting of request headers when using JSONMergePatch. [[GH-14222](https://github.com/hashicorp/vault/pull/14222)]
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13093](https://github.com/hashicorp/vault/pull/13093)]
* core/identity: Address a data race condition between local updates to aliases and invalidations [[GH-13476](https://github.com/hashicorp/vault/pull/13476)]
* core/token: Fix null token panic from 'v1/auth/token/' endpoints and return proper error response. [[GH-13233](https://github.com/hashicorp/vault/pull/13233)]
* core/token: Fix null token_type panic resulting from 'v1/auth/token/roles/{role_name}' endpoint [[GH-13236](https://github.com/hashicorp/vault/pull/13236)]
* core: Fix warnings logged on perf standbys re stored versions [[GH-13042](https://github.com/hashicorp/vault/pull/13042)]
* core: `-output-curl-string` now properly sets cURL options for client and CA
certificates. [[GH-13660](https://github.com/hashicorp/vault/pull/13660)]
* core: add support for go-sockaddr templates in the top-level cluster_addr field [[GH-13678](https://github.com/hashicorp/vault/pull/13678)]
* core: authentication to "login" endpoint for non-existent mount path returns permission denied with status code 403 [[GH-13162](https://github.com/hashicorp/vault/pull/13162)]
* core: revert some unintentionally downgraded dependencies from 1.9.0-rc1 [[GH-13168](https://github.com/hashicorp/vault/pull/13168)]
* ha (enterprise): Prevents performance standby nodes from serving and caching stale data immediately after performance standby election completes
* http (enterprise): Always forward internal/counters endpoints from perf standbys to active node
* http:Fix /sys/monitor endpoint returning streaming not supported [[GH-13200](https://github.com/hashicorp/vault/pull/13200)]
* identity/oidc: Adds support for port-agnostic validation of loopback IP redirect URIs. [[GH-13871](https://github.com/hashicorp/vault/pull/13871)]
* identity/oidc: Check for a nil signing key on rotation to prevent panics. [[GH-13716](https://github.com/hashicorp/vault/pull/13716)]
* identity/oidc: Fixes inherited group membership when evaluating client assignments [[GH-14013](https://github.com/hashicorp/vault/pull/14013)]
* identity/oidc: Fixes potential write to readonly storage on performance secondary clusters during key rotation [[GH-14426](https://github.com/hashicorp/vault/pull/14426)]
* identity/oidc: Make the `nonce` parameter optional for the Authorization Endpoint of OIDC providers. [[GH-13231](https://github.com/hashicorp/vault/pull/13231)]
* identity/token: Fixes a bug where duplicate public keys could appear in the .well-known JWKS [[GH-14543](https://github.com/hashicorp/vault/pull/14543)]
* identity: Fix possible nil pointer dereference. [[GH-13318](https://github.com/hashicorp/vault/pull/13318)]
* identity: Fix regression preventing startup when aliases were created pre-1.9. [[GH-13169](https://github.com/hashicorp/vault/pull/13169)]
* identity: Fixes a panic in the OIDC key rotation due to a missing nil check. [[GH-13298](https://github.com/hashicorp/vault/pull/13298)]
* kmip (enterprise): Fix locate by name operations fail to find key after a rekey operation.
* licensing (enterprise): Revert accidental inclusion of the TDE feature from the `prem` build.
* metrics/autosnapshots (enterprise) : Fix bug that could cause
vault.autosnapshots.save.errors to not be incremented when there is an
autosnapshot save error.
* physical/mysql: Create table with wider `vault_key` column when initializing database tables. [[GH-14231](https://github.com/hashicorp/vault/pull/14231)]
* plugin/couchbase: Fix an issue in which the locking patterns did not allow parallel requests. [[GH-13033](https://github.com/hashicorp/vault/pull/13033)]
* replication (enterprise): When using encrypted secondary tokens, only clear the
private key after a successful connection to the primary cluster
* sdk/framework: Generate proper OpenAPI specs for path patterns that use an alternation as the root. [[GH-13487](https://github.com/hashicorp/vault/pull/13487)]
* sdk/helper/ldaputil: properly escape a trailing escape character to prevent panics. [[GH-13452](https://github.com/hashicorp/vault/pull/13452)]
* sdk/queue: move lock before length check to prevent panics. [[GH-13146](https://github.com/hashicorp/vault/pull/13146)]
* sdk: Fixes OpenAPI to distinguish between paths that can do only List, or both List and Read. [[GH-13643](https://github.com/hashicorp/vault/pull/13643)]
* secrets/azure: Fixed bug where Azure environment did not change Graph URL [[GH-13973](https://github.com/hashicorp/vault/pull/13973)]
* secrets/azure: Fixes service principal generation when assigning roles that have [DataActions](https://docs.microsoft.com/en-us/azure/role-based-access-control/role-definitions#dataactions). [[GH-13277](https://github.com/hashicorp/vault/pull/13277)]
* secrets/azure: Fixes the [rotate root](https://www.vaultproject.io/api-docs/secret/azure#rotate-root) 
operation for upgraded configurations with a `root_password_ttl` of zero. [[GH-14130](https://github.com/hashicorp/vault/pull/14130)]
* secrets/database/cassandra: change connect_timeout to 5s as documentation says [[GH-12443](https://github.com/hashicorp/vault/pull/12443)]
* secrets/database/mssql: Accept a boolean for `contained_db`, rather than just a string. [[GH-13469](https://github.com/hashicorp/vault/pull/13469)]
* secrets/gcp: Fixed bug where error was not reported for invalid bindings [[GH-13974](https://github.com/hashicorp/vault/pull/13974)]
* secrets/gcp: Fixes role bindings for BigQuery dataset resources. [[GH-13548](https://github.com/hashicorp/vault/pull/13548)]
* secrets/openldap: Fix panic from nil logger in backend [[GH-14171](https://github.com/hashicorp/vault/pull/14171)]
* secrets/pki: Default value for key_bits changed to 0, enabling key_type=ec key generation with default value [[GH-13080](https://github.com/hashicorp/vault/pull/13080)]
* secrets/pki: Fix issuance of wildcard certificates matching glob patterns [[GH-14235](https://github.com/hashicorp/vault/pull/14235)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-13759](https://github.com/hashicorp/vault/pull/13759)]
* secrets/pki: Fix regression causing performance secondaries to forward certificate generation to the primary. [[GH-2456](https://github.com/hashicorp/vault/pull/2456)]
* secrets/pki: Fixes around NIST P-curve signature hash length, default value for signature_bits changed to 0. [[GH-12872](https://github.com/hashicorp/vault/pull/12872)]
* secrets/pki: Recognize ed25519 when requesting a response in PKCS8 format [[GH-13257](https://github.com/hashicorp/vault/pull/13257)]
* secrets/pki: Skip signature bits validation for ed25519 curve key type [[GH-13254](https://github.com/hashicorp/vault/pull/13254)]
* secrets/transit: Ensure that Vault does not panic for invalid nonce size when we aren't in convergent encryption mode. [[GH-13690](https://github.com/hashicorp/vault/pull/13690)]
* secrets/transit: Return an error if any required parameter is missing. [[GH-14074](https://github.com/hashicorp/vault/pull/14074)]
* storage/raft: Fix a panic when trying to store a key > 32KB in a transaction. [[GH-13286](https://github.com/hashicorp/vault/pull/13286)]
* storage/raft: Fix a panic when trying to write a key > 32KB [[GH-13282](https://github.com/hashicorp/vault/pull/13282)]
* storage/raft: Fix issues allowing invalid nodes to become leadership candidates. [[GH-13703](https://github.com/hashicorp/vault/pull/13703)]
* storage/raft: Fix regression in 1.9.0-rc1 that changed how time is represented in Raft logs; this prevented using a raft db created pre-1.9. [[GH-13165](https://github.com/hashicorp/vault/pull/13165)]
* storage/raft: On linux, use map_populate for bolt files to improve startup time. [[GH-13573](https://github.com/hashicorp/vault/pull/13573)]
* storage/raft: Units for bolt metrics now given in milliseconds instead of nanoseconds [[GH-13749](https://github.com/hashicorp/vault/pull/13749)]
* ui: Adds pagination to auth methods list view [[GH-13054](https://github.com/hashicorp/vault/pull/13054)]
* ui: Do not show verify connection value on database connection config page [[GH-13152](https://github.com/hashicorp/vault/pull/13152)]
* ui: Fix client count current month data not showing unless monthly history data exists [[GH-13396](https://github.com/hashicorp/vault/pull/13396)]
* ui: Fix default TTL display and set on database role [[GH-14224](https://github.com/hashicorp/vault/pull/14224)]
* ui: Fix incorrect validity message on transit secrets engine [[GH-14233](https://github.com/hashicorp/vault/pull/14233)]
* ui: Fix issue where UI incorrectly handled API errors when mounting backends [[GH-14551](https://github.com/hashicorp/vault/pull/14551)]
* ui: Fix kv engine access bug [[GH-13872](https://github.com/hashicorp/vault/pull/13872)]
* ui: Fixes breadcrumb bug for secrets navigation [[GH-13604](https://github.com/hashicorp/vault/pull/13604)]
* ui: Fixes caching issue on kv new version create [[GH-14489](https://github.com/hashicorp/vault/pull/14489)]
* ui: Fixes displaying empty masked values in PKI engine [[GH-14400](https://github.com/hashicorp/vault/pull/14400)]
* ui: Fixes horizontal bar chart hover issue when filtering namespaces and mounts [[GH-14493](https://github.com/hashicorp/vault/pull/14493)]
* ui: Fixes issue logging out with wrapped token query parameter [[GH-14329](https://github.com/hashicorp/vault/pull/14329)]
* ui: Fixes issue removing raft storage peer via cli not reflected in UI until refresh [[GH-13098](https://github.com/hashicorp/vault/pull/13098)]
* ui: Fixes issue restoring raft storage snapshot [[GH-13107](https://github.com/hashicorp/vault/pull/13107)]
* ui: Fixes issue saving KMIP role correctly [[GH-13585](https://github.com/hashicorp/vault/pull/13585)]
* ui: Fixes issue with OIDC auth workflow when using MetaMask Chrome extension [[GH-13133](https://github.com/hashicorp/vault/pull/13133)]
* ui: Fixes issue with SearchSelect component not holding focus [[GH-13590](https://github.com/hashicorp/vault/pull/13590)]
* ui: Fixes issue with automate secret deletion value not displaying initially if set in secret metadata edit view [[GH-13177](https://github.com/hashicorp/vault/pull/13177)]
* ui: Fixes issue with correct auth method not selected when logging out from OIDC or JWT methods [[GH-14545](https://github.com/hashicorp/vault/pull/14545)]
* ui: Fixes issue with placeholder not displaying for automatically deleted secrets when deletion time has passed [[GH-13166](https://github.com/hashicorp/vault/pull/13166)]
* ui: Fixes issue with the number of PGP Key inputs not matching the key shares number in the initialization form on change [[GH-13038](https://github.com/hashicorp/vault/pull/13038)]
* ui: Fixes long secret key names overlapping masked values [[GH-13032](https://github.com/hashicorp/vault/pull/13032)]
* ui: Fixes node-forge error when parsing EC (elliptical curve) certs [[GH-13238](https://github.com/hashicorp/vault/pull/13238)]
* ui: Redirects to managed namespace if incorrect namespace in URL param [[GH-14422](https://github.com/hashicorp/vault/pull/14422)]
* ui: Removes ability to tune token_type for token auth methods [[GH-12904](https://github.com/hashicorp/vault/pull/12904)]
* ui: trigger token renewal if inactive and half of TTL has passed [[GH-13950](https://github.com/hashicorp/vault/pull/13950)]

## 1.9.9
### August 31, 2022

CHANGES:

* core: Bump Go version to 1.17.13.

BUG FIXES:

* core (enterprise): Fix some races in merkle index flushing code found in testing
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]

## 1.9.8
### July 21, 2022

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
