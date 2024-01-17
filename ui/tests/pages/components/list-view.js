/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { text, isPresent, collection, clickable } from 'ember-cli-page-object';

export default {
  isEmpty: isPresent('[data-test-component="empty-state"]'),
  listItemLinks: collection('[data-test-list-item-link]', {
    text: text(),
    click: clickable(),
    menuToggle: clickable('[data-test-popup-menu-trigger]'),
  }),
  listItems: collection('[data-test-list-item]', {
    text: text(),
    menuToggle: clickable('[data-test-popup-menu-trigger]'),
  }),
  menuItems: collection('.ember-basic-dropdown-content li', {
    testContainer: '#ember-testing',
  }),
  delete: clickable('[data-test-confirm-action-trigger]', {
    testContainer: '#ember-testing',
  }),
  confirmDelete: clickable('[data-test-confirm-button]', {
    testContainer: '#ember-testing',
  }),
};
