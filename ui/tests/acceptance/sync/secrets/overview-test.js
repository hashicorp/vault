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
import { click, visit, currentURL, waitFor } from '@ember/test-helpers';
import { PAGE as ts } from 'vault/tests/helpers/sync/sync-selectors';
import { runCmd } from 'vault/tests/helpers/commands';

// sync is an enterprise feature but since mirage is used the enterprise label has been intentionally omitted from the module name
module('Acceptance | enterprise | sync | overview', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    syncHandlers(this.server);
    this.version = this.owner.lookup('service:version');
    this.version.features = ['Secrets Sync'];

    await authPage.login();
  });

  module('when feature is activated', function (hooks) {
    hooks.beforeEach(async function () {
      syncScenario(this.server);
    });

    test('it fetches destinations and associations', async function (assert) {
      assert.expect(2);

      this.server.get('/sys/sync/destinations', () => {
        assert.true(true, 'destinations is called');
      });
      this.server.get('/sys/sync/associations', () => {
        assert.true(true, 'associations is called');
      });

      await visit('/vault/sync/secrets/overview');
    });

    module('when there are pre-existing destinations', function (hooks) {
      hooks.beforeEach(async function () {
        syncScenario(this.server);
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
  });

  module('when feature is not activated', function (hooks) {
    hooks.beforeEach(async function () {
      let wasActivatePOSTCalled = false;

      // simulate the feature being activated once /secrets-sync/activate has been called
      this.server.get('/sys/activation-flags', () => {
        if (wasActivatePOSTCalled) {
          return {
            data: {
              activated: ['secrets-sync'],
              unactivated: [''],
            },
          };
        } else {
          return {
            data: {
              activated: [''],
              unactivated: ['secrets-sync'],
            },
          };
        }
      });

      this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
        wasActivatePOSTCalled = true;
        return {};
      });
    });

    test('it does not fetch destinations and associations', async function (assert) {
      assert.expect(0);

      this.server.get('/sys/sync/destinations', () => {
        assert.true(false, 'destinations is not called');
      });
      this.server.get('/sys/sync/associations', () => {
        assert.true(false, 'associations is not called');
      });

      await visit('/vault/sync/secrets/overview');
    });

    test('the activation workflow works', async function (assert) {
      await visit('/vault/sync/secrets/overview');

      assert
        .dom(ts.cta.button)
        .doesNotExist('create first destination is not available until feature has been activated');

      assert.dom(ts.overview.optInBanner).exists();
      await click(ts.overview.optInBannerEnable);

      assert.dom(ts.overview.optInModal).exists('modal to opt-in and activate feature is shown');
      await click(ts.overview.optInCheck);
      await click(ts.overview.optInConfirm);

      assert.dom(ts.overview.optInModal).doesNotExist('modal is gone once activation has been submitted');
      assert
        .dom(ts.overview.optInBanner)
        .doesNotExist('opt-in banner is gone once activation has been submitted');

      await click(ts.cta.button);
      assert.strictEqual(
        currentURL(),
        '/vault/sync/secrets/destinations/create',
        'create new destination is available once feature is activated'
      );
    });
  });

  module('enterprise with namespaces', function (hooks) {
    hooks.beforeEach(async function () {
      this.version.features = ['Secrets Sync', 'Namespaces'];
      await runCmd(`write sys/namespaces/admin -f`, false);
      await authPage.loginNs('admin');
      await runCmd(`write sys/namespaces/foo -f`, false);
      await authPage.loginNs('admin/foo');
    });

    test('it should make activation-flag requests to correct namespace', async function (assert) {
      assert.expect(3);

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

      // confirm we're in admin/foo
      assert.dom('[data-test-badge-namespace]').hasText('foo');

      await click(ts.navLink('Secrets Sync'));
      await click(ts.overview.optInBannerEnable);
      await click(ts.overview.optInCheck);
      await click(ts.overview.optInConfirm);
    });

    test.skip('it should make activation-flag requests to correct namespace when managed', async function (assert) {
      // TODO: unskip for 1.16.1 when managed is supported
      assert.expect(3);
      this.owner.lookup('service:flags').setFeatureFlags(['VAULT_CLOUD_ADMIN_NAMESPACE']);

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

      // confirm we're in admin/foo
      assert.dom('[data-test-badge-namespace]').hasText('foo');

      await click(ts.navLink('Secrets Sync'));
      await click(ts.overview.optInBannerEnable);
      await click(ts.overview.optInCheck);
      await click(ts.overview.optInConfirm);
    });
  });
});
