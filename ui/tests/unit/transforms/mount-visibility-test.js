import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

module('Unit | Transform | mount visibility', function (hooks) {
  setupTest(hooks);

  test('it serializes correctly for API', function (assert) {
    const transform = this.owner.lookup('transform:mount-visibility');
    assert.ok(transform);
    let serialized = transform.serialize(true);
    assert.strictEqual(serialized, 'unauth');
    serialized = transform.serialize(false);
    assert.strictEqual(serialized, 'hidden');
  });

  test('it deserializes correctly from API', function (assert) {
    const transform = this.owner.lookup('transform:mount-visibility');
    let deserialized = transform.deserialize('unauth');
    assert.true(deserialized, 'deserializes "unauth" string value to true');
    deserialized = transform.deserialize('hidden');
    assert.false(deserialized, 'deserializes "hidden" string value to false');
    deserialized = transform.deserialize('');
    assert.false(deserialized, 'deserializes empty string to false');
  });
});
