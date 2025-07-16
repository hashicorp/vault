/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { FORM, parseObject } from 'vault/tests/helpers/kv/kv-selectors';
import codemirror from 'vault/tests/helpers/codemirror';

module('Integration | Component | kv | kv-patch/editor/json-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.onSubmit = sinon.spy();
    this.onCancel = sinon.spy();
    this.isSaving = false;
    this.submitError = '';
    this.subkeys = {
      foo: null,
      bar: {
        baz: null,
        quux: {
          hello: null,
        },
      },
    };
    this.renderComponent = async () => {
      return render(
        hbs`
    <KvPatch::JsonForm
      @onSubmit={{this.onSubmit}}
      @onCancel={{this.onCancel}}
      @isSaving={{this.isSaving}}
      @subkeys={{this.subkeys}}
      @submitError={{this.submitError}}
    />`,
        { owner: this.engine }
      );
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();
    assert.propEqual(parseObject(codemirror), { '': '' }, 'json editor initializes with empty object');
    await click(FORM.saveBtn);
    assert.true(this.onSubmit.calledOnce, 'clicking "Save" calls @onSubmit');
    await click(FORM.cancelBtn);
    assert.true(this.onCancel.calledOnce, 'clicking "Cancel" calls @onCancel');
  });

  test('it reveals subkeys', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.toggleInput('Reveal subkeys')).isNotChecked('toggle is initially unchecked');
    assert.dom(GENERAL.codeBlock('subkeys')).doesNotExist();
    await click(GENERAL.toggleInput('Reveal subkeys'));
    assert.dom(GENERAL.toggleInput('Reveal subkeys')).isChecked();
    assert.dom(GENERAL.codeBlock('subkeys')).hasText(JSON.stringify(this.subkeys, null, 2));

    await click(GENERAL.toggleInput('Reveal subkeys'));
    assert.dom(GENERAL.toggleInput('Reveal subkeys')).isNotChecked();
    assert.dom(GENERAL.codeBlock('subkeys')).doesNotExist('unchecking re-hides subkeys');
  });

  test('it renders linting errors', async function (assert) {
    await this.renderComponent();
    await codemirror().setValue('{ "foo3":  }');
    assert
      .dom(GENERAL.inlineError)
      .hasText('JSON is unparsable. Fix linting errors to avoid data discrepancies.');
    await codemirror().setValue('{ "foo": "bar" }');
    assert.dom(GENERAL.inlineError).doesNotExist('error disappears when linting is fixed');
  });

  test('it renders submit error from parent', async function (assert) {
    this.submitError = 'There was a problem';
    await this.renderComponent();
    assert.dom(GENERAL.inlineError).hasText(this.submitError);
  });

  test('it submits data', async function (assert) {
    this.submitError = 'There was a problem';
    await this.renderComponent();
    await codemirror().setValue('{ "foo": "bar" }');
    await click(FORM.saveBtn);
    const [data] = this.onSubmit.lastCall.args;
    assert.propEqual(data, { foo: 'bar' }, `onSubmit called with ${JSON.stringify(data)}`);
  });
});
