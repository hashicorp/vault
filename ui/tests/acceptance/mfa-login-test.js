/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { click, currentRouteName, fillIn, visit, waitUntil, find, waitFor } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import mfaLoginHandler, { validationHandler } from 'vault/mirage/handlers/mfa-login';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';
import { MFA_SELECTORS } from 'vault/tests/helpers/mfa/mfa-selectors';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Acceptance | mfa-login', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    mfaLoginHandler(this.server);
    this.auth = this.owner.lookup('service:auth');
    this.select = async (select = 0, option = 1) => {
      const selector = MFA_SELECTORS.select(select);
      const value = this.element.querySelector(`${selector} option:nth-child(${option + 1})`).value;
      await fillIn(selector, value);
    };
    return visit('/vault/logout');
  });

  hooks.afterEach(function () {
    // Manually clear token after each so that future tests don't get into a weird state
    this.auth.deleteCurrentToken();
  });

  const login = async (user) => {
    await visit('/vault/auth');
    await fillIn(AUTH_FORM.selectMethod, 'userpass');
    await fillIn(GENERAL.inputByAttr('username'), user);
    await fillIn(GENERAL.inputByAttr('password'), 'test');
    await click(GENERAL.submitButton);
  };
  const didLogin = async (assert) => {
    await waitFor('[data-test-dashboard-card-header]', {
      timeout: 5000,
      timeoutMessage: 'timed out waiting for dashboard title to render',
    });
    assert.strictEqual(currentRouteName(), 'vault.cluster.dashboard', 'Route transitions after login');
  };
  const validate = async (multi) => {
    await fillIn(MFA_SELECTORS.passcode(0), 'test');
    if (multi) {
      await fillIn(MFA_SELECTORS.passcode(1), 'test');
    }
    await click(GENERAL.button('Verify'));
  };

  const assertSelfEnroll = async (assert) => {
    await waitFor(MFA_SELECTORS.qrCode);
    assert.dom(MFA_SELECTORS.qrCode).exists('it renders a QR code');
    await click(GENERAL.button('Continue'));
    assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, '1 passcode inputs renders');
  };

  test('it should handle single constraint with passcode method', async function (assert) {
    assert.expect(5);
    await login('mfa-a');
    assert.dom(GENERAL.title).hasText('Sign in to Vault');
    assert
      .dom(MFA_SELECTORS.description)
      .includesText(
        'Enter your authentication code to log in.',
        'Mfa form displays with correct description'
      );
    assert.dom(MFA_SELECTORS.select()).doesNotExist('Select is hidden for single method');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, 'Single passcode input renders');
    await validate();
    await didLogin(assert);
  });

  test('it should handle single constraint with push method', async function (assert) {
    assert.expect(6);

    server.post('/sys/mfa/validate', async (schema, req) => {
      await waitUntil(() => find(MFA_SELECTORS.description));
      assert
        .dom(MFA_SELECTORS.description)
        .hasText(
          'Multi-factor authentication is enabled for your account.',
          'Mfa form displays with correct description'
        );
      assert.dom(MFA_SELECTORS.label).hasText('Okta push notification', 'Correct method renders');
      assert
        .dom(MFA_SELECTORS.push)
        .hasText('Check device for push notification', 'Push notification instruction renders');
      assert.dom(GENERAL.button('Verify')).isDisabled('Button is disabled while validating');
      assert
        .dom(`${GENERAL.button('Verify')} ${GENERAL.icon('loading')}`)
        .exists('Loading icon shows while validating');
      return validationHandler(schema, req);
    });

    await login('mfa-b');
    await didLogin(assert);
  });

  test('it should handle single constraint with 2 passcode methods', async function (assert) {
    assert.expect(6);
    await login('mfa-c');
    assert.dom(GENERAL.title).hasText('Verify your identity');
    assert
      .dom(MFA_SELECTORS.subheader)
      .hasText(
        'Multi-factor authentication is enabled for your account. Choose one of the following methods to continue:',
        'Mfa form displays with correct description'
      );
    assert.dom(GENERAL.button('Verify with Duo')).exists('It renders button for Duo');
    assert.dom(GENERAL.button('Verify with TOTP')).exists('It renders button for TOTP');
    assert.dom(MFA_SELECTORS.passcode()).doesNotExist('Passcode input hidden until selection is made');
    await click(GENERAL.button('Verify with TOTP'));
    await validate();
    await didLogin(assert);
  });

  test('it should handle single constraint with 2 push methods', async function (assert) {
    assert.expect(3);
    await login('mfa-d');
    assert.dom(GENERAL.button('Verify with Okta')).exists('It renders button for Okta');
    assert.dom(GENERAL.button('Verify with Duo')).exists('It renders button for Duo');
    await click(GENERAL.button('Verify with Okta'));
    await didLogin(assert);
  });

  test('it should handle single constraint with 1 passcode and 1 push method', async function (assert) {
    assert.expect(3);
    await login('mfa-e');
    assert.dom(GENERAL.button('Verify with Okta')).exists('It renders button for Okta');
    await click(GENERAL.button('Verify with TOTP'));
    assert.dom(MFA_SELECTORS.passcode()).exists('Passcode input renders');
    await click(GENERAL.button('Try another method'));
    // Clicking "Verify with Okta" automatically starts validation so no need to click "Verify"
    await click(GENERAL.button('Verify with Okta'));
    await didLogin(assert);
  });

  test('it should handle multiple constraints with 1 passcode method each', async function (assert) {
    assert.expect(3);
    await login('mfa-f');
    assert
      .dom(MFA_SELECTORS.description)
      .includesText(
        'Two methods are required for successful authentication.',
        'Mfa form displays with correct description'
      );
    assert.dom(MFA_SELECTORS.select()).doesNotExist('Selects do not render for single methods');
    await validate(true);
    await didLogin(assert);
  });

  test('it should handle multi mfa constraint with 1 push method each', async function (assert) {
    assert.expect(1);
    await login('mfa-g');
    await didLogin(assert);
  });

  test('it should handle multiple constraints with 1 passcode and 1 push method', async function (assert) {
    assert.expect(4);
    await login('mfa-h');
    assert
      .dom(MFA_SELECTORS.description)
      .includesText(
        'Two methods are required for successful authentication.',
        'Mfa form displays with correct description'
      );
    assert.dom(MFA_SELECTORS.select()).doesNotExist('Select is hidden for single method');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, 'Passcode input renders');
    await validate();
    await didLogin(assert);
  });

  test('it should handle multiple constraints with multiple mixed methods', async function (assert) {
    assert.expect(2);
    await login('mfa-i');
    assert
      .dom(MFA_SELECTORS.description)
      .includesText(
        'Two methods are required for successful authentication.',
        'Mfa form displays with correct description'
      );
    await this.select();
    await fillIn(MFA_SELECTORS.passcode(1), 'test');
    await click(GENERAL.button('Verify'));
    await didLogin(assert);
  });

  test('it should render unauthorized message for push failure', async function (assert) {
    await login('mfa-j');
    await waitFor(GENERAL.messageError);
    assert.dom(AUTH_FORM.form).doesNotExist('Auth form does not render');
    // Using hasTextContaining because the UUID is regenerated each time in mirage
    assert
      .dom(GENERAL.messageError)
      .hasTextContaining(
        'Error failed to satisfy enforcement test_0 pingid authentication failed: "Login request denied." login MFA validation failed for methodID:'
      );
    await click(GENERAL.cancelButton);
    assert.dom(AUTH_FORM.form).exists('Auth form renders after mfa error dismissal');
  });

  /*
   * SELF-ENROLLMENT TESTS
   * Even though self-enrollment is an enterprise-only feature, these tests use Mirage so we don't need to filter them out of CE test runs
   */
  test('self-enroll: single constraint with one TOTP passcode', async function (assert) {
    await login('mfa-a-self');
    await waitFor(MFA_SELECTORS.qrCode);
    assert.dom(MFA_SELECTORS.qrCode).exists('it renders QR code');
    assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
    await click(GENERAL.button('Continue'));
    assert.dom(GENERAL.button('Continue')).doesNotExist('"Continue" button is replaced by "Verify"');
    assert.dom(MFA_SELECTORS.qrCode).doesNotExist('Clicking "Continue" removes QR code');
    assert.dom(GENERAL.button('Verify')).exists();
    assert.dom(MFA_SELECTORS.qrCode).doesNotExist('Clicking "Continue" removes QR code');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, 'Single passcode input renders');
    await validate();
    await didLogin(assert);
  });

  test('self-enroll: single constraint with 2 passcode methods', async function (assert) {
    await login('mfa-c-self');
    // Buttons render for both Duo and TOTP, we want to click TOTP to initiate self-enrollment flow
    assert.dom(GENERAL.title).hasText('Verify your identity');
    assert.dom(GENERAL.button('Verify with Duo')).exists();
    await click(GENERAL.button('Setup to verify with TOTP'));
    await waitFor(MFA_SELECTORS.qrCode);
    assert.dom(MFA_SELECTORS.qrCode).exists('it renders QR code');
    assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
    // Click "Continue" for second setup step to verify passcode
    await click(GENERAL.button('Continue'));
    assert
      .dom(MFA_SELECTORS.description)
      .hasText('To verify your device, enter the code generated from your authenticator.');
    assert.dom(GENERAL.button('Continue')).doesNotExist('"Continue" button is replaced by "Verify"');
    assert.dom(MFA_SELECTORS.qrCode).doesNotExist('Clicking "Continue" removes QR code');
    assert.dom(GENERAL.button('Verify')).exists();
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, 'Passcode input renders');
    await validate();
    await didLogin(assert);
  });

  test('self-enroll: multiple constraints with 1 passcode method each', async function (assert) {
    await login('mfa-f-self');
    await waitFor(MFA_SELECTORS.qrCode);
    assert.dom(MFA_SELECTORS.qrCode).exists('it renders QR code');
    assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
    // Click "Continue" for second setup step to verify passcode
    await click(GENERAL.button('Continue'));
    assert.dom(GENERAL.button('Continue')).doesNotExist('"Continue" button is replaced by "Verify"');
    assert.dom(MFA_SELECTORS.qrCode).doesNotExist('Clicking "Continue" removes QR code');
    assert
      .dom(MFA_SELECTORS.description)
      .hasText('To verify your device, enter the code generated from your authenticator.');
    assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, '1 passcode inputs renders');
    // Fill in and click "Verify" which should render passcode for second constraint
    await fillIn(MFA_SELECTORS.passcode(0), 'test');
    await click(GENERAL.button('Verify'));
    assert
      .dom(MFA_SELECTORS.description)
      .hasText(
        'Multi-factor authentication is enabled for your account. Two methods are required for successful authentication.'
      );
    assert.dom(MFA_SELECTORS.verifyBadge('TOTP passcode')).hasText('TOTP passcode');
    assert.dom(GENERAL.button('Verify')).exists();
    assert.dom(MFA_SELECTORS.label).hasText('Duo passcode');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, '1 passcode inputs renders');
    assert.dom('hr').exists({ count: 1 }, 'only one separator renders');
    await fillIn(MFA_SELECTORS.passcode(1), 'test');
    await click(GENERAL.button('Verify'));
    await didLogin(assert);
  });

  test('self-enroll: multiple constraints with 1 passcode and 1 push method', async function (assert) {
    await login('mfa-h-self');
    await waitFor(MFA_SELECTORS.qrCode);
    assert.dom(MFA_SELECTORS.qrCode).exists('it renders QR code');
    assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
    // Click "Continue" for second setup step to verify passcode
    await click(GENERAL.button('Continue'));
    assert.dom(GENERAL.button('Continue')).doesNotExist('"Continue" button is replaced by "Verify"');
    assert.dom(MFA_SELECTORS.qrCode).doesNotExist('Clicking "Continue" removes QR code');
    assert
      .dom(MFA_SELECTORS.description)
      .hasText('To verify your device, enter the code generated from your authenticator.');
    assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, '1 passcode inputs renders');
    // Fill in and click "Verify" which should immediately trigger MFA validation because the
    // second constraint is a push notification and no user input is required.
    await fillIn(MFA_SELECTORS.passcode(0), 'test');
    await click(GENERAL.button('Verify'));
    await didLogin(assert);
  });

  test('self-enroll: multiple constraints with multiple mixed methods', async function (assert) {
    await login('mfa-i-self');
    assert.dom(GENERAL.title).hasText('Verify your identity');
    assert
      .dom(MFA_SELECTORS.subheader)
      .hasText(
        'Multi-factor authentication is enabled for your account. Choose one of the following methods to continue:'
      );
    assert.dom(GENERAL.button('Verify with Okta')).exists();
    await click(GENERAL.button('Setup to verify with TOTP'));
    await waitFor(MFA_SELECTORS.qrCode);
    assert.dom(MFA_SELECTORS.qrCode).exists('it renders QR code');
    assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
    // Click "Continue" to validate TOTP
    await click(GENERAL.button('Continue'));
    assert.dom(GENERAL.button('Continue')).doesNotExist('"Continue" button is replaced by "Verify"');
    assert.dom(MFA_SELECTORS.qrCode).doesNotExist('Clicking "Continue" removes QR code');
    assert
      .dom(MFA_SELECTORS.description)
      .hasText('To verify your device, enter the code generated from your authenticator.');
    assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, '1 passcode inputs renders');
    // Fill in and click "Verify" which should render passcode for second constraint
    await fillIn(MFA_SELECTORS.passcode(0), 'test');
    await click(GENERAL.button('Verify'));
    assert.dom(MFA_SELECTORS.verifyBadge('TOTP passcode')).hasText('TOTP passcode');
    assert.dom(MFA_SELECTORS.select(0)).isDisabled('Select with self-enrolled TOTP is disabled');
    assert.dom(MFA_SELECTORS.passcode()).exists({ count: 1 }, '1 passcode inputs renders');
    assert.dom(MFA_SELECTORS.label).hasText('Duo passcode');
    await fillIn(MFA_SELECTORS.passcode(1), 'test');
    await click(GENERAL.button('Verify'));
    await didLogin(assert);
  });

  test('self-enroll: multiple constraints, 1 with 2 methods (one that supports self-enroll), 1 with push method', async function (assert) {
    await login('mfa-z-self');
    // For the constraint that supports self-enrollment, user must select it first.
    assert.dom(MFA_SELECTORS.select(0)).exists();
    assert.dom(MFA_SELECTORS.select(1)).exists();
    assert.dom('hr').exists({ count: 1 }, 'only one separator renders');
    await this.select(0, 1);
    await assertSelfEnroll(assert);
    await fillIn(MFA_SELECTORS.passcode(0), 'test');
    await click(GENERAL.button('Verify'));
    assert
      .dom(`${MFA_SELECTORS.select(0)} option:nth-child(2)`)
      .hasText('TOTP passcode (supports self-enrollment)', 'TOTP is pre-selected for the first constraint');
    await this.select(1, 2);
    // On second selection we are redirected
    assert.dom(GENERAL.title).hasText('Sign in to Vault');
    assert.dom(MFA_SELECTORS.verifyForm).exists('it renders mfa validation form');
    assert.dom(MFA_SELECTORS.select(0)).isDisabled();
    assert.dom(MFA_SELECTORS.verifyBadge('TOTP passcode')).exists('pending verification badge exists');
    assert.dom(MFA_SELECTORS.select(1)).isNotDisabled();
    await didLogin(assert);
  });

  test('self-enroll: if validation fails it resets the enrollment status', async function (assert) {
    await login('mfa-c-self');
    await click(GENERAL.button('Setup to verify with TOTP'));
    await waitFor(MFA_SELECTORS.qrCode);
    await click(GENERAL.button('Continue'));
    await fillIn(MFA_SELECTORS.passcode(0), '123456'); // not a valid code so it fails
    await click(GENERAL.button('Verify'));
    await waitFor(MFA_SELECTORS.verifyForm);
    assert.dom(MFA_SELECTORS.verifyBadge('TOTP passcode')).doesNotExist();
    assert.dom(MFA_SELECTORS.passcode(0)).exists().hasValue('123456', 'input has last used passcode');
  });

  module('error handling', function (hooks) {
    hooks.beforeEach(function () {
      // TODO confirm with backend what errors could be returned
      this.server.post('/identity/mfa/method/totp/self-enroll', async () => {
        return overrideResponse(500, JSON.stringify({ errors: ['uh oh!'] }));
      });
    });

    test('single constraint with one TOTP passcode', async function (assert) {
      await login('mfa-a-self');
      await waitFor(MFA_SELECTORS.verifyForm);
      assert.dom(MFA_SELECTORS.qrCode).doesNotExist('it does not enter self-enroll workflow');
      assert.dom(GENERAL.button('Continue')).doesNotExist();
      assert.dom(GENERAL.messageError).hasText('Error uh oh!', 'it renders error messages');
    });

    test('single constraint with 2 passcode methods', async function (assert) {
      await login('mfa-c-self');
      // Buttons render for both Duo and TOTP, we want to click TOTP to initiate self-enrollment flow
      assert.dom(GENERAL.button('Verify with Duo')).exists();
      await click(GENERAL.button('Setup to verify with TOTP'));
      await waitFor(MFA_SELECTORS.verifyForm);
      assert.dom(MFA_SELECTORS.qrCode).doesNotExist('it does not enter self-enroll workflow');
      assert.dom(GENERAL.button('Continue')).doesNotExist();
      assert.dom(GENERAL.messageError).hasText('Error uh oh!', 'it renders error messages');
    });

    test('multiple constraints with 1 passcode method each', async function (assert) {
      await login('mfa-f-self');
      await waitFor(MFA_SELECTORS.verifyForm);
      assert.dom(MFA_SELECTORS.qrCode).doesNotExist('it does not enter self-enroll workflow');
      assert.dom(GENERAL.button('Continue')).doesNotExist();
      assert.dom(GENERAL.messageError).hasText('Error uh oh!', 'it renders error messages');
    });
  });
});
