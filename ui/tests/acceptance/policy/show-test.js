/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | policies | show', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    return login();
  });

  test('it redirects to list if navigating to root', async function (assert) {
    await visit('/vault/policy/acl/root');

    assert.strictEqual(
      currentURL(),
      '/vault/policies/acl',
      'navigation to root show redirects you to policy list'
    );
  });

  test('it navigates to edit when the toggle is clicked', async function (assert) {
    await visit('/vault/policy/acl/default');
    await click(GENERAL.button('Edit policy'));
    assert.strictEqual(currentURL(), '/vault/policy/acl/default/edit', 'toggle navigates to edit page');
  });
});
