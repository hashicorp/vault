/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import syncHandler from 'vault/mirage/handlers/sync';
import { STATIC_NOW } from 'vault/mirage/handlers/clients';
import { visit, click, currentURL } from '@ember/test-helpers';
import sinon from 'sinon';
import timestamp from 'core/utils/timestamp';
import authPage from 'vault/tests/pages/auth';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CLIENT_COUNT, CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';

module('Acceptance | clients | sync | activated', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
  });

  hooks.beforeEach(async function () {
    syncHandler(this.server);
    await authPage.login();
    return visit('/vault/clients/counts/sync');
  });

  test('it should render charts when secrets sync is activated', async function (assert) {
    syncHandler(this.server);
    assert.dom(CHARTS.chart('Secrets sync usage')).exists('Secrets sync usage chart is rendered');
    assert.dom(CLIENT_COUNT.statText('Total sync clients')).exists('Total sync clients chart is rendered');
    assert.dom(GENERAL.emptyStateTitle).doesNotExist();
  });
});

module('Acceptance | clients | sync | not activated', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
  });

  hooks.beforeEach(async function () {
    await authPage.login();
    return visit('/vault/clients/counts/sync');
  });

  test('it should show an empty state when secrets sync is not activated', async function (assert) {
    assert.expect(3);

    // ensure secret_syncs clients activity is 0
    this.server.get('/sys/internal/counters/activity', () => {
      // return only the things that determine whether to show/hide secrets sync
      return {
        data: {
          total: {
            secret_syncs: 0,
          },
        },
      };
    });

    this.server.get('/sys/activation-flags', () => {
      assert.true(true, '/sys/activation-flags/ is called to check if secrets-sync is activated');

      return {
        data: {
          activated: [],
          unactivated: ['secrets-sync'],
        },
      };
    });

    assert.dom(GENERAL.emptyStateTitle).exists('Shows empty state when secrets-sync is not activated');

    await click(`${GENERAL.emptyStateActions} .hds-link-standalone`);
    assert.strictEqual(
      currentURL(),
      '/vault/sync/secrets/overview',
      'action button navigates to secrets sync overview page'
    );
  });
});
