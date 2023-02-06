import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | pki/role', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'pki-test';
    this.store.unloadAll('pki/role');
  });

  test('it should make request to correct endpoint when updating a record', async function (assert) {
    assert.expect(1);
    this.server.post('/pki-test/role/pki-role', () => {});

    this.store.pushPayload('pki/role', {
      modelName: 'pki/role',
      backend: 'pki-not-hardcoded',
      id: 'pki-test',
    });
    const record = this.store.peekRecord('pki/role', 'pki-test');
    const response = await record.save().catch((e) => {
      return e.path;
    });
    // because we're inside an engine we can't force permission capabilities to POST.
    // To work around this we look at the 403 error and confirm the path posted includes the custom backend name (pki-test and not hardcoded pki).
    assert.strictEqual(response, '/v1/pki-not-hardcoded/roles');
  });

  test('meep it should make request to correct endpoint on query', async function (assert) {
    // THE store.quey method does not seem to hit a permissions issue.
    assert.expect(1);
    this.server.get(`${this.backend}/roles`, (schema, req) => {
      assert.strictEqual(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
      return {};
    });

    this.store.query('pki/role', { backend: this.backend });
  });
});
