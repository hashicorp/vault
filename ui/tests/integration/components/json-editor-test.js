/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, triggerKeyEvent, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { SELECTORS } from '../../pages/components/json-editor';
import sinon from 'sinon';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { createLongJson } from 'vault/tests/helpers/secret-engine/secret-engine-helpers';
import codemirror, { getCodeEditorValue, setCodeEditorValue } from 'vault/tests/helpers/codemirror';

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
    this.set('long_json', JSON.stringify(createLongJson(), null, `\t`));
    this.set('hashi-read-only-theme', 'hashi-read-only auto-height');
    setRunOptions({
      rules: {
        // CodeMirror has a secret textarea without a label that causes this problem
        label: { enabled: false },
        // TODO: investigate and fix Codemirror styling
        'color-contrast': { enabled: false },
        // failing on .CodeMirror-scroll
        'scrollable-region-focusable': { enabled: false },
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

    assert.dom(SELECTORS.title).hasText('Test title', 'renders the provided title');
    assert.dom(SELECTORS.toolbar).exists('renders the toolbar');
    assert.dom(SELECTORS.copy).exists('renders the copy button');
    assert.dom(SELECTORS.codeBlock).exists('renders the code block');
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

    let view = codemirror();
    await waitFor('.cm-editor');

    let editorValue = getCodeEditorValue(view);
    assert.strictEqual(editorValue, this.example, 'Example renders when there is no value');
    assert.dom('[data-test-restore-example]').isDisabled('Restore button disabled when showing example');

    setCodeEditorValue(view, '');
    setCodeEditorValue(view, 'adding a value should allow the example to be restored');
    await click('[data-test-restore-example]');

    view = codemirror();
    await waitFor('.cm-editor');
    editorValue = getCodeEditorValue(view);

    assert.deepEqual(editorValue, this.example, 'Example is restored');
    assert.strictEqual(this.value, null, 'Value is cleared on restore example');
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

    await waitFor('.cm-editor');
    const view = codemirror();
    setCodeEditorValue(view, '#A comment');
    assert.strictEqual(this.value, '#A comment', 'value is set correctly');

    await triggerKeyEvent('.cm-content', 'keydown', 'Enter');
    assert.strictEqual(
      this.value,
      `\n#A comment`,
      'even after hitting enter the value is still set correctly'
    );
  });
});
