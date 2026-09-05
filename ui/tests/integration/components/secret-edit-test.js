/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render, settled, waitFor } from '@ember/test-helpers';
import { resolve } from 'rsvp';
import { run } from '@ember/runloop';
import Service from '@ember/service';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';

let capabilities;
const storeService = Service.extend({
  queryRecord() {
    return resolve(capabilities);
  },
});
module('Integration | Component | secret edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    capabilities = null;
    this.set('key', { id: 'Foobar' });
    this.root = {
      label: 'kv',
      text: 'kv',
      path: 'vault.cluster.secrets.backend.list-root',
      model: 'kv',
    };
    run(() => {
      this.owner.unregister('service:store');
      this.owner.register('service:store', storeService);
    });
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
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @preferAdvancedEdit={{true}} @key={{this.key}} />`
    );

    await waitFor('.cm-editor');
    const editor = codemirror();
    setCodeEditorValue(editor, JSON.stringify([{ foo: 'bar' }]));
    await settled();
    assert.dom(GENERAL.messageError).includesText('Vault expects data to be formatted as an JSON object');
  });

  test('it allows saving when the model isError', async function (assert) {
    this.set('mode', 'create');
    this.set('model', {
      isError: true,
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

  test('it shows an error when editing and the data is not an object', async function (assert) {
    this.set('mode', 'edit');
    capabilities = {
      canUpdate: true,
    };
    this.set('model', {
      secretData: {
        int: '2',
        null: 'null',
        float: '1.234',
      },
      canReadSecretData: true,
    });

    await render(
      hbs`<SecretEdit @mode={{this.mode}} @root={{this.root}} @model={{this.model}} @preferAdvancedEdit={{true}} @key={{this.key}} />`
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
