import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('b64-toggle', 'Integration | Component | b64 toggle', {
  integration: true,
});

test('it renders', function(assert) {
  this.render(hbs`{{b64-toggle}}`);
  assert.equal(this.$('button').length, 1);
});

test('it toggles encoding on the passed string', function(assert) {
  this.set('value', 'value');
  this.render(hbs`{{b64-toggle value=value}}`);
  this.$('button').click();
  assert.equal(this.get('value'), btoa('value'), 'encodes to base64');
  this.$('button').click();
  assert.equal(this.get('value'), 'value', 'decodes from base64');
});

test('it toggles encoding starting with base64', function(assert) {
  this.set('value', btoa('value'));
  this.render(hbs`{{b64-toggle value=value initialEncoding='base64'}}`);
  assert.ok(this.$('button').text().includes('Decode'), 'renders as on when in b64 mode');
  this.$('button').click();
  assert.equal(this.get('value'), 'value', 'decodes from base64');
});

test('it detects changes to value after encoding', function(assert) {
  this.set('value', btoa('value'));
  this.render(hbs`{{b64-toggle value=value initialEncoding='base64'}}`);
  assert.ok(this.$('button').text().includes('Decode'), 'renders as on when in b64 mode');
  this.set('value', btoa('value') + '=');
  assert.ok(this.$('button').text().includes('Encode'), 'toggles off since value has changed');
  this.set('value', btoa('value'));
  assert.ok(this.$('button').text().includes('Decode'), 'toggles on since value is equal to the original');
});

test('it does not toggle when the value is empty', function(assert) {
  this.set('value', '');
  this.render(hbs`{{b64-toggle value=value}}`);
  this.$('button').click();
  assert.ok(this.$('button').text().includes('Encode'));
});
