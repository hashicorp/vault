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
 */

import Component from '@glimmer/component';

export default class DatabaseRoleSettingForm extends Component {
  get allowedFields() {
    const type = this.args.roleType;
    let fields = ['ttl', 'max_ttl'];
    if (type === 'static') {
      fields = ['username', 'rotation_period'];
    }
    if (!type) return null;
    const filtered = this.args.attrs.filter(a => {
      const includes = fields.includes(a.name);
      return includes;
    });
    return filtered;
  }
}
