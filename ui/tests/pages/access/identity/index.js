/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, text, visitable, collection } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';

export default create({
  visit: visitable('/vault/access/identity/:item_type'),
  flashMessage,
  items: collection('[data-test-identity-row]', {
    menu: clickable('[data-test-popup-menu-trigger]'),
    name: text('[data-test-identity-link]'),
  }),

  delete: clickable('[data-test-popup-menu="delete"]', {
    testContainer: '#ember-testing',
  }),
  confirmDelete: clickable('[data-test-confirm-button]'),
});
