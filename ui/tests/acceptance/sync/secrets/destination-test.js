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
import { click, visit, currentURL } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';

const { breadcrumbAtIdx } = PAGE;

module('Acceptance | sync | destination', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    syncScenario(this.server);
    syncHandlers(this.server);
    return authPage.login();
  });

  test('it should transition to overview route via breadcrumb', async function (assert) {
    await visit('vault/sync/secrets/destinations/aws-sm/destination-aws/secrets');
    await click(breadcrumbAtIdx(0));
    assert.strictEqual(
      currentURL(),
      '/vault/sync/secrets/overview',
      'Transitions to overview on breadcrumb click'
    );
  });
});
