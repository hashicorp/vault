import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { validate } from 'uuid';

module('Unit | Serializer | cluster', function (hooks) {
  setupTest(hooks);

  test('it should generate ids for replication attributes', async function (assert) {
    const serializer = this.owner.lookup('serializer:cluster');
    const data = {};
    serializer.setReplicationId(data);
    assert.true(validate(data.id), 'UUID is generated for replication attribute');
  });
});
