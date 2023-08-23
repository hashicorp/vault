/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, fillable, visitable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/list?tab=alphabet'),
  visitCreate: visitable('/vault/secrets/:backend/create?itemType=alphabet'),
  createLink: clickable('[data-test-secret-create]'),
  editLink: clickable('[data-test-edit-link]'),
  name: fillable('[data-test-input="name"]'),
  alphabet: fillable('[data-test-input="alphabet"'),
  submit: clickable('[data-test-alphabet-transform-create]'),
});
