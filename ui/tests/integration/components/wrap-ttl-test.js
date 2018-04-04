import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('wrap-ttl', 'Integration | Component | wrap ttl', {
  integration: true,
  beforeEach() {
    this.lastOnChangeCall = null;
    this.set('onChange', val => {
      this.lastOnChangeCall = val;
    });
  },
});

test('it requires `onChange`', function(assert) {
  assert.expectAssertion(
    () => this.render(hbs`{{wrap-ttl}}`),
    /`onChange` handler is a required attr in/,
    'asserts without onChange'
  );
});

test('it renders', function(assert) {
  this.render(hbs`{{wrap-ttl onChange=(action onChange)}}`);
  assert.equal(this.lastOnChangeCall, '30m', 'calls onChange with 30m default on first render');
  assert.equal(this.$('label[for=wrap-response]').text().trim(), 'Wrap response');
});

test('it nulls out value when you uncheck wrapResponse', function(assert) {
  this.render(hbs`{{wrap-ttl onChange=(action onChange)}}`);
  this.$('#wrap-response').click().change();
  assert.equal(this.lastOnChangeCall, null, 'calls onChange with null');
});

test('it sends value changes to onChange handler', function(assert) {
  this.render(hbs`{{wrap-ttl onChange=(action onChange)}}`);

  this.$('[data-test-wrap-ttl-picker] input').val('20').trigger('input');
  assert.equal(this.lastOnChangeCall, '20m', 'calls onChange correctly on time input');

  this.$('#unit').val('h').change();
  assert.equal(this.lastOnChangeCall, '20h', 'calls onChange correctly on unit change');

  this.$('#unit').val('d').change();
  assert.equal(this.lastOnChangeCall, '480h', 'converts days to hours correctly');
});
