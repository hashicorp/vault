/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { text, create, collection, clickable, visitable } from 'ember-cli-page-object';
export default create({
  visit: visitable('/vault/policies/:type'),
  policies: collection('[data-test-policy-item]', {
    name: text('[data-test-policy-name]'),
  }),
  row: collection('[data-test-policy-link]', {
    name: text(),
    menu: clickable('[data-test-popup-menu-trigger]'),
  }),
  findPolicyByName(name) {
    return this.policies.filterBy('name', name)[0];
  },
  delete: clickable('[data-test-policy-delete]'),
  confirmDelete: clickable('[data-test-confirm-button]'),
});
