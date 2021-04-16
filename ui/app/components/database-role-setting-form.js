/**
 * @module DatabaseRoleSettingForm
 * DatabaseRoleSettingForm components are used to handle the role settings section on the database/role form
 *
 * @example
 * ```js
 * <DatabaseRoleSettingForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {Array<object>} attrs - all available attrs from the model to iterate over
 * @param {object} model - ember data model which should be updated on change
 * @param {string} [roleType] - role type controls which attributes are shown
 * @param {string} [mode=create] - mode of the form (eg. create or edit)
 * @param {string} [dbType=default] - type of database, eg 'mongodb-database-plugin'
 */

import Component from '@glimmer/component';

// Below fields are intended to be dynamic based on type of role and db.
// example of usage: FIELDS[roleType][db]
const ROLE_FIELDS = {
  static: ['username', 'rotation_period'],
  dynamic: ['ttl', 'max_ttl'],
};

const STATEMENT_FIELDS = {
  static: {
    default: ['rotation_statements'],
    'mongodb-database-plugin': [],
    'mssql-database-plugin': [],
  },
  dynamic: {
    default: ['creation_statements', 'revocation_statements', 'rollback_statements', 'renew_statements'],
    'mongodb-database-plugin': ['creation_statement', 'revocation_statement'],
    'mssql-database-plugin': ['creation_statements', 'revocation_statements'],
  },
};
export default class DatabaseRoleSettingForm extends Component {
  get settingFields() {
    if (!this.args.roleType) return null;
    let dbValidFields = ROLE_FIELDS[this.args.roleType];
    return this.args.attrs.filter(a => {
      return dbValidFields.includes(a.name);
    });
  }

  get statementFields() {
    const type = this.args.roleType;
    const plugin = this.args.dbType;
    if (!type) return null;
    let dbValidFields = STATEMENT_FIELDS[type].default;
    if (STATEMENT_FIELDS[type][plugin]) {
      dbValidFields = STATEMENT_FIELDS[type][plugin];
    }
    return this.args.attrs.filter(a => {
      return dbValidFields.includes(a.name);
    });
  }
}
