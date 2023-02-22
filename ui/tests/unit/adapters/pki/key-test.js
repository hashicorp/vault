import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/key', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.backend = 'pki-test';
    this.secretMountPath.currentPath = this.backend;
    this.data = {
      key_id: '724862ff-6438-bad0-b598-77a6c7f4e934',
      key_type: 'ec',
      key_name: 'test-key',
      key_bits: '256',
    };
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  test('it should make request to correct endpoint on query', async function (assert) {
    assert.expect(1);
    const { key_id, ...otherAttrs } = this.data; // excludes key_id from key_info data
    const key_info = { [key_id]: { ...otherAttrs } };
    this.server.get(`${this.backend}/keys`, (schema, req) => {
      assert.strictEqual(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
      return { data: { keys: [key_id], key_info } };
    });

    this.store.query('pki/key', { backend: this.backend });
  });

  test('it should make request to correct endpoint on queryRecord', async function (assert) {
    assert.expect(1);

    this.server.get(`${this.backend}/key/${this.data.key_id}`, () => {
      assert.ok(true, 'request is made to correct endpoint on query record');
      return { data: this.data };
    });

    this.store.queryRecord('pki/key', { backend: this.backend, id: this.data.key_id });
  });

  test('it should make request to correct endpoint on delete', async function (assert) {
    assert.expect(1);
    this.store.pushPayload('pki/key', { modelName: 'pki/key', ...this.data });
    this.server.get(`${this.backend}/key/${this.data.key_id}`, () => ({ data: this.data }));
    this.server.delete(`${this.backend}/key/${this.data.key_id}`, () => {
      assert.ok(true, 'request made to correct endpoint on delete');
    });

    const model = await this.store.queryRecord('pki/key', { backend: this.backend, id: this.data.key_id });
    await model.destroyRecord();
  });
});
