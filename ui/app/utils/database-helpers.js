export const AVAILABLE_PLUGIN_TYPES = [
  {
    value: 'mongodb-database-plugin',
    displayName: 'MongoDB',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'connection_url' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'write_concern', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'tls', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'tls_ca', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
  {
    value: 'mssql-database-plugin',
    displayName: 'MSSQL',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'connection_url' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
  {
    value: 'mysql-database-plugin',
    displayName: 'MySQL/MariaDB',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'tls', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'tls_ca', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
  {
    value: 'mysql-aurora-database-plugin',
    displayName: 'MySQL (Aurora)',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'tls', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'tls_ca', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
  {
    value: 'mysql-rds-database-plugin',
    displayName: 'MySQL (RDS)',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'tls', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'tls_ca', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
  {
    value: 'mysql-legacy-database-plugin',
    displayName: 'MySQL (Legacy)',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'tls', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'tls_ca', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
  {
    value: 'elasticsearch-database-plugin',
    displayName: 'Elasticsearch',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'url', group: 'pluginConfig' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'ca_cert', group: 'pluginConfig' },
      { attr: 'ca_path', group: 'pluginConfig' },
      { attr: 'client_cert', group: 'pluginConfig' },
      { attr: 'client_key', group: 'pluginConfig' },
      { attr: 'tls_server_name', group: 'pluginConfig' },
      { attr: 'insecure', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
    ],
  },
  {
    value: 'oracle-database-plugin',
    displayName: 'Oracle',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'username', group: 'pluginConfig', show: false },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
      { attr: 'username_template', group: 'pluginConfig' },
      { attr: 'root_rotation_statements', group: 'statements' },
    ],
  },
];

export const ROLE_FIELDS = {
  static: ['username', 'rotation_period'],
  dynamic: ['ttl', 'max_ttl'],
};

export const STATEMENT_FIELDS = {
  static: {
    default: ['rotation_statements'],
    'mongodb-database-plugin': [],
    'mssql-database-plugin': [],
    'mysql-database-plugin': [],
    'mysql-aurora-database-plugin': [],
    'mysql-rds-database-plugin': [],
    'mysql-legacy-database-plugin': [],
    'elasticsearch-database-plugin': [],
    'oracle-database-plugin': [],
  },
  dynamic: {
    default: ['creation_statements', 'revocation_statements', 'rollback_statements', 'renew_statements'],
    'mongodb-database-plugin': ['creation_statement', 'revocation_statement'],
    'mssql-database-plugin': ['creation_statements', 'revocation_statements'],
    'mysql-database-plugin': ['creation_statements', 'revocation_statements'],
    'mysql-aurora-database-plugin': ['creation_statements', 'revocation_statements'],
    'mysql-rds-database-plugin': ['creation_statements', 'revocation_statements'],
    'mysql-legacy-database-plugin': ['creation_statements', 'revocation_statements'],
    'elasticsearch-database-plugin': ['creation_statement'],
    'oracle-database-plugin': ['creation_statements', 'revocation_statements'],
  },
};

export function getStatementFields(type, plugin) {
  if (!type) return null;
  let dbValidFields = STATEMENT_FIELDS[type].default;
  if (STATEMENT_FIELDS[type][plugin]) {
    dbValidFields = STATEMENT_FIELDS[type][plugin];
  }
  return dbValidFields;
}

export function getRoleFields(type) {
  if (!type) return null;
  return ROLE_FIELDS[type];
}
