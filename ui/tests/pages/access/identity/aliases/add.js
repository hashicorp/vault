/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { create, visitable } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/identity/edit-form';

export default create({
  visit: visitable('/vault/access/identity/:item_type/aliases/add/:id'),
  editForm,
});
