/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { currentURL, click, fillIn, waitFor } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { createPolicyCmd, deleteAuthCmd, mountAuthCmd, runCmd } from '../helpers/commands';
import { create } from 'ember-cli-page-object';
import fm from 'vault/tests/pages/components/flash-message';

const flashMessage = create(fm);
const resetPolicy = `
path "auth/userpass/users/reset-me/password" {
  capabilities = ["update", "create"]
}
`;
module('Acceptance | reset password', function (hooks) {
  setupApplicationTest(hooks);

  test('does not allow password reset for non-userpass users', async function (assert) {
    await authPage.login();
    waitFor('[data-test-user-menu-trigger]');

    await click('[data-test-user-menu-trigger]');
    assert.dom('[data-test-user-menu-item="reset-password"]').doesNotExist();
  });

  test('allows password reset for userpass users', async function (assert) {
    await authPage.login();
    await runCmd([
      mountAuthCmd('userpass'),
      createPolicyCmd('userpass', resetPolicy),
      'write auth/userpass/users/reset-me password=password token_policies=userpass',
    ]);

    await authPage.loginUsername('reset-me', 'password');
    if (currentURL() !== '/vault/dashboard') {
      // prevent flakiness -- try again
      await click('[data-test-auth-submit]');
    }
    assert.strictEqual(currentURL(), '/vault/dashboard', 'logs in');
    await click('[data-test-user-menu-trigger]');
    await click('[data-test-user-menu-item="reset-password"]');

    assert.strictEqual(currentURL(), '/vault/access/reset-password', 'links to password reset');

    assert.dom('[data-test-title]').hasText('Reset password', 'page title');
    await fillIn('[data-test-textarea]', 'newpassword');
    await click('[data-test-reset-password-save]');

    await flashMessage.waitForFlash();
    assert.strictEqual(flashMessage.latestMessage, 'Successfully reset password', 'Shows success message');
    assert.dom('[data-test-textarea]').hasValue('', 'Resets input after save');

    // Reset test state
    await authPage.login();
    await runCmd(deleteAuthCmd('userpass'));
  });
});
