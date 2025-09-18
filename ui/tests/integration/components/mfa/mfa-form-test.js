/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled, fillIn, click, waitUntil, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { _cancelTimers as cancelTimers, later } from '@ember/runloop';
import { TOTP_VALIDATION_ERROR } from 'vault/components/mfa/mfa-form';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { MFA_SELECTORS } from 'vault/tests/helpers/mfa/mfa-selectors';
import { QR_CODE_URL } from 'vault/mirage/handlers/mfa-login';

module('Integration | Component | mfa-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.onCancel = sinon.spy();
    this.onSuccess = sinon.spy();
    this.onError = sinon.spy();
    this.authService = this.owner.lookup('service:auth');
    // setup basic totp mfaRequirement
    // override in tests that require different scenarios
    this.totpConstraint = this.server.create('mfa-method', { type: 'totp' });
    const mfaRequirement = this.authService.parseMfaResponse({
      mfa_request_id: 'test-mfa-id',
      mfa_constraints: { test_mfa: { any: [this.totpConstraint] } },
    });

    this.renderComponent = () => {
      return render(hbs`
  <Mfa::MfaForm
    @authData={{this.mfaAuthData}}
    @clusterId="123456"
    @onCancel={{this.onCancel}}
    @onError={{this.onError}}
    @onSuccess={{this.onSuccess}}
  />`);
    };
    this.setMfaAuthData = (mfaRequirement) => {
      this.mfaAuthData = {
        mfaRequirement: mfaRequirement,
        authMethodType: 'userpass',
        authMountPath: 'userpass',
      };
    };
    this.setMfaAuthData(mfaRequirement);
  });

  test('it renders correct text for single passcode', async function (assert) {
    const totpConstraint = this.server.create('mfa-method', { type: 'totp' });
    this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
      mfa_request_id: 'test-mfa-id',
      mfa_constraints: { test_mfa_1: { any: [totpConstraint] } },
    });

    await this.renderComponent();
    assert
      .dom(MFA_SELECTORS.description)
      .hasText(
        'Multi-factor authentication is enabled for your account. Enter your authentication code to log in.'
      );
  });

  test('it renders correct text for multiple methods', async function (assert) {
    const oktaConstraint = this.server.create('mfa-method', { type: 'okta' });
    const duoConstraint = this.server.create('mfa-method', { type: 'duo' });
    this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
      mfa_request_id: 'test-mfa-id',
      mfa_constraints: { test_mfa_1: { any: [duoConstraint, oktaConstraint] } },
    });

    await this.renderComponent();
    assert
      .dom(MFA_SELECTORS.subheader)
      .hasText(
        'Multi-factor authentication is enabled for your account. Choose one of the following methods to continue:'
      );
  });

  test('it renders correct text for multiple constraints', async function (assert) {
    const oktaConstraint = this.server.create('mfa-method', { type: 'okta' });
    const duoConstraint = this.server.create('mfa-method', { type: 'duo' });
    this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
      mfa_request_id: 'test-mfa-id',
      mfa_constraints: { test_mfa_1: { any: [oktaConstraint] }, test_mfa_2: { any: [duoConstraint] } },
    });

    await this.renderComponent();
    assert
      .dom(MFA_SELECTORS.description)
      .hasText(
        'Multi-factor authentication is enabled for your account. Two methods are required for successful authentication.'
      );
  });

  test('it should render a submit button', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.button('Verify')).isNotDisabled('Button is not disabled by default');
  });

  test('it should render method selects and passcode inputs', async function (assert) {
    assert.expect(2);
    const duoConstraint = this.server.create('mfa-method', { type: 'duo', uses_passcode: true });
    const oktaConstraint = this.server.create('mfa-method', { type: 'okta' });
    const pingidConstraint = this.server.create('mfa-method', { type: 'pingid' });
    const mfaRequirement = this.authService.parseMfaResponse({
      mfa_request_id: 'test-mfa-id',
      mfa_constraints: {
        test_mfa_1: {
          any: [pingidConstraint, oktaConstraint],
        },
        test_mfa_2: {
          any: [duoConstraint],
        },
      },
    });
    this.mfaAuthData.mfaRequirement = mfaRequirement;

    this.server.post('/sys/mfa/validate', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      const payload = {
        mfa_request_id: 'test-mfa-id',
        mfa_payload: { [oktaConstraint.id]: [], [duoConstraint.id]: ['passcode=test-code'] },
      };
      assert.deepEqual(json, payload, 'Correct mfa payload passed to validate endpoint');
      return {};
    });

    this.owner.lookup('service:auth').reopen({
      // override to avoid authSuccess method since it expects an auth payload
      async totpValidate({ mfaRequirement }) {
        await this.clusterAdapter().mfaValidate(mfaRequirement);
        return 'test response';
      },
    });

    this.onSuccess = (resp) =>
      assert.strictEqual(resp, 'test response', 'Response is returned in onSuccess callback');

    await this.renderComponent();
    await fillIn(MFA_SELECTORS.select(0), oktaConstraint.id);
    await fillIn(MFA_SELECTORS.passcode(1), 'test-code');
    await click(GENERAL.button('Verify'));
  });

  test('it should validate mfa requirement', async function (assert) {
    assert.expect(5);
    this.server.post('/sys/mfa/validate', (schema, req) => {
      const json = JSON.parse(req.requestBody);
      const payload = {
        mfa_request_id: 'test-mfa-id',
        mfa_payload: { [this.totpConstraint.id]: ['test-code'] },
      };
      assert.deepEqual(json, payload, 'Correct mfa payload passed to validate endpoint');
      return {};
    });
    const expectedAuthData = {
      clusterId: '123456',
      authMethodType: 'userpass',
      authMountPath: 'userpass',
      mfaRequirement: {
        mfa_constraints: [
          {
            methods: [this.totpConstraint],
            passcode: 'test-code', // Added by the MfaForm
            selectedMethod: this.totpConstraint,
          },
        ],
        mfa_request_id: 'test-mfa-id',
      },
    };
    this.owner.lookup('service:auth').reopen({
      // override to avoid authSuccess method since it expects an auth payload
      async totpValidate(authData) {
        await waitUntil(() =>
          assert
            .dom(`${GENERAL.button('Verify')} ${GENERAL.icon('loading')}`)
            .exists('Loading icon shows on button')
        );
        assert.dom(GENERAL.button('Verify')).isDisabled('Button is disabled while loading');
        assert.deepEqual(authData, expectedAuthData, 'Mfa auth data passed to validate method');
        await this.clusterAdapter().mfaValidate(authData.mfaRequirement);
        return 'test response';
      },
    });

    this.onSuccess = (resp) =>
      assert.strictEqual(resp, 'test response', 'Response is returned in onSuccess callback');

    await this.renderComponent();

    await fillIn(MFA_SELECTORS.passcode(), 'test-code');
    await click(GENERAL.button('Verify'));
  });

  test('it should show countdown on passcode already used and rate limit errors', async function (assert) {
    const messages = {
      used: 'code already used; new code is available in 30 seconds',
      // note: the backend returns a duplicate "s" in "30s seconds" in the limit message below. we have intentionally left it as is to ensure our regex for parsing the delay time can handle it
      limit:
        'failed to satisfy enforcement userpass2-not-self-enroll. error: 2 errors occurred:\n\t* maximum TOTP validation attempts 2 exceeded the allowed attempts 1. Please try again in 30s seconds\n\t* login MFA validation failed for methodID: [1f260334-ee5f-6e47-8e86-57be05d457d2]\n\n',
    };
    const codes = ['used', 'limit'];
    for (const code of codes) {
      this.owner.lookup('service:auth').reopen({
        totpValidate() {
          throw new Error(messages[code]);
        },
      });

      await this.renderComponent();
      await fillIn(MFA_SELECTORS.passcode(), 'foo');
      await click(GENERAL.button('Verify'));

      await waitFor(MFA_SELECTORS.countdown);

      assert
        .dom(MFA_SELECTORS.countdown)
        .includesText('30', 'countdown renders with correct initial value from error response');
      assert.dom(GENERAL.button('Verify')).isDisabled();
      assert.dom(GENERAL.cancelButton).isDisabled();
      assert.dom(MFA_SELECTORS.passcode()).isDisabled('Input is disabled during countdown');
      assert.dom(GENERAL.inlineError).exists('Alert message renders');
    }
  });

  test('it defaults countdown to 30 seconds if error message does not indicate when user can try again ', async function (assert) {
    const msg = 'maximum TOTP validation attempts 4 exceeded the allowed attempts 3. Beep-boop.';
    this.owner.lookup('service:auth').reopen({
      totpValidate() {
        throw new Error(msg);
      },
    });
    await this.renderComponent();

    await fillIn(MFA_SELECTORS.passcode(), 'foo');
    await click(GENERAL.button('Verify'));

    await waitFor(MFA_SELECTORS.countdown);

    assert
      .dom(MFA_SELECTORS.countdown)
      .includesText('30', 'countdown renders with correct initial value from error response');
    assert.dom(GENERAL.button('Verify')).isDisabled('Button is disabled during countdown');
    assert.dom(MFA_SELECTORS.passcode()).isDisabled('Input is disabled during countdown');
    assert.dom(GENERAL.inlineError).exists('Alert message renders');
  });

  test('it should show error message for passcode invalid error', async function (assert) {
    this.owner.lookup('service:auth').reopen({
      totpValidate() {
        throw { errors: ['failed to validate'] };
      },
    });
    await this.renderComponent();
    await fillIn(MFA_SELECTORS.passcode(), 'test-code');
    later(() => cancelTimers(), 50);
    await settled();

    await click(GENERAL.button('Verify'));
    assert
      .dom(GENERAL.messageError)
      .includesText(TOTP_VALIDATION_ERROR, 'Generic error message renders for passcode validation error');
  });

  test('it should call onCancel callback', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.cancelButton);
    assert.true(this.onCancel.calledOnce, 'it fires onCancel callback');
  });

  module('self-enrollment', function (hooks) {
    hooks.beforeEach(function () {
      // Self-enrollment is an enterprise only feature
      this.version = this.owner.lookup('service:version');
      this.version.type = 'enterprise';
      this.server.post('/identity/mfa/method/totp/self-enroll', async () => {
        return {
          data: {
            barcode:
              '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',
            url: QR_CODE_URL,
          },
        };
      });
    });

    test('it makes request to self-enroll endpoint when self_enrollment_enabled is true', async function (assert) {
      assert.expect(3);
      const request_id = crypto.randomUUID();
      const totpConstraint = this.server.create('mfa-method', {
        type: 'totp',
        self_enrollment_enabled: true,
      });
      this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
        mfa_request_id: request_id,
        mfa_constraints: { test_mfa_1: { any: [totpConstraint] } },
      });

      this.server.post('/identity/mfa/method/totp/self-enroll', async (schema, req) => {
        const { mfa_method_id, mfa_request_id } = JSON.parse(req.requestBody);
        assert.true(true, 'Request made to /self-enroll');
        assert.strictEqual(mfa_request_id, request_id, 'payload has expected request id');
        assert.strictEqual(mfa_method_id, totpConstraint.id, 'payload has expected method id');
        return {
          data: {
            barcode:
              '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',
            url: QR_CODE_URL,
          },
        };
      });

      await this.renderComponent();
      await click(GENERAL.button('Continue'));
    });

    test('it renders correct text for single passcode', async function (assert) {
      const totpConstraint = this.server.create('mfa-method', {
        type: 'totp',
        self_enrollment_enabled: true,
      });
      this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
        mfa_request_id: 'test-mfa-id',
        mfa_constraints: { test_mfa_1: { any: [totpConstraint] } },
      });

      await this.renderComponent();
      assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText('Your organization has enforced MFA TOTP to protect your accounts. Set up to continue.');
      assert.dom(MFA_SELECTORS.subtitle).hasText('Scan the QR code to continue');
      assert
        .dom(MFA_SELECTORS.description)
        .hasText(
          'Scan the QR code with your authenticator app. If you currently do not have a device on hand, you can copy the MFA secret below and enter it manually.',
          'Correct description renders for single passcode'
        );
      assert.dom(GENERAL.button('Continue')).exists();
      assert.dom(GENERAL.cancelButton).exists();

      // Go on to next step
      await click(GENERAL.button('Continue'));
      assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText('Your organization has enforced MFA TOTP to protect your accounts. Set up to continue.');
      assert.dom(MFA_SELECTORS.subtitle).doesNotExist();
      assert
        .dom(MFA_SELECTORS.description)
        .hasText('To verify your device, enter the code generated from your authenticator.');
      assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
      assert.dom(GENERAL.button('Verify')).exists();
      assert.dom(GENERAL.cancelButton).exists();
    });

    test('it renders correct text for multiple methods (1 passcode 1 push)', async function (assert) {
      const oktaConstraint = this.server.create('mfa-method', { type: 'okta' });
      const totpConstraint = this.server.create('mfa-method', {
        type: 'totp',
        self_enrollment_enabled: true,
      });
      this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
        mfa_request_id: 'test-mfa-id',
        mfa_constraints: { test_mfa_1: { any: [totpConstraint, oktaConstraint] } },
      });
      await this.renderComponent();
      assert.dom(GENERAL.title).hasText('Verify your identity');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText(
          'Multi-factor authentication is enabled for your account. Choose one of the following methods to continue:'
        );
      assert.dom(MFA_SELECTORS.subtitle).doesNotExist();
      assert.dom(GENERAL.button('Verify')).doesNotExist();
      assert.dom(GENERAL.cancelButton).exists();

      // Select TOTP
      await click(GENERAL.button('Setup to verify with TOTP'));
      await waitFor(MFA_SELECTORS.qrCode);
      assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText('Your organization has enforced MFA TOTP to protect your accounts. Set up to continue.');
      assert.dom(MFA_SELECTORS.subtitle).hasText('Scan the QR code to continue');
      assert
        .dom(MFA_SELECTORS.description)
        .hasText(
          'Scan the QR code with your authenticator app. If you currently do not have a device on hand, you can copy the MFA secret below and enter it manually.',
          'Correct description renders for single passcode'
        );
      assert.dom(GENERAL.button('Continue')).exists();
      assert.dom(GENERAL.cancelButton).exists();

      // Go on to next step
      await click(GENERAL.button('Continue'));
      assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText('Your organization has enforced MFA TOTP to protect your accounts. Set up to continue.');
      assert.dom(MFA_SELECTORS.subtitle).doesNotExist();
      assert
        .dom(MFA_SELECTORS.description)
        .hasText('To verify your device, enter the code generated from your authenticator.');
      assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
      assert.dom(GENERAL.button('Verify')).exists();
      assert.dom(GENERAL.cancelButton).exists('it renders "Cancel" after self-enroll workflow');
    });

    test('it renders correct text for multiple methods (2 passcodes)', async function (assert) {
      const duoConstraint = this.server.create('mfa-method', { type: 'duo', uses_passcode: true });
      const totpConstraint = this.server.create('mfa-method', {
        type: 'totp',
        self_enrollment_enabled: true,
      });
      this.mfaAuthData.mfaRequirement = this.authService.parseMfaResponse({
        mfa_request_id: 'test-mfa-id',
        mfa_constraints: { test_mfa_1: { any: [totpConstraint, duoConstraint] } },
      });
      await this.renderComponent();
      assert.dom(GENERAL.title).hasText('Verify your identity');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText(
          'Multi-factor authentication is enabled for your account. Choose one of the following methods to continue:'
        );
      assert.dom(MFA_SELECTORS.subtitle).doesNotExist();
      assert.dom(GENERAL.button('Verify')).doesNotExist();
      assert.dom(GENERAL.cancelButton).exists();

      // Select TOTP
      await click(GENERAL.button('Setup to verify with TOTP'));
      await waitFor(MFA_SELECTORS.qrCode);
      assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText('Your organization has enforced MFA TOTP to protect your accounts. Set up to continue.');
      assert.dom(MFA_SELECTORS.subtitle).hasText('Scan the QR code to continue');
      assert
        .dom(MFA_SELECTORS.description)
        .hasText(
          'Scan the QR code with your authenticator app. If you currently do not have a device on hand, you can copy the MFA secret below and enter it manually.',
          'Correct description renders for single passcode'
        );
      assert.dom(GENERAL.button('Continue')).exists();
      assert.dom(GENERAL.cancelButton).exists();

      // Go on to next step
      await click(GENERAL.button('Continue'));
      assert.dom(GENERAL.title).hasText('Set up MFA TOTP to continue');
      assert
        .dom(MFA_SELECTORS.subheader)
        .hasText('Your organization has enforced MFA TOTP to protect your accounts. Set up to continue.');
      assert.dom(MFA_SELECTORS.subtitle).doesNotExist();
      assert
        .dom(MFA_SELECTORS.description)
        .hasText('To verify your device, enter the code generated from your authenticator.');
      assert.dom(MFA_SELECTORS.label).hasText('Enter your one-time code');
      assert.dom(GENERAL.button('Verify')).exists();
      assert.dom(GENERAL.cancelButton).exists('it renders "Cancel" after self-enroll workflow');
    });

    test('it should render qr code and copy button', async function (assert) {
      const clipboardSpy = sinon.stub(navigator.clipboard, 'writeText').resolves();
      const totpConstraint = this.server.create('mfa-method', {
        type: 'totp',
        self_enrollment_enabled: true,
      });
      const mfaRequirement = this.authService.parseMfaResponse({
        mfa_request_id: 'test-mfa-id',
        mfa_constraints: { test_mfa: { any: [totpConstraint] } },
      });
      this.setMfaAuthData(mfaRequirement);
      await this.renderComponent();
      await waitFor(MFA_SELECTORS.qrCode);
      assert
        .dom(MFA_SELECTORS.mfaForm)
        .hasText(
          'Set up MFA TOTP to continue Your organization has enforced MFA TOTP to protect your accounts. Set up to continue. Scan the QR code to continue Scan the QR code with your authenticator app. If you currently do not have a device on hand, you can copy the MFA secret below and enter it manually. Or Copy TOTP setup URL For your security, this code is only shown once. Please scan or copy the setup URL into your authenticator app now. Continue Cancel',
          'it renders self-enrollment text'
        );
      assert.dom(MFA_SELECTORS.qrCode).exists('it renders qr code');
      assert.dom(GENERAL.cancelButton).exists();
      assert.dom(MFA_SELECTORS.verifyForm).doesNotExist('it does not render input field for TOTP code');
      assert.dom(GENERAL.button('Verify')).doesNotExist('it does not render Validate button');
      await click(GENERAL.copyButton);
      assert.strictEqual(clipboardSpy.firstCall.args[0], QR_CODE_URL, 'copy value is qr code URL');
      // Restore original clipboard
      clipboardSpy.restore(); // cleanup
    });
  });
});
