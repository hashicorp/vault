/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, click, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Acceptance | policies | edit', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
  });

  test('it redirects to list if navigating to root', async function (assert) {
    await visit(`vault/policy/acl/root/edit`);
    assert.strictEqual(
      currentURL(),
      '/vault/policies/acl',
      'navigation to root redirects you to policy list'
    );
  });

  test('it does not show delete for default policy', async function (assert) {
    await visit(`vault/policy/acl/default/edit`);
    assert.dom(GENERAL.confirmButton).doesNotExist('there is no delete button');
  });

  test('it navigates to show when the toggle is clicked', async function (assert) {
    await visit(`vault/policy/acl/default/edit`);
    await click(GENERAL.button('Back to policy'));
    assert.strictEqual(currentURL(), '/vault/policy/acl/default', 'toggle navigates from edit to show');
  });
});
