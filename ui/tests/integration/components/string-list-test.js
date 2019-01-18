import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, find, findAll, fillIn, triggerKeyEvent } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | string list', function(hooks) {
  setupRenderingTest(hooks);

  const assertBlank = function(assert) {
    assert.equal(findAll('[data-test-string-list-input]').length, 1, 'renders 1 input');
    assert.equal(find('[data-test-string-list-input]').value, '', 'the input is blank');
  };

  const assertFoo = function(assert) {
    assert.equal(findAll('[data-test-string-list-input]').length, 2, 'renders 2 inputs');
    assert.equal(find('[data-test-string-list-input="0"]').value, 'foo', 'first input has the inputValue');
    assert.equal(find('[data-test-string-list-input="1"]').value, '', 'second input is blank');
  };

  const assertFooBar = function(assert) {
    assert.equal(findAll('[data-test-string-list-input]').length, 3, 'renders 3 inputs');
    assert.equal(find('[data-test-string-list-input="0"]').value, 'foo');
    assert.equal(find('[data-test-string-list-input="1"]').value, 'bar');
    assert.equal(find('[data-test-string-list-input="2"]').value, '', 'last input is blank');
  };

  test('it renders the label', async function(assert) {
    await render(hbs`{{string-list label="foo"}}`);
    assert.equal(
      find('[data-test-string-list-label]').textContent.trim(),
      'foo',
      'renders the label when provided'
    );

    await render(hbs`{{string-list}}`);
    assert.equal(findAll('[data-test-string-list-label]').length, 0, 'does not render the label');
    assertBlank(assert);
  });

  test('it renders inputValue from empty string', async function(assert) {
    await render(hbs`{{string-list inputValue=""}}`);
    assertBlank(assert);
  });

  test('it renders inputValue from string with one value', async function(assert) {
    await render(hbs`{{string-list inputValue="foo"}}`);
    assertFoo(assert);
  });

  test('it renders inputValue from comma-separated string', async function(assert) {
    await render(hbs`{{string-list inputValue="foo,bar"}}`);
    assertFooBar(assert);
  });

  test('it renders inputValue from a blank array', async function(assert) {
    this.set('inputValue', []);
    await render(hbs`{{string-list inputValue=inputValue}}`);
    assertBlank(assert);
  });

  test('it renders inputValue array with a single item', async function(assert) {
    this.set('inputValue', ['foo']);
    await render(hbs`{{string-list inputValue=inputValue}}`);
    assertFoo(assert);
  });

  test('it renders inputValue array with a multiple items', async function(assert) {
    this.set('inputValue', ['foo', 'bar']);
    await render(hbs`{{string-list inputValue=inputValue}}`);
    assertFooBar(assert);
  });

  test('it adds a new row only when the last row is not blank', async function(assert) {
    await render(hbs`{{string-list inputValue=""}}`);
    await click('[data-test-string-list-button="add"]');
    assertBlank(assert);
    await fillIn('[data-test-string-list-input="0"]', 'foo');
    await triggerKeyEvent('[data-test-string-list-input="0"]', 'keyup', 14);
    await click('[data-test-string-list-button="add"]');
    assertFoo(assert);
  });

  test('it trims input values', async function(assert) {
    await render(hbs`{{string-list inputValue=""}}`);
    await fillIn('[data-test-string-list-input="0"]', 'foo');
    await triggerKeyEvent('[data-test-string-list-input="0"]', 'keyup', 14);
    assert.equal(find('[data-test-string-list-input="0"]').value, 'foo');
  });

  test('it calls onChange with array when editing', async function(assert) {
    this.set('inputValue', ['foo']);
    this.set('onChange', function(val) {
      assert.deepEqual(val, ['foo', 'bar'], 'calls onChange with expected value');
    });
    await render(hbs`{{string-list inputValue=inputValue onChange=(action onChange)}}`);
    await fillIn('[data-test-string-list-input="1"]', 'bar');
    await triggerKeyEvent('[data-test-string-list-input="1"]', 'keyup', 14);
  });

  test('it calls onChange with string when editing', async function(assert) {
    this.set('inputValue', 'foo');
    this.set('onChange', function(val) {
      assert.equal(val, 'foo,bar', 'calls onChange with expected value');
    });
    await render(hbs`{{string-list inputValue=inputValue onChange=(action onChange)}}`);
    await fillIn('[data-test-string-list-input="1"]', 'bar');
    await triggerKeyEvent('[data-test-string-list-input="1"]', 'keyup', 14);
  });

  test('it removes a row', async function(assert) {
    this.set('inputValue', ['foo', 'bar']);
    this.set('onChange', function(val) {
      assert.equal(val, 'bar', 'calls onChange with expected value');
    });
    await render(hbs`{{string-list inputValue=inputValue onChange=(action onChange)}}`);

    await click('[data-test-string-list-row="0"] [data-test-string-list-button="delete"]');
    assert.equal(findAll('[data-test-string-list-input]').length, 2, 'renders 2 inputs');
    assert.equal(find('[data-test-string-list-input="0"]').value, 'bar');
    assert.equal(find('[data-test-string-list-input="1"]').value, '');
  });
});
