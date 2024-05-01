/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { PKI_KEY_FORM } from 'vault/tests/helpers/pki/pki-selectors';

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
    this.onCancel = sinon.spy();
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
    assert.dom(GENERAL.inputByAttr('keyName')).exists('renders name input');
    assert.dom(GENERAL.inputByAttr('type')).exists('renders type input');
    assert.dom(GENERAL.inputByAttr('keyType')).exists('renders key type input');
    assert.dom(GENERAL.inputByAttr('keyBits')).exists('renders key bits input');

    await click(GENERAL.saveButton);
    assert
      .dom(GENERAL.validation('type'))
      .hasTextContaining('Type is required.', 'renders presence validation for type of key');
    assert
      .dom(GENERAL.validation('keyType'))
      .hasTextContaining('Please select a key type.', 'renders selection prompt for key type');
    assert
      .dom(PKI_KEY_FORM.validationError)
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
      return { key_id: 'test' };
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

    await fillIn(GENERAL.inputByAttr('keyName'), 'test-key');
    await fillIn(GENERAL.inputByAttr('type'), 'exported');
    assert.dom(GENERAL.inputByAttr('keyBits')).isDisabled('key bits disabled when no key type selected');
    await fillIn(GENERAL.inputByAttr('keyType'), 'rsa');
    await click(GENERAL.saveButton);
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
      return { key_id: 'test' };
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

    await fillIn(GENERAL.inputByAttr('keyName'), 'test-key');
    await fillIn(GENERAL.inputByAttr('type'), 'internal');
    assert.dom(GENERAL.inputByAttr('keyBits')).isDisabled('key bits disabled when no key type selected');
    await fillIn(GENERAL.inputByAttr('keyType'), 'rsa');
    await click(GENERAL.saveButton);
  });
});
