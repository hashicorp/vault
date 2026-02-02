/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_ROLE_FORM } from 'vault/tests/helpers/pki/pki-selectors';
import PkiRoleForm from 'vault/forms/secrets/pki/role';

module('Integration | Component | pki-key-usage', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.owner.lookup('service:secret-mount-path').update('pki');
    this.form = new PkiRoleForm({}, { isNew: true });

    this.renderComponent = () => render(hbs`<PkiKeyUsage @form={{this.form}} />`, { owner: this.engine });
  });

  test('it should render the component', async function (assert) {
    await this.renderComponent();
    assert.strictEqual(findAll('.b-checkbox').length, 19, 'it render 19 checkboxes');
    assert.dom(PKI_ROLE_FORM.digitalSignature).isChecked('Digital Signature is true by default');
    assert.dom(PKI_ROLE_FORM.keyAgreement).isChecked('Key Agreement is true by default');
    assert.dom(PKI_ROLE_FORM.keyEncipherment).isChecked('Key Encipherment is true by default');
    assert.dom(PKI_ROLE_FORM.any).isNotChecked('Any is false by default');
    assert.dom(GENERAL.inputByAttr('client_flag')).isChecked();
    assert.dom(GENERAL.inputByAttr('server_flag')).isChecked();
    assert.dom(GENERAL.inputByAttr('code_signing_flag')).isNotChecked();
    assert.dom(GENERAL.inputByAttr('email_protection_flag')).isNotChecked();
    assert.dom(GENERAL.inputByAttr('ext_key_usage_oids')).exists('Extended Key usage oids renders');
  });

  test('it should set values of key_usage and ext_key_usage based on the checkbox selections', async function (assert) {
    assert.expect(2);

    await this.renderComponent();
    await click(PKI_ROLE_FORM.digitalSignature);
    await click(PKI_ROLE_FORM.any);
    await click(PKI_ROLE_FORM.serverAuth);

    assert.deepEqual(
      this.form.data.key_usage,
      ['KeyAgreement', 'KeyEncipherment'],
      'removes DigitalSignature from key_usage when unchecked.'
    );
    assert.deepEqual(
      this.form.data.ext_key_usage,
      ['Any', 'ServerAuth'],
      'adds new checkboxes to when checked'
    );
  });
});
