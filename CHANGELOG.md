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
