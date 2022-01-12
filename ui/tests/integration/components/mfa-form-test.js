import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { fillIn, click, waitUntil } from '@ember/test-helpers';
import { run, later } from '@ember/runloop';

module('Integration | Component | mfa-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.clusterId = '123456';
    this.mfaAuthData = {
      backend: 'userpass',
      data: { username: 'foo', password: 'bar' },
      mfa_enforcement: {
        mfa_request_id: 'test-mfa-id',
        mfa_constraints: [this.server.create('mfa-method', { type: 'totp' })],
      },
    };
  });

  test('it should validate passcode', async function (assert) {
    const validate = async (authData, passcode) => {
      await waitUntil(() =>
        assert.dom('[data-test-mfa-validate]').hasClass('is-loading', 'Loading class applied to button')
      );
      assert.dom('[data-test-mfa-validate]').isDisabled('Button is disabled while loading');
      assert.deepEqual(
        authData,
        { clusterId: this.clusterId, ...this.mfaAuthData },
        'Mfa auth data passed to validate method'
      );
      assert.equal(passcode, 'test-code', 'Passcode passed to validate method');
    };
    this.owner.lookup('service:auth').reopen({
      async totpValidate(authData, passcode) {
        await validate(authData, passcode);
        return 'test response';
      },
    });

    this.onSuccess = (resp) =>
      assert.equal(resp, 'test response', 'Response is returned in onSuccess callback');

    await render(hbs`
      <MfaForm
        @clusterId={{this.clusterId}}
        @authData={{this.mfaAuthData}}
        @onSuccess={{this.onSuccess}}
      />
    `);

    await fillIn('[data-test-mfa-passcode]', 'test-code');
    await click('[data-test-mfa-validate]');
  });

  test('it should show countdown on passcode validation failure', async function (assert) {
    this.owner.lookup('service:auth').reopen({
      totpValidate() {
        throw new Error('Incorrect passcode');
      },
    });
    await render(hbs`
      <MfaForm
        @clusterId={{this.clusterId}}
        @authData={{this.mfaAuthData}}
      />
    `);

    await fillIn('[data-test-mfa-passcode]', 'test-code');
    later(() => run.cancelTimers(), 50);
    await click('[data-test-mfa-validate]');
    assert.dom('[data-test-mfa-validate]').isDisabled('Button is disabled during countdown');
    assert.dom('[data-test-mfa-passcode]').isDisabled('Input is disabled during countdown');
    assert.dom('[data-test-mfa-passcode]').hasNoValue('Input value is cleared on error');
    assert.dom('[data-test-inline-error-message]').exists('Alert message renders');
    assert.dom('[data-test-mfa-countdown]').exists('30 second countdown renders');
  });
});
