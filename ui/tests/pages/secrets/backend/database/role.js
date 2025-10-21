/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, fillable, visitable, selectable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets-engines/:backend/list?itemType=role'),
  visitShow: visitable('/vault/secrets-engines/:backend/show/role/:id'),
  visitCreate: visitable('/vault/secrets-engines/:backend/create?itemType=role'),
  createLink: clickable('[data-test-secret-create]'),
  name: fillable('[data-test-input="name"]'),
  roleType: selectable('[data-test-input="type"'),
  edit: clickable('[data-test-edit-link]'),
});
