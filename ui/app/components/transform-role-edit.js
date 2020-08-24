/**
 * @module TransformRoleEdit
 * TransformRoleEdit components are used to...
 *
 * @example
 * ```js
 * <TransformRoleEdit @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

// import Component from '@ember/component';

// export default Component.extend({
// });

import RoleEdit from './role-edit';

export default RoleEdit.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'transform');
  },
});
