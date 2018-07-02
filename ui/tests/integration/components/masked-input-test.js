import { moduleForComponent, test } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import maskedInput from 'vault/tests/pages/components/masked-input';

const component = create(maskedInput);

moduleForComponent('masked-input', 'Integration | Component | masked input', {
  integration: true,

  beforeEach() {
    component.setContext(this);
  },

  afterEach() {
    component.removeContext();
  },
});

const hasClass = (classString = '', classToFind) => {
  return classString.split(' ').contains(classToFind);
}

test('it renders', function(assert) {

  this.render(hbs`{{masked-input}}`);

  assert.ok(hasClass(component.wrapperClass, 'masked'));
});


test('it renders a textarea', function(assert) {

  this.render(hbs`{{masked-input}}`);

  assert.ok(component.textareaIsPresent);
});

test('it does not render a textarea when displayOnly is true', function(assert) {

  this.render(hbs`{{masked-input displayOnly=true}}`);

  assert.notOk(component.textareaIsPresent);
});


test('it unmasks text on focus', function(assert) {

  this.set('value', 'value');
  this.render(hbs`{{masked-input value=value}}`);

  assert.ok(hasClass(component.wrapperClass, 'masked'));

  component.focus();
  assert.notOk(hasClass(component.wrapperClass, 'masked'));
});

test('it remasks text on blur', function(assert) {

  this.set('value', 'value');
  this.render(hbs`{{masked-input value=value}}`);

  assert.ok(hasClass(component.wrapperClass, 'masked'));

  component.focus();
  component.blur();

  assert.ok(hasClass(component.wrapperClass, 'masked'));
});

test('it unmasks text when button is clicked', function(assert) {

  this.set('value', 'value');
  this.render(hbs`{{masked-input value=value}}`);

  assert.ok(hasClass(component.wrapperClass, 'masked'));

  component.toggleMasked();

  assert.notOk(hasClass(component.wrapperClass, 'masked'));

});

test('it remasks text when button is clicked', function(assert) {

  this.set('value', 'value');
  this.render(hbs`{{masked-input value=value}}`);

  component.toggleMasked();
  component.toggleMasked();

  assert.ok(hasClass(component.wrapperClass, 'masked'));

});

