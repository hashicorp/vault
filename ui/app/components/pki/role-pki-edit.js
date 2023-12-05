/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RoleEdit from '../role-edit';

export default RoleEdit.extend({
  init() {
    this._super(...arguments);
    this.set('backendType', 'pki');
  },
});
