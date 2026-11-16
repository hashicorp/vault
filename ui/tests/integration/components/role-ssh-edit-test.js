/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import SshRoleForm from 'vault/forms/ssh/role';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TTL_PICKER } from 'vault/tests/helpers/components/ttl-picker-selectors';

const SELECTORS = {
  submit: '[data-test-role-ssh-create]',
  ttlToggle: `${GENERAL.fieldByAttr('ttl')} ${TTL_PICKER.toggle}`,
  ttlValue: `${GENERAL.fieldByAttr('ttl')} [data-test-ttl-value]`,
  ttlUnit: `${GENERAL.fieldByAttr('ttl')} ${TTL_PICKER.ttlUnit}`,
};

module('Integration | Component | role-ssh-edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.lookup('service:router').transitionTo = sinon.stub();
    this.writeStub = sinon.stub(this.owner.lookup('service:api').secrets, 'sshWriteRole').resolves();
    this.form = new SshRoleForm(
      { backend: 'ssh', key_type: 'ca', not_before_duration: '30s', port: 22 },
      { isNew: true }
    );

    this.renderComponent = () =>
      render(hbs`<RoleSshEdit @form={{this.form}} @mode="create" @root={{hash text="ssh"}} />`);

    this.enterCaTtl = async (value) => {
      await click(GENERAL.button('Options'));
      await fillIn(GENERAL.inputByAttr('name'), 'my-role');
      await click(SELECTORS.ttlToggle);
      await fillIn(SELECTORS.ttlValue, value);
      await fillIn(SELECTORS.ttlUnit, 'd');
    };
  });

  test('it blocks a TTL longer than Vault supports', async function (assert) {
    assert.expect(7);

    await this.renderComponent();
    await this.enterCaTtl('9999999999');
    await click(SELECTORS.submit);

    assert
      .dom(GENERAL.validationErrorByAttr('ttl'))
      .hasText('TTL is too large. Maximum allowed value is 106,751 days.');
    assert.dom('[data-test-ttl-error]').doesNotExist('the picker does not repeat the message');
    assert.dom(GENERAL.invalidFormMessage).hasText('There is an error with this form.');
    assert.true(this.writeStub.notCalled, 'the role is not saved');

    await fillIn(SELECTORS.ttlValue, '106751');
    await click(SELECTORS.submit);
    assert.true(this.writeStub.calledOnce, 'the role saves once the TTL is valid');
    const [name, backend, payload] = this.writeStub.firstCall.args;
    assert.deepEqual([name, backend], ['my-role', 'ssh']);
    assert.strictEqual(payload.ttl, '9223286400s', 'the maximum whole days is sent');
  });

  test('it blocks a TTL too long to be stored as a number', async function (assert) {
    assert.expect(2);

    await this.renderComponent();
    await this.enterCaTtl('9'.repeat(400));
    await click(SELECTORS.submit);

    assert.dom(GENERAL.validationErrorByAttr('ttl')).exists('the field shows an error');
    assert.true(this.writeStub.notCalled, 'the role is not saved');
  });

  test('it does not send CA-only durations for an OTP role', async function (assert) {
    assert.expect(4);

    await this.renderComponent();
    await this.enterCaTtl('9999999999');
    // The TTL fields are hidden once the key type is OTP, so a leftover value must not block or be sent.
    await fillIn(GENERAL.inputByAttr('key_type'), 'otp');
    await fillIn(GENERAL.inputByAttr('default_user'), 'admin');
    await click(SELECTORS.submit);

    assert.true(this.writeStub.calledOnce, 'the role is saved');
    const [, , payload] = this.writeStub.firstCall.args;
    assert.strictEqual(payload.key_type, 'otp');
    assert.notOk('ttl' in payload, 'ttl is not sent');
    assert.notOk('not_before_duration' in payload, 'not_before_duration is not sent');
  });
});
