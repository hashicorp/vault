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

module('Integration | Component | pki-key-usage', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role');
    // add fields that openapi normally hydrates
    // ideally we pull this from the openapi schema in the future
    const openapifields = [
      {
        name: 'clientFlag',
        type: 'boolean',
        options: {
          editType: 'boolean',
          helpText:
            'If set, certificates are flagged for client auth use. Defaults to true. See also RFC 5280 Section 4.2.1.12.',
          fieldGroup: 'default',
          defaultValue: true,
        },
      },
      {
        name: 'serverFlag',
        type: 'boolean',
        options: {
          editType: 'boolean',
          helpText:
            'If set, certificates are flagged for server auth use. Defaults to true. See also,  RFC 5280 Section 4.2.1.12.',
          fieldGroup: 'default',
          defaultValue: true,
        },
      },
      {
        name: 'codeSigningFlag',
        type: 'boolean',
        options: {
          editType: 'boolean',
          helpText:
            'If set, certificates are flagged for code signing use. Defaults to false. See also RFC 5280 Section 4.2.1.12.',
          fieldGroup: 'default',
        },
      },
      {
        name: 'emailProtectionFlag',
        type: 'boolean',
        options: {
          editType: 'boolean',
          helpText:
            'If set, certificates are flagged for email protection use. Defaults to false. See also RFC 5280 Section 4.2.1.12.',
          fieldGroup: 'default',
        },
      },
    ];
    this.model._allByKey = {};
    openapifields.forEach((f) => {
      this.model._allByKey[f.name] = f;
      this.model[f.name] = f.options.defaultValue;
    });
    this.model.backend = 'pki';
  });

  test('it should render the component', async function (assert) {
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

    assert.strictEqual(findAll('.b-checkbox').length, 23, 'it render 23 checkboxes');
    assert.dom(PKI_ROLE_FORM.digitalSignature).isChecked('Digital Signature is true by default');
    assert.dom(PKI_ROLE_FORM.keyAgreement).isChecked('Key Agreement is true by default');
    assert.dom(PKI_ROLE_FORM.keyEncipherment).isChecked('Key Encipherment is true by default');
    assert.dom(PKI_ROLE_FORM.any).isNotChecked('Any is false by default');
    assert.dom(GENERAL.inputByAttr('clientFlag')).isChecked();
    assert.dom(GENERAL.inputByAttr('serverFlag')).isChecked();
    assert.dom(GENERAL.inputByAttr('codeSigningFlag')).isNotChecked();
    assert.dom(GENERAL.inputByAttr('emailProtectionFlag')).isNotChecked();
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
