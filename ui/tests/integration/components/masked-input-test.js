import { moduleForComponent, test } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import maskedInput from 'vault/tests/pages/components/masked-input';

const component = create(maskedInput);

moduleForComponent('masked input', 'Integration | Component | masked input', {
  integration: true,

  beforeEach() {
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

test('it renders', function(assert) {

  this.render(hbs`{{masked-input}}`);

  assert.equal(component.wrapperClass, 'masked');

});

test('it unmasks text on focus', function(assert) {

  this.render(hbs`{{masked-input}}`);

  assert.equal(component.wrapperClass, 'masked');

  component.focus();

  assert.notEqual(component.wrapperClass, 'masked');

});

test('it remasks text on blur', function(assert) {

  this.render(hbs`{{masked-input}}`);

  assert.equal(component.wrapperClass, 'masked');

  component.focus();
  component.blur();

  assert.equal(component.wrapperClass, 'masked');

});

test('it unmasks text when button is clicked', function(assert) {

  this.render(hbs`{{masked-input}}`);

  assert.equal(component.wrapperClass, 'masked');

  component.toggleMasked();

  assert.notEqual(component.wrapperClass, 'masked');

});

test('it remasks text when button is clicked', function(assert) {

  this.render(hbs`{{masked-input}}`);

  component.toggleMasked();
  component.toggleMasked();

  assert.equal(component.wrapperClass, 'masked');

});

