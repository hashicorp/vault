/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, fillable, visitable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets/:backend/kmip/scopes'),
  visitCreate: visitable('/vault/secrets/:backend/kmip/scopes/create'),
  createLink: clickable('[data-test-scope-create]'),
  scopeName: fillable('[data-test-input="name"]'),
  submit: clickable('[data-test-edit-form-submit]'),
  configurationLink: clickable('[data-test-kmip-link-config]'),
  configureLink: clickable('[data-test-kmip-link-configure]'),
  scopesLink: clickable('[data-test-kmip-link-scopes]'),
});
