import { moduleForModel, test } from 'ember-qunit';

moduleForModel('auth-method', 'Unit | Model | auth-method', {
  needs: ['serializer:mount-config', 'model:auth-config', 'model:mount-config'],
});

test('it exists', function(assert) {
  let model = this.subject();
  assert.ok(!!model);
});

test('serialize', function(assert) {
  let model = this.subject();
  Ember.run(() => {
    let config = model.get('config');
    let data = config.serialize();
    debugger;
  });
});
