/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import RoleEdit from '../role-edit';

export default RoleEdit.extend({
  actions: {
    delete() {
      this.model.save({ adapterOptions: { method: 'revoke' } });
    },
  },
});
