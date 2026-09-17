/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentURL, currentRouteName, visit } from '@ember/test-helpers';
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
  test('the persona section renders above the Data & Privacy section', async function (assert) {
    // Verifies that the Persona (Your role) section is present for all users
    // and appears before the Data & Privacy section in the DOM.
    await click(GENERAL.button('user-menu-trigger'));
    await click(GENERAL.menuItem('user-preferences'));

    assert.dom('[data-test-persona-section]').exists('the Persona section renders');
    assert
      .dom('[data-test-data-privacy-section]')
      .exists('the Data & Privacy section renders for non-HVD users');

    const personaTop = document.querySelector('[data-test-persona-section]').getBoundingClientRect().top;
    const privacyTop = document.querySelector('[data-test-data-privacy-section]').getBoundingClientRect().top;
    assert.true(personaTop < privacyTop, 'Persona section appears above Data & Privacy in the DOM');
  });

  test('HVD-managed clusters show the preferences page but hide the Data & Privacy section', async function (assert) {
    // The telemetry toggle is Segment-only and meaningless for HVD (which uses PostHog).
    // The page itself remains accessible so HVD users can still manage other preferences.
    this.owner.lookup('service:flags').featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];

    await visit('/vault/user-preferences');

    assert.strictEqual(currentURL(), '/vault/user-preferences', 'HVD users can access the preferences page');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Preferences', 'page header renders');
    assert
      .dom('[data-test-data-privacy-section]')
      .doesNotExist('Data & Privacy section is hidden for HVD-managed clusters');
  });
});
