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
  static: {
    default: ['ttl', 'max_ttl', 'username', 'rotation_period'],
    'mongodb-database-plugin': ['username', 'rotation_period'],
  },
  dynamic: {
    default: ['ttl', 'max_ttl', 'username', 'rotation_period'],
    'mongodb-database-plugin': ['ttl', 'max_ttl'],
  },
};

const STATEMENT_FIELDS = {
  static: {
    default: ['creation_statements', 'revocation_statements', 'rotation_statements'],
    'mongodb-database-plugin': 'NONE', // will not show the section
  },
  dynamic: {
    default: ['creation_statements', 'revocation_statements', 'rotation_statements'],
    'mongodb-database-plugin': ['creation_statement'],
  },
};

export default class DatabaseRoleSettingForm extends Component {
  get settingFields() {
    const type = this.args.roleType;
    if (!type) return null;
    const db = this.args.dbType || 'default';
    const fields = ROLE_FIELDS[type][db];
    if (!Array.isArray(fields)) return fields;
    const filtered = this.args.attrs.filter(a => {
      const includes = fields.includes(a.name);
      return includes;
    });
    return filtered;
  }

  get statementFields() {
    const type = this.args.roleType;
    if (!type) return null;
    const db = this.args.dbType || 'default';
    const fields = STATEMENT_FIELDS[type][db];
    if (!Array.isArray(fields)) return fields;
    const filtered = this.args.attrs.filter(a => {
      const includes = fields.includes(a.name);
      return includes;
    });
    return filtered;
  }
}
