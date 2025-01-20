/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { render, fillIn, find, waitUntil, click, triggerKeyEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import jsonEditor from '../../pages/components/json-editor';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const component = create(jsonEditor);

module('Integration | Component | json-editor', function (hooks) {
  setupRenderingTest(hooks);

  const JSON_BLOB = `{
    "test": "test"
  }`;
  const BAD_JSON_BLOB = `{
    "test": test
  }`;

  hooks.beforeEach(function () {
    this.set('valueUpdated', sinon.spy());
    this.set('onFocusOut', sinon.spy());
    this.set('json_blob', JSON_BLOB);
    this.set('bad_json_blob', BAD_JSON_BLOB);
    this.set('hashi-read-only-theme', 'hashi-read-only auto-height');
    setRunOptions({
      rules: {
        // CodeMirror has a secret textarea without a label that causes this problem
        label: { enabled: false },
        // TODO: investigate and fix Codemirror styling
        'color-contrast': { enabled: false },
      },
    });
  });

  test('it renders', async function (assert) {
    await render(hbs`<JsonEditor
        @value={{"{}"}}
        @title={{"Test title"}}
        @showToolbar={{true}}
        @readOnly={{true}}
      />`);

    assert.strictEqual(component.title, 'Test title', 'renders the provided title');
    assert.true(component.hasToolbar, 'renders the toolbar');
    assert.true(component.hasCopyButton, 'renders the copy button');
    assert.true(component.hasJSONEditor, 'renders the code mirror modifier');
    assert.ok(component.canEdit, 'json editor can be edited');
  });

  test('it handles editing and linting and styles to json', async function (assert) {
    await render(hbs`<JsonEditor
      @value={{this.json_blob}}
      @readOnly={{false}}
      @valueUpdated={{this.valueUpdated}}
      @onFocusOut={{this.onFocusOut}}
    />`);
    // check for json styling
    assert.dom('.cm-property').hasStyle({
      color: 'rgb(158, 132, 197)',
    });
    assert.dom('.cm-string:nth-child(2)').hasStyle({
      color: 'rgb(29, 219, 163)',
    });

    await fillIn('textarea', this.bad_json_blob);
    await waitUntil(() => find('.CodeMirror-lint-marker-error'));
    assert.dom('.CodeMirror-lint-marker-error').exists('throws linting error');
    assert.dom('.CodeMirror-linenumber').exists('shows line numbers');
  });

  test('it renders the correct theme and expected styling', async function (assert) {
    await render(hbs`<JsonEditor
      @value={{this.json_blob}}
      @theme={{this.hashi-read-only-theme}}
      @readOnly={{true}}
    />`);
    assert.dom('.cm-s-hashi-read-only').hasStyle({
      background: 'rgb(247, 248, 250) none repeat scroll 0% 0% / auto padding-box border-box',
    });
    assert.dom('.CodeMirror-linenumber').doesNotExist('on readOnly does not show line numbers');
  });

  test('it should render example and restore it', async function (assert) {
    this.value = null;
    this.example = 'this is a test example';

    await render(hbs`
      <JsonEditor
        @value={{this.value}}
        @example={{this.example}}
        @mode="ruby"
        @valueUpdated={{fn (mut this.value)}}
      />
    `);

    assert.dom('.CodeMirror-code').hasText(`1${this.example}`, 'Example renders when there is no value');
    assert.dom('[data-test-restore-example]').isDisabled('Restore button disabled when showing example');
    await fillIn('textarea', '');
    await fillIn('textarea', 'adding a value should allow the example to be restored');
    await click('[data-test-restore-example]');
    assert.dom('.CodeMirror-code').hasText(`1${this.example}`, 'Example is restored');
    assert.strictEqual(this.value, null, 'Value is cleared on restore example');
  });

  test('obscure option works correctly', async function (assert) {
    this.set('readOnly', true);
    await render(hbs`<JsonEditor
      @value={{this.json_blob}}
      @allowObscure={{true}}
      @readOnly={{this.readOnly}}
      @valueUpdated={{this.valueUpdated}}
      @onFocusOut={{this.onFocusOut}}
    />`);
    assert.dom('.CodeMirror-code').hasText(`{ "test": "********"}`, 'shows data with obscured values');
    assert.dom('[data-test-toggle-input="revealValues"]').isNotChecked('reveal values toggle is unchecked');
    await click('[data-test-toggle-input="revealValues"]');
    // we are hardcoding the hasText comparison instead of using the JSON_BLOB because we no longer match passed content (ex: @value) to the tracked codemirror instance (ex: this._editor.getVale()) if there are line-breaks or whitespace differences.
    assert.dom('.CodeMirror-code').hasText(`{ "test": "test" }`, 'shows data with real values');
    assert.dom('[data-test-toggle-input="revealValues"]').isChecked('reveal values toggle is checked');
    // turn obscure back on to ensure readonly overrides reveal setting
    await click('[data-test-toggle-input="revealValues"]');
    this.set('readOnly', false);
    assert.dom('[data-test-toggle-input="revealValues"]').doesNotExist('reveal values toggle is hidden');
    assert.dom('.CodeMirror-code').hasText(`{ "test": "test" }`, 'shows data with real values on edit mode');
  });

  test('code-mirror modifier sets value correctly on non json object', async function (assert) {
    // this.value is a tracked property, so anytime it changes the modifier is called. We're testing non-json content by setting the mode to ruby and adding a comment
    this.value = null;
    await render(hbs`
      <JsonEditor
        @value={{this.value}}
        @mode="ruby"
        @valueUpdated={{fn (mut this.value)}}
      />
    `);
    await fillIn('textarea', '#A comment');
    assert.strictEqual(this.value, '#A comment', 'value is set correctly');
    await triggerKeyEvent('textarea', 'keydown', 'Enter');
    assert.strictEqual(
      this.value,
      `#A comment\n`,
      'even after hitting enter the value is still set correctly'
    );
  });
});
