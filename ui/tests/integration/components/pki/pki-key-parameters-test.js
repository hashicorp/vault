/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-role-form';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | pki key parameters', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role', { backend: 'pki' });
    [this.fields] = Object.values(this.model.formFieldGroups.find((g) => g['Key parameters']));
    // TODO: remove Tooltip/ember-basic-dropdown
    setRunOptions({
      rules: {
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should render the component and display the correct defaults', async function (assert) {
    assert.expect(3);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiKeyParameters
          @model={{this.model}}
          @fields={{this.fields}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.keyType).hasValue('rsa');
    assert.dom(SELECTORS.keyBits).hasValue('2048');
    assert.dom(SELECTORS.signatureBits).hasValue('0');
  });

  test('it should set the model properties of key_type and key_bits when key_type changes', async function (assert) {
    assert.expect(11);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiKeyParameters
          @model={{this.model}}
          @fields={{this.fields}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.strictEqual(this.model.keyType, 'rsa', 'sets the default value for key_type on the model.');
    assert.strictEqual(this.model.keyBits, '2048', 'sets the default value for key_bits on the model.');
    assert.strictEqual(
      this.model.signatureBits,
      '0',
      'sets the default value for signature_bits on the model.'
    );
    await fillIn(SELECTORS.keyType, 'ec');
    assert.strictEqual(this.model.keyType, 'ec', 'sets the new selected value for key_type on the model.');
    assert.strictEqual(
      this.model.keyBits,
      '256',
      'sets the new selected value for key_bits on the model based on the selection of key_type.'
    );

    await fillIn(SELECTORS.keyType, 'ed25519');
    assert.strictEqual(
      this.model.keyType,
      'ed25519',
      'sets the new selected value for key_type on the model.'
    );
    assert.strictEqual(
      this.model.keyBits,
      '0',
      'sets the new selected value for key_bits on the model based on the selection of key_type.'
    );

    await fillIn(SELECTORS.keyType, 'ec');
    await fillIn(SELECTORS.keyBits, '384');
    assert.strictEqual(this.model.keyType, 'ec', 'sets the new selected value for key_type on the model.');
    assert.strictEqual(
      this.model.keyBits,
      '384',
      'sets the new selected value for key_bits on the model based on the selection of key_type.'
    );

    await fillIn(SELECTORS.signatureBits, '384');
    assert.strictEqual(
      this.model.signatureBits,
      '384',
      'sets the new selected value for signature_bits on the model.'
    );

    await fillIn(SELECTORS.signatureBits, '0');
    assert.strictEqual(
      this.model.signatureBits,
      '0',
      'sets the default value for signature_bits on the model.'
    );
  });
});
