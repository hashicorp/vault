import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { click } from '@ember/test-helpers';
import { TOTP_NOT_CONFIGURED } from 'vault/services/auth';
import { TOTP_NA_MSG, MFA_ERROR_MSG } from 'vault/components/mfa-error';
const UNAUTH = 'MFA authorization failed';

module('Integration | Component | mfa-error', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    const auth = this.owner.lookup('service:auth');
    auth.set('mfaErrors', [TOTP_NOT_CONFIGURED]);

    this.onClose = () => assert.ok(true, 'onClose event is triggered');

    await render(hbs`<MfaError @onClose={{this.onClose}}/>`);

    assert.dom('[data-test-empty-state-title]').hasText('TOTP not set up', 'Title renders for TOTP error');
    assert
      .dom('[data-test-empty-state-subText]')
      .hasText(TOTP_NOT_CONFIGURED, 'Error message renders for TOTP error');
    assert.dom('[data-test-empty-state-message]').hasText(TOTP_NA_MSG, 'Description renders for TOTP error');

    auth.set('mfaErrors', [UNAUTH]);
    await render(hbs`<MfaError @onClose={{this.onClose}}/>`);

    assert.dom('[data-test-empty-state-title]').hasText('Unauthorized', 'Title renders for mfa error');
    assert.dom('[data-test-empty-state-subText]').hasText(UNAUTH, 'Error message renders for mfa error');
    assert.dom('[data-test-empty-state-message]').hasText(MFA_ERROR_MSG, 'Description renders for mfa error');

    await click('[data-test-go-back]');

    assert.equal(auth.mfaErrors, null, 'mfaErrors unset in auth service');
  });
});
