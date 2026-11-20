/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, settled, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import SecretForm from 'vault/forms/secret';
import sinon from 'sinon';

module('Integration | Component | secret edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('key', { id: 'Foobar' });
    this.set('capabilities', null);
    this.root = {
      label: 'kv',
      text: 'kv',
      path: 'vault.cluster.secrets.backend.list-root',
      model: 'kv',
    };
  });

  test('it disables the UI view button in show mode when data is an advanced format', async function (assert) {
    this.set('mode', 'show');
    this.set('model', {
      secretData: {
        int: 2,
        null: null,
        float: 1.234,
      },
    });

    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );
    // Non-string values are "advanced" and can't be shown in the key/value UI,
    // so the UI view button is disabled (JSON/YAML remain available).
    assert.dom(GENERAL.button('ui')).isDisabled();
  });

  test('it does JSON toggle in show mode when showing string data', async function (assert) {
    this.set('mode', 'show');
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });

    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );
    assert.dom(GENERAL.button('ui')).isNotDisabled();
  });

  test('it shows an error when creating and data is not an object', async function (assert) {
    this.set('mode', 'create');
    this.set('model', {
      secretData: null,
    });

    await render(
      hbs`<SecretEdit
        @mode={{this.mode}}
        @root={{this.root}}
        @model={{this.model}}
        @preferAdvancedEdit={{true}}
        @key={{this.key}}
      />`
    );

    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom(GENERAL.messageError).includesText('Vault expects data to be formatted as an JSON object');
  });

  test('it requires a path before creating a secret', async function (assert) {
    const write = sinon.stub(this.owner.lookup('service:api').secrets, 'kvV1Write').resolves({});
    this.set('mode', 'create');
    this.set('model', {
      id: '',
      backend: 'generic',
      path: '',
      pathAttr: 'path',
      secretData: {},
    });

    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );
    await click(GENERAL.submitButton);

    assert.dom(this.element).includesText("Path can't be blank.", 'the create form reports the missing path');
    assert.false(write.called, 'the secret is not written without a path');
  });

  test('it enables saving a valid secret', async function (assert) {
    this.set('mode', 'create');
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
    });
    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );
    assert.dom(GENERAL.submitButton).isNotDisabled();
  });

  test('it displays a permission error when secret creation is denied', async function (assert) {
    const api = this.owner.lookup('service:api');
    const write = sinon.stub(api.secrets, 'kvV1Write').rejects(new Error('request failed'));
    sinon.stub(api, 'parseError').resolves({ message: 'permission denied' });
    this.set('mode', 'create');
    this.set('model', {
      id: '',
      backend: 'generic',
      path: '',
      pathAttr: 'path',
      secretData: { key: 'value' },
    });

    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );
    await fillIn('[data-test-secret-path="create"]', 'denied-create');
    await click(GENERAL.submitButton);
    await waitFor(GENERAL.messageError);

    assert.dom(GENERAL.messageError).includesText('permission denied');
    assert.true(write.calledOnce, 'the denied secret is submitted once');
  });

  test('it renders a blank row for a new secret and saves only the entered data', async function (assert) {
    const write = sinon.stub(this.owner.lookup('service:api').secrets, 'kvV1Write').resolves({});
    sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.model = new SecretForm({ backend: 'generic', secretData: {} }, { isNew: true });

    await render(
      hbs`<SecretEdit @mode="create" @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );

    assert.dom('[data-test-secret-key="create"]').hasValue('', 'a new secret starts with a blank key');
    await fillIn('[data-test-secret-path="create"]', 'my-secret');
    await fillIn('[data-test-secret-key="create"]', 'username');
    await fillIn('[data-test-secret-value] textarea', 'example');
    await click('[data-test-secret-add-row]');
    await click(GENERAL.submitButton);

    assert.deepEqual(
      write.args,
      [['my-secret', 'generic', { username: 'example' }]],
      'blank rows are not sent'
    );
  });

  test('it shows an error when editing and the data is not an object', async function (assert) {
    this.set('mode', 'edit');
    this.set('capabilities', {
      canUpdate: true,
    });
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
      canReadSecretData: true,
    });

    await render(
      hbs`<SecretEdit
        @mode={{this.mode}}
        @root={{this.root}}
        @model={{this.model}}
        @capabilities={{this.capabilities}}
        @preferAdvancedEdit={{true}}
        @key={{this.key}}
      />`
    );

    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom(GENERAL.messageError).includesText('Vault expects data to be formatted as an JSON object');
  });

  // Permission gating
  // A user who can write but not read a secret (failedServerRead)
  // must never see the secret data in ANY view (UI/JSON/YAML).
  test('it hides the format toggle and does not leak secret data for write-without-read (show mode)', async function (assert) {
    this.set('mode', 'show');
    this.set('model', {
      failedServerRead: true,
      // even if data were somehow present on the model, it must not be rendered
      secretData: { password: 'super-secret-value' },
    });

    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @key={{this.key}} />`
    );

    assert
      .dom('[data-test-write-without-read-empty-message]')
      .exists('shows the "no permission to read this secret" empty state');
    assert.dom('[data-test-button]').doesNotExist('the UI/JSON/YAML format toggle is hidden');
    assert.dom('.cm-editor').doesNotExist('no secret data editor is rendered');
    assert
      .dom(this.element)
      .doesNotIncludeText('super-secret-value', 'the secret value is not leaked in any format');
  });

  test('it renders secret data as YAML when the YAML view is selected (show mode)', async function (assert) {
    this.set('mode', 'show');
    this.set('onToggleAdvancedEdit', () => {});
    this.set('model', {
      secretData: { password: 'super-secret-value', count: '2' },
    });

    await render(
      hbs`<SecretEdit
        @mode={{this.mode}}
        @root={{this.root}}
        @model={{this.model}}
        @key={{this.key}}
        @preferAdvancedEdit={{true}}
        @onToggleAdvancedEdit={{this.onToggleAdvancedEdit}}
      />`
    );

    await click(GENERAL.button('yaml'));

    const text = this.element.textContent;

    assert.true(text.includes('password: super-secret-value'), 'renders YAML-formatted secret data');
    assert.false(text.includes('"password"'), 'output is YAML, not JSON');
  });

  test('it renders secret data as JSON when the JSON view is selected (show mode)', async function (assert) {
    this.set('mode', 'show');
    this.set('onToggleAdvancedEdit', () => {});
    this.set('model', {
      secretData: { password: 'super-secret-value', count: '2' },
    });

    await render(
      hbs`<SecretEdit
        @mode={{this.mode}}
        @root={{this.root}}
        @model={{this.model}}
        @key={{this.key}}
        @preferAdvancedEdit={{true}}
        @onToggleAdvancedEdit={{this.onToggleAdvancedEdit}}
      />`
    );

    await click(GENERAL.button('json'));
    const text = this.element.textContent;

    assert.true(text.includes('"password"'), 'renders JSON-formatted secret data (quoted key)');
    assert.true(text.includes('super-secret-value'), 'renders the secret value');
  });
});
