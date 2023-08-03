/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { hbs } from 'ember-cli-htmlbars';
import { click, fillIn, findAll, render, typeIn } from '@ember/test-helpers';
import codemirror from 'vault/tests/helpers/codemirror';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv-v2 | KvSecretForm', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = this.store.createRecord('kv/data', { backend: this.backend });
    this.onSave = () => {};
    this.onCancel = () => {};
  });

  test('it saves a secret', async function (assert) {
    assert.expect(3);
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { foo: 'bar' },
        options: { cas: 0 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onSave={{this.onSave}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    await fillIn(FORM.inputByAttr('path'), this.path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');
    await click(FORM.saveBtn);
  });

  test('it saves nested secrets', async function (assert) {
    assert.expect(4);
    const pathToSecret = 'path/to/secret/';
    this.secret.path = pathToSecret;
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.server.post(`${this.backend}/data/${pathToSecret + this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { foo: 'bar' },
        options: { cas: 0 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onSave={{this.onSave}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    assert.dom(FORM.inputByAttr('path')).hasValue(pathToSecret);
    await typeIn(FORM.inputByAttr('path'), this.path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');
    await click(FORM.saveBtn);
  });

  test('it renders API errors', async function (assert) {
    assert.expect(2);
    this.server.post(`${this.backend}/data/${this.path}`, () => {
      return new Response(500, {}, { errors: ['nope'] });
    });

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onSave={{this.onSave}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    await fillIn(FORM.inputByAttr('path'), this.path);
    await click(FORM.saveBtn);
    assert.dom(FORM.messageError).hasText('Error nope', 'it renders API error');
    assert.dom(FORM.inlineAlert).hasText('There was an error submitting this form.');
  });

  test('it renders kv secret validations', async function (assert) {
    assert.expect(6);

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onSave={{this.onSave}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    await typeIn(FORM.inputByAttr('path'), 'space ');
    assert
      .dom(FORM.validation('path'))
      .hasText(
        `Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.`
      );

    await fillIn(FORM.inputByAttr('path'), ''); // clear input
    await typeIn(FORM.inputByAttr('path'), 'slash/');
    assert.dom(FORM.validation('path')).hasText(`Path can't end in forward slash '/'.`);

    await typeIn(FORM.inputByAttr('path'), 'secret');
    assert
      .dom(FORM.validation('path'))
      .doesNotExist('it removes validation on key up when secret contains slash but does not end in one');

    await click(FORM.toggleJson);
    codemirror().setValue('i am a string and not JSON');
    assert
      .dom(FORM.inlineAlert)
      .hasText('JSON is unparsable. Fix linting errors to avoid data discrepancies.');

    codemirror().setValue('{}'); // clear linting error
    await fillIn(FORM.inputByAttr('path'), '');
    await click(FORM.saveBtn);
    const [pathValidation, formAlert] = findAll(FORM.inlineAlert);
    assert.dom(pathValidation).hasText(`Path can't be blank.`);
    assert.dom(formAlert).hasText('There is an error with this form.');
  });

  // TODO validate onConfirmLeave() unloads model on cancel
  test('it toggles JSON view and editor modifies secretData', async function (assert) {
    assert.expect(6);
    this.onCancel = () => assert.ok(true, 'onCancel callback fires on save success');

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    assert.dom(FORM.dataInputLabel()).hasText('Secret data');
    await click(FORM.toggleJson);
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Secret data');

    assert.strictEqual(
      codemirror().getValue(' '),
      `{   \"\": \"\" }`, // eslint-disable-line no-useless-escape
      'json editor initializes with empty object'
    );
    await fillIn(`${FORM.jsonEditor} textarea`, 'blah');
    assert.strictEqual(codemirror().state.lint.marked.length, 1, 'codemirror lints input');
    codemirror().setValue(`{ "hello": "there"}`);
    assert.propEqual(this.secret.secretData, { hello: 'there' }, 'json editor updates secret data');
    await click(FORM.cancelBtn);
  });

  test('it disables path and prefills secret data when creating a new secret version', async function (assert) {
    assert.expect(6);
    this.secret.secretData = { foo: 'bar' };
    this.secret.path = this.path;

    this.newVersion = this.store.createRecord('kv/data', {
      backend: this.backend,
      path: this.path,
      secretData: this.secret.secretData,
    });

    await render(
      hbs`
        <KvSecretForm
          @previousVersion={{this.secret}}
          @secret={{this.newVersion}}
          @onSave={{this.onSave}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    assert.dom(FORM.inputByAttr('path')).isDisabled();
    assert.dom(FORM.inputByAttr('path')).hasValue(this.path);
    assert.dom(FORM.keyInput()).hasValue('foo');
    assert.dom(FORM.maskedValueInput()).hasValue('bar');

    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Version data');
    await click(FORM.toggleJson);
    assert.dom(FORM.dataInputLabel({ isJson: true })).hasText('Version data');
  });

  test('it renders alert when creating a new secret version from an old version', async function (assert) {
    assert.expect(1);
    const metadata = this.server.create('kv-metadatum');
    metadata.id = 'my-metadata';
    metadata.backend = this.backend;
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...metadata,
    });
    this.metadata = this.store.peekRecord('kv/metadata', 'my-metadata');
    // mimics createRecord in model hook of details/edit route
    this.newVersion = this.store.createRecord('kv/data', {
      backend: this.backend,
      path: this.path,
      secretData: { foo: 'bar' },
    });
    await render(
      hbs`
        <KvSecretForm
          @previousVersion={{2}}
          @metadata={{this.metadata}}
          @secret={{this.newVersion}}
          @onSave={{this.onSave}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    assert
      .dom(FORM.versionAlert)
      .hasText(
        `Warning You are creating a new version based on data from Version 2. The current version for my-secret is Version ${this.metadata.currentVersion}.`
      );
  });
});
