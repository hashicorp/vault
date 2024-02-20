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
import { click, waitFor } from '@ember/test-helpers';
import { PAGE as ts } from 'vault/tests/helpers/sync/sync-selectors';

// sync is an enterprise feature but since mirage is used the enterprise label has been intentionally omitted from the module name
module('Acceptance | sync | destination', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    syncScenario(this.server);
    syncHandlers(this.server);
    return authPage.login();
  });

  test('it should transition to correct routes when performing actions', async function (assert) {
    await click(ts.navLink('Secrets Sync'));
    await click(ts.destinations.list.create);
    await click(ts.createCancel);
    await click(ts.overviewCard.actionLink('Create new'));
    await click(ts.createCancel);
    await waitFor(ts.overview.table.actionToggle(0));
    await click(ts.overview.table.actionToggle(0));
    await click(ts.overview.table.action('sync'));
    await click(ts.destinations.sync.cancel);
    await click(ts.breadcrumbLink('Secrets Sync'));
    await waitFor(ts.overview.table.actionToggle(0));
    await click(ts.overview.table.actionToggle(0));
    await click(ts.overview.table.action('details'));
    assert.dom(ts.tab('Secrets')).hasClass('active', 'Navigates to secrets view for destination');
  });
});
