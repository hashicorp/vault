/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import RoleEdit from '../role-edit';

export default RoleEdit.extend({
  actions: {
    delete() {
      this.model.save({ adapterOptions: { method: 'revoke' } });
    },
  },
});
