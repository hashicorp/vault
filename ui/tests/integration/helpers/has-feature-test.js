import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import Ember from 'ember';

const versionStub = Ember.Service.extend({
  features: null,
});

moduleForComponent('has-feature', 'helper:has-feature', {
  integration: true,
  beforeEach: function() {
    this.register('service:version', versionStub);
    this.inject.service('version', { as: 'versionService' });
  },
});

test('it asserts on unknown features', function(assert) {
  assert.expectAssertion(() => {
    this.render(hbs`{{has-feature 'New Feature'}}`);
  }, 'New Feature is not one of the available values for Vault Enterprise features.');
});

test('it is true with existing features', function(assert) {
  this.set('versionService.features', ['HSM']);
  this.render(hbs`{{if (has-feature 'HSM') 'It works' null}}`);
  assert.dom(this._element).hasText('It works', 'present features evaluate to true');
});

test('it is false with missing features', function(assert) {
  this.set('versionService.features', ['MFA']);
  this.render(hbs`{{if (has-feature 'HSM') 'It works' null}}`);
  assert.dom(this._element).hasText('', 'missing features evaluate to false');
});
