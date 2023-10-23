import { module, test } from 'qunit';
import { currentURL, click, fillIn, settled } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { createPolicyCmd, mountAuthCmd, runCmd } from '../helpers/commands';

const resetPolicy = `
    path "sys/control-group/authorize" {
    capabilities = ["update"]
  }

  path "sys/control-group/request" {
    capabilities = ["update"]
  }
`;
module('Acceptance | reset password', function (hooks) {
  setupApplicationTest(hooks);

  test('does not allow password reset for non-userpass users', async function (assert) {
    await authPage.login();
    await settled();

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

    await click('[data-test-user-menu-trigger]');
    await click('[data-test-user-menu-item="reset-password"]');

    assert.strictEqual(currentURL(), '/vault/access/reset-password', 'links to password reset');

    assert.dom('[data-test-title]').hasText('Reset Password', 'page title');
    await fillIn('[data-test-textarea]', 'newpassword');
    await click('[data-test-reset-password-save]');
    assert.dom('[data-test-reset-password-success]').exists('shows success');
  });
});
