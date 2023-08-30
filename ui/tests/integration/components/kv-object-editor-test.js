/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

import { create } from 'ember-cli-page-object';
import kvObjectEditor from '../../pages/components/kv-object-editor';

import sinon from 'sinon';
const component = create(kvObjectEditor);

module('Integration | Component | kv-object-editor', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy();
  });

  test('it renders with no initial value', async function (assert) {
    await render(hbs`{{kv-object-editor onChange=this.spy}}`);
    assert.strictEqual(component.rows.length, 1, 'renders a single row');
    await component.addRow();
    assert.strictEqual(component.rows.length, 1, 'will only render row with a blank key');
  });

  test('it calls onChange when the val changes', async function (assert) {
    await render(hbs`{{kv-object-editor onChange=this.spy}}`);
    await component.rows.objectAt(0).kvKey('foo').kvVal('bar');
    assert.strictEqual(this.spy.callCount, 2, 'calls onChange each time change is triggered');
    assert.deepEqual(
      this.spy.lastCall.args[0],
      { foo: 'bar' },
      'calls onChange with the JSON representation of the data'
    );
    await component.addRow();
    await assert.strictEqual(component.rows.length, 2, 'adds a row when there is no blank one');
    await component.rows.objectAt(1).kvKey('another').kvVal('row');
    assert.propEqual(
      this.spy.lastCall.args[0],
      { foo: 'bar', another: 'row' },
      'calls onChange with second row of data'
    );
  });

  test('it renders passed data', async function (assert) {
    const metadata = { foo: 'bar', baz: 'bop' };
    this.set('value', metadata);
    await render(hbs`{{kv-object-editor value=this.value}}`);
    assert.strictEqual(
      component.rows.length,
      Object.keys(metadata).length + 1,
      'renders both rows of the metadata, plus an empty one'
    );
  });

  test('it deletes a row', async function (assert) {
    await render(hbs`{{kv-object-editor onChange=this.spy}}`);
    await component.rows.objectAt(0).kvKey('foo').kvVal('bar');
    await component.addRow();
    assert.strictEqual(component.rows.length, 2);
    assert.strictEqual(this.spy.callCount, 2, 'calls onChange for editing');
    await component.rows.objectAt(0).deleteRow();

    assert.strictEqual(component.rows.length, 1, 'only the blank row left');
    assert.strictEqual(this.spy.callCount, 3, 'calls onChange deleting row');
    assert.deepEqual(this.spy.lastCall.args[0], {}, 'last call to onChange is an empty object');
  });

  test('it shows a warning if there are duplicate keys', async function (assert) {
    const metadata = { foo: 'bar', baz: 'bop' };
    this.set('value', metadata);
    await render(hbs`{{kv-object-editor value=this.value onChange=this.spy}}`);
    await component.rows.objectAt(0).kvKey('foo');

    assert.ok(component.showsDuplicateError, 'duplicate keys are allowed but an error message is shown');
  });

  test('it supports custom placeholders', async function (assert) {
    await render(hbs`<KvObjectEditor @keyPlaceholder="foo" @valuePlaceholder="bar" />`);
    assert.dom('input').hasAttribute('placeholder', 'foo', 'Placeholder applied to key input');
    assert.dom('textarea').hasAttribute('placeholder', 'bar', 'Placeholder applied to value input');
  });

  test('it yields block in place of value input', async function (assert) {
    await render(
      hbs`
        <KvObjectEditor>
          <span data-test-yield></span>
        </KvObjectEditor>
      `
    );
    assert.dom('textarea').doesNotExist('Value input hidden when block is provided');
    assert.dom('[data-test-yield]').exists('Component yields block');
  });
});
