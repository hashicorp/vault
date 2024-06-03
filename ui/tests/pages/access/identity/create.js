/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/identity/edit-form';

export default create({
  visit: visitable('/vault/access/identity/:item_type/create'),
  editForm,
  async createItem(item_type, type) {
    if (type) {
      return await this.visit({ item_type }).editForm.type(type).submit();
    }
    return await this.visit({ item_type }).editForm.submit();
  },
});
