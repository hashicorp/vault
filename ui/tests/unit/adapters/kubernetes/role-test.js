import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | kubernetes/role', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.unloadAll('kubernetes/role');
  });

  test('it should make request to correct endpoint when listing records', async function (assert) {
    assert.expect(1);
    this.server.get('/kubernetes-test/roles', (schema, req) => {
      assert.ok(req.queryParams.list, 'GET request made to correct endpoint when listing records');
      return { data: { keys: ['test-role'] } };
    });
    await this.store.query('kubernetes/role', { backend: 'kubernetes-test' });
  });

  test('it should make request to correct endpoint when querying record', async function (assert) {
    assert.expect(1);
    this.server.get('/kubernetes-test/roles/test-role', () => {
      assert.ok('GET request made to correct endpoint when querying record');
      return { data: {} };
    });
    await this.store.queryRecord('kubernetes/role', { backend: 'kubernetes-test', name: 'test-role' });
  });

  test('it should make request to correct endpoint when creating new record', async function (assert) {
    assert.expect(1);
    this.server.post('/kubernetes-test/roles/test-role', () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });
    const record = this.store.createRecord('kubernetes/role', {
      backend: 'kubernetes-test',
      name: 'test-role',
    });
    await record.save();
  });

  test('it should make request to correct endpoint when updating record', async function (assert) {
    assert.expect(1);
    this.server.post('/kubernetes-test/roles/test-role', () => {
      assert.ok('POST request made to correct endpoint when updating record');
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      name: 'test-role',
    });
    const record = this.store.peekRecord('kubernetes/role', 'test-role');
    await record.save();
  });

  test('it should make request to correct endpoint when deleting record', async function (assert) {
    assert.expect(1);
    this.server.delete('/kubernetes-test/roles/test-role', () => {
      assert.ok('DELETE request made to correct endpoint when deleting record');
    });
    this.store.pushPayload('kubernetes/role', {
      modelName: 'kubernetes/role',
      backend: 'kubernetes-test',
      name: 'test-role',
    });
    const record = this.store.peekRecord('kubernetes/role', 'test-role');
    await record.destroyRecord();
  });
});
