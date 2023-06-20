/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, currentURL } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';

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

  test('shows the secrets engines card', async function (assert) {
    await enablePage.enable('pki', 'pki');
    await visit('/vault/dashboard');
    assert.dom(SELECTORS.cardTitle).hasText('Secrets Engines');
    assert.dom(SELECTORS.getSecretEngineAccessor('pki')).exists();
  });
});
