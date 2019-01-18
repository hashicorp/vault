import { run } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Model | secret-engine', function(hooks) {
  setupTest(hooks);

  test('modelTypeForKV is secret by default', function(assert) {
    let model;
    run(() => {
      model = run(() => this.owner.lookup('service:store').createRecord('secret-engine'));
      assert.equal(model.get('modelTypeForKV'), 'secret');
    });
  });

  test('modelTypeForKV is secret-v2 for kv v2', function(assert) {
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          options: { version: 2 },
          type: 'kv',
        })
      );
      assert.equal(model.get('modelTypeForKV'), 'secret-v2');
    });
  });

  test('modelTypeForKV is secret-v2 for generic v2', function(assert) {
    let model;
    run(() => {
      model = run(() =>
        this.owner.lookup('service:store').createRecord('secret-engine', {
          options: { version: 2 },
          type: 'kv',
        })
      );
      assert.equal(model.get('modelTypeForKV'), 'secret-v2');
    });
  });
});
