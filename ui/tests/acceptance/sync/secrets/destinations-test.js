/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import syncScenario from 'vault/mirage/scenarios/sync';
import syncHandlers from 'vault/mirage/handlers/sync';
import authPage from 'vault/tests/pages/auth';
import { click, visit, fillIn } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';

const { searchSelect, filter, listItem } = PAGE;

module('Acceptance | sync | destinations', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    syncScenario(this.server);
    syncHandlers(this.server);
    return authPage.login();
  });

  test('it should filter destinations list', async function (assert) {
    await visit('vault/sync/secrets/destinations');
    assert.dom(listItem).exists({ count: 6 }, 'All destinations render');
    await click(`${filter('type')} .ember-basic-dropdown-trigger`);
    await click(searchSelect.option());
    assert.dom(listItem).exists({ count: 2 }, 'Destinations are filtered by type');
    await fillIn(filter('name'), 'new');
    assert.dom(listItem).exists({ count: 1 }, 'Destinations are filtered by type and name');
    await click(searchSelect.removeSelected);
    await fillIn(filter('name'), 'gcp');
    assert.dom(listItem).exists({ count: 1 }, 'Destinations are filtered by name');
  });
});
