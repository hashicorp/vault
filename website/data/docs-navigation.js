// The root folder for this documentation category is `pages/docs`
//
// - A string refers to the name of a file
// - A "category" value refers to the name of a directory
// - All directories must have an "index.mdx" file to serve as
//   the landing page for the category

export default [
  {
    category: 'install'
  },
  {
    category: 'internals',
    content: [
      'architecture',
      'high-availability',
      'integrated-storage',
      'security',
      'telemetry',
      'token',
      'rotation',
      'replication',
      'plugins'
    ]
  },
  {
    category: 'concepts',
    content: [
      'dev-server',
      'seal',
      'lease',
      'auth',
      'tokens',
      'response-wrapping',
      'policies',
      'ha',
      'pgp-gpg-keybase',
      'recovery-mode'
    ]
  },
  {
    category: 'configuration',
    content: [
      {
        category: 'listener',
        content: ['tcp']
      },
      {
        category: 'seal',
        content: [
          'alicloudkms',
          'awskms',
          'azurekeyvault',
          'gcpckms',
          'ocikms',
          'pkcs11',
          'transit'
        ]
      },
      {
        category: 'storage',
        content: [
          'azure',
          'cassandra',
          'cockroachdb',
          'consul',
          'couchdb',
          'dynamodb',
          'etcd',
          'filesystem',
          'foundationdb',
          'google-cloud-spanner',
          'google-cloud-storage',
          'in-memory',
          'manta',
          'mssql',
          'mysql',
          'oci-object-storage',
          'postgresql',
          'raft',
          's3',
          'swift',
          'zookeeper'
        ]
      },
      {
        category: 'service-registration',
        content: [
          'consul',
          'kubernetes'
        ]
      },
      'telemetry',
      { category: 'ui' },
      { category: 'entropy-augmentation' }
    ]
  },
  {
    category: 'commands',
    content: [
      'agent',
      {
        category: 'audit',
        content: ['disable', 'enable', 'list']
      },
      {
        category: 'auth',
        content: ['disable', 'enable', 'help', 'list', 'tune']
      },
      'debug',
      'delete',
      {
        category: 'kv',
        content: [
          'delete',
          'destroy',
          'enable-versioning',
          'get',
          'list',
          'metadata',
          'patch',
          'put',
          'rollback',
          'undelete'
        ]
      },
      {
        category: 'lease',
        content: ['renew', 'revoke']
      },
      'list',
      'login',
      'namespace',
      {
        category: 'operator',
        content: [
          'generate-root',
          'init',
          'key-status',
          'migrate',
          'rekey',
          'rotate',
          'seal',
          'step-down',
          'unseal'
        ]
      },
      'path-help',
      {
        category: 'plugin',
        content: ['deregister', 'info', 'list', 'register']
      },
      {
        category: 'policy',
        content: ['delete', 'fmt', 'list', 'read', 'write']
      },
      'read',
      {
        category: 'secrets',
        content: ['disable', 'enable', 'list', 'move', 'tune']
      },
      'server',
      'ssh',
      'status',
      {
        category: 'token',
        content: ['capabilities', 'create', 'lookup', 'renew', 'revoke']
      },
      'unwrap',
      'version',
      'write',
      'token-helper'
    ]
  },
  {
    category: 'agent',
    content: [
      {
        category: 'autoauth',
        content: [
          {
            category: 'methods',
            content: [
              'alicloud',
              'approle',
              'aws',
              'azure',
              'cert',
              'cf',
              'gcp',
              'jwt',
              'kerberos',
              'kubernetes'
            ]
          },
          {
            category: 'sinks',
            content: ['file']
          }
        ]
      },
      { category: 'caching' },
      { category: 'template' }
    ]
  },
  '----------------',
  {
    category: 'secrets',
    content: [
      { category: 'ad' },
      { category: 'alicloud' },
      { category: 'aws' },
      { category: 'azure' },
      { category: 'consul' },
      { category: 'cubbyhole' },
      {
        category: 'databases',
        content: [
          'cassandra',
          'elasticdb',
          'influxdb',
          'hanadb',
          'mongodb',
          'mongodbatlas',
          'mssql',
          'mysql-maria',
          'postgresql',
          'oracle',
          'custom'
        ]
      },
      { category: 'gcp' },
      { category: 'gcpkms' },
      { category: 'kmip' },
      {
        category: 'kv',
        content: ['kv-v1', 'kv-v2']
      },
      { category: 'identity' },
      { category: 'mongodbatlas' },
      { category: 'nomad' },
      { category: 'openldap' },
      { category: 'pki' },
      { category: 'rabbitmq' },
      {
        category: 'ssh',
        content: [
          'signed-ssh-certificates',
          'one-time-ssh-passwords',
          'dynamic-ssh-keys'
        ]
      },
      { category: 'totp' },
      { category: 'transit' },
      '------------------------',
      { category: 'cassandra' },
      { category: 'mongodb' },
      { category: 'mssql' },
      { category: 'mysql' }
    ]
  },
  {
    category: 'auth',
    content: [
      'approle',
      'alicloud',
      'aws',
      'azure',
      'cf',
      'gcp',
      'jwt',
      'kerberos',
      'kubernetes',
      'github',
      'ldap',
      'oci',
      'okta',
      'radius',
      'cert',
      'token',
      'userpass',
      '---------',
      'app-id',
      'mfa'
    ]
  },
  {
    category: 'audit',
    content: ['file', 'syslog', 'socket']
  },
  {
    category: 'plugin'
  },
  '----------------',
  {
    category: 'what-is-vault'
  },
  {
    category: 'use-cases'
  },
  {
    category: 'partnerships'
  },
  '----------------',
  {
    category: 'platform',
    content: [
      {
        category: 'k8s',
        content: [
          {
            category: 'helm',
            content: [
              'run',
              'configuration',
              {
                category: 'examples',
                content: [
                  'standalone-load-balanced-ui',
                  'standalone-tls',
                  'standalone-audit',
                  'ha-with-consul',
                  'kubernetes-auth',
                  'enterprise-best-practice'
                ]
              },
            ],
          },
          {
            category: 'injector',
            content: ['annotations', 'installation', 'examples']
          }
        ]
      },
      {
        category: 'aws-mp',
        content: ['run']
      }
    ]
  },
  '----------------',
  {
    category: 'upgrading',
    content: [
      'plugins',
      'upgrade-to-0.5.0',
      'upgrade-to-0.5.1',
      'upgrade-to-0.6.0',
      'upgrade-to-0.6.1',
      'upgrade-to-0.6.2',
      'upgrade-to-0.6.3',
      'upgrade-to-0.6.4',
      'upgrade-to-0.7.0',
      'upgrade-to-0.8.0',
      'upgrade-to-0.9.0',
      'upgrade-to-0.9.1',
      'upgrade-to-0.9.2',
      'upgrade-to-0.9.3',
      'upgrade-to-0.9.6',
      'upgrade-to-0.10.0',
      'upgrade-to-0.10.2',
      'upgrade-to-0.10.4',
      'upgrade-to-0.11.0',
      'upgrade-to-0.11.2',
      'upgrade-to-0.11.6',
      'upgrade-to-1.0.0',
      'upgrade-to-1.1.0',
      'upgrade-to-1.1.1',
      'upgrade-to-1.1.2',
      'upgrade-to-1.2.0',
      'upgrade-to-1.2.1',
      'upgrade-to-1.2.4',
      'upgrade-to-1.3.0'
    ]
  },
  '----------------',
  {
    category: 'enterprise',
    content: [
      { category: 'replication' },
      {
        category: 'hsm',
        content: ['behavior', 'security']
      },
      { category: 'entropy-augmentation' },
      { category: 'sealwrap' },
      { category: 'namespaces' },
      { category: 'performance-standby' },
      { category: 'control-groups' },
      {
        category: 'mfa',
        content: ['mfa-duo', 'mfa-okta', 'mfa-pingid', 'mfa-totp']
      },
      {
        category: 'sentinel',
        content: ['examples', 'properties']
      }
    ]
  }
]
