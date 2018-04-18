import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('string-list', 'Integration | Component | string list', {
  integration: true,
});

const assertBlank = function(assert) {
  assert.equal(this.$('[data-test-string-list-input]').length, 1, 'renders 1 input');
  assert.equal(this.$('[data-test-string-list-input]').val(), '', 'the input is blank');
};

const assertFoo = function(assert) {
  assert.equal(this.$('[data-test-string-list-input]').length, 2, 'renders 2 inputs');
  assert.equal(this.$('[data-test-string-list-input="0"]').val(), 'foo', 'first input has the inputValue');
  assert.equal(this.$('[data-test-string-list-input="1"]').val(), '', 'second input is blank');
};

const assertFooBar = function(assert) {
  assert.equal(this.$('[data-test-string-list-input]').length, 3, 'renders 3 inputs');
  assert.equal(this.$('[data-test-string-list-input="0"]').val(), 'foo');
  assert.equal(this.$('[data-test-string-list-input="1"]').val(), 'bar');
  assert.equal(this.$('[data-test-string-list-input="2"]').val(), '', 'last input is blank');
};

test('it renders the label', function(assert) {
  this.render(hbs`{{string-list label="foo"}}`);
  assert.equal(
    this.$('[data-test-string-list-label]').text().trim(),
    'foo',
    'renders the label when provided'
  );

  this.render(hbs`{{string-list}}`);
  assert.equal(this.$('[data-test-string-list-label]').length, 0, 'does not render the label');
  assertBlank.call(this, assert);
});

test('it renders inputValue from empty string', function(assert) {
  this.render(hbs`{{string-list inputValue=""}}`);
  assertBlank.call(this, assert);
});

test('it renders inputValue from string with one value', function(assert) {
  this.render(hbs`{{string-list inputValue="foo"}}`);
  assertFoo.call(this, assert);
});

test('it renders inputValue from comma-separated string', function(assert) {
  this.render(hbs`{{string-list inputValue="foo,bar"}}`);
  assertFooBar.call(this, assert);
});

test('it renders inputValue from a blank array', function(assert) {
  this.set('inputValue', []);
  this.render(hbs`{{string-list inputValue=inputValue}}`);
  assertBlank.call(this, assert);
});

test('it renders inputValue array with a single item', function(assert) {
  this.set('inputValue', ['foo']);
  this.render(hbs`{{string-list inputValue=inputValue}}`);
  assertFoo.call(this, assert);
});

test('it renders inputValue array with a multiple items', function(assert) {
  this.set('inputValue', ['foo', 'bar']);
  this.render(hbs`{{string-list inputValue=inputValue}}`);
  assertFooBar.call(this, assert);
});

test('it adds a new row only when the last row is not blank', function(assert) {
  this.render(hbs`{{string-list inputValue=""}}`);
  this.$('[data-test-string-list-button="add"]').click();
  assertBlank.call(this, assert);
  this.$('[data-test-string-list-input="0"]').val('foo').keyup();
  this.$('[data-test-string-list-button="add"]').click();
  assertFoo.call(this, assert);
});

test('it trims input values', function(assert) {
  this.render(hbs`{{string-list inputValue=""}}`);
  this.$('[data-test-string-list-input="0"]').val(' foo ').keyup();
  assert.equal(this.$('[data-test-string-list-input="0"]').val(), 'foo');
});

test('it calls onChange with array when editing', function(assert) {
  assert.expect(1);
  this.set('inputValue', ['foo']);
  this.set('onChange', function(val) {
    assert.deepEqual(val, ['foo', 'bar'], 'calls onChange with expected value');
  });
  this.render(hbs`{{string-list inputValue=inputValue onChange=(action onChange)}}`);
  this.$('[data-test-string-list-input="1"]').val('bar').keyup();
});

test('it calls onChange with string when editing', function(assert) {
  assert.expect(1);
  this.set('inputValue', 'foo');
  this.set('onChange', function(val) {
    assert.equal(val, 'foo,bar', 'calls onChange with expected value');
  });
  this.render(hbs`{{string-list inputValue=inputValue onChange=(action onChange)}}`);
  this.$('[data-test-string-list-input="1"]').val('bar').keyup();
});

test('it removes a row', function(assert) {
  this.set('inputValue', ['foo', 'bar']);
  this.set('onChange', function(val) {
    assert.equal(val, 'bar', 'calls onChange with expected value');
  });
  this.render(hbs`{{string-list inputValue=inputValue onChange=(action onChange)}}`);

  this.$('[data-test-string-list-row="0"] [data-test-string-list-button="delete"]').click();
  assert.equal(this.$('[data-test-string-list-input]').length, 2, 'renders 2 inputs');
  assert.equal(this.$('[data-test-string-list-input="0"]').val(), 'bar');
  assert.equal(this.$('[data-test-string-list-input="1"]').val(), '');
});
