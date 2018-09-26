import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';

module('Unit | Adapter | secret engine', function(hooks) {
  setupTest(hooks);

  hooks.beforeEach(function() {
    this.server = apiStub();
  });

  hooks.afterEach(function() {
    this.server.shutdown();
  });

  const storeStub = {
    serializerFor() {
      return {
        serializeIntoHash() {},
      };
    },
  };
  const type = {
    modelName: 'secret-engine',
  };

  const cases = [
    {
      description: 'Empty query',
      adapterMethod: 'query',
      args: [storeStub, type, {}],
      url: '/v1/sys/internal/ui/mounts',
      method: 'GET',
    },
    {
      description: 'Query with a path',
      adapterMethod: 'query',
      args: [storeStub, type, { path: 'foo' }],
      url: '/v1/sys/internal/ui/mounts/foo',
      method: 'GET',
    },

    {
      description: 'Query with nested path',
      adapterMethod: 'query',
      args: [storeStub, type, { path: 'foo/bar/baz' }],
      url: '/v1/sys/internal/ui/mounts/foo/bar/baz',
      method: 'GET',
    },
  ];
  cases.forEach(testCase => {
    test(`secret-engine: ${testCase.description}`, function(assert) {
      assert.expect(2);
      let adapter = this.owner.lookup('adapter:secret-engine');
      adapter[testCase.adapterMethod](...testCase.args);
      let { url, method } = this.server.handledRequests[0];
      assert.equal(url, testCase.url, `${testCase.adapterMethod} calls the correct url: ${testCase.url}`);
      assert.equal(
        method,
        testCase.method,
        `${testCase.adapterMethod} uses the correct http verb: ${testCase.method}`
      );
    });
  });
});
