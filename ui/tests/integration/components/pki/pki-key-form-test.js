/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-key-form';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Integration | Component | pki key form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki'); // https://github.com/ember-engines/ember-engines/pull/653

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/key');
    this.backend = 'pki-test';
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = this.backend;
  });

  test('it should render fields and show validation messages', async function (assert) {
    assert.expect(7);
    await render(
      hbs`
      <PkiKeyForm
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
      `,
      { owner: this.engine }
    );
    assert.dom(SELECTORS.keyNameInput).exists('renders name input');
    assert.dom(SELECTORS.typeInput).exists('renders type input');
    assert.dom(SELECTORS.keyTypeInput).exists('renders key type input');
    assert.dom(SELECTORS.keyBitsInput).exists('renders key bits input');

    await click(SELECTORS.keyCreateButton);
    assert
      .dom(SELECTORS.fieldErrorByName('type'))
      .hasTextContaining('Type is required.', 'renders presence validation for type of key');
    assert
      .dom(SELECTORS.fieldErrorByName('keyType'))
      .hasTextContaining('Please select a key type.', 'renders selection prompt for key type');
    assert
      .dom(SELECTORS.validationError)
      .hasTextContaining('There are 2 errors with this form.', 'renders correct form error count');
  });

  test('it generates a key type=exported', async function (assert) {
    assert.expect(4);
    this.server.post(`/${this.backend}/keys/generate/exported`, (schema, req) => {
      assert.ok(true, 'Request made to the correct endpoint to generate exported key');
      const request = JSON.parse(req.requestBody);
      assert.propEqual(
        request,
        {
          key_name: 'test-key',
          key_type: 'rsa',
          key_bits: '2048',
        },
        'sends params in correct type'
      );
      return {};
    });

    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(
      hbs`
      <PkiKeyForm
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
      `,
      { owner: this.engine }
    );

    await fillIn(SELECTORS.keyNameInput, 'test-key');
    await fillIn(SELECTORS.typeInput, 'exported');
    assert.dom(SELECTORS.keyBitsInput).isDisabled('key bits disabled when no key type selected');
    await fillIn(SELECTORS.keyTypeInput, 'rsa');
    await click(SELECTORS.keyCreateButton);
  });

  test('it generates a key type=internal', async function (assert) {
    assert.expect(4);
    this.server.post(`/${this.backend}/keys/generate/internal`, (schema, req) => {
      assert.ok(true, 'Request made to the correct endpoint to generate internal key');
      const request = JSON.parse(req.requestBody);
      assert.propEqual(
        request,
        {
          key_name: 'test-key',
          key_type: 'rsa',
          key_bits: '2048',
        },
        'sends params in correct type'
      );
      return {};
    });
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(
      hbs`
      <PkiKeyForm
         @model={{this.model}}
         @onCancel={{this.onCancel}}
         @onSave={{this.onSave}}
       />
      `,
      { owner: this.engine }
    );

    await fillIn(SELECTORS.keyNameInput, 'test-key');
    await fillIn(SELECTORS.typeInput, 'internal');
    assert.dom(SELECTORS.keyBitsInput).isDisabled('key bits disabled when no key type selected');
    await fillIn(SELECTORS.keyTypeInput, 'rsa');
    await click(SELECTORS.keyCreateButton);
  });

  test('it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(5);
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');
    await render(
      hbs`
        <PkiKeyForm
          @model={{this.model}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
        />
      `,
      { owner: this.engine }
    );

    await click(SELECTORS.keyCancelButton);
    assert.true(this.model.isDestroyed, 'new model is unloaded on cancel');

    this.store.pushPayload('pki/key', {
      modelName: 'pki/key',
      key_name: 'test-key',
      type: 'exported',
      key_id: 'some-key-id',
      key_type: 'rsa',
      key_bits: '2048',
    });
    this.model = this.store.peekRecord('pki/key', 'some-key-id');

    await render(
      hbs`
          <PkiKeyForm
            @model={{this.model}}
            @onCancel={{this.onCancel}}
            @onSave={{this.onSave}}
          />
        `,
      { owner: this.engine }
    );

    await fillIn(SELECTORS.keyNameInput, 'new-name');
    await click(SELECTORS.keyCancelButton);
    assert.strictEqual(this.model.keyName, 'test-key', 'Model name rolled back on cancel');
    await fillIn(SELECTORS.keyNameInput, 'new-name');
    await click(SELECTORS.keyCreateButton);
    assert.strictEqual(this.model.keyName, 'new-name', 'Model name correctly save on create');
  });
});
