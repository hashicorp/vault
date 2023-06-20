/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';

module('Acceptance | landing page dashboard', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return authPage.login();
  });

  // TODO LANDING PAGE: create a test that will navigate to dashboard if user opts into new dashboard ui
  test('does not navigate to dashboard on login when user has not opted into dashboard ui', async function (assert) {
    assert.strictEqual(currentURL(), '/vault/secrets');

    await visit('/vault/dashboard');

    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  module('secrets engines card', function () {
    test('shows a secrets engine card', async function (assert) {
      await mountSecrets.enable('pki', 'pki');
      await visit('/vault/dashboard');
      assert.dom(SELECTORS.cardTitle).hasText('Secrets Engines');
      assert.dom(SELECTORS.getSecretEngineAccessor('pki')).exists();
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await mountSecrets.enable('nomad', 'nomad');
      await visit('/vault/dashboard');
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
    });
  });
});
