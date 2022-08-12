import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/client', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/client';
    this.data = {
      name: 'client-1',
      key: 'test-key',
      access_token_ttl: '30m',
      id_token_ttl: '1h',
    };
    this.path = '/identity/oidc/client/client-1';
  });

  testHelper(test);

  test('it filters list response when passed query containing clientIds', async function (assert) {
    assert.expect(2);

    const keys = ['client-1', 'client-2', 'client-3'];
    const key_info = [
      {
        'client-1': {
          key: 'test-key',
          access_token_ttl: '30m',
          id_token_ttl: '1h',
        },
      },
      {
        'client-2': {
          key: 'test-key',
          access_token_ttl: '30m',
          id_token_ttl: '1h',
        },
      },
      {
        'client-3': {
          key: 'test-key',
          access_token_ttl: '30m',
          id_token_ttl: '1h',
        },
      },
    ];

    this.server.get(`/identity/${this.modelName}`, (schema, req) => {
      assert.equal(req.queryParams.list, 'true', 'request is made to correct endpoint on query');
      return { data: { keys, key_info } };
    });

    let testQuery = ['*', 'client-1'];
    await this.store
      .query(this.modelName, { filterIds: testQuery })
      .then((resp) => assert.equal(resp.content.length, 3, 'returns all clients when ids include glob (*)'));

    testQuery = ['*'];
    await this.store
      .query(this.modelName, { filterIds: testQuery })
      .then((resp) => assert.equal(resp.content.length, 3, 'returns all clients when glob (*) is only id'));

    testQuery = ['client-2'];
    await this.store.query(this.modelName, { filterIds: testQuery }).then((resp) => {
      assert.equal(resp.content.length, 1, 'filters response and returns only matching id');
      assert.equal(resp.firstObject.name, 'client-2', 'response contains correct model');
    });

    testQuery = ['client-2', 'client-3'];
    await this.store.query(this.modelName, { filterIds: testQuery }).then((resp) => {
      assert.equal(resp.content.length, 2, 'filters response when passed multiple ids');
    });
  });
});
