/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, collection, fillable, text, visitable, value, clickable } from 'ember-cli-page-object';
import { getter } from 'ember-cli-page-object/macros';

export default create({
  visit: visitable('/vault/secrets-engines/:backend/list/:id'),
  visitRoot: visitable('/vault/secrets-engines/:backend/list'),
  tabs: collection('[data-test-secret-list-tab]'),
  filterInput: fillable('[data-test-nav-input] input'),
  filterInputValue: value('[data-test-nav-input] input'),
  secrets: collection('[data-test-secret-link]', {
    menuToggle: clickable('[data-test-popup-menu-trigger]'),
    id: text(),
    click: clickable(),
  }),
  menuItems: collection('.ember-basic-dropdown-content li', {
    testContainer: '#ember-testing',
  }),
  backendIsEmpty: getter(function () {
    return this.secrets.length === 0;
  }),
});
