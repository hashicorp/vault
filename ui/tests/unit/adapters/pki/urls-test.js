import { module, test } from 'qunit';
import { resolve } from 'rsvp';
import { setupTest } from 'vault/tests/helpers';

const storeStub = {
  pushPayload() {},
  serializerFor() {
    return {
      serializeIntoHash() {},
    };
  },
};
const makeSnapshot = (obj) => {
  const snapshot = {
    id: obj.id,
    record: {
      ...obj,
    },
  };
  snapshot.attr = (attr) => snapshot[attr];
  return snapshot;
};

module('Unit | Adapter | pki/urls', function (hooks) {
  setupTest(hooks);

  test('pki url endpoints', function (assert) {
    let url, method;
    const adapter = this.owner.factoryFor('adapter:pki/urls').create({
      ajax: (...args) => {
        [url, method] = args;
        return resolve({});
      },
    });

    adapter.createRecord(storeStub, 'pki/urls', makeSnapshot({ id: 'pki-create' }));
    assert.strictEqual(url, '/v1/pki-create/config/urls', 'create url OK');
    assert.strictEqual(method, 'POST', 'create method OK');

    adapter.updateRecord(storeStub, 'pki/urls', makeSnapshot({ id: 'pki-update' }));
    assert.strictEqual(url, '/v1/pki-update/config/urls', 'update url OK');
    assert.strictEqual(method, 'PUT', 'update method OK');

    adapter.findRecord(null, 'capabilities', 'pki-find');
    assert.strictEqual(url, '/v1/pki-find/config/urls', 'find url OK');
    assert.strictEqual(method, 'GET', 'find method OK');
  });
});
