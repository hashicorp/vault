/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { PKI_ROLE_FORM } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | pki-key-usage', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role');
    this.model.backend = 'pki';
  });

  test('it should render the component', async function (assert) {
    assert.expect(6);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiKeyUsage
          @model={{this.model}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.strictEqual(findAll('.b-checkbox').length, 19, 'it render 19 checkboxes');
    assert.dom(PKI_ROLE_FORM.digitalSignature).isChecked('Digital Signature is true by default');
    assert.dom(PKI_ROLE_FORM.keyAgreement).isChecked('Key Agreement is true by default');
    assert.dom(PKI_ROLE_FORM.keyEncipherment).isChecked('Key Encipherment is true by default');
    assert.dom(PKI_ROLE_FORM.any).isNotChecked('Any is false by default');
    assert.dom(GENERAL.inputByAttr('extKeyUsageOids')).exists('Extended Key usage oids renders');
  });

  test('it should set the model properties of key_usage and ext_key_usage based on the checkbox selections', async function (assert) {
    assert.expect(2);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiKeyUsage
          @model={{this.model}}
        />
       </div>
  `,
      { owner: this.engine }
    );

    await click(PKI_ROLE_FORM.digitalSignature);
    await click(PKI_ROLE_FORM.any);
    await click(PKI_ROLE_FORM.serverAuth);

    assert.deepEqual(
      this.model.keyUsage,
      ['KeyAgreement', 'KeyEncipherment'],
      'removes digitalSignature from the model when unchecked.'
    );
    assert.deepEqual(this.model.extKeyUsage, ['Any', 'ServerAuth'], 'adds new checkboxes to when checked');
  });
});
