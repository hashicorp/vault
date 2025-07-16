/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { getStatementFields, getRoleFields } from '../utils/model-helpers/database-helpers';

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
 * @param {object} dbParams - holds database config values, { plugin_name: string [eg 'mongodb-database-plugin'], skip_static_role_rotation_import: boolean }
 */

export default class DatabaseRoleSettingForm extends Component {
  get dbConfig() {
    return this.args.dbParams;
  }

  get settingFields() {
    const dbValues = this.args.dbParams;
    if (!this.args.roleType) return null;
    const dbValidFields = getRoleFields(this.args.roleType);
    return this.args.attrs.filter((a) => {
      // Sets default value for skip_import_rotation based on parent db config value
      if (a.name === 'skip_import_rotation' && this.args.mode === 'create') {
        a.options.defaultValue = dbValues?.skip_static_role_rotation_import;
      }
      return dbValidFields.includes(a.name);
    });
  }

  get statementFields() {
    const type = this.args.roleType;
    if (!type) return null;
    const dbValidFields = getStatementFields(type, this.dbConfig ? this.dbConfig.plugin_name : null);
    return this.args.attrs.filter((a) => {
      return dbValidFields.includes(a.name);
    });
  }

  get isOverridden() {
    if (this.args.mode !== 'create' || !this.dbConfig) return null;

    const dbSkip = this.dbConfig.skip_static_role_rotation_import;
    const staticVal = this.args.model.get('skip_import_rotation');
    return this.args.mode === 'create' && dbSkip !== staticVal;
  }
}
