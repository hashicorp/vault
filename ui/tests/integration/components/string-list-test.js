import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, triggerKeyEvent, triggerEvent } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | string list', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy();
  });

  const assertBlank = function (assert) {
    assert.dom('[data-test-string-list-input]').exists({ count: 1 }, 'renders 1 input');
    assert.dom('[data-test-string-list-input]').hasValue('', 'the input is blank');
  };

  const assertFoo = function (assert) {
    assert.dom('[data-test-string-list-input]').exists({ count: 2 }, 'renders 2 inputs');
    assert.dom('[data-test-string-list-input="0"]').hasValue('foo', 'first input has the inputValue');
    assert.dom('[data-test-string-list-input="1"]').hasValue('', 'second input is blank');
  };

  const assertFooBar = function (assert) {
    assert.dom('[data-test-string-list-input]').exists({ count: 3 }, 'renders 3 inputs');
    assert.dom('[data-test-string-list-input="0"]').hasValue('foo');
    assert.dom('[data-test-string-list-input="1"]').hasValue('bar');
    assert.dom('[data-test-string-list-input="2"]').hasValue('', 'last input is blank');
  };

  test('it renders the label', async function (assert) {
    assert.expect(4);
    await render(hbs`<StringList @label="foo" @onChange={{this.spy}} />`);
    assert.dom('[data-test-string-list-label]').hasText('foo', 'renders the label when provided');
    await render(hbs`<StringList />`);
    assert.dom('[data-test-string-list-label]').doesNotExist('does not render the label');
    assertBlank(assert);
  });

  test('it renders inputValue from empty string', async function (assert) {
    assert.expect(2);
    await render(hbs`<StringList @inputValue="" />`);
    assertBlank(assert);
  });

  test('it renders inputValue from string with one value', async function (assert) {
    assert.expect(3);
    await render(hbs`<StringList @inputValue="foo" />`);
    assertFoo(assert);
  });

  test('it renders inputValue from comma-separated string', async function (assert) {
    assert.expect(4);
    await render(hbs`<StringList @inputValue="foo,bar" />`);
    assertFooBar(assert);
  });

  test('it renders inputValue from a blank array', async function (assert) {
    assert.expect(2);
    this.set('inputValue', []);
    await render(hbs`<StringList @inputValue={{this.inputValue}} />`);
    assertBlank(assert);
  });

  test('it renders inputValue array with a single item', async function (assert) {
    assert.expect(3);
    this.set('inputValue', ['foo']);
    await render(hbs`<StringList @inputValue={{this.inputValue}} />`);
    assertFoo(assert);
  });

  test('it renders inputValue array with a multiple items', async function (assert) {
    assert.expect(4);
    this.set('inputValue', ['foo', 'bar']);
    await render(hbs`<StringList @inputValue={{this.inputValue}} />`);
    assertFooBar(assert);
  });

  test('it adds a new row only when the last row is not blank', async function (assert) {
    assert.expect(5);
    await render(hbs`<StringList @inputValue="" @onChange={{this.spy}}/>`);
    await click('[data-test-string-list-button="add"]');
    assertBlank(assert);
    await fillIn('[data-test-string-list-input="0"]', 'foo');
    await triggerKeyEvent('[data-test-string-list-input="0"]', 'keyup', 14);
    await click('[data-test-string-list-button="add"]');
    assertFoo(assert);
  });

  test('it trims input values', async function (assert) {
    await render(hbs`<StringList @inputValue="" @onChange={{this.spy}}/>`);
    await fillIn('[data-test-string-list-input="0"]', 'foo');
    await triggerKeyEvent('[data-test-string-list-input="0"]', 'keyup', 14);
    assert.dom('[data-test-string-list-input="0"]').hasValue('foo');
  });

  test('it calls onChange with array when editing', async function (assert) {
    assert.expect(2);
    this.set('inputValue', ['foo']);
    this.set('onChange', function (val) {
      assert.deepEqual(val, ['foo', 'bar'], 'calls onChange with expected value');
    });
    await render(hbs`<StringList @inputValue={{this.inputValue}} @onChange={{this.onChange}} />`);
    await fillIn('[data-test-string-list-input="1"]', 'bar');
    await triggerKeyEvent('[data-test-string-list-input="1"]', 'keyup', 14);
  });

  test('it calls onChange with string when editing', async function (assert) {
    assert.expect(2);
    this.set('inputValue', 'foo');
    this.set('onChange', function (val) {
      assert.strictEqual(val, 'foo,bar', 'calls onChange with expected value');
    });
    await render(
      hbs`<StringList @type="string" @inputValue={{this.inputValue}} @onChange={{this.onChange}}/>`
    );
    await fillIn('[data-test-string-list-input="1"]', 'bar');
    await triggerKeyEvent('[data-test-string-list-input="1"]', 'keyup', 14);
  });

  test('it removes a row', async function (assert) {
    assert.expect(4);
    this.set('inputValue', ['foo', 'bar']);
    this.set('onChange', function (val) {
      assert.deepEqual(val, ['bar'], 'calls onChange with expected value');
    });
    await render(hbs`<StringList @inputValue={{this.inputValue}} @onChange={{this.onChange}} />`);

    await click('[data-test-string-list-row="0"] [data-test-string-list-button="delete"]');
    assert.dom('[data-test-string-list-input]').exists({ count: 2 }, 'renders 2 inputs');
    assert.dom('[data-test-string-list-input="0"]').hasValue('bar');
    assert.dom('[data-test-string-list-input="1"]').hasValue('');
  });

  test('it replaces helpText if name is tokenBoundCidrs', async function (assert) {
    assert.expect(1);
    await render(hbs`<StringList @label={{'blah'}} @helpText={{'blah'}} @attrName={{'tokenBoundCidrs'}} />`);
    const tooltipTrigger = document.querySelector('[data-test-tool-tip-trigger]');
    await triggerEvent(tooltipTrigger, 'mouseenter');
    assert
      .dom('[data-test-info-tooltip-content]')
      .hasText(
        'Specifies the blocks of IP addresses which are allowed to use the generated token. One entry per row.'
      );
  });
});
