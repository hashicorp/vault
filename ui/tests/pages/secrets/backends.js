/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, collection, clickable } from 'ember-cli-page-object';
import uiPanel from 'vault/tests/pages/components/console/ui-panel';

export default create({
  consoleToggle: clickable('[data-test-console-toggle]'),
  visit: visitable('/vault/secrets'),
  rows: collection('[data-test-secrets-backend-link]', {
    menu: clickable('[data-test-popup-menu-trigger]'),
  }),
  configLink: clickable('[data-test-engine-config]', {
    testContainer: '#ember-testing',
  }),
  disableButton: clickable('[data-test-confirm-action-trigger]'),
  confirmDisable: clickable('[data-test-confirm-button]'),
  console: uiPanel,
});
