/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { click, fillIn, find, render, settled, waitFor, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { TTL_PICKER as TTL } from 'vault/tests/helpers/components/ttl-picker-selectors';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';
import codemirror, { getCodeEditorValue, setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import sinon from 'sinon';

async function setEditorValue(value) {
  await waitFor('.cm-editor');
  const editor = codemirror();
  setCodeEditorValue(editor, value);
  return settled();
}

module('Integration | Component | tools/wrap', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Wrap />`);
    };
    this.wrapData = `{"foo": "bar"}`;
    this.token = 'blah.jhfel7SmsVeZwihaGiIKHGh2cy5XZWtEeEt5WmRwS1VYSTNDb1BBVUNsVFAQ3JIK';
    // default mirage response here is overridden in some tests
    this.server.post('sys/wrapping/wrap', () => {
      // removed superfluous response data for this test
      return { wrap_info: { token: this.token } };
    });
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    await waitFor('.cm-editor');

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Wrap data', 'Title renders');
    assert.dom('[data-test-toggle-label="json"]').hasText('JSON');
    assert.dom('[data-test-component="json-editor-title"]').hasText('Data to wrap');
    assert.dom('.hds-code-editor__description').hasText('json-formatted');

    const editor = codemirror();
    const editorValue = getCodeEditorValue(editor);

    assert.strictEqual(
      editorValue,
      `{
  "": ""
}`,
      'json editor initializes with empty object that includes whitespace'
    );
    assert.dom(TTL.toggleByLabel('Wrap TTL')).isNotChecked('Wrap TTL defaults to unchecked');
    assert.dom(GENERAL.submitButton).isEnabled();
    assert.dom(TS.toolsInput('wrapping-token')).doesNotExist();
    assert.dom(GENERAL.button('Back')).doesNotExist();
    assert.dom(GENERAL.button('Done')).doesNotExist();

    await click(TTL.toggleByLabel('Wrap TTL'));
    assert.dom(TTL.valueInputByLabel('Wrap TTL')).hasValue('30', 'ttl defaults to 30 when toggled');
    assert.dom(TTL.ttlUnit).hasValue('m', 'ttl defaults to minutes when toggled');
  });

  test('it renders errors', async function (assert) {
    this.server.post('sys/wrapping/wrap', () => new Response(500, {}, { errors: ['Something is wrong'] }));
    await this.renderComponent();
    await click(GENERAL.submitButton);
    await waitUntil(() => find(GENERAL.messageError));
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong', 'Error renders');
  });

  test('it submits with defaults', async function (assert) {
    assert.expect(6);
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');

    this.server.post('sys/wrapping/wrap', (schema, { requestBody, requestHeaders }) => {
      const payload = JSON.parse(requestBody);
      assert.propEqual(payload, JSON.parse(this.wrapData), `payload contains data: ${requestBody}`);
      assert.strictEqual(requestHeaders['x-vault-wrap-ttl'], '30m', 'request header has default wrap ttl');
      return {
        wrap_info: {
          token: this.token,
          accessor: '5yjKx6Om9NmBx1mjiN1aIrnm',
          ttl: 1800,
          creation_time: '2024-06-07T12:02:22.096254-07:00',
          creation_path: 'sys/wrapping/wrap',
        },
      };
    });

    await this.renderComponent();
    await setEditorValue(this.wrapData);
    await click(GENERAL.submitButton);
    await waitUntil(() => find(TS.toolsInput('wrapping-token')));
    assert.true(flashSpy.calledWith('Wrap was successful.'), 'it renders success flash');
    assert.dom(TS.toolsInput('wrapping-token')).hasText(this.token);
    assert.dom('label').hasText('Wrapped token');
    assert.dom('.CodeMirror').doesNotExist();
  });

  test('it submits with updated ttl', async function (assert) {
    assert.expect(2);
    this.server.post('sys/wrapping/wrap', (schema, { requestBody, requestHeaders }) => {
      const payload = JSON.parse(requestBody);
      assert.propEqual(payload, JSON.parse(this.wrapData), `payload contains data: ${requestBody}`);
      assert.strictEqual(requestHeaders['x-vault-wrap-ttl'], '1200s', 'request header has updated wrap ttl');
      // only testing payload/header assertions, no need for return here
      return {};
    });

    await this.renderComponent();
    await setEditorValue(this.wrapData);
    await click(TTL.toggleByLabel('Wrap TTL'));
    await fillIn(TTL.valueInputByLabel('Wrap TTL'), '20');
    await click(GENERAL.submitButton);
  });

  test('it toggles between views and preserves input data', async function (assert) {
    assert.expect(7);
    await this.renderComponent();
    await setEditorValue(this.wrapData);
    assert.dom('[data-test-component="json-editor-title"]').hasText('Data to wrap');
    assert.dom('.hds-code-editor__description').hasText('json-formatted');
    await click(GENERAL.toggleInput('json'));

    assert.dom('[data-test-component="json-editor-title"]').doesNotExist();
    assert.dom(GENERAL.kvObjectEditor.key('0')).hasValue('foo');
    assert.dom(GENERAL.kvObjectEditor.value('0')).hasValue('bar');
    await click(GENERAL.toggleInput('json'));
    assert.dom('[data-test-component="json-editor-title"]').exists();

    await waitFor('.cm-editor');
    const editor = codemirror();
    const editorValue = getCodeEditorValue(editor);
    assert.strictEqual(
      editorValue,
      `{
  "foo": "bar"
}`,
      'json editor has original data'
    );
  });

  test('it submits from kv view', async function (assert) {
    assert.expect(6);

    const multilineData = `this is a multi-line secret
      that contains
      some seriously important config`;
    const flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'success');
    const updatedWrapData = JSON.stringify({
      ...JSON.parse(this.wrapData),
      foo: 'bar',
      foo2: multilineData,
    });

    this.server.post('sys/wrapping/wrap', (schema, { requestBody, requestHeaders }) => {
      const payload = JSON.parse(requestBody);
      assert.propEqual(payload, JSON.parse(updatedWrapData), `payload contains data: ${requestBody}`);
      assert.strictEqual(requestHeaders['x-vault-wrap-ttl'], '30m', 'request header has default wrap ttl');
      return {
        wrap_info: {
          token: this.token,
          accessor: '5yjKx6Om9NmBx1mjiN1aIrnm',
          ttl: 1800,
          creation_time: '2024-06-07T12:02:22.096254-07:00',
          creation_path: 'sys/wrapping/wrap',
        },
      };
    });

    await this.renderComponent();
    await click(GENERAL.toggleInput('json'));
    await fillIn(GENERAL.kvObjectEditor.key('0'), 'foo');
    await fillIn(GENERAL.kvObjectEditor.value('0'), 'bar');
    await click('[data-test-kv-add-row="0"]');
    await fillIn(GENERAL.kvObjectEditor.key('1'), 'foo2');
    await fillIn(GENERAL.kvObjectEditor.value('1'), multilineData);
    await click(GENERAL.submitButton);
    await waitUntil(() => find(TS.toolsInput('wrapping-token')));
    assert.true(flashSpy.calledWith('Wrap was successful.'), 'it renders success flash');
    assert.dom(TS.toolsInput('wrapping-token')).hasText(this.token);
    assert.dom('label').hasText('Wrapped token');
    assert.dom('.CodeMirror').doesNotExist();
  });

  test('it resets on done', async function (assert) {
    await this.renderComponent();
    await setEditorValue(this.wrapData);
    await click(TTL.toggleByLabel('Wrap TTL'));
    await fillIn(TTL.valueInputByLabel('Wrap TTL'), '20');
    await click(GENERAL.submitButton);

    await waitUntil(() => find(GENERAL.button('Done')));
    await click(GENERAL.button('Done'));
    await waitFor('.cm-editor');
    const editor = codemirror();
    const editorValue = getCodeEditorValue(editor);
    assert.strictEqual(
      editorValue,
      `{
  "": ""
}`,
      'json editor initializes with empty object that includes whitespace'
    );
    assert.dom(TTL.toggleByLabel('Wrap TTL')).isNotChecked('Wrap TTL resets to unchecked');
    await click(TTL.toggleByLabel('Wrap TTL'));
    assert.dom(TTL.valueInputByLabel('Wrap TTL')).hasValue('30', 'ttl resets to default when toggled');
  });

  test('it preserves input data on back', async function (assert) {
    await this.renderComponent();
    await setEditorValue(this.wrapData);
    await click(GENERAL.submitButton);

    await waitUntil(() => find(GENERAL.button('Back')));
    await click(GENERAL.button('Back'));
    await waitFor('.cm-editor');
    const editor = codemirror();
    const editorValue = getCodeEditorValue(editor);
    assert.strictEqual(
      editorValue,
      `{
  "foo": "bar"
}`,
      'json editor has original data'
    );
    assert.dom(TTL.toggleByLabel('Wrap TTL')).isNotChecked('Wrap TTL defaults to unchecked');
  });

  test('it renders/hides warning based on json linting', async function (assert) {
    await this.renderComponent();
    await setEditorValue(`{bad json}`);
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'JSON is unparsable. Fix linting errors to avoid data discrepancies.',
        'Linting error message is shown for json view'
      );
    await setEditorValue(this.wrapData);
    assert.dom(GENERAL.inlineAlert).doesNotExist();
  });

  test('it hides json warning on back and on done', async function (assert) {
    await this.renderComponent();
    await setEditorValue(`{bad json}`);
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'JSON is unparsable. Fix linting errors to avoid data discrepancies.',
        'Linting error message is shown for json view'
      );
    await click(GENERAL.submitButton);
    await waitUntil(() => find(GENERAL.button('Done')));
    await click(GENERAL.button('Done'));
    assert.dom(GENERAL.inlineAlert).doesNotExist();

    await setEditorValue(`{bad json}`);
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'JSON is unparsable. Fix linting errors to avoid data discrepancies.',
        'Linting error message is shown for json view'
      );
    await click(GENERAL.submitButton);
    await waitUntil(() => find(GENERAL.button('Back')));
    await click(GENERAL.button('Back'));
    assert.dom(GENERAL.inlineAlert).doesNotExist();
  });
});
