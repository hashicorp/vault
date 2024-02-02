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
import { click, visit, currentURL, fillIn } from '@ember/test-helpers';
import { PAGE as ts } from 'vault/tests/helpers/sync/sync-selectors';

module('Acceptance | enterprise | sync | destination', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    syncScenario(this.server);
    syncHandlers(this.server);
    return authPage.login();
  });

  test('it should transition to overview route via breadcrumb', async function (assert) {
    await visit('vault/sync/secrets/destinations/aws-sm/destination-aws/secrets');
    await click(ts.breadcrumbAtIdx(0));
    assert.strictEqual(
      currentURL(),
      '/vault/sync/secrets/overview',
      'Transitions to overview on breadcrumb click'
    );
  });

  test('it should transition to correct routes when performing actions', async function (assert) {
    await click(ts.navLink('Secrets Sync'));
    await click(ts.tab('Destinations'));
    await click(ts.listItem);
    assert.dom(ts.tab('Secrets')).hasClass('active', 'Secrets tab is active');

    await click(ts.tab('Details'));
    assert.dom(ts.infoRowLabel('Name')).exists('Destination details display');

    await click(ts.toolbar('Sync secrets'));
    await click(ts.destinations.sync.cancel);

    await click(ts.toolbar('Edit destination'));
    assert.dom(ts.inputByAttr('name')).isDisabled('Edit view renders with disabled name field');
    await click(ts.cancelButton);
    assert.dom(ts.tab('Details')).hasClass('active', 'Details view is active');
  });

  test('it should delete destination', async function (assert) {
    await visit('vault/sync/secrets/destinations/aws-sm/destination-aws/details');
    await click(ts.toolbar('Delete destination'));
    await fillIn(ts.confirmModalInput, 'DELETE');
    await click(ts.confirmButton);
    assert.dom(ts.destinations.deleteBanner).exists('Delete banner renders');
  });
});
