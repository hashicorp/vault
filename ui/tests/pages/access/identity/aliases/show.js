/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, collection, contains, visitable } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';
import infoTableRow from 'vault/tests/pages/components/info-table-row';

export default create({
  visit: visitable('/vault/access/identity/:item_type/aliases/:alias_id'),
  flashMessage,
  nameContains: contains('[data-test-alias-name]'),
  rows: collection('[data-test-component="info-table-row"]', infoTableRow),
  edit: clickable('[data-test-alias-edit-link]'),
});
