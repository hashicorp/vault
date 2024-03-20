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
import { runCmd } from 'vault/tests/helpers/commands';

// sync is an enterprise feature but since mirage is used the enterprise label has been intentionally omitted from the module name
module('Acceptance | sync | overview', function (hooks) {
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

  module('enterprise with namespaces', function (hooks) {
    hooks.beforeEach(async function () {
      await runCmd(`write sys/namespaces/admin -f`, false);
      await authPage.loginNs('admin');
      await runCmd(`write sys/namespaces/foo -f`, false);
      await authPage.loginNs('admin/foo');
    });
    test('it should make activation-flag requests to correct namespace', async function (assert) {
      assert.expect(6);
      this.server.get('/sys/activation-flags', (_, req) => {
        assert.deepEqual(req.requestHeaders, {}, 'Request is unauthenticated and in root namespace');
        return {
          data: {
            activated: [''],
            unactivated: ['secrets-sync'],
          },
        };
      });
      this.server.post('/sys/activation-flags/secrets-sync/activate', (_, req) => {
        assert.strictEqual(
          req.requestHeaders['X-Vault-Namespace'],
          undefined,
          'Request is made to root namespace'
        );
        return {};
      });

      assert.dom('[data-test-badge-namespace]').hasText('foo'); // confirm we're in admin/foo
      await click(ts.navLink('Secrets Sync'));
      assert.dom(ts.overview.optInBanner).exists('Opt-in banner is shown');
      await click(ts.overview.optInBannerEnable);
      assert.dom(ts.overview.optInModal).exists('Opt-in modal is shown');
      assert
        .dom(ts.overview.optInConfirm)
        .isDisabled('Confirm button is disabled when checkbox is unchecked');
      await click(ts.overview.optInCheck);
      await click(ts.overview.optInConfirm);
    });

    test('it should make activation-flag requests to correct namespace when managed', async function (assert) {
      assert.expect(6);
      this.owner.lookup('service:feature-flag').setFeatureFlags(['VAULT_CLOUD_ADMIN_NAMESPACE']);
      this.server.get('/sys/activation-flags', (_, req) => {
        assert.deepEqual(req.requestHeaders, {}, 'Request is unauthenticated and in root namespace');
        return {
          data: {
            activated: [''],
            unactivated: ['secrets-sync'],
          },
        };
      });
      this.server.post('/sys/activation-flags/secrets-sync/activate', (_, req) => {
        assert.strictEqual(
          req.requestHeaders['X-Vault-Namespace'],
          'admin',
          'Request is made to admin namespace'
        );
        return {};
      });

      assert.dom('[data-test-badge-namespace]').hasText('foo'); // confirm we're in admin/foo
      await click(ts.navLink('Secrets Sync'));
      assert.dom(ts.overview.optInBanner).exists('Opt-in banner is shown');
      await click(ts.overview.optInBannerEnable);
      assert.dom(ts.overview.optInModal).exists('Opt-in modal is shown');
      assert
        .dom(ts.overview.optInConfirm)
        .isDisabled('Confirm button is disabled when checkbox is unchecked');
      await click(ts.overview.optInCheck);
      await click(ts.overview.optInConfirm);
    });
  });
});
