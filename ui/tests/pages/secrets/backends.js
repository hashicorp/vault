/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable, collection, clickable, text } from 'ember-cli-page-object';
import uiPanel from 'vault/tests/pages/components/console/ui-panel';

export default create({
  consoleToggle: clickable('[data-test-console-toggle]'),
  visit: visitable('/vault/secrets'),
  rows: collection('[data-test-table-row]', {
    path: text('[data-test-table-data]'),
    menu: clickable('[data-test-popup-menu-trigger]'),
  }),
  console: uiPanel,
});
