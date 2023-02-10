import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { render, typeIn, find, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import jsonEditor from '../../pages/components/json-editor';
import sinon from 'sinon';

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

    await typeIn('textarea', this.bad_json_blob);
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
});
