/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import PkiRoleForm from 'vault/forms/secrets/pki/role';

module('Integration | Component | pki key parameters', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.data = { key_type: 'rsa', key_bits: 2048, signature_bits: 0 };
    this.form = new PkiRoleForm(this.data, { isNew: true });
    this.fields = this.form.formFieldGroups.find((g) => g['Key parameters'])['Key parameters'];
    this.renderComponent = () =>
      render(
        hbs`<PkiKeyParameters @form={{this.form}} @fields={{this.fields}} @modelValidations={{this.modelValidations}} />`,
        { owner: this.engine }
      );
  });

  test('it should render the component and display the correct values', async function (assert) {
    assert.expect(3);

    await this.renderComponent();
    assert.dom(GENERAL.inputByAttr('key_type')).hasValue('rsa');
    assert.dom(GENERAL.inputByAttr('key_bits')).hasValue('2048');
    assert.dom(GENERAL.inputByAttr('signature_bits')).hasValue('0');
  });

  test('it should set values of key_type and key_bits when key_type changes', async function (assert) {
    assert.expect(8);

    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('key_type'), 'ec');
    assert.strictEqual(
      this.form.data.key_type,
      'ec',
      'sets the new selected value for key_type on the model.'
    );
    assert.strictEqual(
      this.form.data.key_bits,
      256,
      'sets the new selected value for key_bits on the model based on the selection of key_type.'
    );

    await fillIn(GENERAL.inputByAttr('key_type'), 'ed25519');
    assert.strictEqual(
      this.form.data.key_type,
      'ed25519',
      'sets the new selected value for key_type on the model.'
    );
    assert.strictEqual(
      this.form.data.key_bits,
      0,
      'sets the new selected value for key_bits on the model based on the selection of key_type.'
    );

    await fillIn(GENERAL.inputByAttr('key_type'), 'ec');
    await fillIn(GENERAL.inputByAttr('key_bits'), '384');
    assert.strictEqual(
      this.form.data.key_type,
      'ec',
      'sets the new selected value for key_type on the model.'
    );
    assert.strictEqual(
      this.form.data.key_bits,
      '384',
      'sets the new selected value for key_bits on the model based on the selection of key_type.'
    );

    await fillIn(GENERAL.inputByAttr('signature_bits'), '384');
    assert.strictEqual(
      this.form.data.signature_bits,
      '384',
      'sets the new selected value for signature_bits on the model.'
    );

    await fillIn(GENERAL.inputByAttr('signature_bits'), '0');
    assert.strictEqual(
      this.form.data.signature_bits,
      '0',
      'sets the default value for signature_bits on the model.'
    );
  });
});
