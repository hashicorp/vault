module.exports = [
  {
    source: '/home',
    destination: '/',
    permanent: true,
  },
  {
    source: '/trial',
    destination: 'https://www.hashicorp.com/products/vault/trial',
    permanent: true,
  },
  {
    source: '/intro',
    destination: '/intro/getting-started',
    permanent: false,
  },
  {
    source: '/docs/release-notes/1.10',
    destination: '/docs/release-notes/1.10.0',
    permanent: true,
  },
  {
    source: '/api/secret/generic',
    destination: '/api-docs/secret/kv',
    permanent: true,
  },
  {
    source: '/api/system/renew',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/api/system/revoke-force',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/api/system/revoke-prefix',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/api/system/revoke',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/docs/auth/aws-ec2',
    destination: '/docs/auth/aws',
    permanent: true,
  },
  {
    source: '/docs/auth/jwt_oidc_providers',
    destination: '/docs/auth/jwt/oidc-providers',
    permanent: true,
  },
  {
    source: '/docs/auth/jwt/oidc_providers',
    destination: '/docs/auth/jwt/oidc-providers',
    permanent: true,
  },
  {
    source: '/docs/commands/environment',
    destination: '/docs/commands/#environment-variables',
    permanent: true,
  },
  {
    source: '/docs/commands/help',
    destination: '/docs/commands/path-help',
    permanent: true,
  },
  {
    source: '/docs/commands/read-write',
    destination: '/docs/commands#reading-and-writing-data',
    permanent: true,
  },
  {
    source: '/docs/config',
    destination: '/docs/configuration',
    permanent: true,
  },
  {
    source: '/docs/configuration/storage/google-cloud',
    destination: '/docs/configuration/storage/google-cloud-storage',
    permanent: true,
  },
  {
    source: '/docs/configuration/storage/spanner',
    destination: '/docs/configuration/storage/google-cloud-spanner',
    permanent: true,
  },
  {
    source: '/docs/enterprise/auto-unseal',
    destination: '/docs/concepts/seal.html',
    permanent: true,
  },
  {
    source: '/docs/enterprise/license/faqs',
    destination: '/docs/enterprise/license/faq',
    permanent: true,
  },
  {
    source: '/docs/enterprise/hsm/configuration',
    destination: '/docs/configuration/seal/pkcs11',
    permanent: true,
  },
  {
    source: '/docs/enterprise/ui',
    destination: '/docs/configuration/ui',
    permanent: true,
  },
  {
    source: '/docs/enterprise/automated-raft-snapshots',
    destination: '/docs/enterprise/automated-integrated-storage-snapshots',
    permanent: true,
  },
  {
    source: '/docs/guides/generate-root',
    destination: '/guides/operations/generate-root',
    permanent: true,
  },
  { source: '/docs/guides', destination: '/guides', permanent: true },
  {
    source: '/docs/guides/production',
    destination: '/guides/operations/production',
    permanent: true,
  },
  {
    source: '/docs/guides/replication',
    destination: '/guides/operations/replication',
    permanent: true,
  },
  {
    source: '/docs/guides/upgrading',
    destination: '/docs/upgrading',
    permanent: true,
  },
  {
    source: '/docs/http/sys-audit',
    destination: '/api-docs/system/audit',
    permanent: true,
  },
  {
    source: '/docs/http/sys-auth',
    destination: '/api-docs/system/auth',
    permanent: true,
  },
  {
    source: '/docs/http/sys-health',
    destination: '/api-docs/system/health',
    permanent: true,
  },
  {
    source: '/docs/http/sys-init',
    destination: '/api-docs/system/init',
    permanent: true,
  },
  {
    source: '/docs/http/sys-key-status',
    destination: '/api-docs/system/key-status',
    permanent: true,
  },
  {
    source: '/docs/http/sys-leader',
    destination: '/api-docs/system/leader',
    permanent: true,
  },
  {
    source: '/docs/http/sys-mounts',
    destination: '/api-docs/system/mounts',
    permanent: true,
  },
  {
    source: '/docs/http/sys-policy',
    destination: '/api-docs/system/policy',
    permanent: true,
  },
  {
    source: '/docs/http/sys-raw',
    destination: '/api-docs/system/raw',
    permanent: true,
  },
  {
    source: '/docs/http/sys-rekey',
    destination: '/api-docs/system/rekey',
    permanent: true,
  },
  {
    source: '/docs/http/sys-remount',
    destination: '/api-docs/system/remount',
    permanent: true,
  },
  {
    source: '/docs/http/sys-renew',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/docs/http/sys-revoke-prefix',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/docs/http/sys-revoke',
    destination: '/api-docs/system/leases',
    permanent: true,
  },
  {
    source: '/docs/http/sys-rotate',
    destination: '/api-docs/system/rotate',
    permanent: true,
  },
  {
    source: '/docs/http/sys-seal-status',
    destination: '/api-docs/system/seal-status',
    permanent: true,
  },
  {
    source: '/docs/http/sys-seal',
    destination: '/api-docs/system/seal',
    permanent: true,
  },
  {
    source: '/docs/http/sys-unseal',
    destination: '/api-docs/system/unseal',
    permanent: true,
  },
  {
    source: '/docs/http/sys-version-history',
    destination: '/api-docs/system/version-history',
    permanent: true,
  },
  {
    source: '/docs/install/install',
    destination: '/docs/install',
    permanent: true,
  },
  {
    source: '/docs/secrets/custom',
    destination: '/docs/plugins/plugin-management',
    permanent: true,
  },
  {
    source: '/docs/internals/plugins',
    destination: '/docs/plugins',
    permanent: true,
  },
  {
    source: '/docs/plugin-portal',
    destination: '/docs/plugins/plugin-portal',
    permanent: true,
  },
  {
    source: '/docs/secrets/generic',
    destination: '/docs/secrets/kv',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/hsm/behavior',
    destination: '/docs/enterprise/hsm/behavior',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/hsm/configuration',
    destination: '/docs/enterprise/hsm/configuration',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/hsm',
    destination: '/docs/enterprise/hsm',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/hsm/security',
    destination: '/docs/enterprise/hsm/security',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/identity',
    destination: '/docs/enterprise/identity',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise',
    destination: '/docs/enterprise',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/mfa',
    destination: '/docs/enterprise/mfa',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/mfa/mfa-duo',
    destination: '/docs/enterprise/mfa/mfa-duo',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/mfa/mfa-okta',
    destination: '/docs/enterprise/mfa/mfa-okta',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/mfa/mfa-pingid',
    destination: '/docs/enterprise/mfa/mfa-pingid',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/mfa/mfa-totp',
    destination: '/docs/enterprise/mfa/mfa-totp',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/replication',
    destination: '/docs/enterprise/replication',
    permanent: true,
  },
  {
    source: '/docs/vault-enterprise/ui',
    destination: '/docs/configuration/ui',
    permanent: true,
  },
  {
    source: '/docs/secrets/cassandra',
    destination: '/docs/secrets/databases/cassandra',
    permanent: true,
  },
  {
    source: '/docs/secrets/mongodb',
    destination: '/docs/secrets/databases/mongodb',
    permanent: true,
  },
  {
    source: '/docs/secrets/mssql',
    destination: '/docs/secrets/databases/mssql',
    permanent: true,
  },
  {
    source: '/docs/secrets/mysql',
    destination: '/docs/secrets/databases/mysql-maria',
    permanent: true,
  },
  {
    source: '/docs/secrets/postgresql',
    destination: '/docs/secrets/databases/postgresql',
    permanent: true,
  },
  {
    source: '/guides/authentication',
    destination: '/guides/identity/authentication',
    permanent: true,
  },
  {
    source: '/guides/configuration/authentication',
    destination: '/guides/identity/authentication',
    permanent: true,
  },
  {
    source: '/guides/configuration/generate-root',
    destination: '/guides/operations/generate-root',
    permanent: true,
  },
  {
    source: '/guides/configuration/lease',
    destination: '/guides/identity/lease',
    permanent: true,
  },
  {
    source: '/guides/configuration/plugin-backends',
    destination: '/guides/operations/plugin-backends',
    permanent: true,
  },
  {
    source: '/guides/configuration/policies',
    destination: '/guides/identity/policies',
    permanent: true,
  },
  {
    source: '/guides/cubbyhole',
    destination: '/guides/secret-mgmt/cubbyhole',
    permanent: true,
  },
  {
    source: '/guides/dynamic-secret',
    destination: '/guides/secret-mgmt/dynamic-secret',
    permanent: true,
  },
  {
    source: '/guides/generate-root',
    destination: '/guides/operations/generate-root',
    permanent: true,
  },
  {
    source: '/guides/lease',
    destination: '/guides/identity/lease',
    permanent: true,
  },
  {
    source: '/guides/plugin-backends',
    destination: '/guides/operations/plugin-backends',
    permanent: true,
  },
  {
    source: '/guides/policies',
    destination: '/guides/identity/policies',
    permanent: true,
  },
  {
    source: '/guides/production',
    destination: '/guides/operations/production',
    permanent: true,
  },
  {
    source: '/guides/rekeying-and-rotating',
    destination: '/guides/operations/rekeying-and-rotating',
    permanent: true,
  },
  {
    source: '/guides/replication',
    destination: '/guides/operations/replication',
    permanent: true,
  },
  {
    source: '/guides/static-secrets',
    destination: '/guides/secret-mgmt/static-secrets',
    permanent: true,
  },
  {
    source: '/intro/getting-started/acl',
    destination: '/intro/getting-started/policies',
    permanent: true,
  },
  {
    source: '/intro/getting-started/secret-backends',
    destination: '/intro/getting-started/secrets-engines',
    permanent: true,
  },
  {
    source: '/guides/configuration/rekeying-and-rotating',
    destination: '/guides/operations/rekeying-and-rotating',
    permanent: true,
  },
  {
    source: '/docs/guides/upgrading/:path*',
    destination: '/docs/upgrading/:path*',
    permanent: true,
  },
  {
    source: '/guides/upgrading/:path*',
    destination: '/docs/upgrading/:path*',
    permanent: true,
  },
  {
    source: '/docs/install/upgrade-:path*',
    destination: '/docs/upgrading/upgrade-:path*',
    permanent: true,
  },
  {
    source: '/docs/install/upgrade',
    destination: '/docs/upgrading',
    permanent: true,
  },
  {
    source: '/docs/platform/aws/lambda-extension-cache',
    destination: '/docs/platform/aws/lambda-extension',
    permanent: true,
  },
  // Guides and Intro redirects to Learn
  {
    source: '/guides',
    destination: 'https://learn.hashicorp.com/vault',
    permanent: true,
  },
  {
    source: '/guides/getting-started',
    destination: 'https://learn.hashicorp.com/vault',
    permanent: true,
  },
  {
    source: '/guides/operations',
    destination: 'https://learn.hashicorp.com/collections/vault/operations',
    permanent: true,
  },
  {
    source: '/guides/operations/reference-architecture',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/reference-architecture',
    permanent: true,
  },
  {
    source: '/guides/operations/deployment-guide',
    destination: 'https://learn.hashicorp.com/tutorials/vault/deployment-guide',
    permanent: true,
  },
  {
    source: '/guides/operations/vault-ha-consul',
    destination: 'https://learn.hashicorp.com/tutorials/vault/ha-with-consul',
    permanent: true,
  },
  {
    source: '/guides/operations/production',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/production-hardening',
    permanent: true,
  },
  {
    source: '/guides/operations/generate-root',
    destination: 'https://learn.hashicorp.com/tutorials/vault/generate-root',
    permanent: true,
  },
  {
    source: '/guides/operations/rekeying-and-rotating',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/rekeying-and-rotating',
    permanent: true,
  },
  {
    source: '/guides/operations/plugin-backends',
    destination: 'https://learn.hashicorp.com/tutorials/vault/plugin-backends',
    permanent: true,
  },
  {
    source: '/guides/operations/replication',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/performance-replication',
    permanent: true,
  },
  {
    source: '/guides/operations/disaster-recovery',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/disaster-recovery',
    permanent: true,
  },
  {
    source: '/guides/operations/mount-filter',
    destination: 'https://learn.hashicorp.com/tutorials/vault/paths-filter',
    permanent: true,
  },
  {
    source: '/guides/operations/performance-nodes',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/performance-standbys',
    permanent: true,
  },
  {
    source: '/guides/operations/multi-tenant',
    destination: 'https://learn.hashicorp.com/tutorials/vault/namespaces',
    permanent: true,
  },
  {
    source: '/guides/operations/autounseal-aws-kms',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/autounseal-aws-kms',
    permanent: true,
  },
  {
    source: '/guides/operations/seal-wrap',
    destination: 'https://learn.hashicorp.com/tutorials/vault/seal-wrap',
    permanent: true,
  },
  {
    source: '/guides/operations/monitoring',
    destination: 'https://learn.hashicorp.com/tutorials/vault/monitoring',
    permanent: true,
  },
  // Identity
  {
    source: '/guides/identity',
    destination: 'https://learn.hashicorp.com/collections/vault/operations',
    permanent: true,
  },
  {
    source: '/guides/identity/secure-intro',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/secure-introduction',
    permanent: true,
  },
  {
    source: '/guides/identity/policies',
    destination: 'https://learn.hashicorp.com/tutorials/vault/policies',
    permanent: true,
  },
  {
    source: '/guides/identity/policy-templating',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/policy-templating',
    permanent: true,
  },
  {
    source: '/guides/identity/authentication',
    destination: 'https://learn.hashicorp.com/tutorials/vault/approle',
    permanent: true,
  },
  {
    source: '/guides/identity/approle-trusted-entities',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/approle-trusted-entities',
    permanent: true,
  },
  {
    source: '/guides/identity/lease',
    destination: 'https://learn.hashicorp.com/tutorials/vault/tokens',
    permanent: true,
  },
  {
    source: '/guides/identity/identity',
    destination: 'https://learn.hashicorp.com/tutorials/vault/identity',
    permanent: true,
  },
  {
    source: '/guides/identity/sentinel',
    destination: 'https://learn.hashicorp.com/tutorials/vault/sentinel',
    permanent: true,
  },
  {
    source: '/guides/identity/control-groups',
    destination: 'https://learn.hashicorp.com/tutorials/vault/control-groups',
    permanent: true,
  },
  // Secrets management
  {
    source: '/guides/secret-mgmt/index.html',
    destination:
      'https://learn.hashicorp.com/collections/vault/secrets-management',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/static-secrets',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/sm-static-secrets',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/versioned-kv',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/sm-versioned-kv',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/dynamic-secrets',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/sm-dynamic-secrets',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/db-root-rotation',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/db-root-rotation',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/cubbyhole',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/sm-cubbyhole',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/ssh-otp',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/sm-ssh-otp',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/pki-engine',
    destination:
      'https://learn.hashicorp.com/vault/secrets-management/sm-pki-engine',
    permanent: true,
  },
  {
    source: '/guides/secret-mgmt/app-integration',
    destination:
      'https://learn.hashicorp.com/vault/developer/sm-app-integration',
    permanent: true,
  },
  // Encryption
  {
    source: '/guides/encryption',
    destination:
      'https://learn.hashicorp.com/collections/vault/encryption-as-a-service',
    permanent: true,
  },
  {
    source: '/guides/encryption/transit',
    destination:
      'https://learn.hashicorp.com/vault/encryption-as-a-service/eaas-transit',
    permanent: true,
  },
  {
    source: '/guides/encryption/spring-demo',
    destination:
      'https://learn.hashicorp.com/vault/encryption-as-a-service/eaas-spring-demo',
    permanent: true,
  },
  {
    source: '/guides/encryption/transit-rewrap',
    destination:
      'https://learn.hashicorp.com/vault/encryption-as-a-service/eaas-transit-rewrap',
    permanent: true,
  },
  // Intro getting started content -> Learn
  {
    source: '/intro',
    destination:
      'https://learn.hashicorp.com/collections/vault/getting-started',
    permanent: true,
  },
  {
    source: '/intro/getting-started',
    destination: 'https://learn.hashicorp.com/vault/getting-started/install',
    permanent: true,
  },
  {
    source: '/intro/getting-started/install',
    destination:
      'https://learn.hashicorp.com/tutorials/vault/getting-started-install',
    permanent: true,
  },
  {
    source: '/intro/getting-started/dev-server',
    destination: 'https://learn.hashicorp.com/vault/getting-started/dev-server',
    permanent: true,
  },
  {
    source: '/intro/getting-started/first-secret',
    destination:
      'https://learn.hashicorp.com/vault/getting-started/first-secret',
    permanent: true,
  },
  {
    source: '/intro/getting-started/secrets-engines',
    destination:
      'https://learn.hashicorp.com/vault/getting-started/secrets-engines',
    permanent: true,
  },
  {
    source: '/intro/getting-started/dynamic-secrets',
    destination:
      'https://learn.hashicorp.com/vault/getting-started/dynamic-secrets',
    permanent: true,
  },
  {
    source: '/intro/getting-started/help',
    destination: 'https://learn.hashicorp.com/vault/getting-started/help',
    permanent: true,
  },
  {
    source: '/intro/getting-started/authentication',
    destination:
      'https://learn.hashicorp.com/vault/getting-started/authentication',
    permanent: true,
  },
  {
    source: '/intro/getting-started/policies',
    destination: 'https://learn.hashicorp.com/vault/getting-started/policies',
    permanent: true,
  },
  {
    source: '/intro/getting-started/deploy',
    destination: 'https://learn.hashicorp.com/vault/getting-started/deploy',
    permanent: true,
  },
  {
    source: '/intro/getting-started/apis',
    destination: 'https://learn.hashicorp.com/vault/getting-started/apis',
    permanent: true,
  },
  {
    source: '/intro/getting-started/next-steps',
    destination: 'https://learn.hashicorp.com/vault/getting-started/next-steps',
    permanent: true,
  },
  // Rearranged out of `guides` but still on `.io`
  {
    source: '/guides/partnerships',
    destination: '/docs/partnerships',
    permanent: true,
  },
  {
    source: '/intro/use-cases',
    destination: '/docs/use-cases',
    permanent: true,
  },
  { source: '/intro/vs', destination: '/docs/vs', permanent: true },
  { source: '/intro/vs/:path*', destination: '/docs/vs', permanent: true },
  {
    source: '/intro/what-is-vault',
    destination: '/docs/what-is-vault',
    permanent: true,
  },
  {
    source: '/api/:path*',
    destination: '/api-docs/:path*',
    permanent: true,
  },
  // disallow '.html' or '/index.html' in favor of cleaner, simpler paths
  { source: '/:path*/index', destination: '/:path*', permanent: true },
  { source: '/:path*.html', destination: '/:path*', permanent: true },
]
