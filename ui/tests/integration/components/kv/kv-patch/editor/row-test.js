/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { blur, click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';
import sinon from 'sinon';

module('Integration | Component | kv | kv-patch/editor/row', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    // for simplicity using an object to test views
    this.updateKey = sinon.spy();
    this.undoKey = sinon.spy();
    this.updateValue = sinon.spy();
    this.updateState = sinon.spy();
    this.kvClass = {
      key: 'foo',
      value: undefined,
      state: 'disabled',
      updateValue: this.updateValue,
      updateState: this.updateState,
    };

    this.renderComponent = async () => {
      return render(
        hbs`
      <KvPatch::Editor::Row
        @idx={{0}}
        @kvClass={{this.kvClass}}
        @isOriginalSubkey={{this.isOriginalSubkey}}
        @updateKey={{this.updateKey}}
        @undoKey={{this.undoKey}}
      />`,
        { owner: this.engine }
      );
    };
  });

  module('it renders original subkeys', function (hooks) {
    hooks.beforeEach(function () {
      this.isOriginalSubkey = () => true;
    });

    test('in disabled state', async function (assert) {
      await this.renderComponent();
      assert.dom(FORM.keyInput()).hasValue('foo');
      assert.dom(FORM.keyInput()).isDisabled();

      assert.dom(FORM.valueInput()).hasValue('');
      assert.dom(FORM.valueInput()).isDisabled();

      assert.dom(FORM.patchEdit()).exists();
      assert.dom(FORM.patchDelete()).exists();
      assert.dom(FORM.patchUndo()).doesNotExist();
    });

    test('in enabled state', async function (assert) {
      this.kvClass.state = 'enabled';
      await this.renderComponent();
      assert.dom(FORM.keyInput()).hasValue('foo');
      assert.dom(FORM.keyInput()).hasAttribute('readonly');

      assert.dom(FORM.valueInput()).hasValue('');
      assert.dom(FORM.valueInput()).isEnabled();

      assert.dom(FORM.patchEdit()).doesNotExist();
      assert.dom(FORM.patchDelete()).doesNotExist();
      assert.dom(FORM.patchUndo()).hasText('Cancel');
    });

    test('in deleted state', async function (assert) {
      this.kvClass.state = 'deleted';
      await this.renderComponent();
      assert.dom(FORM.keyInput()).hasValue('foo');
      assert.dom(FORM.keyInput()).hasAttribute('readonly');

      assert.dom(FORM.valueInput()).hasValue('');
      assert.dom(FORM.keyInput()).hasAttribute('readonly');

      assert.dom(FORM.patchEdit()).doesNotExist();
      assert.dom(FORM.patchDelete()).doesNotExist();
      assert.dom(FORM.patchUndo()).hasText('Cancel');
      assert.dom(FORM.patchAlert('delete', 0)).hasText('This key value pair is marked for deletion.');
    });

    test('it clicks undo', async function (assert) {
      this.kvClass.state = 'enabled';
      await this.renderComponent();
      await click(FORM.patchUndo());

      const [arg] = this.undoKey.lastCall.args;
      assert.propEqual(arg, this.kvClass, 'undoKey is called with class');
    });
  });

  module('it renders new subkeys', function (hooks) {
    hooks.beforeEach(function () {
      this.isOriginalSubkey = () => false;
      this.kvClass = { ...this.kvClass, value: 'bar', state: 'enabled' };
    });

    // only test this state because new keys are only ever 'enabled'
    test('in enabled state', async function (assert) {
      await this.renderComponent();
      assert.dom(FORM.keyInput()).hasValue('foo');
      assert.dom(FORM.keyInput()).isNotDisabled();

      assert.dom(FORM.valueInput()).hasValue('bar');
      assert.dom(FORM.valueInput()).isNotDisabled();

      assert.dom(FORM.patchEdit()).doesNotExist();
      assert.dom(FORM.patchDelete()).doesNotExist();
      assert.dom(FORM.patchUndo()).hasText('Remove');
    });

    test('it updates key', async function (assert) {
      this.kvClass.key = '';
      await this.renderComponent();
      await fillIn(FORM.keyInput(), 'foo');
      await blur(FORM.keyInput());

      const [arg, event] = this.updateKey.lastCall.args;
      assert.propEqual(arg, this.kvClass, 'updateKey is called with class object');
      assert.strictEqual(event.target.value, 'foo', 'updateKey is called with event');
    });

    test('it clicks undo', async function (assert) {
      await this.renderComponent();
      await click(FORM.patchUndo());

      const [arg] = this.undoKey.lastCall.args;
      assert.propEqual(arg, this.kvClass, 'undoKey is called with class');
    });
  });

  test('it updates value', async function (assert) {
    this.kvClass.state = 'enabled';
    await this.renderComponent();
    await fillIn(FORM.valueInput(), 'bar');
    await blur(FORM.valueInput());

    const [event] = this.updateValue.lastCall.args;
    assert.strictEqual(event.target.value, 'bar', 'updateValue is called with blur event');
  });

  test('it clicks enable', async function (assert) {
    await this.renderComponent();
    await click(FORM.patchEdit());

    const [state] = this.updateState.lastCall.args;
    assert.strictEqual(state, 'enabled', 'updateState is called with "enabled"');
  });

  test('it clicks delete', async function (assert) {
    await this.renderComponent();
    await click(FORM.patchDelete());

    const [state] = this.updateState.lastCall.args;
    assert.strictEqual(state, 'deleted', 'updateState is called with "deleted"');
  });
});
