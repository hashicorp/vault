/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, currentRouteName } from '@ember/test-helpers';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | user-preferences', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Logging in lands the user on the cluster dashboard.
    await login();
  });

  test('a user navigates to User Preferences from the account menu', async function (assert) {
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'starts on the dashboard');

    await click(GENERAL.button('user-menu-trigger'));
    assert
      .dom(GENERAL.menuItem('user-preferences'))
      .hasText('User preferences', 'the account menu shows the User preferences item');

    await click(GENERAL.menuItem('user-preferences'));

    assert.strictEqual(currentURL(), '/vault/user-preferences', 'lands on the user-preferences route');
    assert
      .dom(GENERAL.button('user-menu-trigger'))
      .hasAttribute('aria-expanded', 'false', 'the dropdown closes after navigating');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('User preferences', 'the page header renders its title');
  });
});
