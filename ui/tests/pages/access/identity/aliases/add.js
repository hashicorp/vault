/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/identity/edit-form';

export default create({
  visit: visitable('/vault/access/identity/:item_type/aliases/add/:id'),
  editForm,
});
