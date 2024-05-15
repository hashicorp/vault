/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { currentURL, click, fillIn, settled, waitFor } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { createPolicyCmd, deleteAuthCmd, mountAuthCmd, runCmd } from '../helpers/commands';
import { v4 as uuidv4 } from 'uuid';

const SUCCESS_MESSAGE = 'Successfully reset password';

module('Acceptance | reset password', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    this.userpass = `userpass-${uid}`;
    this.policy = `
    path "auth/${this.userpass}/users/reset-me/password" {
      capabilities = ["update", "create"]
    }
    `;
  });

  hooks.afterEach(async function () {
    await authPage.login();
    await runCmd([deleteAuthCmd(this.userpass), `delete sys/policies/acl/${this.userpass}`], false);
  });

  test('does not allow password reset for non-userpass users', async function (assert) {
    await authPage.login();
    await settled();

    await click('[data-test-user-menu-trigger]');
    assert.dom('[data-test-user-menu-item="reset-password"]').doesNotExist();
  });

  test('allows password reset for userpass users logged in via dropdown', async function (assert) {
    await authPage.login();
    await runCmd([
      mountAuthCmd('userpass', this.userpass),
      createPolicyCmd(this.userpass, this.policy),
      `write auth/${this.userpass}/users/reset-me password=password token_policies=${this.userpass}`,
    ]);
    await authPage.loginUsername('reset-me', 'password', this.userpass);

    await click('[data-test-user-menu-trigger]');
    await click('[data-test-user-menu-item="reset-password"]');

    assert.strictEqual(currentURL(), '/vault/access/reset-password', 'links to password reset');
    assert
      .dom('[data-test-current-user-banner]')
      .hasText(
        `You are updating the password for reset-me on the ${this.userpass} auth mount.`,
        'shows correct banner text'
      );

    assert.dom('[data-test-title]').hasText('Reset password', 'page title');
    await fillIn('[data-test-textarea]', 'newpassword');
    await click('[data-test-reset-password-save]');
    await waitFor('[data-test-flash-message]');
    assert.dom('[data-test-flash-message]').hasText(`Success ${SUCCESS_MESSAGE}`);
    assert.dom('[data-test-textarea]').hasValue('', 'Resets input after save');
  });

  test('allows password reset for userpass users logged in via tab', async function (assert) {
    await authPage.login();
    await runCmd([
      mountAuthCmd('userpass', this.userpass),
      `write sys/auth/${this.userpass}/tune listing_visibility="unauth"`,
      createPolicyCmd(this.userpass, this.policy),
      `write auth/${this.userpass}/users/reset-me password=password token_policies=${this.userpass}`,
    ]);
    await authPage.logout();
    await click(`[data-test-auth-method="${this.userpass}"] a`);
    await authPage.usernameInput('reset-me');
    await authPage.passwordInput('password').submit();

    await click('[data-test-user-menu-trigger]');
    await click('[data-test-user-menu-item="reset-password"]');

    assert.strictEqual(currentURL(), '/vault/access/reset-password', 'links to password reset');
    assert
      .dom('[data-test-current-user-banner]')
      .hasText(
        `You are updating the password for reset-me on the ${this.userpass} auth mount.`,
        'shows correct banner text'
      );

    assert.dom('[data-test-title]').hasText('Reset password', 'page title');
    await fillIn('[data-test-textarea]', 'newpassword');
    await click('[data-test-reset-password-save]');
    await waitFor('[data-test-flash-message]');
    assert.dom('[data-test-flash-message]').hasText(`Success ${SUCCESS_MESSAGE}`);
    assert.dom('[data-test-textarea]').hasValue('', 'Resets input after save');
  });
});
