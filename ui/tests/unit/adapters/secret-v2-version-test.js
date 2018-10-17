import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

module('Unit | Adapter | secret-v2-version', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  let fakeStore = {
    peekRecord() {
      return {
        reload() {},
      };
    },
  };
  [
    [
      'findRecord with version',
      'findRecord',
      [null, {}, JSON.stringify(['secret', 'foo', '2']), {}],
      'GET',
      '/v1/secret/data/foo?version=2',
    ],
    [
      'v2DeleteOperation with delete',
      'v2DeleteOperation',
      [fakeStore, JSON.stringify(['secret', 'foo', '2']), 'delete'],
      'POST',
      '/v1/secret/delete/foo',
      { versions: ['2'] },
    ],
    [
      'v2DeleteOperation with destroy',
      'v2DeleteOperation',
      [fakeStore, JSON.stringify(['secret', 'foo', '2']), 'destroy'],
      'POST',
      '/v1/secret/destroy/foo',
      { versions: ['2'] },
    ],
    [
      'v2DeleteOperation with destroy',
      'v2DeleteOperation',
      [fakeStore, JSON.stringify(['secret', 'foo', '2']), 'undelete'],
      'POST',
      '/v1/secret/undelete/foo',
      { versions: ['2'] },
    ],
    [
      'updateRecord makes calls to correct url',
      'updateRecord',
      [
        {
          serializerFor() {
            return { serializeIntoHash() {} };
          },
        },
        {},
        { id: JSON.stringify(['secret', 'foo', '2']) },
      ],
      'PUT',
      '/v1/secret/data/foo',
    ],
  ].forEach(([testName, adapterMethod, args, expectedHttpVerb, expectedURL, exptectedRequestBody]) => {
    test(`${testName}`, function(assert) {
      let adapter = this.owner.lookup('adapter:secret-v2-version');
      adapter[adapterMethod](...args);
      let { url, method, requestBody } = this.server.handledRequests[0];
      assert.equal(url, expectedURL, `${adapterMethod} calls the correct url: ${expectedURL}`);
      assert.equal(
        method,
        expectedHttpVerb,
        `${adapterMethod} uses the correct http verb: ${expectedHttpVerb}`
      );
      if (exptectedRequestBody) {
        assert.deepEqual(JSON.parse(requestBody), exptectedRequestBody);
      }
    });
  });
});
