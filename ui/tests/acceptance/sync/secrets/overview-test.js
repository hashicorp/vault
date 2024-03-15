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
import { click } from '@ember/test-helpers';
import { PAGE as ts } from 'vault/tests/helpers/sync/sync-selectors';
import AdapterError from '@ember-data/adapter/error';

// sync is an enterprise feature but since mirage is used the enterprise label has been intentionally omitted from the module name
module('Acceptance | enterprise | sync | overview', function (hooks) {
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
    await click(ts.overview.table.actionToggle(0));
    await click(ts.overview.table.action('sync'));
    await click(ts.destinations.sync.cancel);
    await click(ts.breadcrumbLink('Secrets Sync'));
    await click(ts.overview.table.actionToggle(0));
    await click(ts.overview.table.action('details'));
    assert.dom(ts.tab('Secrets')).hasClass('active', 'Navigates to secrets view for destination');
  });

  test('it should show opt-in banner and modal if secrets-sync is not activated', async function (assert) {
    assert.expect(6);
    this.server.get('/sys/activation-flags', () => {
      assert.ok(true, 'Request on initial load to check if secrets-sync is activated');
      return {
        data: {
          activated: [''],
          unactivated: ['secrets-sync'],
        },
      };
    });
    this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
      assert.ok(true, 'Request made to activate secrets-sync');
      return {};
    });
    await click(ts.navLink('Secrets Sync'));
    assert.dom(ts.overview.optInBanner).exists('Opt-in banner is shown');
    await click(ts.overview.optInBannerEnable);
    assert.dom(ts.overview.optInModal).exists('Opt-in modal is shown');
    assert.dom(ts.overview.optInConfirm).isDisabled('Confirm button is disabled when checkbox is unchecked');
    await click(ts.overview.optInCheck);
    await click(ts.overview.optInConfirm);
  });

  test('it should show adapter error if call to activated-features fails', async function (assert) {
    assert.expect(2);
    this.server.get('/sys/activation-flags', () => {
      assert.ok(true, 'Request on initial load to check if secrets-sync is activated');
      return AdapterError.create();
    });
    await click(ts.navLink('Secrets Sync'));
    assert.dom(ts.overview.optInBannerEnableError).exists('Adapter error message is shown');
  });
});
