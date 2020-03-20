## 1.4 RC1 (March 19th, 2020)

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
* storage/raft: Fix a potential deadlock that could occure on leadership transition [[GH-8547](https://github.com/hashicorp/vault/pull/8547)]
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
* replication: Fix issue causing cubbyholes in namespaces on performance secondaries to not work.
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

## 1.3.4 (March 19th, 2020)

SECURITY:

* A vulnerability was identified in Vault and Vault Enterprise such that, under certain circumstances,  an Entity's Group membership may inadvertently include Groups the Entity no longer has permissions to. This vulnerability, CVE-2020-10660, affects Vault and Vault Enterprise versions 0.9.0 and newer, and is fixed in 1.3.4.
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
 * auth/aws: Add aws metadata to identity alias [[GH-7975](https://github.com/hashicorp/vault/pull/7975)]
 * auth/kubernetes: Allow both names and namespaces to be set to "*" [[GH-78](https://github.com/hashicorp/vault-plugin-auth-kubernetes/pull/78)]

BUG FIXES:

* auth/azure: Fix Azure compute client to use correct base URL [[GH-8072](https://github.com/hashicorp/vault/pull/8072)]
* auth/ldap: Fix renewal of tokens without cofigured policies that are
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

## 0.11.6 (December 14th, 2018)

This release contains the three security fixes from 1.0.0 and 1.0.1 and the
following bug fixes from 1.0.0/1.0.1:

 * namespaces: Correctly reload the proper mount when tuning or reloading the
   mount [[GH-5937](https://github.com/hashicorp/vault/pull/5937)]
 * replication/perfstandby: Fix audit table upgrade on standbys [[GH-5811](https://github.com/hashicorp/vault/pull/5811)]
 * replication/perfstandby: Fix redirect on approle update [[GH-5820](https://github.com/hashicorp/vault/pull/5820)]
 * secrets/kv: Fix issue where storage version would get incorrectly downgraded
   [[GH-5809](https://github.com/hashicorp/vault/pull/5809)]

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

## 0.11.5 (November 13th, 2018)

BUG FIXES:

 * agent: Fix issue when specifying two file sinks [[GH-5610](https://github.com/hashicorp/vault/pull/5610)]
 * auth/userpass: Fix minor timing issue that could leak the presence of a
   username [[GH-5614](https://github.com/hashicorp/vault/pull/5614)]
 * autounseal/alicloud: Fix issue interacting with the API (Enterprise)
 * autounseal/azure: Fix key version tracking (Enterprise)
 * cli: Fix panic that could occur if parameters were not provided [[GH-5603](https://github.com/hashicorp/vault/pull/5603)]
 * core: Fix buggy behavior if trying to remount into a namespace
 * identity: Fix duplication of entity alias entity during alias transfer
   between entities [[GH-5733](https://github.com/hashicorp/vault/pull/5733)]
 * namespaces: Fix tuning of auth mounts in a namespace
 * ui: Fix bug where editing secrets as JSON doesn't save properly [[GH-5660](https://github.com/hashicorp/vault/pull/5660)]
 * ui: Fix issue where IE 11 didn't render the UI and also had a broken form
   when trying to use tool/hash [[GH-5714](https://github.com/hashicorp/vault/pull/5714)]

## 0.11.4 (October 23rd, 2018)

CHANGES:

 * core: HA lock file is no longer copied during `operator migrate` [[GH-5503](https://github.com/hashicorp/vault/pull/5503)].
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

 * core: Add last WAL in leader/health output for easier debugging [[GH-5523](https://github.com/hashicorp/vault/pull/5523)]
 * identity: Identity names will now be handled case insensitively by default.
   This includes names of entities, aliases and groups [[GH-5404](https://github.com/hashicorp/vault/pull/5404)]
 * secrets/aws: Added role-option max_sts_ttl to cap TTL for AWS STS
   credentials [[GH-5500](https://github.com/hashicorp/vault/pull/5500)]
 * secret/database: Allow Cassandra user to be non-superuser so long as it has
   role creation permissions [[GH-5402](https://github.com/hashicorp/vault/pull/5402)]
 * secret/radius: Allow setting the NAS Identifier value in the generated
   packet [[GH-5465](https://github.com/hashicorp/vault/pull/5465)]
 * secret/ssh: Allow usage of JSON arrays when setting zero addresses [[GH-5528](https://github.com/hashicorp/vault/pull/5528)]
 * secret/transit: Allow trimming unused keys [[GH-5388](https://github.com/hashicorp/vault/pull/5388)]
 * ui: Support KVv2 [[GH-5547](https://github.com/hashicorp/vault/pull/5547)], [[GH-5563](https://github.com/hashicorp/vault/pull/5563)]
 * ui: Allow viewing and updating Vault license via the UI
 * ui: Onboarding will now display your progress through the chosen tutorials
 * ui: Dynamic secret backends obfuscate sensitive data by default and
   visibility is toggleable

BUG FIXES:

 * agent: Fix potential hang during agent shutdown [[GH-5026](https://github.com/hashicorp/vault/pull/5026)]
 * auth/ldap: Fix listing of users/groups that contain slashes [[GH-5537](https://github.com/hashicorp/vault/pull/5537)]
 * core: Fix memory leak during some expiration calls [[GH-5505](https://github.com/hashicorp/vault/pull/5505)]
 * core: Fix generate-root operations requiring empty `otp` to be provided
   instead of an empty body [[GH-5495](https://github.com/hashicorp/vault/pull/5495)]
 * identity: Remove lookup check during alias removal from entity [[GH-5524](https://github.com/hashicorp/vault/pull/5524)]
 * secret/pki: Fix TTL/MaxTTL check when using `sign-verbatim` [[GH-5549](https://github.com/hashicorp/vault/pull/5549)]
 * secret/pki: Fix regression in 0.11.2+ causing the NotBefore value of
   generated certificates to be set to the Unix epoch if the role value was not
   set, instead of using the default of 30 seconds [[GH-5481](https://github.com/hashicorp/vault/pull/5481)]
 * storage/mysql: Use `varbinary` instead of `varchar` when creating HA tables
   [[GH-5529](https://github.com/hashicorp/vault/pull/5529)]

## 0.11.3 (October 8th, 2018)

SECURITY:

 * Revocation: A regression in 0.11.2 (OSS) and 0.11.0 (Enterprise) caused
   lease IDs containing periods (`.`) to not be revoked properly. Upon startup
   when revocation is tried again these should now revoke successfully.

IMPROVEMENTS:

 * auth/ldap: Listing of users and groups return absolute paths [[GH-5537](https://github.com/hashicorp/vault/pull/5537)]
 * secret/pki: OID SANs can now specify `*` to allow any value [[GH-5459](https://github.com/hashicorp/vault/pull/5459)]

BUG FIXES:

 * auth/ldap: Fix panic if specific values were given to be escaped [[GH-5471](https://github.com/hashicorp/vault/pull/5471)]
 * cli/auth: Fix panic if `vault auth` was given no parameters [[GH-5473](https://github.com/hashicorp/vault/pull/5473)]
 * secret/database/mongodb: Fix panic that could occur at high load [[GH-5463](https://github.com/hashicorp/vault/pull/5463)]
 * secret/pki: Fix CA generation not allowing OID SANs [[GH-5459](https://github.com/hashicorp/vault/pull/5459)]

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
   credentials it is using [[GH-5140](https://github.com/hashicorp/vault/pull/5140)]
 * **Storage Backend Migrator**: A new `operator migrate` command allows offline
   migration of data between two storage backends
 * **AliCloud KMS Auto Unseal and Seal Wrap Support (Enterprise)**: AliCloud KMS can now be used a support seal for
   Auto Unseal and Seal Wrapping

BUG FIXES:

 * auth/okta: Fix reading deprecated `token` parameter if a token was
   previously set in the configuration [[GH-5409](https://github.com/hashicorp/vault/pull/5409)]
 * core: Re-add deprecated capabilities information for now [[GH-5360](https://github.com/hashicorp/vault/pull/5360)]
 * core: Fix handling of cyclic token relationships [[GH-4803](https://github.com/hashicorp/vault/pull/4803)]
 * storage/mysql: Fix locking on MariaDB [[GH-5343](https://github.com/hashicorp/vault/pull/5343)]
 * replication: Fix DR API when using a token [[GH-5398](https://github.com/hashicorp/vault/pull/5398)]
 * identity: Ensure old group alias is removed when a new one is written [[GH-5350](https://github.com/hashicorp/vault/pull/5350)]
 * storage/alicloud: Don't call uname on package init [[GH-5358](https://github.com/hashicorp/vault/pull/5358)]
 * secrets/jwt: Fix issue where request context would be canceled too early
 * ui: fix need to have update for aws iam creds generation [GF-5294]
 * ui: fix calculation of token expiry [[GH-5435](https://github.com/hashicorp/vault/pull/5435)]

IMPROVEMENTS:

 * auth/aws: The identity alias name can now configured to be either IAM unique
   ID of the IAM Principal, or ARN of the caller identity [[GH-5247](https://github.com/hashicorp/vault/pull/5247)]
 * auth/cert: Add allowed_organizational_units support [[GH-5252](https://github.com/hashicorp/vault/pull/5252)]
 * cli: Format TTLs for non-secret responses [[GH-5367](https://github.com/hashicorp/vault/pull/5367)]
 * identity: Support operating on entities and groups by their names [[GH-5355](https://github.com/hashicorp/vault/pull/5355)]
 * plugins: Add `env` parameter when registering plugins to the catalog to allow
   operators to include environment variables during plugin execution. [[GH-5359](https://github.com/hashicorp/vault/pull/5359)]
 * secrets/aws: WAL Rollback improvements [[GH-5202](https://github.com/hashicorp/vault/pull/5202)]
 * secrets/aws: Allow specifying STS role-default TTLs [[GH-5138](https://github.com/hashicorp/vault/pull/5138)]
 * secrets/pki: Add configuration support for setting NotBefore [[GH-5325](https://github.com/hashicorp/vault/pull/5325)]
 * core: Support for passing the Vault token via an Authorization Bearer header [[GH-5397](https://github.com/hashicorp/vault/pull/5397)]
 * replication: Reindex process now runs in the background and does not block other
   vault operations
 * storage/zookeeper: Enable TLS based communication with Zookeeper [[GH-4856](https://github.com/hashicorp/vault/pull/4856)]
 * ui: you can now init a cluster with a seal config [[GH-5428](https://github.com/hashicorp/vault/pull/5428)]
 * ui: added the option to force promote replication clusters [[GH-5438](https://github.com/hashicorp/vault/pull/5438)]
 * replication: Allow promotion of a secondary when data is syncing with a "force" flag

## 0.11.1.1 (September 17th, 2018) (Enterprise Only)

BUG FIXES:

 * agent: Fix auth handler-based wrapping of output tokens [[GH-5316](https://github.com/hashicorp/vault/pull/5316)]
 * core: Properly store the replication checkpoint file if it's larger than the
   storage engine's per-item limit
 * core: Improve WAL deletion rate
 * core: Fix token creation on performance standby nodes
 * core: Fix unwrapping inside a namespace
 * core: Always forward tidy operations from performance standby nodes

IMPROVEMENTS:

 * auth/aws: add support for key/value pairs or JSON values for
   `iam_request_headers` with IAM auth method [[GH-5320](https://github.com/hashicorp/vault/pull/5320)]
 * auth/aws, secret/aws: Throttling errors from the AWS API will now be
   reported as 502 errors by Vault, along with the original error [[GH-5270](https://github.com/hashicorp/vault/pull/5270)]
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
   Vault's config file [[GH-5280](https://github.com/hashicorp/vault/pull/5280)]

BUG FIXES:

 * core: Ensure we use a background context when stepping down [[GH-5290](https://github.com/hashicorp/vault/pull/5290)]
 * core: Properly check error return from random byte reading [[GH-5277](https://github.com/hashicorp/vault/pull/5277)]
 * core: Re-add `sys/` top-route injection for now [[GH-5241](https://github.com/hashicorp/vault/pull/5241)]
 * core: Policies stored in minified JSON would return an error [[GH-5229](https://github.com/hashicorp/vault/pull/5229)]
 * core: Evaluate templated policies in capabilities check [[GH-5250](https://github.com/hashicorp/vault/pull/5250)]
 * identity: Update MemDB with identity group alias while loading groups [[GH-5289](https://github.com/hashicorp/vault/pull/5289)]
 * secrets/database: Fix nil pointer when revoking some leases [[GH-5262](https://github.com/hashicorp/vault/pull/5262)]
 * secrets/pki: Fix sign-verbatim losing extra Subject attributes [[GH-5245](https://github.com/hashicorp/vault/pull/5245)]
 * secrets/pki: Remove certificates from store when tidying revoked
   certificates and simplify API [[GH-5231](https://github.com/hashicorp/vault/pull/5231)]
 * ui: JSON editor will not coerce input to an object, and will now show an
   error about Vault expecting an object [[GH-5271](https://github.com/hashicorp/vault/pull/5271)]
 * ui: authentication form will now default to any methods that have been tuned
   to show up for unauthenticated users [[GH-5281](https://github.com/hashicorp/vault/pull/5281)]


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
   guidance, linking out to relevant links or guides on vaultproject.io for
   various workflows in Vault.

IMPROVEMENTS:

 * agent: Add `exit_after_auth` to be able to use the Agent for a single
   authentication [[GH-5013](https://github.com/hashicorp/vault/pull/5013)]
 * auth/approle: Add ability to set token bound CIDRs on individual Secret IDs
   [[GH-5034](https://github.com/hashicorp/vault/pull/5034)]
 * cli: Add support for passing parameters to `vault read` operations [[GH-5093](https://github.com/hashicorp/vault/pull/5093)]
 * secrets/aws: Make credential types more explicit [[GH-4360](https://github.com/hashicorp/vault/pull/4360)]
 * secrets/nomad: Support for longer token names [[GH-5117](https://github.com/hashicorp/vault/pull/5117)]
 * secrets/pki: Allow disabling CRL generation [[GH-5134](https://github.com/hashicorp/vault/pull/5134)]
 * storage/azure: Add support for different Azure environments [[GH-4997](https://github.com/hashicorp/vault/pull/4997)]
 * storage/file: Sort keys in list responses [[GH-5141](https://github.com/hashicorp/vault/pull/5141)]
 * storage/mysql: Support special characters in database and table names.

BUG FIXES:

 * auth/jwt: Always validate `aud` claim even if `bound_audiences` isn't set
   (IOW, error in this case)
 * core: Prevent Go's HTTP library from interspersing logs in a different
   format and/or interleaved [[GH-5135](https://github.com/hashicorp/vault/pull/5135)]
 * identity: Properly populate `mount_path` and `mount_type` on group lookup
   [[GH-5074](https://github.com/hashicorp/vault/pull/5074)]
 * identity: Fix persisting alias metadata [[GH-5188](https://github.com/hashicorp/vault/pull/5188)]
 * identity: Fix carryover issue from previously fixed race condition that
   could cause Vault not to start up due to two entities referencing the same
   alias. These entities are now merged. [[GH-5000](https://github.com/hashicorp/vault/pull/5000)]
 * replication: Fix issue causing some pages not to flush to storage
 * secrets/database: Fix inability to update custom SQL statements on
   database roles. [[GH-5080](https://github.com/hashicorp/vault/pull/5080)]
 * secrets/pki: Disallow putting the CA's serial on its CRL. While technically
   legal, doing so inherently means the CRL can't be trusted anyways, so it's
   not useful and easy to footgun. [[GH-5134](https://github.com/hashicorp/vault/pull/5134)]
 * storage/gcp,spanner: Fix data races [[GH-5081](https://github.com/hashicorp/vault/pull/5081)]

## 0.10.4 (July 25th, 2018)

SECURITY:

 * Control Groups: The associated Identity entity with a request was not being
   properly persisted. As a result, the same authorizer could provide more than
   one authorization.

DEPRECATIONS/CHANGES:

 * Revocations of dynamic secrets leases are now queued/asynchronous rather
   than synchronous. This allows Vault to take responsibility for revocation
   even if the initial attempt fails. The previous synchronous behavior can be
   attained via the `-sync` CLI flag or `sync` API parameter. When in
   synchronous mode, if the operation results in failure it is up to the user
   to retry.
 * CLI Retries: The CLI will no longer retry commands on 5xx errors. This was a
   source of confusion to users as to why Vault would "hang" before returning a
   5xx error. The Go API client still defaults to two retries.
 * Identity Entity Alias metadata: You can no longer manually set metadata on
   entity aliases. All alias data (except the canonical entity ID it refers to)
   is intended to be managed by the plugin providing the alias information, so
   allowing it to be set manually didn't make sense.

FEATURES:

 * **JWT/OIDC Auth Method**: The new `jwt` auth method accepts JWTs and either
   validates signatures locally or uses OIDC Discovery to fetch the current set
   of keys for signature validation. Various claims can be specified for
   validation (in addition to the cryptographic signature) and a user and
   optional groups claim can be used to provide Identity information.
 * **FoundationDB Storage**: You can now use FoundationDB for storing Vault
   data.
 * **UI Control Group Workflow (enterprise)**: The UI will now detect control
   group responses and provides a workflow to view the status of the request
   and to authorize requests.
 * **Vault Agent (Beta)**: Vault Agent is a daemon that can automatically
   authenticate for you across a variety of authentication methods, provide
   tokens to clients, and keep the tokens renewed, reauthenticating as
   necessary.

IMPROVEMENTS:

 * auth/azure: Add support for virtual machine scale sets
 * auth/gcp: Support multiple bindings for region, zone, and instance group
 * cli: Add subcommands for interacting with the plugin catalog [[GH-4911](https://github.com/hashicorp/vault/pull/4911)]
 * cli: Add a `-description` flag to secrets and auth tune subcommands to allow
   updating an existing secret engine's or auth method's description. This
   change also allows the description to be unset by providing an empty string.
 * core: Add config flag to disable non-printable character check [[GH-4917](https://github.com/hashicorp/vault/pull/4917)]
 * core: A `max_request_size` parameter can now be set per-listener to adjust
   the maximum allowed size per request [[GH-4824](https://github.com/hashicorp/vault/pull/4824)]
 * core: Add control group request endpoint to default policy [[GH-4904](https://github.com/hashicorp/vault/pull/4904)]
 * identity: Identity metadata is now passed through to plugins [[GH-4967](https://github.com/hashicorp/vault/pull/4967)]
 * replication: Add additional saftey checks and logging when replication is
   in a bad state
 * secrets/kv: Add support for using `-field=data` to KVv2 when using `vault
   kv` [[GH-4895](https://github.com/hashicorp/vault/pull/4895)]
 * secrets/pki: Add the ability to tidy revoked but unexpired certificates
   [[GH-4916](https://github.com/hashicorp/vault/pull/4916)]
 * secrets/ssh: Allow Vault to work with single-argument SSH flags [[GH-4825](https://github.com/hashicorp/vault/pull/4825)]
 * secrets/ssh: SSH executable path can now be configured in the CLI [[GH-4937](https://github.com/hashicorp/vault/pull/4937)]
 * storage/swift: Add additional configuration options [[GH-4901](https://github.com/hashicorp/vault/pull/4901)]
 * ui: Choose which auth methods to show to unauthenticated users via
   `listing_visibility` in the auth method edit forms [[GH-4854](https://github.com/hashicorp/vault/pull/4854)]
 * ui: Authenticate users automatically by passing a wrapped token to the UI via
   the new `wrapped_token` query parameter [[GH-4854](https://github.com/hashicorp/vault/pull/4854)]

BUG FIXES:

 * api: Fix response body being cleared too early [[GH-4987](https://github.com/hashicorp/vault/pull/4987)]
 * auth/approle: Fix issue with tidy endpoint that would unnecessarily remove
   secret accessors [[GH-4981](https://github.com/hashicorp/vault/pull/4981)]
 * auth/aws: Fix updating `max_retries` [[GH-4980](https://github.com/hashicorp/vault/pull/4980)]
 * auth/kubernetes: Trim trailing whitespace when sending JWT
 * cli: Fix parsing of environment variables for integer flags [[GH-4925](https://github.com/hashicorp/vault/pull/4925)]
 * core: Fix returning 500 instead of 503 if a rekey is attempted when Vault is
   sealed [[GH-4874](https://github.com/hashicorp/vault/pull/4874)]
 * core: Fix issue releasing the leader lock in some circumstances [[GH-4915](https://github.com/hashicorp/vault/pull/4915)]
 * core: Fix a panic that could happen if the server was shut down while still
   starting up
 * core: Fix deadlock that would occur if a leadership loss occurs at the same
   time as a seal operation [[GH-4932](https://github.com/hashicorp/vault/pull/4932)]
 * core: Fix issue with auth mounts failing to renew tokens due to policies
   changing [[GH-4960](https://github.com/hashicorp/vault/pull/4960)]
 * auth/radius: Fix issue where some radius logins were being canceled too early
   [[GH-4941](https://github.com/hashicorp/vault/pull/4941)]
 * core: Fix accidental seal of vault of we lose leadership during startup
   [[GH-4924](https://github.com/hashicorp/vault/pull/4924)]
 * core: Fix standby not being able to forward requests larger than 4MB
   [[GH-4844](https://github.com/hashicorp/vault/pull/4844)]
 * core: Avoid panic while processing group memberships [[GH-4841](https://github.com/hashicorp/vault/pull/4841)]
 * identity: Fix a race condition creating aliases [[GH-4965](https://github.com/hashicorp/vault/pull/4965)]
 * plugins: Fix being unable to send very large payloads to or from plugins
   [[GH-4958](https://github.com/hashicorp/vault/pull/4958)]
 * physical/azure: Long list responses would sometimes be truncated [[GH-4983](https://github.com/hashicorp/vault/pull/4983)]
 * replication: Allow replication status requests to be processed while in
   merkle sync
 * replication: Ensure merkle reindex flushes all changes to storage immediately
 * replication: Fix a case where a network interruption could cause a secondary
   to be unable to reconnect to a primary
 * secrets/pki: Fix permitted DNS domains performing improper validation
   [[GH-4863](https://github.com/hashicorp/vault/pull/4863)]
 * secrets/database: Fix panic during DB creds revocation [[GH-4846](https://github.com/hashicorp/vault/pull/4846)]
 * ui: Fix usage of cubbyhole backend in the UI [[GH-4851](https://github.com/hashicorp/vault/pull/4851)]
 * ui: Fix toggle state when a secret is JSON-formatted [[GH-4913](https://github.com/hashicorp/vault/pull/4913)]
 * ui: Fix coercion of falsey values to empty string when editing secrets as
   JSON [[GH-4977](https://github.com/hashicorp/vault/pull/4977)]

## 0.10.3 (June 20th, 2018)

DEPRECATIONS/CHANGES:

 * In the audit log and in client responses, policies are now split into three
   parameters: policies that came only from tokens, policies that came only
   from Identity, and the combined set. Any previous location of policies via
   the API now contains the full, combined set.
 * When a token is tied to an Identity entity and the entity is deleted, the
   token will no longer be usable, regardless of the validity of the token
   itself.
 * When authentication succeeds but no policies were defined for that specific
   user, most auth methods would allow a token to be generated but a few would
   reject the authentication, namely `ldap`, `okta`, and `radius`. Since the
   `default` policy is added by Vault's core, this would incorrectly reject
   valid authentications before they would in fact be granted policies. This
   inconsistency has been addressed; valid authentications for these methods
   now succeed even if no policy was specifically defined in that method for
   that user.

FEATURES:

 * Root Rotation for Active Directory: You can now command Vault to rotate the
   configured root credentials used in the AD secrets engine, to ensure that
   only Vault knows the credentials it's using.
 * URI SANs in PKI: You can now configure URI Subject Alternate Names in the
   `pki` backend. Roles can limit which SANs are allowed via globbing.
 * `kv rollback` Command: You can now use `vault kv rollback` to roll a KVv2
   path back to a previous non-deleted/non-destroyed version. The previous
   version becomes the next/newest version for the path.
 * Token Bound CIDRs in AppRole: You can now add CIDRs to which a token
   generated from AppRole will be bound.

IMPROVEMENTS:

 * approle: Return 404 instead of 202 on invalid role names during POST
   operations [[GH-4778](https://github.com/hashicorp/vault/pull/4778)]
 * core: Add idle and initial header read/TLS handshake timeouts to connections
   to ensure server resources are cleaned up [[GH-4760](https://github.com/hashicorp/vault/pull/4760)]
 * core: Report policies in token, identity, and full sets [[GH-4747](https://github.com/hashicorp/vault/pull/4747)]
 * secrets/databases: Add `create`/`update` distinction for connection
   configurations [[GH-3544](https://github.com/hashicorp/vault/pull/3544)]
 * secrets/databases: Add `create`/`update` distinction for role configurations
   [[GH-3544](https://github.com/hashicorp/vault/pull/3544)]
 * secrets/databases: Add best-effort revocation logic for use when a role has
   been deleted [[GH-4782](https://github.com/hashicorp/vault/pull/4782)]
 * secrets/kv: Add `kv rollback` [[GH-4774](https://github.com/hashicorp/vault/pull/4774)]
 * secrets/pki: Add URI SANs support [[GH-4675](https://github.com/hashicorp/vault/pull/4675)]
 * secrets/ssh: Allow standard SSH command arguments to be used, without
   requiring username@hostname syntax [[GH-4710](https://github.com/hashicorp/vault/pull/4710)]
 * storage/consul: Add context support so that requests are cancelable
   [[GH-4739](https://github.com/hashicorp/vault/pull/4739)]
 * sys: Added `hidden` option to `listing_visibility` field on `sys/mounts`
   API [[GH-4827](https://github.com/hashicorp/vault/pull/4827)]
 * ui: Secret values are obfuscated by default and visibility is toggleable [[GH-4422](https://github.com/hashicorp/vault/pull/4422)]

BUG FIXES:

 * auth/approle: Fix panic due to metadata being nil [[GH-4719](https://github.com/hashicorp/vault/pull/4719)]
 * auth/aws: Fix delete path for tidy operations [[GH-4799](https://github.com/hashicorp/vault/pull/4799)]
 * core: Optimizations to remove some speed regressions due to the
   security-related changes in 0.10.2
 * storage/dynamodb: Fix errors seen when reading existing DynamoDB data [[GH-4721](https://github.com/hashicorp/vault/pull/4721)]
 * secrets/database: Fix default MySQL root rotation statement [[GH-4748](https://github.com/hashicorp/vault/pull/4748)]
 * secrets/gcp: Fix renewal for GCP account keys
 * secrets/kv: Fix writing to the root of a KVv2 mount from `vault kv` commands
   incorrectly operating on a root+mount path instead of being an error
   [[GH-4726](https://github.com/hashicorp/vault/pull/4726)]
 * seal/pkcs11: Add `CKK_SHA256_HMAC` to the search list when finding HMAC
   keys, fixing lookup on some Thales devices
 * replication: Fix issue enabling replication when a non-auth mount and auth
   mount have the same name
 * auth/kubernetes: Fix issue verifying ECDSA signed JWTs
 * ui: add missing edit mode for auth method configs [[GH-4770](https://github.com/hashicorp/vault/pull/4770)]

## 0.10.2 (June 6th, 2018)

SECURITY:

 * Tokens: A race condition was identified that could occur if a token's
   lease expired while Vault was not running. In this case, when Vault came
   back online, sometimes it would properly revoke the lease but other times it
   would not, leading to a Vault token that no longer had an expiration and had
   essentially unlimited lifetime. This race was per-token, not all-or-nothing
   for all tokens that may have expired during Vault's downtime. We have fixed
   the behavior and put extra checks in place to help prevent any similar
   future issues. In addition, the logic we have put in place ensures that such
   lease-less tokens can no longer be used (unless they are root tokens that
   never had an expiration to begin with).
 * Convergent Encryption: The version 2 algorithm used in `transit`'s
   convergent encryption feature is susceptible to offline
   plaintext-confirmation attacks. As a result, we are introducing a version 3
   algorithm that mitigates this. If you are currently using convergent
   encryption, we recommend upgrading, rotating your encryption key (the new
   key version will use the new algorithm), and rewrapping your data (the
   `rewrap` endpoint can be used to allow a relatively non-privileged user to
   perform the rewrapping while never divulging the plaintext).
 * AppRole case-sensitive role name secret-id leaking: When using a mixed-case
   role name via AppRole, deleting a secret-id via accessor or other operations
   could end up leaving the secret-id behind and valid but without an accessor.
   This has now been fixed, and we have put checks in place to prevent these
   secret-ids from being used.

DEPRECATIONS/CHANGES:

 * PKI duration return types: The PKI backend now returns durations (e.g. when
   reading a role) as an integer number of seconds instead of a Go-style
   string, in line with how the rest of Vault's API returns durations.

FEATURES:

 * Active Directory Secrets Engine: A new `ad` secrets engine has been created
   which allows Vault to rotate and provide credentials for configured AD
   accounts.
 * Rekey Verification: Rekey operations can now require verification. This
   turns on a two-phase process where the existing key shares authorize
   generating a new master key, and a threshold of the new, returned key shares
   must be provided to verify that they have been successfully received in
   order for the actual master key to be rotated.
 * CIDR restrictions for `cert`, `userpass`, and `kubernetes` auth methods:
   You can now limit authentication to specific CIDRs; these will also be
   encoded in resultant tokens to limit their use.
 * Vault UI Browser CLI: The UI now supports usage of read/write/list/delete
   commands in a CLI that can be accessed from the nav bar. Complex inputs such
   as JSON files are not currently supported. This surfaces features otherwise
   unsupported in Vault's UI.
 * Azure Key Vault Auto Unseal/Seal Wrap Support (Enterprise): Azure Key Vault
   can now be used a support seal for Auto Unseal and Seal Wrapping.

IMPROVEMENTS:

 * api: Close renewer's doneCh when the renewer is stopped, so that programs
   expecting a final value through doneCh behave correctly [[GH-4472](https://github.com/hashicorp/vault/pull/4472)]
 * auth/cert: Break out `allowed_names` into component parts and add
   `allowed_uri_sans` [[GH-4231](https://github.com/hashicorp/vault/pull/4231)]
 * auth/ldap: Obfuscate error messages pre-bind for greater security [[GH-4700](https://github.com/hashicorp/vault/pull/4700)]
 * cli: `vault login` now supports a `-no-print` flag to suppress printing
   token information but still allow storing into the token helper [[GH-4454](https://github.com/hashicorp/vault/pull/4454)]
 * core/pkcs11 (enterprise): Add support for CKM_AES_CBC_PAD, CKM_RSA_PKCS, and
   CKM_RSA_PKCS_OAEP mechanisms
 * core/pkcs11 (enterprise): HSM slots can now be selected by token label
   instead of just slot number
 * core/token: Optimize token revocation by removing unnecessary list call
   against the storage backend when calling revoke-orphan on tokens [[GH-4465](https://github.com/hashicorp/vault/pull/4465)]
 * core/token: Refactor token revocation logic to not block on the call when
   underlying leases are pending revocation by moving the expiration logic to
   the expiration manager [[GH-4512](https://github.com/hashicorp/vault/pull/4512)]
 * expiration: Allow revoke-prefix and revoke-force to work on single leases as
   well as prefixes [[GH-4450](https://github.com/hashicorp/vault/pull/4450)]
 * identity: Return parent group info when reading a group [[GH-4648](https://github.com/hashicorp/vault/pull/4648)]
 * identity: Provide more contextual key information when listing entities,
   groups, and aliases
 * identity: Passthrough EntityID to backends [[GH-4663](https://github.com/hashicorp/vault/pull/4663)]
 * identity: Adds ability to request entity information through system view
   [GH_4681]
 * secret/pki: Add custom extended key usages [[GH-4667](https://github.com/hashicorp/vault/pull/4667)]
 * secret/pki: Add custom PKIX serial numbers [[GH-4694](https://github.com/hashicorp/vault/pull/4694)]
 * secret/ssh: Use hostname instead of IP in OTP mode, similar to CA mode
   [[GH-4673](https://github.com/hashicorp/vault/pull/4673)]
 * storage/file: Attempt in some error conditions to do more cleanup [[GH-4684](https://github.com/hashicorp/vault/pull/4684)]
 * ui: wrapping lookup now distplays the path [[GH-4644](https://github.com/hashicorp/vault/pull/4644)]
 * ui: Identity interface now has more inline actions to make editing and adding
   aliases to an entity or group easier [[GH-4502](https://github.com/hashicorp/vault/pull/4502)]
 * ui: Identity interface now lists groups by name [[GH-4655](https://github.com/hashicorp/vault/pull/4655)]
 * ui: Permission denied errors still render the sidebar in the Access section
   [[GH-4658](https://github.com/hashicorp/vault/pull/4658)]
 * replication: Improve performance of index page flushes and WAL garbage
   collecting

BUG FIXES:

 * auth/approle: Make invalid role_id a 400 error instead of 500 [[GH-4470](https://github.com/hashicorp/vault/pull/4470)]
 * auth/cert: Fix Identity alias using serial number instead of common name
   [[GH-4475](https://github.com/hashicorp/vault/pull/4475)]
 * cli: Fix panic running `vault token capabilities` with multiple paths
   [[GH-4552](https://github.com/hashicorp/vault/pull/4552)]
 * core: When using the `use_always` option with PROXY protocol support, do not
   require `authorized_addrs` to be set [[GH-4065](https://github.com/hashicorp/vault/pull/4065)]
 * core: Fix panic when certain combinations of policy paths and allowed/denied
   parameters were used [[GH-4582](https://github.com/hashicorp/vault/pull/4582)]
 * secret/gcp: Make `bound_region` able to use short names
 * secret/kv: Fix response wrapping for KV v2 [[GH-4511](https://github.com/hashicorp/vault/pull/4511)]
 * secret/kv: Fix address flag not being honored correctly [[GH-4617](https://github.com/hashicorp/vault/pull/4617)]
 * secret/pki: Fix `safety_buffer` for tidy being allowed to be negative,
   clearing all certs [[GH-4641](https://github.com/hashicorp/vault/pull/4641)]
 * secret/pki: Fix `key_type` not being allowed to be set to `any` [[GH-4595](https://github.com/hashicorp/vault/pull/4595)]
 * secret/pki: Fix path length parameter being ignored when using
   `use_csr_values` and signing an intermediate CA cert [[GH-4459](https://github.com/hashicorp/vault/pull/4459)]
 * secret/ssh: Only append UserKnownHostsFile to args when configured with a
   value [[GH-4674](https://github.com/hashicorp/vault/pull/4674)]
 * storage/dynamodb: Fix listing when one child is left within a nested path
   [[GH-4570](https://github.com/hashicorp/vault/pull/4570)]
 * storage/gcs: Fix swallowing an error on connection close [[GH-4691](https://github.com/hashicorp/vault/pull/4691)]
 * ui: Fix HMAC algorithm in transit [[GH-4604](https://github.com/hashicorp/vault/pull/4604)]
 * ui: Fix unwrap of auth responses via the UI's unwrap tool [[GH-4611](https://github.com/hashicorp/vault/pull/4611)]
 * ui (enterprise): Fix parsing of version string that blocked some users from seeing
   enterprise-specific pages in the UI [[GH-4547](https://github.com/hashicorp/vault/pull/4547)]
 * ui: Fix incorrect capabilities path check when viewing policies [[GH-4566](https://github.com/hashicorp/vault/pull/4566)]
 * replication: Fix error while running plugins on a newly created replication
   secondary
 * replication: Fix issue with token store lookups after a secondary's mount table
   is invalidated.
 * replication: Improve startup time when a large merkle index is in use.
 * replication: Fix panic when storage becomes unreachable during unseal.

## 0.10.1/0.9.7 (April 25th, 2018)

The following two items are in both 0.9.7 and 0.10.1. They only affect
Enterprise, and as such 0.9.7 is an Enterprise-only release:

SECURITY:

 * EGPs: A regression affecting 0.9.6 and 0.10.0 causes EGPs to not be applied
   correctly if an EGP is updated in a running Vault after initial write or
   after it is loaded on unseal. This has been fixed.

BUG FIXES:

 * Fixed an upgrade issue affecting performance secondaries when migrating from
   a version that did not include Identity to one that did.

All other content in this release is for 0.10.1 only.

DEPRECATIONS/CHANGES:

 * `vault kv` and Vault versions: In 0.10.1 some issues with `vault kv` against
   v1 K/V engine mounts are fixed. However, using 0.10.1 for both the server
   and CLI versions is required.
 * Mount information visibility: Users that have access to any path within a
   mount can now see information about that mount, such as its type and
   options, via some API calls.
 * Identity and Local Mounts: Local mounts would allow creating Identity
   entities but these would not be able to be used successfully (even locally)
   in replicated scenarios. We have now disallowed entities and groups from
   being created for local mounts in the first place.

FEATURES:

 * X-Forwarded-For support: `X-Forwarded-For` headers can now be used to set the
   client IP seen by Vault. See the [TCP listener configuration
   page](https://www.vaultproject.io/docs/configuration/listener/tcp.html) for
   details.
 * CIDR IP Binding for Tokens: Tokens now support being bound to specific
   CIDR(s) for usage. Currently this is implemented in Token Roles; usage can be
   expanded to other authentication backends over time.
 * `vault kv patch` command: A new `kv patch` helper command that allows
   modifying only some values in existing data at a K/V path, but uses
   check-and-set to ensure that this modification happens safely.
 * AppRole Local Secret IDs: Roles can now be configured to generate secret IDs
   local to the cluster. This enables performance secondaries to generate and
   consume secret IDs without contacting the primary.
 * AES-GCM Support for PKCS#11 [BETA] (Enterprise): For supporting HSMs,
   AES-GCM can now be used in lieu of AES-CBC/HMAC-SHA256. This has currently
   only been fully tested on AWS CloudHSM.
 * Auto Unseal/Seal Wrap Key Rotation Support (Enterprise): Auto Unseal
   mechanisms, including PKCS#11 HSMs, now support rotation of encryption keys,
   and migration between key and encryption types, such as from AES-CBC to
   AES-GCM, can be performed at the same time (where supported).

IMPROVEMENTS:

 * auth/approle: Support for cluster local secret IDs. This enables secondaries
   to generate secret IDs without contacting the primary [[GH-4427](https://github.com/hashicorp/vault/pull/4427)]
 * auth/token: Add to the token lookup response, the policies inherited due to
   identity associations [[GH-4366](https://github.com/hashicorp/vault/pull/4366)]
 * auth/token: Add CIDR binding to token roles [[GH-815](https://github.com/hashicorp/vault/pull/815)]
 * cli: Add `vault kv patch` [[GH-4432](https://github.com/hashicorp/vault/pull/4432)]
 * core: Add X-Forwarded-For support [[GH-4380](https://github.com/hashicorp/vault/pull/4380)]
 * core: Add token CIDR-binding support [[GH-815](https://github.com/hashicorp/vault/pull/815)]
 * identity: Add the ability to disable an entity. Disabling an entity does not
   revoke associated tokens, but while the entity is disabled they cannot be
   used. [[GH-4353](https://github.com/hashicorp/vault/pull/4353)]
 * physical/consul: Allow tuning of session TTL and lock wait time [[GH-4352](https://github.com/hashicorp/vault/pull/4352)]
 * replication: Dynamically adjust WAL cleanup over a period of time based on
   the rate of writes committed
 * secret/ssh: Update dynamic key install script to use shell locking to avoid
   concurrent modifications [[GH-4358](https://github.com/hashicorp/vault/pull/4358)]
 * ui: Access to `sys/mounts` is no longer needed to use the UI - the list of
   engines will show you the ones you implicitly have access to (because you have
   access to to secrets in those engines) [[GH-4439](https://github.com/hashicorp/vault/pull/4439)]

BUG FIXES:

 * cli: Fix `vault kv` backwards compatibility with KV v1 engine mounts
   [[GH-4430](https://github.com/hashicorp/vault/pull/4430)]
 * identity: Persist entity memberships in external identity groups across
   mounts [[GH-4365](https://github.com/hashicorp/vault/pull/4365)]
 * identity: Fix error preventing authentication using local mounts on
   performance secondary replication clusters [[GH-4407](https://github.com/hashicorp/vault/pull/4407)]
 * replication: Fix issue causing secondaries to not connect properly to a
   pre-0.10 primary until the primary was upgraded
 * secret/gcp: Fix panic on rollback when a roleset wasn't created properly
   [[GH-4344](https://github.com/hashicorp/vault/pull/4344)]
 * secret/gcp: Fix panic on renewal
 * ui: Fix IE11 form submissions in a few parts of the application [[GH-4378](https://github.com/hashicorp/vault/pull/4378)]
 * ui: Fix IE file saving on policy pages and init screens [[GH-4376](https://github.com/hashicorp/vault/pull/4376)]
 * ui: Fixed an issue where the AWS secret backend would show the wrong menu
   [[GH-4371](https://github.com/hashicorp/vault/pull/4371)]
 * ui: Fixed an issue where policies with commas would not render in the
   interface properly [[GH-4398](https://github.com/hashicorp/vault/pull/4398)]
 * ui: Corrected the saving of mount tune ttls for auth methods [[GH-4431](https://github.com/hashicorp/vault/pull/4431)]
 * ui: Credentials generation no longer checks capabilities before making
   api calls. This should fix needing "update" capabilites to read IAM
   credentials in the AWS secrets engine [[GH-4446](https://github.com/hashicorp/vault/pull/4446)]

## 0.10.0 (April 10th, 2018)

SECURITY:

 * Log sanitization for Combined Database Secret Engine: In certain failure
   scenarios with incorrectly formatted connection urls, the raw connection
   errors were being returned to the user with the configured database
   credentials. Errors are now sanitized before being returned to the user.

DEPRECATIONS/CHANGES:

 * Database plugin compatibility: The database plugin interface was enhanced to
   support some additional functionality related to root credential rotation
   and supporting templated URL strings. The changes were made in a
   backwards-compatible way and all builtin plugins were updated with the new
   features. Custom plugins not built into Vault will need to be upgraded to
   support templated URL strings and root rotation. Additionally, the
   Initialize method was deprecated in favor of a new Init method that supports
   configuration modifications that occur in the plugin back to the primary
   data store.
 * Removal of returned secret information: For a long time Vault has returned
   configuration given to various secret engines and auth methods with secret
   values (such as secret API keys or passwords) still intact, and with a
   warning to the user on write that anyone with read access could see the
   secret. This was mostly done to make it easy for tools like Terraform to
   judge whether state had drifted. However, it also feels quite un-Vault-y to
   do this and we've never felt very comfortable doing so. In 0.10 we have gone
   through and removed this behavior from the various backends; fields which
   contained secret values are simply no longer returned on read. We are
   working with the Terraform team to make changes to their provider to
   accommodate this as best as possible, and users of other tools may have to
   make adjustments, but in the end we felt that the ends did not justify the
   means and we needed to prioritize security over operational convenience.
 * LDAP auth method case sensitivity: We now treat usernames and groups
   configured locally for policy assignment in a case insensitive fashion by
   default. Existing configurations will continue to work as they do now;
   however, the next time a configuration is written `case_sensitive_names`
   will need to be explicitly set to `true`.
 * TTL handling within core: All lease TTL handling has been centralized within
   the core of Vault to ensure consistency across all backends. Since this was
   previously delegated to individual backends, there may be some slight
   differences in TTLs generated from some backends.
 * Removal of default `secret/` mount: In 0.12 we will stop mounting `secret/`
   by default at initialization time (it will still be available in `dev`
   mode).

FEATURES:

 * OSS UI: The Vault UI is now fully open-source. Similarly to the CLI, some
   features are only available with a supporting version of Vault, but the code
   base is entirely open.
 * Versioned K/V: The `kv` backend has been completely revamped, featuring
   flexible versioning of values, check-and-set protections, and more. A new
   `vault kv` subcommand allows friendly interactions with it. Existing mounts
   of the `kv` backend can be upgraded to the new versioned mode (downgrades
   are not currently supported). The old "passthrough" mode is still the
   default for new mounts; versioning can be turned on by setting the
   `-version=2` flag for the `vault secrets enable` command.
 * Database Root Credential Rotation: Database configurations can now rotate
   their own configured admin/root credentials, allowing configured credentials
   for a database connection to be rotated immediately after sending them into
   Vault, invalidating the old credentials and ensuring only Vault knows the
   actual valid values.
 * Azure Authentication Plugin: There is now a plugin (pulled in to Vault) that
   allows authenticating Azure machines to Vault using Azure's Managed Service
   Identity credentials. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-auth-azure) for more
   information.
 * GCP Secrets Plugin: There is now a plugin (pulled in to Vault) that allows
   generating secrets to allow access to GCP. See the [plugin
   repository](https://github.com/hashicorp/vault-plugin-secrets-gcp) for more
   information.
 * Selective Audit HMACing of Request and Response Data Keys: HMACing in audit
   logs can be turned off for specific keys in the request input map and
   response `data` map on a per-mount basis.
 * Passthrough Request Headers: Request headers can now be selectively passed
   through to backends on a per-mount basis. This is useful in various cases
   when plugins are interacting with external services.
 * HA for Google Cloud Storage: The GCS storage type now supports HA.
 * UI support for identity: Add and edit entities, groups, and their associated
   aliases.
 * UI auth method support: Enable, disable, and configure all of the built-in
   authentication methods.
 * UI (Enterprise): View and edit Sentinel policies.

IMPROVEMENTS:

 * core: Centralize TTL generation for leases in core [[GH-4230](https://github.com/hashicorp/vault/pull/4230)]
 * identity: API to update group-alias by ID [[GH-4237](https://github.com/hashicorp/vault/pull/4237)]
 * secret/cassandra: Update Cassandra storage delete function to not use batch
   operations [[GH-4054](https://github.com/hashicorp/vault/pull/4054)]
 * storage/mysql: Allow setting max idle connections and connection lifetime
   [[GH-4211](https://github.com/hashicorp/vault/pull/4211)]
 * storage/gcs: Add HA support [[GH-4226](https://github.com/hashicorp/vault/pull/4226)]
 * ui: Add Nomad to the list of available secret engines
 * ui: Adds ability to set static headers to be returned by the UI

BUG FIXES:

 * api: Fix retries not working [[GH-4322](https://github.com/hashicorp/vault/pull/4322)]
 * auth/gcp: Invalidate clients on config change
 * auth/token: Revoke-orphan and tidy operations now correctly cleans up the
   parent prefix entry in the underlying storage backend. These operations also
   mark corresponding child tokens as orphans by removing the parent/secondary
   index from the entries. [[GH-4193](https://github.com/hashicorp/vault/pull/4193)]
 * command: Re-add `-mfa` flag and migrate to OSS binary [[GH-4223](https://github.com/hashicorp/vault/pull/4223)]
 * core: Fix issue occurring from mounting two auth backends with the same path
   with one mount having `auth/` in front [[GH-4206](https://github.com/hashicorp/vault/pull/4206)]
 * mfa: Invalidation of MFA configurations (Enterprise)
 * replication: Fix a panic on some non-64-bit platforms
 * replication: Fix invalidation of policies on performance secondaries
 * secret/pki: When tidying if a value is unexpectedly nil, delete it and move
   on [[GH-4214](https://github.com/hashicorp/vault/pull/4214)]
 * storage/s3: Fix panic if S3 returns no Content-Length header [[GH-4222](https://github.com/hashicorp/vault/pull/4222)]
 * ui: Fixed an issue where the UI was checking incorrect paths when operating
   on transit keys. Capabilities are now checked when attempting to encrypt /
   decrypt, etc.
 * ui: Fixed IE 11 layout issues and JS errors that would stop the application
   from running.
 * ui: Fixed the link that gets rendered when a user doesn't have permissions
   to view the root of a secret engine. The link now sends them back to the list
   of secret engines.
 * replication: Fix issue with DR secondaries when using mount specified local
   paths.
 * cli: Fix an issue where generating a dr operation token would not output the
   token [[GH-4328](https://github.com/hashicorp/vault/pull/4328)]

## 0.9.6 (March 20th, 2018)

DEPRECATIONS/CHANGES:

 * The AWS authentication backend now allows binds for inputs as either a
   comma-delimited string or a string array. However, to keep consistency with
   input and output, when reading a role the binds will now be returned as
   string arrays rather than strings.
 * In order to prefix-match IAM role and instance profile ARNs in AWS auth
   backend, you now must explicitly opt-in by adding a `*` to the end of the
   ARN. Existing configurations will be upgraded automatically, but when
   writing a new role configuration the updated behavior will be used.

FEATURES:

 * Replication Activation Enhancements: When activating a replication
   secondary, a public key can now be fetched first from the target cluster.
   This public key can be provided to the primary when requesting the
   activation token. If provided, the public key will be used to perform a
   Diffie-Hellman key exchange resulting in a shared key that encrypts the
   contents of the activation token. The purpose is to protect against
   accidental disclosure of the contents of the token if unwrapped by the wrong
   party, given that the contents of the token are highly sensitive. If
   accidentally unwrapped, the contents of the token are not usable by the
   unwrapping party. It is important to note that just as a malicious operator
   could unwrap the contents of the token, a malicious operator can pretend to
   be a secondary and complete the Diffie-Hellman exchange on their own; this
   feature provides defense in depth but still requires due diligence around
   replication activation, including multiple eyes on the commands/tokens and
   proper auditing.

IMPROVEMENTS:

 * api: Update renewer grace period logic. It no longer is static, but rather
   dynamically calculates one based on the current lease duration after each
   renew. [[GH-4090](https://github.com/hashicorp/vault/pull/4090)]
 * auth/approle: Allow array input for bound_cidr_list [4078]
 * auth/aws: Allow using lists in role bind parameters [[GH-3907](https://github.com/hashicorp/vault/pull/3907)]
 * auth/aws: Allow binding by EC2 instance IDs [[GH-3816](https://github.com/hashicorp/vault/pull/3816)]
 * auth/aws: Allow non-prefix-matched IAM role and instance profile ARNs
   [[GH-4071](https://github.com/hashicorp/vault/pull/4071)]
 * auth/ldap: Set a very large size limit on queries [[GH-4169](https://github.com/hashicorp/vault/pull/4169)]
 * core: Log info notifications of revoked leases for all leases/reasons, not
   just expirations [[GH-4164](https://github.com/hashicorp/vault/pull/4164)]
 * physical/couchdb: Removed limit on the listing of items [[GH-4149](https://github.com/hashicorp/vault/pull/4149)]
 * secret/pki: Support certificate policies [[GH-4125](https://github.com/hashicorp/vault/pull/4125)]
 * secret/pki: Add ability to have CA:true encoded into intermediate CSRs, to
   improve compatibility with some ADFS scenarios [[GH-3883](https://github.com/hashicorp/vault/pull/3883)]
 * secret/transit: Allow selecting signature algorithm as well as hash
   algorithm when signing/verifying [[GH-4018](https://github.com/hashicorp/vault/pull/4018)]
 * server: Make sure `tls_disable_client_cert` is actually a true value rather
   than just set [[GH-4049](https://github.com/hashicorp/vault/pull/4049)]
 * storage/dynamodb: Allow specifying max retries for dynamo client [[GH-4115](https://github.com/hashicorp/vault/pull/4115)]
 * storage/gcs: Allow specifying chunk size for transfers, which can reduce
   memory utilization [[GH-4060](https://github.com/hashicorp/vault/pull/4060)]
 * sys/capabilities: Add the ability to use multiple paths for capability
   checking [[GH-3663](https://github.com/hashicorp/vault/pull/3663)]

BUG FIXES:

 * auth/aws: Fix honoring `max_ttl` when a corresponding role `ttl` is not also
   set [[GH-4107](https://github.com/hashicorp/vault/pull/4107)]
 * auth/okta: Fix honoring configured `max_ttl` value [[GH-4110](https://github.com/hashicorp/vault/pull/4110)]
 * auth/token: If a periodic token being issued has a period greater than the
   max_lease_ttl configured on the token store mount, truncate it. This matches
   renewal behavior; before it was inconsistent between issuance and renewal.
   [[GH-4112](https://github.com/hashicorp/vault/pull/4112)]
 * cli: Improve error messages around `vault auth help` when there is no CLI
   helper for a particular method [[GH-4056](https://github.com/hashicorp/vault/pull/4056)]
 * cli: Fix autocomplete installation when using Fish as the shell [[GH-4094](https://github.com/hashicorp/vault/pull/4094)]
 * secret/database: Properly honor mount-tuned max TTL [[GH-4051](https://github.com/hashicorp/vault/pull/4051)]
 * secret/ssh: Return `key_bits` value when reading a role [[GH-4098](https://github.com/hashicorp/vault/pull/4098)]
 * sys: When writing policies on a performance replication secondary, properly
   forward requests to the primary [[GH-4129](https://github.com/hashicorp/vault/pull/4129)]

## 0.9.5 (February 26th, 2018)

IMPROVEMENTS:

 * auth: Allow sending default_lease_ttl and max_lease_ttl values when enabling
   auth methods. [[GH-4019](https://github.com/hashicorp/vault/pull/4019)]
 * secret/database: Add list functionality to `database/config` endpoint
   [[GH-4026](https://github.com/hashicorp/vault/pull/4026)]
 * physical/consul: Allow setting a specific service address [[GH-3971](https://github.com/hashicorp/vault/pull/3971)]
 * replication: When bootstrapping a new secondary, if the initial cluster
   connection fails, Vault will attempt to roll back state so that
   bootstrapping can be tried again, rather than having to recreate the
   downstream cluster. This will still require fetching a new secondary
   activation token.

BUG FIXES:

 * auth/aws: Update libraries to fix regression verifying PKCS#7 identity
   documents [[GH-4014](https://github.com/hashicorp/vault/pull/4014)]
 * listener: Revert to Go 1.9 for now to allow certificates with non-DNS names
   in their DNS SANs to be used for Vault's TLS connections [[GH-4028](https://github.com/hashicorp/vault/pull/4028)]
 * replication: Fix issue with a performance secondary/DR primary node losing
   its DR primary status when performing an update-primary operation
 * replication: Fix issue where performance secondaries could be unable to
   automatically connect to a performance primary after that performance
   primary has been promoted to a DR primary from a DR secondary
 * ui: Fix behavior when a value contains a `.`

## 0.9.4 (February 20th, 2018)

SECURITY:

 * Role Tags used with the EC2 style of AWS auth were being improperly parsed;
   as a result they were not being used to properly restrict values.
   Implementations following our suggestion of using these as defense-in-depth
   rather than the only source of restriction should not have significant
   impact.

FEATURES:

 * **ChaCha20-Poly1305 support in `transit`**: You can now encrypt and decrypt
   with ChaCha20-Poly1305 in `transit`. Key derivation and convergent
   encryption is also supported.
 * **Okta Push support in Okta Auth Backend**: If a user account has MFA
   required within Okta, an Okta Push MFA flow can be used to successfully
   finish authentication.
 * **PKI Improvements**: Custom OID subject alternate names can now be set,
   subject to allow restrictions that support globbing. Additionally, Country,
   Locality, Province, Street Address, and Postal Code can now be set in
   certificate subjects.
 * **Manta Storage**: Joyent Triton Manta can now be used for Vault storage
 * **Google Cloud Spanner Storage**: Google Cloud Spanner can now be used for
   Vault storage

IMPROVEMENTS:

 * auth/centrify: Add CLI helper
 * audit: Always log failure metrics, even if zero, to ensure the values appear
   on dashboards [[GH-3937](https://github.com/hashicorp/vault/pull/3937)]
 * cli: Disable color when output is not a TTY [[GH-3897](https://github.com/hashicorp/vault/pull/3897)]
 * cli: Add `-format` flag to all subcommands [[GH-3897](https://github.com/hashicorp/vault/pull/3897)]
 * cli: Do not display deprecation warnings when the format is not table
   [[GH-3897](https://github.com/hashicorp/vault/pull/3897)]
 * core: If over a predefined lease count (256k), log a warning not more than
   once a minute. Too many leases can be problematic for many of the storage
   backends and often this number of leases is indicative of a need for
   workflow improvements. [[GH-3957](https://github.com/hashicorp/vault/pull/3957)]
 * secret/nomad: Have generated ACL tokens cap out at 64 characters [[GH-4009](https://github.com/hashicorp/vault/pull/4009)]
 * secret/pki: Country, Locality, Province, Street Address, and Postal Code can
   now be set on certificates [[GH-3992](https://github.com/hashicorp/vault/pull/3992)]
 * secret/pki: UTF-8 Other Names can now be set in Subject Alternate Names in
   issued certs; allowed values can be set per role and support globbing
   [[GH-3889](https://github.com/hashicorp/vault/pull/3889)]
 * secret/pki: Add a flag to make the common name optional on certs [[GH-3940](https://github.com/hashicorp/vault/pull/3940)]
 * secret/pki: Ensure only DNS-compatible names go into DNS SANs; additionally,
   properly handle IDNA transformations for these DNS names [[GH-3953](https://github.com/hashicorp/vault/pull/3953)]
 * secret/ssh: Add `valid-principles` flag to CLI for CA mode [[GH-3922](https://github.com/hashicorp/vault/pull/3922)]
 * storage/manta: Add Manta storage [[GH-3270](https://github.com/hashicorp/vault/pull/3270)]
 * ui (Enterprise): Support for ChaCha20-Poly1305 keys in the transit engine.

BUG FIXES:
 * api/renewer: Honor increment value in renew auth calls [[GH-3904](https://github.com/hashicorp/vault/pull/3904)]
 * auth/approle: Fix inability to use limited-use-count secret IDs on
   replication performance secondaries
 * auth/approle: Cleanup of secret ID accessors during tidy and removal of
   dangling accessor entries [[GH-3924](https://github.com/hashicorp/vault/pull/3924)]
 * auth/aws-ec2: Avoid masking of role tag response [[GH-3941](https://github.com/hashicorp/vault/pull/3941)]
 * auth/cert: Verify DNS SANs in the authenticating certificate [[GH-3982](https://github.com/hashicorp/vault/pull/3982)]
 * auth/okta: Return configured durations as seconds, not nanoseconds [[GH-3871](https://github.com/hashicorp/vault/pull/3871)]
 * auth/okta: Get all okta groups for a user vs. default 200 limit [[GH-4034](https://github.com/hashicorp/vault/pull/4034)]
 * auth/token: Token creation via the CLI no longer forces periodic token
   creation. Passing an explicit zero value for the period no longer create
   periodic tokens. [[GH-3880](https://github.com/hashicorp/vault/pull/3880)]
 * command: Fix interpreted formatting directives when printing raw fields
   [[GH-4005](https://github.com/hashicorp/vault/pull/4005)]
 * command: Correctly format output when using -field and -format flags at the
   same time [[GH-3987](https://github.com/hashicorp/vault/pull/3987)]
 * command/rekey: Re-add lost `stored-shares` parameter [[GH-3974](https://github.com/hashicorp/vault/pull/3974)]
 * command/ssh: Create and reuse the api client [[GH-3909](https://github.com/hashicorp/vault/pull/3909)]
 * command/status: Fix panic when status returns 500 from leadership lookup
   [[GH-3998](https://github.com/hashicorp/vault/pull/3998)]
 * identity: Fix race when creating entities [[GH-3932](https://github.com/hashicorp/vault/pull/3932)]
 * plugin/gRPC: Fixed an issue with list requests and raw responses coming from
   plugins using gRPC transport [[GH-3881](https://github.com/hashicorp/vault/pull/3881)]
 * plugin/gRPC: Fix panic when special paths are not set [[GH-3946](https://github.com/hashicorp/vault/pull/3946)]
 * secret/pki: Verify a name is a valid hostname before adding to DNS SANs
   [[GH-3918](https://github.com/hashicorp/vault/pull/3918)]
 * secret/transit: Fix auditing when reading a key after it has been backed up
   or restored [[GH-3919](https://github.com/hashicorp/vault/pull/3919)]
 * secret/transit: Fix storage/memory consistency when persistence fails
   [[GH-3959](https://github.com/hashicorp/vault/pull/3959)]
 * storage/consul: Validate that service names are RFC 1123 compliant [[GH-3960](https://github.com/hashicorp/vault/pull/3960)]
 * storage/etcd3: Fix memory ballooning with standby instances [[GH-3798](https://github.com/hashicorp/vault/pull/3798)]
 * storage/etcd3: Fix large lists (like token loading at startup) not being
   handled [[GH-3772](https://github.com/hashicorp/vault/pull/3772)]
 * storage/postgresql: Fix compatibility with versions using custom string
   version tags [[GH-3949](https://github.com/hashicorp/vault/pull/3949)]
 * storage/zookeeper: Update vendoring to fix freezing issues [[GH-3896](https://github.com/hashicorp/vault/pull/3896)]
 * ui (Enterprise): Decoding the replication token should no longer error and
   prevent enabling of a secondary replication cluster via the ui.
 * plugin/gRPC: Add connection info to the request object [[GH-3997](https://github.com/hashicorp/vault/pull/3997)]

## 0.9.3 (January 28th, 2018)

A regression from a feature merge disabled the Nomad secrets backend in 0.9.2.
This release re-enables the Nomad secrets backend; it is otherwise identical to
0.9.2.

## 0.9.2 (January 26th, 2018)

SECURITY:

 * Okta Auth Backend: While the Okta auth backend was successfully verifying
   usernames and passwords, it was not checking the returned state of the
   account, so accounts that had been marked locked out could still be used to
   log in. Only accounts in SUCCESS or PASSWORD_WARN states are now allowed.
 * Periodic Tokens: A regression in 0.9.1 meant that periodic tokens created by
   the AppRole, AWS, and Cert auth backends would expire when the max TTL for
   the backend/mount/system was hit instead of their stated behavior of living
   as long as they are renewed. This is now fixed; existing tokens do not have
   to be reissued as this was purely a regression in the renewal logic.
 * Seal Wrapping: During certain replication states values written marked for
   seal wrapping may not be wrapped on the secondaries. This has been fixed,
   and existing values will be wrapped on next read or write. This does not
   affect the barrier keys.

DEPRECATIONS/CHANGES:

 * `sys/health` DR Secondary Reporting: The `replication_dr_secondary` bool
   returned by `sys/health` could be misleading since it would be `false` both
   when a cluster was not a DR secondary but also when the node is a standby in
   the cluster and has not yet fully received state from the active node. This
   could cause health checks on LBs to decide that the node was acceptable for
   traffic even though DR secondaries cannot handle normal Vault traffic. (In
   other words, the bool could only convey "yes" or "no" but not "not sure
   yet".) This has been replaced by `replication_dr_mode` and
   `replication_perf_mode` which are string values that convey the current
   state of the node; a value of `disabled` indicates that replication is
   disabled or the state is still being discovered. As a result, an LB check
   can positively verify that the node is both not `disabled` and is not a DR
   secondary, and avoid sending traffic to it if either is true.
 * PKI Secret Backend Roles parameter types: For `ou` and `organization`
   in role definitions in the PKI secret backend, input can now be a
   comma-separated string or an array of strings. Reading a role will
   now return arrays for these parameters.
 * Plugin API Changes: The plugin API has been updated to utilize golang's
   context.Context package. Many function signatures now accept a context
   object as the first parameter. Existing plugins will need to pull in the
   latest Vault code and update their function signatures to begin using
   context and the new gRPC transport.

FEATURES:

 * **gRPC Backend Plugins**: Backend plugins now use gRPC for transport,
   allowing them to be written in other languages.
 * **Brand New CLI**: Vault has a brand new CLI interface that is significantly
   streamlined, supports autocomplete, and is almost entirely backwards
   compatible.
 * **UI: PKI Secret Backend (Enterprise)**: Configure PKI secret backends,
   create and browse roles and certificates, and issue and sign certificates via
   the listed roles.

IMPROVEMENTS:

 * auth/aws: Handle IAM headers produced by clients that formulate numbers as
   ints rather than strings [[GH-3763](https://github.com/hashicorp/vault/pull/3763)]
 * auth/okta: Support JSON lists when specifying groups and policies [[GH-3801](https://github.com/hashicorp/vault/pull/3801)]
 * autoseal/hsm: Attempt reconnecting to the HSM on certain kinds of issues,
   including HA scenarios for some Gemalto HSMs.
   (Enterprise)
 * cli: Output password prompts to stderr to make it easier to pipe an output
   token to another command [[GH-3782](https://github.com/hashicorp/vault/pull/3782)]
 * core: Report replication status in `sys/health` [[GH-3810](https://github.com/hashicorp/vault/pull/3810)]
 * physical/s3: Allow using paths with S3 for non-AWS deployments [[GH-3730](https://github.com/hashicorp/vault/pull/3730)]
 * physical/s3: Add ability to disable SSL for non-AWS deployments [[GH-3730](https://github.com/hashicorp/vault/pull/3730)]
 * plugins: Args for plugins can now be specified separately from the command,
   allowing the same output format and input format for plugin information
   [[GH-3778](https://github.com/hashicorp/vault/pull/3778)]
 * secret/pki: `ou` and `organization` can now be specified as a
   comma-separated string or an array of strings [[GH-3804](https://github.com/hashicorp/vault/pull/3804)]
 * plugins: Plugins will fall back to using netrpc as the communication protocol
   on older versions of Vault [[GH-3833](https://github.com/hashicorp/vault/pull/3833)]

BUG FIXES:

 * auth/(approle,aws,cert): Fix behavior where periodic tokens generated by
   these backends could not have their TTL renewed beyond the system/mount max
   TTL value [[GH-3803](https://github.com/hashicorp/vault/pull/3803)]
 * auth/aws: Fix error returned if `bound_iam_principal_arn` was given to an
   existing role update [[GH-3843](https://github.com/hashicorp/vault/pull/3843)]
 * core/sealwrap: Speed improvements and bug fixes (Enterprise)
 * identity: Delete group alias when an external group is deleted [[GH-3773](https://github.com/hashicorp/vault/pull/3773)]
 * legacymfa/duo: Fix intermittent panic when Duo could not be reached
   [[GH-2030](https://github.com/hashicorp/vault/pull/2030)]
 * secret/database: Fix a location where a lock could potentially not be
   released, leading to deadlock [[GH-3774](https://github.com/hashicorp/vault/pull/3774)]
 * secret/(all databases) Fix behavior where if a max TTL was specified but no
   default TTL was specified the system/mount default TTL would be used but not
   be capped by the local max TTL [[GH-3814](https://github.com/hashicorp/vault/pull/3814)]
 * secret/database: Fix an issue where plugins were not closed properly if they
   failed to initialize [[GH-3768](https://github.com/hashicorp/vault/pull/3768)]
 * ui: mounting a secret backend will now properly set `max_lease_ttl` and
   `default_lease_ttl` when specified - previously both fields set
   `default_lease_ttl`.

## 0.9.1 (December 21st, 2017)

DEPRECATIONS/CHANGES:

 * AppRole Case Sensitivity: In prior versions of Vault, `list` operations
   against AppRole roles would require preserving case in the role name, even
   though most other operations within AppRole are case-insensitive with
   respect to the role name. This has been fixed; existing roles will behave as
   they have in the past, but new roles will act case-insensitively in these
   cases.
 * Token Auth Backend Roles parameter types: For `allowed_policies` and
   `disallowed_policies` in role definitions in the token auth backend, input
   can now be a comma-separated string or an array of strings. Reading a role
   will now return arrays for these parameters.
 * Transit key exporting: You can now mark a key in the `transit` backend as
   `exportable` at any time, rather than just at creation time; however, once
   this value is set, it still cannot be unset.
 * PKI Secret Backend Roles parameter types: For `allowed_domains` and
   `key_usage` in role definitions in the PKI secret backend, input
   can now be a comma-separated string or an array of strings. Reading a role
   will now return arrays for these parameters.
 * SSH Dynamic Keys Method Defaults to 2048-bit Keys: When using the dynamic
   key method in the SSH backend, the default is now to use 2048-bit keys if no
   specific key bit size is specified.
 * Consul Secret Backend lease handling: The `consul` secret backend can now
   accept both strings and integer numbers of seconds for its lease value. The
   value returned on a role read will be an integer number of seconds instead
   of a human-friendly string.
 * Unprintable characters not allowed in API paths: Unprintable characters are
   no longer allowed in names in the API (paths and path parameters), with an
   extra restriction on whitespace characters. Allowed characters are those
   that are considered printable by Unicode plus spaces.

FEATURES:

 * **Transit Backup/Restore**: The `transit` backend now supports a backup
   operation that can export a given key, including all key versions and
   configuration, as well as a restore operation allowing import into another
   Vault.
 * **gRPC Database Plugins**: Database plugins now use gRPC for transport,
   allowing them to be written in other languages.
 * **Nomad Secret Backend**: Nomad ACL tokens can now be generated and revoked
   using Vault.
 * **TLS Cert Auth Backend Improvements**: The `cert` auth backend can now
   match against custom certificate extensions via exact or glob matching, and
   additionally supports max_ttl and periodic token toggles.

IMPROVEMENTS:

 * auth/cert: Support custom certificate constraints [[GH-3634](https://github.com/hashicorp/vault/pull/3634)]
 * auth/cert: Support setting `max_ttl` and `period` [[GH-3642](https://github.com/hashicorp/vault/pull/3642)]
 * audit/file: Setting a file mode of `0000` will now disable Vault from
   automatically `chmod`ing the log file [[GH-3649](https://github.com/hashicorp/vault/pull/3649)]
 * auth/github: The legacy MFA system can now be used with the GitHub auth
   backend [[GH-3696](https://github.com/hashicorp/vault/pull/3696)]
 * auth/okta: The legacy MFA system can now be used with the Okta auth backend
   [[GH-3653](https://github.com/hashicorp/vault/pull/3653)]
 * auth/token: `allowed_policies` and `disallowed_policies` can now be specified
   as a comma-separated string or an array of strings [[GH-3641](https://github.com/hashicorp/vault/pull/3641)]
 * command/server: The log level can now be specified with `VAULT_LOG_LEVEL`
   [[GH-3721](https://github.com/hashicorp/vault/pull/3721)]
 * core: Period values from auth backends will now be checked and applied to the
   TTL value directly by core on login and renewal requests [[GH-3677](https://github.com/hashicorp/vault/pull/3677)]
 * database/mongodb: Add optional `write_concern` parameter, which can be set
   during database configuration. This establishes a session-wide [write
   concern](https://docs.mongodb.com/manual/reference/write-concern/) for the
   lifecycle of the mount [[GH-3646](https://github.com/hashicorp/vault/pull/3646)]
 * http: Request path containing non-printable characters will return 400 - Bad
   Request [[GH-3697](https://github.com/hashicorp/vault/pull/3697)]
 * mfa/okta: Filter a given email address as a login filter, allowing operation
   when login email and account email are different
 * plugins: Make Vault more resilient when unsealing when plugins are
   unavailable [[GH-3686](https://github.com/hashicorp/vault/pull/3686)]
 * secret/pki: `allowed_domains` and `key_usage` can now be specified
   as a comma-separated string or an array of strings [[GH-3642](https://github.com/hashicorp/vault/pull/3642)]
 * secret/ssh: Allow 4096-bit keys to be used in dynamic key method [[GH-3593](https://github.com/hashicorp/vault/pull/3593)]
 * secret/consul: The Consul secret backend now uses the value of `lease` set
   on the role, if set, when renewing a secret. [[GH-3796](https://github.com/hashicorp/vault/pull/3796)]
 * storage/mysql: Don't attempt database creation if it exists, which can help
   under certain permissions constraints [[GH-3716](https://github.com/hashicorp/vault/pull/3716)]

BUG FIXES:

 * api/status (enterprise): Fix status reporting when using an auto seal
 * auth/approle: Fix case-sensitive/insensitive comparison issue [[GH-3665](https://github.com/hashicorp/vault/pull/3665)]
 * auth/cert: Return `allowed_names` on role read [[GH-3654](https://github.com/hashicorp/vault/pull/3654)]
 * auth/ldap: Fix incorrect control information being sent [[GH-3402](https://github.com/hashicorp/vault/pull/3402)] [[GH-3496](https://github.com/hashicorp/vault/pull/3496)]
   [[GH-3625](https://github.com/hashicorp/vault/pull/3625)] [[GH-3656](https://github.com/hashicorp/vault/pull/3656)]
 * core: Fix seal status reporting when using an autoseal
 * core: Add creation path to wrap info for a control group token
 * core: Fix potential panic that could occur using plugins when a node
   transitioned from active to standby [[GH-3638](https://github.com/hashicorp/vault/pull/3638)]
 * core: Fix memory ballooning when a connection would connect to the cluster
   port and then go away -- redux! [[GH-3680](https://github.com/hashicorp/vault/pull/3680)]
 * core: Replace recursive token revocation logic with depth-first logic, which
   can avoid hitting stack depth limits in extreme cases [[GH-2348](https://github.com/hashicorp/vault/pull/2348)]
 * core: When doing a read on configured audited-headers, properly handle case
   insensitivity [[GH-3701](https://github.com/hashicorp/vault/pull/3701)]
 * core/pkcs11 (enterprise): Fix panic when PKCS#11 library is not readable
 * database/mysql: Allow the creation statement to use commands that are not yet
   supported by the prepare statement protocol [[GH-3619](https://github.com/hashicorp/vault/pull/3619)]
 * plugin/auth-gcp: Fix IAM roles when using `allow_gce_inference` [VPAG-19]

## 0.9.0.1 (November 21st, 2017) (Enterprise Only)

IMPROVEMENTS:

 * auth/gcp: Support seal wrapping of configuration parameters
 * auth/kubernetes: Support seal wrapping of configuration parameters

BUG FIXES:

 * Fix an upgrade issue with some physical backends when migrating from legacy
   HSM stored key support to the new Seal Wrap mechanism (Enterprise)
 * mfa: Add the 'mfa' flag that was removed by mistake [[GH-4223](https://github.com/hashicorp/vault/pull/4223)]

## 0.9.0 (November 14th, 2017)

DEPRECATIONS/CHANGES:

 * HSM config parameter requirements: When using Vault with an HSM, a new
   parameter is required: `hmac_key_label`.  This performs a similar function to
   `key_label` but for the HMAC key Vault will use. Vault will generate a
   suitable key if this value is specified and `generate_key` is set true.
 * API HTTP client behavior: When calling `NewClient` the API no longer
   modifies the provided client/transport. In particular this means it will no
   longer enable redirection limiting and HTTP/2 support on custom clients. It
   is suggested that if you want to make changes to an HTTP client that you use
   one created by `DefaultConfig` as a starting point.
 * AWS EC2 client nonce behavior: The client nonce generated by the backend
   that gets returned along with the authentication response will be audited in
   plaintext. If this is undesired, the clients can choose to supply a custom
   nonce to the login endpoint. The custom nonce set by the client will from
   now on, not be returned back with the authentication response, and hence not
   audit logged.
 * AWS Auth role options: The API will now error when trying to create or
   update a role with the mutually-exclusive options
   `disallow_reauthentication` and `allow_instance_migration`.
 * SSH CA role read changes: When reading back a role from the `ssh` backend,
   the TTL/max TTL values will now be an integer number of seconds rather than
   a string. This better matches the API elsewhere in Vault.
 * SSH role list changes: When listing roles from the `ssh` backend via the API,
   the response data will additionally return a `key_info` map that will contain
   a map of each key with a corresponding object containing the `key_type`.
 * More granularity in audit logs: Audit request and response entries are still
   in RFC3339 format but now have a granularity of nanoseconds.
 * High availability related values have been moved out of the `storage` and
   `ha_storage` stanzas, and into the top-level configuration. `redirect_addr`
   has been renamed to `api_addr`. The stanzas still support accepting
   HA-related values to maintain backward compatibility, but top-level values
   will take precedence.
 * A new `seal` stanza has been added to the configuration file, which is
   optional and enables configuration of the seal type to use for additional
   data protection, such as using HSM or Cloud KMS solutions to encrypt and
   decrypt data.

FEATURES:

 * **RSA Support for Transit Backend**: Transit backend can now generate RSA
   keys which can be used for encryption and signing. [[GH-3489](https://github.com/hashicorp/vault/pull/3489)]
 * **Identity System**: Now in open source and with significant enhancements,
   Identity is an integrated system for understanding users across tokens and
   enabling easier management of users directly and via groups.
 * **External Groups in Identity**: Vault can now automatically assign users
   and systems to groups in Identity based on their membership in external
   groups.
 * **Seal Wrap / FIPS 140-2 Compatibility (Enterprise)**: Vault can now take
   advantage of FIPS 140-2-certified HSMs to ensure that Critical Security
   Parameters are protected in a compliant fashion. Vault's implementation has
   received a statement of compliance from Leidos.
 * **Control Groups (Enterprise)**: Require multiple members of an Identity
   group to authorize a requested action before it is allowed to run.
 * **Cloud Auto-Unseal (Enterprise)**: Automatically unseal Vault using AWS KMS
   and GCP CKMS.
 * **Sentinel Integration (Enterprise)**: Take advantage of HashiCorp Sentinel
   to create extremely flexible access control policies -- even on
   unauthenticated endpoints.
 * **Barrier Rekey Support for Auto-Unseal (Enterprise)**: When using auto-unsealing
   functionality, the `rekey` operation is now supported; it uses recovery keys
   to authorize the master key rekey.
 * **Operation Token for Disaster Recovery Actions (Enterprise)**: When using
   Disaster Recovery replication, a token can be created that can be used to
   authorize actions such as promotion and updating primary information, rather
   than using recovery keys.
 * **Trigger Auto-Unseal with Recovery Keys (Enterprise)**: When using
   auto-unsealing, a request to unseal Vault can be triggered by a threshold of
   recovery keys, rather than requiring the Vault process to be restarted.
 * **UI Redesign (Enterprise)**: All new experience for the Vault Enterprise
   UI. The look and feel has been completely redesigned to give users a better
   experience and make managing secrets fast and easy.
 * **UI: SSH Secret Backend (Enterprise)**: Configure an SSH secret backend,
   create and browse roles. And use them to sign keys or generate one time
   passwords.
 * **UI: AWS Secret Backend (Enterprise)**: You can now configure the AWS
   backend via the Vault Enterprise UI. In addition you can create roles,
   browse the roles and Generate IAM Credentials from them in the UI.

IMPROVEMENTS:

 * api: Add ability to set custom headers on each call [[GH-3394](https://github.com/hashicorp/vault/pull/3394)]
 * command/server: Add config option to disable requesting client certificates
   [[GH-3373](https://github.com/hashicorp/vault/pull/3373)]
 * auth/aws: Max retries can now be customized for the AWS client [[GH-3965](https://github.com/hashicorp/vault/pull/3965)]
 * core: Disallow mounting underneath an existing path, not just over [[GH-2919](https://github.com/hashicorp/vault/pull/2919)]
 * physical/file: Use `700` as permissions when creating directories. The files
   themselves were `600` and are all encrypted, but this doesn't hurt.
 * secret/aws: Add ability to use custom IAM/STS endpoints [[GH-3416](https://github.com/hashicorp/vault/pull/3416)]
 * secret/aws: Max retries can now be customized for the AWS client [[GH-3965](https://github.com/hashicorp/vault/pull/3965)]
 * secret/cassandra: Work around Cassandra ignoring consistency levels for a
   user listing query [[GH-3469](https://github.com/hashicorp/vault/pull/3469)]
 * secret/pki: Private keys can now be marshalled as PKCS#8 [[GH-3518](https://github.com/hashicorp/vault/pull/3518)]
 * secret/pki: Allow entering URLs for `pki` as both comma-separated strings and JSON
   arrays [[GH-3409](https://github.com/hashicorp/vault/pull/3409)]
 * secret/ssh: Role TTL/max TTL can now be specified as either a string or an
   integer [[GH-3507](https://github.com/hashicorp/vault/pull/3507)]
 * secret/transit: Sign and verify operations now support a `none` hash
   algorithm to allow signing/verifying pre-hashed data [[GH-3448](https://github.com/hashicorp/vault/pull/3448)]
 * secret/database: Add the ability to glob allowed roles in the Database Backend [[GH-3387](https://github.com/hashicorp/vault/pull/3387)]
 * ui (enterprise): Support for RSA keys in the transit backend
 * ui (enterprise): Support for DR Operation Token generation, promoting, and
   updating primary on DR Secondary clusters

BUG FIXES:

 * api: Fix panic when setting a custom HTTP client but with a nil transport
   [[GH-3435](https://github.com/hashicorp/vault/pull/3435)] [[GH-3437](https://github.com/hashicorp/vault/pull/3437)]
 * api: Fix authing to the `cert` backend when the CA for the client cert is
   not known to the server's listener [[GH-2946](https://github.com/hashicorp/vault/pull/2946)]
 * auth/approle: Create role ID index during read if a role is missing one [[GH-3561](https://github.com/hashicorp/vault/pull/3561)]
 * auth/aws: Don't allow mutually exclusive options [[GH-3291](https://github.com/hashicorp/vault/pull/3291)]
 * auth/radius: Fix logging in in some situations [[GH-3461](https://github.com/hashicorp/vault/pull/3461)]
 * core: Fix memleak when a connection would connect to the cluster port and
   then go away [[GH-3513](https://github.com/hashicorp/vault/pull/3513)]
 * core: Fix panic if a single-use token is used to step-down or seal [[GH-3497](https://github.com/hashicorp/vault/pull/3497)]
 * core: Set rather than add headers to prevent some duplicated headers in
   responses when requests were forwarded to the active node [[GH-3485](https://github.com/hashicorp/vault/pull/3485)]
 * physical/etcd3: Fix some listing issues due to how etcd3 does prefix
   matching [[GH-3406](https://github.com/hashicorp/vault/pull/3406)]
 * physical/etcd3: Fix case where standbys can lose their etcd client lease
   [[GH-3031](https://github.com/hashicorp/vault/pull/3031)]
 * physical/file: Fix listing when underscores are the first component of a
   path [[GH-3476](https://github.com/hashicorp/vault/pull/3476)]
 * plugins: Allow response errors to be returned from backend plugins [[GH-3412](https://github.com/hashicorp/vault/pull/3412)]
 * secret/transit: Fix panic if the length of the input ciphertext was less
   than the expected nonce length [[GH-3521](https://github.com/hashicorp/vault/pull/3521)]
 * ui (enterprise): Reinstate support for generic secret backends - this was
   erroneously removed in a previous release

## 0.8.3 (September 19th, 2017)

CHANGES:

 * Policy input/output standardization: For all built-in authentication
   backends, policies can now be specified as a comma-delimited string or an
   array if using JSON as API input; on read, policies will be returned as an
   array; and the `default` policy will not be forcefully added to policies
   saved in configurations. Please note that the `default` policy will continue
   to be added to generated tokens, however, rather than backends adding
   `default` to the given set of input policies (in some cases, and not in
   others), the stored set will reflect the user-specified set.
 * `sign-self-issued` modifies Issuer in generated certificates: In 0.8.2 the
   endpoint would not modify the Issuer in the generated certificate, leaving
   the output self-issued. Although theoretically valid, in practice crypto
   stacks were unhappy validating paths containing such certs. As a result,
   `sign-self-issued` now encodes the signing CA's Subject DN into the Issuer
   DN of the generated certificate.
 * `sys/raw` requires enabling: While the `sys/raw` endpoint can be extremely
   useful in break-glass or support scenarios, it is also extremely dangerous.
   As of now, a configuration file option `raw_storage_endpoint` must be set in
   order to enable this API endpoint. Once set, the available functionality has
   been enhanced slightly; it now supports listing and decrypting most of
   Vault's core data structures, except for the encryption keyring itself.
 * `generic` is now `kv`: To better reflect its actual use, the `generic`
   backend is now `kv`. Using `generic` will still work for backwards
   compatibility.

FEATURES:

 * **GCE Support for GCP Auth**: GCE instances can now authenticate to Vault
   using machine credentials.
 * **Support for Kubernetes Service Account Auth**: Kubernetes Service Accounts
   can now authenticate to vault using JWT tokens.

IMPROVEMENTS:

 * configuration: Provide a config option to store Vault server's process ID
   (PID) in a file [[GH-3321](https://github.com/hashicorp/vault/pull/3321)]
 * mfa (Enterprise): Add the ability to use identity metadata in username format
 * mfa/okta (Enterprise): Add support for configuring base_url for API calls
 * secret/pki: `sign-intermediate` will now allow specifying a `ttl` value
   longer than the signing CA certificate's NotAfter value. [[GH-3325](https://github.com/hashicorp/vault/pull/3325)]
 * sys/raw: Raw storage access is now disabled by default [[GH-3329](https://github.com/hashicorp/vault/pull/3329)]

BUG FIXES:

 * auth/okta: Fix regression that removed the ability to set base_url [[GH-3313](https://github.com/hashicorp/vault/pull/3313)]
 * core: Fix panic while loading leases at startup on ARM processors
   [[GH-3314](https://github.com/hashicorp/vault/pull/3314)]
 * secret/pki: Fix `sign-self-issued` encoding the wrong subject public key
   [[GH-3325](https://github.com/hashicorp/vault/pull/3325)]

## 0.8.2.1 (September 11th, 2017) (Enterprise Only)

BUG FIXES:

 * Fix an issue upgrading to 0.8.2 for Enterprise customers.

## 0.8.2 (September 5th, 2017)

SECURITY:

* In prior versions of Vault, if authenticating via AWS IAM and requesting a
  periodic token, the period was not properly respected. This could lead to
  tokens expiring unexpectedly, or a token lifetime being longer than expected.
  Upon token renewal with Vault 0.8.2 the period will be properly enforced.

DEPRECATIONS/CHANGES:

* `vault ssh` users should supply `-mode` and `-role` to reduce the number of
  API calls. A future version of Vault will mark these optional values are
  required. Failure to supply `-mode` or `-role` will result in a warning.
* Vault plugins will first briefly run a restricted version of the plugin to
  fetch metadata, and then lazy-load the plugin on first request to prevent
  crash/deadlock of Vault during the unseal process. Plugins will need to be
  built with the latest changes in order for them to run properly.

FEATURES:

* **Lazy Lease Loading**: On startup, Vault will now load leases from storage
  in a lazy fashion (token checks and revocation/renewal requests still force
  an immediate load). For larger installations this can significantly reduce
  downtime when switching active nodes or bringing Vault up from cold start.
* **SSH CA Login with `vault ssh`**: `vault ssh` now supports the SSH CA
  backend for authenticating to machines. It also supports remote host key
  verification through the SSH CA backend, if enabled.
* **Signing of Self-Issued Certs in PKI**: The `pki` backend now supports
  signing self-issued CA certs. This is useful when switching root CAs.

IMPROVEMENTS:

 * audit/file: Allow specifying `stdout` as the `file_path` to log to standard
   output [[GH-3235](https://github.com/hashicorp/vault/pull/3235)]
 * auth/aws: Allow wildcards in `bound_iam_principal_arn` [[GH-3213](https://github.com/hashicorp/vault/pull/3213)]
 * auth/okta: Compare groups case-insensitively since Okta is only
   case-preserving [[GH-3240](https://github.com/hashicorp/vault/pull/3240)]
 * auth/okta: Standardize Okta configuration APIs across backends [[GH-3245](https://github.com/hashicorp/vault/pull/3245)]
 * cli: Add subcommand autocompletion that can be enabled with
   `vault -autocomplete-install` [[GH-3223](https://github.com/hashicorp/vault/pull/3223)]
 * cli: Add ability to handle wrapped responses when using `vault auth`. What
   is output depends on the other given flags; see the help output for that
   command for more information. [[GH-3263](https://github.com/hashicorp/vault/pull/3263)]
 * core: TLS cipher suites used for cluster behavior can now be set via
   `cluster_cipher_suites` in configuration [[GH-3228](https://github.com/hashicorp/vault/pull/3228)]
 * core: The `plugin_name` can now either be specified directly as part of the
   parameter or within the `config` object when mounting a secret or auth backend
   via `sys/mounts/:path` or `sys/auth/:path` respectively [[GH-3202](https://github.com/hashicorp/vault/pull/3202)]
 * core: It is now possible to update the `description` of a mount when
   mount-tuning, although this must be done through the HTTP layer [[GH-3285](https://github.com/hashicorp/vault/pull/3285)]
 * secret/databases/mongo: If an EOF is encountered, attempt reconnecting and
   retrying the operation [[GH-3269](https://github.com/hashicorp/vault/pull/3269)]
 * secret/pki: TTLs can now be specified as a string or an integer number of
   seconds [[GH-3270](https://github.com/hashicorp/vault/pull/3270)]
 * secret/pki: Self-issued certs can now be signed via
   `pki/root/sign-self-issued` [[GH-3274](https://github.com/hashicorp/vault/pull/3274)]
 * storage/gcp: Use application default credentials if they exist [[GH-3248](https://github.com/hashicorp/vault/pull/3248)]

BUG FIXES:

 * auth/aws: Properly use role-set period values for IAM-derived token renewals
   [[GH-3220](https://github.com/hashicorp/vault/pull/3220)]
 * auth/okta: Fix updating organization/ttl/max_ttl after initial setting
   [[GH-3236](https://github.com/hashicorp/vault/pull/3236)]
 * core: Fix PROXY when underlying connection is TLS [[GH-3195](https://github.com/hashicorp/vault/pull/3195)]
 * core: Policy-related commands would sometimes fail to act case-insensitively
   [[GH-3210](https://github.com/hashicorp/vault/pull/3210)]
 * storage/consul: Fix parsing TLS configuration when using a bare IPv6 address
   [[GH-3268](https://github.com/hashicorp/vault/pull/3268)]
 * plugins: Lazy-load plugins to prevent crash/deadlock during unseal process.
   [[GH-3255](https://github.com/hashicorp/vault/pull/3255)]
 * plugins: Skip mounting plugin-based secret and credential mounts when setting
   up mounts if the plugin is no longer present in the catalog. [[GH-3255](https://github.com/hashicorp/vault/pull/3255)]

## 0.8.1 (August 16th, 2017)

DEPRECATIONS/CHANGES:

 * PKI Root Generation: Calling `pki/root/generate` when a CA cert/key already
   exists will now return a `204` instead of overwriting an existing root. If
   you want to recreate the root, first run a delete operation on `pki/root`
   (requires `sudo` capability), then generate it again.

FEATURES:

 * **Oracle Secret Backend**: There is now an external plugin to support leased
   credentials for Oracle databases (distributed separately).
 * **GCP IAM Auth Backend**: There is now an authentication backend that allows
   using GCP IAM credentials to retrieve Vault tokens. This is available as
   both a plugin and built-in to Vault.
 * **PingID Push Support for Path-Based MFA (Enterprise)**: PingID Push can
   now be used for MFA with the new path-based MFA introduced in Vault
   Enterprise 0.8.
 * **Permitted DNS Domains Support in PKI**: The `pki` backend now supports
   specifying permitted DNS domains for CA certificates, allowing you to
   narrowly scope the set of domains for which a CA can issue or sign child
   certificates.
 * **Plugin Backend Reload Endpoint**: Plugin backends can now be triggered to
   reload using the `sys/plugins/reload/backend` endpoint and providing either
   the plugin name or the mounts to reload.
 * **Self-Reloading Plugins**: The plugin system will now attempt to reload a
   crashed or stopped plugin, once per request.

IMPROVEMENTS:

 * auth/approle: Allow array input for policies in addition to comma-delimited
   strings [[GH-3163](https://github.com/hashicorp/vault/pull/3163)]
 * plugins: Send logs through Vault's logger rather than stdout [[GH-3142](https://github.com/hashicorp/vault/pull/3142)]
 * secret/pki: Add `pki/root` delete operation [[GH-3165](https://github.com/hashicorp/vault/pull/3165)]
 * secret/pki: Don't overwrite an existing root cert/key when calling generate
   [[GH-3165](https://github.com/hashicorp/vault/pull/3165)]

BUG FIXES:

 * aws: Don't prefer a nil HTTP client over an existing one [[GH-3159](https://github.com/hashicorp/vault/pull/3159)]
 * core: If there is an error when checking for create/update existence, return
   500 instead of 400 [[GH-3162](https://github.com/hashicorp/vault/pull/3162)]
 * secret/database: Avoid creating usernames that are too long for legacy MySQL
   [[GH-3138](https://github.com/hashicorp/vault/pull/3138)]

## 0.8.0 (August 9th, 2017)

SECURITY:

 * We've added a note to the docs about the way the GitHub auth backend works
   as it may not be readily apparent that GitHub personal access tokens, which
   are used by the backend, can be used for unauthorized access if they are
   stolen from third party services and access to Vault is public.

DEPRECATIONS/CHANGES:

 * Database Plugin Backends: Passwords generated for these backends now
   enforce stricter password requirements, as opposed to the previous behavior
   of returning a randomized UUID. Passwords are of length 20, and have a `A1a-`
   characters prepended to ensure stricter requirements. No regressions are
   expected from this change. (For database backends that were previously
   substituting underscores for hyphens in passwords, this will remain the
   case.)
 * Lease Endpoints: The endpoints `sys/renew`, `sys/revoke`, `sys/revoke-prefix`,
   `sys/revoke-force` have been deprecated and relocated under `sys/leases`.
   Additionally, the deprecated path `sys/revoke-force` now requires the `sudo`
   capability.
 * Response Wrapping Lookup Unauthenticated: The `sys/wrapping/lookup` endpoint
   is now unauthenticated. This allows introspection of the wrapping info by
   clients that only have the wrapping token without then invalidating the
   token. Validation functions/checks are still performed on the token.

FEATURES:

 * **Cassandra Storage**: Cassandra can now be used for Vault storage
 * **CockroachDB Storage**: CockroachDB can now be used for Vault storage
 * **CouchDB Storage**: CouchDB can now be used for Vault storage
 * **SAP HANA Database Plugin**: The `databases` backend can now manage users
   for SAP HANA databases
 * **Plugin Backends**: Vault now supports running secret and auth backends as
   plugins. Plugins can be mounted like normal backends and can be developed
   independently from Vault.
 * **PROXY Protocol Support** Vault listeners can now be configured to honor
   PROXY protocol v1 information to allow passing real client IPs into Vault. A
   list of authorized addresses (IPs or subnets) can be defined and
   accept/reject behavior controlled.
 * **Lease Lookup and Browsing in the Vault Enterprise UI**: Vault Enterprise UI
   now supports lookup and listing of leases and the associated actions from the
   `sys/leases` endpoints in the API. These are located in the new top level
   navigation item "Leases".
 * **Filtered Mounts for Performance Mode Replication**: Whitelists or
   blacklists of mounts can be defined per-secondary to control which mounts
   are actually replicated to that secondary. This can allow targeted
   replication of specific sets of data to specific geolocations/datacenters.
 * **Disaster Recovery Mode Replication (Enterprise Only)**: There is a new
   replication mode, Disaster Recovery (DR), that performs full real-time
   replication (including tokens and leases) to DR secondaries. DR secondaries
   cannot handle client requests, but can be promoted to primary as needed for
   failover.
 * **Manage New Replication Features in the Vault Enterprise UI**: Support for
   Replication features in Vault Enterprise UI has expanded to include new DR
   Replication mode and management of Filtered Mounts in Performance Replication
   mode.
 * **Vault Identity (Enterprise Only)**: Vault's new Identity system allows
   correlation of users across tokens. At present this is only used for MFA,
   but will be the foundation of many other features going forward.
 * **Duo Push, Okta Push, and TOTP MFA For All Authenticated Paths (Enterprise
   Only)**: A brand new MFA system built on top of Identity allows MFA
   (currently Duo Push, Okta Push, and TOTP) for any authenticated path within
   Vault. MFA methods can be configured centrally, and TOTP keys live within
   the user's Identity information to allow using the same key across tokens.
   Specific MFA method(s) required for any given path within Vault can be
   specified in normal ACL path statements.

IMPROVEMENTS:

 * api: Add client method for a secret renewer background process [[GH-2886](https://github.com/hashicorp/vault/pull/2886)]
 * api: Add `RenewTokenAsSelf` [[GH-2886](https://github.com/hashicorp/vault/pull/2886)]
 * api: Client timeout can now be adjusted with the `VAULT_CLIENT_TIMEOUT` env
   var or with a new API function [[GH-2956](https://github.com/hashicorp/vault/pull/2956)]
 * api/cli: Client will now attempt to look up SRV records for the given Vault
   hostname [[GH-3035](https://github.com/hashicorp/vault/pull/3035)]
 * audit/socket: Enhance reconnection logic and don't require the connection to
   be established at unseal time [[GH-2934](https://github.com/hashicorp/vault/pull/2934)]
 * audit/file: Opportunistically try re-opening the file on error [[GH-2999](https://github.com/hashicorp/vault/pull/2999)]
 * auth/approle: Add role name to token metadata [[GH-2985](https://github.com/hashicorp/vault/pull/2985)]
 * auth/okta: Allow specifying `ttl`/`max_ttl` inside the mount [[GH-2915](https://github.com/hashicorp/vault/pull/2915)]
 * cli: Client timeout can now be adjusted with the `VAULT_CLIENT_TIMEOUT` env
   var [[GH-2956](https://github.com/hashicorp/vault/pull/2956)]
 * command/auth: Add `-token-only` flag to `vault auth` that returns only the
   token on stdout and does not store it via the token helper [[GH-2855](https://github.com/hashicorp/vault/pull/2855)]
 * core: CORS allowed origins can now be configured [[GH-2021](https://github.com/hashicorp/vault/pull/2021)]
 * core: Add metrics counters for audit log failures [[GH-2863](https://github.com/hashicorp/vault/pull/2863)]
 * cors: Allow setting allowed headers via the API instead of always using
   wildcard [[GH-3023](https://github.com/hashicorp/vault/pull/3023)]
 * secret/ssh: Allow specifying the key ID format using template values for CA
   type [[GH-2888](https://github.com/hashicorp/vault/pull/2888)]
 * server: Add `tls_client_ca_file` option for specifying a CA file to use for
   client certificate verification when `tls_require_and_verify_client_cert` is
   enabled [[GH-3034](https://github.com/hashicorp/vault/pull/3034)]
 * storage/cockroachdb: Add CockroachDB storage backend [[GH-2713](https://github.com/hashicorp/vault/pull/2713)]
 * storage/couchdb: Add CouchDB storage backend [[GH-2880](https://github.com/hashicorp/vault/pull/2880)]
 * storage/mssql: Add `max_parallel` [[GH-3026](https://github.com/hashicorp/vault/pull/3026)]
 * storage/postgresql: Add `max_parallel` [[GH-3026](https://github.com/hashicorp/vault/pull/3026)]
 * storage/postgresql: Improve listing speed [[GH-2945](https://github.com/hashicorp/vault/pull/2945)]
 * storage/s3: More efficient paging when an object has a lot of subobjects
   [[GH-2780](https://github.com/hashicorp/vault/pull/2780)]
 * sys/wrapping: Make `sys/wrapping/lookup` unauthenticated [[GH-3084](https://github.com/hashicorp/vault/pull/3084)]
 * sys/wrapping: Wrapped tokens now store the original request path of the data
   [[GH-3100](https://github.com/hashicorp/vault/pull/3100)]
 * telemetry: Add support for DogStatsD [[GH-2490](https://github.com/hashicorp/vault/pull/2490)]

BUG FIXES:

 * api/health: Don't treat standby `429` codes as an error [[GH-2850](https://github.com/hashicorp/vault/pull/2850)]
 * api/leases: Fix lease lookup returning lease properties at the top level
 * audit: Fix panic when audit logging a read operation on an asymmetric
   `transit` key [[GH-2958](https://github.com/hashicorp/vault/pull/2958)]
 * auth/approle: Fix panic when secret and cidr list not provided in role
   [[GH-3075](https://github.com/hashicorp/vault/pull/3075)]
 * auth/aws: Look up proper account ID on token renew [[GH-3012](https://github.com/hashicorp/vault/pull/3012)]
 * auth/aws: Store IAM header in all cases when it changes [[GH-3004](https://github.com/hashicorp/vault/pull/3004)]
 * auth/ldap: Verify given certificate is PEM encoded instead of failing
   silently [[GH-3016](https://github.com/hashicorp/vault/pull/3016)]
 * auth/token: Don't allow using the same token ID twice when manually
   specifying [[GH-2916](https://github.com/hashicorp/vault/pull/2916)]
 * cli: Fix issue with parsing keys that start with special characters [[GH-2998](https://github.com/hashicorp/vault/pull/2998)]
 * core: Relocated `sys/leases/renew` returns same payload as original
   `sys/leases` endpoint [[GH-2891](https://github.com/hashicorp/vault/pull/2891)]
 * secret/ssh: Fix panic when signing with incorrect key type [[GH-3072](https://github.com/hashicorp/vault/pull/3072)]
 * secret/totp: Ensure codes can only be used once. This makes some automated
   workflows harder but complies with the RFC. [[GH-2908](https://github.com/hashicorp/vault/pull/2908)]
 * secret/transit: Fix locking when creating a key with unsupported options
   [[GH-2974](https://github.com/hashicorp/vault/pull/2974)]

## 0.7.3 (June 7th, 2017)

SECURITY:

 * Cert auth backend now checks validity of individual certificates: In
   previous versions of Vault, validity (e.g. expiration) of individual leaf
   certificates added for authentication was not checked. This was done to make
   it easier for administrators to control lifecycles of individual
   certificates added to the backend, e.g. the authentication material being
   checked was access to that specific certificate's private key rather than
   all private keys signed by a CA. However, this behavior is often unexpected
   and as a result can lead to insecure deployments, so we are now validating
   these certificates as well.
 * App-ID path salting was skipped in 0.7.1/0.7.2: A regression in 0.7.1/0.7.2
   caused the HMACing of any App-ID information stored in paths (including
   actual app-IDs and user-IDs) to be unsalted and written as-is from the API.
   In 0.7.3 any such paths will be automatically changed to salted versions on
   access (e.g. login or read); however, if you created new app-IDs or user-IDs
   in 0.7.1/0.7.2, you may want to consider whether any users with access to
   Vault's underlying data store may have intercepted these values, and
   revoke/roll them.

DEPRECATIONS/CHANGES:

 * Step-Down is Forwarded: When a step-down is issued against a non-active node
   in an HA cluster, it will now forward the request to the active node.

FEATURES:

 * **ed25519 Signing/Verification in Transit with Key Derivation**: The
   `transit` backend now supports generating
   [ed25519](https://ed25519.cr.yp.to/) keys for signing and verification
   functionality. These keys support derivation, allowing you to modify the
   actual encryption key used by supplying a `context` value.
 * **Key Version Specification for Encryption in Transit**: You can now specify
   the version of a key you use to wish to generate a signature, ciphertext, or
   HMAC. This can be controlled by the `min_encryption_version` key
   configuration property.
 * **Replication Primary Discovery (Enterprise)**: Replication primaries will
   now advertise the addresses of their local HA cluster members to replication
   secondaries. This helps recovery if the primary active node goes down and
   neither service discovery nor load balancers are in use to steer clients.

IMPROVEMENTS:

 * api/health: Add Sys().Health() [[GH-2805](https://github.com/hashicorp/vault/pull/2805)]
 * audit: Add auth information to requests that error out [[GH-2754](https://github.com/hashicorp/vault/pull/2754)]
 * command/auth: Add `-no-store` option that prevents the auth command from
   storing the returned token into the configured token helper [[GH-2809](https://github.com/hashicorp/vault/pull/2809)]
 * core/forwarding: Request forwarding now heartbeats to prevent unused
   connections from being terminated by firewalls or proxies
 * plugins/databases: Add MongoDB as an internal database plugin [[GH-2698](https://github.com/hashicorp/vault/pull/2698)]
 * storage/dynamodb: Add a method for checking the existence of children,
   speeding up deletion operations in the DynamoDB storage backend [[GH-2722](https://github.com/hashicorp/vault/pull/2722)]
 * storage/mysql: Add max_parallel parameter to MySQL backend [[GH-2760](https://github.com/hashicorp/vault/pull/2760)]
 * secret/databases: Support listing connections [[GH-2823](https://github.com/hashicorp/vault/pull/2823)]
 * secret/databases: Support custom renewal statements in Postgres database
   plugin [[GH-2788](https://github.com/hashicorp/vault/pull/2788)]
 * secret/databases: Use the role name as part of generated credentials
   [[GH-2812](https://github.com/hashicorp/vault/pull/2812)]
 * ui (Enterprise): Transit key and secret browsing UI handle large lists better
 * ui (Enterprise): root tokens are no longer persisted
 * ui (Enterprise): support for mounting Database and TOTP secret backends

BUG FIXES:

 * auth/app-id: Fix regression causing loading of salts to be skipped
 * auth/aws: Improve EC2 describe instances performance [[GH-2766](https://github.com/hashicorp/vault/pull/2766)]
 * auth/aws: Fix lookup of some instance profile ARNs [[GH-2802](https://github.com/hashicorp/vault/pull/2802)]
 * auth/aws: Resolve ARNs to internal AWS IDs which makes lookup at various
   points (e.g. renewal time) more robust [[GH-2814](https://github.com/hashicorp/vault/pull/2814)]
 * auth/aws: Properly honor configured period when using IAM authentication
   [[GH-2825](https://github.com/hashicorp/vault/pull/2825)]
 * auth/aws: Check that a bound IAM principal is not empty (in the current
   state of the role) before requiring it match the previously authenticated
   client [[GH-2781](https://github.com/hashicorp/vault/pull/2781)]
 * auth/cert: Fix panic on renewal [[GH-2749](https://github.com/hashicorp/vault/pull/2749)]
 * auth/cert: Certificate verification for non-CA certs [[GH-2761](https://github.com/hashicorp/vault/pull/2761)]
 * core/acl: Prevent race condition when compiling ACLs in some scenarios
   [[GH-2826](https://github.com/hashicorp/vault/pull/2826)]
 * secret/database: Increase wrapping token TTL; in a loaded scenario it could
   be too short
 * secret/generic: Allow integers to be set as the value of `ttl` field as the
   documentation claims is supported [[GH-2699](https://github.com/hashicorp/vault/pull/2699)]
 * secret/ssh: Added host key callback to ssh client config [[GH-2752](https://github.com/hashicorp/vault/pull/2752)]
 * storage/s3: Avoid a panic when some bad data is returned [[GH-2785](https://github.com/hashicorp/vault/pull/2785)]
 * storage/dynamodb: Fix list functions working improperly on Windows [[GH-2789](https://github.com/hashicorp/vault/pull/2789)]
 * storage/file: Don't leak file descriptors in some error cases
 * storage/swift: Fix pre-v3 project/tenant name reading [[GH-2803](https://github.com/hashicorp/vault/pull/2803)]

## 0.7.2 (May 8th, 2017)

BUG FIXES:

 * audit: Fix auditing entries containing certain kinds of time values
   [[GH-2689](https://github.com/hashicorp/vault/pull/2689)]

## 0.7.1 (May 5th, 2017)

DEPRECATIONS/CHANGES:

 * LDAP Auth Backend: Group membership queries will now run as the `binddn`
   user when `binddn`/`bindpass` are configured, rather than as the
   authenticating user as was the case previously.

FEATURES:

 * **AWS IAM Authentication**: IAM principals can get Vault tokens
   automatically, opening AWS-based authentication to users, ECS containers,
   Lambda instances, and more. Signed client identity information retrieved
   using the AWS API `sts:GetCallerIdentity` is validated against the AWS STS
   service before issuing a Vault token. This backend is unified with the
   `aws-ec2` authentication backend under the name `aws`, and allows additional
   EC2-related restrictions to be applied during the IAM authentication; the
   previous EC2 behavior is also still available. [[GH-2441](https://github.com/hashicorp/vault/pull/2441)]
 * **MSSQL Physical Backend**: You can now use Microsoft SQL Server as your
   Vault physical data store [[GH-2546](https://github.com/hashicorp/vault/pull/2546)]
 * **Lease Listing and Lookup**: You can now introspect a lease to get its
   creation and expiration properties via `sys/leases/lookup`; with `sudo`
   capability you can also list leases for lookup, renewal, or revocation via
   that endpoint. Various lease functions (renew, revoke, revoke-prefix,
   revoke-force) have also been relocated to `sys/leases/`, but they also work
   at the old paths for compatibility. Reading (but not listing) leases via
   `sys/leases/lookup` is now a part of the current `default` policy. [[GH-2650](https://github.com/hashicorp/vault/pull/2650)]
 * **TOTP Secret Backend**: You can now store multi-factor authentication keys
   in Vault and use the API to retrieve time-based one-time use passwords on
   demand. The backend can also be used to generate a new key and validate
   passwords generated by that key. [[GH-2492](https://github.com/hashicorp/vault/pull/2492)]
 * **Database Secret Backend & Secure Plugins (Beta)**: This new secret backend
   combines the functionality of the MySQL, PostgreSQL, MSSQL, and Cassandra
   backends. It also provides a plugin interface for extendability through
   custom databases. [[GH-2200](https://github.com/hashicorp/vault/pull/2200)]

IMPROVEMENTS:

 * auth/cert: Support for constraints on subject Common Name and DNS/email
   Subject Alternate Names in certificates [[GH-2595](https://github.com/hashicorp/vault/pull/2595)]
 * auth/ldap: Use the binding credentials to search group membership rather
   than the user credentials [[GH-2534](https://github.com/hashicorp/vault/pull/2534)]
 * cli/revoke: Add `-self` option to allow revoking the currently active token
   [[GH-2596](https://github.com/hashicorp/vault/pull/2596)]
 * core: Randomize x coordinate in Shamir shares [[GH-2621](https://github.com/hashicorp/vault/pull/2621)]
 * replication: Fix a bug when enabling `approle` on a primary before
   secondaries were connected
 * replication: Add heartbeating to ensure firewalls don't kill connections to
   primaries
 * secret/pki: Add `no_store` option that allows certificates to be issued
   without being stored. This removes the ability to look up and/or add to a
   CRL but helps with scaling to very large numbers of certificates. [[GH-2565](https://github.com/hashicorp/vault/pull/2565)]
 * secret/pki: If used with a role parameter, the `sign-verbatim/<role>`
   endpoint honors the values of `generate_lease`, `no_store`, `ttl` and
   `max_ttl` from the given role [[GH-2593](https://github.com/hashicorp/vault/pull/2593)]
 * secret/pki: Add role parameter `allow_glob_domains` that enables defining
   names in `allowed_domains` containing `*` glob patterns [[GH-2517](https://github.com/hashicorp/vault/pull/2517)]
 * secret/pki: Update certificate storage to not use characters that are not
   supported on some filesystems [[GH-2575](https://github.com/hashicorp/vault/pull/2575)]
 * storage/etcd3: Add `discovery_srv` option to query for SRV records to find
   servers [[GH-2521](https://github.com/hashicorp/vault/pull/2521)]
 * storage/s3: Support `max_parallel` option to limit concurrent outstanding
   requests [[GH-2466](https://github.com/hashicorp/vault/pull/2466)]
 * storage/s3: Use pooled transport for http client [[GH-2481](https://github.com/hashicorp/vault/pull/2481)]
 * storage/swift: Allow domain values for V3 authentication [[GH-2554](https://github.com/hashicorp/vault/pull/2554)]
 * tidy: Improvements to `auth/token/tidy` and `sys/leases/tidy` to handle more
   cleanup cases [[GH-2452](https://github.com/hashicorp/vault/pull/2452)]

BUG FIXES:

 * api: Respect a configured path in Vault's address [[GH-2588](https://github.com/hashicorp/vault/pull/2588)]
 * auth/aws-ec2: New bounds added as criteria to allow role creation [[GH-2600](https://github.com/hashicorp/vault/pull/2600)]
 * auth/ldap: Don't lowercase groups attached to users [[GH-2613](https://github.com/hashicorp/vault/pull/2613)]
 * cli: Don't panic if `vault write` is used with the `force` flag but no path
   [[GH-2674](https://github.com/hashicorp/vault/pull/2674)]
 * core: Help operations should request forward since standbys may not have
   appropriate info [[GH-2677](https://github.com/hashicorp/vault/pull/2677)]
 * replication: Fix enabling secondaries when certain mounts already existed on
   the primary
 * secret/mssql: Update mssql driver to support queries with colons [[GH-2610](https://github.com/hashicorp/vault/pull/2610)]
 * secret/pki: Don't lowercase O/OU values in certs [[GH-2555](https://github.com/hashicorp/vault/pull/2555)]
 * secret/pki: Don't attempt to validate IP SANs if none are provided [[GH-2574](https://github.com/hashicorp/vault/pull/2574)]
 * secret/ssh: Don't automatically lowercase principles in issued SSH certs
   [[GH-2591](https://github.com/hashicorp/vault/pull/2591)]
 * storage/consul: Properly handle state events rather than timing out
   [[GH-2548](https://github.com/hashicorp/vault/pull/2548)]
 * storage/etcd3: Ensure locks are released if client is improperly shut down
   [[GH-2526](https://github.com/hashicorp/vault/pull/2526)]

## 0.7.0 (March 21th, 2017)

SECURITY:

 * Common name not being validated when `exclude_cn_from_sans` option used in
   `pki` backend: When using a role in the `pki` backend that specified the
   `exclude_cn_from_sans` option, the common name would not then be properly
   validated against the role's constraints. This has been fixed. We recommend
   any users of this feature to upgrade to 0.7 as soon as feasible.

DEPRECATIONS/CHANGES:

 * List Operations Always Use Trailing Slash: Any list operation, whether via
   the `GET` or `LIST` HTTP verb, will now internally canonicalize the path to
   have a trailing slash. This makes policy writing more predictable, as it
   means clients will no longer work or fail based on which client they're
   using or which HTTP verb they're using. However, it also means that policies
   allowing `list` capability must be carefully checked to ensure that they
   contain a trailing slash; some policies may need to be split into multiple
   stanzas to accommodate.
 * PKI Defaults to Unleased Certificates: When issuing certificates from the
   PKI backend, by default, no leases will be issued. If you want to manually
   revoke a certificate, its serial number can be used with the `pki/revoke`
   endpoint. Issuing leases is still possible by enabling the `generate_lease`
   toggle in PKI role entries (this will default to `true` for upgrades, to
   keep existing behavior), which will allow using lease IDs to revoke
   certificates. For installations issuing large numbers of certificates (tens
   to hundreds of thousands, or millions), this will significantly improve
   Vault startup time since leases associated with these certificates will not
   have to be loaded; however note that it also means that revocation of a
   token used to issue
