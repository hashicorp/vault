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
