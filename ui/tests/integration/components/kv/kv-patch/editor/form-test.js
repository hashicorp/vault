/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { blur, click, fillIn, typeIn, render, focus } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';
import { NON_STRING_WARNING, WHITESPACE_WARNING } from 'vault/utils/model-helpers/validators';

module('Integration | Component | kv | kv-patch/editor/form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.subkeys = {
      foo: null,
      baz: {
        nested: null,
        bar: {
          hello: null,
        },
      },
    };
    this.onSubmit = sinon.spy();
    this.onCancel = sinon.spy();
    this.isSaving = false;
    this.submitError = '';

    this.renderComponent = async () => {
      return render(
        hbs`
    <KvPatch::Editor::Form
      @subkeys={{this.subkeys}}
      @onSubmit={{this.onSubmit}}
      @onCancel={{this.onCancel}}
      @isSaving={{this.isSaving}}
      @submitError={{this.submitError}}
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
    assert.true(this.onSubmit.calledOnce, 'clicking "Save" calls @onSubmit');
    await click(FORM.cancelBtn);
    assert.true(this.onCancel.calledOnce, 'clicking "Cancel" calls @onCancel');
  });

  test('it renders with no subkeys', async function (assert) {
    this.subkeys = {};
    await this.renderComponent();

    this.assertEmptyRow(assert);
  });

  test('it renders submit error from parent', async function (assert) {
    this.submitError = 'There was a problem submitting this form.';
    await this.renderComponent();
    assert.dom(GENERAL.inlineError).hasText(this.submitError);
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

    await click(FORM.patchAdd);
    assert
      .dom('[data-test-kv-key]')
      .exists({ count: 3 }, 'clicking add does not create a new row if key input is empty');

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
      await click(FORM.patchEdit());
      await fillIn(FORM.valueInput(), 'bar');
      // in qunit we have to unfocus the input so the following click event works on first try
      await blur(FORM.valueInput());

      // delete existing key
      await click(FORM.patchDelete(1));
      assert.dom(FORM.patchAlert('delete', 1)).hasText('This key value pair is marked for deletion.');
      assert.dom(FORM.keyInput(1)).hasClass('line-through');
      assert.dom(`${FORM.patchAlert('delete', 1)} ${GENERAL.icon('trash')}`).exists();
      // value is set to null under the hood, confirm the non-string warning doesn't display
      assert
        .dom(FORM.patchAlert('value-warning', 1))
        .doesNotExist('non-string warning does not render for null values');

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

      await click(FORM.patchEdit());
      await fillIn(FORM.valueInput(), 'bar');
      // in qunit we have to unfocus the input so the following click event works on the first try
      await blur(FORM.valueInput());
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

  module('it does not submit', function () {
    test('new keys that duplicate original subkeys', async function (assert) {
      await this.renderComponent();
      // patch existing key
      await click(FORM.patchEdit());
      await fillIn(FORM.valueInput(), 'bar');
      // add duplicate
      await fillIn(FORM.keyInput('new'), 'foo');
      await fillIn(FORM.valueInput('new'), 'duplicate');
      await click(FORM.saveBtn);

      assert
        .dom(GENERAL.inlineError)
        .hasText('This form contains validations errors, please resolve those before submitting.');
    });

    test('newly added keys edited to duplicate original subkeys', async function (assert) {
      await this.renderComponent();

      // patch existing key
      await click(FORM.patchEdit());
      await fillIn(FORM.valueInput(), 'bar');
      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);
      // go back and edit "aKey" to match pre-existing subkey 'foo'
      await fillIn(FORM.keyInput(2), 'foo');
      await click(FORM.saveBtn);

      assert
        .dom(GENERAL.inlineError)
        .hasText('This form contains validations errors, please resolve those before submitting.');
    });

    test('new keys that duplicate recently added keys', async function (assert) {
      await this.renderComponent();

      // create new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);
      // add same key name as above
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'duplicate');
      await click(FORM.saveBtn);

      assert
        .dom(GENERAL.inlineError)
        .hasText('This form contains validations errors, please resolve those before submitting.');
    });
  });

  module('duplicate keys error', function () {
    const validationMessage = (name) =>
      `"${name}" key already exists. Update the value of the existing key or rename this one.`;

    test('it renders for new keys that duplicate original subkeys', async function (assert) {
      await this.renderComponent();

      await fillIn(FORM.keyInput('new'), 'foo');
      await blur(FORM.keyInput('new')); // unfocus input to fire input change event and validation
      assert.dom(FORM.patchAlert('validation', 'new')).hasText(validationMessage('foo'));

      await click(FORM.patchAdd);
      assert
        .dom(FORM.keyInput('new'))
        .hasValue('foo', 'clicking "Add" is a noop, new row still has invalid value');

      await typeIn(FORM.keyInput('new'), '2'); // input value is now "foo2"
      await blur(FORM.keyInput('new')); // unfocus input
      assert
        .dom(FORM.patchAlert('validation', 'new'))
        .doesNotExist('error disappears when key no longer matches');

      await click(FORM.patchAdd);
      assert.dom(FORM.keyInput('new')).hasValue('', 'clicking "Add" creates a new row');
    });

    test('it renders for newly added keys edited to duplicate original subkeys', async function (assert) {
      // if a key is a duplicate then clicking "Add" does not work
      // this test asserts an error appears if a user goes back to rename a previously added key
      // to a duplicate and that it is not added to the payload
      await this.renderComponent();

      // add new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // add another
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');

      // go back and update "aKey" to match a pre-existing subkey
      await fillIn(FORM.keyInput(2), 'foo');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('validation', 2)).hasText(validationMessage('foo'));
      await typeIn(FORM.keyInput(2), '2');
      await blur(FORM.keyInput(2)); // unfocus input
      assert
        .dom(FORM.patchAlert('validation', 2))
        .doesNotExist('error disappears when key no longer matches');
    });

    test('it renders for new keys that duplicate recently added keys', async function (assert) {
      await this.renderComponent();

      // create new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // add same key name as above
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'bValue');
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('validation', 'new')).hasText(validationMessage('aKey'));

      await typeIn(FORM.keyInput('new'), '2');
      await blur(FORM.keyInput('new')); // unfocus input
      assert
        .dom(FORM.patchAlert('validation', 'new'))
        .doesNotExist('error disappears when key no longer matches');
    });

    test('it disappears after clicking "Remove" for duplicate', async function (assert) {
      await this.renderComponent();

      // create new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      // add same key name as above
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'bValue');
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('validation', 'new')).hasText(validationMessage('aKey'));

      await click(FORM.patchUndo(2));
      assert.dom(FORM.patchAlert('validation', 'new')).doesNotExist('error clears when duplicate is removed');
      await click(FORM.patchAdd);
      // assert a new row is added
      this.assertEmptyRow(assert);
    });

    test('it disappears after clicking "Remove" for multiple duplicates', async function (assert) {
      await this.renderComponent();

      // create new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);

      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');
      await click(FORM.patchAdd);

      // and add another duplicate
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');

      // edit key to be same as first
      await fillIn(FORM.keyInput(3), 'aKey');

      // remove all but latest key
      await click(FORM.patchUndo(3));
      await click(FORM.patchUndo(2));
      await click(FORM.patchAdd);

      // assert a new row is added
      this.assertEmptyRow(assert);
    });

    test('it does not render when refocusing a previously inputted key', async function (assert) {
      await this.renderComponent();

      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);
      // focus and unfocus
      await focus(FORM.keyInput(2));
      await blur(FORM.keyInput(2));
      assert.dom(FORM.patchAlert('validation', 2)).doesNotExist();
    });

    // accounts for an edge case where not setting invalid key values caused
    // error to show for outdated keys and then updating the key did not remove error
    test('it disappears for new key when another duplicate key is edited', async function (assert) {
      await this.renderComponent();

      // create new key and click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);
      // edit key to be a duplicate of original subkey, "foo"
      // since "foo" is invalid it does not update the tracked KV.key value.
      // the input value reads "foo" but underlying KV class key value is still "aKey"
      await fillIn(FORM.keyInput(2), 'foo');
      // fill in new key that matches underlying value of input above
      await fillIn(FORM.keyInput('new'), 'aKey');
      // validation errors now show for both inputs even though no visible input reads "aKey" (while strange UX, it's a super edge case)
      // editing input at index 2 ("foo") should make both disappear
      await fillIn(FORM.keyInput(2), 'foo2');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('validation', 2)).doesNotExist();
      assert.dom(FORM.patchAlert('validation', 'new')).doesNotExist();
    });
  });

  module('it shows whitespace warning', function () {
    test('for new keys with whitespace', async function (assert) {
      await this.renderComponent();

      await fillIn(FORM.keyInput('new'), 'a space');
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('key-warning', 'new')).hasText(WHITESPACE_WARNING('this key'));

      await fillIn(FORM.keyInput('new'), 'nospace');
      await blur(FORM.keyInput('new')); // unfocus input
      assert.dom(FORM.patchAlert('key-warning', 'new')).doesNotExist('warning disappears when key updates');
    });

    test('for newly added keys edited to have whitespace', async function (assert) {
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
      assert.dom(FORM.patchAlert('key-warning', 2)).hasText(WHITESPACE_WARNING('this key'));

      await fillIn(FORM.keyInput(2), 'aKey');
      await blur(FORM.keyInput(2)); // unfocus input
      assert.dom(FORM.patchAlert('key-warning', 2)).doesNotExist('warning disappears when key updates');
    });

    test('for keys with whitespace after clicking "Add"', async function (assert) {
      await this.renderComponent();

      // add new key with space and click add
      await fillIn(FORM.keyInput('new'), 'a key');
      await click(FORM.patchAdd);
      assert
        .dom(FORM.patchAlert('key-warning', 2))
        .hasText(WHITESPACE_WARNING('this key'), 'warning is attached to relevant key');
      assert
        .dom(FORM.patchAlert('key-warning', 'new'))
        .doesNotExist('there is no whitespace warning for the new empty row');

      // add another
      await fillIn(FORM.keyInput('new'), 'b key');
      await blur(FORM.keyInput('new')); // unfocus input
      assert
        .dom(FORM.patchAlert('key-warning', 2))
        .hasText(WHITESPACE_WARNING('this key'), 'warning is still attached to relevant key');
      assert
        .dom(FORM.patchAlert('key-warning', 'new'))
        .hasText(WHITESPACE_WARNING('this key'), 'new key also has whitespace warning');
    });
  });

  module('it shows non-string warning', function () {
    const NON_STRING_VALUES = [0, 123, '{ "a": "b" }', 'null'];

    NON_STRING_VALUES.forEach((value) => {
      test(`for new non-string values: ${value}`, async function (assert) {
        await this.renderComponent();
        await fillIn(FORM.keyInput('new'), 'aKey');
        await fillIn(FORM.valueInput('new'), value);
        await blur(FORM.valueInput('new')); // unfocus input
        assert.dom(FORM.patchAlert('value-warning', 'new')).hasText(NON_STRING_WARNING);

        await typeIn(FORM.valueInput('new'), 'abc');
        await blur(FORM.valueInput('new')); // unfocus input
        assert
          .dom(FORM.patchAlert('value-warning', 'new'))
          .doesNotExist(`warning disappears when ${value} includes a non-parsable string`);
      });
    });

    NON_STRING_VALUES.forEach((value) => {
      test(`for newly added values edited to non-string values: ${value}`, async function (assert) {
        await this.renderComponent();

        // add new key and click add
        await fillIn(FORM.keyInput('new'), 'aKey');
        await fillIn(FORM.valueInput('new'), 'aValue');
        await click(FORM.patchAdd);

        // add another
        await fillIn(FORM.keyInput('new'), 'bKey');
        await fillIn(FORM.valueInput('new'), 'bValue');

        // go back and change "aKey" to have a non-string
        await fillIn(FORM.valueInput(2), value);
        await blur(FORM.valueInput(2)); // unfocus input
        assert.dom(FORM.patchAlert('value-warning', 2)).hasText(NON_STRING_WARNING);

        await fillIn(FORM.valueInput(2), 'abc');
        await blur(FORM.valueInput(2)); // unfocus input
        assert
          .dom(FORM.patchAlert('value-warning', 2))
          .doesNotExist(`warning disappears when ${value} is replaced with a string`);
      });
    });

    NON_STRING_VALUES.forEach((value) => {
      test(`for non-string values after clicking "Add": ${value}`, async function (assert) {
        await this.renderComponent();

        // add non-string value and click add
        await fillIn(FORM.keyInput('new'), 'aKey');
        await fillIn(FORM.valueInput('new'), value);
        await click(FORM.patchAdd);
        assert
          .dom(FORM.patchAlert('value-warning', 2))
          .hasText(NON_STRING_WARNING, 'warning is attached to relevant row');
        assert
          .dom(FORM.patchAlert('value-warning', 'new'))
          .doesNotExist('there is no non-string warning for the new empty row');

        // add another
        await fillIn(FORM.keyInput('new'), 'bKey');
        await fillIn(FORM.valueInput('new'), value);
        await blur(FORM.valueInput('new')); // unfocus input
        assert
          .dom(FORM.patchAlert('value-warning', 2))
          .hasText(NON_STRING_WARNING, 'warning is still attached to relevant row');
        assert
          .dom(FORM.patchAlert('value-warning', 'new'))
          .hasText(NON_STRING_WARNING, 'new row also has non-string warning');
      });
    });
  });
});
