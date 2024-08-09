/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { blur, click, fillIn, typeIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv | kv-patch-editor', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.subkeyArray = ['foo', 'baz'];
    this.onSubmit = sinon.spy();
    this.onCancel = sinon.spy();
    this.isSaving = false;

    this.renderComponent = async () => {
      return render(
        hbs`
    <KvPatchEditor
      @subkeyArray={{this.subkeyArray}}
      @onSubmit={{this.onSubmit}}
      @onCancel={{this.onCancel}}
      @isSaving={{this.isSaving}}
    />`,
        { owner: this.engine }
      );
    };

    // HELPERS
    this.assertDefaultRow = (key, assert) => {
      assert.dom(FORM.keyInput(key)).hasValue(key);
      assert.dom(FORM.keyInput(key)).isDisabled();
      assert.dom(FORM.valueInput(key)).hasValue('');
      assert.dom(FORM.valueInput(key)).isDisabled();
      assert.dom(FORM.patchEdit(key)).exists();
      assert.dom(FORM.patchDelete(key)).exists();
    };

    this.assertEmptyRow = (assert) => {
      assert.dom(FORM.keyInput('new')).hasValue('');
      assert.dom(FORM.keyInput('new')).isNotDisabled();
      assert.dom(FORM.keyInput('new')).hasAttribute('placeholder', 'key');
      assert.dom(FORM.valueInput('new')).hasValue('');
      assert.dom(FORM.valueInput('new')).isNotDisabled();
      assert.dom(FORM.patchAdd).exists({ count: 1 });
      assert.dom(FORM.patchAdd).isDisabled();
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();

    this.assertDefaultRow('foo', assert);
    this.assertDefaultRow('baz', assert);
    this.assertEmptyRow(assert);

    await click(FORM.saveBtn);
    assert.true(this.onSubmit.calledOnce);
    await click(FORM.cancelBtn);
    assert.true(this.onCancel.calledOnce);
  });

  test('it enables and disables inputs', async function (assert) {
    await this.renderComponent();

    const enableAndAssert = async (key) => {
      await click(FORM.patchEdit(key));
      assert.dom(FORM.valueInput(key)).isEnabled('clicking edit enables value input');
      assert.dom(FORM.keyInput(key)).isDisabled(`${key} key input is still disabled`);
      assert.dom(FORM.patchEdit(key)).doesNotExist('edit button disappears');
      assert.dom(FORM.patchDelete(key)).doesNotExist('delete button disappears');
      assert
        .dom(FORM.patchUndo(key))
        .hasText('Cancel', 'Undo button reads "Cancel" and replaces edit and delete');
    };

    await enableAndAssert('foo');
    await click(FORM.patchUndo('foo'));
    this.assertDefaultRow('foo', assert);

    await enableAndAssert('baz');
    await click(FORM.patchUndo('baz'));
    this.assertDefaultRow('baz', assert);
  });

  test('it adds a new row', async function (assert) {
    await this.renderComponent();

    await fillIn(FORM.keyInput('new'), 'a');
    assert.dom(FORM.patchAdd).isDisabled(); // only enables when both key and value exist
    await fillIn(FORM.valueInput('new'), 'b');
    await click(FORM.patchAdd);

    assert.dom(FORM.keyInput('a')).hasValue('a');
    assert.dom(FORM.keyInput('a')).isEnabled('new key inputs are enabled');
    assert.dom(FORM.patchUndo('a')).hasText('Remove', 'Undo button reads "Remove" for new keys');

    // assert a new row is added
    this.assertEmptyRow(assert);
  });

  test('it renders loading state', async function (assert) {
    this.isSaving = true;
    await this.renderComponent();

    assert.dom(FORM.saveBtn).isDisabled();
    assert.dom(FORM.cancelBtn).isDisabled();
    assert.dom(`${FORM.saveBtn} ${GENERAL.icon('loading')}`).exists();
  });

  module('it submits', function () {
    test('patch data for existing, deleted and new keys', async function (assert) {
      await this.renderComponent();

      // patch existing key
      await click(FORM.patchEdit('foo'));
      await fillIn(FORM.valueInput('foo'), 'bar');
      await blur(FORM.valueInput('foo')); // unfocus input so click event below fires

      // delete existing key
      await click(FORM.patchDelete('baz'));
      assert.dom(FORM.patchAlert('delete', 'baz')).hasText('This key value pair is marked for deletion.');
      assert.dom(`${FORM.patchAlert('delete', 'baz')} ${GENERAL.icon('alert-diamond-fill')}`).exists();

      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'a');
      await fillIn(FORM.valueInput('new'), 'b');
      await click(FORM.patchAdd);

      // add new key and do NOT click add
      await fillIn(FORM.keyInput('new'), 'c');
      await fillIn(FORM.valueInput('new'), 'd');

      await click(FORM.saveBtn);

      const [data] = this.onSubmit.lastCall.args;
      assert.propEqual(
        data,
        { baz: null, foo: 'bar', a: 'b', c: 'd' },
        `onSubmit called with ${JSON.stringify(data)}`
      );
    });

    test('patch data when every action is cancelled', async function (assert) {
      await this.renderComponent();

      await click(FORM.patchEdit('foo'));
      await fillIn(FORM.valueInput('foo'), 'bar');
      await blur(FORM.valueInput('foo')); // unfocus input so click event below fires

      await click(FORM.patchDelete('baz'));

      await fillIn(FORM.keyInput('new'), 'a');
      await fillIn(FORM.valueInput('new'), 'b');
      await click(FORM.patchAdd);

      // undo every action
      await click(FORM.patchUndo('foo'));
      await click(FORM.patchUndo('baz'));
      await click(FORM.patchUndo('a'));
      await click(FORM.saveBtn);

      const [data] = this.onSubmit.lastCall.args;
      assert.propEqual(data, {}, `onSubmit called with ${JSON.stringify(data)}`);
    });
  });

  module('it validates', function () {
    const validationMessage =
      '"foo" key already exists. Patch the value of the existing key or rename this one.';

    test('duplicate keys', async function (assert) {
      await this.renderComponent();

      await typeIn(FORM.keyInput('new'), 'foo');
      await fillIn(FORM.valueInput('new'), 'value');
      assert.dom(FORM.patchAlert('validation', 'new')).hasText(validationMessage);
      assert.dom(FORM.patchAdd).isDisabled('add button disables for invalid inputs');
      await typeIn(FORM.keyInput('new'), '2');
      assert.dom(FORM.patchAlert('validation', 'new')).doesNotExist('warning disappears when key updates');
      assert.dom(FORM.patchAdd).isEnabled('add button is enabled when input is fixed');
    });

    test('duplicate keys after they have been added', async function (assert) {
      // check for new key and for recently added key
      // check for changing newly added key to duplicate
      await this.renderComponent();

      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'a');
      await fillIn(FORM.valueInput('new'), 'b');
      await blur(FORM.valueInput('new')); // unfocus input so click event below fires
      await click(FORM.patchAdd);

      // add another
      await fillIn(FORM.keyInput('new'), 'c');
      await fillIn(FORM.valueInput('new'), 'd');

      // update previous key to match a subkey
      await fillIn(FORM.keyInput('a'), ''); // clear input
      await typeIn(FORM.keyInput('a'), 'foo'); // clear input
      assert.dom(FORM.patchAlert('validation', 2)).hasText(validationMessage);
      await typeIn(FORM.keyInput('new'), '2');
      assert.dom(FORM.patchAlert('validation', 'new')).doesNotExist('warning disappears when key updates');
    });
  });
});
