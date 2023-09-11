/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentRouteName, fillIn, visit, waitUntil, find } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { validationHandler } from '../../mirage/handlers/mfa-login';

module('Acceptance | mfa-login', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'mfaLogin';
  });
  hooks.beforeEach(function () {
    this.auth = this.owner.lookup('service:auth');
    this.select = async (select = 0, option = 1) => {
      const selector = `[data-test-mfa-select="${select}"]`;
      const value = this.element.querySelector(`${selector} option:nth-child(${option + 1})`).value;
      await fillIn(`${selector} select`, value);
    };
    return visit('/vault/logout');
  });
  hooks.afterEach(function () {
    // Manually clear token after each so that future tests don't get into a weird state
    this.auth.deleteCurrentToken();
  });
  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  const login = async (user) => {
    await visit('/vault/auth');
    await fillIn('[data-test-select="auth-method"]', 'userpass');
    await fillIn('[data-test-username]', user);
    await fillIn('[data-test-password]', 'test');
    await click('[data-test-auth-submit]');
  };
  const didLogin = (assert) => {
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'Route transitions after login');
  };
  const validate = async (multi) => {
    await fillIn('[data-test-mfa-passcode="0"]', 'test');
    if (multi) {
      await fillIn('[data-test-mfa-passcode="1"]', 'test');
    }
    await click('[data-test-mfa-validate]');
  };

  test('it should handle single mfa constraint with passcode method', async function (assert) {
    assert.expect(4);
    await login('mfa-a');
    assert
      .dom('[data-test-mfa-description]')
      .includesText(
        'Enter your authentication code to log in.',
        'Mfa form displays with correct description'
      );
    assert.dom('[data-test-mfa-select]').doesNotExist('Select is hidden for single method');
    assert.dom('[data-test-mfa-passcode]').exists({ count: 1 }, 'Single passcode input renders');
    await validate();
    didLogin(assert);
  });

  test('it should handle single mfa constraint with push method', async function (assert) {
    assert.expect(6);

    server.post('/sys/mfa/validate', async (schema, req) => {
      await waitUntil(() => find('[data-test-mfa-description]'));
      assert
        .dom('[data-test-mfa-description]')
        .hasText(
          'Multi-factor authentication is enabled for your account.',
          'Mfa form displays with correct description'
        );
      assert.dom('[data-test-mfa-label]').hasText('Okta push notification', 'Correct method renders');
      assert
        .dom('[data-test-mfa-push-instruction]')
        .hasText('Check device for push notification', 'Push notification instruction renders');
      assert.dom('[data-test-mfa-validate]').isDisabled('Button is disabled while validating');
      assert
        .dom('[data-test-mfa-validate]')
        .hasClass('is-loading', 'Loading class applied to button while validating');
      return validationHandler(schema, req);
    });

    await login('mfa-b');
    didLogin(assert);
  });

  test('it should handle single mfa constraint with 2 passcode methods', async function (assert) {
    assert.expect(4);
    await login('mfa-c');
    assert
      .dom('[data-test-mfa-description]')
      .includesText('Select the MFA method you wish to use.', 'Mfa form displays with correct description');
    assert
      .dom('[data-test-mfa-select]')
      .exists({ count: 1 }, 'Select renders for single constraint with multiple methods');
    assert.dom('[data-test-mfa-passcode]').doesNotExist('Passcode input hidden until selection is made');
    await this.select();
    await validate();
    didLogin(assert);
  });

  test('it should handle single mfa constraint with 2 push methods', async function (assert) {
    assert.expect(1);
    await login('mfa-d');
    await this.select();
    await click('[data-test-mfa-validate]');
    didLogin(assert);
  });

  test('it should handle single mfa constraint with 1 passcode and 1 push method', async function (assert) {
    assert.expect(3);
    await login('mfa-e');
    await this.select(0, 2);
    assert.dom('[data-test-mfa-passcode]').exists('Passcode input renders');
    await this.select();
    assert.dom('[data-test-mfa-passcode]').doesNotExist('Passcode input is hidden for push method');
    await click('[data-test-mfa-validate]');
    didLogin(assert);
  });

  test('it should handle multiple mfa constraints with 1 passcode method each', async function (assert) {
    assert.expect(3);
    await login('mfa-f');
    assert
      .dom('[data-test-mfa-description]')
      .includesText(
        'Two methods are required for successful authentication.',
        'Mfa form displays with correct description'
      );
    assert.dom('[data-test-mfa-select]').doesNotExist('Selects do not render for single methods');
    await validate(true);
    didLogin(assert);
  });

  test('it should handle multi mfa constraint with 1 push method each', async function (assert) {
    assert.expect(1);
    await login('mfa-g');
    didLogin(assert);
  });

  test('it should handle multiple mfa constraints with 1 passcode and 1 push method', async function (assert) {
    assert.expect(4);
    await login('mfa-h');
    assert
      .dom('[data-test-mfa-description]')
      .includesText(
        'Two methods are required for successful authentication.',
        'Mfa form displays with correct description'
      );
    assert.dom('[data-test-mfa-select]').doesNotExist('Select is hidden for single method');
    assert.dom('[data-test-mfa-passcode]').exists({ count: 1 }, 'Passcode input renders');
    await validate();
    didLogin(assert);
  });

  test('it should handle multiple mfa constraints with multiple mixed methods', async function (assert) {
    assert.expect(2);
    await login('mfa-i');
    assert
      .dom('[data-test-mfa-description]')
      .includesText(
        'Two methods are required for successful authentication.',
        'Mfa form displays with correct description'
      );
    await this.select();
    await fillIn('[data-test-mfa-passcode="1"]', 'test');
    await click('[data-test-mfa-validate]');
    didLogin(assert);
  });

  test('it should render unauthorized message for push failure', async function (assert) {
    await login('mfa-j');
    assert.dom('[data-test-auth-form]').doesNotExist('Auth form hidden when mfa fails');
    assert.dom('[data-test-empty-state-title]').hasText('Unauthorized', 'Error title renders');
    assert
      .dom('[data-test-empty-state-subText]')
      .hasText('PingId MFA validation failed', 'Error message from server renders');
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'Multi-factor authentication is required, but failed. Go back and try again, or contact your administrator.',
        'Error description renders'
      );
    await click('[data-test-mfa-error] button');
    assert.dom('[data-test-auth-form]').exists('Auth form renders after mfa error dismissal');
  });
});
