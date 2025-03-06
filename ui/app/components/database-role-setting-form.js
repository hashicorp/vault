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
 * @param {object} dbParams - holds database config values, (plugin_name [eg 'mongodb-database-plugin'], skip_static_role_rotation_import)
 */

export default class DatabaseRoleSettingForm extends Component {
  get settingFields() {
    this.setSkipImport();
    if (!this.args.roleType) return null;
    const dbValidFields = getRoleFields(this.args.roleType);
    return this.args.attrs.filter((a) => {
      return dbValidFields.includes(a.name);
    });
  }

  get statementFields() {
    const type = this.args.roleType;
    const params = this.args.dbParams;
    if (!type || !params) return null;
    const dbValidFields = getStatementFields(type, params.plugin_name);
    return this.args.attrs.filter((a) => {
      return dbValidFields.includes(a.name);
    });
  }

  /**
   * Sets default value for skip_import_rotation based on parent db config value
   */
  setSkipImport() {
    const params = this.args.dbParams;
    if (!params) return;
    const skipInput = this.args.attrs.find((x) => x.name === 'skip_import_rotation');
    skipInput.options.defaultValue = params.skip_static_role_rotation_import;
  }

  get isOverridden() {
    const params = this.args.dbParams;
    if (this.args.mode !== 'create' || !params) return null;

    const dbSkip = params.skip_static_role_rotation_import;
    const staticVal = this.args.model.get('skip_import_rotation');
    return this.args.mode === 'create' && dbSkip !== staticVal;
  }
}
