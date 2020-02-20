// The root folder for this documentation category is `pages/api-docs`
//
// - A string refers to the name of a file
// - A "category" value refers to the name of a directory
// - All directories must have an "index.mdx" file to serve as
//   the landing page for the category

export default [
  'index',
  'libraries',
  'relatedtools',
  '------------',
  {
    category: 'secret',
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
          'mssql',
          'mysql-maria',
          'postgresql',
          'oracle'
        ]
      },
      { category: 'gcp' },
      { category: 'gcpkms' },
      { category: 'kmip' },
      {
        category: 'kv',
        content: ['kv-v1', 'kv-v2']
      },
      {
        category: 'identity',
        content: [
          'entity',
          'entity-alias',
          'group',
          'group-alias',
          'tokens',
          'lookup'
        ]
      },
      { category: 'nomad' },
	  { category: 'openldap' },
      { category: 'pki' },
      { category: 'rabbitmq' },
      { category: 'ssh' },
      { category: 'totp' },
      { category: 'transit' },
      '-----------------------',
      { category: 'cassandra' },
      { category: 'mongodb' },
      { category: 'mssql' },
      { category: 'mysql' },
      { category: 'postgresql' }
    ]
  },
  {
    category: 'auth',
    content: [
      { category: 'alicloud' },
      { category: 'approle' },
      { category: 'aws' },
      { category: 'azure' },
      { category: 'cf' },
      { category: 'github' },
      { category: 'gcp' },
      { category: 'jwt' },
      { category: 'kerberos' },
      { category: 'kubernetes' },
      { category: 'ldap' },
      { category: 'oci' },
      { category: 'okta' },
      { category: 'radius' },
      { category: 'cert' },
      { category: 'token' },
      { category: 'userpass' },
      { category: 'app-id' }
    ]
  },
  {
    category: 'system',
    content: [
      'audit',
      'audit-hash',
      'auth',
      'capabilities',
      'capabilities-accessor',
      'capabilities-self',
      'config-auditing',
      'config-control-group',
      'config-cors',
      'config-state',
      'config-ui',
      'control-group',
      'generate-root',
      'health',
      'host-info',
      'init',
      'internal-specs-openapi',
      'internal-ui-mounts',
      'key-status',
      'leader',
      'leases',
      'license',
      'metrics',
      {
        category: 'mfa',
        content: ['duo', 'okta', 'pingid', 'totp']
      },
      'mounts',
      'namespaces',
      'plugins-reload-backend',
      'plugins-catalog',
      'policy',
      'policies',
      'pprof',
      'raw',
      'rekey',
      'rekey-recovery-key',
      'remount',
      {
        category: 'replication',
        content: ['replication-performance', 'replication-dr']
      },
      'rotate',
      'seal',
      'seal-status',
      'sealwrap-rewrap',
      'step-down',
      {
        category: 'storage',
        content: ['raft']
      },
      'tools',
      'unseal',
      'wrapping-lookup',
      'wrapping-rewrap',
      'wrapping-unwrap',
      'wrapping-wrap'
    ]
  }
]
