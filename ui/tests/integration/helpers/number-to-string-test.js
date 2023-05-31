import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { numberToString } from 'core/helpers/number-to-string';

module('Integration | Helper | number-to-string', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.inputValue = 1234;
    await render(hbs`{{number-to-string this.inputValue}}`);
    assert.dom(this.element).hasText('1234');

    this.inputValue = '4567';
    await render(hbs`{{number-to-string this.inputValue}}`);
    assert.dom(this.element).hasText('4567');
  });

  test('it converts numbers to string type with no options', async function (assert) {
    assert.strictEqual(numberToString([0], {}), '0');
    assert.strictEqual(numberToString([123], {}), '123');
    assert.strictEqual(numberToString(['1,234'], {}), '1,234');
    assert.strictEqual(numberToString(['456'], {}), '456', 'it returns non-integer values as-is');
    assert.strictEqual(numberToString(['0'], {}), '0', 'it returns string 0 as-is');
    assert.strictEqual(numberToString(['abc'], {}), 'abc', 'it returns string of characters as-is');
  });

  test('it converts to string only if options are strings', async function (assert) {
    assert.strictEqual(numberToString([0], { options: ['0', '1', '2'] }), '0');
    assert.strictEqual(numberToString([123], { options: ['123', '456'] }), '123');
    assert.strictEqual(numberToString([0], { options: [0] }), 0);
    assert.strictEqual(numberToString([123], { options: [123] }), 123);
    assert.strictEqual(numberToString(['456'], { options: ['456'] }), '456');
    assert.strictEqual(numberToString(['0'], { options: ['0'] }), '0');
    assert.strictEqual(numberToString(['abc'], { options: ['abc'] }), 'abc');
  });
});
