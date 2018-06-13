import { moduleForModel, test } from 'ember-qunit';
import Ember from 'ember';

moduleForModel('secret-engine', 'Unit | Model | secret-engine', {
  needs: ['model:mount-options'],
});

test('modelTypeForKV is secret by default', function(assert) {
  let model;
  Ember.run(() => {
    model = this.subject();
    assert.equal(model.get('modelTypeForKV'), 'secret');
  });
});

test('modelTypeForKV is secret-v2 for kv v2', function(assert) {
  let model;
  Ember.run(() => {
    model = this.subject({
      options: { version: 2 },
      type: 'kv',
    });
    assert.equal(model.get('modelTypeForKV'), 'secret-v2');
  });
});

test('modelTypeForKV is secret-v2 for generic v2', function(assert) {
  let model;
  Ember.run(() => {
    model = this.subject({
      options: { version: 2 },
      type: 'kv',
    });
    assert.equal(model.get('modelTypeForKV'), 'secret-v2');
  });
});
