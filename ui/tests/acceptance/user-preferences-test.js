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
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Preferences', 'the page header renders its title');
  });

  test('the page header states browser-only storage and does not claim entity/cross-device persistence', async function (assert) {
    const terms = ['identity', 'entity', 'user', 'cross-device'];

    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'starts on the dashboard');

    await click(GENERAL.button('user-menu-trigger'));
    await click(GENERAL.menuItem('user-preferences'));

    assert.dom(GENERAL.hdsPageHeaderDescription).includesText('Privacy Policy');
    assert
      .dom(`${GENERAL.hdsPageHeaderDescription} a`)
      .hasAttribute('href', 'https://www.hashicorp.com/privacy');
    assert
      .dom(GENERAL.hdsPageHeaderDescription)
      .includesText('stored in this browser only', 'page header states browser-only storage');

    terms.forEach((s) =>
      assert.dom(GENERAL.hdsPageHeaderDescription).doesNotIncludeText(s, `page header makes no ${s} claim`)
    );
  });
});
