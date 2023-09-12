## Previous versions
- [v1.0.0 - v1.9.10](CHANGELOG-pre-v1.10.md)
- [v0.11.6 and earlier](CHANGELOG-v0.md)

## 1.14.2
### August 30, 2023

CHANGES:

* auth/azure: Update plugin to v0.16.0 [[GH-22277](https://github.com/hashicorp/vault/pull/22277)]
* core: Bump Go version to 1.20.7.
* database/snowflake: Update plugin to v0.9.0 [[GH-22516](https://github.com/hashicorp/vault/pull/22516)]

IMPROVEMENTS:

* auto-auth/azure: Added Azure Workload Identity Federation support to auto-auth (for Vault Agent and Vault Proxy). [[GH-22264](https://github.com/hashicorp/vault/pull/22264)]
* core: Log rollback manager failures during unmount, remount to prevent replication failures on secondary clusters. [[GH-22235](https://github.com/hashicorp/vault/pull/22235)]
* kmip (enterprise): Add namespace lock and unlock support [[GH-21925](https://github.com/hashicorp/vault/pull/21925)]
* replication (enterprise): Make reindex less disruptive by allowing writes during the flush phase.
* secrets/database: Improves error logging for static role rotations by including the database and role names. [[GH-22253](https://github.com/hashicorp/vault/pull/22253)]
* storage/raft: Cap the minimum dead_server_last_contact_threshold to 1m. [[GH-22040](https://github.com/hashicorp/vault/pull/22040)]
* ui: KV View Secret card will link to list view if input ends in "/" [[GH-22502](https://github.com/hashicorp/vault/pull/22502)]
* ui: adds allowed_user_ids field to create role form and user_ids to generate certificates form in pki [[GH-22191](https://github.com/hashicorp/vault/pull/22191)]
* ui: enables create and update KV secret workflow when control group present [[GH-22471](https://github.com/hashicorp/vault/pull/22471)]
* website/docs: Fix link formatting in Vault lambda extension docs [[GH-22396](https://github.com/hashicorp/vault/pull/22396)]

BUG FIXES:

* activity (enterprise): Fix misattribution of entities to no or child namespace auth methods [[GH-18809](https://github.com/hashicorp/vault/pull/18809)]
* agent: Environment variable VAULT_CACERT_BYTES now works for Vault Agent templates. [[GH-22322](https://github.com/hashicorp/vault/pull/22322)]
* api: Fix breakage with UNIX domain socket addresses introduced by newest Go versions as a security fix. [[GH-22523](https://github.com/hashicorp/vault/pull/22523)]
* core (enterprise): Remove MFA Configuration for namespace when deleting namespace
* core/metrics: vault.raft_storage.bolt.write.time should be a counter not a summary [[GH-22468](https://github.com/hashicorp/vault/pull/22468)]
* core/quotas (enterprise): Fix a case where we were applying login roles to lease count quotas in a non-login context.
Also fix a related potential deadlock. [[GH-21110](https://github.com/hashicorp/vault/pull/21110)]
* core:  Remove "expiration manager is nil on tokenstore" error log for unauth requests on DR secondary as they do not have expiration manager. [[GH-22137](https://github.com/hashicorp/vault/pull/22137)]
* core: Fix bug where background thread to update locked user entries runs on DR secondaries. [[GH-22355](https://github.com/hashicorp/vault/pull/22355)]
* core: Fix readonly errors that could occur while loading mounts/auths during unseal [[GH-22362](https://github.com/hashicorp/vault/pull/22362)]
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-21470](https://github.com/hashicorp/vault/pull/21470)]
* expiration: Fix a deadlock that could occur when a revocation failure happens while restoring leases on startup. [[GH-22374](https://github.com/hashicorp/vault/pull/22374)]
* license: Add autoloaded license path to the cache exempt list. This is to ensure the license changes on the active node is observed on the perfStandby node. [[GH-22363](https://github.com/hashicorp/vault/pull/22363)]
* replication (enterprise): Fix bug sync invalidate CoreReplicatedClusterInfoPath
* replication (enterprise): Fix panic when update-primary was called on demoted clusters using update_primary_addrs
* replication (enterprise): Fixing a bug by which the atomicity of a merkle diff result could be affected. This means it could be a source of a merkle-diff & sync process failing to switch into stream-wal mode afterwards.
* sdk/ldaputil: Properly escape user filters when using UPN domains
sdk/ldaputil: use EscapeLDAPValue implementation from cap/ldap [[GH-22249](https://github.com/hashicorp/vault/pull/22249)]
* secrets/ldap: Fix bug causing schema and password_policy to be overwritten in config. [[GH-22330](https://github.com/hashicorp/vault/pull/22330)]
* secrets/transform (enterprise): Batch items with repeated tokens in the tokenization decode api will now contain the decoded_value element
* secrets/transform (enterprise): Fix nil panic when encoding a tokenization transformation on a non-active node
* secrets/transform (enterprise): Tidy operations will be re-scheduled at a minimum of every minute, not a maximum of every minute
* storage/raft: Fix race where new follower joining can get pruned by dead server cleanup. [[GH-20986](https://github.com/hashicorp/vault/pull/20986)]
* ui: Fix blank page or ghost secret when canceling KV secret create [[GH-22541](https://github.com/hashicorp/vault/pull/22541)]
* ui: fixes `max_versions` default for secret metadata unintentionally overriding kv engine defaults [[GH-22394](https://github.com/hashicorp/vault/pull/22394)]
* ui: fixes model defaults overwriting input value when user tries to clear form input [[GH-22458](https://github.com/hashicorp/vault/pull/22458)]
* ui: fixes text readability issue in revoke token confirmation dialog [[GH-22390](https://github.com/hashicorp/vault/pull/22390)]

## 1.14.1
### July 25, 2023

CHANGES:

* auth/ldap: Normalize HTTP response codes when invalid credentials are provided [[GH-21282](https://github.com/hashicorp/vault/pull/21282)]
* core/namespace (enterprise): Introduce the concept of high-privilege namespace (administrative namespace),
which will have access to some system backend paths that were previously only accessible in the root namespace. [[GH-21215](https://github.com/hashicorp/vault/pull/21215)]
* secrets/transform (enterprise): Enforce a transformation role's max_ttl setting on encode requests, a warning will be returned if max_ttl was applied.
* storage/aerospike: Aerospike storage shouldn't be used on 32-bit architectures and is now unsupported on them. [[GH-20825](https://github.com/hashicorp/vault/pull/20825)]

IMPROVEMENTS:

* core/fips: Add RPM, DEB packages of FIPS 140-2 and HSM+FIPS 140-2 Vault Enterprise.
* eventbus: updated go-eventlogger library to allow removal of nodes referenced by pipelines (used for subscriptions) [[GH-21623](https://github.com/hashicorp/vault/pull/21623)]
* openapi: Better mount points for kv-v1 and kv-v2 in openapi.json [[GH-21563](https://github.com/hashicorp/vault/pull/21563)]
* replication (enterprise): Avoid logging warning if request is forwarded from a performance standby and not a performance secondary
* secrets/pki: Add a parameter to allow ExtKeyUsage field usage from a role within ACME. [[GH-21702](https://github.com/hashicorp/vault/pull/21702)]
* secrets/transform (enterprise): Switch to pgx PostgreSQL driver for better timeout handling
* sys/metrics (enterprise): Adds a gauge metric that tracks whether enterprise builtin secret plugins are enabled. [[GH-21681](https://github.com/hashicorp/vault/pull/21681)]

BUG FIXES:

* agent: Fix "generate-config" command documentation URL [[GH-21466](https://github.com/hashicorp/vault/pull/21466)]
* auth/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21800](https://github.com/hashicorp/vault/pull/21800)]
* auth/token, sys: Fix path-help being unavailable for some list-only endpoints [[GH-18571](https://github.com/hashicorp/vault/pull/18571)]
* auth/token: Fix parsing of `auth/token/create` fields to avoid incorrect warnings about ignored parameters [[GH-18556](https://github.com/hashicorp/vault/pull/18556)]
* awsutil: Update awsutil to v0.2.3 to fix a regression where Vault no longer
respects `AWS_ROLE_ARN`, `AWS_WEB_IDENTITY_TOKEN_FILE`, and `AWS_ROLE_SESSION_NAME`. [[GH-21951](https://github.com/hashicorp/vault/pull/21951)]
* core/managed-keys (enterprise): Allow certain symmetric PKCS#11 managed key mechanisms (AES CBC with and without padding) to operate without an HMAC.
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-24170](https://github.com/hashicorp/vault/pull/24170)]
* core: Fixed issue with some durations not being properly parsed to include days. [[GH-21357](https://github.com/hashicorp/vault/pull/21357)]
* identity: Remove caseSensitivityKey to prevent errors while loading groups which could result in missing groups in memDB when duplicates are found. [[GH-20965](https://github.com/hashicorp/vault/pull/20965)]
* openapi: Fix response schema for PKI Issue requests [[GH-21449](https://github.com/hashicorp/vault/pull/21449)]
* openapi: Fix schema definitions for PKI EAB APIs [[GH-21458](https://github.com/hashicorp/vault/pull/21458)]
* replication (enterprise): update primary cluster address after DR failover
* secrets/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21631](https://github.com/hashicorp/vault/pull/21631)]
* secrets/pki: Fix bug with ACME tidy, 'unable to determine acme base folder path'. [[GH-21870](https://github.com/hashicorp/vault/pull/21870)]
* secrets/pki: Fix preserving acme_account_safety_buffer on config/auto-tidy. [[GH-21870](https://github.com/hashicorp/vault/pull/21870)]
* secrets/pki: Prevent deleted issuers from reappearing when migrating from a version 1 bundle to a version 2 bundle (versions including 1.13.0, 1.12.2, and 1.11.6); when managed keys were removed but referenced in the Vault 1.10 legacy CA bundle, this the error: `no managed key found with uuid`. [[GH-21316](https://github.com/hashicorp/vault/pull/21316)]
* secrets/transform (enterprise): Fix nil panic when deleting a template with tokenization transformations present
* secrets/transform (enterprise): Grab shared locks for various read operations, only escalating to write locks if work is required
* serviceregistration: Fix bug where multiple nodes in a secondary cluster could be labelled active after updating the cluster's primary [[GH-21642](https://github.com/hashicorp/vault/pull/21642)]
* ui: Adds missing values to details view after generating PKI certificate [[GH-21635](https://github.com/hashicorp/vault/pull/21635)]
* ui: Fixed an issue where editing an SSH role would clear `default_critical_options` and `default_extension` if left unchanged. [[GH-21739](https://github.com/hashicorp/vault/pull/21739)]
* ui: Fixed secrets, leases, and policies filter dropping focus after a single character [[GH-21767](https://github.com/hashicorp/vault/pull/21767)]
* ui: Fixes issue with certain navigational links incorrectly displaying in child namespaces [[GH-21562](https://github.com/hashicorp/vault/pull/21562)]
* ui: Fixes login screen display issue with Safari browser [[GH-21582](https://github.com/hashicorp/vault/pull/21582)]
* ui: Fixes problem displaying certificates issued with unsupported signature algorithms (i.e. ed25519) [[GH-21926](https://github.com/hashicorp/vault/pull/21926)]
* ui: Fixes styling of private key input when configuring an SSH key [[GH-21531](https://github.com/hashicorp/vault/pull/21531)]
* ui: Surface DOMException error when browser settings prevent localStorage. [[GH-21503](https://github.com/hashicorp/vault/pull/21503)]

## 1.14.0
### June 21, 2023

BREAKING CHANGES:

* secrets/pki: Maintaining running count of certificates will be turned off by default.
To re-enable keeping these metrics available on the tidy status endpoint, enable
maintain_stored_certificate_counts on tidy-config, to also publish them to the
metrics consumer, enable publish_stored_certificate_count_metrics . [[GH-18186](https://github.com/hashicorp/vault/pull/18186)]

CHANGES:

* auth/alicloud: Updated plugin from v0.14.0 to v0.15.0 [[GH-20758](https://github.com/hashicorp/vault/pull/20758)]
* auth/azure: Updated plugin from v0.13.0 to v0.15.0 [[GH-20816](https://github.com/hashicorp/vault/pull/20816)]
* auth/centrify: Updated plugin from v0.14.0 to v0.15.1 [[GH-20745](https://github.com/hashicorp/vault/pull/20745)]
* auth/gcp: Updated plugin from v0.15.0 to v0.16.0 [[GH-20725](https://github.com/hashicorp/vault/pull/20725)]
* auth/jwt: Updated plugin from v0.15.0 to v0.16.0 [[GH-20799](https://github.com/hashicorp/vault/pull/20799)]
* auth/kubernetes: Update plugin to v0.16.0 [[GH-20802](https://github.com/hashicorp/vault/pull/20802)]
* core: Bump Go version to 1.20.5.
* core: Remove feature toggle for SSCTs, i.e. the env var VAULT_DISABLE_SERVER_SIDE_CONSISTENT_TOKENS. [[GH-20834](https://github.com/hashicorp/vault/pull/20834)]
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]
* database/couchbase: Updated plugin from v0.9.0 to v0.9.2 [[GH-20764](https://github.com/hashicorp/vault/pull/20764)]
* database/redis-elasticache: Updated plugin from v0.2.0 to v0.2.1 [[GH-20751](https://github.com/hashicorp/vault/pull/20751)]
* replication (enterprise): Add a new parameter for the update-primary API call
that allows for setting of the primary cluster addresses directly, instead of
via a token.
* secrets/ad: Updated plugin from v0.10.1-0.20230329210417-0b2cdb26cf5d to v0.16.0 [[GH-20750](https://github.com/hashicorp/vault/pull/20750)]
* secrets/alicloud: Updated plugin from v0.5.4-beta1.0.20230330124709-3fcfc5914a22 to v0.15.0 [[GH-20787](https://github.com/hashicorp/vault/pull/20787)]
* secrets/aure: Updated plugin from v0.15.0 to v0.16.0 [[GH-20777](https://github.com/hashicorp/vault/pull/20777)]
* secrets/database/mongodbatlas: Updated plugin from v0.9.0 to v0.10.0 [[GH-20882](https://github.com/hashicorp/vault/pull/20882)]
* secrets/database/snowflake: Updated plugin from v0.7.0 to v0.8.0 [[GH-20807](https://github.com/hashicorp/vault/pull/20807)]
* secrets/gcp: Updated plugin from v0.15.0 to v0.16.0 [[GH-20818](https://github.com/hashicorp/vault/pull/20818)]
* secrets/keymgmt: Updated plugin to v0.9.1
* secrets/kubernetes: Update plugin to v0.5.0 [[GH-20802](https://github.com/hashicorp/vault/pull/20802)]
* secrets/mongodbatlas: Updated plugin from v0.9.1 to v0.10.0 [[GH-20742](https://github.com/hashicorp/vault/pull/20742)]
* secrets/pki: Allow issuance of root CAs without AIA, when templated AIA information includes issuer_id. [[GH-21209](https://github.com/hashicorp/vault/pull/21209)]
* secrets/pki: Warning when issuing leafs from CSRs with basic constraints. In the future, issuance of non-CA leaf certs from CSRs with asserted IsCA Basic Constraints will be prohibited. [[GH-20654](https://github.com/hashicorp/vault/pull/20654)]

FEATURES:

* **AWS Static Roles**: The AWS Secrets Engine can manage static roles configured by users. [[GH-20536](https://github.com/hashicorp/vault/pull/20536)]
* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* **Environment Variables through Vault Agent**: Introducing a new process-supervisor mode for Vault Agent which allows injecting secrets as environment variables into a child process using a new `env_template` configuration stanza. The process-supervisor configuration can be generated with a new `vault agent generate-config` helper tool. [[GH-20530](https://github.com/hashicorp/vault/pull/20530)]
* **MongoDB Atlas Database Secrets**: Adds support for client certificate credentials [[GH-20425](https://github.com/hashicorp/vault/pull/20425)]
* **MongoDB Atlas Database Secrets**: Adds support for generating X.509 certificates on dynamic roles for user authentication [[GH-20882](https://github.com/hashicorp/vault/pull/20882)]
* **NEW PKI Workflow in UI**: Completes generally available rollout of new PKI UI that provides smoother mount configuration and a more guided user experience [[GH-pki-ui-improvements](https://github.com/hashicorp/vault/pull/pki-ui-improvements)]
* **Secrets/Auth Plugin Multiplexing**: The plugin will be multiplexed when run
as an external plugin by vault versions that support secrets/auth plugin
multiplexing (> 1.12) [[GH-19215](https://github.com/hashicorp/vault/pull/19215)]
* **Sidebar Navigation in UI**: A new sidebar navigation panel has been added in the UI to replace the top navigation bar. [[GH-19296](https://github.com/hashicorp/vault/pull/19296)]
* **Vault PKI ACME Server**: Support for the ACME certificate lifecycle management protocol has been added to the Vault PKI Plugin. This allows standard ACME clients, such as the EFF's certbot and the CNCF's k8s cert-manager, to request certificates from a Vault server with no knowledge of Vault APIs or authentication mechanisms. For public-facing Vault instances, we recommend requiring External Account Bindings (EAB) to limit the ability to request certificates to only authenticated clients. [[GH-20752](https://github.com/hashicorp/vault/pull/20752)]
* **Vault Proxy**: Introduced Vault Proxy, a new subcommand of the Vault binary that can be invoked using `vault proxy -config=config.hcl`. It currently has the same feature set as Vault Agent's API proxy, but the two may diverge in the future. We plan to deprecate the API proxy functionality of Vault Agent in a future release. [[GH-20548](https://github.com/hashicorp/vault/pull/20548)]
* **OCI Auto-Auth**: Add OCI (Oracle Cloud Infrastructure) auto-auth method [[GH-19260](https://github.com/hashicorp/vault/pull/19260)]

IMPROVEMENTS:

* * api: Add Config.TLSConfig method to fetch the TLS configuration from a client config. [[GH-20265](https://github.com/hashicorp/vault/pull/20265)]
* * physical/etcd: Upgrade etcd3 client to v3.5.7 [[GH-20261](https://github.com/hashicorp/vault/pull/20261)]
* activitylog: EntityRecord protobufs now contain a ClientType field for
distinguishing client sources. [[GH-20626](https://github.com/hashicorp/vault/pull/20626)]
* agent: Add integration tests for agent running in process supervisor mode [[GH-20741](https://github.com/hashicorp/vault/pull/20741)]
* agent: Add logic to validate env_template entries in configuration [[GH-20569](https://github.com/hashicorp/vault/pull/20569)]
* agent: Added `reload` option to cert auth configuration in case of external renewals of local x509 key-pairs. [[GH-19002](https://github.com/hashicorp/vault/pull/19002)]
* agent: JWT auto-auth has a new config option, `remove_jwt_follows_symlinks` (default: false), that, if set to true will now remove the JWT, instead of the symlink to the JWT, if a symlink to a JWT has been provided in the `path` option, and the `remove_jwt_after_reading` config option is set to true (default). [[GH-18863](https://github.com/hashicorp/vault/pull/18863)]
* agent: Vault Agent now reports its name and version as part of the User-Agent header in all requests issued. [[GH-19776](https://github.com/hashicorp/vault/pull/19776)]
* agent: initial implementation of a process runner for injecting secrets via environment variables via vault agent [[GH-20628](https://github.com/hashicorp/vault/pull/20628)]
* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* api: property based testing for LifetimeWatcher sleep duration calculation [[GH-17919](https://github.com/hashicorp/vault/pull/17919)]
* audit: add plugin metadata, including plugin name, type, version, sha256, and whether plugin is external, to audit logging [[GH-19814](https://github.com/hashicorp/vault/pull/19814)]
* audit: forwarded requests can now contain host metadata on the node it was sent 'from' or a flag to indicate that it was forwarded.
* auth/cert: Better return OCSP validation errors during login to the caller. [[GH-20234](https://github.com/hashicorp/vault/pull/20234)]
* auth/kerberos: Enable plugin multiplexing
auth/kerberos: Upgrade plugin dependencies [[GH-20771](https://github.com/hashicorp/vault/pull/20771)]
* auth/ldap: allow configuration of alias dereferencing in LDAP search [[GH-18230](https://github.com/hashicorp/vault/pull/18230)]
* auth/ldap: allow providing the LDAP password via an env var when authenticating via the CLI [[GH-18225](https://github.com/hashicorp/vault/pull/18225)]
* auth/oidc: Adds support for group membership parsing when using IBM ISAM as an OIDC provider. [[GH-19247](https://github.com/hashicorp/vault/pull/19247)]
* build: Prefer GOBIN when set over GOPATH/bin when building the binary [[GH-19862](https://github.com/hashicorp/vault/pull/19862)]
* cli: Add walkSecretsTree helper function, which recursively walks secrets rooted at the given path [[GH-20464](https://github.com/hashicorp/vault/pull/20464)]
* cli: Improve addPrefixToKVPath helper [[GH-20488](https://github.com/hashicorp/vault/pull/20488)]
* command/server (enterprise): -dev-three-node now creates perf standbys instead of regular standbys. [[GH-20629](https://github.com/hashicorp/vault/pull/20629)]
* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* command/server: New -dev-cluster-json writes a file describing the dev cluster in -dev and -dev-three-node modes, plus -dev-three-node now enables unauthenticated metrics and pprof requests. [[GH-20224](https://github.com/hashicorp/vault/pull/20224)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core, secrets/pki, audit: Update dependency go-jose to v3 due to v2 deprecation. [[GH-20559](https://github.com/hashicorp/vault/pull/20559)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* core: Add possibility to decode a generated encoded root token via the rest API [[GH-20595](https://github.com/hashicorp/vault/pull/20595)]
* core: include namespace path in granting_policies block of audit log
* core: include reason for ErrReadOnly on PBPWF writing failures
* core: report intermediate error messages during request forwarding [[GH-20643](https://github.com/hashicorp/vault/pull/20643)]
* core:provide more descriptive error message when calling enterprise feature paths in open-source [[GH-18870](https://github.com/hashicorp/vault/pull/18870)]
* database/elasticsearch: Upgrade plugin dependencies [[GH-20767](https://github.com/hashicorp/vault/pull/20767)]
* database/mongodb: upgrade mongo driver to 1.11 [[GH-19954](https://github.com/hashicorp/vault/pull/19954)]
* database/redis: Upgrade plugin dependencies [[GH-20763](https://github.com/hashicorp/vault/pull/20763)]
* http: Support responding to HEAD operation from plugins [[GH-19520](https://github.com/hashicorp/vault/pull/19520)]
* openapi: Add openapi response definitions to /sys defined endpoints. [[GH-18633](https://github.com/hashicorp/vault/pull/18633)]
* openapi: Add openapi response definitions to pki/config_*.go [[GH-18376](https://github.com/hashicorp/vault/pull/18376)]
* openapi: Add openapi response definitions to vault/logical_system_paths.go defined endpoints. [[GH-18515](https://github.com/hashicorp/vault/pull/18515)]
* openapi: Consistently stop Vault server on exit in gen_openapi.sh [[GH-19252](https://github.com/hashicorp/vault/pull/19252)]
* openapi: Improve operationId/request/response naming strategy [[GH-19319](https://github.com/hashicorp/vault/pull/19319)]
* openapi: add openapi response definitions to /sys/internal endpoints [[GH-18542](https://github.com/hashicorp/vault/pull/18542)]
* openapi: add openapi response definitions to /sys/rotate endpoints [[GH-18624](https://github.com/hashicorp/vault/pull/18624)]
* openapi: add openapi response definitions to /sys/seal endpoints [[GH-18625](https://github.com/hashicorp/vault/pull/18625)]
* openapi: add openapi response definitions to /sys/tool endpoints [[GH-18626](https://github.com/hashicorp/vault/pull/18626)]
* openapi: add openapi response definitions to /sys/version-history, /sys/leader, /sys/ha-status, /sys/host-info, /sys/in-flight-req [[GH-18628](https://github.com/hashicorp/vault/pull/18628)]
* openapi: add openapi response definitions to /sys/wrapping endpoints [[GH-18627](https://github.com/hashicorp/vault/pull/18627)]
* openapi: add openapi response defintions to /sys/auth endpoints [[GH-18465](https://github.com/hashicorp/vault/pull/18465)]
* openapi: add openapi response defintions to /sys/capabilities endpoints [[GH-18468](https://github.com/hashicorp/vault/pull/18468)]
* openapi: add openapi response defintions to /sys/config and /sys/generate-root endpoints [[GH-18472](https://github.com/hashicorp/vault/pull/18472)]
* openapi: added ability to validate response structures against openapi schema for test clusters [[GH-19043](https://github.com/hashicorp/vault/pull/19043)]
* sdk/framework: Fix non-deterministic ordering of 'required' fields in OpenAPI spec [[GH-20881](https://github.com/hashicorp/vault/pull/20881)]
* sdk: Add new docker-based cluster testing framework to the sdk. [[GH-20247](https://github.com/hashicorp/vault/pull/20247)]
* secrets/ad: upgrades dependencies [[GH-19829](https://github.com/hashicorp/vault/pull/19829)]
* secrets/alicloud: upgrades dependencies [[GH-19846](https://github.com/hashicorp/vault/pull/19846)]
* secrets/consul: Improve error message when ACL bootstrapping fails. [[GH-20891](https://github.com/hashicorp/vault/pull/20891)]
* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]
* secrets/gcpkms: Enable plugin multiplexing
secrets/gcpkms: Upgrade plugin dependencies [[GH-20784](https://github.com/hashicorp/vault/pull/20784)]
* secrets/mongodbatlas: upgrades dependencies [[GH-19861](https://github.com/hashicorp/vault/pull/19861)]
* secrets/openldap: upgrades dependencies [[GH-19993](https://github.com/hashicorp/vault/pull/19993)]
* secrets/pki: Add missing fields to tidy-status, include new last_auto_tidy_finished field. [[GH-20442](https://github.com/hashicorp/vault/pull/20442)]
* secrets/pki: Add warning when issuer lacks KeyUsage during CRL rebuilds; expose in logs and on rotation. [[GH-20253](https://github.com/hashicorp/vault/pull/20253)]
* secrets/pki: Allow determining existing issuers and keys on import. [[GH-20441](https://github.com/hashicorp/vault/pull/20441)]
* secrets/pki: Include CA serial number, key UUID on issuers list endpoint. [[GH-20276](https://github.com/hashicorp/vault/pull/20276)]
* secrets/pki: Limit ACME issued certificates NotAfter TTL to a maximum of 90 days [[GH-20981](https://github.com/hashicorp/vault/pull/20981)]
* secrets/pki: Support TLS-ALPN-01 challenge type in ACME for DNS certificate identifiers. [[GH-20943](https://github.com/hashicorp/vault/pull/20943)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]
* secrets/postgresql: Add configuration to scram-sha-256 encrypt passwords on Vault before sending them to PostgreSQL [[GH-19616](https://github.com/hashicorp/vault/pull/19616)]
* secrets/terraform: upgrades dependencies [[GH-19798](https://github.com/hashicorp/vault/pull/19798)]
* secrets/transit: Add support to import public keys in transit engine and allow encryption and verification of signed data [[GH-17934](https://github.com/hashicorp/vault/pull/17934)]
* secrets/transit: Allow importing RSA-PSS OID (1.2.840.113549.1.1.10) private keys via BYOK. [[GH-19519](https://github.com/hashicorp/vault/pull/19519)]
* secrets/transit: Respond to writes with updated key policy, cache configuration. [[GH-20652](https://github.com/hashicorp/vault/pull/20652)]
* secrets/transit: Support BYOK-encrypted export of keys to securely allow synchronizing specific keys and version across clusters. [[GH-20736](https://github.com/hashicorp/vault/pull/20736)]
* ui: Add download button for each secret value in KV v2 [[GH-20431](https://github.com/hashicorp/vault/pull/20431)]
* ui: Add filtering by auth type and auth name to the Authentication Method list view. [[GH-20747](https://github.com/hashicorp/vault/pull/20747)]
* ui: Add filtering by engine type and engine name to the Secret Engine list view. [[GH-20481](https://github.com/hashicorp/vault/pull/20481)]
* ui: Adds whitespace warning to secrets engine and auth method path inputs [[GH-19913](https://github.com/hashicorp/vault/pull/19913)]
* ui: Remove the Bulma CSS framework. [[GH-19878](https://github.com/hashicorp/vault/pull/19878)]
* ui: Update Web CLI with examples and a new `kv-get` command for reading kv v2 data and metadata [[GH-20590](https://github.com/hashicorp/vault/pull/20590)]
* ui: Updates UI javascript dependencies [[GH-19901](https://github.com/hashicorp/vault/pull/19901)]
* ui: add allowed_managed_keys field to secret engine mount options [[GH-19791](https://github.com/hashicorp/vault/pull/19791)]
* ui: adds warning for commas in stringArray inputs and updates tooltip help text to remove references to comma separation [[GH-20163](https://github.com/hashicorp/vault/pull/20163)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]
* website/docs: Add rotate root documentation for azure secrets engine [[GH-19187](https://github.com/hashicorp/vault/pull/19187)]
* website/docs: fix database static-user sample payload [[GH-19170](https://github.com/hashicorp/vault/pull/19170)]

BUG FIXES:

* agent: Fix agent generate-config to accept -namespace, VAULT_NAMESPACE, and other client-modifying flags. [[GH-21297](https://github.com/hashicorp/vault/pull/21297)]
* agent: Fix bug with 'cache' stanza validation [[GH-20934](https://github.com/hashicorp/vault/pull/20934)]
* api: Addressed a couple of issues that arose as edge cases for the -output-policy flag. Specifically around properly handling list commands, distinguishing kv V1/V2, and correctly recognizing protected paths. [[GH-19160](https://github.com/hashicorp/vault/pull/19160)]
* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* auth/token: Fix cubbyhole and revocation for legacy service tokens [[GH-19416](https://github.com/hashicorp/vault/pull/19416)]
* cli/kv: add -mount flag to kv list [[GH-19378](https://github.com/hashicorp/vault/pull/19378)]
* core (enterprise): Don't delete backend stored data that appears to be filterable
on this secondary if we don't have a corresponding mount entry.
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry.  Fix wrong lock used in ListAuths link meta interface implementation. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix Forwarded Writer construction to correctly find active nodes, allowing PKI cross-cluster functionality to succeed on existing mounts.
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* core: Fix writes to readonly storage on performance standbys when user lockout feature is enabled. [[GH-20783](https://github.com/hashicorp/vault/pull/20783)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* license (enterprise): Fix bug where license would update even if the license didn't change.
* openapi: Small fixes for OpenAPI display attributes. Changed "log-in" to "login" [[GH-20285](https://github.com/hashicorp/vault/pull/20285)]
* plugin/reload:  Fix a possible data race with rollback manager and plugin reload [[GH-19468](https://github.com/hashicorp/vault/pull/19468)]
* replication (enterprise): Fix a caching issue when replicating filtered data to
a performance secondary. This resulted in the data being set to nil in the cache
and a "invalid value" error being returned from the API.
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs
* replication (enterprise): Fix regression causing token creation against a role
with a new entity alias to be incorrectly forwarded from perf standbys. [[GH-21100](https://github.com/hashicorp/vault/pull/21100)]
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* sdk/backend: prevent panic when computing the zero value for a `TypeInt64` schema field. [[GH-18729](https://github.com/hashicorp/vault/pull/18729)]
* secrets/pki: Support setting both maintain_stored_certificate_counts=false and publish_stored_certificate_count_metrics=false explicitly in tidy config. [[GH-20664](https://github.com/hashicorp/vault/pull/20664)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* secrets/transform: Added importing of keys and key versions into the Transform secrets engine using the command 'vault transform import' and 'vault transform import-version'. [[GH-20668](https://github.com/hashicorp/vault/pull/20668)]
* secrets/transit: Fix export of HMAC-only key, correctly exporting the key used for sign operations. For consumers of the previously incorrect key, use the plaintext export to retrieve these incorrect keys and import them as new versions.
* secrets/transit: Fix bug related to shorter dedicated HMAC key sizing.
* sdk/helper/keysutil: New HMAC type policies will have HMACKey equal to Key and be copied over on import. [[GH-20864](https://github.com/hashicorp/vault/pull/20864)]
* shamir: change mul and div implementations to be constant-time [[GH-19495](https://github.com/hashicorp/vault/pull/19495)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix secret render when path includes %. Resolves #11616. [[GH-20430](https://github.com/hashicorp/vault/pull/20430)]
* ui: Fixes issue unsealing cluster for seal types other than shamir [[GH-20897](https://github.com/hashicorp/vault/pull/20897)]
* ui: fixes auto_rotate_period ttl input for transit keys [[GH-20731](https://github.com/hashicorp/vault/pull/20731)]
* ui: fixes bug in kmip role form that caused `operation_all` to persist after deselecting all operation checkboxes [[GH-19139](https://github.com/hashicorp/vault/pull/19139)]
* ui: fixes key_bits and signature_bits reverting to default values when editing a pki role [[GH-20907](https://github.com/hashicorp/vault/pull/20907)]
* ui: wait for wanted message event during OIDC callback instead of using the first message event [[GH-18521](https://github.com/hashicorp/vault/pull/18521)]

## 1.13.6
### August 30, 2023

CHANGES:

* core: Bump Go version to 1.20.7.

IMPROVEMENTS:

* core: Log rollback manager failures during unmount, remount to prevent replication failures on secondary clusters. [[GH-22235](https://github.com/hashicorp/vault/pull/22235)]
* replication (enterprise): Make reindex less disruptive by allowing writes during the flush phase.
* secrets/database: Improves error logging for static role rotations by including the database and role names. [[GH-22253](https://github.com/hashicorp/vault/pull/22253)]
* storage/raft: Cap the minimum dead_server_last_contact_threshold to 1m. [[GH-22040](https://github.com/hashicorp/vault/pull/22040)]
* ui: KV View Secret card will link to list view if input ends in "/" [[GH-22502](https://github.com/hashicorp/vault/pull/22502)]
* ui: enables create and update KV secret workflow when control group present [[GH-22471](https://github.com/hashicorp/vault/pull/22471)]

BUG FIXES:

* activity (enterprise): Fix misattribution of entities to no or child namespace auth methods [[GH-18809](https://github.com/hashicorp/vault/pull/18809)]
* api: Fix breakage with UNIX domain socket addresses introduced by newest Go versions as a security fix. [[GH-22523](https://github.com/hashicorp/vault/pull/22523)]
* core (enterprise): Remove MFA Configuration for namespace when deleting namespace
* core/quotas (enterprise): Fix a case where we were applying login roles to lease count quotas in a non-login context.
Also fix a related potential deadlock. [[GH-21110](https://github.com/hashicorp/vault/pull/21110)]
* core:  Remove "expiration manager is nil on tokenstore" error log for unauth requests on DR secondary as they do not have expiration manager. [[GH-22137](https://github.com/hashicorp/vault/pull/22137)]
* core: Fix bug where background thread to update locked user entries runs on DR secondaries. [[GH-22355](https://github.com/hashicorp/vault/pull/22355)]
* core: Fix readonly errors that could occur while loading mounts/auths during unseal [[GH-22362](https://github.com/hashicorp/vault/pull/22362)]
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-21470](https://github.com/hashicorp/vault/pull/21470)]
* expiration: Fix a deadlock that could occur when a revocation failure happens while restoring leases on startup. [[GH-22374](https://github.com/hashicorp/vault/pull/22374)]
* license: Add autoloaded license path to the cache exempt list. This is to ensure the license changes on the active node is observed on the perfStandby node. [[GH-22363](https://github.com/hashicorp/vault/pull/22363)]
* replication (enterprise): Fix bug sync invalidate CoreReplicatedClusterInfoPath
* replication (enterprise): Fix panic when update-primary was called on demoted clusters using update_primary_addrs
* replication (enterprise): Fixing a bug by which the atomicity of a merkle diff result could be affected. This means it could be a source of a merkle-diff & sync process failing to switch into stream-wal mode afterwards.
* sdk/ldaputil: Properly escape user filters when using UPN domains
sdk/ldaputil: use EscapeLDAPValue implementation from cap/ldap [[GH-22249](https://github.com/hashicorp/vault/pull/22249)]
* secrets/ldap: Fix bug causing schema and password_policy to be overwritten in config. [[GH-22331](https://github.com/hashicorp/vault/pull/22331)]
* secrets/transform (enterprise): Tidy operations will be re-scheduled at a minimum of every minute, not a maximum of every minute
* ui: Fix blank page or ghost secret when canceling KV secret create [[GH-22541](https://github.com/hashicorp/vault/pull/22541)]
* ui: fixes `max_versions` default for secret metadata unintentionally overriding kv engine defaults [[GH-22394](https://github.com/hashicorp/vault/pull/22394)]
* ui: fixes model defaults overwriting input value when user tries to clear form input [[GH-22458](https://github.com/hashicorp/vault/pull/22458)]
  
## 1.13.5
### July 25, 2023

CHANGES:

* auth/ldap: Normalize HTTP response codes when invalid credentials are provided [[GH-21282](https://github.com/hashicorp/vault/pull/21282)]
* core/namespace (enterprise): Introduce the concept of high-privilege namespace (administrative namespace),
which will have access to some system backend paths that were previously only accessible in the root namespace. [[GH-21215](https://github.com/hashicorp/vault/pull/21215)]
* secrets/transform (enterprise): Enforce a transformation role's max_ttl setting on encode requests, a warning will be returned if max_ttl was applied.

IMPROVEMENTS:

* core/fips: Add RPM, DEB packages of FIPS 140-2 and HSM+FIPS 140-2 Vault Enterprise.
* core: Add a new periodic metric to track the number of available policies, `vault.policy.configured.count`. [[GH-21010](https://github.com/hashicorp/vault/pull/21010)]
* replication (enterprise): Avoid logging warning if request is forwarded from a performance standby and not a performance secondary
* secrets/transform (enterprise): Switch to pgx PostgreSQL driver for better timeout handling
* sys/metrics (enterprise): Adds a gauge metric that tracks whether enterprise builtin secret plugins are enabled. [[GH-21681](https://github.com/hashicorp/vault/pull/21681)]

BUG FIXES:

* auth/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21799](https://github.com/hashicorp/vault/pull/21799)]
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-24170](https://github.com/hashicorp/vault/pull/24170)]
* identity: Remove caseSensitivityKey to prevent errors while loading groups which could result in missing groups in memDB when duplicates are found. [[GH-20965](https://github.com/hashicorp/vault/pull/20965)]
* replication (enterprise): update primary cluster address after DR failover
* secrets/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21632](https://github.com/hashicorp/vault/pull/21632)]
* secrets/pki: Prevent deleted issuers from reappearing when migrating from a version 1 bundle to a version 2 bundle (versions including 1.13.0, 1.12.2, and 1.11.6); when managed keys were removed but referenced in the Vault 1.10 legacy CA bundle, this the error: `no managed key found with uuid`. [[GH-21316](https://github.com/hashicorp/vault/pull/21316)]
* secrets/pki: Support setting both maintain_stored_certificate_counts=false and publish_stored_certificate_count_metrics=false explicitly in tidy config. [[GH-20664](https://github.com/hashicorp/vault/pull/20664)]
* secrets/transform (enterprise): Fix nil panic when deleting a template with tokenization transformations present
* secrets/transform (enterprise): Grab shared locks for various read operations, only escalating to write locks if work is required
* serviceregistration: Fix bug where multiple nodes in a secondary cluster could be labelled active after updating the cluster's primary [[GH-21642](https://github.com/hashicorp/vault/pull/21642)]
* ui: Fixed an issue where editing an SSH role would clear `default_critical_options` and `default_extension` if left unchanged. [[GH-21739](https://github.com/hashicorp/vault/pull/21739)]
* ui: Surface DOMException error when browser settings prevent localStorage. [[GH-21503](https://github.com/hashicorp/vault/pull/21503)]

## 1.13.4
### June 21, 2023
BREAKING CHANGES:

* secrets/pki: Maintaining running count of certificates will be turned off by default.
To re-enable keeping these metrics available on the tidy status endpoint, enable
maintain_stored_certificate_counts on tidy-config, to also publish them to the
metrics consumer, enable publish_stored_certificate_count_metrics . [[GH-18186](https://github.com/hashicorp/vault/pull/18186)]

CHANGES:

* core: Bump Go version to 1.20.5.

FEATURES:

* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* core (enterprise): Add background worker for automatic reporting of billing
information. [[GH-19625](https://github.com/hashicorp/vault/pull/19625)]

IMPROVEMENTS:

* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]

BUG FIXES:

* agent: Fix bug with 'cache' stanza validation [[GH-20934](https://github.com/hashicorp/vault/pull/20934)]
* core (enterprise): Don't delete backend stored data that appears to be filterable
on this secondary if we don't have a corresponding mount entry.
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry.  Fix wrong lock used in ListAuths link meta interface implementation. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs
* replication (enterprise): Fix regression causing token creation against a role
with a new entity alias to be incorrectly forwarded from perf standbys. [[GH-21100](https://github.com/hashicorp/vault/pull/21100)]
* storage/raft: Fix race where new follower joining can get pruned by dead server cleanup. [[GH-20986](https://github.com/hashicorp/vault/pull/20986)]

## 1.13.3
### June 08, 2023

CHANGES:

* core: Bump Go version to 1.20.4.
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]
* replication (enterprise): Add a new parameter for the update-primary API call
that allows for setting of the primary cluster addresses directly, instead of
via a token.
* storage/aerospike: Aerospike storage shouldn't be used on 32-bit architectures and is now unsupported on them. [[GH-20825](https://github.com/hashicorp/vault/pull/20825)]

IMPROVEMENTS:

* Add debug symbols back to builds to fix Dynatrace support [[GH-20519](https://github.com/hashicorp/vault/pull/20519)]
* audit: add a `mount_point` field to audit requests and response entries [[GH-20411](https://github.com/hashicorp/vault/pull/20411)]
* autopilot: Update version to v0.2.0 to add better support for respecting min quorum [[GH-19472](https://github.com/hashicorp/vault/pull/19472)]
* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* core: Add possibility to decode a generated encoded root token via the rest API [[GH-20595](https://github.com/hashicorp/vault/pull/20595)]
* core: include namespace path in granting_policies block of audit log
* core: report intermediate error messages during request forwarding [[GH-20643](https://github.com/hashicorp/vault/pull/20643)]
* openapi: Fix generated types for duration strings [[GH-20841](https://github.com/hashicorp/vault/pull/20841)]
* sdk/framework: Fix non-deterministic ordering of 'required' fields in OpenAPI spec [[GH-20881](https://github.com/hashicorp/vault/pull/20881)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]

BUG FIXES:

* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* cli: CLI should take days as a unit of time for ttl like flags [[GH-20477](https://github.com/hashicorp/vault/pull/20477)]
* cli: disable printing flags warnings messages for the ssh command [[GH-20502](https://github.com/hashicorp/vault/pull/20502)]
* command/server: fixes panic in Vault server command when running in recovery mode [[GH-20418](https://github.com/hashicorp/vault/pull/20418)]
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core/identity: Allow updates of only the custom-metadata for entity alias. [[GH-20368](https://github.com/hashicorp/vault/pull/20368)]
* core: Fix Forwarded Writer construction to correctly find active nodes, allowing PKI cross-cluster functionality to succeed on existing mounts.
* core: Fix writes to readonly storage on performance standbys when user lockout feature is enabled. [[GH-20783](https://github.com/hashicorp/vault/pull/20783)]
* core: prevent panic on login after namespace is deleted that had mfa enforcement [[GH-20375](https://github.com/hashicorp/vault/pull/20375)]
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* secrets/pki: Include per-issuer enable_aia_url_templating in issuer read endpoint. [[GH-20354](https://github.com/hashicorp/vault/pull/20354)]
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation
* secrets/transform: Added importing of keys and key versions into the Transform secrets engine using the command 'vault transform import' and 'vault transform import-version'. [[GH-20668](https://github.com/hashicorp/vault/pull/20668)]
* secrets/transit: Fix export of HMAC-only key, correctly exporting the key used for sign operations. For consumers of the previously incorrect key, use the plaintext export to retrieve these incorrect keys and import them as new versions.
secrets/transit: Fix bug related to shorter dedicated HMAC key sizing.
sdk/helper/keysutil: New HMAC type policies will have HMACKey equal to Key and be copied over on import. [[GH-20864](https://github.com/hashicorp/vault/pull/20864)]
* ui: Fixes issue unsealing cluster for seal types other than shamir [[GH-20897](https://github.com/hashicorp/vault/pull/20897)]
* ui: fixes issue creating mfa login enforcement from method enforcements tab [[GH-20603](https://github.com/hashicorp/vault/pull/20603)]
* ui: fixes key_bits and signature_bits reverting to default values when editing a pki role [[GH-20907](https://github.com/hashicorp/vault/pull/20907)]

## 1.13.2
### April 26, 2023

CHANGES:

* core: Bump Go version to 1.20.3.

SECURITY:

* core/seal: Fix handling of HMACing of seal-wrapped storage entries from HSMs using CKM_AES_CBC or CKM_AES_CBC_PAD which may have allowed an attacker to conduct a padding oracle attack. This vulnerability, CVE-2023-2197, affects Vault from 1.13.0 up to 1.13.1 and was fixed in 1.13.2. [[HCSEC-2023-14](https://discuss.hashicorp.com/t/hcsec-2023-14-vault-enterprise-vulnerable-to-padding-oracle-attacks-when-using-a-cbc-based-encryption-mechanism-with-a-hsm/53322)]

IMPROVEMENTS:

* Add debug symbols back to builds to fix Dynatrace support [[GH-20294](https://github.com/hashicorp/vault/pull/20294)]
* cli/namespace: Add detailed flag to output additional namespace information
such as namespace IDs and custom metadata. [[GH-20243](https://github.com/hashicorp/vault/pull/20243)]
* core/activity: add an endpoint to write test activity log data, guarded by a build flag [[GH-20019](https://github.com/hashicorp/vault/pull/20019)]
* core: Add a `raft` sub-field to the `storage` and `ha_storage` details provided by the
`/sys/config/state/sanitized` endpoint in order to include the `max_entry_size`. [[GH-20044](https://github.com/hashicorp/vault/pull/20044)]
* core: include reason for ErrReadOnly on PBPWF writing failures
* sdk/ldaputil: added `connection_timeout` to tune connection timeout duration
for all LDAP plugins. [[GH-20144](https://github.com/hashicorp/vault/pull/20144)]
* secrets/pki: Decrease size and improve compatibility of OCSP responses by removing issuer certificate. [[GH-20201](https://github.com/hashicorp/vault/pull/20201)]
* sys/wrapping: Add example how to unwrap without authentication in Vault [[GH-20109](https://github.com/hashicorp/vault/pull/20109)]
* ui: Allows license-banners to be dismissed. Saves preferences in localStorage. [[GH-19116](https://github.com/hashicorp/vault/pull/19116)]

BUG FIXES:

* auth/ldap: Add max_page_size configurable to LDAP configuration [[GH-19032](https://github.com/hashicorp/vault/pull/19032)]
* command/server: Fix incorrect paths in generated config for `-dev-tls` flag on Windows [[GH-20257](https://github.com/hashicorp/vault/pull/20257)]
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Fix regression breaking non-raft clusters whose nodes share the same cluster_addr/api_addr. [[GH-19721](https://github.com/hashicorp/vault/pull/19721)]
* helper/random: Fix race condition in string generator helper [[GH-19875](https://github.com/hashicorp/vault/pull/19875)]
* kmip (enterprise): Fix a problem decrypting with keys that have no Process Start Date attribute.
* pki: Fix automatically turning off CRL signing on upgrade to Vault >= 1.12, if CA Key Usage disallows it [[GH-20220](https://github.com/hashicorp/vault/pull/20220)]
* replication (enterprise): Fix a caching issue when replicating filtered data to
a performance secondary. This resulted in the data being set to nil in the cache
and a "invalid value" error being returned from the API.
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* sdk/helper/ocsp: Workaround bug in Go's ocsp.ParseResponse(...), causing validation to fail with embedded CA certificates.
auth/cert: Fix OCSP validation against Vault's PKI engine. [[GH-20181](https://github.com/hashicorp/vault/pull/20181)]
* secrets/aws: Revert changes that removed the lease on STS credentials, while leaving the new ttl field in place. [[GH-20034](https://github.com/hashicorp/vault/pull/20034)]
* secrets/pki: Ensure cross-cluster delta WAL write failure only logs to avoid unattended forwarding. [[GH-20057](https://github.com/hashicorp/vault/pull/20057)]
* secrets/pki: Fix building of unified delta CRLs and recovery during unified delta WAL write failures. [[GH-20058](https://github.com/hashicorp/vault/pull/20058)]
* secrets/pki: Fix patching of leaf_not_after_behavior on issuers. [[GH-20341](https://github.com/hashicorp/vault/pull/20341)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* ui: Fix OIDC provider logo showing when domain doesn't match [[GH-20263](https://github.com/hashicorp/vault/pull/20263)]
* ui: Fix bad link to namespace when namespace name includes `.` [[GH-19799](https://github.com/hashicorp/vault/pull/19799)]
* ui: fixes browser console formatting for help command output [[GH-20064](https://github.com/hashicorp/vault/pull/20064)]
* ui: fixes remaining doc links to include /vault in path [[GH-20070](https://github.com/hashicorp/vault/pull/20070)]
* ui: remove use of htmlSafe except when first sanitized [[GH-20235](https://github.com/hashicorp/vault/pull/20235)]
* website/docs: Fix Kubernetes Auth Code Example to use the correct whitespace in import. [[GH-20216](https://github.com/hashicorp/vault/pull/20216)]

## 1.13.1
### March 29, 2023

SECURITY:

* storage/mssql: When using Vault’s community-supported Microsoft SQL (MSSQL) database storage backend, a privileged attacker with the ability to write arbitrary data to Vault’s configuration may be able to perform arbitrary SQL commands on the underlying database server through Vault. This vulnerability, CVE-2023-0620, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-12](https://discuss.hashicorp.com/t/hcsec-2023-12-vault-s-microsoft-sql-database-storage-backend-vulnerable-to-sql-injection-via-configuration-file/52080)]
* secrets/pki: Vault’s PKI mount issuer endpoints did not correctly authorize access to remove an issuer or modify issuer metadata, potentially resulting in denial of service of the PKI mount. This bug did not affect public or private key material, trust chains or certificate issuance. This vulnerability, CVE-2023-0665, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-11](https://discuss.hashicorp.com/t/hcsec-2023-11-vault-s-pki-issuer-endpoint-did-not-correctly-authorize-access-to-issuer-metadata/52079)]
* core: HashiCorp Vault’s implementation of Shamir’s secret sharing used precomputed table lookups, and was vulnerable to cache-timing attacks. An attacker with access to, and the ability to observe a large number of unseal operations on the host through a side channel may reduce the search space of a brute force effort to recover the Shamir shares. This vulnerability, CVE-2023-25000, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-10](https://discuss.hashicorp.com/t/hcsec-2023-10-vault-vulnerable-to-cache-timing-attacks-during-seal-and-unseal-operations/52078)]

IMPROVEMENTS:

* auth/github: Allow for an optional Github auth token environment variable to make authenticated requests when fetching org id
website/docs: Add docs for `VAULT_AUTH_CONFIG_GITHUB_TOKEN` environment variable when writing Github config [[GH-19244](https://github.com/hashicorp/vault/pull/19244)]
* core: Allow overriding gRPC connect timeout via VAULT_GRPC_MIN_CONNECT_TIMEOUT. This is an env var rather than a config setting because we don't expect this to ever be needed.  It's being added as a last-ditch
option in case all else fails for some replication issues we may not have fully reproduced. [[GH-19676](https://github.com/hashicorp/vault/pull/19676)]
* core: validate name identifiers in mssql physical storage backend prior use [[GH-19591](https://github.com/hashicorp/vault/pull/19591)]
* database/elasticsearch: Update error messages resulting from Elasticsearch API errors [[GH-19545](https://github.com/hashicorp/vault/pull/19545)]
* events: Suppress log warnings triggered when events are sent but the events system is not enabled. [[GH-19593](https://github.com/hashicorp/vault/pull/19593)]

BUG FIXES:

* agent: Fix panic when SIGHUP is issued to Agent while it has a non-TLS listener. [[GH-19483](https://github.com/hashicorp/vault/pull/19483)]
* core (enterprise): Attempt to reconnect to a PKCS#11 HSM if we retrieve a CKR_FUNCTION_FAILED error.
* core: Fixed issue with remounting mounts that have a non-trailing space in the 'to' or 'from' paths. [[GH-19585](https://github.com/hashicorp/vault/pull/19585)]
* kmip (enterprise): Do not require attribute Cryptographic Usage Mask when registering Secret Data managed objects.
* kmip (enterprise): Fix a problem forwarding some requests to the active node.
* openapi: Fix logic for labeling unauthenticated/sudo paths. [[GH-19600](https://github.com/hashicorp/vault/pull/19600)]
* secrets/ldap: Invalidates WAL entry for static role if `password_policy` has changed. [[GH-19640](https://github.com/hashicorp/vault/pull/19640)]
* secrets/pki: Fix PKI revocation request forwarding from standby nodes due to an error wrapping bug [[GH-19624](https://github.com/hashicorp/vault/pull/19624)]
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* ui: Fixes crypto.randomUUID error in unsecure contexts from third party ember-data library [[GH-19428](https://github.com/hashicorp/vault/pull/19428)]
* ui: fixes SSH engine config deletion [[GH-19448](https://github.com/hashicorp/vault/pull/19448)]
* ui: fixes issue navigating back a level using the breadcrumb from secret metadata view [[GH-19703](https://github.com/hashicorp/vault/pull/19703)]
* ui: fixes oidc tabs in auth form submitting with the root's default_role value after a namespace has been inputted [[GH-19541](https://github.com/hashicorp/vault/pull/19541)]
* ui: pass encodeBase64 param to HMAC transit-key-actions. [[GH-19429](https://github.com/hashicorp/vault/pull/19429)]
* ui: use URLSearchParams interface to capture namespace param from SSOs (ex. ADFS) with decoded state param in callback url [[GH-19460](https://github.com/hashicorp/vault/pull/19460)]

## 1.13.0
### March 01, 2023

SECURITY:

* secrets/ssh: removal of the deprecated dynamic keys mode. **When any remaining dynamic key leases expire**, an error stating `secret is unsupported by this backend` will be thrown by the lease manager. [[GH-18874](https://github.com/hashicorp/vault/pull/18874)]
* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* auth/alicloud: require the `role` field on login [[GH-19005](https://github.com/hashicorp/vault/pull/19005)]
* auth/approle: Add maximum length of 4096 for approle role_names, as this value results in HMAC calculation [[GH-17768](https://github.com/hashicorp/vault/pull/17768)]
* auth: Returns invalid credentials for ldap, userpass and approle when wrong credentials are provided for existent users.
This will only be used internally for implementing user lockout. [[GH-17104](https://github.com/hashicorp/vault/pull/17104)]
* core: Bump Go version to 1.20.1.
* core: Vault version has been moved out of sdk and into main vault module.
Plugins using sdk/useragent.String must instead use sdk/useragent.PluginString. [[GH-14229](https://github.com/hashicorp/vault/pull/14229)]
* logging: Removed legacy environment variable for log format ('LOGXI_FORMAT'), should use 'VAULT_LOG_FORMAT' instead [[GH-17822](https://github.com/hashicorp/vault/pull/17822)]
* plugins: Mounts can no longer be pinned to a specific _builtin_ version. Mounts previously pinned to a specific builtin version will now automatically upgrade to the latest builtin version, and may now be overridden if an unversioned plugin of the same name and type is registered. Mounts using plugin versions without `builtin` in their metadata remain unaffected. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* plugins: `GET /database/config/:name` endpoint now returns an additional `plugin_version` field in the response data. [[GH-16982](https://github.com/hashicorp/vault/pull/16982)]
* plugins: `GET /sys/auth/:path/tune` and `GET /sys/mounts/:path/tune` endpoints may now return an additional `plugin_version` field in the response data if set. [[GH-17167](https://github.com/hashicorp/vault/pull/17167)]
* plugins: `GET` for `/sys/auth`, `/sys/auth/:path`, `/sys/mounts`, and `/sys/mounts/:path` paths now return additional `plugin_version`, `running_plugin_version` and `running_sha256` fields in the response data for each mount. [[GH-17167](https://github.com/hashicorp/vault/pull/17167)]
* sdk: Remove version package, make useragent.String versionless. [[GH-19068](https://github.com/hashicorp/vault/pull/19068)]
* secrets/aws: do not create leases for non-renewable/non-revocable STS credentials to reduce storage calls [[GH-15869](https://github.com/hashicorp/vault/pull/15869)]
* secrets/gcpkms: Updated plugin from v0.13.0 to v0.14.0 [[GH-19063](https://github.com/hashicorp/vault/pull/19063)]
* sys/internal/inspect: Turns of this endpoint by default. A SIGHUP can now be used to reload the configs and turns this endpoint on.
* ui: Upgrade Ember to version 4.4.0 [[GH-17086](https://github.com/hashicorp/vault/pull/17086)]

FEATURES:

* **User lockout**: Ignore repeated bad credentials from the same user for a configured period of time. Enabled by default.
* **Azure Auth Managed Identities**: Allow any Azure resource that supports managed identities to authenticate with Vault [[GH-19077](https://github.com/hashicorp/vault/pull/19077)]
* **Azure Auth Rotate Root**: Add support for rotate root in Azure Auth engine [[GH-19077](https://github.com/hashicorp/vault/pull/19077)]
* **Event System (Alpha)**: Vault has a new opt-in experimental event system. Not yet suitable for production use. Events are currently only generated on writes to the KV secrets engine, but external plugins can also be updated to start generating events. [[GH-19194](https://github.com/hashicorp/vault/pull/19194)]
* **GCP Secrets Impersonated Account Support**: Add support for GCP service account impersonation, allowing callers to generate a GCP access token without requiring Vault to store or retrieve a GCP service account key for each role. [[GH-19018](https://github.com/hashicorp/vault/pull/19018)]
* **Kubernetes Secrets Engine UI**: Kubernetes is now available in the UI as a supported secrets engine. [[GH-17893](https://github.com/hashicorp/vault/pull/17893)]
* **New PKI UI**: Add beta support for new and improved PKI UI [[GH-18842](https://github.com/hashicorp/vault/pull/18842)]
* **PKI Cross-Cluster Revocations**: Revocation information can now be
synchronized across primary and performance replica clusters offering
a unified CRL/OCSP view of revocations across cluster boundaries. [[GH-19196](https://github.com/hashicorp/vault/pull/19196)]
* **Server UDS Listener**: Adding listener to Vault server to serve http request via unix domain socket [[GH-18227](https://github.com/hashicorp/vault/pull/18227)]
* **Transit managed keys**: The transit secrets engine now supports configuring and using managed keys
* **User Lockout**: Adds support to configure the user-lockout behaviour for failed logins to prevent
brute force attacks for userpass, approle and ldap auth methods. [[GH-19230](https://github.com/hashicorp/vault/pull/19230)]
* **VMSS Flex Authentication**: Adds support for Virtual Machine Scale Set Flex Authentication [[GH-19077](https://github.com/hashicorp/vault/pull/19077)]
* **Namespaces (enterprise)**: Added the ability to allow access to secrets and more to be shared across namespaces that do not share a namespace hierarchy. Using the new `sys/config/group-policy-application` API, policies can be configured to apply outside of namespace hierarchy, allowing this kind of cross-namespace sharing.
* **OpenAPI-based Go & .NET Client Libraries (Beta)**: We have now made available two new [[OpenAPI-based Go](https://github.com/hashicorp/vault-client-go/)] & [[OpenAPI-based .NET](https://github.com/hashicorp/vault-client-dotnet/)] Client libraries (beta). You can use them to perform various secret management operations easily from your applications.

IMPROVEMENTS:

* **Redis ElastiCache DB Engine**: Renamed configuration parameters for disambiguation; old parameters still supported for compatibility. [[GH-18752](https://github.com/hashicorp/vault/pull/18752)]
* Bump github.com/hashicorp/go-plugin version from 1.4.5 to 1.4.8 [[GH-19100](https://github.com/hashicorp/vault/pull/19100)]
* Reduced binary size [[GH-17678](https://github.com/hashicorp/vault/pull/17678)]
* agent/config: Allow config directories to be specified with -config, and allow multiple -configs to be supplied. [[GH-18403](https://github.com/hashicorp/vault/pull/18403)]
* agent: Add note in logs when starting Vault Agent indicating if the version differs to the Vault Server. [[GH-18684](https://github.com/hashicorp/vault/pull/18684)]
* agent: Added `token_file` auto-auth configuration to allow using a pre-existing token for Vault Agent. [[GH-18740](https://github.com/hashicorp/vault/pull/18740)]
* agent: Agent listeners can now be to be the `metrics_only` role, serving only metrics, as part of the listener's new top level `role` option. [[GH-18101](https://github.com/hashicorp/vault/pull/18101)]
* agent: Configured Vault Agent listeners now listen without the need for caching to be configured. [[GH-18137](https://github.com/hashicorp/vault/pull/18137)]
* agent: allows some parts of config to be reloaded without requiring a restart. [[GH-18638](https://github.com/hashicorp/vault/pull/18638)]
* agent: fix incorrectly used loop variables in parallel tests and when finalizing seals [[GH-16872](https://github.com/hashicorp/vault/pull/16872)]
* api: Remove dependency on sdk module. [[GH-18962](https://github.com/hashicorp/vault/pull/18962)]
* api: Support VAULT_DISABLE_REDIRECTS environment variable (and --disable-redirects flag) to disable default client behavior and prevent the client following any redirection responses. [[GH-17352](https://github.com/hashicorp/vault/pull/17352)]
* audit: Add `elide_list_responses` option, providing a countermeasure for a common source of oversized audit log entries [[GH-18128](https://github.com/hashicorp/vault/pull/18128)]
* audit: Include stack trace when audit logging recovers from a panic. [[GH-18121](https://github.com/hashicorp/vault/pull/18121)]
* auth/alicloud: upgrades dependencies [[GH-18021](https://github.com/hashicorp/vault/pull/18021)]
* auth/azure: Adds support for authentication with Managed Service Identity (MSI) from a
Virtual Machine Scale Set (VMSS) in flexible orchestration mode. [[GH-17540](https://github.com/hashicorp/vault/pull/17540)]
* auth/azure: upgrades dependencies [[GH-17857](https://github.com/hashicorp/vault/pull/17857)]
* auth/cert: Add configurable support for validating client certs with OCSP. [[GH-17093](https://github.com/hashicorp/vault/pull/17093)]
* auth/cert: Support listing provisioned CRLs within the mount. [[GH-18043](https://github.com/hashicorp/vault/pull/18043)]
* auth/cf: Remove incorrect usage of CreateOperation from path_config [[GH-19098](https://github.com/hashicorp/vault/pull/19098)]
* auth/gcp: Upgrades dependencies [[GH-17858](https://github.com/hashicorp/vault/pull/17858)]
* auth/oidc: Adds `abort_on_error` parameter to CLI login command to help in non-interactive contexts [[GH-19076](https://github.com/hashicorp/vault/pull/19076)]
* auth/oidc: Adds ability to set Google Workspace domain for groups search [[GH-19076](https://github.com/hashicorp/vault/pull/19076)]
* auth/token (enterprise): Allow batch token creation in perfStandby nodes
* auth: Allow naming login MFA methods and using those names instead of IDs in satisfying MFA requirement for requests.
Make passcode arguments consistent across login MFA method types. [[GH-18610](https://github.com/hashicorp/vault/pull/18610)]
* auth: Provide an IP address of the requests from Vault to a Duo challenge after successful authentication. [[GH-18811](https://github.com/hashicorp/vault/pull/18811)]
* autopilot: Update version to v.0.2.0 to add better support for respecting min quorum
* cli/kv: improve kv CLI to remove data or custom metadata using kv patch [[GH-18067](https://github.com/hashicorp/vault/pull/18067)]
* cli/pki: Add List-Intermediates functionality to pki client. [[GH-18463](https://github.com/hashicorp/vault/pull/18463)]
* cli/pki: Add health-check subcommand to evaluate the health of a PKI instance. [[GH-17750](https://github.com/hashicorp/vault/pull/17750)]
* cli/pki: Add pki issue command, which creates a CSR, has a vault mount sign it, then reimports it. [[GH-18467](https://github.com/hashicorp/vault/pull/18467)]
* cli/pki: Added "Reissue" command which allows extracting fields from an existing certificate to create a new certificate. [[GH-18499](https://github.com/hashicorp/vault/pull/18499)]
* cli/pki: Change the pki health-check --list default config output to JSON so it's a usable configuration file [[GH-19269](https://github.com/hashicorp/vault/pull/19269)]
* cli: Add support for creating requests to existing non-KVv2 PATCH-capable endpoints. [[GH-17650](https://github.com/hashicorp/vault/pull/17650)]
* cli: Add transit import key helper commands for BYOK to Transit/Transform. [[GH-18887](https://github.com/hashicorp/vault/pull/18887)]
* cli: Support the -format=raw option, to read non-JSON Vault endpoints and original response bodies. [[GH-14945](https://github.com/hashicorp/vault/pull/14945)]
* cli: updated `vault operator rekey` prompts to describe recovery keys when `-target=recovery` [[GH-18892](https://github.com/hashicorp/vault/pull/18892)]
* client/pki: Add a new command verify-sign which checks the relationship between two certificates. [[GH-18437](https://github.com/hashicorp/vault/pull/18437)]
* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* core/identity: Add machine-readable output to body of response upon alias clash during entity merge [[GH-17459](https://github.com/hashicorp/vault/pull/17459)]
* core/server: Added an environment variable to write goroutine stacktraces to a
temporary file for SIGUSR2 signals. [[GH-17929](https://github.com/hashicorp/vault/pull/17929)]
* core: Add RPCs to read and update userFailedLoginInfo map
* core: Add experiments system and `events.alpha1` experiment. [[GH-18682](https://github.com/hashicorp/vault/pull/18682)]
* core: Add read support to `sys/loggers` and `sys/loggers/:name` endpoints [[GH-17979](https://github.com/hashicorp/vault/pull/17979)]
* core: Add user lockout field to config and configuring this for auth mount using auth tune to prevent brute forcing in auth methods [[GH-17338](https://github.com/hashicorp/vault/pull/17338)]
* core: Add vault.core.locked_users telemetry metric to emit information about total number of locked users. [[GH-18718](https://github.com/hashicorp/vault/pull/18718)]
* core: Added sys/locked-users endpoint to list locked users. Changed api endpoint from
sys/lockedusers/[mount_accessor]/unlock/[alias_identifier] to sys/locked-users/[mount_accessor]/unlock/[alias_identifier]. [[GH-18675](https://github.com/hashicorp/vault/pull/18675)]
* core: Added sys/lockedusers/[mount_accessor]/unlock/[alias_identifier] endpoint to unlock an user
with given mount_accessor and alias_identifier if locked [[GH-18279](https://github.com/hashicorp/vault/pull/18279)]
* core: Added warning to /sys/seal-status and vault status command if potentially dangerous behaviour overrides are being used. [[GH-17855](https://github.com/hashicorp/vault/pull/17855)]
* core: Implemented background thread to update locked user entries every 15 minutes to prevent brute forcing in auth methods. [[GH-18673](https://github.com/hashicorp/vault/pull/18673)]
* core: License location is no longer cache exempt, meaning sys/health will not contribute as greatly to storage load when using consul as a storage backend. [[GH-17265](https://github.com/hashicorp/vault/pull/17265)]
* core: Update protoc from 3.21.5 to 3.21.7 [[GH-17499](https://github.com/hashicorp/vault/pull/17499)]
* core: add `detect_deadlocks` config to optionally detect core state deadlocks [[GH-18604](https://github.com/hashicorp/vault/pull/18604)]
* core: added changes for user lockout workflow. [[GH-17951](https://github.com/hashicorp/vault/pull/17951)]
* core: parallelize backend initialization to improve startup time for large numbers of mounts. [[GH-18244](https://github.com/hashicorp/vault/pull/18244)]
* database/postgres: Support multiline strings for revocation statements. [[GH-18632](https://github.com/hashicorp/vault/pull/18632)]
* database/redis-elasticache: changed config argument names for disambiguation [[GH-19044](https://github.com/hashicorp/vault/pull/19044)]
* database/snowflake: Allow parallel requests to Snowflake [[GH-17593](https://github.com/hashicorp/vault/pull/17593)]
* hcp/connectivity: Add foundational OSS support for opt-in secure communication between self-managed Vault nodes and [HashiCorp Cloud Platform](https://cloud.hashicorp.com) [[GH-18228](https://github.com/hashicorp/vault/pull/18228)]
* hcp/connectivity: Include HCP organization, project, and resource ID in server startup logs [[GH-18315](https://github.com/hashicorp/vault/pull/18315)]
* hcp/connectivity: Only update SCADA session metadata if status changes [[GH-18585](https://github.com/hashicorp/vault/pull/18585)]
* hcp/status: Add cluster-level status information [[GH-18351](https://github.com/hashicorp/vault/pull/18351)]
* hcp/status: Expand node-level status information [[GH-18302](https://github.com/hashicorp/vault/pull/18302)]
* logging: Vault Agent supports logging to a specified file path via environment variable, CLI or config [[GH-17841](https://github.com/hashicorp/vault/pull/17841)]
* logging: Vault agent and server commands support log file and log rotation. [[GH-18031](https://github.com/hashicorp/vault/pull/18031)]
* migration: allow parallelization of key migration for `vault operator migrate` in order to speed up a migration. [[GH-18817](https://github.com/hashicorp/vault/pull/18817)]
* namespaces (enterprise): Add new API, `sys/config/group-policy-application`, to allow group policies to be configurable
to apply to a group in `any` namespace. The default, `within_namespace_hierarchy`, is the current behaviour.
* openapi: Add default values to thing_mount_path parameters [[GH-18935](https://github.com/hashicorp/vault/pull/18935)]
* openapi: Add logic to generate openapi response structures [[GH-18192](https://github.com/hashicorp/vault/pull/18192)]
* openapi: Add openapi response definitions to approle/path_login.go & approle/path_tidy_user_id.go [[GH-18772](https://github.com/hashicorp/vault/pull/18772)]
* openapi: Add openapi response definitions to approle/path_role.go [[GH-18198](https://github.com/hashicorp/vault/pull/18198)]
* openapi: Change gen_openapi.sh to generate schema with generic mount paths [[GH-18934](https://github.com/hashicorp/vault/pull/18934)]
* openapi: Mark request body objects as required [[GH-17909](https://github.com/hashicorp/vault/pull/17909)]
* openapi: add openapi response defintions to /sys/audit endpoints [[GH-18456](https://github.com/hashicorp/vault/pull/18456)]
* openapi: generic_mount_paths: Move implementation fully into server, rather than partially in plugin framework; recognize all 4 singleton mounts (auth/token, cubbyhole, identity, system) rather than just 2; change parameter from `{mountPath}` to `{<type>_mount_path}` [[GH-18663](https://github.com/hashicorp/vault/pull/18663)]
* plugins: Add plugin version information to key plugin lifecycle log lines. [[GH-17430](https://github.com/hashicorp/vault/pull/17430)]
* plugins: Allow selecting builtin plugins by their reported semantic version of the form `vX.Y.Z+builtin` or `vX.Y.Z+builtin.vault`. [[GH-17289](https://github.com/hashicorp/vault/pull/17289)]
* plugins: Let Vault unseal and mount deprecated builtin plugins in a
deactivated state if this is not the first unseal after an upgrade. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* plugins: Mark app-id auth method Removed and remove the plugin code. [[GH-18039](https://github.com/hashicorp/vault/pull/18039)]
* plugins: Mark logical database plugins Removed and remove the plugin code. [[GH-18039](https://github.com/hashicorp/vault/pull/18039)]
* sdk/ldap: Added support for paging when searching for groups using group filters [[GH-17640](https://github.com/hashicorp/vault/pull/17640)]
* sdk: Add response schema validation method framework/FieldData.ValidateStrict and two test helpers (ValidateResponse, ValidateResponseData) [[GH-18635](https://github.com/hashicorp/vault/pull/18635)]
* sdk: Adding FindResponseSchema test helper to assist with response schema validation in tests [[GH-18636](https://github.com/hashicorp/vault/pull/18636)]
* secrets/aws: Update dependencies [[PR-17747](https://github.com/hashicorp/vault/pull/17747)] [[GH-17747](https://github.com/hashicorp/vault/pull/17747)]
* secrets/azure: Adds ability to persist an application for the lifetime of a role. [[GH-19096](https://github.com/hashicorp/vault/pull/19096)]
* secrets/azure: upgrades dependencies [[GH-17964](https://github.com/hashicorp/vault/pull/17964)]
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* secrets/gcp: Upgrades dependencies [[GH-17871](https://github.com/hashicorp/vault/pull/17871)]
* secrets/kubernetes: Add /check endpoint to determine if environment variables are set [[GH-18](https://github.com/hashicorp/vault-plugin-secrets-kubernetes/pull/18)] [[GH-18587](https://github.com/hashicorp/vault/pull/18587)]
* secrets/kubernetes: add /check endpoint to determine if environment variables are set [[GH-19084](https://github.com/hashicorp/vault/pull/19084)]
* secrets/kv: Emit events on write if events system enabled [[GH-19145](https://github.com/hashicorp/vault/pull/19145)]
* secrets/kv: make upgrade synchronous when no keys to upgrade [[GH-19056](https://github.com/hashicorp/vault/pull/19056)]
* secrets/kv: new KVv2 mounts and KVv1 mounts without any keys will upgrade synchronously, allowing for instant use [[GH-17406](https://github.com/hashicorp/vault/pull/17406)]
* secrets/pki: Add a new API that returns the serial numbers of revoked certificates on the local cluster [[GH-17779](https://github.com/hashicorp/vault/pull/17779)]
* secrets/pki: Add support to specify signature bits when generating CSRs through intermediate/generate apis [[GH-17388](https://github.com/hashicorp/vault/pull/17388)]
* secrets/pki: Added a new API that allows external actors to craft a CRL through JSON parameters [[GH-18040](https://github.com/hashicorp/vault/pull/18040)]
* secrets/pki: Allow UserID Field (https://www.rfc-editor.org/rfc/rfc1274#section-9.3.1) to be set on Certificates when
allowed by role [[GH-18397](https://github.com/hashicorp/vault/pull/18397)]
* secrets/pki: Allow issuer creation, import to change default issuer via `default_follows_latest_issuer`. [[GH-17824](https://github.com/hashicorp/vault/pull/17824)]
* secrets/pki: Allow templating performance replication cluster- and issuer-specific AIA URLs. [[GH-18199](https://github.com/hashicorp/vault/pull/18199)]
* secrets/pki: Allow tidying of expired issuer certificates. [[GH-17823](https://github.com/hashicorp/vault/pull/17823)]
* secrets/pki: Allow tidying of the legacy ca_bundle, improving startup on post-migrated, seal-wrapped PKI mounts. [[GH-18645](https://github.com/hashicorp/vault/pull/18645)]
* secrets/pki: Respond with written data to `config/auto-tidy`, `config/crl`, and `roles/:role`. [[GH-18222](https://github.com/hashicorp/vault/pull/18222)]
* secrets/pki: Return issuer_id and issuer_name on /issuer/:issuer_ref/json endpoint. [[GH-18482](https://github.com/hashicorp/vault/pull/18482)]
* secrets/pki: Return new fields revocation_time_rfc3339 and issuer_id to existing certificate serial lookup api if it is revoked [[GH-17774](https://github.com/hashicorp/vault/pull/17774)]
* secrets/ssh: Allow removing SSH host keys from the dynamic keys feature. [[GH-18939](https://github.com/hashicorp/vault/pull/18939)]
* secrets/ssh: Evaluate ssh validprincipals user template before splitting [[GH-16622](https://github.com/hashicorp/vault/pull/16622)]
* secrets/transit: Add an optional reference field to batch operation items
which is repeated on batch responses to help more easily correlate inputs with outputs. [[GH-18243](https://github.com/hashicorp/vault/pull/18243)]
* secrets/transit: Add associated_data parameter for additional authenticated data in AEAD ciphers [[GH-17638](https://github.com/hashicorp/vault/pull/17638)]
* secrets/transit: Add support for PKCSv1_5_NoOID RSA signatures [[GH-17636](https://github.com/hashicorp/vault/pull/17636)]
* secrets/transit: Allow configuring whether upsert of keys is allowed. [[GH-18272](https://github.com/hashicorp/vault/pull/18272)]
* storage/raft: Add `retry_join_as_non_voter` config option. [[GH-18030](https://github.com/hashicorp/vault/pull/18030)]
* storage/raft: add additional raft metrics relating to applied index and heartbeating; also ensure OSS standbys emit periodic metrics. [[GH-12166](https://github.com/hashicorp/vault/pull/12166)]
* sys/internal/inspect: Creates an endpoint to look to inspect internal subsystems. [[GH-17789](https://github.com/hashicorp/vault/pull/17789)]
* sys/internal/inspect: Creates an endpoint to look to inspect internal subsystems.
* ui: Add algorithm-signer as a SSH Secrets Engine UI field [[GH-10299](https://github.com/hashicorp/vault/pull/10299)]
* ui: Add inline policy creation when creating an identity entity or group [[GH-17749](https://github.com/hashicorp/vault/pull/17749)]
* ui: Added JWT authentication warning message about blocked pop-up windows and web browser settings. [[GH-18787](https://github.com/hashicorp/vault/pull/18787)]
* ui: Enable typescript for future development [[GH-17927](https://github.com/hashicorp/vault/pull/17927)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]
* ui: adds allowed_response_headers as param for secret engine mount config [[GH-19216](https://github.com/hashicorp/vault/pull/19216)]
* ui: consolidate all <a> tag usage [[GH-17866](https://github.com/hashicorp/vault/pull/17866)]
* ui: mfa: use proper request id generation [[GH-17835](https://github.com/hashicorp/vault/pull/17835)]
* ui: remove wizard [[GH-19220](https://github.com/hashicorp/vault/pull/19220)]
* ui: update DocLink component to use new host url: developer.hashicorp.com [[GH-18374](https://github.com/hashicorp/vault/pull/18374)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]
* ui: use the combined activity log (partial + historic) API for client count dashboard and remove use of monthly endpoint [[GH-17575](https://github.com/hashicorp/vault/pull/17575)]
* vault/diagnose: Upgrade `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, `go.opentelemetry.io/otel/trace` to v1.11.2 [[GH-18589](https://github.com/hashicorp/vault/pull/18589)]

DEPRECATIONS:

* secrets/ad: Marks the Active Directory (AD) secrets engine as deprecated. [[GH-19334](https://github.com/hashicorp/vault/pull/19334)]

BUG FIXES:

* api: Remove timeout logic from ReadRaw functions and add ReadRawWithContext [[GH-18708](https://github.com/hashicorp/vault/pull/18708)]
* auth/alicloud: fix regression in vault login command that caused login to fail [[GH-19005](https://github.com/hashicorp/vault/pull/19005)]
* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/cert: Address a race condition accessing the loaded crls without a lock [[GH-18945](https://github.com/hashicorp/vault/pull/18945)]
* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#173](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/173)] [[GH-18716](https://github.com/hashicorp/vault/pull/18716)]
* auth/kubernetes: fixes and dep updates for the auth-kubernetes plugin (see plugin changelog for details) [[GH-19094](https://github.com/hashicorp/vault/pull/19094)]
* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* cli/pki: Decode integer values properly in health-check configuration file [[GH-19265](https://github.com/hashicorp/vault/pull/19265)]
* cli/pki: Fix path for role health-check warning messages [[GH-19274](https://github.com/hashicorp/vault/pull/19274)]
* cli/pki: Properly report permission issues within health-check mount tune checks [[GH-19276](https://github.com/hashicorp/vault/pull/19276)]
* cli/transit: Fix import, import-version command invocation [[GH-19373](https://github.com/hashicorp/vault/pull/19373)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* cli: Remove empty table heading for `vault secrets list -detailed` output. [[GH-17577](https://github.com/hashicorp/vault/pull/17577)]
* command/namespace: Fix vault cli namespace patch examples in help text. [[GH-18143](https://github.com/hashicorp/vault/pull/18143)]
* core (enterprise): Fix missing quotation mark in error message
* core (enterprise): Fix panic that could occur with SSCT alongside invoking external plugins for revocation.
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core (enterprise): Supported storage check in `vault server` command will no longer prevent startup. Instead, a warning will be logged if configured to use storage backend other than `raft` or `consul`.
* core/activity: add namespace breakdown for new clients when date range spans multiple months, including the current month. [[GH-18766](https://github.com/hashicorp/vault/pull/18766)]
* core/activity: de-duplicate namespaces when historical and current month data are mixed [[GH-18452](https://github.com/hashicorp/vault/pull/18452)]
* core/activity: fix the end_date returned from the activity log endpoint when partial counts are computed [[GH-17856](https://github.com/hashicorp/vault/pull/17856)]
* core/activity: include mount counts when de-duplicating current and historical month data [[GH-18598](https://github.com/hashicorp/vault/pull/18598)]
* core/activity: report mount paths (rather than mount accessors) in current month activity log counts and include deleted mount paths in precomputed queries. [[GH-18916](https://github.com/hashicorp/vault/pull/18916)]
* core/activity: return partial month counts when querying a historical date range and no historical data exists. [[GH-17935](https://github.com/hashicorp/vault/pull/17935)]
* core/auth: Return a 403 instead of a 500 for wrapping requests when token is not provided [[GH-18859](https://github.com/hashicorp/vault/pull/18859)]
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): Switch to using hash length as PSS Salt length within the test/sign api for better PKCS#11 compatibility
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core/seal: Fix regression handling of the key_id parameter in seal configuration HCL. [[GH-17612](https://github.com/hashicorp/vault/pull/17612)]
* core: Fix panic caused in Vault Agent when rendering certificate templates [[GH-17419](https://github.com/hashicorp/vault/pull/17419)]
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: Fix spurious `permission denied` for all HelpOperations on sudo-protected paths [[GH-18568](https://github.com/hashicorp/vault/pull/18568)]
* core: Fix vault operator init command to show the right curl string with -output-curl-string and right policy hcl with -output-policy [[GH-17514](https://github.com/hashicorp/vault/pull/17514)]
* core: Fixes spurious warnings being emitted relating to "unknown or unsupported fields" for JSON config [[GH-17660](https://github.com/hashicorp/vault/pull/17660)]
* core: Linux packages now have vendor label and set the default label to HashiCorp.
This fix is implemented for any future releases, but will not be updated for historical releases.
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* core: Refactor lock grabbing code to simplify stateLock deadlock investigations [[GH-17187](https://github.com/hashicorp/vault/pull/17187)]
* core: fix GPG encryption to support subkeys. [[GH-16224](https://github.com/hashicorp/vault/pull/16224)]
* core: fix a start up race condition where performance standbys could go into a
mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: fix race when using SystemView.ReplicationState outside of a request context [[GH-17186](https://github.com/hashicorp/vault/pull/17186)]
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* core: trying to unseal with the wrong key now returns HTTP 400 [[GH-17836](https://github.com/hashicorp/vault/pull/17836)]
* credential/cert: adds error message if no tls connection is found during the AliasLookahead operation [[GH-17904](https://github.com/hashicorp/vault/pull/17904)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* expiration: Prevent panics on perf standbys when an irrevocable lease gets deleted. [[GH-18401](https://github.com/hashicorp/vault/pull/18401)]
* kmip (enterprise): Fix a problem with some multi-part MAC Verify operations.
* kmip (enterprise): Only require data to be full blocks on encrypt/decrypt operations using CBC and ECB block cipher modes.
* license (enterprise): Fix bug where license would update even if the license didn't change.
* licensing (enterprise): update autoloaded license cache after reload
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* openapi: Fix many incorrect details in generated API spec, by using better techniques to parse path regexps [[GH-18554](https://github.com/hashicorp/vault/pull/18554)]
* openapi: fix gen_openapi.sh script to correctly load vault plugins [[GH-17752](https://github.com/hashicorp/vault/pull/17752)]
* plugins/kv: KV v2 returns 404 instead of 500 for request paths that incorrectly include a trailing slash. [[GH-17339](https://github.com/hashicorp/vault/pull/17339)]
* plugins: Allow running external plugins which override deprecated builtins. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* plugins: Corrected the path to check permissions on when the registered plugin name does not match the plugin binary's filename. [[GH-17340](https://github.com/hashicorp/vault/pull/17340)]
* plugins: Listing all plugins while audit logging is enabled will no longer result in an internal server error. [[GH-18173](https://github.com/hashicorp/vault/pull/18173)]
* plugins: Only report deprecation status for builtin plugins. [[GH-17816](https://github.com/hashicorp/vault/pull/17816)]
* plugins: Skip loading but still mount data associated with missing plugins on unseal. [[GH-18189](https://github.com/hashicorp/vault/pull/18189)]
* plugins: Vault upgrades will no longer fail if a mount has been created using an explicit builtin plugin version. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* sdk: Don't panic if system view or storage methods called during plugin setup. [[GH-18210](https://github.com/hashicorp/vault/pull/18210)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/ad: Fix bug where updates to config would fail if password isn't provided [[GH-19061](https://github.com/hashicorp/vault/pull/19061)]
* secrets/gcp: fix issue where IAM bindings were not preserved during policy update [[GH-19018](https://github.com/hashicorp/vault/pull/19018)]
* secrets/mongodb-atlas: Fix a bug that did not allow WAL rollback to handle partial failures when creating API keys [[GH-19111](https://github.com/hashicorp/vault/pull/19111)]
* secrets/pki: Address nil panic when an empty POST request is sent to the OCSP handler [[GH-18184](https://github.com/hashicorp/vault/pull/18184)]
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17385](https://github.com/hashicorp/vault/pull/17385)]
* secrets/pki: Fix upgrade of missing expiry, delta_rebuild_interval by setting them to the default. [[GH-17693](https://github.com/hashicorp/vault/pull/17693)]
* secrets/pki: Fixes duplicate otherName in certificates created by the sign-verbatim endpoint. [[GH-16700](https://github.com/hashicorp/vault/pull/16700)]
* secrets/pki: OCSP GET request parameter was not being URL unescaped before processing. [[GH-18938](https://github.com/hashicorp/vault/pull/18938)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* secrets/pki: Revert fix for PR [18938](https://github.com/hashicorp/vault/pull/18938) [[GH-19037](https://github.com/hashicorp/vault/pull/19037)]
* secrets/pki: consistently use UTC for CA's notAfter exceeded error message [[GH-18984](https://github.com/hashicorp/vault/pull/18984)]
* secrets/pki: fix race between tidy's cert counting and tidy status reporting. [[GH-18899](https://github.com/hashicorp/vault/pull/18899)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* secrets/transit: Honor `partial_success_response_code` on decryption failures. [[GH-18310](https://github.com/hashicorp/vault/pull/18310)]
* server/config:  Use file.Stat when checking file permissions when VAULT_ENABLE_FILE_PERMISSIONS_CHECK is enabled [[GH-19311](https://github.com/hashicorp/vault/pull/19311)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* storage/raft: Fix race with follower heartbeat tracker during teardown. [[GH-18704](https://github.com/hashicorp/vault/pull/18704)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]
* ui: Remove `default` and add `default-service` and `default-batch` to UI token_type for auth mount and tuning. [[GH-19290](https://github.com/hashicorp/vault/pull/19290)]
* ui: Remove default value of 30 to TtlPicker2 if no value is passed in. [[GH-17376](https://github.com/hashicorp/vault/pull/17376)]
* ui: allow selection of "default" for ssh algorithm_signer in web interface [[GH-17894](https://github.com/hashicorp/vault/pull/17894)]
* ui: cleanup unsaved auth method ember data record when navigating away from mount backend form [[GH-18651](https://github.com/hashicorp/vault/pull/18651)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]
* ui: fixes reliance on secure context (https) by removing methods using the Crypto interface [[GH-19403](https://github.com/hashicorp/vault/pull/19403)]
* ui: show Get credentials button for static roles detail page when a user has the proper permissions. [[GH-19190](https://github.com/hashicorp/vault/pull/19190)]

## 1.12.10
### August 30, 2023

CHANGES:

* core: Bump Go version to 1.19.12.

IMPROVEMENTS:

* core: Log rollback manager failures during unmount, remount to prevent replication failures on secondary clusters. [[GH-22235](https://github.com/hashicorp/vault/pull/22235)]
* replication (enterprise): Make reindex less disruptive by allowing writes during the flush phase.
* storage/raft: Cap the minimum dead_server_last_contact_threshold to 1m. [[GH-22040](https://github.com/hashicorp/vault/pull/22040)]
* ui: enables create and update KV secret workflow when control group present [[GH-22471](https://github.com/hashicorp/vault/pull/22471)]

BUG FIXES:

* api: Fix breakage with UNIX domain socket addresses introduced by newest Go versions as a security fix. [[GH-22523](https://github.com/hashicorp/vault/pull/22523)]
* core (enterprise): Remove MFA Configuration for namespace when deleting namespace
* core/quotas (enterprise): Fix a case where we were applying login roles to lease count quotas in a non-login context.
Also fix a related potential deadlock. [[GH-21110](https://github.com/hashicorp/vault/pull/21110)]
* core:  Remove "expiration manager is nil on tokenstore" error log for unauth requests on DR secondary as they do not have expiration manager. [[GH-22137](https://github.com/hashicorp/vault/pull/22137)]
* core: Fix readonly errors that could occur while loading mounts/auths during unseal [[GH-22362](https://github.com/hashicorp/vault/pull/22362)]
* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-21470](https://github.com/hashicorp/vault/pull/21470)]
* expiration: Fix a deadlock that could occur when a revocation failure happens while restoring leases on startup. [[GH-22374](https://github.com/hashicorp/vault/pull/22374)]
* license: Add autoloaded license path to the cache exempt list. This is to ensure the license changes on the active node is observed on the perfStandby node. [[GH-22363](https://github.com/hashicorp/vault/pull/22363)]
* replication (enterprise): Fix bug sync invalidate CoreReplicatedClusterInfoPath
* replication (enterprise): Fixing a bug by which the atomicity of a merkle diff result could be affected. This means it could be a source of a merkle-diff & sync process failing to switch into stream-wal mode afterwards.
* sdk/ldaputil: Properly escape user filters when using UPN domains
sdk/ldaputil: use EscapeLDAPValue implementation from cap/ldap [[GH-22249](https://github.com/hashicorp/vault/pull/22249)]
* secrets/ldap: Fix bug causing schema and password_policy to be overwritten in config. [[GH-22332](https://github.com/hashicorp/vault/pull/22332)]
* secrets/transform (enterprise): Tidy operations will be re-scheduled at a minimum of every minute, not a maximum of every minute
* ui: Fix blank page or ghost secret when canceling KV secret create [[GH-22541](https://github.com/hashicorp/vault/pull/22541)]
* ui: fixes `max_versions` default for secret metadata unintentionally overriding kv engine defaults [[GH-22394](https://github.com/hashicorp/vault/pull/22394)]

## 1.12.9
### July 25, 2023

CHANGES:

* secrets/transform (enterprise): Enforce a transformation role's max_ttl setting on encode requests, a warning will be returned if max_ttl was applied.

IMPROVEMENTS:

* core/fips: Add RPM, DEB packages of FIPS 140-2 and HSM+FIPS 140-2 Vault Enterprise.
* replication (enterprise): Avoid logging warning if request is forwarded from a performance standby and not a performance secondary
* secrets/transform (enterprise): Switch to pgx PostgreSQL driver for better timeout handling

BUG FIXES:

* core: Fixed an instance where incorrect route entries would get tainted. We now pre-calculate namespace specific paths to avoid this. [[GH-24170](https://github.com/hashicorp/vault/pull/24170)]
* identity: Remove caseSensitivityKey to prevent errors while loading groups which could result in missing groups in memDB when duplicates are found. [[GH-20965](https://github.com/hashicorp/vault/pull/20965)]
* replication (enterprise): update primary cluster address after DR failover
* secrets/azure: Fix intermittent 401s by preventing performance secondary clusters from rotating root credentials. [[GH-21633](https://github.com/hashicorp/vault/pull/21633)]
* secrets/pki: Prevent deleted issuers from reappearing when migrating from a version 1 bundle to a version 2 bundle (versions including 1.13.0, 1.12.2, and 1.11.6); when managed keys were removed but referenced in the Vault 1.10 legacy CA bundle, this the error: `no managed key found with uuid`. [[GH-21316](https://github.com/hashicorp/vault/pull/21316)]
* secrets/pki: Support setting both maintain_stored_certificate_counts=false and publish_stored_certificate_count_metrics=false explicitly in tidy config. [[GH-20664](https://github.com/hashicorp/vault/pull/20664)]
* secrets/transform (enterprise): Fix nil panic when deleting a template with tokenization transformations present
* secrets/transform (enterprise): Grab shared locks for various read operations, only escalating to write locks if work is required
* serviceregistration: Fix bug where multiple nodes in a secondary cluster could be labelled active after updating the cluster's primary [[GH-21642](https://github.com/hashicorp/vault/pull/21642)]
* ui: Fixed an issue where editing an SSH role would clear `default_critical_options` and `default_extension` if left unchanged. [[GH-21739](https://github.com/hashicorp/vault/pull/21739)]

## 1.12.8
### June 21, 2023
BREAKING CHANGES:

* secrets/pki: Maintaining running count of certificates will be turned off by default.
To re-enable keeping these metrics available on the tidy status endpoint, enable
maintain_stored_certificate_counts on tidy-config, to also publish them to the
metrics consumer, enable publish_stored_certificate_count_metrics . [[GH-18186](https://github.com/hashicorp/vault/pull/18186)]

CHANGES:

* core: Bump Go version to 1.19.10.

FEATURES:

* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* core (enterprise): Add background worker for automatic reporting of billing
information. [[GH-19625](https://github.com/hashicorp/vault/pull/19625)]

IMPROVEMENTS:

* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]

BUG FIXES:

* core (enterprise): Don't delete backend stored data that appears to be filterable
on this secondary if we don't have a corresponding mount entry.
* core/activity: add namespace breakdown for new clients when date range spans multiple months, including the current month. [[GH-18766](https://github.com/hashicorp/vault/pull/18766)]
* core/activity: de-duplicate namespaces when historical and current month data are mixed [[GH-18452](https://github.com/hashicorp/vault/pull/18452)]
* core/activity: fix the end_date returned from the activity log endpoint when partial counts are computed [[GH-17856](https://github.com/hashicorp/vault/pull/17856)]
* core/activity: include mount counts when de-duplicating current and historical month data [[GH-18598](https://github.com/hashicorp/vault/pull/18598)]
* core/activity: report mount paths (rather than mount accessors) in current month activity log counts and include deleted mount paths in precomputed queries. [[GH-18916](https://github.com/hashicorp/vault/pull/18916)]
* core/activity: return partial month counts when querying a historical date range and no historical data exists. [[GH-17935](https://github.com/hashicorp/vault/pull/17935)]
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry.  Fix wrong lock used in ListAuths link meta interface implementation. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs
* storage/raft: Fix race where new follower joining can get pruned by dead server cleanup. [[GH-20986](https://github.com/hashicorp/vault/pull/20986)]

## 1.12.7
### June 08, 2023

CHANGES:

* core: Bump Go version to 1.19.9.
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]

IMPROVEMENTS:

* audit: add a `mount_point` field to audit requests and response entries [[GH-20411](https://github.com/hashicorp/vault/pull/20411)]
* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* core: include namespace path in granting_policies block of audit log
* openapi: Fix generated types for duration strings [[GH-20841](https://github.com/hashicorp/vault/pull/20841)]
* sdk/framework: Fix non-deterministic ordering of 'required' fields in OpenAPI spec [[GH-20881](https://github.com/hashicorp/vault/pull/20881)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]

BUG FIXES:

* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* cli: CLI should take days as a unit of time for ttl like flags [[GH-20477](https://github.com/hashicorp/vault/pull/20477)]
* cli: disable printing flags warnings messages for the ssh command [[GH-20502](https://github.com/hashicorp/vault/pull/20502)]
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core: prevent panic on login after namespace is deleted that had mfa enforcement [[GH-20375](https://github.com/hashicorp/vault/pull/20375)]
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation
* secrets/transit: Fix export of HMAC-only key, correctly exporting the key used for sign operations. For consumers of the previously incorrect key, use the plaintext export to retrieve these incorrect keys and import them as new versions.
secrets/transit: Fix bug related to shorter dedicated HMAC key sizing.
sdk/helper/keysutil: New HMAC type policies will have HMACKey equal to Key and be copied over on import. [[GH-20864](https://github.com/hashicorp/vault/pull/20864)]
* ui: Fixes issue unsealing cluster for seal types other than shamir [[GH-20897](https://github.com/hashicorp/vault/pull/20897)]

## 1.12.6
### April 26, 2023

CHANGES:

* core: Bump Go version to 1.19.8.

IMPROVEMENTS:

* cli/namespace: Add detailed flag to output additional namespace information
such as namespace IDs and custom metadata. [[GH-20243](https://github.com/hashicorp/vault/pull/20243)]
* core/activity: add an endpoint to write test activity log data, guarded by a build flag [[GH-20019](https://github.com/hashicorp/vault/pull/20019)]
* core: Add a `raft` sub-field to the `storage` and `ha_storage` details provided by the
`/sys/config/state/sanitized` endpoint in order to include the `max_entry_size`. [[GH-20044](https://github.com/hashicorp/vault/pull/20044)]
* sdk/ldaputil: added `connection_timeout` to tune connection timeout duration
for all LDAP plugins. [[GH-20144](https://github.com/hashicorp/vault/pull/20144)]
* secrets/pki: Decrease size and improve compatibility of OCSP responses by removing issuer certificate. [[GH-20201](https://github.com/hashicorp/vault/pull/20201)]

BUG FIXES:

* auth/ldap: Add max_page_size configurable to LDAP configuration [[GH-19032](https://github.com/hashicorp/vault/pull/19032)]
* command/server: Fix incorrect paths in generated config for `-dev-tls` flag on Windows [[GH-20257](https://github.com/hashicorp/vault/pull/20257)]
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Fix regression breaking non-raft clusters whose nodes share the same cluster_addr/api_addr. [[GH-19721](https://github.com/hashicorp/vault/pull/19721)]
* helper/random: Fix race condition in string generator helper [[GH-19875](https://github.com/hashicorp/vault/pull/19875)]
* kmip (enterprise): Fix a problem decrypting with keys that have no Process Start Date attribute.
* openapi: Fix many incorrect details in generated API spec, by using better techniques to parse path regexps [[GH-18554](https://github.com/hashicorp/vault/pull/18554)]
* pki: Fix automatically turning off CRL signing on upgrade to Vault >= 1.12, if CA Key Usage disallows it [[GH-20220](https://github.com/hashicorp/vault/pull/20220)]
* replication (enterprise): Fix a caching issue when replicating filtered data to
a performance secondary. This resulted in the data being set to nil in the cache
and a "invalid value" error being returned from the API.
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* secrets/pki: Fix patching of leaf_not_after_behavior on issuers. [[GH-20341](https://github.com/hashicorp/vault/pull/20341)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* ui: Fix OIDC provider logo showing when domain doesn't match [[GH-20263](https://github.com/hashicorp/vault/pull/20263)]
* ui: Fix bad link to namespace when namespace name includes `.` [[GH-19799](https://github.com/hashicorp/vault/pull/19799)]
* ui: fixes browser console formatting for help command output [[GH-20064](https://github.com/hashicorp/vault/pull/20064)]
* ui: remove use of htmlSafe except when first sanitized [[GH-20235](https://github.com/hashicorp/vault/pull/20235)]

## 1.12.5
### March 29, 2023

SECURITY:

* storage/mssql: When using Vault’s community-supported Microsoft SQL (MSSQL) database storage backend, a privileged attacker with the ability to write arbitrary data to Vault’s configuration may be able to perform arbitrary SQL commands on the underlying database server through Vault. This vulnerability, CVE-2023-0620, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-12](https://discuss.hashicorp.com/t/hcsec-2023-12-vault-s-microsoft-sql-database-storage-backend-vulnerable-to-sql-injection-via-configuration-file/52080)]
* secrets/pki: Vault’s PKI mount issuer endpoints did not correctly authorize access to remove an issuer or modify issuer metadata, potentially resulting in denial of service of the PKI mount. This bug did not affect public or private key material, trust chains or certificate issuance. This vulnerability, CVE-2023-0665, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-11](https://discuss.hashicorp.com/t/hcsec-2023-11-vault-s-pki-issuer-endpoint-did-not-correctly-authorize-access-to-issuer-metadata/52079)]
* core: HashiCorp Vault’s implementation of Shamir’s secret sharing used precomputed table lookups, and was vulnerable to cache-timing attacks. An attacker with access to, and the ability to observe a large number of unseal operations on the host through a side channel may reduce the search space of a brute force effort to recover the Shamir shares. This vulnerability, CVE-2023-25000, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-10](https://discuss.hashicorp.com/t/hcsec-2023-10-vault-vulnerable-to-cache-timing-attacks-during-seal-and-unseal-operations/52078)]

IMPROVEMENTS:

* auth/github: Allow for an optional Github auth token environment variable to make authenticated requests when fetching org id
website/docs: Add docs for `VAULT_AUTH_CONFIG_GITHUB_TOKEN` environment variable when writing Github config [[GH-19244](https://github.com/hashicorp/vault/pull/19244)]
* core: Allow overriding gRPC connect timeout via VAULT_GRPC_MIN_CONNECT_TIMEOUT. This is an env var rather than a config setting because we don't expect this to ever be needed.  It's being added as a last-ditch
option in case all else fails for some replication issues we may not have fully reproduced. [[GH-19676](https://github.com/hashicorp/vault/pull/19676)]
* core: validate name identifiers in mssql physical storage backend prior use [[GH-19591](https://github.com/hashicorp/vault/pull/19591)]

BUG FIXES:

* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* core (enterprise): Attempt to reconnect to a PKCS#11 HSM if we retrieve a CKR_FUNCTION_FAILED error.
* core: Fixed issue with remounting mounts that have a non-trailing space in the 'to' or 'from' paths. [[GH-19585](https://github.com/hashicorp/vault/pull/19585)]
* kmip (enterprise): Do not require attribute Cryptographic Usage Mask when registering Secret Data managed objects.
* kmip (enterprise): Fix a problem forwarding some requests to the active node.
* openapi: Fix logic for labeling unauthenticated/sudo paths. [[GH-19600](https://github.com/hashicorp/vault/pull/19600)]
* secrets/ldap: Invalidates WAL entry for static role if `password_policy` has changed. [[GH-19641](https://github.com/hashicorp/vault/pull/19641)]
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* ui: fixes issue navigating back a level using the breadcrumb from secret metadata view [[GH-19703](https://github.com/hashicorp/vault/pull/19703)]
* ui: pass encodeBase64 param to HMAC transit-key-actions. [[GH-19429](https://github.com/hashicorp/vault/pull/19429)]
* ui: use URLSearchParams interface to capture namespace param from SSOs (ex. ADFS) with decoded state param in callback url [[GH-19460](https://github.com/hashicorp/vault/pull/19460)]

## 1.12.4
### March 01, 2023

SECURITY:
* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* core: Bump Go version to 1.19.6.

IMPROVEMENTS:

* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]
* ui: remove wizard [[GH-19220](https://github.com/hashicorp/vault/pull/19220)]

BUG FIXES:

* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* license (enterprise): Fix bug where license would update even if the license didn't change.
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* secrets/ad: Fix bug where config couldn't be updated unless binddn/bindpass were included in the update. [[GH-18207](https://github.com/hashicorp/vault/pull/18207)]
* secrets/pki: Revert fix for PR [18938](https://github.com/hashicorp/vault/pull/18938) [[GH-19037](https://github.com/hashicorp/vault/pull/19037)]
* server/config:  Use file.Stat when checking file permissions when VAULT_ENABLE_FILE_PERMISSIONS_CHECK is enabled [[GH-19311](https://github.com/hashicorp/vault/pull/19311)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]
* ui: fixes reliance on secure context (https) by removing methods using the Crypto interface [[GH-19410](https://github.com/hashicorp/vault/pull/19410)]
* ui: show Get credentials button for static roles detail page when a user has the proper permissions. [[GH-19190](https://github.com/hashicorp/vault/pull/19190)]

## 1.12.3
### February 6, 2023

CHANGES:

* core: Bump Go version to 1.19.4.

IMPROVEMENTS:

* audit: Include stack trace when audit logging recovers from a panic. [[GH-18121](https://github.com/hashicorp/vault/pull/18121)]
* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* core: Add read support to `sys/loggers` and `sys/loggers/:name` endpoints [[GH-17979](https://github.com/hashicorp/vault/pull/17979)]
* plugins: Let Vault unseal and mount deprecated builtin plugins in a
deactivated state if this is not the first unseal after an upgrade. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* secrets/kv: new KVv2 mounts and KVv1 mounts without any keys will upgrade synchronously, allowing for instant use [[GH-17406](https://github.com/hashicorp/vault/pull/17406)]
* storage/raft: add additional raft metrics relating to applied index and heartbeating; also ensure OSS standbys emit periodic metrics. [[GH-12166](https://github.com/hashicorp/vault/pull/12166)]
* ui: Added JWT authentication warning message about blocked pop-up windows and web browser settings. [[GH-18787](https://github.com/hashicorp/vault/pull/18787)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]

BUG FIXES:

* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/cert: Address a race condition accessing the loaded crls without a lock [[GH-18945](https://github.com/hashicorp/vault/pull/18945)]
* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#173](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/173)] [[GH-18716](https://github.com/hashicorp/vault/pull/18716)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* command/namespace: Fix vault cli namespace patch examples in help text. [[GH-18143](https://github.com/hashicorp/vault/pull/18143)]
* core (enterprise): Fix a race condition resulting in login errors to PKCS#11 modules under high concurrency.
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core/seal: Fix regression handling of the key_id parameter in seal configuration HCL. [[GH-17612](https://github.com/hashicorp/vault/pull/17612)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* expiration: Prevent panics on perf standbys when an irrevocable lease gets deleted. [[GH-18401](https://github.com/hashicorp/vault/pull/18401)]
* kmip (enterprise): Fix Destroy operation response that omitted Unique Identifier on some batched responses.
* kmip (enterprise): Fix Locate operation response incompatibility with clients using KMIP versions prior to 1.3.
* kmip (enterprise): Fix Query operation response that omitted streaming capability and supported profiles.
* licensing (enterprise): update autoloaded license cache after reload
* plugins: Allow running external plugins which override deprecated builtins. [[GH-17879](https://github.com/hashicorp/vault/pull/17879)]
* plugins: Listing all plugins while audit logging is enabled will no longer result in an internal server error. [[GH-18173](https://github.com/hashicorp/vault/pull/18173)]
* plugins: Skip loading but still mount data associated with missing plugins on unseal. [[GH-18189](https://github.com/hashicorp/vault/pull/18189)]
* sdk: Don't panic if system view or storage methods called during plugin setup. [[GH-18210](https://github.com/hashicorp/vault/pull/18210)]
* secrets/pki: Address nil panic when an empty POST request is sent to the OCSP handler [[GH-18184](https://github.com/hashicorp/vault/pull/18184)]
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/pki: OCSP GET request parameter was not being URL unescaped before processing. [[GH-18938](https://github.com/hashicorp/vault/pull/18938)]
* secrets/pki: fix race between tidy's cert counting and tidy status reporting. [[GH-18899](https://github.com/hashicorp/vault/pull/18899)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* secrets/transit: Honor `partial_success_response_code` on decryption failures. [[GH-18310](https://github.com/hashicorp/vault/pull/18310)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* ui: cleanup unsaved auth method ember data record when navigating away from mount backend form [[GH-18651](https://github.com/hashicorp/vault/pull/18651)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]
## 1.12.2
### November 30, 2022

CHANGES:

* core: Bump Go version to 1.19.3.
* plugins: Mounts can no longer be pinned to a specific _builtin_ version. Mounts previously pinned to a specific builtin version will now automatically upgrade to the latest builtin version, and may now be overridden if an unversioned plugin of the same name and type is registered. Mounts using plugin versions without `builtin` in their metadata remain unaffected. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]

IMPROVEMENTS:

* secrets/pki: Allow issuer creation, import to change default issuer via `default_follows_latest_issuer`. [[GH-17824](https://github.com/hashicorp/vault/pull/17824)]
* storage/raft: Add `retry_join_as_non_voter` config option. [[GH-18030](https://github.com/hashicorp/vault/pull/18030)]

BUG FIXES:

* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* core (enterprise): Supported storage check in `vault server` command will no longer prevent startup. Instead, a warning will be logged if configured to use storage backend other than `raft` or `consul`.
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: fix a start up race condition where performance standbys could go into a
  mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* plugins: Only report deprecation status for builtin plugins. [[GH-17816](https://github.com/hashicorp/vault/pull/17816)]
* plugins: Vault upgrades will no longer fail if a mount has been created using an explicit builtin plugin version. [[GH-18051](https://github.com/hashicorp/vault/pull/18051)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/azure: add WAL to clean up role assignments if errors occur [[GH-18086](https://github.com/hashicorp/vault/pull/18086)]
* secrets/gcp: Fixes duplicate service account key for rotate root on standby or secondary [[GH-18111](https://github.com/hashicorp/vault/pull/18111)]
* secrets/pki: Fix upgrade of missing expiry, delta_rebuild_interval by setting them to the default. [[GH-17693](https://github.com/hashicorp/vault/pull/17693)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]

## 1.12.1
### November 2, 2022

IMPROVEMENTS:

* api: Support VAULT_DISABLE_REDIRECTS environment variable (and --disable-redirects flag) to disable default client behavior and prevent the client following any redirection responses. [[GH-17352](https://github.com/hashicorp/vault/pull/17352)]
* database/snowflake: Allow parallel requests to Snowflake [[GH-17593](https://github.com/hashicorp/vault/pull/17593)]
* plugins: Add plugin version information to key plugin lifecycle log lines. [[GH-17430](https://github.com/hashicorp/vault/pull/17430)]
* sdk/ldap: Added support for paging when searching for groups using group filters [[GH-17640](https://github.com/hashicorp/vault/pull/17640)]

BUG FIXES:

* cli: Remove empty table heading for `vault secrets list -detailed` output. [[GH-17577](https://github.com/hashicorp/vault/pull/17577)]
* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): Switch to using hash length as PSS Salt length within the test/sign api for better PKCS#11 compatibility
* core: Fix panic caused in Vault Agent when rendering certificate templates [[GH-17419](https://github.com/hashicorp/vault/pull/17419)]
* core: Fixes spurious warnings being emitted relating to "unknown or unsupported fields" for JSON config [[GH-17660](https://github.com/hashicorp/vault/pull/17660)]
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* kmip (enterprise): Fix a problem in the handling of attributes that caused Import operations to fail.
* kmip (enterprise): Fix selection of Cryptographic Parameters for Encrypt/Decrypt operations.
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]

## 1.12.0
### October 13, 2022

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

CHANGES:

* api: Exclusively use `GET /sys/plugins/catalog` endpoint for listing plugins, and add `details` field to list responses. [[GH-17347](https://github.com/hashicorp/vault/pull/17347)]
* auth: `GET /sys/auth/:name` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* auth: `GET /sys/auth` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* auth: `POST /sys/auth/:type` endpoint response contains a warning for `Deprecated` auth methods. [[GH-17058](https://github.com/hashicorp/vault/pull/17058)]
* auth: `auth enable` returns an error and `POST /sys/auth/:type` endpoint reports an error for `Pending Removal` auth methods. [[GH-17005](https://github.com/hashicorp/vault/pull/17005)]
* core/entities: Fixed stranding of aliases upon entity merge, and require explicit selection of which aliases should be kept when some must be deleted [[GH-16539](https://github.com/hashicorp/vault/pull/16539)]
* core: Bump Go version to 1.19.2.
* core: Validate input parameters for vault operator init command. Vault 1.12 CLI version is needed to run operator init now. [[GH-16379](https://github.com/hashicorp/vault/pull/16379)]
* identity: a request to `/identity/group` that includes `member_group_ids` that contains a cycle will now be responded to with a 400 rather than 500 [[GH-15912](https://github.com/hashicorp/vault/pull/15912)]
* licensing (enterprise): Terminated licenses will no longer result in shutdown. Instead, upgrades will not be allowed if the license expiration time is before the build date of the binary.
* plugins: Add plugin version to auth register, list, and mount table [[GH-16856](https://github.com/hashicorp/vault/pull/16856)]
* plugins: `GET /sys/plugins/catalog/:type/:name` endpoint contains deprecation status for builtin plugins. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* plugins: `GET /sys/plugins/catalog/:type/:name` endpoint now returns an additional `version` field in the response data. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* plugins: `GET /sys/plugins/catalog/` endpoint contains deprecation status in `detailed` list. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* plugins: `GET /sys/plugins/catalog` endpoint now returns an additional `detailed` field in the response data with a list of additional plugin metadata. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* plugins: `plugin info` displays deprecation status for builtin plugins. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* plugins: `plugin list` now accepts a `-detailed` flag, which display deprecation status and version info. [[GH-17077](https://github.com/hashicorp/vault/pull/17077)]
* secrets/azure: Removed deprecated AAD graph API support from the secrets engine. [[GH-17180](https://github.com/hashicorp/vault/pull/17180)]
* secrets: All database-specific (standalone DB) secrets engines are now marked `Pending Removal`. [[GH-17038](https://github.com/hashicorp/vault/pull/17038)]
* secrets: `GET /sys/mounts/:name` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* secrets: `GET /sys/mounts` endpoint now returns an additional `deprecation_status` field in the response data for builtins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* secrets: `POST /sys/mounts/:type` endpoint response contains a warning for `Deprecated` secrets engines. [[GH-17058](https://github.com/hashicorp/vault/pull/17058)]
* secrets: `secrets enable` returns an error and `POST /sys/mount/:type` endpoint reports an error for `Pending Removal` secrets engines. [[GH-17005](https://github.com/hashicorp/vault/pull/17005)]

FEATURES:

* **GCP Cloud KMS support for managed keys**: Managed keys now support using GCP Cloud KMS keys
* **LDAP Secrets Engine**: Adds the `ldap` secrets engine with service account check-out functionality for all supported schemas. [[GH-17152](https://github.com/hashicorp/vault/pull/17152)]
* **OCSP Responder**: PKI mounts now have an OCSP responder that implements a subset of RFC6960, answering single serial number OCSP requests for a specific cluster's revoked certificates in a mount. [[GH-16723](https://github.com/hashicorp/vault/pull/16723)]
* **Redis DB Engine**: Adding the new Redis database engine that supports the generation of static and dynamic user roles and root credential rotation on a stand alone Redis server. [[GH-17070](https://github.com/hashicorp/vault/pull/17070)]
* **Redis ElastiCache DB Plugin**: Added Redis ElastiCache as a built-in plugin. [[GH-17075](https://github.com/hashicorp/vault/pull/17075)]
* **Secrets/auth plugin multiplexing**: manage multiple plugin configurations with a single plugin process [[GH-14946](https://github.com/hashicorp/vault/pull/14946)]
* **Transform Key Import (BYOK)**: The transform secrets engine now supports importing keys for tokenization and FPE transformations
* HCP (enterprise): Adding foundational support for self-managed vault nodes to securely communicate with [HashiCorp Cloud Platform](https://cloud.hashicorp.com) as an opt-in feature
* ui: UI support for Okta Number Challenge. [[GH-15998](https://github.com/hashicorp/vault/pull/15998)]
* **Plugin Versioning**: Vault supports registering, managing, and running plugins with semantic versions specified.

IMPROVEMENTS:

* core/managed-keys (enterprise): Allow operators to specify PSS signatures and/or hash algorithm for the test/sign api
* activity (enterprise): Added new clients unit tests to test accuracy of estimates
* agent/auto-auth: Add `exit_on_err` which when set to true, will cause Agent to exit if any errors are encountered during authentication. [[GH-17091](https://github.com/hashicorp/vault/pull/17091)]
* agent: Added `disable_idle_connections` configuration to disable leaving idle connections open in auto-auth, caching and templating. [[GH-15986](https://github.com/hashicorp/vault/pull/15986)]
* agent: Added `disable_keep_alives` configuration to disable keep alives in auto-auth, caching and templating. [[GH-16479](https://github.com/hashicorp/vault/pull/16479)]
* agent: JWT auto auth now supports a `remove_jwt_after_reading` config option which defaults to true. [[GH-11969](https://github.com/hashicorp/vault/pull/11969)]
* agent: Send notifications to systemd on start and stop. [[GH-9802](https://github.com/hashicorp/vault/pull/9802)]
* api/mfa: Add namespace path to the MFA read/list endpoint [[GH-16911](https://github.com/hashicorp/vault/pull/16911)]
* api: Add a sentinel error for missing KV secrets [[GH-16699](https://github.com/hashicorp/vault/pull/16699)]
* auth/alicloud: Enables AliCloud roles to be compatible with Vault's role based quotas. [[GH-17251](https://github.com/hashicorp/vault/pull/17251)]
* auth/approle: SecretIDs can now be generated with an per-request specified TTL and num_uses.
When either the ttl and num_uses fields are not specified, the role's configuration is used. [[GH-14474](https://github.com/hashicorp/vault/pull/14474)]
* auth/aws: PKCS7 signatures will now use SHA256 by default in prep for Go 1.18 [[GH-16455](https://github.com/hashicorp/vault/pull/16455)]
* auth/azure: Enables Azure roles to be compatible with Vault's role based quotas. [[GH-17194](https://github.com/hashicorp/vault/pull/17194)]
* auth/cert: Add metadata to identity-alias [[GH-14751](https://github.com/hashicorp/vault/pull/14751)]
* auth/cert: Operators can now specify a CRL distribution point URL, in which case the cert auth engine will fetch and use the CRL from that location rather than needing to push CRLs directly to auth/cert. [[GH-17136](https://github.com/hashicorp/vault/pull/17136)]
* auth/cf: Enables CF roles to be compatible with Vault's role based quotas. [[GH-17196](https://github.com/hashicorp/vault/pull/17196)]
* auth/gcp: Add support for GCE regional instance groups [[GH-16435](https://github.com/hashicorp/vault/pull/16435)]
* auth/gcp: Updates dependencies: `google.golang.org/api@v0.83.0`, `github.com/hashicorp/go-gcp-common@v0.8.0`. [[GH-17160](https://github.com/hashicorp/vault/pull/17160)]
* auth/jwt: Adds support for Microsoft US Gov L4 to the Azure provider for groups fetching. [[GH-16525](https://github.com/hashicorp/vault/pull/16525)]
* auth/jwt: Improves detection of Windows Subsystem for Linux (WSL) for CLI-based logins. [[GH-16525](https://github.com/hashicorp/vault/pull/16525)]
* auth/kerberos: add `add_group_aliases` config to include LDAP groups in Vault group aliases [[GH-16890](https://github.com/hashicorp/vault/pull/16890)]
* auth/kerberos: add `remove_instance_name` parameter to the login CLI and the Kerberos config in Vault. This removes any instance names found in the keytab service principal name. [[GH-16594](https://github.com/hashicorp/vault/pull/16594)]
* auth/kubernetes: Role resolution for K8S Auth [[GH-156](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/156)] [[GH-17161](https://github.com/hashicorp/vault/pull/17161)]
* auth/oci: Add support for role resolution. [[GH-17212](https://github.com/hashicorp/vault/pull/17212)]
* auth/oidc: Adds support for group membership parsing when using SecureAuth as an OIDC provider. [[GH-16274](https://github.com/hashicorp/vault/pull/16274)]
* cli: CLI commands will print a warning if flags will be ignored because they are passed after positional arguments. [[GH-16441](https://github.com/hashicorp/vault/pull/16441)]
* cli: `auth` and `secrets` list `-detailed` commands now show Deprecation Status for builtin plugins. [[GH-16849](https://github.com/hashicorp/vault/pull/16849)]
* cli: `vault plugin list` now has a `details` field in JSON format, and version and type information in table format. [[GH-17347](https://github.com/hashicorp/vault/pull/17347)]
* command/audit: Improve missing type error message [[GH-16409](https://github.com/hashicorp/vault/pull/16409)]
* command/server: add `-dev-tls` and `-dev-tls-cert-dir` subcommands to create a Vault dev server with generated certificates and private key. [[GH-16421](https://github.com/hashicorp/vault/pull/16421)]
* command: Fix shell completion for KV v2 mounts [[GH-16553](https://github.com/hashicorp/vault/pull/16553)]
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
* core: Activity log goroutine management improvements to allow tests to be more deterministic. [[GH-17028](https://github.com/hashicorp/vault/pull/17028)]
* core: Add `sys/loggers` and `sys/loggers/:name` endpoints to provide ability to modify logging verbosity [[GH-16111](https://github.com/hashicorp/vault/pull/16111)]
* core: Handle and log deprecated builtin mounts. Introduces `VAULT_ALLOW_PENDING_REMOVAL_MOUNTS` to override shutdown and error when attempting to mount `Pending Removal` builtin plugins. [[GH-17005](https://github.com/hashicorp/vault/pull/17005)]
* core: Limit activity log client count usage by namespaces [[GH-16000](https://github.com/hashicorp/vault/pull/16000)]
* core: Upgrade github.com/hashicorp/raft [[GH-16609](https://github.com/hashicorp/vault/pull/16609)]
* core: remove gox [[GH-16353](https://github.com/hashicorp/vault/pull/16353)]
* docs: Clarify the behaviour of local mounts in the context of DR replication [[GH-16218](https://github.com/hashicorp/vault/pull/16218)]
* identity/oidc: Adds support for detailed listing of clients and providers. [[GH-16567](https://github.com/hashicorp/vault/pull/16567)]
* identity/oidc: Adds the `client_secret_post` token endpoint authentication method. [[GH-16598](https://github.com/hashicorp/vault/pull/16598)]
* identity/oidc: allows filtering the list providers response by an allowed_client_id [[GH-16181](https://github.com/hashicorp/vault/pull/16181)]
* identity: Prevent possibility of data races on entity creation. [[GH-16487](https://github.com/hashicorp/vault/pull/16487)]
* physical/postgresql: pass context to queries to propagate timeouts and cancellations on requests. [[GH-15866](https://github.com/hashicorp/vault/pull/15866)]
* plugins/multiplexing: Added multiplexing support to database plugins if run as external plugins [[GH-16995](https://github.com/hashicorp/vault/pull/16995)]
* plugins: Add Deprecation Status method to builtinregistry. [[GH-16846](https://github.com/hashicorp/vault/pull/16846)]
* plugins: Added environment variable flag to opt-out specific plugins from multiplexing [[GH-16972](https://github.com/hashicorp/vault/pull/16972)]
* plugins: Adding version to plugin GRPC interface [[GH-17088](https://github.com/hashicorp/vault/pull/17088)]
* plugins: Plugin catalog supports registering and managing plugins with semantic version information. [[GH-16688](https://github.com/hashicorp/vault/pull/16688)]
* replication (enterprise): Fix race in merkle sync that can prevent streaming by returning key value matching provided hash if found in log shipper buffer.
* secret/nomad: allow reading CA and client auth certificate from /nomad/config/access [[GH-15809](https://github.com/hashicorp/vault/pull/15809)]
* secret/pki: Add RSA PSS signature support for issuing certificates, signing CRLs [[GH-16519](https://github.com/hashicorp/vault/pull/16519)]
* secret/pki: Add signature_bits to sign-intermediate, sign-verbatim endpoints [[GH-16124](https://github.com/hashicorp/vault/pull/16124)]
* secret/pki: Allow issuing certificates with non-domain, non-email Common Names from roles, sign-verbatim, and as issuers (`cn_validations`). [[GH-15996](https://github.com/hashicorp/vault/pull/15996)]
* secret/pki: Allow specifying SKID for cross-signed issuance from older Vault versions. [[GH-16494](https://github.com/hashicorp/vault/pull/16494)]
* secret/transit: Allow importing Ed25519 keys from PKCS#8 with inner RFC 5915 ECPrivateKey blobs (NSS-wrapped keys). [[GH-15742](https://github.com/hashicorp/vault/pull/15742)]
* secrets/ad: set config default length only if password_policy is missing [[GH-16140](https://github.com/hashicorp/vault/pull/16140)]
* secrets/azure: Adds option to permanently delete AzureAD objects created by Vault. [[GH-17045](https://github.com/hashicorp/vault/pull/17045)]
* secrets/database/hana: Add ability to customize dynamic usernames [[GH-16631](https://github.com/hashicorp/vault/pull/16631)]
* secrets/database/snowflake: Add multiplexing support [[GH-17159](https://github.com/hashicorp/vault/pull/17159)]
* secrets/gcp: Updates dependencies: `google.golang.org/api@v0.83.0`, `github.com/hashicorp/go-gcp-common@v0.8.0`. [[GH-17174](https://github.com/hashicorp/vault/pull/17174)]
* secrets/gcpkms: Update dependencies: google.golang.org/api@v0.83.0. [[GH-17199](https://github.com/hashicorp/vault/pull/17199)]
* secrets/kubernetes: upgrade to v0.2.0 [[GH-17164](https://github.com/hashicorp/vault/pull/17164)]
* secrets/pki/tidy: Add another pair of metrics counting certificates not deleted by the tidy operation. [[GH-16702](https://github.com/hashicorp/vault/pull/16702)]
* secrets/pki: Add a new flag to issue/sign APIs which can filter out root CAs from the returned ca_chain field [[GH-16935](https://github.com/hashicorp/vault/pull/16935)]
* secrets/pki: Add a warning to any successful response when the requested TTL is overwritten by MaxTTL [[GH-17073](https://github.com/hashicorp/vault/pull/17073)]
* secrets/pki: Add ability to cancel tidy operations, control tidy resource usage. [[GH-16958](https://github.com/hashicorp/vault/pull/16958)]
* secrets/pki: Add ability to periodically rebuild CRL before expiry [[GH-16762](https://github.com/hashicorp/vault/pull/16762)]
* secrets/pki: Add ability to periodically run tidy operations to remove expired certificates. [[GH-16900](https://github.com/hashicorp/vault/pull/16900)]
* secrets/pki: Add support for per-issuer Authority Information Access (AIA) URLs [[GH-16563](https://github.com/hashicorp/vault/pull/16563)]
* secrets/pki: Add support to specify signature bits when generating CSRs through intermediate/generate apis [[GH-17388](https://github.com/hashicorp/vault/pull/17388)]
* secrets/pki: Added gauge metrics "secrets.pki.total_revoked_certificates_stored" and "secrets.pki.total_certificates_stored" to track the number of certificates in storage. [[GH-16676](https://github.com/hashicorp/vault/pull/16676)]
* secrets/pki: Allow revocation of certificates with explicitly provided certificate (bring your own certificate / BYOC). [[GH-16564](https://github.com/hashicorp/vault/pull/16564)]
* secrets/pki: Allow revocation via proving possession of certificate's private key [[GH-16566](https://github.com/hashicorp/vault/pull/16566)]
* secrets/pki: Allow tidy to associate revoked certs with their issuers for OCSP performance [[GH-16871](https://github.com/hashicorp/vault/pull/16871)]
* secrets/pki: Honor If-Modified-Since header on CA, CRL fetch; requires passthrough_request_headers modification on the mount point. [[GH-16249](https://github.com/hashicorp/vault/pull/16249)]
* secrets/pki: Improve stability of association of revoked cert with its parent issuer; when an issuer loses crl-signing usage, do not place certs on default issuer's CRL. [[GH-16874](https://github.com/hashicorp/vault/pull/16874)]
* secrets/pki: Support generating delta CRLs for up-to-date CRLs when auto-building is enabled. [[GH-16773](https://github.com/hashicorp/vault/pull/16773)]
* secrets/ssh: Add allowed_domains_template to allow templating of allowed_domains. [[GH-16056](https://github.com/hashicorp/vault/pull/16056)]
* secrets/ssh: Allow additional text along with a template definition in defaultExtension value fields. [[GH-16018](https://github.com/hashicorp/vault/pull/16018)]
* secrets/ssh: Allow the use of Identity templates in the `default_user` field [[GH-16351](https://github.com/hashicorp/vault/pull/16351)]
* secrets/transit: Add a dedicated HMAC key type, which can be used with key import. [[GH-16668](https://github.com/hashicorp/vault/pull/16668)]
* secrets/transit: Added a parameter to encrypt/decrypt batch operations to allow the caller to override the HTTP response code in case of partial user-input failures. [[GH-17118](https://github.com/hashicorp/vault/pull/17118)]
* secrets/transit: Allow configuring the possible salt lengths for RSA PSS signatures. [[GH-16549](https://github.com/hashicorp/vault/pull/16549)]
* ssh: Addition of an endpoint `ssh/issue/:role` to allow the creation of signed key pairs [[GH-15561](https://github.com/hashicorp/vault/pull/15561)]
* storage/cassandra: tuning parameters for clustered environments `connection_timeout`, `initial_connection_timeout`, `simple_retry_policy_retries`. [[GH-10467](https://github.com/hashicorp/vault/pull/10467)]
* storage/gcs: Add documentation explaining how to configure the gcs backend using environment variables instead of options in the configuration stanza [[GH-14455](https://github.com/hashicorp/vault/pull/14455)]
* ui: Changed the tokenBoundCidrs tooltip content to clarify that comma separated values are not accepted in this field. [[GH-15852](https://github.com/hashicorp/vault/pull/15852)]
* ui: Prevents requests to /sys/internal/ui/resultant-acl endpoint when unauthenticated [[GH-17139](https://github.com/hashicorp/vault/pull/17139)]
* ui: Removed deprecated version of core-js 2.6.11 [[GH-15898](https://github.com/hashicorp/vault/pull/15898)]
* ui: Renamed labels under Tools for wrap, lookup, rewrap and unwrap with description. [[GH-16489](https://github.com/hashicorp/vault/pull/16489)]
* ui: Replaces non-inclusive terms [[GH-17116](https://github.com/hashicorp/vault/pull/17116)]
* ui: redirect_to param forwards from auth route when authenticated [[GH-16821](https://github.com/hashicorp/vault/pull/16821)]
* website/docs: API generate-recovery-token documentation. [[GH-16213](https://github.com/hashicorp/vault/pull/16213)]
* website/docs: Add documentation around the expensiveness of making lots of lease count quotas in a short period [[GH-16950](https://github.com/hashicorp/vault/pull/16950)]
* website/docs: Removes mentions of unauthenticated from internal ui resultant-acl doc [[GH-17139](https://github.com/hashicorp/vault/pull/17139)]
* website/docs: Update replication docs to mention Integrated Storage [[GH-16063](https://github.com/hashicorp/vault/pull/16063)]
* website/docs: changed to echo for all string examples instead of (<<<) here-string. [[GH-9081](https://github.com/hashicorp/vault/pull/9081)]

BUG FIXES:

* agent/template: Fix parsing error for the exec stanza [[GH-16231](https://github.com/hashicorp/vault/pull/16231)]
* agent: Agent will now respect `max_retries` retry configuration even when caching is set. [[GH-16970](https://github.com/hashicorp/vault/pull/16970)]
* agent: Update consul-template for pkiCert bug fixes [[GH-16087](https://github.com/hashicorp/vault/pull/16087)]
* api/sys/internal/specs/openapi: support a new "dynamic" query parameter to generate generic mountpaths [[GH-15835](https://github.com/hashicorp/vault/pull/15835)]
* api: Fixed erroneous warnings of unrecognized parameters when unwrapping data. [[GH-16794](https://github.com/hashicorp/vault/pull/16794)]
* api: Fixed issue with internal/ui/mounts and internal/ui/mounts/(?P<path>.+) endpoints where it was not properly handling /auth/ [[GH-15552](https://github.com/hashicorp/vault/pull/15552)]
* api: properly handle switching to/from unix domain socket when changing client address [[GH-11904](https://github.com/hashicorp/vault/pull/11904)]
* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* auth/kerberos: Maintain headers set by the client [[GH-16636](https://github.com/hashicorp/vault/pull/16636)]
* auth/kubernetes: Restore support for JWT signature algorithm ES384 [[GH-160](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/160)] [[GH-17161](https://github.com/hashicorp/vault/pull/17161)]
* auth/token: Fix ignored parameter warnings for valid parameters on token create [[GH-16938](https://github.com/hashicorp/vault/pull/16938)]
* command/debug: fix bug where monitor was not honoring configured duration [[GH-16834](https://github.com/hashicorp/vault/pull/16834)]
* core (enterprise): Fix bug where wrapping token lookup does not work within namespaces. [[GH-15583](https://github.com/hashicorp/vault/pull/15583)]
* core (enterprise): Fix creation of duplicate entities via alias metadata changes on local auth mounts.
* core/auth: Return a 403 instead of a 500 for a malformed SSCT [[GH-16112](https://github.com/hashicorp/vault/pull/16112)]
* core/identity: Replicate member_entity_ids and policies in identity/group across nodes identically [[GH-16088](https://github.com/hashicorp/vault/pull/16088)]
* core/license (enterprise): Always remove stored license and allow unseal to complete when license cleanup fails
* core/managed-keys (enterprise): fix panic when having `cache_disable` true
* core/quotas (enterprise): Fixed issue with improper counting of leases if lease count quota created after leases
* core/quotas: Added globbing functionality on the end of path suffix quota paths [[GH-16386](https://github.com/hashicorp/vault/pull/16386)]
* core/quotas: Fix goroutine leak caused by the seal process not fully cleaning up Rate Limit Quotas. [[GH-17281](https://github.com/hashicorp/vault/pull/17281)]
* core/replication (enterprise): Don't flush merkle tree pages to disk after losing active duty
* core/seal: Fix possible keyring truncation when using the file backend. [[GH-15946](https://github.com/hashicorp/vault/pull/15946)]
* core: Fix panic when the plugin catalog returns neither a plugin nor an error. [[GH-17204](https://github.com/hashicorp/vault/pull/17204)]
* core: Fixes parsing boolean values for ha_storage backends in config [[GH-15900](https://github.com/hashicorp/vault/pull/15900)]
* core: Increase the allowed concurrent gRPC streams over the cluster port. [[GH-16327](https://github.com/hashicorp/vault/pull/16327)]
* core: Prevent two or more DR failovers from invalidating SSCT tokens generated on the previous primaries. [[GH-16956](https://github.com/hashicorp/vault/pull/16956)]
* database: Invalidate queue should cancel context first to avoid deadlock [[GH-15933](https://github.com/hashicorp/vault/pull/15933)]
* debug: Fix panic when capturing debug bundle on Windows [[GH-14399](https://github.com/hashicorp/vault/pull/14399)]
* debug: Remove extra empty lines from vault.log when debug command is run [[GH-16714](https://github.com/hashicorp/vault/pull/16714)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* identity/oidc: Adds `claims_supported` to discovery document. [[GH-16992](https://github.com/hashicorp/vault/pull/16992)]
* identity/oidc: Change the `state` parameter of the Authorization Endpoint to optional. [[GH-16599](https://github.com/hashicorp/vault/pull/16599)]
* identity/oidc: Detect invalid `redirect_uri` values sooner in validation of the Authorization Endpoint. [[GH-16601](https://github.com/hashicorp/vault/pull/16601)]
* identity/oidc: Fixes validation of the `request` and `request_uri` parameters. [[GH-16600](https://github.com/hashicorp/vault/pull/16600)]
* openapi: Fixed issue where information about /auth/token endpoints was not present with explicit policy permissions [[GH-15552](https://github.com/hashicorp/vault/pull/15552)]
* plugin/multiplexing: Fix panic when id doesn't exist in connection map [[GH-16094](https://github.com/hashicorp/vault/pull/16094)]
* plugin/secrets/auth: Fix a bug with aliased backends such as aws-ec2 or generic [[GH-16673](https://github.com/hashicorp/vault/pull/16673)]
* plugins: Corrected the path to check permissions on when the registered plugin name does not match the plugin binary's filename. [[GH-17340](https://github.com/hashicorp/vault/pull/17340)]
* quotas/lease-count: Fix lease-count quotas on mounts not properly being enforced when the lease generating request is a read [[GH-15735](https://github.com/hashicorp/vault/pull/15735)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* replication (enterprise): Fix data race in saveCheckpoint.
* replication (enterprise): Fix possible data race during merkle diff/sync
* secret/pki: Do not fail validation with a legacy key_bits default value and key_type=any when signing CSRs [[GH-16246](https://github.com/hashicorp/vault/pull/16246)]
* secrets/database: Fix a bug where the secret engine would queue up a lot of WAL deletes during startup. [[GH-16686](https://github.com/hashicorp/vault/pull/16686)]
* secrets/gcp: Fixes duplicate static account key creation from performance secondary clusters. [[GH-16534](https://github.com/hashicorp/vault/pull/16534)]
* secrets/kv: Fix `kv get` issue preventing the ability to read a secret when providing a leading slash [[GH-16443](https://github.com/hashicorp/vault/pull/16443)]
* secrets/pki: Allow import of issuers without CRLSign KeyUsage; prohibit setting crl-signing usage on such issuers [[GH-16865](https://github.com/hashicorp/vault/pull/16865)]
* secrets/pki: Do not ignore provided signature bits value when signing intermediate and leaf certificates with a managed key [[GH-17328](https://github.com/hashicorp/vault/pull/17328)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17385](https://github.com/hashicorp/vault/pull/17385)]
* secrets/pki: Fix migration to properly handle mounts that contain only keys, no certificates [[GH-16813](https://github.com/hashicorp/vault/pull/16813)]
* secrets/pki: Ignore EC PARAMETER PEM blocks during issuer import (/config/ca, /issuers/import/*, and /intermediate/set-signed) [[GH-16721](https://github.com/hashicorp/vault/pull/16721)]
* secrets/pki: LIST issuers endpoint is now unauthenticated. [[GH-16830](https://github.com/hashicorp/vault/pull/16830)]
* secrets/transform (enterprise): Fix an issue loading tokenization transform configuration after a specific sequence of reconfigurations.
* secrets/transform (enterprise): Fix persistence problem with tokenization store credentials.
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* storage/raft (enterprise): Prevent unauthenticated voter status change with rejoin [[GH-16324](https://github.com/hashicorp/vault/pull/16324)]
* storage/raft: Fix retry_join initialization failure [[GH-16550](https://github.com/hashicorp/vault/pull/16550)]
* storage/raft: Nodes no longer get demoted to nonvoter if we don't know their version due to missing heartbeats. [[GH-17019](https://github.com/hashicorp/vault/pull/17019)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fix OIDC callback to accept namespace flag in different formats [[GH-16886](https://github.com/hashicorp/vault/pull/16886)]
* ui: Fix info tooltip submitting form [[GH-16659](https://github.com/hashicorp/vault/pull/16659)]
* ui: Fix issue logging in with JWT auth method [[GH-16466](https://github.com/hashicorp/vault/pull/16466)]
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]
* ui: Fix naming of permitted_dns_domains form parameter on CA creation (root generation and sign intermediate). [[GH-16739](https://github.com/hashicorp/vault/pull/16739)]
* ui: Fixed bug where red spellcheck underline appears in sensitive/secret kv values when it should not appear [[GH-15681](https://github.com/hashicorp/vault/pull/15681)]
* ui: Fixes secret version and status menu links transitioning to auth screen [[GH-16983](https://github.com/hashicorp/vault/pull/16983)]
* ui: OIDC login type uses localStorage instead of sessionStorage [[GH-16170](https://github.com/hashicorp/vault/pull/16170)]
* vault: Fix a bug where duplicate policies could be added to an identity group. [[GH-15638](https://github.com/hashicorp/vault/pull/15638)]

## 1.11.12
### June 21, 2023

CHANGES:

* core: Bump Go version to 1.19.10.
* licensing (enterprise): Terminated licenses will no longer result in shutdown. Instead, upgrades
will not be allowed if the license termination time is before the build date of the binary.

FEATURES:

* **Automated License Utilization Reporting**: Added automated license
utilization reporting, which sends minimal product-license [metering
data](https://developer.hashicorp.com/vault/docs/enterprise/license/utilization-reporting)
to HashiCorp without requiring you to manually collect and report them.
* core (enterprise): Add background worker for automatic reporting of billing
information. [[GH-19625](https://github.com/hashicorp/vault/pull/19625)]

IMPROVEMENTS:

* api: GET ... /sys/internal/counters/activity?current_billing_period=true now
results in a response which contains the full billing period [[GH-20694](https://github.com/hashicorp/vault/pull/20694)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`minimum_retention_months`. [[GH-20150](https://github.com/hashicorp/vault/pull/20150)]
* api: `/sys/internal/counters/config` endpoint now contains read-only
`reporting_enabled` and `billing_start_timestamp` fields. [[GH-20086](https://github.com/hashicorp/vault/pull/20086)]
* core (enterprise): add configuration for license reporting [[GH-19891](https://github.com/hashicorp/vault/pull/19891)]
* core (enterprise): license updates trigger a reload of reporting and the activity log [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): support reloading configuration for automated reporting via SIGHUP [[GH-20680](https://github.com/hashicorp/vault/pull/20680)]
* core (enterprise): vault server command now allows for opt-out of automated
reporting via the `OPTOUT_LICENSE_REPORTING` environment variable. [[GH-3939](https://github.com/hashicorp/vault/pull/3939)]
* core/activity: error when attempting to update retention configuration below the minimum [[GH-20078](https://github.com/hashicorp/vault/pull/20078)]
* core/activity: generate hyperloglogs containing clientIds for each month during precomputation [[GH-16146](https://github.com/hashicorp/vault/pull/16146)]
* core/activity: refactor activity log api to reuse partial api functions in activity endpoint when current month is specified [[GH-16162](https://github.com/hashicorp/vault/pull/16162)]
* core/activity: refactor the activity log's generation of precomputed queries [[GH-20073](https://github.com/hashicorp/vault/pull/20073)]
* core/activity: use monthly hyperloglogs to calculate new clients approximation for current month [[GH-16184](https://github.com/hashicorp/vault/pull/16184)]
* core: Activity log goroutine management improvements to allow tests to be more deterministic. [[GH-17028](https://github.com/hashicorp/vault/pull/17028)]
* core: Limit activity log client count usage by namespaces [[GH-16000](https://github.com/hashicorp/vault/pull/16000)]
* storage/raft: add additional raft metrics relating to applied index and heartbeating; also ensure OSS standbys emit periodic metrics. [[GH-12166](https://github.com/hashicorp/vault/pull/12166)]
* ui: updates clients configuration edit form state based on census reporting configuration [[GH-20125](https://github.com/hashicorp/vault/pull/20125)]

BUG FIXES:

* core/activity: add namespace breakdown for new clients when date range spans multiple months, including the current month. [[GH-18766](https://github.com/hashicorp/vault/pull/18766)]
* core/activity: de-duplicate namespaces when historical and current month data are mixed [[GH-18452](https://github.com/hashicorp/vault/pull/18452)]
* core/activity: fix the end_date returned from the activity log endpoint when partial counts are computed [[GH-17856](https://github.com/hashicorp/vault/pull/17856)]
* core/activity: include mount counts when de-duplicating current and historical month data [[GH-18598](https://github.com/hashicorp/vault/pull/18598)]
* core/activity: report mount paths (rather than mount accessors) in current month activity log counts and include deleted mount paths in precomputed queries. [[GH-18916](https://github.com/hashicorp/vault/pull/18916)]
* core/activity: return partial month counts when querying a historical date range and no historical data exists. [[GH-17935](https://github.com/hashicorp/vault/pull/17935)]
* core: Change where we evaluate filtered paths as part of mount operations; this is part of an enterprise bugfix that will
have its own changelog entry. [[GH-21260](https://github.com/hashicorp/vault/pull/21260)]
* core: Do not cache seal configuration to fix a bug that resulted in sporadic auto unseal failures. [[GH-21223](https://github.com/hashicorp/vault/pull/21223)]
* core: Don't exit just because we think there's a potential deadlock. [[GH-21342](https://github.com/hashicorp/vault/pull/21342)]
* core: Fix panic in sealed nodes using raft storage trying to emit raft metrics [[GH-21249](https://github.com/hashicorp/vault/pull/21249)]
* identity: Fixes duplicate groups creation with the same name but unique IDs. [[GH-20964](https://github.com/hashicorp/vault/pull/20964)]
* replication (enterprise): Fix a race condition with update-primary that could result in data loss after a DR failover
* replication (enterprise): Fix path filters deleting data right after it's written by backend Initialize funcs

## 1.11.11
### June 08, 2023

CHANGES:

* core: Bump Go version to 1.19.9.
* core: Revert #19676 (VAULT_GRPC_MIN_CONNECT_TIMEOUT env var) as we decided it was unnecessary. [[GH-20826](https://github.com/hashicorp/vault/pull/20826)]

IMPROVEMENTS:

* command/server: Add support for dumping pprof files to the filesystem via SIGUSR2 when
`VAULT_PPROF_WRITE_TO_FILE=true` is set on the server. [[GH-20609](https://github.com/hashicorp/vault/pull/20609)]
* secrets/pki: add subject key identifier to read key response [[GH-20642](https://github.com/hashicorp/vault/pull/20642)]
* ui: update TTL picker for consistency [[GH-18114](https://github.com/hashicorp/vault/pull/18114)]

BUG FIXES:

* api: Properly Handle nil identity_policies in Secret Data [[GH-20636](https://github.com/hashicorp/vault/pull/20636)]
* auth/ldap: Set default value for `max_page_size` properly [[GH-20453](https://github.com/hashicorp/vault/pull/20453)]
* cli: CLI should take days as a unit of time for ttl like flags [[GH-20477](https://github.com/hashicorp/vault/pull/20477)]
* core (enterprise): Fix log shipper buffer size overflow issue for 32 bit architecture.
* core (enterprise): Fix logshipper buffer size to default to DefaultBufferSize only when reported system memory is zero.
* core (enterprise): Remove MFA Enforcment configuration for namespace when deleting namespace
* core: prevent panic on login after namespace is deleted that had mfa enforcement [[GH-20375](https://github.com/hashicorp/vault/pull/20375)]
* replication (enterprise): Fix a race condition with invalid tokens during WAL streaming that was causing Secondary clusters to be unable to connect to a Primary.
* replication (enterprise): fix bug where secondary grpc connections would timeout when connecting to a primary host that no longer exists.
* secrets/transform (enterprise): Fix a caching bug affecting secondary nodes after a tokenization key rotation

## 1.11.10
### April 26, 2023

CHANGES:

* core: Bump Go version to 1.19.8.

IMPROVEMENTS:

* cli/namespace: Add detailed flag to output additional namespace information
such as namespace IDs and custom metadata. [[GH-20243](https://github.com/hashicorp/vault/pull/20243)]
* core/activity: add an endpoint to write test activity log data, guarded by a build flag [[GH-20019](https://github.com/hashicorp/vault/pull/20019)]
* core: Add a `raft` sub-field to the `storage` and `ha_storage` details provided by the
`/sys/config/state/sanitized` endpoint in order to include the `max_entry_size`. [[GH-20044](https://github.com/hashicorp/vault/pull/20044)]
* sdk/ldaputil: added `connection_timeout` to tune connection timeout duration
for all LDAP plugins. [[GH-20144](https://github.com/hashicorp/vault/pull/20144)]

BUG FIXES:

* auth/ldap: Add max_page_size configurable to LDAP configuration [[GH-19032](https://github.com/hashicorp/vault/pull/19032)]
* core (enterprise): Fix intermittent issue with token entries sometimes not being found when using a newly created token in a request to a secondary, even when SSCT `new_token` forwarding is set. When this occurred, this would result in the following error to the client: `error performing token check: no lease entry found for token that ought to have one, possible eventual consistency issue`.
* core (enterprise): Fix read on perf standbys failing with 412 after leadership change, unseal, restores or restarts when no writes occur
* core/ssct (enterprise): Fixed race condition where a newly promoted DR may revert `sscGenCounter`
resulting in 412 errors.
* core: Fix regression breaking non-raft clusters whose nodes share the same cluster_addr/api_addr. [[GH-19721](https://github.com/hashicorp/vault/pull/19721)]
* helper/random: Fix race condition in string generator helper [[GH-19875](https://github.com/hashicorp/vault/pull/19875)]
* openapi: Fix many incorrect details in generated API spec, by using better techniques to parse path regexps [[GH-18554](https://github.com/hashicorp/vault/pull/18554)]
* replication (enterprise): Fix replication status for Primary clusters showing its primary cluster's information (in case of DR) in secondaries field when known_secondaries field is nil
* secrets/pki: Fix patching of leaf_not_after_behavior on issuers. [[GH-20341](https://github.com/hashicorp/vault/pull/20341)]
* secrets/transform (enterprise): Address SQL connection leak when cleaning expired tokens
* ui: Fix OIDC provider logo showing when domain doesn't match [[GH-20263](https://github.com/hashicorp/vault/pull/20263)]
* ui: Fix bad link to namespace when namespace name includes `.` [[GH-19799](https://github.com/hashicorp/vault/pull/19799)]
* ui: fixes browser console formatting for help command output [[GH-20064](https://github.com/hashicorp/vault/pull/20064)]
* ui: remove use of htmlSafe except when first sanitized [[GH-20235](https://github.com/hashicorp/vault/pull/20235)]

## 1.11.9
### March 29, 2023

SECURITY:

* storage/mssql: When using Vault’s community-supported Microsoft SQL (MSSQL) database storage backend, a privileged attacker with the ability to write arbitrary data to Vault’s configuration may be able to perform arbitrary SQL commands on the underlying database server through Vault. This vulnerability, CVE-2023-0620, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-12](https://discuss.hashicorp.com/t/hcsec-2023-12-vault-s-microsoft-sql-database-storage-backend-vulnerable-to-sql-injection-via-configuration-file/52080)]
* secrets/pki: Vault’s PKI mount issuer endpoints did not correctly authorize access to remove an issuer or modify issuer metadata, potentially resulting in denial of service of the PKI mount. This bug did not affect public or private key material, trust chains or certificate issuance. This vulnerability, CVE-2023-0665, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-11](https://discuss.hashicorp.com/t/hcsec-2023-11-vault-s-pki-issuer-endpoint-did-not-correctly-authorize-access-to-issuer-metadata/52079)]
* core: HashiCorp Vault’s implementation of Shamir’s secret sharing used precomputed table lookups, and was vulnerable to cache-timing attacks. An attacker with access to, and the ability to observe a large number of unseal operations on the host through a side channel may reduce the search space of a brute force effort to recover the Shamir shares. This vulnerability, CVE-2023-25000, is fixed in Vault 1.13.1, 1.12.5, and 1.11.9. [[HCSEC-2023-10](https://discuss.hashicorp.com/t/hcsec-2023-10-vault-vulnerable-to-cache-timing-attacks-during-seal-and-unseal-operations/52078)]

IMPROVEMENTS:

* auth/github: Allow for an optional Github auth token environment variable to make authenticated requests when fetching org id
website/docs: Add docs for `VAULT_AUTH_CONFIG_GITHUB_TOKEN` environment variable when writing Github config [[GH-19244](https://github.com/hashicorp/vault/pull/19244)]
* core: Allow overriding gRPC connect timeout via VAULT_GRPC_MIN_CONNECT_TIMEOUT. This is an env var rather than a config setting because we don't expect this to ever be needed.  It's being added as a last-ditch
option in case all else fails for some replication issues we may not have fully reproduced. [[GH-19676](https://github.com/hashicorp/vault/pull/19676)]
* core: validate name identifiers in mssql physical storage backend prior use [[GH-19591](https://github.com/hashicorp/vault/pull/19591)]

BUG FIXES:

* auth/kubernetes: Ensure a consistent TLS configuration for all k8s API requests [[#190](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/190)] [[GH-19720](https://github.com/hashicorp/vault/pull/19720)]
* cli: Fix vault read handling to return raw data as secret.Data when there is no top-level data object from api response. [[GH-17913](https://github.com/hashicorp/vault/pull/17913)]
* core (enterprise): Attempt to reconnect to a PKCS#11 HSM if we retrieve a CKR_FUNCTION_FAILED error.
* core: Fixed issue with remounting mounts that have a non-trailing space in the 'to' or 'from' paths. [[GH-19585](https://github.com/hashicorp/vault/pull/19585)]
* openapi: Fix logic for labeling unauthenticated/sudo paths. [[GH-19600](https://github.com/hashicorp/vault/pull/19600)]
* secrets/transform (enterprise): Fix persistence problem with rotated tokenization key versions
* ui: fixes issue navigating back a level using the breadcrumb from secret metadata view [[GH-19703](https://github.com/hashicorp/vault/pull/19703)]
* ui: pass encodeBase64 param to HMAC transit-key-actions. [[GH-19429](https://github.com/hashicorp/vault/pull/19429)]
* ui: use URLSearchParams interface to capture namespace param from SSOs (ex. ADFS) with decoded state param in callback url [[GH-19460](https://github.com/hashicorp/vault/pull/19460)]

## 1.11.8
### March 01, 2023

SECURITY:

* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* core: Bump Go version to 1.19.6.

IMPROVEMENTS:

* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]

BUG FIXES:

* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core (enterprise): Fix perf standby WAL streaming silently failures when replication setup happens at a bad time.
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* license (enterprise): Fix bug where license would update even if the license didn't change.
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* secrets/ad: Fix bug where config couldn't be updated unless binddn/bindpass were included in the update. [[GH-18208](https://github.com/hashicorp/vault/pull/18208)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]
* ui: show Get credentials button for static roles detail page when a user has the proper permissions. [[GH-19190](https://github.com/hashicorp/vault/pull/19190)]

## 1.11.7
### February 6, 2023

CHANGES:

* core: Bump Go version to 1.19.4.

IMPROVEMENTS:

* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]

BUG FIXES:

* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* core (enterprise): Fix a race condition resulting in login errors to PKCS#11 modules under high concurrency.
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* kmip (enterprise): Fix Destroy operation response that omitted Unique Identifier on some batched responses.
* kmip (enterprise): Fix Locate operation response incompatibility with clients using KMIP versions prior to 1.3.
* kmip (enterprise): Fix Query operation response that omitted streaming capability and supported profiles.
* licensing (enterprise): update autoloaded license cache after reload
* secrets/pki: Allow patching issuer to set an empty issuer name. [[GH-18466](https://github.com/hashicorp/vault/pull/18466)]
* secrets/transit: Do not warn about unrecognized parameter 'batch_input' [[GH-18299](https://github.com/hashicorp/vault/pull/18299)]
* storage/raft (enterprise): An already joined node can rejoin by wiping storage
and re-issueing a join request, but in doing so could transiently become a
non-voter.  In some scenarios this resulted in loss of quorum. [[GH-18263](https://github.com/hashicorp/vault/pull/18263)]
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* storage/raft: Don't panic on unknown raft ops [[GH-17732](https://github.com/hashicorp/vault/pull/17732)]
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]

## 1.11.6
### November 30, 2022

IMPROVEMENTS:

* secrets/pki: Allow issuer creation, import to change default issuer via `default_follows_latest_issuer`. [[GH-17824](https://github.com/hashicorp/vault/pull/17824)]

BUG FIXES:

* auth/okta: fix a panic for AuthRenew in Okta [[GH-18011](https://github.com/hashicorp/vault/pull/18011)]
* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* cli: Fix issue preventing kv commands from executing properly when the mount path provided by `-mount` flag and secret key path are the same. [[GH-17679](https://github.com/hashicorp/vault/pull/17679)]
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: fix a start up race condition where performance standbys could go into a
  mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* secret/pki: fix bug with initial legacy bundle migration (from < 1.11 into 1.11+) and missing issuers from ca_chain [[GH-17772](https://github.com/hashicorp/vault/pull/17772)]
* secrets/azure: add WAL to clean up role assignments if errors occur [[GH-18085](https://github.com/hashicorp/vault/pull/18085)]
* secrets/gcp: Fixes duplicate service account key for rotate root on standby or secondary [[GH-18110](https://github.com/hashicorp/vault/pull/18110)]
* ui: Fixes issue with not being able to download raft snapshot via service worker [[GH-17769](https://github.com/hashicorp/vault/pull/17769)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]

## 1.11.5
### November 2, 2022

IMPROVEMENTS:

* database/snowflake: Allow parallel requests to Snowflake [[GH-17594](https://github.com/hashicorp/vault/pull/17594)]
* sdk/ldap: Added support for paging when searching for groups using group filters [[GH-17640](https://github.com/hashicorp/vault/pull/17640)]

BUG FIXES:

* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): fix panic when having `cache_disable` true
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* kmip (enterprise): Fix a problem in the handling of attributes that caused Import operations to fail.
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* secrets/pki: Do not ignore provided signature bits value when signing intermediate and leaf certificates with a managed key [[GH-17328](https://github.com/hashicorp/vault/pull/17328)]
* secrets/pki: Do not read revoked certificates from backend when CRL is disabled [[GH-17384](https://github.com/hashicorp/vault/pull/17384)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* ui/keymgmt: Sets the defaultValue for type when creating a key. [[GH-17407](https://github.com/hashicorp/vault/pull/17407)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]

## 1.11.4
### September 30, 2022

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

IMPROVEMENTS:

* agent/auto-auth: Add `exit_on_err` which when set to true, will cause Agent to exit if any errors are encountered during authentication. [[GH-17091](https://github.com/hashicorp/vault/pull/17091)]
* agent: Send notifications to systemd on start and stop. [[GH-9802](https://github.com/hashicorp/vault/pull/9802)]

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* auth/kubernetes: Restore support for JWT signature algorithm ES384 [[GH-160](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/160)] [[GH-17162](https://github.com/hashicorp/vault/pull/17162)]
* auth/token: Fix ignored parameter warnings for valid parameters on token create [[GH-16938](https://github.com/hashicorp/vault/pull/16938)]
* core/quotas: Fix goroutine leak caused by the seal process not fully cleaning up Rate Limit Quotas. [[GH-17281](https://github.com/hashicorp/vault/pull/17281)]
* core: Prevent two or more DR failovers from invalidating SSCT tokens generated on the previous primaries. [[GH-16956](https://github.com/hashicorp/vault/pull/16956)]
* identity/oidc: Adds `claims_supported` to discovery document. [[GH-16992](https://github.com/hashicorp/vault/pull/16992)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* secrets/transform (enterprise): Fix an issue loading tokenization transform configuration after a specific sequence of reconfigurations.
* secrets/transform (enterprise): Fix persistence problem with tokenization store credentials.
* ui: Fixes secret version and status menu links transitioning to auth screen [[GH-16983](https://github.com/hashicorp/vault/pull/16983)]
* ui: Fixes secret version and status menu links transitioning to auth screen [[GH-16983](https://github.com/hashicorp/vault/pull/16983)]

## 1.11.3
### August 31, 2022

SECURITY:

* core: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. This vulnerability, CVE-2022-40186, is fixed in 1.11.3, 1.10.6, and 1.9.9. [[HSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

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

SECURITY:

* identity/entity: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. [[HCSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

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

SECURITY:

* storage/raft: Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

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

## 1.10.11
### March 01, 2023

SECURITY:

* auth/approle: When using the Vault and Vault Enterprise (Vault) approle auth method, any authenticated user with access to the /auth/approle/role/:role_name/secret-id-accessor/destroy endpoint can destroy the secret ID of any other role by providing the secret ID accessor. This vulnerability, CVE-2023-24999 has been fixed in Vault 1.13.0, 1.12.4, 1.11.8, 1.10.11 and above. [[HSEC-2023-07](https://discuss.hashicorp.com/t/hcsec-2023-07-vault-fails-to-verify-if-approle-secretid-belongs-to-role-during-a-destroy-operation/51305)]

CHANGES:

* core: Bump Go version to 1.19.6.

IMPROVEMENTS:

* secrets/database: Adds error message requiring password on root crednetial rotation. [[GH-19103](https://github.com/hashicorp/vault/pull/19103)]

BUG FIXES:

* auth/approle: Add nil check for the secret ID entry when deleting via secret id accessor preventing cross role secret id deletion [[GH-19186](https://github.com/hashicorp/vault/pull/19186)]
* core (enterprise): Fix panic when using invalid accessor for control-group request
* core: Prevent panics in `sys/leases/lookup`, `sys/leases/revoke`, and `sys/leases/renew` endpoints if provided `lease_id` is null [[GH-18951](https://github.com/hashicorp/vault/pull/18951)]
* replication (enterprise): Fix bug where reloading external plugin on a secondary would
break replication.
* secrets/ad: Fix bug where config couldn't be updated unless binddn/bindpass were included in the update. [[GH-18209](https://github.com/hashicorp/vault/pull/18209)]
* ui (enterprise): Fix cancel button from transform engine role creation page [[GH-19135](https://github.com/hashicorp/vault/pull/19135)]
* ui: Fix bug where logging in via OIDC fails if browser is in fullscreen mode [[GH-19071](https://github.com/hashicorp/vault/pull/19071)]

## 1.10.10
### February 6, 2023

CHANGES:

* core: Bump Go version to 1.19.4.

IMPROVEMENTS:

* command/server: Environment variable keys are now logged at startup. [[GH-18125](https://github.com/hashicorp/vault/pull/18125)]
* core/fips: use upstream toolchain for FIPS 140-2 compliance again; this will appear as X=boringcrypto on the Go version in Vault server logs.
* secrets/db/mysql: Add `tls_server_name` and `tls_skip_verify` parameters [[GH-18799](https://github.com/hashicorp/vault/pull/18799)]
* ui: Prepends "passcode=" if not provided in user input for duo totp mfa method authentication [[GH-18342](https://github.com/hashicorp/vault/pull/18342)]
* ui: Update language on database role to "Connection name" [[GH-18261](https://github.com/hashicorp/vault/issues/18261)] [[GH-18350](https://github.com/hashicorp/vault/pull/18350)]

BUG FIXES:

* auth/approle: Fix `token_bound_cidrs` validation when using /32 blocks for role and secret ID [[GH-18145](https://github.com/hashicorp/vault/pull/18145)]
* auth/token: Fix ignored parameter warnings for valid parameters on token create [[GH-16938](https://github.com/hashicorp/vault/pull/16938)]
* cli/kv: skip formatting of nil secrets for patch and put with field parameter set [[GH-18163](https://github.com/hashicorp/vault/pull/18163)]
* core (enterprise): Fix a race condition resulting in login errors to PKCS#11 modules under high concurrency.
* core/managed-keys (enterprise): Limit verification checks to mounts in a key's namespace
* core/quotas (enterprise): Fix a potential deadlock that could occur when using lease count quotas.
* core/quotas: Fix issue with improper application of default rate limit quota exempt paths [[GH-18273](https://github.com/hashicorp/vault/pull/18273)]
* core: fix bug where context cancellations weren't forwarded to active node from performance standbys.
* core: prevent panic in login mfa enforcement delete after enforcement's namespace is deleted [[GH-18923](https://github.com/hashicorp/vault/pull/18923)]
* database/mongodb: Fix writeConcern set to be applied to any query made on the database [[GH-18546](https://github.com/hashicorp/vault/pull/18546)]
* identity (enterprise): Fix a data race when creating an entity for a local alias.
* kmip (enterprise): Fix Destroy operation response that omitted Unique Identifier on some batched responses.
* kmip (enterprise): Fix Locate operation response incompatibility with clients using KMIP versions prior to 1.3.
* licensing (enterprise): update autoloaded license cache after reload
* storage/raft (enterprise): Fix some storage-modifying RPCs used by perf standbys that weren't returning the resulting WAL state.
* ui: fixes query parameters not passed in api explorer test requests [[GH-18743](https://github.com/hashicorp/vault/pull/18743)]

## 1.10.9
### November 30, 2022

BUG FIXES:

* auth: Deduplicate policies prior to ACL generation [[GH-17914](https://github.com/hashicorp/vault/pull/17914)]
* core/quotas (enterprise): Fix a lock contention issue that could occur and cause Vault to become unresponsive when creating, changing, or deleting lease count quotas.
* core: Fix potential deadlock if barrier ciphertext is less than 4 bytes. [[GH-17944](https://github.com/hashicorp/vault/pull/17944)]
* core: fix a start up race condition where performance standbys could go into a
  mount loop if default policies are not yet synced from the active node. [[GH-17801](https://github.com/hashicorp/vault/pull/17801)]
* secrets/azure: add WAL to clean up role assignments if errors occur [[GH-18084](https://github.com/hashicorp/vault/pull/18084)]
* secrets/gcp: Fixes duplicate service account key for rotate root on standby or secondary [[GH-18109](https://github.com/hashicorp/vault/pull/18109)]
* ui: fix entity policies list link to policy show page [[GH-17950](https://github.com/hashicorp/vault/pull/17950)]

## 1.10.8
### November 2, 2022

BUG FIXES:

* core/managed-keys (enterprise): Return better error messages when encountering key creation failures
* core/managed-keys (enterprise): fix panic when having `cache_disable` true
* core: prevent memory leak when using control group factors in a policy [[GH-17532](https://github.com/hashicorp/vault/pull/17532)]
* core: prevent panic during mfa after enforcement's namespace is deleted [[GH-17562](https://github.com/hashicorp/vault/pull/17562)]
* login: Store token in tokenhelper for interactive login MFA [[GH-17040](https://github.com/hashicorp/vault/pull/17040)]
* secrets/pki: Do not ignore provided signature bits value when signing intermediate and leaf certificates with a managed key [[GH-17328](https://github.com/hashicorp/vault/pull/17328)]
* secrets/pki: Respond to tidy-status, tidy-cancel on PR Secondary clusters. [[GH-17497](https://github.com/hashicorp/vault/pull/17497)]
* ui: Fixes oidc/jwt login issue with alternate mount path and jwt login via mount path tab [[GH-17661](https://github.com/hashicorp/vault/pull/17661)]

## 1.10.7
### September 30, 2022

SECURITY:

* secrets/pki: Vault’s TLS certificate auth method did not initially load the optionally-configured CRL issued by the role’s CA into memory on startup, resulting in the revocation list not being checked, if the CRL has not yet been retrieved. This vulnerability, CVE-2022-41316, is fixed in Vault 1.12.0, 1.11.4, 1.10.7, and 1.9.10. [[HSEC-2022-24](https://discuss.hashicorp.com/t/hcsec-2022-24-vaults-tls-cert-auth-method-only-loaded-crl-after-first-request/45483)]

BUG FIXES:

* auth/cert: Vault does not initially load the CRLs in cert auth unless the read/write CRL endpoint is hit. [[GH-17138](https://github.com/hashicorp/vault/pull/17138)]
* core/quotas: Fix goroutine leak caused by the seal process not fully cleaning up Rate Limit Quotas. [[GH-17281](https://github.com/hashicorp/vault/pull/17281)]
* core: Prevent two or more DR failovers from invalidating SSCT tokens generated on the previous primaries. [[GH-16956](https://github.com/hashicorp/vault/pull/16956)]
* identity/oidc: Adds `claims_supported` to discovery document. [[GH-16992](https://github.com/hashicorp/vault/pull/16992)]
* replication (enterprise): Fix data race in SaveCheckpoint()
* secrets/transform (enterprise): Fix an issue loading tokenization transform configuration after a specific sequence of reconfigurations.
* secrets/transform (enterprise): Fix persistence problem with tokenization store credentials.
* ui: Fix lease force revoke action [[GH-16930](https://github.com/hashicorp/vault/pull/16930)]

## 1.10.6
### August 31, 2022

SECURITY:

* core: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. This vulnerability, CVE-2022-40186, is fixed in 1.11.3, 1.10.6, and 1.9.9. [[HSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

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

SECURITY:

* identity/entity: When entity aliases mapped to a single entity share the same alias name, but have different mount accessors, Vault can leak metadata between the aliases. This metadata leak may result in unexpected access if templated policies are using alias metadata for path names. [[HCSEC-2022-18](https://discuss.hashicorp.com/t/hcsec-2022-18-vault-entity-alias-metadata-may-leak-between-aliases-with-the-same-name-assigned-to-the-same-entity/44550)]

## 1.10.5
### July 21, 2022

SECURITY:

* storage/raft: Vault Enterprise (“Vault”) clusters using Integrated Storage expose an unauthenticated API endpoint that could be abused to override the voter status of a node within a Vault HA cluster, introducing potential for future data loss or catastrophic failure. This vulnerability, CVE-2022-36129, was fixed in Vault 1.9.8, 1.10.5, and 1.11.1. [[HSEC-2022-15](https://discuss.hashicorp.com/t/hcsec-2022-15-vault-enterprise-does-not-verify-existing-voter-status-when-joining-an-integrated-storage-ha-node/42420)]

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
