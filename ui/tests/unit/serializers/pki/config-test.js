import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Serializer | pki/config', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const store = this.owner.lookup('service:store');
    const serializer = store.serializerFor('pki/config');

    assert.ok(serializer);
  });

  test('it serializes records', function (assert) {
    const store = this.owner.lookup('service:store');
    const record = store.createRecord('pki/config', {});

    const serializedRecord = record.serialize();

    assert.ok(serializedRecord);
  });
});
