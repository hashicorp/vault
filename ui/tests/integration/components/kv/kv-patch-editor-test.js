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
    this.assertDefaultRow = (idx, key, assert) => {
      assert.dom(FORM.keyInput(idx)).hasValue(key);
      assert.dom(FORM.keyInput(idx)).isDisabled();
      assert.dom(FORM.valueInput(idx)).hasValue('');
      assert.dom(FORM.valueInput(idx)).isDisabled();
      assert.dom(FORM.patchEdit(idx)).exists();
      assert.dom(FORM.patchDelete(idx)).exists();
    };

    this.assertEmptyRow = (assert) => {
      assert.dom(FORM.keyInput('new')).hasValue('');
      assert.dom(FORM.keyInput('new')).isNotDisabled();
      assert.dom(FORM.keyInput('new')).hasAttribute('placeholder', 'key');
      assert.dom(FORM.valueInput('new')).hasValue('');
      assert.dom(FORM.valueInput('new')).isNotDisabled();
      assert.dom(FORM.patchAdd).exists({ count: 1 });
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();

    this.assertDefaultRow(0, 'foo', assert);
    this.assertDefaultRow(1, 'baz', assert);
    this.assertEmptyRow(assert);

    await click(FORM.saveBtn);
    assert.true(this.onSubmit.calledOnce);
    await click(FORM.cancelBtn);
    assert.true(this.onCancel.calledOnce);
  });

  test('it enables and disables inputs', async function (assert) {
    await this.renderComponent();

    const enableAndAssert = async (idx, key) => {
      await click(FORM.patchEdit(idx));
      assert.dom(FORM.valueInput(idx)).isEnabled('clicking edit enables value input');
      assert.dom(FORM.keyInput(idx)).hasAttribute('readonly', '', `${key} input updates to readonly`);
      assert.dom(FORM.patchEdit(idx)).doesNotExist('edit button disappears');
      assert.dom(FORM.patchDelete(idx)).doesNotExist('delete button disappears');
      assert
        .dom(FORM.patchUndo(idx))
        .hasText('Cancel', 'Undo button reads "Cancel" and replaces edit and delete');
    };

    await enableAndAssert(0, 'foo');
    await click(FORM.patchUndo(0));
    this.assertDefaultRow(0, 'foo', assert);

    await enableAndAssert(1, 'baz');
    await click(FORM.patchUndo(1));
    this.assertDefaultRow(1, 'baz', assert);
  });

  test('it adds a new row', async function (assert) {
    await this.renderComponent();

    await fillIn(FORM.keyInput('new'), 'aKey');
    await fillIn(FORM.valueInput('new'), 'aValue');
    await click(FORM.patchAdd);

    assert.dom(FORM.keyInput(2)).hasValue('aKey');
    assert.dom(FORM.keyInput(2)).isEnabled('new key inputs are enabled');
    assert.dom(FORM.valueInput(2)).isEnabled('new value inputs are enabled');
    assert.dom(FORM.patchUndo(2)).hasText('Remove', 'Undo button reads "Remove" for new keys');

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
      await click(FORM.patchEdit(0));
      await fillIn(FORM.valueInput(0), 'bar');
      // in qunit we have to unfocus the input so the following click event works on first try
      await blur(FORM.valueInput(0));

      // delete existing key
      await click(FORM.patchDelete(1));
      assert.dom(FORM.patchAlert('delete', 1)).hasText('This key value pair is marked for deletion.');
      assert.dom(`${FORM.patchAlert('delete', 1)} ${GENERAL.icon('trash')}`).exists();

      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // add new key and do NOT click add
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');

      await click(FORM.saveBtn);

      const [data] = this.onSubmit.lastCall.args;
      assert.propEqual(
        data,
        { baz: null, foo: 'bar', aKey: 'aValue', bKey: 'bValue' },
        `onSubmit called with ${JSON.stringify(data)}`
      );
    });

    test('patch data when every action is canceled', async function (assert) {
      await this.renderComponent();

      await click(FORM.patchEdit(0));
      await fillIn(FORM.valueInput(0), 'bar');
      // in qunit we have to unfocus the input so the following click event works on the first try
      await blur(FORM.valueInput(0));
      await click(FORM.patchDelete(1));

      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // undo every action
      await click(FORM.patchUndo(0)); // undo edit
      await click(FORM.patchUndo(1)); // undo delete
      await click(FORM.patchUndo(2)); // remove new row
      await click(FORM.saveBtn);

      const [data] = this.onSubmit.lastCall.args;
      assert.propEqual(data, {}, `onSubmit called with ${JSON.stringify(data)}`);
    });
  });

  module('it validates', function () {
    const validationMessage =
      '"foo" key already exists. Update the value of the existing key or rename this one.';
    const whitespaceWarning =
      "This key contains whitespace. If this is desired, you'll need to encode it with %20 in APi requests.";

    test('new duplicate keys', async function (assert) {
      await this.renderComponent();

      await fillIn(FORM.keyInput('new'), 'foo');
      await blur(FORM.keyInput('new')); // unfocus input to fire input change event and validation
      assert.dom(FORM.patchAlert('validation', 'new')).hasText(validationMessage);

      await click(FORM.patchAdd);
      assert
        .dom(FORM.keyInput('new'))
        .hasValue('foo', 'clicking "Add" is a noop, new row still has invalid value');

      await typeIn(FORM.keyInput('new'), '2'); // input value is now "foo2"
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('validation', 'new')).doesNotExist('error disappears when key updates');

      await click(FORM.patchAdd);
      assert.dom(FORM.keyInput('new')).hasValue('', 'clicking "Add" creates a new row');
    });

    test('existing duplicate keys', async function (assert) {
      await this.renderComponent();

      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // add another
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');

      // go back and update "aKey" to match an existing subkey
      await fillIn(FORM.keyInput(2), 'foo');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('validation', 2)).hasText(validationMessage);

      await typeIn(FORM.keyInput(2), '2');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('validation', '2')).doesNotExist('error disappears when key updates');
    });

    test('new keys with white space', async function (assert) {
      await this.renderComponent();

      await fillIn(FORM.keyInput('new'), 'a space');
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('warning', 'new')).hasText(whitespaceWarning);

      await fillIn(FORM.keyInput('new'), 'nospace');
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('warning', 'new')).doesNotExist('warning disappears when key updates');
    });

    test('existing keys with whitespace', async function (assert) {
      await this.renderComponent();

      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // add another
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');

      // go back and change "aKey" to have a space
      await fillIn(FORM.keyInput(2), 'a key');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('warning', 2)).hasText(whitespaceWarning);

      await fillIn(FORM.keyInput(2), 'aKey');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('validation', '2')).doesNotExist('warning disappears when key updates');
    });

    test('keys with whitespace after clicking "Add"', async function (assert) {
      await this.renderComponent();

      // add new key with space and click add
      await fillIn(FORM.keyInput('new'), 'a key');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);
      assert
        .dom(FORM.patchAlert('warning', 2))
        .hasText(whitespaceWarning, 'warning is attached to relevant key');
      assert
        .dom(FORM.patchAlert('warning', 'new'))
        .doesNotExist('there is no whitespace warning for the empty row');

      // add another
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');

      // update "aKey" to have a space
      await fillIn(FORM.keyInput(2), 'a key');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('warning', 2)).hasText(whitespaceWarning);

      await fillIn(FORM.keyInput(2), 'aKey');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('validation', '2')).doesNotExist('warning disappears when key updates');
    });
  });
});
