/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, fillable, visitable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets-engines/:backend/kmip/scopes/:scope/roles'),
  visitDetail: visitable('/vault/secrets-engines/:backend/kmip/scopes/:scope/roles/:role'),
  create: clickable('[data-test-role-create]'),
  roleName: fillable('[data-test-input="name"]'),
  detailEditLink: clickable('[data-test-kmip-link-edit-role]'),
  cancelLink: clickable('[data-test-edit-form-cancel]'),
});
